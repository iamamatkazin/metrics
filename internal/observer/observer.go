// Package observer реализует паттерн Наблюдатель для системы аудита.
// Предоставляет механизм публикации-подписки для уведомления нескольких подписчиков
// (файловые и URL приемники) о событиях обновления метрик.
package observer

import (
	"github.com/iamamatkazin/metrics.git/internal/model"
	"github.com/iamamatkazin/metrics.git/internal/observer/file"
	"github.com/iamamatkazin/metrics.git/internal/observer/url"
	sconfig "github.com/iamamatkazin/metrics.git/pkg/config/server"
)

// Observer определяет интерфейс для подписчиков на события аудита.
type Observer interface {
	Send(model.Message)
	GetID() string
}

// Publisher определяет интерфейс для управления наблюдателями.
type Publisher interface {
	register(Observer)
	deregister(Observer)
	notify()
}

// Event представляет реализацию Publisher для рассылки аудит-сообщений.
type Event struct {
	observers map[string]Observer
	Message   model.Message
}

// New создает новый экземпляр Event с настройками аудита.
func New(cfg *sconfig.Config) (*Event, error) {
	e := Event{}

	if cfg.FileAudit != "" {
		sub, err := file.New(cfg, "file")
		if err != nil {
			return nil, err
		}

		e.register(sub)
	}

	if cfg.URLAudit != "" {
		e.register(url.New(cfg, "url"))
	}

	return &e, nil
}

// register добавляет нового наблюдателя в список подписчиков.
func (e *Event) register(o Observer) {
	if e.observers == nil {
		e.observers = make(map[string]Observer)
	}
	e.observers[o.GetID()] = o
}

// deregister удаляет наблюдателя из списка подписчиков.
func (e *Event) deregister(o Observer) {
	delete(e.observers, o.GetID())
}

// notify отправляет сообщение всем зарегистрированным наблюдателям.
func (e *Event) notify() {
	for _, observer := range e.observers {
		observer.Send(e.Message)
	}
}

// Send отправляет сообщение всем подписчикам.
func (e *Event) Send(mes model.Message) {
	e.Message = mes
	e.notify()
}
