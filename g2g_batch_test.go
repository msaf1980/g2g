package g2g

import (
	"testing"
	"time"

	"github.com/msaf1980/g2g/pkg/expvars"
)

var testExpvar2 = expvars.NewInt("j")

func TestPublishBatch(t *testing.T) {
	testPubBatch(t, "localhost:2003", NewMockGraphite(t, "tcp://:2003"))
	testPubBatch(t, "unix:///tmp/test.sock", NewMockGraphite(t, "unix:///tmp/test.sock"))
}

func testPubBatch(t *testing.T, address string, mock *MockGraphite) {
	// setup
	d := 50 * time.Millisecond
	var g *Graphite
	g = NewGraphiteBatch(address, d, d, 512)

	// register, wait, check
	testExpvar.Set(34)
	g.Register("test.foo.i", testExpvar)
	testExpvar2.Set(2)
	g.Register("test2.foo.j", testExpvar2)

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
