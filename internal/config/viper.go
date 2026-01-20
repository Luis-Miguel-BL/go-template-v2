package config

import (
	"context"
	"fmt"
	"strings"

	"github.com/Luis-Miguel-BL/go-lm-template/internal/infrastructure/aws"
	"github.com/spf13/viper"
)

func Load(cfg *BootstrapConfig, ssmClient *aws.SSMClient) (*Config, error) {
	v := viper.New()

	v.SetConfigName("config")
	v.AddConfigPath("./config/")

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	if cfg.AWS.SSM.LoadFromSSM {
		basePath := fmt.Sprintf("/%s/%s", cfg.AWS.SSM.ParameterPrefix, cfg.Environment)
		if err := loadFromSSM(v, ssmClient, basePath); err != nil {
			return nil, err
		}
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func LoadBootstrapConfig() (*BootstrapConfig, error) {
	v := viper.New()

	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.SetEnvPrefix("APP")

	v.SetDefault("environment", "")
	v.SetDefault("aws.region", "")
	v.SetDefault("aws.endpoint", "")
	v.SetDefault("aws.ssm.load-from-ssm", "")
	v.SetDefault("aws.ssm.parameter-prefix", "")

	var config BootstrapConfig
	if err := v.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func loadFromSSM(
	v *viper.Viper,
	ssmClient *aws.SSMClient,
	basePath string,
) error {
	params, err := ssmClient.GetParametersByPath(context.Background(), basePath)
	if err != nil {
		return err
	}

	for fullKey, value := range params {
		key := strings.TrimPrefix(fullKey, basePath+"/")
		v.Set(key, value)
	}

	return nil
}
