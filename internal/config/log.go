package config

import (
	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/params"
)

type Log struct {
	Level logging.Level
}

func (l *Log) get(env params.Interface) (err error) {
	l.Level, err = l.getLevel(env)
	if err != nil {
		return err
	}
	return nil
}

func (l *Log) getLevel(env params.Interface) (level logging.Level, err error) {
	const envKey = "LOG_LEVEL"
	options := []params.OptionSetter{
		params.Default("info"),
	}
	return env.LogLevel(envKey, options...)
}
