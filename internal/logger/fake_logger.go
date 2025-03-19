package logger

type FakeLogger struct {
	CorrelationId string
	Message       string
	Fields        []interface{}
}

func NewFakeLogger() Logger {
	return &FakeLogger{
		CorrelationId: "",
	}
}

func (fl *FakeLogger) Info(msg string, fields ...interface{}) {
	fl.Message = msg
	fl.Fields = fields
}
func (fl FakeLogger) Warn(msg string, fields ...interface{})  {}
func (fl FakeLogger) Debug(msg string, fields ...interface{}) {}
func (fl FakeLogger) Error(msg string, fields ...interface{}) {}
func (fl FakeLogger) GetCorrelationId() string {
	return fl.CorrelationId
}
func (fl *FakeLogger) SetCorrelationId(correlationId string) {
	fl.CorrelationId = correlationId
}
