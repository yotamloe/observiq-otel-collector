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

var errUnsupportedType = errors.New("unsupported type")

type ParamType string

const (
	StringType    ParamType = "string"
	IntType       ParamType = "int"
	UintType      ParamType = "uint"
	FloatType     ParamType = "float"
	BoolType      ParamType = "bool"
	MapType       ParamType = "map"
	InterfaceType ParamType = "interface{}"
	DurationType  ParamType = "duration"
	StructType    ParamType = "struct"
	ArrayType     ParamType = "array"
)

func ParamTypeDefault(paramType ParamType) interface{} {
	switch paramType {
	case StringType:
		return ""
	case IntType:
		return 0
	case UintType:
		return uint(0)
	case FloatType:
		return float64(0)
	case BoolType:
		return false
	case DurationType:
		return time.Duration(0)
	case ArrayType:
		return []interface{}{}
	default:
		return make(map[string]interface{})
	}
}

func IsBasicParamType(paramType ParamType) bool {
	switch paramType {
	case StringType, IntType, UintType, FloatType, BoolType, DurationType, ArrayType:
		return true
	default:
		return false
	}
}

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

func extractParameters(defaults map[string]interface{}) (parameters []*Paramaeter, err error) {
	parameters = make([]*Paramaeter, 0, len(defaults))
	for k, v := range defaults {
		// Weird case with for some components. Can't resolve these
		if k == "" || k == "metrics" {
			continue
		}

		param := &Paramaeter{
			Name:         k,
			DefaultValue: v,
		}

		// Look up known values
		knownParam, knownParamFound := knownParams[param.Name]
		if knownParamFound {
			param.Type = knownParam.Type
			param.Required = knownParam.Required
		}

		// Load param details
		if err := determineParamDetails(param, v); err != nil {
			if errors.Is(err, errUnsupportedType) {
				continue
			}

			return nil, err
		}

		// If something is not a map and is not a pointer but has a zero value it's required
		// We also don't want to override known params
		if param.Type != MapType && param.DefaultValue != nil && !knownParamFound {
			reflectVal := reflect.ValueOf(param.DefaultValue)
			if reflectVal.IsZero() {
				param.Required = true
				param.DefaultValue = nil
			}
		}

		parameters = append(parameters, param)
	}

	return parameters, nil
}

func determineParamDetails(param *Paramaeter, v interface{}) (err error) {
	switch i := v.(type) {
	case time.Duration, *time.Duration:
		param.Type = DurationType

		// If it's a pointer check if it's nill
		if isNil(v) {
			param.DefaultValue = nil
		}
	case map[string]interface{}:
		param.Type = MapType
		param.DefaultValue, err = processNested(i)
		if err != nil {
			return
		}
	default:
		reflectVal := reflect.ValueOf(v)
		reflectType := reflect.TypeOf(v)
		if reflectType == nil {
			// Only empty interfaces have nil reflect.Types
			param.Type = InterfaceType
			param.Required = true
			return
		}

		if reflectType.Kind() == reflect.Pointer {
			// If it's a pointer see if it's nil or not
			if reflectVal.IsNil() {
				param.DefaultValue = nil
			}
		}

		// This should unwrap any type alias or pointers
		param.Type, err = determineUnderlyingType(reflectType)
		if err != nil {
			return err
		}

		// Do post processing on special embedded types
		switch param.Type {
		case MapType:
			subMap, ok := v.(map[string]interface{})
			if ok {
				param.DefaultValue, err = processNested(subMap)
				if err != nil {
					return
				}
			}
		case StructType:
			subMap, getErr := getDefaultValues(v)
			if getErr != nil {
				return fmt.Errorf("bad struct value: %w", getErr)
			}

			param.DefaultValue, err = processNested(subMap)
			if err != nil {
				return
			}
			param.Type = MapType
		}
	}

	return
}

func isNil(v interface{}) bool {
	reflectType := reflect.TypeOf(v)
	if reflectType.Kind() != reflect.Pointer {
		return false
	}

	reflectVal := reflect.ValueOf(v)
	return reflectVal.IsNil()
}

func determineUnderlyingType(reflectType reflect.Type) (kind ParamType, err error) {
	switch reflectType.Kind() {
	case reflect.Pointer:
		kind, err = determineUnderlyingType(reflectType.Elem())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		kind = IntType
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		kind = UintType
	case reflect.Bool:
		kind = BoolType
	case reflect.Float32, reflect.Float64:
		kind = FloatType
	case reflect.String:
		kind = StringType
	case reflect.Map:
		kind = MapType
	case reflect.Struct:
		kind = StructType
	case reflect.Slice, reflect.Array:
		kind = ArrayType
	default:
		err = errUnsupportedType
	}

	return
}

func processNested(subMap map[string]interface{}) ([]*Paramaeter, error) {
	subParams, err := extractParameters(subMap)
	if err != nil {
		return nil, err
	}

	markSubParams(&subParams)

	return subParams, nil
}

// Maps are never required to we mark all their non-default value required params as not required
func markSubParams(subParams *[]*Paramaeter) {
	for _, subParam := range *subParams {
		if subParam.Required && subParam.DefaultValue == nil {
			subParam.Required = false
		}
	}
}
