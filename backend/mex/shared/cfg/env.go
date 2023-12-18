package cfg

import "os"

type EnvVars interface {
	GetValue(key string) string
}

type OSEnvs struct{}

func (envs *OSEnvs) GetValue(key string) string {
	return os.Getenv(key)
}

type MockedEnvs map[string]string

func (envs MockedEnvs) GetValue(key string) string {
	value, ok := envs[key]
	if !ok {
		return ""
	}
	return value
}
