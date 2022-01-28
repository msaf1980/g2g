package g2g

import (
	"io"
	"net"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/msaf1980/g2g/pkg/expvars"
)

var testExpvar = expvars.NewInt("i")

func TestPublish(t *testing.T) {
	testPub(t, "localhost:2003", NewMockGraphite(t, "tcp://:2003"))
	testPub(t, "unix:///tmp/test.sock", NewMockGraphite(t, "unix:///tmp/test.sock"))
}

func testPub(t *testing.T, address string, mock *MockGraphite) {
	// setup
	d := 25 * time.Millisecond
	var g *Graphite
	g = NewGraphite(address, d, d)

	// register, wait, check
	testExpvar.Set(34)
	g.Register("test.foo.i", testExpvar)

	time.Sleep(2*d + d/10)
	count := mock.Count()
	if count < 1 || count > 2 {
		t.Errorf("expected 1 <= count <= 2, got %d", count)
	}
	// t.Logf("after %s, count=%d", 2*d, count)

	time.Sleep(2*d + d/10)
	count = mock.Count()
	if count < 3 || count > 4 {
		t.Errorf("expected 3 <= count <= 4, got %d", count)
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

//
//
//

type MockGraphite struct {
	t       *testing.T
	network string
	address string
	count   int
	mtx     sync.Mutex
	ln      net.Listener
	done    chan bool
}

func NewMockGraphite(t *testing.T, address string) *MockGraphite {
	network, address := splitEndpoint(address)
	m := &MockGraphite{
		t:       t,
		network: network,
		address: address,
		count:   0,
		mtx:     sync.Mutex{},
		ln:      nil,
		done:    make(chan bool, 1),
	}
	go m.loop()
	return m
}

func (m *MockGraphite) Count() int {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	return m.count
}

func (m *MockGraphite) Shutdown() {
	if m.ln != nil {
		m.ln.Close()
		<-m.done
	}
}

func (m *MockGraphite) loop() {
	ln, err := net.Listen(m.network, m.address)
	if err != nil {
		panic(err)
	}
	m.ln = ln
	for {
		conn, err := m.ln.Accept()
		if err != nil {
			m.done <- true
			return
		}
		go m.handle(conn)
	}
}

func (m *MockGraphite) handle(conn net.Conn) {
	b := make([]byte, 1024)
	for {
		n, err := conn.Read(b)
		if err != nil {
			if err == io.EOF {
				return
			}
			m.t.Logf("Mock Graphite: read error: %s", err)
			return
		}
		if n > 256 {
			m.t.Errorf("Mock Graphite: read %dB: too much data", n)
			return
		}
		s := string(b[:n])
		count := strings.Count(s, "\n")
		s = strings.TrimSpace(s)
		m.t.Logf("Mock Graphite: read %dB/%dM: %s", n, count, s)
		m.mtx.Lock()
		m.count += count
		m.mtx.Unlock()
	}
}
