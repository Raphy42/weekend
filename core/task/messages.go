package task

import (
	"github.com/rs/xid"

	"github.com/Raphy42/weekend/core"
	"github.com/Raphy42/weekend/core/message"
)

const (
	MWorkerAnnounce = "wk.work.announce"
	MTaskExecute    = "wk.task.execute"
	MTaskExecuting  = "wk.task#executing"
	MTaskExecuted   = "wk.task#executed"
)

func init() {
	core.RegisterOnStartHook(func() error {
		message.GlobalDecoderRegistry.
			Register(MWorkerAnnounce, message.NewDecoder[WorkerAnnounceMessage]()).
			Register(MTaskExecute, message.NewDecoder[ExecuteMessage]()).
			Register(MTaskExecuting, message.NewDecoder[ExecutingMessage]()).
			Register(MTaskExecuted, message.NewDecoder[ExecutedMessage]())
		return nil
	})
}

type WorkerAnnounceMessage struct {
	WorkerID xid.ID
	Tasks    []string
}

func NewWorkerAnnounceMessage(workerID xid.ID, tasks ...string) message.Message {
	return message.New(MWorkerAnnounce, &WorkerAnnounceMessage{
		WorkerID: workerID,
		Tasks:    tasks,
	})
}

type ExecuteMessage struct {
	Manifest
}

func NewTaskExecuteMessage(name string, args []byte, payloadType string) message.Message {
	return message.New(MTaskExecute, &ExecuteMessage{
		Manifest: Manifest{
			Name:        name,
			Payload:     args,
			PayloadType: payloadType,
			Options: Options{
				Priority: 0,
			},
		},
	})
}

type ExecutingMessage struct {
	ManifestID xid.ID
	FutureID   xid.ID
}

func NewTaskExecutingMessage(manifestID, futureID xid.ID) message.Message {
	return message.New(MTaskExecuting, &ExecutingMessage{
		ManifestID: manifestID,
		FutureID:   futureID,
	})
}

type ExecutedMessage struct {
	FutureID xid.ID
	Result   any
	Error    *string
}

func NewTaskExecutedMessage(id xid.ID, result any, err error) message.Message {
	var reason *string
	if err != nil {
		err := err.Error()
		reason = &err
	}
	return message.New(MTaskExecuted, &ExecutedMessage{
		FutureID: id,
		Result:   result,
		Error:    reason,
	})
}
