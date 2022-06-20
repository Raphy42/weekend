package scheduler

import (
	"github.com/rs/xid"

	"github.com/Raphy42/weekend/core"
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

func init() {
	core.RegisterOnStartHook(func() error {
		message.GlobalDecoderRegistry.
			Register(MSchedule, message.NewDecoder[ScheduleMessagePayload]()).
			Register(MScheduled, message.NewDecoder[ScheduledMessagePayload]()).
			Register(MProgress, message.NewDecoder[ProgressMessagePayload]()).
			Register(MSuccess, message.NewDecoder[SuccessMessagePayload]()).
			Register(MFailure, message.NewDecoder[FailureMessagePayload]())

		return nil
	})
}

type ScheduleMessagePayload struct {
	ManifestID xid.ID
	Args       any
}

func NewScheduleMessage(id xid.ID, args any) message.Message {
	return message.New(MSchedule, &ScheduleMessagePayload{
		ManifestID: id,
		Args:       args,
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
	ManifestID xid.ID
	HandleID   xid.ID
	Current    int
	Max        int
}

func NewProgressMessage(id, handleID xid.ID, current, max int) message.Message {
	return message.New(MProgress, &ProgressMessagePayload{ManifestID: id, Current: current, Max: max})
}

type SuccessMessagePayload struct {
	ManifestID xid.ID
	HandleID   xid.ID
}

func NewSuccessMessage(id, handleID xid.ID) message.Message {
	return message.New(MSuccess, &SuccessMessagePayload{
		ManifestID: id,
		HandleID:   handleID,
	})
}

type FailureMessagePayload struct {
	ManifestID xid.ID
	HandleID   xid.ID
	Error      error
}

func NewFailureMessage(id, handleID xid.ID, err error) message.Message {
	return message.New(MFailure, &FailureMessagePayload{
		ManifestID: id,
		HandleID:   handleID,
		Error:      err,
	})
}
