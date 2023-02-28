package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/looplab/fsm"
	"github.com/yvv4git/task-wow/internal/storage"
	"github.com/yvv4git/task-wow/pkg/pow"
	"github.com/yvv4git/task-wow/pkg/utils"
)

const (
	MessageChallenge = "challenge\n"

	EventSendChallenge = "event-send-challenge"
	EventApprove       = "event-approve"
	EventDisapprove    = "event-disapprove"

	EventSendWOW = "event-send-wow"

	StateBegin            = "begin"
	StateChallengeSent    = "challenge-sent"
	StateChallengeApprove = "challenge-approve"

	DefaultDifficulty    = 2
	DefaultChallengeSize = 5
)

var (
	ErrWrongEventState = errors.New("wrong event state")
	ErrEndOfIteration  = errors.New("end of iteration")
	ErrChallengeFiled  = errors.New("challenge failed")
)

type WOW struct {
	fsmProcessor *fsm.FSM
	powProcessor pow.POW
	challenge    string
	wowStorage   storage.WOW
}

func NewWOW(pow pow.POW) *WOW {
	instanceFSM := fsm.NewFSM(
		StateBegin,
		fsm.Events{
			{Name: EventSendChallenge, Src: []string{StateBegin}, Dst: StateChallengeSent},
			{Name: EventApprove, Src: []string{StateChallengeSent}, Dst: StateChallengeApprove},
			{Name: EventDisapprove, Src: []string{StateChallengeSent}, Dst: StateBegin},
			{Name: EventSendWOW, Src: []string{StateChallengeApprove}, Dst: StateBegin},
		},
		fsm.Callbacks{},
	)

	return &WOW{
		fsmProcessor: instanceFSM,
		powProcessor: pow,
		wowStorage:   storage.NewWOW(),
	}
}

func (w *WOW) Send(ctx context.Context) (string, error) {
	switch w.fsmProcessor.Current() {
	case StateBegin:
		if err := w.fsmProcessor.Event(ctx, EventSendChallenge); err != nil {
			return "", fmt.Errorf("error on event send challenge: %w", ErrWrongEventState)
		}

		// send challenge to client
		challenge, err := utils.GenerateStr(DefaultChallengeSize)
		if err != nil {
			return "", fmt.Errorf("error on get random string: %w", err)
		}

		w.challenge = challenge

		return fmt.Sprintf("%s\n", w.challenge), nil
	case StateChallengeApprove:
		if err := w.fsmProcessor.Event(ctx, EventSendWOW); err != nil {
			return "", fmt.Errorf("error on event send wow: %w", ErrWrongEventState)
		}

		wowMessagePhrase := w.wowStorage.LoadPhrase()

		return fmt.Sprintf("%s\n", wowMessagePhrase), nil
	default:
		return "", fmt.Errorf("error on unknown event: %w", ErrWrongEventState)
	}
}

func (w *WOW) Receive(ctx context.Context, value string) error {
	processingChallenge := func(value string) error {
		var proof pow.Proof
		if err := json.Unmarshal([]byte(value), &proof); err != nil {
			return fmt.Errorf("error on unmarshal challenge: %w", err)
		}

		if w.powProcessor.Check([]byte(w.challenge), proof) {
			if err := w.fsmProcessor.Event(ctx, EventApprove); err != nil {
				return fmt.Errorf("error on send wow: %w", err)
			}
		} else {
			return fmt.Errorf("error on receive challenge result: %w", ErrChallengeFiled)
		}

		return nil
	}

	switch w.fsmProcessor.Current() {
	case StateChallengeSent:
		return processingChallenge(value)
	case StateBegin:
		return fmt.Errorf("error on recive in begin state: %w", ErrEndOfIteration)
	default:
		return fmt.Errorf("error with processing unknown state: %w", ErrWrongEventState)
	}
}

func (w *WOW) State() string {
	return w.fsmProcessor.Current()
}
