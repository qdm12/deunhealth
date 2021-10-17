package env

import (
	"os"
)

type Reader struct {
	getEnv func(key string) (value string)
}

func New() *Reader {
	return &Reader{
		getEnv: os.Getenv,
	}
}
