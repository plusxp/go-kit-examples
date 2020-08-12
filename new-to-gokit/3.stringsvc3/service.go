package main

import (
	"context"
	"errors"
)

// StringService provides operations on strings.
type StringService interface {
	Uppercase(string) (string, error)
	Count(string) int
}

type stringService struct {
	instances func() string
}

func (svc *stringService) Uppercase(s string) (string, error) {
	//if s == "" {
	//	return "", ErrEmpty
	//}
	//return strings.ToUpper(s), nil
	ep := svc.getRPCEndpoint(route_uppercase)
	response, err := ep(context.Background(), uppercaseRequest{S: s})
	if err != nil {
		return "", err
	}

	resp := response.(uppercaseResponse)
	if resp.Err != "" {
		return resp.V, errors.New(resp.Err)
	}
	return resp.V, nil
}

func (stringService) Count(s string) int {
	return len(s)
}

// ErrEmpty is returned when an input string is empty.
var ErrEmpty = errors.New("empty string")

// ServiceMiddleware is a chainable behavior modifier for StringService.
type ServiceMiddleware func(StringService) StringService
