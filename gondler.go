package gondler

import (
	"github.com/lrita/cmap"
)

type Gondler[T any] struct {
	source    chan T
	callbacks cmap.Map[string, func(T)]
	match     func(T) string
}

func New[T any](source chan T, match func(T) string) *Gondler[T] {
	return &Gondler[T]{
		source:    source,
		callbacks: cmap.Map[string, func(T)]{},
		match:     match,
	}
}

func (a *Gondler[T]) RunSync() {
	for message := range a.source {
		a.handle(message)
	}
}

func (a *Gondler[T]) RunAsync() {
	for message := range a.source {
		go a.handle(message)
	}
}

func (a *Gondler[T]) On(match string, callback func(T)) {
	a.callbacks.Store(match, callback)
}

func (a *Gondler[T]) handle(message T) {
	callback, ok := a.callbacks.Load(a.match(message))
	if !ok {
		return
	}

	callback(message)
}
