package logger

import (
	"context"
	"os"

	"github.com/kolzxx/html2pdf/configs"
	"go.elastic.co/ecszap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const correlationIdKey string = "x-correlation-id"

var (
	LoggerKey = "logger"
)

type labels struct {
	Application string `json:"application"`
	Environment string `json:"environment"`
}
type Trace struct {
	ID string `json:"id"`
}

type ecsLogger struct {
	logger *zap.Logger
	ctx    context.Context
}

// An Option configures a EcsLogger.
type Option interface {
	apply(*ecsLogger)
}

// optionFunc wraps a func so it satisfies the Option interface.
type optionFunc func(*ecsLogger)

func NewEcsLogger(ctx context.Context, opts ...Option) Logger {
	// @see https://www.elastic.co/guide/en/ecs-logging/go-zap/current/intro.html
	encoderConfig := ecszap.EncoderConfig{
		EncodeDuration: zapcore.MillisDurationEncoder,
	}
	core := ecszap.NewCore(
		encoderConfig,
		os.Stdout,
		zap.DebugLevel,
	)

	cfg := configs.GetConfig()

	logger := zap.New(core, zap.AddCaller()).With(zap.Any("labels", labels{
		Application: cfg.Log.Application,
		Environment: cfg.Log.Environment,
	}))

	defer logger.Sync()

	log := &ecsLogger{
		logger: logger,
		ctx:    ctx,
	}

	return log.WithOptions(opts...)
}

func (ecs *ecsLogger) WithOptions(opts ...Option) *ecsLogger {
	c := ecs.clone()
	for _, opt := range opts {
		opt.apply(c)
	}
	return c
}

func (ecs ecsLogger) Info(msg string, fields ...interface{}) {
	ecs.logger.Info(msg, ecs.toZapFields(fields)...)
}

func (ecs ecsLogger) Warn(msg string, fields ...interface{}) {
	ecs.logger.Warn(msg, ecs.toZapFields(fields)...)
}

func (ecs ecsLogger) Debug(msg string, fields ...interface{}) {
	ecs.logger.Debug(msg, ecs.toZapFields(fields)...)
}

func (ecs ecsLogger) Error(msg string, fields ...interface{}) {
	ecs.logger.Error(msg, ecs.toZapFields(fields)...)
}

func (ecs *ecsLogger) SetCorrelationId(correlationId string) {
	ecs.ctx = context.WithValue(ecs.ctx, correlationIdKey, correlationId)
}

func (ecs ecsLogger) GetCorrelationId() string {
	correlationId := ecs.ctx.Value(correlationIdKey)
	if correlationId != nil {
		return correlationId.(string)
	}
	return ""
}

func WithLogger(logger *zap.Logger) Option {
	return optionFunc(func(ecs *ecsLogger) {
		ecs.logger = logger
	})
}

func (ecs ecsLogger) toZapFields(fields []interface{}) (zapFields []zap.Field) {
	zapFields = make([]zap.Field, len(fields))
	for i, v := range fields {
		zapFields[i] = v.(zap.Field)
	}
	if correlationId := ecs.GetCorrelationId(); correlationId != "" {
		zapFields = append(zapFields, zap.Any("trace", Trace{
			ID: correlationId,
		}))
	}
	return
}

func (f optionFunc) apply(log *ecsLogger) {
	f(log)
}

func (ecs *ecsLogger) clone() *ecsLogger {
	copy := *ecs
	return &copy
}
