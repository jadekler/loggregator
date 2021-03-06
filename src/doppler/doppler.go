package main

import (
	"diodes"
	"fmt"
	"log"
	"sync"
	"time"

	"doppler/config"
	"doppler/grpcmanager"
	"doppler/sinkserver"
	"doppler/sinkserver/blacklist"
	"doppler/sinkserver/sinkmanager"
	"doppler/sinkserver/websocketserver"

	"doppler/listeners"
	"monitor"

	"github.com/cloudfoundry/dropsonde"
	"github.com/cloudfoundry/dropsonde/dropsonde_unmarshaller"
	"github.com/cloudfoundry/dropsonde/metric_sender"
	"github.com/cloudfoundry/dropsonde/metricbatcher"
	"github.com/cloudfoundry/dropsonde/metrics"
	"github.com/cloudfoundry/dropsonde/signature"
	"github.com/cloudfoundry/loggregatorlib/appservice"
	"github.com/cloudfoundry/loggregatorlib/store"
	"github.com/cloudfoundry/loggregatorlib/store/cache"
	"github.com/cloudfoundry/sonde-go/events"
	"github.com/cloudfoundry/storeadapter"
)

type Doppler struct {
	batcher *metricbatcher.MetricBatcher

	appStoreWatcher *store.AppServiceStoreWatcher

	errChan         chan error
	envelopeBuffer  *diodes.ManyToOneEnvelope
	udpListener     *listeners.UDPListener
	tcpListener     *listeners.TCPListener
	tlsListener     *listeners.TCPListener
	grpcListener    *listeners.GRPCListener
	sinkManager     *sinkmanager.SinkManager
	messageRouter   *sinkserver.MessageRouter
	websocketServer *websocketserver.WebsocketServer

	dropsondeUnmarshallerCollection *dropsonde_unmarshaller.DropsondeUnmarshallerCollection
	dropsondeBytesChan              <-chan []byte
	dropsondeVerifiedBytesChan      chan []byte
	signatureVerifier               *signature.Verifier

	storeAdapter storeadapter.StoreAdapter

	uptimeMonitor   *monitor.Uptime
	openFileMonitor *monitor.LinuxFileDescriptor

	newAppServiceChan, deletedAppServiceChan <-chan appservice.AppService
	wg                                       sync.WaitGroup
}

func New(
	host string,
	conf *config.Config,
	storeAdapter storeadapter.StoreAdapter,
	messageDrainBufferSize uint,
	dropsondeOrigin string,
	websocketWriteTimeout time.Duration,
	dialTimeout time.Duration,
) (*Doppler, error) {
	doppler := &Doppler{
		storeAdapter:               storeAdapter,
		dropsondeVerifiedBytesChan: make(chan []byte),
	}

	keepAliveInterval := 30 * time.Second

	appStoreCache := cache.NewAppServiceCache()
	doppler.appStoreWatcher, doppler.newAppServiceChan, doppler.deletedAppServiceChan = store.NewAppServiceStoreWatcher(storeAdapter, appStoreCache)

	doppler.batcher = initializeMetrics(conf.MetricBatchIntervalMilliseconds)

	doppler.udpListener, doppler.dropsondeBytesChan = listeners.NewUDPListener(
		fmt.Sprintf("%s:%d", host, conf.IncomingUDPPort),
		doppler.batcher,
		"udpListener",
	)

	doppler.envelopeBuffer = diodes.NewManyToOneEnvelope(10000, doppler)

	var err error
	if conf.EnableTLSTransport {
		tlsConfig := &conf.TLSListenerConfig
		addr := fmt.Sprintf("%s:%d", host, tlsConfig.Port)
		contextName := "tlsListener"
		doppler.tlsListener, err = listeners.NewTCPListener(contextName, addr, tlsConfig, doppler.envelopeBuffer, doppler.batcher, TCPTimeout)
		if err != nil {
			return nil, err
		}
	}

	addr := fmt.Sprintf("%s:%d", host, conf.IncomingTCPPort)
	contextName := "tcpListener"
	doppler.tcpListener, err = listeners.NewTCPListener(contextName, addr, nil, doppler.envelopeBuffer, doppler.batcher, TCPTimeout)

	doppler.signatureVerifier = signature.NewVerifier(conf.SharedSecret)

	doppler.dropsondeUnmarshallerCollection = dropsonde_unmarshaller.NewDropsondeUnmarshallerCollection(conf.UnmarshallerCount)

	blacklist := blacklist.New(conf.BlackListIps)
	metricTTL := time.Duration(conf.ContainerMetricTTLSeconds) * time.Second
	sinkTimeout := time.Duration(conf.SinkInactivityTimeoutSeconds) * time.Second
	sinkIOTimeout := time.Duration(conf.SinkIOTimeoutSeconds) * time.Second
	doppler.sinkManager = sinkmanager.New(
		conf.MaxRetainedLogMessages,
		conf.SinkSkipCertVerify,
		blacklist,
		messageDrainBufferSize,
		dropsondeOrigin,
		sinkTimeout,
		sinkIOTimeout,
		metricTTL,
		dialTimeout,
	)

	grpcRouter := grpcmanager.NewRouter()
	doppler.grpcListener, err = listeners.NewGRPCListener(grpcRouter, doppler.sinkManager, conf.GRPC, doppler.envelopeBuffer, doppler.batcher)
	if err != nil {
		return nil, err
	}

	doppler.messageRouter = sinkserver.NewMessageRouter(doppler.sinkManager, grpcRouter)

	doppler.websocketServer, err = websocketserver.New(
		fmt.Sprintf("%s:%d", conf.WebsocketHost, conf.OutgoingPort),
		doppler.sinkManager,
		websocketWriteTimeout,
		keepAliveInterval,
		conf.MessageDrainBufferSize,
		dropsondeOrigin,
		doppler.batcher,
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to create the websocket server: %s", err.Error())
	}

	monitorInterval := time.Duration(conf.MonitorIntervalSeconds) * time.Second
	doppler.openFileMonitor = monitor.NewLinuxFD(monitorInterval)
	doppler.uptimeMonitor = monitor.NewUptime(monitorInterval)

	return doppler, nil
}

