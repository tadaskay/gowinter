package event

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type MarshalError struct {
	input interface{}
}

func (err MarshalError) Error() string {
	return fmt.Sprintf("input: [%v]", err.input)
}

func Marshal(event interface{}) (msg string, err error) {
	refType := reflect.TypeOf(event)
	var eventId, found = typeToId(refType)
	if !found {
		return "", MarshalError{event}
	}
	refValue := reflect.ValueOf(event)

	messageParts := []string{string(eventId)}
	for i := 0; i < refValue.NumField(); i++ {
		field := refValue.Field(i)
		var val string
		switch field.Kind() {
		case reflect.Int:
			val = strconv.FormatInt(field.Int(), 10)
		case reflect.Bool:
			val = strconv.FormatBool(field.Bool())
		default:
			val = field.String()
		}
		messageParts = append(messageParts, val)
	}

	return strings.Join(messageParts, " "), nil
}
