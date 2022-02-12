package omnitrack_test

import (
	"io"
	"strings"
	"testing"

	"github.com/edsonmichaque/omnitrack"
)

func TestNewDispatcher(t *testing.T) {
	d := omnitrack.NewDispatcher()

	if d == nil {
		t.Fatal("dispatcher should not be not")
	}
}

type dummyProvider struct{}

func (d *dummyProvider) Handshake(r io.Reader) error {
	return nil
}

func (d *dummyProvider) Handle() error {
	return nil
}

func TestNewDispatcherWithDriver(t *testing.T) {
	d := omnitrack.NewDispatcher(
		omnitrack.WithProvider(&dummyProvider{}),
	)

	if d == nil {
		t.Fatal("dispatcher should not be not")
	}
}

func TestDispatchWithNoProviderAvailable(t *testing.T) {
	d := omnitrack.NewDispatcher()

	reader := strings.NewReader("abcd")

	_, err := d.Dispatch(reader)
	if err == nil {
		t.Fatal("should have returned an error")
	} else {
		if _, ok := err.(omnitrack.NoProviderError); !ok {
			t.Fatal("should have returned an error")
		}
	}
}

func TestDispatch_WithDummyProvider(t *testing.T) {
	d := omnitrack.NewDispatcher(
		omnitrack.WithProvider(&dummyProvider{}),
	)

	reader := strings.NewReader("abcd")

	if provider, err := d.Dispatch(reader); err != nil {
		t.Fatal("shouldn't have returned an error")
	} else {
		if provider == nil {
			t.Fatal("provider shouldn't be nil")
		}
	}
}

type dummyConn struct{}

func (d *dummyConn) Read(b []byte) (int, error) {
	return 0, nil
}

func (d *dummyConn) Write(b []byte) (int, error) {
	return 0, nil
}

func (d *dummyConn) Close() error {
	return nil
}

func TestHandle_WithDummyProvider(t *testing.T) {
	d := omnitrack.NewDispatcher(
		omnitrack.WithProvider(&dummyProvider{}),
	)

	conn := &dummyConn{}

	if err := d.Handle(conn); err != nil {
		t.Fatal("should not have returned error")
	}

}
