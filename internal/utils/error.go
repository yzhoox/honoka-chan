package honokautils

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
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
	Detail    string   `json:"detail,omitempty"`
	Exception string   `json:"exception,omitempty"`
	Message   string   `json:"message,omitempty"`
	Traceback []string `json:"traceback,omitempty"`
}

func IsMainPHPRequest(path string) bool {
	return path == "/main.php" || strings.HasPrefix(path, "/main.php/")
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

func NewInternalErrorContent(err error) ErrorContent {
	return ErrorContent{
		Exception: exceptionName(err),
		Message:   err.Error(),
	}
}

func NewErrorContent(exception, message string) ErrorContent {
	return ErrorContent{
		Exception: exception,
		Message:   message,
	}
}

func NewPanicContent(recovered any, stack []byte) ErrorContent {
	return ErrorContent{
		Exception: "panic",
		Message:   fmt.Sprint(recovered),
		Traceback: splitLines(stack),
	}
}

func exceptionName(err error) string {
	if err == nil {
		return http.StatusText(http.StatusInternalServerError)
	}

	typ := reflect.TypeOf(err)
	if typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}
	if name := typ.Name(); name != "" {
		return name
	}
	return typ.String()
}

func splitLines(stack []byte) []string {
	if len(stack) == 0 {
		return nil
	}
	return strings.Split(strings.TrimSpace(string(stack)), "\n")
}
