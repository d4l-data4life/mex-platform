package try

import (
	"fmt"
	"time"
)

type PauseStrategy interface {
	Pause() error
}

type ConstantPauseStrategy time.Duration

func (c ConstantPauseStrategy) Pause() error {
	time.Sleep(time.Duration(c))
	return nil
}

type MaxAttemptsConstantPauseStrategy struct {
	MaxAttempts   uint32
	PauseStrategy PauseStrategy
}

func (s *MaxAttemptsConstantPauseStrategy) Pause() error {
	if s.MaxAttempts == 0 {
		return fmt.Errorf("max attempts exhausted")
	}

	err := s.PauseStrategy.Pause()
	if err != nil {
		return err
	}

	s.MaxAttempts--
	return nil
}

func NewMaxAttemptsConstantPauseStrategy(maxAttempts uint32, pause time.Duration) PauseStrategy {
	return &MaxAttemptsConstantPauseStrategy{maxAttempts, ConstantPauseStrategy(pause)}
}
