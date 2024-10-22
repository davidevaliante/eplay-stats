package env

import (
	"fmt"
	"os"
)

type EnvContainer struct {
	SubriptionKey string
	ServiceToken  string
	PartnerId     string
	Port          string
}

var Env *EnvContainer

func Read() error {
	// Get environment variables
	subscriptionKey := os.Getenv("SUBSCRIPTION_KEY")
	serviceToken := os.Getenv("SERVICE_TOKEN")
	partnerID := os.Getenv("PARTNER_ID")
	port := os.Getenv("PORT")

	env := &EnvContainer{
		SubriptionKey: subscriptionKey,
		ServiceToken:  serviceToken,
		PartnerId:     partnerID,
		Port:          port,
	}

	fmt.Println(env)

	Env = env

	return nil
}
