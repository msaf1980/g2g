package g2g

import (
	"testing"
	"time"
)

type intLoop struct {
}

func (v *intLoop) Strings() []string {
	return []string{"1", "2"}
}

func NewIntLoop(name string) *intLoop {
	v := new(intLoop)
	Publish(name, v)
	return v
}

var testExploop = NewIntLoop("loop")

func testPubMVar(t *testing.T, address string, mock *MockGraphite) {
	// setup
	d := 25 * time.Millisecond
	var g *Graphite
	g = NewGraphite(address, d, d)

	// register, wait, check
	g.Register("test.foo.loop", testExploop)

	time.Sleep(2*d + d/10)
	count := mock.Count()
	if count < 1 || count > 4 {
		t.Errorf("expected 2 <= count <= 4, got %d", count)
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
