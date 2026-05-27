package main

import (
	"context"

	v1 "github.com/tyzerrr/aws-log-practice/server/gen/greet/v1"
)

type GreetHandler struct{}

func NewGreetHandler() *GreetHandler {
	return &GreetHandler{}
}

func (s *GreetHandler) Greet(ctx context.Context, req *v1.GreetRequest) (*v1.GreetResponse, error) {
	return nil, nil
}
