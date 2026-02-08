package pool

import "sync"

// Resetter интерфейс для типов, поддерживающих сброс состояния
type Resetter interface {
	Reset()
}

// Pool потокобезопасный пул объектов типа T, реализующих метод Reset()
type Pool[T Resetter] struct {
	mu    sync.Mutex
	items []T
	newFn func() T
}

// New создаёт новый пул объектов типа T.
// Параметр newFn — фабричная функция для создания нового объекта, когда пул пуст.
func New[T Resetter](newFn func() T) *Pool[T] {
	return &Pool[T]{
		items: make([]T, 0),
		newFn: newFn,
	}
}

// Get возвращает объект из пула.
// Если пул пуст — создаёт новый объект через фабричную функцию.
func (p *Pool[T]) Get() T {
	p.mu.Lock()
	defer p.mu.Unlock()

	var item T
	if len(p.items) > 0 {
		item = p.items[len(p.items)-1]
		p.items = p.items[:len(p.items)-1]
	} else {
		item = p.newFn()
	}

	return item
}

// Put возвращает объект в пул после сброса его состояния.
// Перед помещением в пул автоматически вызывается метод Reset().
func (p *Pool[T]) Put(obj T) {
	obj.Reset()

	p.mu.Lock()
	defer p.mu.Unlock()

	p.items = append(p.items, obj)
}
