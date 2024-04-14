package gondler

import (
	"github.com/lrita/cmap"
)

type Gondler[M comparable, F any] struct {
	source    chan F
	callbacks cmap.Map[M, func(F)]
	match     func(F) M
}

func New[M comparable, F any](source chan F, match func(F) M) *Gondler[M, F] {
	return &Gondler[M, F]{
		source:    source,
		callbacks: cmap.Map[M, func(F)]{},
		match:     match,
	}
}

func (a *Gondler[M, F]) RunSync() {
	for message := range a.source {
		a.handle(message)
	}
}

func (a *Gondler[M, F]) RunAsync() {
	for message := range a.source {
		go a.handle(message)
	}
}

func (a *Gondler[M, F]) On(match M, callback func(F)) {
	a.callbacks.Store(match, callback)
}

func (a *Gondler[M, F]) handle(message F) {
	callback, ok := a.callbacks.Load(a.match(message))
	if !ok {
		return
	}

	callback(message)
}
