package grpcmanager_test

import (
	"context"
	"diodes"
	"doppler/grpcmanager"
	"errors"
	"io"
	"net"
	"plumbing"

	"google.golang.org/grpc"

	"github.com/apoydence/eachers/testhelpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("IngestorManager", func() {
	var startGRPCServer = func(ds plumbing.DopplerIngestorServer) (*grpc.Server, string) {
		lis, err := net.Listen("tcp", ":0")
		Expect(err).ToNot(HaveOccurred())
		s := grpc.NewServer()
		plumbing.RegisterDopplerIngestorServer(s, ds)
		go s.Serve(lis)

		return s, lis.Addr().String()
	}

	var establishClient = func(dopplerAddr string) (plumbing.DopplerIngestorClient, io.Closer) {
		conn, err := grpc.Dial(dopplerAddr, grpc.WithInsecure())
		Expect(err).ToNot(HaveOccurred())
		c := plumbing.NewDopplerIngestorClient(conn)

		return c, conn
	}

	var (
		outgoingMsgs  *diodes.ManyToOneEnvelope
		manager       *grpcmanager.IngestorManager
		server        *grpc.Server
		connCloser    io.Closer
		dopplerClient plumbing.DopplerIngestorClient
	)

	BeforeEach(func() {
		var grpcAddr string
		outgoingMsgs = diodes.NewManyToOneEnvelope(5, nil)
		mockBatcher := newMockBatcher()
		mockChainer := newMockBatchCounterChainer()
		testhelpers.AlwaysReturn(mockBatcher.BatchCounterOutput, mockChainer)
		testhelpers.AlwaysReturn(mockChainer.SetTagOutput, mockChainer)
		manager = grpcmanager.NewIngestor(outgoingMsgs, mockBatcher)
		server, grpcAddr = startGRPCServer(manager)
		dopplerClient, connCloser = establishClient(grpcAddr)
	})

	AfterEach(func() {
		server.Stop()
		connCloser.Close()
	})

	It("reads envelopes from ingestor client", func() {
		pusherClient, err := dopplerClient.Pusher(context.TODO())
		Expect(err).ToNot(HaveOccurred())

		someEnvelope, data := buildContainerMetric()
		pusherClient.Send(&plumbing.EnvelopeData{data})

		Eventually(outgoingMsgs.Next).Should(Equal(someEnvelope))
	})

	Context("With an unsupported envelope payload", func() {
		It("does not forward the message to the sender", func() {
			pusherClient, err := dopplerClient.Pusher(context.TODO())
			Expect(err).ToNot(HaveOccurred())

			err = pusherClient.Send(&plumbing.EnvelopeData{[]byte("unsupported envelope")})
			Expect(err).ToNot(HaveOccurred())
			Consistently(func() bool {
				_, ok := outgoingMsgs.TryNext()
				return ok
			}).Should(BeFalse())

			err = pusherClient.Send(&plumbing.EnvelopeData{nil})
			Expect(err).ToNot(HaveOccurred())
			Consistently(func() bool {
				_, ok := outgoingMsgs.TryNext()
				return ok
			}).Should(BeFalse())
		})
	})

	Context("When the Recv returns an EOF error", func() {
		It("exits the function gracefully", func() {
			fakeStream := newMockIngestorGRPCServer()
			fakeStream.RecvOutput.Ret0 <- nil
			fakeStream.RecvOutput.Ret1 <- io.EOF
			fakeStream.ContextOutput.Ret0 <- context.TODO()

			Eventually(func() error {
				return manager.Pusher(fakeStream)
			}).Should(Succeed())
			Consistently(func() bool {
				_, ok := outgoingMsgs.TryNext()
				return ok
			}).Should(BeFalse())
		})
	})

	Context("When the Recv returns an error", func() {
		It("does not forward the message to the sender", func() {
			fakeStream := newMockIngestorGRPCServer()
			fakeStream.RecvOutput.Ret0 <- nil
			fakeStream.RecvOutput.Ret1 <- errors.New("fake error")
			fakeStream.ContextOutput.Ret0 <- context.TODO()

			go manager.Pusher(fakeStream)
			Consistently(func() bool {
				_, ok := outgoingMsgs.TryNext()
				return ok
			}).Should(BeFalse())
		})
	})

	Context("When the pusher context finishes", func() {
		It("returns the error from the context", func() {
			fakeStream := newMockIngestorGRPCServer()

			for i := 0; i < 100; i++ {
				fakeStream.RecvOutput.Ret0 <- nil
				fakeStream.RecvOutput.Ret1 <- errors.New("fake error")
			}

			context, cancelCtx := context.WithCancel(context.Background())
			fakeStream.ContextOutput.Ret0 <- context
			cancelCtx()

			err := manager.Pusher(fakeStream)

			Expect(err).To(HaveOccurred())
		})
	})
})
