package fastlystats

type Config struct {
	FastlyAPIKey      string `env:"FASTLY_API_KEY"`
	FastlyService     string `env:"FASTLY_SERVICE"`
	NewRelicInsertKey string `env:"NEWRELIC_INSERT_KEY"`
}
