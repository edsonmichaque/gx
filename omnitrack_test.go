package omnitrack_test

import (
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

func (d *dummy) Handshake() error {
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

func TestDispatchWithErr(t *testing.T) {
	d := omnitrack.NewDispatcher()

	reader := strings.NewReader("abcd")

	if err := d.Dispatch(reader); err == nil {
		t.Fatal("should have returned an error")
	}
}
