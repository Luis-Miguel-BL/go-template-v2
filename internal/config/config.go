package config

type BootstrapConfig struct {
	Environment string    `mapstructure:"environment"`
	ConfigPath  string    `mapstructure:"config-path"`
	AWS         AWSConfig `mapstructure:"aws"`
}

type Config struct {
	Environment string            `mapstructure:"environment"`
	App         AppConfig         `mapstructure:"app"`
	Logger      LoggerConfig      `mapstructure:"logger"`
	Server      ServerConfig      `mapstructure:"server"`
	Consumer    ConsumerConfig    `mapstructure:"consumer"`
	AWS         AWSConfig         `mapstructure:"aws"`
	Monitor     MonitorConfig     `mapstructure:"monitor"`
	Integration IntegrationConfig `mapstructure:"integration"`
}

type ServerConfig struct {
	Port      int    `mapstructure:"port"`
	Prefix    string `mapstructure:"prefix"`
	AppKey    string `mapstructure:"app-key"`
	JWTSecret string `mapstructure:"jwt-secret"`
}

type ConsumerConfig struct {
	Enabled     bool   `mapstructure:"enabled"`
	SQSQueueURL string `mapstructure:"sqs-queue-url"`
}

type LoggerConfig struct {
	Level string `mapstructure:"level"`
}

type AppConfig struct {
	Name       string `mapstructure:"name"`
	InMemoryDB bool   `mapstructure:"in-memory-db"`
}

type MonitorConfig struct {
	Enabled        bool           `mapstructure:"enabled"`
	NewRelicConfig NewRelicConfig `mapstructure:"new-relic"`
}

type NewRelicConfig struct {
	AppKey                 string `mapstructure:"app-key"`
	CustomEventPrefix      string `mapstructure:"custom-event-prefix"`
	Endpoint               string `mapstructure:"endpoint"`
	ShutdownTimeoutSeconds int    `mapstructure:"shutdown-timeout-seconds"`
}

type AWSConfig struct {
	Region   string         `mapstructure:"region"`
	Endpoint string         `mapstructure:"endpoint"`
	DynamoDB DynamoDBConfig `mapstructure:"dynamodb"`
	SSM      SSMConfig      `mapstructure:"ssm"`
}

type DynamoDBConfig struct {
	LeadTableName string `mapstructure:"lead-table-name"`
}

type SSMConfig struct {
	LoadFromSSM     bool   `mapstructure:"load-from-ssm"`
	ParameterPrefix string `mapstructure:"parameter-prefix"`
}

type IntegrationConfig struct {
	ExampleAPI IntegrationExampleAPIConfig `mapstructure:"example-api"`
}

type IntegrationExampleAPIConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	BaseURL string `mapstructure:"base-url"`
	APIKey  string `mapstructure:"api-key"`
}

func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}
func (c *Config) IsLocal() bool {
	return c.Environment == "local" || c.Environment == "test"
}
