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

func (c *Config) Client() (*didiyun.Client, error) {
	client := didiyun.Client{
		AccessToken: c.AccessToken,
	}

	return &CombinedConfig{
		client: didiyun.Client,
	}
}
