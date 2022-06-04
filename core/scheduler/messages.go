package scheduler

import (
	"github.com/rs/xid"

	"github.com/Raphy42/weekend/core/message"
)

const (
	MSchedule           = "wk.scheduler.schedule"
	MScheduled          = "wk.scheduler#scheduled"
	MProgress           = "wk.scheduler#progress"
	MSuccess            = "wk.scheduler#success"
	MFailure            = "wk.scheduler#failure"
	MSupervise          = "wk.supervisor.supervise"
	MSupervised         = "wk.supervisor#scheduled"
	MSupervisionFailed  = "wk.supervisor#failed"
	MSupervisionSuccess = "wk.supervisor#success"
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

func NewSuperviseMessage(name string, args interface{}) message.Message {
	return message.New(MSupervise, &SuperviseMessagePayload{ScheduleMessagePayload{
		Name: name,
		Args: args,
	}})
}

type SupervisedMessagePayload struct {
	ScheduledMessagePayload
	SupervisorID xid.ID
}

func NewSupervisedMessage(name string, id, parent, supervisorID xid.ID) message.Message {
	return message.New(MSupervised, &SupervisedMessagePayload{
		ScheduledMessagePayload: ScheduledMessagePayload{
			Name:     name,
			ID:       id,
			ParentID: parent,
		},
		SupervisorID: supervisorID,
	})
}

type SupervisionSuccessMessagePayload struct {
	SuccessMessagePayload
	SupervisorID xid.ID
}

func NewSupervisionSuccessMessage(id, supervisorID xid.ID) message.Message {
	return message.New(MSupervisionSuccess, &SupervisionSuccessMessagePayload{
		SuccessMessagePayload: SuccessMessagePayload{
			ID: id,
		},
		SupervisorID: supervisorID,
	})
}

type SupervisionFailureMessagePayload struct {
	FailureMessagePayload
	SupervisorID xid.ID
}

func NewSupervisionFailureMessage(id, supervisorID xid.ID, err error) message.Message {
	return message.New(MSupervisionFailed, &SupervisionFailureMessagePayload{
		FailureMessagePayload: FailureMessagePayload{
			ID:     id,
			Reason: err.Error(),
		},
		SupervisorID: supervisorID,
	})
}
