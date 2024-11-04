package config

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type SecretsManagerInterface interface {
	GetSecretValue(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error)
}

// SecretManagerFunc allows for injecting a custom Secrets Manager function for testing.
var SecretManagerFunc = func() (SecretsManagerInterface, error) {
	cfg, err := loadAWSConfig(context.Background())
	if err != nil {
		return nil, err
	}
	return secretsmanager.NewFromConfig(cfg), nil
}

var loadAWSConfig = config.LoadDefaultConfig

func LoadConfig() (*Config, error) {
	secretName := "testing/splunkToken"

	svc, err := SecretManagerFunc()
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	}

	result, err := svc.GetSecretValue(context.Background(), input)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve secret: %w", err)
	}

	secretString := *result.SecretString
	config := &Config{}

	err = json.Unmarshal([]byte(secretString), config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal secret string: %w", err)
	}

	return config, nil
}
