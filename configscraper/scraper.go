package configscraper

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/mitchellh/mapstructure"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
)

type ParamType string

const (
	stringType    ParamType = "string"
	intType       ParamType = "int"
	uintType      ParamType = "uint"
	floatType     ParamType = "float"
	boolType      ParamType = "bool"
	mapType       ParamType = "map"
	interfaceType ParamType = "interface{}"
	durationType  ParamType = "duration"
	structType    ParamType = "struct"
	arrayType     ParamType = "array"
)

type Paramaeter struct {
	Name         string      `yaml:"name"`
	Type         ParamType   `yaml:"type"`
	Required     bool        `yaml:"required"`
	DefaultValue interface{} `yaml:"default_value,omitempty"`
}

func GetConfigMeta(componenetID config.Type, factories *component.Factories) ([]*Paramaeter, error) {
	// Search Receivers
	if rcvFactory, ok := factories.Receivers[componenetID]; ok {
		return scrapeReceiverConfig(rcvFactory)
	}

	// Search Processors
	if procFactory, ok := factories.Processors[componenetID]; ok {
		return scrapeProcessorConfig(procFactory)
	}

	// Search Exporters
	if exptFactory, ok := factories.Exporters[componenetID]; ok {
		return scrapeExporterConfig(exptFactory)
	}

	// Search Extensions
	if extnFactory, ok := factories.Extensions[componenetID]; ok {
		return scrapeExtensionConfig(extnFactory)
	}

	return nil, fmt.Errorf("unknown component %s", componenetID)
}

func getDefaultValues(cfg interface{}) (map[string]interface{}, error) {
	defaults := make(map[string]interface{})
	if err := mapstructure.Decode(cfg, &defaults); err != nil {
		return nil, err
	}

	return defaults, nil
}

func extractParameters(defaults map[string]interface{}) ([]*Paramaeter, error) {
	parameters := make([]*Paramaeter, 0, len(defaults))
	for k, v := range defaults {
		// Weird case with for some components. Can't resolve these
		if k == "" {
			continue
		}
		var paramType ParamType
		defaultVal := v
		required := false
		var err error
		// switch i := v.(type) {
		// case int, int8, int16, int32, int64:
		// 	paramType = intType
		// case uint, uint8, uint16, uint32, uint64:
		// 	paramType = uintType
		// case string:
		// 	paramType = stringType
		// case bool:
		// 	paramType = boolType
		// case float32, float64:
		// 	paramType = floatType
		// case time.Duration, *time.Duration:
		// 	paramType = durationType
		// case map[string]interface{}:
		// 	paramType = mapType
		// 	defaultVal, err = extractParameters(i)
		// 	if err != nil {
		// 		return nil, err
		// 	}

		// 	// Maps are never required to we mark all their non-default value required params as not required
		// 	subParams, ok := defaultVal.([]*Paramaeter)
		// 	if !ok {
		// 		return nil, errors.New("bad map value")
		// 	}

		// 	for _, subParam := range subParams {
		// 		if subParam.Required && subParam.DefaultValue == nil {
		// 			subParam.Required = false
		// 		}
		// 	}
		// default:
		// 	reflectVal := reflect.ValueOf(v)
		// 	reflectType := reflect.TypeOf(v)
		// 	if reflectType == nil {
		// 		// Only interfaces have nil reflect.Types
		// 		paramType = interfaceType
		// 		required = true
		// 	} else {
		// 		if reflectType.Kind() == reflect.Pointer {
		// 			// If it's a pointer see if it's nil or not
		// 			if reflectVal.IsNil() {
		// 				defaultVal = nil
		// 			}
		// 		}

		// 		// This should unwrap any type alias or pointers
		// 		paramType, err = determineUnderlyingType(reflectType)
		// 		if err != nil {
		// 			fmt.Println(err)
		// 			continue
		// 		}
		// 	}
		// }

		switch i := v.(type) {
		case time.Duration, *time.Duration:
			paramType = durationType
		case map[string]interface{}:
			paramType = mapType
			defaultVal, err = extractParameters(i)
			if err != nil {
				return nil, err
			}

			// Maps are never required to we mark all their non-default value required params as not required
			subParams, ok := defaultVal.([]*Paramaeter)
			if !ok {
				return nil, errors.New("bad map value")
			}

			makeSubParams(&subParams)
		default:
			reflectVal := reflect.ValueOf(v)
			reflectType := reflect.TypeOf(v)
			if reflectType == nil {
				// Only interfaces have nil reflect.Types
				paramType = interfaceType
				required = true
			} else {
				if reflectType.Kind() == reflect.Pointer {
					// If it's a pointer see if it's nil or not
					if reflectVal.IsNil() {
						defaultVal = nil
					}
				}

				// This should unwrap any type alias or pointers
				paramType, err = determineUnderlyingType(reflectType)
				if err != nil {
					fmt.Println(err)
					continue
				}

				switch paramType {
				case mapType:
					subMap, ok := v.(map[string]interface{})
					if ok {
						defaultVal, err = extractParameters(subMap)
						if err != nil {
							return nil, err
						}

						// Maps are never required to we mark all their non-default value required params as not required
						subParams, ok := defaultVal.([]*Paramaeter)
						if !ok {
							return nil, errors.New("bad map value")
						}

						makeSubParams(&subParams)
					}
				case structType:
					subMap, err := getDefaultValues(v)
					if err != nil {
						return nil, errors.New("bad struct value")
					}

					defaultVal, err = extractParameters(subMap)
					if err != nil {
						return nil, err
					}
					paramType = mapType

					subParams, ok := defaultVal.([]*Paramaeter)
					if !ok {
						return nil, errors.New("bad struct value")
					}

					makeSubParams(&subParams)
				}
			}
		}

		// If something is not a map and is not a pointer but has a zero value it's not required
		if paramType != mapType && defaultVal != nil {
			reflectVal := reflect.ValueOf(v)
			if reflectVal.IsZero() {
				required = true
				defaultVal = nil
			}
		}

		param := &Paramaeter{
			Name:         k,
			Type:         paramType,
			Required:     required,
			DefaultValue: defaultVal,
		}

		parameters = append(parameters, param)
	}

	return parameters, nil
}

func isNull(v interface{}) bool {
	reflectVal := reflect.ValueOf(v)
	return reflectVal.IsNil()
}

func determineUnderlyingType(reflectType reflect.Type) (kind ParamType, err error) {
	switch reflectType.Kind() {
	case reflect.Pointer:
		kind, err = determineUnderlyingType(reflectType.Elem())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		kind = intType
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		kind = uintType
	case reflect.Bool:
		kind = boolType
	case reflect.Float32, reflect.Float64:
		kind = floatType
	case reflect.String:
		kind = stringType
	case reflect.Map:
		kind = mapType
	case reflect.Struct:
		kind = structType
	case reflect.Slice, reflect.Array:
		kind = arrayType
	default:
		err = fmt.Errorf("Unsupported type: %s", reflectType.Kind())
	}

	return
}

// Maps are never required to we mark all their non-default value required params as not required
func makeSubParams(subParams *[]*Paramaeter) {
	for _, subParam := range *subParams {
		if subParam.Required && subParam.DefaultValue == nil {
			subParam.Required = false
		}
	}
}
