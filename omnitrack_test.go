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

type dummy struct{}

func (d *dummy) Handshake(r io.Reader) error {
	return nil
}

func TestNewDispatcherWithDriver(t *testing.T) {
	d := omnitrack.NewDispatcher(
		omnitrack.WithProvider(&dummy{}),
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
