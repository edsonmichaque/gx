package stdlib

import "log"

type Logger struct{}

func (l Logger) Error(args ...interface{}) {
	log.Println(args...)
}

func (l Logger) Errorf(args ...interface{}) {
	log.Println(args...)
}

func (l Logger) Debug(args ...interface{}) {
	log.Println(args...)
}

func (l Logger) Debugf(args ...interface{}) {
	log.Println(args...)
}

func (l Logger) Warn(args ...interface{}) {
	log.Println(args...)
}

func (l Logger) Warnf(args ...interface{}) {
	log.Println(args...)
}

func (l Logger) Info(args ...interface{}) {
	log.Println(args...)
}

func (l Logger) Infof(args ...interface{}) {
	log.Println(args...)
}
