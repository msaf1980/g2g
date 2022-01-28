package g2g

import (
	"testing"
	"time"

	"github.com/msaf1980/g2g/pkg/expvars"
)

type loopSlice struct {
}

func (v *loopSlice) Strings() []expvars.MValue {
	return []expvars.MValue{
		{Name: "min", V: "1"},
		{Name: "max", V: "2"},
	}
}

var testExpvarLoop = &loopSlice{}

func TestPubMulti(t *testing.T) {
	address := "localhost:2003"
	mock := NewMockGraphite(t, "tcp://:2003")

	// setup
	d := 50 * time.Millisecond
	var g *Graphite
	g = NewGraphite(address, d, d)

	// register, wait, check
	g.MRegister("test2.foo.j", testExpvarLoop)

	time.Sleep(2*d + d/10)
	count := mock.Count()
	if count < 3 || count > 4 {
		t.Errorf("expected 3 <= count <= 4, got %d", count)
	}
	// t.Logf("after %s, count=%d", 2*d, count)

	time.Sleep(2*d + d/10)
	count = mock.Count()
	if count < 6 || count > 8 {
		t.Errorf("expected 6 <= count <= 8, got %d", count)
	}
	// t.Logf("after second %s, count=%d", 2*d, count)

	// teardown
	ok := make(chan bool)
	go func() {
		g.Shutdown()
		mock.Shutdown()
		ok <- true
	}()
	select {
	case <-ok:
		t.Logf("shutdown OK")
	case <-time.After(d):
		t.Errorf("timeout during shutdown")
	}

}

func TestPubMultiBatch(t *testing.T) {
	address := "localhost:2003"
	mock := NewMockGraphite(t, "tcp://:2003")

	// setup
	d := 50 * time.Millisecond
	var g *Graphite
	g = NewGraphiteBatch(address, d, d, 512)

	// register, wait, check
	g.MRegister("test2.foo.j", testExpvarLoop)

	time.Sleep(2*d + d/10)
	count := mock.Count()
	if count < 3 || count > 4 {
		t.Errorf("expected 3 <= count <= 4, got %d", count)
	}
	// t.Logf("after %s, count=%d", 2*d, count)

	time.Sleep(2*d + d/10)
	count = mock.Count()
	if count < 6 || count > 8 {
		t.Errorf("expected 6 <= count <= 8, got %d", count)
	}
	// t.Logf("after second %s, count=%d", 2*d, count)

	// teardown
	ok := make(chan bool)
	go func() {
		g.Shutdown()
		mock.Shutdown()
		ok <- true
	}()
	select {
	case <-ok:
		t.Logf("shutdown OK")
	case <-time.After(d):
		t.Errorf("timeout during shutdown")
	}

}