func (doppler *Doppler) Start() {
	doppler.errChan = make(chan error)

	doppler.wg.Add(7 + doppler.dropsondeUnmarshallerCollection.Size())

	go func() {
		defer doppler.wg.Done()
		doppler.grpcListener.Start()
	}()

	go func() {
		defer doppler.wg.Done()
		doppler.appStoreWatcher.Run()
	}()

	go func() {
		defer doppler.wg.Done()
		doppler.udpListener.Start()
	}()

	go func() {
		defer doppler.wg.Done()
		doppler.tcpListener.Start()
	}()

	if doppler.tlsListener != nil {
		doppler.wg.Add(1)
		go func() {
			defer doppler.wg.Done()
			doppler.tlsListener.Start()
		}()
	}

	udpEnvelopes := make(chan *events.Envelope)
	doppler.dropsondeUnmarshallerCollection.Run(doppler.dropsondeVerifiedBytesChan, udpEnvelopes, &doppler.wg)
	go func() {
		for {
			env := <-udpEnvelopes
			doppler.batcher.BatchCounter("listeners.receivedEnvelopes").
				SetTag("protocol", "udp").
				SetTag("event_type", env.GetEventType().String()).
				Increment()
			doppler.envelopeBuffer.Set(env)
		}
	}()

	go func() {
		defer func() {
			doppler.wg.Done()
			close(doppler.dropsondeVerifiedBytesChan)
		}()
		doppler.signatureVerifier.Run(doppler.dropsondeBytesChan, doppler.dropsondeVerifiedBytesChan)
	}()

	go func() {
		defer doppler.wg.Done()
		doppler.sinkManager.Start(doppler.newAppServiceChan, doppler.deletedAppServiceChan)
	}()

	go func() {
		defer doppler.wg.Done()
		doppler.messageRouter.Start(doppler.envelopeBuffer)
	}()

	go func() {
		defer doppler.wg.Done()
		doppler.websocketServer.Start()
	}()

	go doppler.uptimeMonitor.Start()
	go doppler.openFileMonitor.Start()

	// The following runs forever. Put all startup functions above here.
	for err := range doppler.errChan {
		log.Printf("Got error %s", err)
	}
}

func (doppler *Doppler) Stop() {
	go doppler.udpListener.Stop()
	go doppler.tcpListener.Stop()
	go doppler.tlsListener.Stop()
	go doppler.sinkManager.Stop()
	go doppler.websocketServer.Stop()
	doppler.appStoreWatcher.Stop()
	doppler.wg.Wait()

	doppler.storeAdapter.Disconnect()
	close(doppler.errChan)
	doppler.uptimeMonitor.Stop()
	doppler.openFileMonitor.Stop()
}

func (doppler *Doppler) Alert(missed int) {
	log.Printf("Shed %d envelopes", missed)
	doppler.batcher.BatchCounter("doppler.shedEnvelopes").Add(uint64(missed))
}

func initializeMetrics(batchIntervalMilliseconds uint) *metricbatcher.MetricBatcher {
	eventEmitter := dropsonde.AutowiredEmitter()
	metricSender := metric_sender.NewMetricSender(eventEmitter)
	metricBatcher := metricbatcher.New(metricSender, time.Duration(batchIntervalMilliseconds)*time.Millisecond)
	metrics.Initialize(metricSender, metricBatcher)
	return metricBatcher
}
