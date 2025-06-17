package configs

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	DBConnStr    string `mapstructure:"DB_CONN_STR"`
}

type Environment struct {
	DBConfig DBConfig `mapstructure:"db_config" json:"db_config"`
}

type DBConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

func LoadEnvironment() (*Config, error) {
	v := viper.New()
	v.AddConfigPath("./configs")
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading base config: %w", err)
	}
	var env Environment
	if err := v.Unmarshal(&env); err != nil {
		return nil, fmt.Errorf("error unmarshalling base config: %w", err)
	}
	return &Config{
		DBConnStr: ConnString(env.DBConfig),
	}, nil
}

func ConnString(c DBConfig) string {
	log.Println(c)
	fmt.Printf("Connecting with host=%s port=%d user=%s db=%s\n",
		c.Host, c.Port, c.User, c.DBName)

	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.DBName, c.SSLMode)
}
