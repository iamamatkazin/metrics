package observer

import (
	"github.com/iamamatkazin/metrics.git/internal/model"
	"github.com/iamamatkazin/metrics.git/internal/observer/file"
	"github.com/iamamatkazin/metrics.git/internal/observer/url"
	sconfig "github.com/iamamatkazin/metrics.git/pkg/config/server"
)

// Observer - интерфейс патерна Наблюдатель.
type Observer interface {
	Send(model.Message)
	GetID() string
}

// Publisher - интерфейс подписчиков на событие.
type Publisher interface {
	register(Observer)
	deregister(Observer)
	notify()
}

// Event - реализация publisher.
type Event struct {
	observers map[string]Observer
	Message   model.Message
}

// New - контструктор для структуры Event.
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

// register - регистрация нового наблюдателя.
func (e *Event) register(o Observer) {
	if e.observers == nil {
		e.observers = make(map[string]Observer)
	}
	e.observers[o.GetID()] = o
}

// deregister - отписка наблюдателя.
func (e *Event) deregister(o Observer) {
	delete(e.observers, o.GetID())
}

// notify - рассылка сообщения.
func (e *Event) notify() {
	for _, observer := range e.observers {
		observer.Send(e.Message)
	}
}

// Send - отослать всем наблюдателям новое сообщение.
func (e *Event) Send(mes model.Message) {
	e.Message = mes
	e.notify()
}
