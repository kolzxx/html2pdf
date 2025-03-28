package logger

type Logger interface {
	Info(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
	Debug(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
	SetCorrelationId(correlationId string)
	GetCorrelationId() string
}
