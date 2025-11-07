// saga/remote_step.go
package saga

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/xfrr/go-cqrsify/messaging"
)

const defaultRemoteStepTimeout = 30 * time.Second

type RemoteSubjects struct {
	Action     string // e.g. "order.reserve.cmd"
	Compensate string // e.g. "order.reserve.compensate.cmd"
	Timeout    time.Duration
}

type RemotePayload struct {
	messaging.BaseCommand
	RemotePayloadData
}

type RemotePayloadData struct {
	SagaID         string            `json:"sagaId"`
	Step           string            `json:"step"`
	Type           string            `json:"type"` // ACTION|COMPENSATE
	Attempt        int               `json:"attempt"`
	IdempotencyKey string            `json:"idempotencyKey"`
	Input          map[string]any    `json:"input"`
	StepData       map[string]any    `json:"stepData"`
	Metadata       map[string]string `json:"metadata"`
}

type RemoteResult struct {
	messaging.BaseEvent
	RemoteResultData
}

type RemoteResultData struct {
	SagaID  string         `json:"sagaId"`
	Step    string         `json:"step"`
	OK      bool           `json:"ok"`
	Error   string         `json:"error,omitempty"`
	Outputs map[string]any `json:"outputs,omitempty"`
}

func RemoteAction(bus messaging.CommandBusReplier, subj RemoteSubjects) StepAction {
	return func(ctx context.Context, ex *Execution) error {
		cmd := RemotePayload{
			BaseCommand: messaging.NewBaseCommand(subj.Action),
			RemotePayloadData: RemotePayloadData{
				SagaID:         ex.SagaID,
				Step:           ex.Def.Steps[ex.StepIndex].Name,
				Type:           "ACTION",
				Attempt:        ex.Instance.Steps[ex.StepIndex].Attempt,
				IdempotencyKey: ex.IdempotencyKey(),
				Input:          ex.Instance.Input,
				StepData:       ex.StepData,
				Metadata:       ex.Instance.Metadata,
			},
		}

		if subj.Timeout <= 0 {
			subj.Timeout = defaultRemoteStepTimeout
		}

		requestCtx, cancel := context.WithTimeout(ctx, subj.Timeout)
		defer cancel()

		resp, err := messaging.Request[RemoteResult](requestCtx, bus, cmd)
		if err != nil {
			return fmt.Errorf("remote action request failed: %w", err)
		}

		if !resp.OK {
			return errors.New(resp.Error)
		}

		// merge outputs into step data
		for k, v := range resp.Outputs {
			ex.Set(k, v)
		}
		return nil
	}
}

func RemoteCompensation(bus messaging.CommandBusReplier, subj RemoteSubjects) StepCompensation {
	return func(ctx context.Context, ex *Execution) error {
		cmd := RemotePayload{
			BaseCommand: messaging.NewBaseCommand(subj.Compensate),
			RemotePayloadData: RemotePayloadData{
				SagaID:         ex.SagaID,
				Step:           ex.Def.Steps[ex.StepIndex].Name,
				Type:           "COMPENSATE",
				Attempt:        1,
				IdempotencyKey: ex.IdempotencyKey(),
				Input:          ex.Instance.Input,
				StepData:       ex.StepData,
				Metadata:       ex.Instance.Metadata,
			},
		}

		if subj.Timeout <= 0 {
			subj.Timeout = defaultRemoteStepTimeout
		}

		requestCtx, cancel := context.WithTimeout(ctx, subj.Timeout)
		defer cancel()

		resp, err := messaging.Request[RemoteResult](requestCtx, bus, cmd)
		if err != nil {
			return fmt.Errorf("remote compensation request failed: %w", err)
		}
		if !resp.OK {
			return errors.New(resp.Error)
		}
		return nil
	}
}

func RemoteResultOK(cmd RemotePayload, outputs map[string]any) RemoteResult {
	return RemoteResult{
		BaseEvent: messaging.NewBaseEvent(cmd.MessageType() + ".result"),
		RemoteResultData: RemoteResultData{
			Step:    cmd.Step,
			SagaID:  cmd.SagaID,
			OK:      true,
			Outputs: outputs,
		},
	}
}

func RemoteResultError(cmd RemotePayload, errMsg string) RemoteResult {
	return RemoteResult{
		BaseEvent: messaging.NewBaseEvent(cmd.MessageType() + ".error"),
		RemoteResultData: RemoteResultData{
			Step:   cmd.Step,
			SagaID: cmd.SagaID,
			OK:     false,
			Error:  errMsg,
		},
	}
}

func NewRemotePayloadFromJSON(msg messaging.JSONMessage[RemotePayloadData]) messaging.Message {
	return RemotePayload{
		BaseCommand: messaging.NewCommandFromJSON(msg),
		RemotePayloadData: RemotePayloadData{
			SagaID:         msg.Payload.SagaID,
			Step:           msg.Payload.Step,
			Type:           msg.Payload.Type,
			Attempt:        msg.Payload.Attempt,
			IdempotencyKey: msg.Payload.IdempotencyKey,
			Input:          msg.Payload.Input,
			StepData:       msg.Payload.StepData,
			Metadata:       msg.Payload.Metadata,
		},
	}
}

func NewRemoteResultFromJSON(msg messaging.JSONMessage[RemoteResultData]) messaging.Message {
	return RemoteResult{
		BaseEvent: messaging.NewEventFromJSON(msg),
		RemoteResultData: RemoteResultData{
			SagaID:  msg.Payload.SagaID,
			Step:    msg.Payload.Step,
			OK:      msg.Payload.OK,
			Error:   msg.Payload.Error,
			Outputs: msg.Payload.Outputs,
		},
	}
}

func RemotePayloadFromJSONEncoder(payload RemotePayload) messaging.JSONMessage[RemotePayloadData] {
	return messaging.NewJSONMessage(
		messaging.Message(payload.BaseCommand),
		RemotePayloadData{
			SagaID:         payload.SagaID,
			Step:           payload.Step,
			Type:           payload.Type,
			Attempt:        payload.Attempt,
			IdempotencyKey: payload.IdempotencyKey,
			Input:          payload.Input,
			StepData:       payload.StepData,
			Metadata:       payload.Metadata,
		},
	)
}

func RemoteResultFromJSONEncoder(payload RemoteResult) messaging.JSONMessage[RemoteResultData] {
	return messaging.NewJSONMessage(
		messaging.Message(payload.BaseEvent),
		RemoteResultData{
			SagaID:  payload.SagaID,
			Step:    payload.Step,
			OK:      payload.OK,
			Error:   payload.Error,
			Outputs: payload.Outputs,
		},
	)
}

func RemotePayloadFromJSONDecoder(msg messaging.JSONMessage[RemotePayloadData]) (messaging.Message, error) {
	return NewRemotePayloadFromJSON(msg), nil
}

func RemoteResultFromJSONDecoder(msg messaging.JSONMessage[RemoteResultData]) (messaging.Message, error) {
	return NewRemoteResultFromJSON(msg), nil
}
