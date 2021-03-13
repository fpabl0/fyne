package widget

import (
	"fmt"
	"testing"

	"fyne.io/fyne/v2/internal/cache"
	"github.com/stretchr/testify/assert"
)

type first struct {
	cache.ExtendedWidget
}

func (f *first) Hello() string {
	return fmt.Sprintf("Hello from %T", cache.ExtendedSuper(f))
}

type second struct {
	first
}

type third struct {
	second
}

func TestExtendedWidget(t *testing.T) {
	w := &second{}
	// because we don't extend it, Hello() should print *widget.first
	assert.Equal(t, w.Hello(), "Hello from *widget.first")
	// now try extending it
	w.ExtendBaseWidget(w)
	assert.Equal(t, w.Hello(), "Hello from *widget.second")

	w3 := &third{}
	assert.Equal(t, w3.Hello(), "Hello from *widget.first")
	w3.ExtendBaseWidget(w3)
	assert.Equal(t, w3.Hello(), "Hello from *widget.third")
}
