package event

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type UnmarshalError struct {
	input string
	err   error
}

func (err UnmarshalError) Error() string {
	return fmt.Sprintf("%v, input: [%v]", err.err.Error(), err.input)
}

func Unmarshal(msg string) (event interface{}, err error) {
	refType, args, err := messageParts(msg)
	if err != nil {
		return nil, UnmarshalError{err: err, input: msg}
	}

	eventValue := reflect.New(refType).Elem()
	err = unmarshalEventArgs(&eventValue, args)
	if err != nil {
		return nil, UnmarshalError{err: err, input: msg}
	}

	return eventValue.Interface(), nil
}

func messageParts(msg string) (eventType reflect.Type, args []string, err error) {
	split := strings.Split(msg, " ")

	eventId := id(strings.ToUpper(split[0]))
	eventType, found := idToType[eventId]
	if !found {
		err = fmt.Errorf("unrecognized event id [%v]", eventId)
		return nil, nil, err
	}

	args = split[1:]
	return eventType, args, err
}

func unmarshalEventArgs(refValue *reflect.Value, args []string) error {
	if n := refValue.NumField(); len(args) < n {
		return fmt.Errorf("%v argument(s) required", n)
	}

	for i := 0; i < refValue.NumField(); i++ {
		arg := args[i]
		field := refValue.Field(i)
		switch field.Kind() {
		case reflect.Int:
			val, err := strconv.Atoi(arg)
			if err != nil {
				return fmt.Errorf("cannot read int argument [%v]", arg)
			}
			field.SetInt(int64(val))
		case reflect.String:
			field.SetString(arg)
		case reflect.Bool:
			val, err := strconv.ParseBool(arg)
			if err != nil {
				return fmt.Errorf("cannot read bool argument [%v]", arg)
			}
			field.SetBool(val)
		default:
			return fmt.Errorf("unsupported field type [%v]", field.Kind())
		}
	}

	return nil
}
