package logger

import (
	"github.com/nathan-osman/go-herald"
)

const messageContainerUpdateType = "container-update"

type messageContainerUpdate struct {
	ID     string `json:"id"`
	Action string `json:"action"`
}

func (l *Logger) LogContainerUpdate(id, action string) {
	if l.herald != nil {
		m, err := herald.NewMessage(
			messageContainerUpdateType,
			&messageContainerUpdate{
				ID:     id,
				Action: action,
			},
		)
		if err != nil {
			// TODO: log error
		}
		l.herald.Send(m, nil)
	}
}
