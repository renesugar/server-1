package error

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gotify/server/model"
	"github.com/stretchr/testify/assert"
)

func TestDefaultErrorInternal(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	ctx.AbortWithError(500, errors.New("something went wrong"))

	Handler()(ctx)

	assertJSONResponse(t, rec, 500, `{"errorCode":500, "errorDescription":"something went wrong", "error":"Internal Server Error"}`)
}

func TestBindingErrorDefault(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	ctx.AbortWithError(400, errors.New("you need todo something")).SetType(gin.ErrorTypeBind)

	Handler()(ctx)

	assertJSONResponse(t, rec, 400, `{"errorCode":400, "errorDescription":"you need todo something", "error":"Bad Request"}`)
}

func TestDefaultErrorBadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	ctx.AbortWithError(400, errors.New("you need todo something"))

	Handler()(ctx)

	assertJSONResponse(t, rec, 400, `{"errorCode":400, "errorDescription":"you need todo something", "error":"Bad Request"}`)
}

type testValidate struct {
	Username string `json:"username" binding:"required"`
	Mail     string `json:"mail" binding:"email"`
}

func TestValidationError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	ctx.Request = httptest.NewRequest("GET", "/uri", nil)

	assert.NotNil(t, ctx.Bind(&testValidate{}))
	Handler()(ctx)

	err := new(model.Error)
	json.NewDecoder(rec.Body).Decode(err)
	assert.Equal(t, 400, rec.Code)
	assert.Equal(t, "Bad Request", err.Error)
	assert.Equal(t, 400, err.ErrorCode)
	assert.Contains(t, err.ErrorDescription, "Field 'username' is required")
	assert.Contains(t, err.ErrorDescription, "Field 'mail' is not valid")
}

func assertJSONResponse(t *testing.T, rec *httptest.ResponseRecorder, code int, json string) {
	bytes, _ := ioutil.ReadAll(rec.Body)
	assert.Equal(t, code, rec.Code)
	assert.JSONEq(t, json, string(bytes))
}
