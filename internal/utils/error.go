package honokautils

import "errors"

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
