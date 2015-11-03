package listeners_test

import (
	"crypto/tls"
	"encoding/binary"
	"strings"

	"github.com/cloudfoundry/sonde-go/events"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"doppler/config"
	"doppler/listeners"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/cloudfoundry/dropsonde/emitter"
	"github.com/cloudfoundry/dropsonde/factories"
	"github.com/cloudfoundry/loggregatorlib/loggertesthelper"
	"github.com/gogo/protobuf/proto"
	"github.com/nu7hatch/gouuid"
)

var _ = Describe("TLSlistener", func() {
	var envelopeChan chan *events.Envelope
	var tlsListener listeners.Listener
	var tlsListenerConfig config.TLSListenerConfig
	var tlsClientConfig *tls.Config

	BeforeEach(func() {
		tlsListenerConfig = config.TLSListenerConfig{
			CertFile: "fixtures/server.crt",
			KeyFile:  "fixtures/server.key",
			CAFile:   "fixtures/loggregator-ca.crt",
		}

		var err error
		tlsClientConfig, err = listeners.NewTLSConfig("fixtures/client.crt", "fixtures/client.key", "fixtures/loggregator-ca.crt")
		Expect(err).NotTo(HaveOccurred())
		tlsClientConfig.ServerName = "doppler"

		envelopeChan = make(chan *events.Envelope)
	})

	JustBeforeEach(func() {
		var err error
		tlsListener, err = listeners.NewTLSListener("aname", "127.0.0.1:0", tlsListenerConfig, envelopeChan, loggertesthelper.Logger())
		Expect(err).NotTo(HaveOccurred())
		go tlsListener.Start()
	})

	AfterEach(func() {
		tlsListener.Stop()
	})

	Context("With invalid client configuration", func() {
		JustBeforeEach(func() {
			conn := openTLSConnection(tlsListener.Address(), tlsClientConfig)
			conn.Close()
		})

		Context("without a CA file", func() {
			It("fails", func() {
				tlsClientConfig, err := listeners.NewTLSConfig("fixtures/client.crt", "fixtures/client.key", "")
				Expect(err).NotTo(HaveOccurred())
				tlsClientConfig.ServerName = "doppler"

				_, err = tls.Dial("tcp", tlsListener.Address(), tlsClientConfig)
				Expect(err).To(MatchError("x509: certificate signed by unknown authority"))
			})
		})

		Context("without a server name", func() {
			It("fails", func() {
				tlsClientConfig.ServerName = ""
				_, err := tls.Dial("tcp", tlsListener.Address(), tlsClientConfig)
				Expect(err).To(MatchError("x509: cannot validate certificate for 127.0.0.1 because it doesn't contain any IP SANs"))
			})
		})
	})

	Context("dropsonde metric emission", func() {
		BeforeEach(func() {
			fakeEventEmitter.Reset()
			metricBatcher.Reset()
		})

		It("sends all types of messages as a protobuf", func() {
			for name, eventType := range events.Envelope_EventType_value {
				envelope := createEnvelope(events.Envelope_EventType(eventType))
				conn := openTLSConnection(tlsListener.Address(), tlsClientConfig)

				err := send(conn, envelope)
				Expect(err).ToNot(HaveOccurred())

				Eventually(envelopeChan).Should(Receive(Equal(envelope)), fmt.Sprintf("did not receive expected event: %s", name))
				conn.Close()
			}
		})

		It("sends all types of messages over multiple connections", func() {
			for _, eventType := range events.Envelope_EventType_value {
				envelope1 := createEnvelope(events.Envelope_EventType(eventType))
				conn1 := openTLSConnection(tlsListener.Address(), tlsClientConfig)

				envelope2 := createEnvelope(events.Envelope_EventType(eventType))
				conn2 := openTLSConnection(tlsListener.Address(), tlsClientConfig)

				err := send(conn1, envelope1)
				Expect(err).ToNot(HaveOccurred())
				err = send(conn2, envelope2)
				Expect(err).ToNot(HaveOccurred())

				envelopes := readMessages(envelopeChan, 2)
				Expect(envelopes).To(ContainElement(envelope1))
				Expect(envelopes).To(ContainElement(envelope2))

				conn1.Close()
				conn2.Close()
			}
		})

		It("issues intended metrics", func() {
			envelope := createEnvelope(events.Envelope_LogMessage)
			conn := openTLSConnection(tlsListener.Address(), tlsClientConfig)

			err := send(conn, envelope)
			Expect(err).ToNot(HaveOccurred())
			conn.Close()

			Eventually(envelopeChan).Should(Receive())

			Eventually(func() int {
				return len(fakeEventEmitter.GetMessages())
			}).Should(BeNumerically(">", 2))

			var counterEvents []*events.CounterEvent
			for _, e := range fakeEventEmitter.GetMessages() {
				if ce, ok := e.Event.(*events.CounterEvent); ok {
					if strings.HasPrefix(ce.GetName(), "aname.") {
						counterEvents = append(counterEvents, ce)
					}
				}
			}

			Expect(counterEvents).To(ConsistOf(
				&events.CounterEvent{
					Name:  proto.String("aname.receivedMessageCount"),
					Delta: proto.Uint64(1),
				},
				&events.CounterEvent{
					Name:  proto.String("aname.receivedByteCount"),
					Delta: proto.Uint64(67),
				},
			))
		})

		Context("receiveErrors count", func() {
			expectReceiveErrorCount := func() {
				Eventually(func() int {
					return len(fakeEventEmitter.GetMessages())
				}).Should(Equal(1))
				errorEvents := fakeEventEmitter.GetEvents()
				Expect(errorEvents).To(ConsistOf(
					&events.CounterEvent{
						Name:  proto.String("aname.receiveErrorCount"),
						Delta: proto.Uint64(1),
					},
				))
			}

			It("does not increment error count for a valid message", func() {
				envelope := createEnvelope(events.Envelope_LogMessage)
				conn := openTLSConnection(tlsListener.Address(), tlsClientConfig)

				err := send(conn, envelope)
				Expect(err).ToNot(HaveOccurred())
			})

			It("increments when tls handshake fails", func() {
				tlsConfig, err := listeners.NewTLSConfig("fixtures/bad_client.crt", "fixtures/bad_client.key", "fixtures/badCA.crt")
				Expect(err).NotTo(HaveOccurred())

				_, err = tls.Dial("tcp", tlsListener.Address(), tlsConfig)
				Expect(err).To(HaveOccurred())
				expectReceiveErrorCount()
			})

			It("increments when size is greater than payload", func() {
				conn := openTLSConnection(tlsListener.Address(), tlsClientConfig)
				bytes := []byte("invalid payload")
				var n uint32
				n = 1000
				err := binary.Write(conn, binary.LittleEndian, n)
				Expect(err).ToNot(HaveOccurred())
				_, err = conn.Write(bytes)
				Expect(err).ToNot(HaveOccurred())
				conn.Close()
				expectReceiveErrorCount()
			})
		})
	})

	Context("Start Stop", func() {
		It("panics if you start again", func() {
			conn := openTLSConnection(tlsListener.Address(), tlsClientConfig)
			defer conn.Close()

			Expect(tlsListener.Start).Should(Panic())
		})

		It("panics if you start after a stop", func() {
			conn := openTLSConnection(tlsListener.Address(), tlsClientConfig)
			defer conn.Close()

			tlsListener.Stop()
			Expect(tlsListener.Start).Should(Panic())
		})

		It("fails to send message after listener has been stopped", func() {
			logMessage := factories.NewLogMessage(events.LogMessage_OUT, "some message", "appId", "source")
			envelope, _ := emitter.Wrap(logMessage, "origin")
			conn := openTLSConnection(tlsListener.Address(), tlsClientConfig)

			err := send(conn, envelope)
			Expect(err).ToNot(HaveOccurred())

			tlsListener.Stop()

			Eventually(func() error {
				return send(conn, envelope)
			}).Should(HaveOccurred())

			conn.Close()
		})
	})
})

