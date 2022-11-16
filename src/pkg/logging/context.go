package logging

import (
	"context"
)

type fieldsKey struct{}

// ContextWithFields adds logger fields to fields in context
func ContextWithFields(parent context.Context, fields Fields) context.Context {
	var newFields Fields

	val := parent.Value(fieldsKey{})
	if val == nil {
		newFields = fields
	} else {
		newFields = make(Fields)
		oldFields, _ := val.(Fields)
		for k, v := range oldFields {
			newFields[k] = v
		}
		for k, v := range fields {
			newFields[k] = v
		}
	}

	return context.WithValue(parent, fieldsKey{}, newFields)
}

// ContextWithField is like ContextWithFields but adds only one field
func ContextWithField(ctx context.Context, key string, value interface{}) context.Context {
	return ContextWithFields(ctx, Fields{key: value})
}

// FieldsFromContext returns logging fields from the context
func FieldsFromContext(ctx context.Context) Fields {
	if ctx == nil {
		return nil
	}

	fields, _ := ctx.Value(fieldsKey{}).(Fields)
	return fields
}
