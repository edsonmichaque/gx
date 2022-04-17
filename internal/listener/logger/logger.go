package logger

type Logger interface {
	Error(...interface{})
	Errorf(...interface{})
	Debug(...interface{})
	Debugf(...interface{})
	Warn(...interface{})
	Warnf(...interface{})
	Info(...interface{})
	Infof(...interface{})
}
