package main

/*
#include <stdlib.h>
*/
import "C"

import (
	"context"
	"encoding/json"
	"honoka-chan/internal/app"
	systemhandler "honoka-chan/internal/handler/system"
	"time"
	"unsafe"
)

//export ServerStart
func ServerStart(workDir *C.char) *C.char {
	if err := app.Start(C.GoString(workDir)); err != nil {
		return C.CString(err.Error())
	}
	return nil
}

//export ServerStop
func ServerStop() *C.char {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.Stop(ctx); err != nil {
		return C.CString(err.Error())
	}
	return nil
}

//export ServerStatusJSON
func ServerStatusJSON() *C.char {
	return C.CString(app.GetStatusJSON())
}

//export ServerHealthJSON
func ServerHealthJSON() *C.char {
	return C.CString(systemhandler.HealthJSON())
}

//export ServerReload
func ServerReload() *C.char {
	resp, _, err := systemhandler.Reload("", false)
	if err == nil {
		return nil
	}

	data, marshalErr := json.Marshal(resp)
	if marshalErr != nil {
		return C.CString(err.Error())
	}
	return C.CString(string(data))
}

//export ServerFreeString
func ServerFreeString(str *C.char) {
	if str == nil {
		return
	}
	C.free(unsafe.Pointer(str))
}

func main() {}
