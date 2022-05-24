package main

import (
	"encoding"
	"errors"
	"fmt"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/maps"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/mitchellh/mapstructure"
	"reflect"
	"strings"
)

const typeAndNameSeparator = "/"

type configSettings struct {
	Receivers  map[ComponentID]map[string]interface{} `mapstructure:"receivers"`
	Processors map[ComponentID]map[string]interface{} `mapstructure:"processors"`
	Exporters  map[ComponentID]map[string]interface{} `mapstructure:"exporters"`
	Extensions map[ComponentID]map[string]interface{} `mapstructure:"extensions"`
	Service    map[string]interface{}                 `mapstructure:"service"`
}

type Type string

type ComponentID struct {
	typeVal Type   `mapstructure:"-"`
	nameVal string `mapstructure:"-"`
}

var k = koanf.New(".")

func main() {
	err := k.Load(file.Provider("mock/mock.yaml"), yaml.Parser())
	if err != nil {
		return
	}
	input := maps.Unflatten(k.All(), ".")

	var rawConf = configSettings{}
	dc := decoderConfig(&rawConf)
	dc.ErrorUnused = true

	decoder, err := mapstructure.NewDecoder(dc)
	if err != nil {
		return
	}

	decoder.Decode(input)

	fmt.Println(&rawConf)

}

func decoderConfig(result interface{}) *mapstructure.DecoderConfig {
	return &mapstructure.DecoderConfig{
		Result:           result,
		Metadata:         nil,
		TagName:          "mapstructure",
		WeaklyTypedInput: true,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapKeyStringToMapKeyTextUnmarshalerHookFunc(),
			mapstructure.StringToSliceHookFunc(","),
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.TextUnmarshallerHookFunc(),
		),
	}
}

func mapKeyStringToMapKeyTextUnmarshalerHookFunc() mapstructure.DecodeHookFuncType {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {

		if f.Kind() != reflect.Map || f.Key().Kind() != reflect.String {
			return data, nil
		}

		if t.Kind() != reflect.Map {
			return data, nil
		}

		if _, ok := reflect.New(t.Key()).Interface().(encoding.TextUnmarshaler); !ok {
			return data, nil
		}

		m := reflect.MakeMap(reflect.MapOf(t.Key(), reflect.TypeOf(true)))
		for k := range data.(map[string]interface{}) {
			tKey := reflect.New(t.Key())
			if err := tKey.Interface().(encoding.TextUnmarshaler).UnmarshalText([]byte(k)); err != nil {
				return nil, err
			}

			if m.MapIndex(reflect.Indirect(tKey)).IsValid() {
				return nil, fmt.Errorf("duplicate name %q after unmarshaling %v", k, tKey)
			}
			m.SetMapIndex(reflect.Indirect(tKey), reflect.ValueOf(true))
		}
		return data, nil
	}
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (id *ComponentID) UnmarshalText(text []byte) error {
	idStr := string(text)
	items := strings.SplitN(idStr, typeAndNameSeparator, 2)
	if len(items) >= 1 {
		id.typeVal = Type(strings.TrimSpace(items[0]))
	}

	if len(items) == 1 && id.typeVal == "" {
		return errors.New("id must not be empty")
	}

	if id.typeVal == "" {
		return fmt.Errorf("in %q id: the part before %s should not be empty", idStr, typeAndNameSeparator)
	}

	if len(items) > 1 {
		// "name" part is present.
		id.nameVal = strings.TrimSpace(items[1])
		if id.nameVal == "" {
			return fmt.Errorf("in %q id: the part after %s should not be empty", idStr, typeAndNameSeparator)
		}
	}

	return nil
}
