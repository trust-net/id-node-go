// Copyright 2019 The trust-net Authors
// ID Application specific world state wrapper
package app

import (
	"fmt"
	"github.com/trust-net/dag-lib-go/stack/state"
)

type idState struct {
	SubmitterId []byte
	ws          state.State
}

func NewIdState(submitterId []byte, state state.State) *idState {
	return &idState{
		SubmitterId: submitterId,
		ws:          state,
	}
}

func (s *idState) Prefixed(attribute string) []byte {
	return append(s.SubmitterId, []byte(attribute)...)
}

func (s *idState) Get(attribute string) ([]byte, error) {
	if res, err := s.ws.Get(s.Prefixed(attribute)); err != nil {
		return nil, err
	} else if string(s.SubmitterId) != string(res.Owner) {
		return nil, fmt.Errorf("resource not owned")
	} else {
		return res.Value, nil
	}
}

func (s *idState) Put(attribute string, value []byte) error {
	return s.ws.Put(&state.Resource{
		Key:   s.Prefixed(attribute),
		Owner: s.SubmitterId,
		Value: value,
	})
}
