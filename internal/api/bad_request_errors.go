package api

import (
	"errors"
	"fmt"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
)

type badRequestError struct {
	Property string `json:"property"`
	Reason   string `json:"reason"`
}

func collectBadRequestErrors(err error) []badRequestError {
	var result []badRequestError
	collectBadRequestErrorsInto(err, "", &result)
	return result
}

func collectBadRequestErrorsInto(err error, property string, result *[]badRequestError) {
	if err == nil {
		return
	}

	switch typedErr := err.(type) {
	case *openapi3filter.RequestError:
		if typedErr.Parameter != nil {
			property = typedErr.Parameter.Name
		}
		collectBadRequestErrorsInto(typedErr.Err, property, result)
		return
	case *openapi3.SchemaError:
		if path := strings.Join(typedErr.JSONPointer(), "."); path != "" {
			property = path
		}
		*result = append(*result, badRequestError{Property: property, Reason: schemaErrorReason(typedErr)})
		return
	case *openapi3filter.ParseError:
		if path := formatErrorPath(typedErr.Path()); path != "" {
			property = path
		}
		*result = append(*result, badRequestError{Property: property, Reason: typedErr.Error()})
		return
	case openapi3.MultiError:
		for _, nestedErr := range typedErr {
			collectBadRequestErrorsInto(nestedErr, property, result)
		}
		return
	}

	if unwrapper, ok := err.(interface{ Unwrap() []error }); ok {
		for _, nestedErr := range unwrapper.Unwrap() {
			collectBadRequestErrorsInto(nestedErr, property, result)
		}
		return
	}
	if nestedErr := errors.Unwrap(err); nestedErr != nil {
		collectBadRequestErrorsInto(nestedErr, property, result)
		return
	}

	*result = append(*result, badRequestError{Property: property, Reason: err.Error()})
}

func formatErrorPath(path []any) string {
	parts := make([]string, len(path))
	for index, part := range path {
		parts[index] = fmt.Sprintf("%v", part)
	}
	return strings.Join(parts, ".")
}

func schemaErrorReason(err *openapi3.SchemaError) string {
	if err.Origin != nil {
		if rootErr := errors.Unwrap(err.Origin); rootErr != nil {
			return rootErr.Error()
		}
		return err.Origin.Error()
	}
	if err.Reason != "" {
		return err.Reason
	}
	return err.Error()
}
