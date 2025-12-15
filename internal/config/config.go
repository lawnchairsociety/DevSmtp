package config

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Auth     AuthConfig     `mapstructure:"auth"`
	TLS      TLSConfig      `mapstructure:"tls"`
}

type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type DatabaseConfig struct {
	Path string `mapstructure:"path"`
}

type AuthConfig struct {
	Required bool   `mapstructure:"required"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type TLSConfig struct {
	Cert string `mapstructure:"cert"`
	Key  string `mapstructure:"key"`
}

func Load(cfgFile string, cmd *cobra.Command) (*Config, error) {
	v := viper.New()

	// Set defaults
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.port", 587)
	v.SetDefault("database.path", "./devsmtp.db")
	v.SetDefault("auth.required", false)
	v.SetDefault("auth.username", "")
	v.SetDefault("auth.password", "")
	v.SetDefault("tls.cert", "")
	v.SetDefault("tls.key", "")

	// Config file
	if cfgFile != "" {
		v.SetConfigFile(cfgFile)
	} else {
		v.SetConfigName("devsmtp")
		v.SetConfigType("yaml")
		v.AddConfigPath(".")
	}

	// Environment variables
	v.SetEnvPrefix("DEVSMTP")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Read config file (ignore if not found)
	_ = v.ReadInConfig()

	// Bind CLI flags if command provided
	if cmd != nil {
		if flag := cmd.Flags().Lookup("host"); flag != nil {
			_ = v.BindPFlag("server.host", flag)
		}
		if flag := cmd.Flags().Lookup("port"); flag != nil {
			_ = v.BindPFlag("server.port", flag)
		}
		if flag := cmd.Flags().Lookup("db"); flag != nil {
			_ = v.BindPFlag("database.path", flag)
		}
		if flag := cmd.Flags().Lookup("auth-required"); flag != nil {
			_ = v.BindPFlag("auth.required", flag)
		}
		if flag := cmd.Flags().Lookup("auth-user"); flag != nil {
			_ = v.BindPFlag("auth.username", flag)
		}
		if flag := cmd.Flags().Lookup("auth-pass"); flag != nil {
			_ = v.BindPFlag("auth.password", flag)
		}
		if flag := cmd.Flags().Lookup("tls-cert"); flag != nil {
			_ = v.BindPFlag("tls.cert", flag)
		}
		if flag := cmd.Flags().Lookup("tls-key"); flag != nil {
			_ = v.BindPFlag("tls.key", flag)
		}
	}

	// Handle environment variable overrides explicitly
	if val := os.Getenv("DEVSMTP_HOST"); val != "" {
		v.Set("server.host", val)
	}
	if val := os.Getenv("DEVSMTP_PORT"); val != "" {
		v.Set("server.port", val)
	}
	if val := os.Getenv("DEVSMTP_DB"); val != "" {
		v.Set("database.path", val)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
