package internal

import (
	"fmt"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/utils/pointer"
)

func HelmValuesToSchema(values map[string]interface{}) (*apiextensionsv1.JSONSchemaProps, error) {
	schema := &apiextensionsv1.JSONSchemaProps{
		Type:       "object",
		Properties: map[string]apiextensionsv1.JSONSchemaProps{},
	}
	for k, v := range values {
		t, err := getJSONSchema(v)
		if err != nil {
			return nil, err
		}
		schema.Properties[k] = *t
	}
	return schema, nil
}

func getJSONSchema(value interface{}) (*apiextensionsv1.JSONSchemaProps, error) {
	switch valueType := value.(type) {
	case string:
		return &apiextensionsv1.JSONSchemaProps{
			Type: "string",
		}, nil
	case int:
		return &apiextensionsv1.JSONSchemaProps{
			Type: "integer",
		}, nil
	case float64:
		return &apiextensionsv1.JSONSchemaProps{
			Type: "number",
		}, nil
	case bool:
		return &apiextensionsv1.JSONSchemaProps{
			Type: "boolean",
		}, nil
	case map[string]interface{}:
		jsonSchema := map[string]apiextensionsv1.JSONSchemaProps{}
		for k, v := range valueType {
			t, err := getJSONSchema(v)
			if err != nil {
				return nil, err
			}
			jsonSchema[k] = *t
		}
		return &apiextensionsv1.JSONSchemaProps{
			Type:                   "object",
			Properties:             jsonSchema,
			XPreserveUnknownFields: pointer.Bool(true),
		}, nil
	case []interface{}:
		v := value.([]interface{})
		var schemaV *apiextensionsv1.JSONSchemaProps
		if len(v) > 0 {
			var err error
			schemaV, err = getJSONSchema(v[0])
			if err != nil {
				return nil, err
			}
		} else {
			schemaV = &apiextensionsv1.JSONSchemaProps{
				XIntOrString: true,
			}
		}
		return &apiextensionsv1.JSONSchemaProps{
			Type: "array",
			Items: &apiextensionsv1.JSONSchemaPropsOrArray{
				Schema: schemaV,
			},
		}, nil
	default:
		return nil, fmt.Errorf("unsupported type %T found in helm chart values for %v", value, value)
	}
}
