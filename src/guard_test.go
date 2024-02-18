package main

import (
	"privacy-guard/src/tv"
	"testing"
)

type SpySleeper struct {
	Calls      int
	Break      bool
	Iterations int
}

func (s *SpySleeper) Sleep() {
	s.Calls++
}

func (s *SpySleeper) Stop() bool {
	if s.Calls < s.Iterations {
		return false
	} else {
		return true
	}
}

type TvMock struct {
	Status tv.Status
	Rule   string
}

func (t *TvMock) GetAddress() string {
	return t.Rule
}

func (t *TvMock) GetStatus() tv.Status {
	if t.Status == tv.Active {
		t.Status = tv.StandBy
	} else {
		t.Status = tv.Active
	}

	return t.Status
}

type RulesArray map[int]string

type BlockerSpy struct {
	SetRuleCount    int
	SetRuleArgs     RulesArray
	RemoveRuleCount int
	RemoveRuleArgs  RulesArray
}

func (b *BlockerSpy) SetRule(rule string) {
	b.SetRuleArgs[b.SetRuleCount] = rule
	b.SetRuleCount++
}

func (b *BlockerSpy) RemoveRule(rule string) {
	b.RemoveRuleArgs[b.RemoveRuleCount] = rule
	b.RemoveRuleCount++
}

func TestInit(t *testing.T) {
	tvMock := &TvMock{
		Status: tv.StandBy,
		Rule:   "tvRule",
	}
	blocker := &BlockerSpy{
		SetRuleArgs:    make(map[int]string),
		RemoveRuleArgs: make(map[int]string),
	}

	Init(tvMock, blocker, tvMock.Status)

	if blocker.SetRuleArgs[0] != tvMock.Rule {
		t.Errorf("Expected to call SetRule on init: %s", blocker.SetRuleArgs[0])
	}

	tvMock.Status = tv.Active
	Init(tvMock, blocker, tvMock.Status)

	if blocker.RemoveRuleArgs[0] != tvMock.Rule {
		t.Errorf("Expected to call RemoveRule on init: %s", blocker.SetRuleArgs[0])
	}
}

func TestWatchInLoop(t *testing.T) {
	tvMock := &TvMock{
		Status: tv.StandBy,
		Rule:   "tvRule",
	}
	blocker := &BlockerSpy{
		SetRuleArgs:    make(map[int]string),
		RemoveRuleArgs: make(map[int]string),
	}
	sleeper := &SpySleeper{
		Iterations: 3,
	}

	WatchInLoop(tvMock, blocker, sleeper, tvMock.Status)

	if len(blocker.SetRuleArgs) != sleeper.Iterations-1 {
		t.Errorf("Expected to call SetRule %d times but found %d", sleeper.Iterations-1, len(blocker.SetRuleArgs))
	}

	if len(blocker.RemoveRuleArgs) != sleeper.Iterations-1 {
		t.Errorf("Expected to RemoveRule %d times but found %d", sleeper.Iterations-1, len(blocker.RemoveRuleArgs))
	}
}
