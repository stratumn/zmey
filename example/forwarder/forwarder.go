package forwarder

import (
	"github.com/stratumn/zmey"
)

// Forwarder implements Process interface and runs simple forwarding algorithm
type Forwarder struct {
	pid int
	api zmey.API
}

// FCall represents a message exchanged between forwarder process and client
type FCall struct {
	// SequenceNumber is always-increasing message id
	SequenceNumber int
	// To is recepient id
	To int
	// Payload is just few randomly-generated bytes
	Payload []byte
}

// NewForwarder creates and returns the instance of Forwarder process
func NewForwarder(pid int) zmey.Process {
	return &Forwarder{pid: pid}
}

// Bind implements Process.Bind
func (f *Forwarder) Bind(api zmey.API) {
	f.api = api
}

// ReceiveNet implements Process.ReceiveNet
func (f *Forwarder) ReceiveNet(from int, payload interface{}) {
	t := zmey.NewTracer("ReceiveNet")
	fcall, ok := payload.(FCall)
	if !ok {
		f.api.ReportError(t.Errorf("cannot coerce to FCall: %+v", payload))
		return
	}

	if fcall.To != f.pid {
		f.api.ReportError(t.Errorf("%d: incorrect recepient, should be %d", fcall.To, f.pid))
		return
	}

	f.api.Trace(t.Logf("%d return net", fcall.SequenceNumber))
	f.api.Return(payload)
}

// ReceiveCall implements Process.ReceiveCall
func (f *Forwarder) ReceiveCall(c interface{}) {
	t := zmey.NewTracer("ReceiveCall")
	fcall, ok := c.(FCall)
	if !ok {
		f.api.ReportError(t.Errorf("cannot coerce to FCall: %+v", c))
		return
	}

	t = t.Fork("FCall %d", fcall.SequenceNumber)

	f.api.Trace(t.Logf("receive"))

	if fcall.To == f.pid { // Local call
		f.api.Trace(t.Logf("return local"))
		f.api.Return(c)
		return
	}

	f.api.Trace(t.Logf("forward"))
	f.api.Send(fcall.To, fcall)
}

// Tick implements Process.Tick
func (f *Forwarder) Tick(uint) {

}
