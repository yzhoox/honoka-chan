package banner

import "reflect"

func extractResultField(data any) any {
	if data == nil {
		return nil
	}

	value := reflect.ValueOf(data)
	if value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return nil
		}
		value = value.Elem()
	}
	if value.Kind() != reflect.Struct {
		return data
	}

	field := value.FieldByName("Result")
	if !field.IsValid() || !field.CanInterface() {
		return data
	}
	return field.Interface()
}
