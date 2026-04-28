package config

import (
	"os"
	"strconv"
)

type Config struct {
	SEM_MAX int64
	TIMEOUT int64
	MONGODB_URI string
}

func Load() *Config {
	cfg := &Config{
		SEM_MAX: getEnvInt("SEM_MAX", 100),
		TIMEOUT: getEnvInt("TIMEOUT", 30),
		MONGODB_URI: getEnvStr("MONGODB_URI", ""),
	}
	return cfg
}

func getEnvStr(key, fallback string) string {
	v, ok := os.LookupEnv(key)
	if ok {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int64) int64 {
	v, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	i, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return fallback
	}
	return i
}
