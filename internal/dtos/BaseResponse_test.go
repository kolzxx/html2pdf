package dtos_test

import (
	"testing"

	"github.com/kolzxx/html2pdf/internal/dtos"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type baseResponseContractTest struct {
	suite.Suite
}

func Test_BaseResponseContract(t *testing.T) {
	t.Parallel()

	t.Run("Test_BaseResponseContract_WithSuccess_ReturnsExpectedResponse", func(t *testing.T) {
		message := "Success"
		code := 200
		result := "Result Data"

		response := dtos.WithSuccess(message, code, result)

		assert.Equal(t, message, response.Message)
		assert.Equal(t, true, response.Success)
		assert.Equal(t, code, response.Code)
		assert.Equal(t, result, response.Result)
		assert.Nil(t, response.Errors)
	})

	t.Run("Test_BaseResponseContract_WithError_ReturnsExpectedResponse", func(t *testing.T) {
		message := "Error"
		code := 400
		errors := []dtos.Error{{Title: "title1", Detail: "detail1"}, {Title: "title2", Detail: "detail2"}}

		response := dtos.WithError(message, code, errors...)

		assert.Equal(t, message, response.Message)
		assert.Equal(t, false, response.Success)
		assert.Equal(t, code, response.Code)
		assert.Equal(t, errors, response.Errors)
		assert.Nil(t, response.Result)
	})

	t.Run("Test_BaseResponseContract_WithError_NoErrors_ReturnsExpectedResponse", func(t *testing.T) {
		message := "Error"
		code := 400

		response := dtos.WithError(message, code)

		assert.Equal(t, message, response.Message)
		assert.Equal(t, false, response.Success)
		assert.Equal(t, code, response.Code)
		assert.Empty(t, response.Errors)
		assert.Nil(t, response.Result)
	})
}
