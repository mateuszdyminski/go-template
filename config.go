package main

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type config struct {
	httpPort            int
	httpPprofPort       int
	httpGracefulTimeout int
	httpGracefulSleep   int
	pgHost              string
	pgPort              int
	pgUser              string
	pgDBName            string
	pgPassword          string
}

func loadConfig(l *zap.Logger) (*config, error) {
	viper.SetEnvPrefix("APP") // Set the environment prefix to APP_*
	viper.AutomaticEnv()      // Automatically search for environment variables

	config := &config{
		httpPort:            viper.GetInt("http_port"),
		httpPprofPort:       viper.GetInt("http_pprof_port"),
		httpGracefulTimeout: viper.GetInt("http_graceful_timeout"),
		httpGracefulSleep:   viper.GetInt("http_graceful_sleep"),
		pgHost:              viper.GetString("postgres_host"),
		pgPort:              viper.GetInt("postgres_port"),
		pgUser:              viper.GetString("postgres_user"),
		pgDBName:            viper.GetString("postgres_dbname"),
		pgPassword:          viper.GetString("postgres_password"),
	}

	config.print(l.Sugar())

	return config, nil
}

func (c *config) print(l *zap.SugaredLogger) {
	l.Infow("config value", "http_port", c.httpPort)
	l.Infow("config value", "http_pprof_port", c.httpPprofPort)
	l.Infow("config value", "http_graceful_timeout", c.httpGracefulTimeout)
	l.Infow("config value", "http_graceful_sleep", c.httpGracefulSleep)
	l.Infow("config value", "postgres_host", c.pgHost)
	l.Infow("config value", "postgres_port", c.pgPort)
	l.Infow("config value", "postgres_user", c.pgUser)
	l.Infow("config value", "postgres_dbname", c.pgDBName)
	l.Infow("config value", "postgres_password", maskLeft(c.pgPassword, 4))
}

func maskLeft(s string, l int) string {
	rs := []rune(s)
	for i := 0; i < len(rs)-l; i++ {
		rs[i] = 'X'
	}
	return string(rs)
}
