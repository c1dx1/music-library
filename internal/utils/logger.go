package utils

import (
	"github.com/sirupsen/logrus"
	"music-library/config"
	"os"
)

func LoadLogger(cfg *config.Config) *logrus.Logger {
	log := logrus.New()

	log.SetOutput(os.Stdout)
	log.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})

	level, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		log.Warn("Invalid LOG_LEVEL, defaulting to INFO")
		level = logrus.InfoLevel
	}

	log.SetLevel(level)

	return log
}
