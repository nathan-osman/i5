package logger

import (
	"github.com/nathan-osman/go-herald"
	"github.com/nathan-osman/i5/container"
)

const (
	messageContainerCreatedType   = "container-created"
	messageContainerDestroyedType = "container-destroyed"
	messageContainerStartedType   = "container-started"
	messageContainerStoppedType   = "container-stopped"
)

func (l *Logger) logContainerAction(messageType string, v interface{}) {
	if l.herald != nil {
		m, err := herald.NewMessage(
			messageType,
			v,
		)
		if err != nil {
			// TODO: log error
		}
		l.herald.Send(m, nil)
	}
}

func (l *Logger) LogContainerCreated(c *container.Container) {
	l.logContainerAction(messageContainerCreatedType, c)
}

type messageContainerAction struct {
	ID string `json:"id"`
}

func (l *Logger) LogContainerDestroyed(id string) {
	l.logContainerAction(
		messageContainerDestroyedType,
		&messageContainerAction{ID: id},
	)
}

func (l *Logger) LogContainerStarted(id string) {
	l.logContainerAction(
		messageContainerStartedType,
		&messageContainerAction{ID: id},
	)
}

func (l *Logger) LogContainerStopped(id string) {
	l.logContainerAction(
		messageContainerStoppedType,
		&messageContainerAction{ID: id},
	)
}
