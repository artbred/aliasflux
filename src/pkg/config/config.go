package config

import (
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

var (
	Debug bool

	ApiPort string

	Domain string
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	ApiPort = os.Getenv("API_PORT")
	Domain = os.Getenv("DOMAIN")

	Debug, _ = strconv.ParseBool(os.Getenv("DEBUG"))
}
