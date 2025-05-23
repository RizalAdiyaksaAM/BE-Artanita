package config

import (
	"log"
	"os"
)

type MidtransConfig struct {
	ServerKey string
	ClientKey string
}

// InitConfigMidtrans initializes the Midtrans configuration from environment variables
func InitConfigMidtrans() MidtransConfig {
	serverKey := os.Getenv("MIDTRANS_SERVER_KEY")
	clientKey := os.Getenv("MIDTRANS_CLIENT_KEY")

	if serverKey == "" || clientKey == "" {
		log.Fatal("Missing Midtrans server or client key. Please set the MIDTRANS_SERVER_KEY and MIDTRANS_CLIENT_KEY environment variables.")
	}

	return MidtransConfig{
		ServerKey: serverKey,
		ClientKey: clientKey,
	}
}
