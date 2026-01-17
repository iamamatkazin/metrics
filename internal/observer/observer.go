package observer

import (
	"github.com/iamamatkazin/metrics.git/internal/model"
	"github.com/iamamatkazin/metrics.git/internal/observer/file"
	"github.com/iamamatkazin/metrics.git/internal/observer/url"
	sconfig "github.com/iamamatkazin/metrics.git/pkg/config/server"
)

type Observer interface {
	Send(model.Message)
	GetID() string
}

// интерфейсы
type Publisher interface {
	register(Observer)
	deregister(Observer)
	notify()
}

// реализация publisher
type Event struct {
	observers map[string]Observer
	Message   model.Message
}

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

func (e *Event) register(o Observer) {
	if e.observers == nil {
		e.observers = make(map[string]Observer)
	}
	e.observers[o.GetID()] = o
}

func (e *Event) deregister(o Observer) {
	delete(e.observers, o.GetID())
}

func (e *Event) notify() {
	for _, observer := range e.observers {
		observer.Send(e.Message)
	}
}

func (e *Event) Send(mes model.Message) {
	e.Message = mes
	e.notify()
}
