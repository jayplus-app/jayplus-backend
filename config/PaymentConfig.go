package config

import "os"

var (
	StripeSecretKey     string
	StripeWebhookSecret string
)

func LoadPaymentConfig() {
	StripeSecretKey = os.Getenv("STRIPE_SECRET_KEY")
	StripeWebhookSecret = os.Getenv("STRIPE_WEBHOOK_SECRET")
}
