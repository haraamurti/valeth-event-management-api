package config

import (
	"os"

	"github.com/joho/godotenv"
)

func Load() {
    _ = godotenv.Load()
}

func Get(key, def string) string {
    v := os.Getenv(key)
    if v == "" {
        return def
    }
    return v
}
