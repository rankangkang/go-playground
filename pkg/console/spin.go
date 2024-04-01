package console

import (
	"strings"
	"time"

	"github.com/theckman/yacspin"
)

type consoleSpin struct {
	spinner *yacspin.Spinner
}

func Spinner() *consoleSpin {
	cfg := yacspin.Config{
		Frequency:         100 * time.Millisecond,
		CharSet:           yacspin.CharSets[14],
		Suffix:            " ",
		Message:           "loading",
		SuffixAutoColon:   true,
		ColorAll:          true,
		Colors:            []string{"fgYellow"},
		StopCharacter:     IconSuccess,
		StopColors:        []string{"fgGreen"},
		StopMessage:       "done",
		StopFailCharacter: IconFailure,
		StopFailColors:    []string{"fgRed"},
		StopFailMessage:   "failed",
	}

	spinner, _ := yacspin.New(cfg)
	return &consoleSpin{spinner: spinner}
}

func (s *consoleSpin) Update(msgs ...string) *consoleSpin {
	msg := strings.Join(msgs, " ")
	if len(msg) > 0 {
		s.spinner.Message(msg)
	}
	return s
}

func (s *consoleSpin) Start(msgs ...string) *consoleSpin {
	s.Update(msgs...)
	_ = s.spinner.Start()
	return s
}

func (s *consoleSpin) Success(msgs ...string) *consoleSpin {
	msg := strings.Join(msgs, " ")
	if len(msg) > 0 {
		s.spinner.StopMessage(msg)
	}

	_ = s.spinner.Stop()
	return s
}

func (s *consoleSpin) Fail(msgs ...string) *consoleSpin {
	msg := strings.Join(msgs, " ")
	if len(msg) > 0 {
		s.spinner.StopFailMessage(msg)
	}

	_ = s.spinner.StopFail()
	return s
}

func (s *consoleSpin) Pause(msgs ...string) *consoleSpin {
	s.Update(msgs...)
	_ = s.spinner.Pause()
	return s
}

func (s *consoleSpin) Unpause(msgs ...string) *consoleSpin {
	s.Update(msgs...)
	_ = s.spinner.Unpause()
	return s
}
