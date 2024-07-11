package initializers

import (
	"github.com/joho/godotenv"
	"github.com/zetamatta/go-outputdebug"
	"time"
)

func LoadEnvVariables() {
	err := godotenv.Load()
	if err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())

	}
}