func readMessages(envelopeChan chan *events.Envelope, n int) []*events.Envelope {
	var envelopes []*events.Envelope
	for i := 0; i < n; i++ {
		var envelope *events.Envelope
		Eventually(envelopeChan).Should(Receive(&envelope))
		envelopes = append(envelopes, envelope)
	}
	return envelopes
}

func openTLSConnection(address string, tlsConfig *tls.Config) net.Conn {
	var conn net.Conn
	var err error
	Eventually(func() error {
		conn, err = tls.Dial("tcp", address, tlsConfig)
		return err
	}).ShouldNot(HaveOccurred())

	return conn
}

func send(conn net.Conn, envelope *events.Envelope) error {
	bytes, err := proto.Marshal(envelope)
	if err != nil {
		return err
	}

	var n uint32
	n = uint32(len(bytes))
	err = binary.Write(conn, binary.LittleEndian, n)
	if err != nil {
		return err
	}

	_, err = conn.Write(bytes)
	return err
}

func createEnvelope(eventType events.Envelope_EventType) *events.Envelope {

	envelope := &events.Envelope{Origin: proto.String("origin"), EventType: &eventType, Timestamp: proto.Int64(time.Now().UnixNano())}

	switch eventType {
	case events.Envelope_HttpStart:
		req, _ := http.NewRequest("GET", "http://www.example.com", nil)
		req.RemoteAddr = "www.example.com"
		req.Header.Add("User-Agent", "user-agent")
		uuid, _ := uuid.NewV4()
		envelope.HttpStart = factories.NewHttpStart(req, events.PeerType_Client, uuid)
	case events.Envelope_HttpStop:
		req, _ := http.NewRequest("GET", "http://www.example.com", nil)
		uuid, _ := uuid.NewV4()
		envelope.HttpStop = factories.NewHttpStop(req, http.StatusOK, 128, events.PeerType_Client, uuid)
	case events.Envelope_HttpStartStop:
		req, _ := http.NewRequest("GET", "http://www.example.com", nil)
		req.RemoteAddr = "www.example.com"
		req.Header.Add("User-Agent", "user-agent")
		uuid, _ := uuid.NewV4()
		envelope.HttpStartStop = factories.NewHttpStartStop(req, http.StatusOK, 128, events.PeerType_Client, uuid)
	case events.Envelope_ValueMetric:
		envelope.ValueMetric = factories.NewValueMetric("some-value-metric", 123, "km")
	case events.Envelope_CounterEvent:
		envelope.CounterEvent = factories.NewCounterEvent("some-counter-event", 123)
	case events.Envelope_LogMessage:
		envelope.LogMessage = factories.NewLogMessage(events.LogMessage_OUT, "some message", "appId", "source")
	case events.Envelope_ContainerMetric:
		envelope.ContainerMetric = factories.NewContainerMetric("appID", 123, 1, 5, 5)
	case events.Envelope_Error:
		envelope.Error = factories.NewError("source", 123, "message")
	default:
		fmt.Printf("Unknown event %v\n", eventType)
		return nil
	}

	return envelope
}
