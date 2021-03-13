package cache

import (
	"sync"
)

type ExtendedWidget struct {
	impl     interface{}
	implLock sync.RWMutex
}

// ExtendBaseWidget is used by an extending widget to make use of BaseWidget functionality.
func (w *ExtendedWidget) ExtendBaseWidget(wid interface{}) {
	impl := w.super()
	if impl != nil {
		return
	}

	w.implLock.Lock()
	w.impl = wid
	w.implLock.Unlock()
}

func (w *ExtendedWidget) super() interface{} {
	w.implLock.RLock()
	impl := w.impl
	w.implLock.RUnlock()

	if impl == nil {
		return nil
	}

	return impl
}

type isExtendedWidget interface {
	ExtendBaseWidget(wid interface{})
	super() interface{}
}

func ExtendedSuper(wid interface{}) interface{} {
	if wid == nil {
		return nil
	}

	if wd, ok := wid.(isExtendedWidget); ok {
		if wd.super() != nil {
			wid = wd.super()
		}
	}

	return wid
}
