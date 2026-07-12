package honokautils

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type UnimplementedError struct {
	Kind   string
	Module string
	Name   string
}

func (e *UnimplementedError) Error() string {
	if e == nil {
		return "unimplemented"
	}
	return "unimplemented " + e.Kind + ": " + e.Module + ": " + e.Name
}

func NewUnimplementedModuleError(module string) error {
	return &UnimplementedError{
		Kind:   "api module",
		Module: module,
		Name:   module,
	}
}

func NewUnimplementedActionError(module string, action string) error {
	return &UnimplementedError{
		Kind:   "action",
		Module: module,
		Name:   action,
	}
}

func IsUnimplementedError(err error) bool {
	var target *UnimplementedError
	return err != nil && errors.As(err, &target)
}

const maintenanceHeader = "Maintenance"

type ErrorContent struct {
	Detail    string `json:"detail,omitempty"`
	Exception string `json:"exception,omitempty"`
	Message   string `json:"message,omitempty"`
}

func IsMainPHPRequest(path string) bool {
	return path == "/main.php" || strings.HasPrefix(path, "/main.php/")
}

func IsMainLoginEndpoint(path string) bool {
	return path == "/main.php/login/authkey" || path == "/main.php/login/login"
}

func WriteMaintenanceJSON(setHeader func(string, string), writeJSON func(int, any), status int, content any) {
	setHeader("Content-Type", "application/json; charset=utf-8")
	setHeader("status_code", strconv.Itoa(status))
	setHeader(maintenanceHeader, "1")
	writeJSON(status, content)
}

func AbortMaintenanceJSON(ctx *gin.Context, status int, content any) {
	WriteMaintenanceJSON(ctx.Header, ctx.JSON, status, content)
	ctx.Abort()
}

func NewDetailContent(detail string) ErrorContent {
	return ErrorContent{
		Detail: detail,
	}
}

func NewNotFoundContent(path string) ErrorContent {
	return NewDetailContent(fmt.Sprintf("Endpoint not found: %s", path))
}

func NewInternalErrorContent() ErrorContent {
	return ErrorContent{
		Exception: "InternalServerError",
		Message:   http.StatusText(http.StatusInternalServerError),
	}
}

func NewPanicContent() ErrorContent {
	return NewInternalErrorContent()
}
