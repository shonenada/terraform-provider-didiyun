package didiyun

import (
	didiyun "github.com/shonenada/didiyun-go"
)

type Config struct {
	AccessToken string
}

type CombinedConfig struct {
	client *didiyun.Client
}

func (c *CombinedConfig) Client() *didiyun.Client { return c.client }

func (c *Config) Client() *CombinedConfig {
	client := didiyun.Client{
		AccessToken: c.AccessToken,
	}

	return &CombinedConfig{
		client: &client,
	}
}
