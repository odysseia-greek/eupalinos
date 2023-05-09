package eupalinos

import (
	"fmt"
	"os"
	"strings"
)

const (
	rabbitDefaultAddress  = "localhost:5672"
	rabbitUsernameDefault = "guest"
	rabbitPasswordDefault = "guest"
	EnvRabbitService      = "RABBIT_SERVICE"
	EnvRabbitUser         = "RABBIT_USER"
	EnvRabbitPassword     = "RABBIT_PASSWORD"
)

func CreateQueueConnectionStringFromEnv() string {
	user := os.Getenv(EnvRabbitUser)
	if user == "" {
		user = rabbitUsernameDefault
	}
	password := os.Getenv(EnvRabbitPassword)
	if password == "" {
		password = rabbitPasswordDefault
	}

	service := os.Getenv(EnvRabbitService)
	if service == "" {
		service = rabbitDefaultAddress
	}

	split := strings.Split(service, ":")

	return fmt.Sprintf("amqp://%s:%s@%s:%s/", user, password, split[0], split[1])
}
