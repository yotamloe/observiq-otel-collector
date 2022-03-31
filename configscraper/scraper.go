package configscraper

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/mitchellh/mapstructure"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/config/configunmarshaler"
	"gopkg.in/yaml.v3"
)

type ParamType string

const (
	stringType   ParamType = "string"
	intType      ParamType = "int"
	uintType     ParamType = "uint"
	floatType    ParamType = "float"
	boolType     ParamType = "bool"
	mapType      ParamType = "map"
	durationType ParamType = "duration"
)

type Paramaeter struct {
	Name         string      `yaml:"name"`
	Type         ParamType   `yaml:"type"`
	Required     bool        `yaml:"required"`
	DefaultValue interface{} `yaml:"default_value,omitempty"`
}

func ScrapeReceiverConfig(factory component.ReceiverFactory) (*string, error) {
	// Load up default config
	cfg, err := configunmarshaler.LoadReceiver(config.NewMap(), factory.CreateDefaultConfig().ID(), factory)
	if err != nil {
		return nil, err
	}

	// Get default values as a map we can lookup
	defaults, err := getDefaultValues(cfg)
	if err != nil {
		return nil, err
	}

	parameters, err := extractParameters(defaults)
	if err != nil {
		return nil, err
	}

	//Yaml out
	data, err := yaml.Marshal(&parameters)
	if err != nil {
		return nil, err
	}
	out := string(data)

	return &out, nil
}

func getDefaultValues(cfg config.Receiver) (map[string]interface{}, error) {
	defaults := make(map[string]interface{})
	if err := mapstructure.Decode(cfg, &defaults); err != nil {
		return nil, err
	}

	return defaults, nil
}

func extractParameters(defaults map[string]interface{}) ([]*Paramaeter, error) {
	parameters := make([]*Paramaeter, 0, len(defaults))
	for k, v := range defaults {
		var paramType ParamType
		defaultVal := v
		required := false
		var err error
		switch i := v.(type) {
		case int, int8, int16, int32, int64:
			paramType = intType
		case uint, uint8, uint16, uint32, uint64:
			paramType = uintType
		case string:
			paramType = stringType
		case bool:
			paramType = boolType
		case float32, float64:
			paramType = floatType
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

			for _, subParam := range subParams {
				if subParam.Required && subParam.DefaultValue == nil {
					subParam.Required = false
				}
			}
		default:
			reflectVal := reflect.ValueOf(v)
			if reflectVal.Type().Kind() == reflect.Pointer {
				if reflectVal.IsNil() {
					defaultVal = nil
				}

			}

			// This should unwrap any type alias or pointers
			paramType, err = determineUnderlyingType(reflectVal.Type())
			if err != nil {
				fmt.Println(err)
				continue
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
	default:
		err = fmt.Errorf("Unsupported type: %s", reflectType.Kind())
	}

	return
}
