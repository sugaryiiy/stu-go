package common

import (
	"fmt"
	"reflect"
)

// FieldValue captures the name and value of a struct field.
type FieldValue struct {
	Name  string
	Value interface{}
}

// ExtractFields returns exported field names and values from the provided struct.
// It follows pointers and ignores unexported fields.
func ExtractFields(input interface{}) []FieldValue {
	v := reflect.ValueOf(input)
	if !v.IsValid() {
		return nil
	}

	// Dereference pointers
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil
	}

	t := v.Type()
	var fields []FieldValue
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		// Skip unexported fields
		if field.PkgPath != "" {
			continue
		}
		fields = append(fields, FieldValue{Name: field.Name, Value: v.Field(i).Interface()})
	}

	return fields
}

// PrintFields prints struct field names and values for quick debugging.
func PrintFields(input interface{}) {
	for _, field := range ExtractFields(input) {
		fmt.Printf("%s: %#v\n", field.Name, field.Value)
	}
}
