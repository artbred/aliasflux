package config

import (
	"fmt"
	"os"
)

func ConnectionURLBuilder(n string) string {
	switch n {
		case "postgres":
			return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
				os.Getenv("POSTGRES_USER"),
				os.Getenv("POSTGRES_PASSWORD"),
				os.Getenv("POSTGRES_HOST"),
				os.Getenv("POSTGRES_PORT"),
				os.Getenv("POSTGRES_DB"),
				os.Getenv("POSTGRES_SSLMODE"),
			)
		case "redis":
			return fmt.Sprintf(
				"%s:%s",
				os.Getenv("REDIS_HOST"),
				os.Getenv("REDIS_PORT"),
			)
		case "nats":
			return fmt.Sprintf(
				"%s:%s",
				os.Getenv("NATS_HOST"),
				os.Getenv("NATS_PORT"),
			)
		default:
			return ""
	}
}

