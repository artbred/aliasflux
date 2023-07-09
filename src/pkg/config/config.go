package config

import (
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

var (
	Debug bool

	ApiPort     string
	FrontDomain string

	GoDaddyApiKey    string
	GoDaddyApiSecret string
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	ApiPort = os.Getenv("API_PORT")
	FrontDomain = os.Getenv("DOMAIN")

	Debug, _ = strconv.ParseBool(os.Getenv("DEBUG"))

	GoDaddyApiKey = os.Getenv("GODADDY_API_KEY")
	GoDaddyApiSecret = os.Getenv("GODADDY_API_SECRET")
}
