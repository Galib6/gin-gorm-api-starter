package base

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var basicErrorMessages = map[string]string{
	"excluded_unless": "The value of '%s' must be empty.",
	"alphanum":        "The value of '%s' must contain only letters and numbers, without spaces or symbols.",
	"boolean":         "The value of '%s' must be either true or false.",
	"email":           "The value of '%s' must be a valid email address.",
	"url":             "The value of '%s' must be a valid URL.",
	"uuid":            "The value of '%s' must be a valid UUID.",
	"ip":              "The value of '%s' must be a valid IP address.",
	"alpha":           "The value of '%s' must contain only letters without numbers or symbols.",
	"alphaunicode":    "The value of '%s' must contain only Unicode letters without numbers or symbols.",
	"ascii":           "The value of '%s' must contain only ASCII characters.",
	"base64":          "The value of '%s' must be a valid Base64 string.",
}

var paramErrorMessages = map[string]string{
	"len":      "The value of '%s' must be exactly %s characters long.",
	"max":      "The value of '%s' must be at most %s.",
	"min":      "The value of '%s' must be at least %s.",
	"gte":      "The value of '%s' must be greater than or equal to %s.",
	"lte":      "The value of '%s' must be less than or equal to %s.",
	"gt":       "The value of '%s' must be greater than %s.",
	"lt":       "The value of '%s' must be less than %s.",
	"datetime": "The value of '%s' must be in a valid date-time format (example: %s).",
	"oneof":    "The value of '%s' must be one of the following: %s.",
	"unique":   "The value of '%s' which is '%s' must be unique and not duplicate any other value.",
}

func getErrorMessage(fieldName string, fe validator.FieldError) string {
	if msgTemplate, exists := basicErrorMessages[fe.Tag()]; exists {
		return fmt.Sprintf(msgTemplate, fieldName)
	} else if msgTemplate, exists := paramErrorMessages[fe.Tag()]; exists {
		return fmt.Sprintf(msgTemplate, fieldName, fe.Param())
	}

	switch fe.Tag() {
	case "required", "required_if", "required_with":
		return "The value of '" + fieldName + "' is required."
	case "numeric", "number":
		return "The value of '" + fieldName + "' must contain only numbers."
	default:
		return "The value of '" + fieldName + "' is invalid."
	}
}

func GetValidationErrorMessage(err error, reqStruct interface{}, defaultMsg string) string {
	valMsgs := FormatValidationErrors(err, reqStruct)
	if len(valMsgs) > 0 {
		return valMsgs[0]
	}
	return defaultMsg
}

func FormatValidationErrors(err error, reqStruct interface{}, customAdditionalErrs ...string) []string {
	jsonTagMap := getTagValueFromStruct(reqStruct, "json")
	fieldNameMap := getTagValueFromStruct(reqStruct, "field")

	messages := []string{}
	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, fieldErr := range validationErrs {
			fieldPath := fieldErr.Namespace()
			if idx := strings.Index(fieldPath, "."); idx != -1 {
				fieldPath = fieldPath[idx+1:]
			}

			cleanFieldPath := removeBracketsAndContent(fieldPath)

			jsonPath := convertStructPathToJSON(fieldPath, jsonTagMap)
			fieldName, exists := fieldNameMap[cleanFieldPath]
			if !exists {
				fieldName = strings.ToLower(jsonPath)
			} else {
				fieldName = strings.ToLower(fieldName)
			}

			message := getErrorMessage(fieldName, fieldErr)
			messages = append(messages, message)
		}
	} else if err != nil {
		return []string{err.Error()}
	}

	messages = append(messages, customAdditionalErrs...)
	return messages
}

func CreateValidationErrorMessage(errMsg string) []string {
	var errors []string
	errors = append(errors, errMsg)
	return errors
}

func convertStructPathToJSON(structPath string, jsonTagMap map[string]string) string {
	parts := strings.Split(structPath, ".")
	jsonPathParts := []string{}

	for i := range parts {
		fullKey := strings.Join(parts[:i+1], ".")

		if jsonTag, exists := jsonTagMap[fullKey]; exists {
			jsonPathParts = append(jsonPathParts, jsonTag)
		} else {
			jsonPathParts = append(jsonPathParts, parts[i])
		}
	}

	return strings.ToLower(strings.Join(jsonPathParts, "."))
}

func getTagValueFromStruct(reqStruct interface{}, tagName string) map[string]string {
	tagValueMap := make(map[string]string)
	parseStruct(reflect.TypeOf(reqStruct), tagName, tagValueMap, "")
	return tagValueMap
}

func parseStruct(t reflect.Type, tagName string, tagMap map[string]string, parent string) {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldType := field.Type
		tagValue := field.Tag.Get(tagName)

		fullFieldName := field.Name
		if parent != "" {
			fullFieldName = parent + "." + field.Name
		}

		if tagValue != "" {
			tagMap[fullFieldName] = tagValue
		}

		if fieldType.Kind() == reflect.Slice && fieldType.Elem().Kind() == reflect.Struct {
			parseStruct(fieldType.Elem(), tagName, tagMap, fullFieldName)
		}

		if fieldType.Kind() == reflect.Struct {
			parseStruct(fieldType, tagName, tagMap, fullFieldName)
		}
	}
}

func removeBracketsAndContent(input string) string {
	// Regular expression to match '[' followed by anything, then ']'
	re := regexp.MustCompile(`\[.*?\]`)
	// Replace the matched content with an empty string
	return re.ReplaceAllString(input, "")
}

func FinalizeErrorMessage(msg string) string {
	if msg == "" {
		return msg
	}

	msg = strings.ToUpper(msg[:1]) + msg[1:]

	if !strings.HasSuffix(msg, ".") {
		msg += "."
	}

	return msg
}
