package fastlystackdriver

type Config struct {
	FastlyAPIKey  string `env:"FASTLY_API_KEY"`
	FastlyService string `env:"FASTLY_SERVICE"`
}
