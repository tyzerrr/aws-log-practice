package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type awsClient struct {
	smc    *secretsmanager.Client
	logger *slog.Logger
}

func (c *awsClient) GetSecretValue(ctx context.Context, input *secretsmanager.GetSecretValueInput) ([]byte, error) {
	result, err := c.smc.GetSecretValue(ctx, input)
	if err != nil {
		return nil, err
	}
	if result.SecretString == nil {
		return nil, fmt.Errorf("secret string is empty")
	}
	return []byte(*result.SecretString), nil
}
