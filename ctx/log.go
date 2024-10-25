package ctx

import (
	"errors"
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

// SetLogger handle setup log level, and given hook config
func SetLogger(configs ...*Config) error {
	if settleLogger != nil && settleLogger.Hooks != nil {
		newHooks := make(logrus.LevelHooks)
		settleLogger.ReplaceHooks(newHooks)
	}
	settleLogger = nil
	if configs == nil {
		return errors.New("set logger without config")
	}
	settleLogger = logrus.StandardLogger()
	logHooks := []logrus.Hook{}
	if len(configs) == 0 {
		conf := &Config{
			Environment: "localhost",
			LogFormat:   &logrus.TextFormatter{ForceColors: true},
		}
		h, _ := NewHook(conf)
		settleLogger.AddHook(h)
		return nil
	}

	for _, config := range configs {
		settleLogger.SetOutput(io.Discard)
		if config.ExportToFile == nil && !config.SetAsStdOutput {
			h, _ := NewHook(config)
			logHooks = append(logHooks, h)
			continue
		}

		if config.ExportToFile != nil {
			_, err := os.OpenFile(*config.ExportToFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				settleLogger.WithField("err", err).Panic("ctx set logger export file path failed")
				return err
			}
			h, err := NewHook(config)
			if err != nil {
				settleLogger.WithField("err", err).Fatal("setup logger with config failed")
				return err
			}

			logHooks = append(logHooks, h)
			continue
		}

		if config.SetAsStdOutput {
			logLevels := config.ExportLevel
			infoConfig := *config
			infoConfig.ExportLevel = []logrus.Level{}
			errConfig := *config
			errConfig.ExportLevel = []logrus.Level{}
			for _, lvl := range logLevels {
				if int(lvl) < int(logrus.InfoLevel) {
					errConfig.ExportLevel = append(errConfig.ExportLevel, lvl)
					continue
				}
				infoConfig.ExportLevel = append(infoConfig.ExportLevel, lvl)
			}

			if len(errConfig.ExportLevel) > 0 {
				f := os.Stderr.Name()
				errConfig.ExportToFile = &f
				errHook, _ := NewHook(&errConfig)

				logHooks = append(logHooks, errHook)
			}

			if len(infoConfig.ExportLevel) > 0 {
				f := os.Stdout.Name()
				infoConfig.ExportToFile = &f
				normalH, _ := NewHook(&infoConfig)

				logHooks = append(logHooks, normalH)
			}
			continue
		}
	}

	for _, h := range logHooks {
		settleLogger.AddHook(h)
	}

	return nil
}

// GetLogLevel returns log set level
func GetLogLevel() logrus.Level {
	l := logrus.GetLevel()

	return l
}

// SetDebugLevel for developer debug use case
func SetDebugLevel() {
	logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true, FullTimestamp: true})
}
