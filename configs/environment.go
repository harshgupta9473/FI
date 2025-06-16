package configs

import (
	"fmt"
	"github.com/spf13/viper"
)

type environment struct {
	DBConnStr string `mapstructure:"DB_CONN_STR"`
}

func LoadEnvironment() (*environment, error) {
	v := viper.New()
	v.AddConfigPath("./configs")
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading base config: %w", err)
	}
	var env environment
	if err := v.Unmarshal(&env); err != nil {
		return nil, fmt.Errorf("error unmarshalling base config: %w", err)
	}
	return &env, nil
}
