// This file was generated by github.com/nelsam/hel.  Do not
// edit this code by hand unless you *really* know what you're
// doing.  Expect any changes made manually to be overwritten
// the next time hel regenerates this file.

package grpcconnector_test

import (
	"doppler/dopplerservice"
	"plumbing"

	"github.com/cloudfoundry/dropsonde/metricbatcher"
)

type mockPlumbingReceiver struct {
	RecvCalled chan bool
	RecvOutput struct {
		Ret0 chan *plumbing.Response
		Ret1 chan error
	}
}

func newMockPlumbingReceiver() *mockPlumbingReceiver {
	m := &mockPlumbingReceiver{}
	m.RecvCalled = make(chan bool, 100)
	m.RecvOutput.Ret0 = make(chan *plumbing.Response, 100)
	m.RecvOutput.Ret1 = make(chan error, 100)
	return m
}
func (m *mockPlumbingReceiver) Recv() (*plumbing.Response, error) {
	m.RecvCalled <- true
	return <-m.RecvOutput.Ret0, <-m.RecvOutput.Ret1
}

type mockFinder struct {
	NextCalled chan bool
	NextOutput struct {
		Ret0 chan dopplerservice.Event
	}
}

func newMockFinder() *mockFinder {
	m := &mockFinder{}
	m.NextCalled = make(chan bool, 100)
	m.NextOutput.Ret0 = make(chan dopplerservice.Event, 100)
	return m
}
func (m *mockFinder) Next() dopplerservice.Event {
	m.NextCalled <- true
	return <-m.NextOutput.Ret0
}

type mockMetaMetricBatcher struct {
	BatchCounterCalled chan bool
	BatchCounterInput  struct {
		Name chan string
	}
	BatchCounterOutput struct {
		Ret0 chan metricbatcher.BatchCounterChainer
	}
	BatchAddCounterCalled chan bool
	BatchAddCounterInput  struct {
		Name  chan string
		Delta chan uint64
	}
}

func newMockMetaMetricBatcher() *mockMetaMetricBatcher {
	m := &mockMetaMetricBatcher{}
	m.BatchCounterCalled = make(chan bool, 100)
	m.BatchCounterInput.Name = make(chan string, 100)
	m.BatchCounterOutput.Ret0 = make(chan metricbatcher.BatchCounterChainer, 100)
	m.BatchAddCounterCalled = make(chan bool, 100)
	m.BatchAddCounterInput.Name = make(chan string, 100)
	m.BatchAddCounterInput.Delta = make(chan uint64, 100)
	return m
}
func (m *mockMetaMetricBatcher) BatchCounter(name string) metricbatcher.BatchCounterChainer {
	m.BatchCounterCalled <- true
	m.BatchCounterInput.Name <- name
	return <-m.BatchCounterOutput.Ret0
}
func (m *mockMetaMetricBatcher) BatchAddCounter(name string, delta uint64) {
	m.BatchAddCounterCalled <- true
	m.BatchAddCounterInput.Name <- name
	m.BatchAddCounterInput.Delta <- delta
}

type mockReceiver struct {
	RecvCalled chan bool
	RecvOutput struct {
		Ret0 chan []byte
		Ret1 chan error
	}
}

func newMockReceiver() *mockReceiver {
	m := &mockReceiver{}
	m.RecvCalled = make(chan bool, 100)
	m.RecvOutput.Ret0 = make(chan []byte, 100)
	m.RecvOutput.Ret1 = make(chan error, 100)
	return m
}
func (m *mockReceiver) Recv() ([]byte, error) {
	m.RecvCalled <- true
	return <-m.RecvOutput.Ret0, <-m.RecvOutput.Ret1
}

type mockBatchCounterChainer struct {
	SetTagCalled chan bool
	SetTagInput  struct {
		Key, Value chan string
	}
	SetTagOutput struct {
		Ret0 chan metricbatcher.BatchCounterChainer
	}
	IncrementCalled chan bool
	AddCalled       chan bool
	AddInput        struct {
		Value chan uint64
	}
}

func newMockBatchCounterChainer() *mockBatchCounterChainer {
	m := &mockBatchCounterChainer{}
	m.SetTagCalled = make(chan bool, 100)
	m.SetTagInput.Key = make(chan string, 100)
	m.SetTagInput.Value = make(chan string, 100)
	m.SetTagOutput.Ret0 = make(chan metricbatcher.BatchCounterChainer, 100)
	m.IncrementCalled = make(chan bool, 100)
	m.AddCalled = make(chan bool, 100)
	m.AddInput.Value = make(chan uint64, 100)
	return m
}
func (m *mockBatchCounterChainer) SetTag(key, value string) metricbatcher.BatchCounterChainer {
	m.SetTagCalled <- true
	m.SetTagInput.Key <- key
	m.SetTagInput.Value <- value
	return <-m.SetTagOutput.Ret0
}
func (m *mockBatchCounterChainer) Increment() {
	m.IncrementCalled <- true
}
func (m *mockBatchCounterChainer) Add(value uint64) {
	m.AddCalled <- true
	m.AddInput.Value <- value
}
