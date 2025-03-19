package logger_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/kolzxx/html2pdf/internal/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func opts(t *testing.T, opts ...logger.Option) []logger.Option {
	t.Helper()

	return opts
}

func makeSuit(t *testing.T) (logger.Logger, *observer.ObservedLogs) {
	t.Helper()

	observedZapCore, observedLogs := observer.New(zap.DebugLevel)
	observedLogger := zap.New(observedZapCore)
	ctx := context.Background()

	opts := opts(t, logger.WithLogger(observedLogger))
	ecsLogger := logger.NewEcsLogger(ctx, opts...)

	return ecsLogger, observedLogs
}

func verifyLogAndLevel(t *testing.T, observedLogs *observer.ObservedLogs, method string, msg string, fields ...zap.Field) {
	t.Helper()

	var level zapcore.Level
	switch method {
	case "Warn":
		level = zap.WarnLevel
	case "Debug":
		level = zap.DebugLevel
	case "Error":
		level = zap.ErrorLevel
	default:
		level = zap.InfoLevel
	}

	allLogs := observedLogs.All()
	require.Equal(t, 1, observedLogs.Len())
	firstLog := allLogs[0]
	assert.Equal(t, msg, firstLog.Message)
	assert.Equal(t, level, firstLog.Level)
	assert.ElementsMatch(t, fields, firstLog.Context)
}

// @see https://medium.com/go-for-punks/handle-zap-log-messages-in-a-test-8503b25fe38f
func TestLogLevels(t *testing.T) {
	t.Parallel()

	type logLevelTestCases struct {
		Method string
	}

	for _, scenario := range []logLevelTestCases{
		{
			Method: "Info",
		},
		{
			Method: "Warn",
		},
		{
			Method: "Debug",
		},
		{
			Method: "Error",
		},
	} {

		t.Run(fmt.Sprintf("should display a message in the level log: %s", scenario.Method), func(t *testing.T) {
			// Give
			ecsLogger, observedLogs := makeSuit(t)

			// Get a reflect.Value representing the instance
			value := reflect.ValueOf(ecsLogger)

			msg := "foo"
			val := reflect.ValueOf(msg)
			args := []reflect.Value{val}

			// Get a reflect.Method representing the Foo method
			method := value.MethodByName(scenario.Method)

			// When
			method.Call(args)

			// Then
			verifyLogAndLevel(t, observedLogs, scenario.Method, msg)
		})

		t.Run(fmt.Sprintf("should display a message in the level log: %s (with correlation ID)", scenario.Method), func(t *testing.T) {
			// Give
			ecsLogger, observedLogs := makeSuit(t)
			msg := "foo"
			correlationID := "bar"

			// When
			ecsLogger.SetCorrelationId(correlationID)
			value := reflect.ValueOf(ecsLogger)
			val := reflect.ValueOf(msg)
			args := []reflect.Value{val}
			method := value.MethodByName(scenario.Method)
			method.Call(args)

			// Then
			verifyLogAndLevel(t, observedLogs, scenario.Method, msg, zap.Any("trace", logger.Trace{
				ID: correlationID,
			}))
		})

		t.Run(fmt.Sprintf("should display a message in the level log: %s (with fields)", scenario.Method), func(t *testing.T) {
			// Give
			ecsLogger, observedLogs := makeSuit(t)

			// When
			msg := "foo"
			value := reflect.ValueOf(ecsLogger)
			args := []reflect.Value{
				reflect.ValueOf(msg),
				reflect.ValueOf(zap.Field{Key: "keyOne", String: "valueOne"}),
				reflect.ValueOf(zap.Field{Key: "keyTwo", String: "valueTwo"}),
			}
			method := value.MethodByName(scenario.Method)
			method.Call(args)

			// Then
			verifyLogAndLevel(
				t,
				observedLogs,
				scenario.Method,
				msg,
				zap.Field{Key: "keyOne", String: "valueOne"},
				zap.Field{Key: "keyTwo", String: "valueTwo"},
			)
		})
	}
}

func TestCorrelationId(t *testing.T) {
	t.Parallel()

	t.Run("should set the correlation ID", func(t *testing.T) {
		// Give
		ecsLogger, _ := makeSuit(t)

		// When
		correlationID := "foo"
		ecsLogger.SetCorrelationId(correlationID)

		// Then
		assert.Equal(t, correlationID, ecsLogger.GetCorrelationId())
	})
}
