package scheduler

import (
	"github.com/rs/xid"

	"github.com/Raphy42/weekend/core/message"
)

// scheduler specifics
const (
	MSchedule  = "wk.scheduler.schedule"
	MScheduled = "wk.scheduler#scheduled"
	MProgress  = "wk.scheduler#progress"
	MSuccess   = "wk.scheduler#success"
	MFailure   = "wk.scheduler#failure"
)

type ScheduleMessagePayload struct {
	Name string
	Args interface{}
}

func NewScheduleMessage(name string, args interface{}) message.Message {
	return message.New(MSchedule, &ScheduleMessagePayload{
		Name: name,
		Args: args,
	})
}

type ScheduledMessagePayload struct {
	Name     string
	ID       xid.ID
	ParentID xid.ID
}

func NewScheduledMessage(name string, id, parentID xid.ID) message.Message {
	return message.New(MScheduled, &ScheduledMessagePayload{Name: name, ID: id, ParentID: parentID})
}

type ProgressMessagePayload struct {
	ID      xid.ID
	Current int
	Max     int
}

func NewProgressMessage(id xid.ID, current, max int) message.Message {
	return message.New(MProgress, &ProgressMessagePayload{ID: id, Current: current, Max: max})
}

type SuccessMessagePayload struct {
	ID xid.ID
}

func NewSuccessMessage(id xid.ID) message.Message {
	return message.New(MSuccess, &SuccessMessagePayload{ID: id})
}

type FailureMessagePayload struct {
	ID     xid.ID
	Reason string
}

func NewFailureMessage(id xid.ID, err error) message.Message {
	return message.New(MFailure, &FailureMessagePayload{ID: id, Reason: err.Error()})
}

type SuperviseMessagePayload struct {
	ScheduleMessagePayload
}
