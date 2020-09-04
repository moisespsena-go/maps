package maps

import (
	"reflect"

	"github.com/mitchellh/mapstructure"
	"github.com/moisespsena-go/xbindata"
)

var UnmarshalerType = reflect.TypeOf((*MapUnmarshaler)(nil)).Elem()

type MapUnmarshaler interface {
	UnmarshalMap(value interface{}) (err error)
}

// A DecoderConfigOption can be passed to viper.Unmarshal to configure
// mapstructure.DecoderConfig options
type DecoderConfigOption func(*mapstructure.DecoderConfig)

// DecodeHook returns a DecoderConfigOption which overrides the default
// DecoderConfig.DecodeHook value, the default is:
//
//  mapstructure.ComposeDecodeHookFunc(
//		mapstructure.StringToTimeDurationHookFunc(),
//		mapstructure.StringToSliceHookFunc(","),
//	)
func DecodeHook(hook mapstructure.DecodeHookFunc) DecoderConfigOption {
	return func(c *mapstructure.DecoderConfig) {
		c.DecodeHook = hook
	}
}

func DecoderConfig(config *mapstructure.DecoderConfig) {
	oldHook := config.DecodeHook
	config.DecodeHook = mapstructure.ComposeDecodeHookFunc(oldHook, func(from reflect.Type, to reflect.Type, v interface{}) (interface{}, error) {
		if to.Kind() == reflect.Struct {
			if reflect.PtrTo(to).Implements(UnmarshalerType) {
				unmh := reflect.New(to).Interface().(xbindata.MapUnmarshaler)
				if err := unmh.UnmarshalMap(v); err != nil {
					return nil, err
				}
				return unmh, nil
			}
		}
		return v, nil
	})
}

// Unmarshal unmarshals the config into a Struct. Make sure that the tags
// on the fields of the structure are properly set.
func (this MapSI) Unmarshal(rawVal interface{}, opts ...DecoderConfigOption) error {
	err := decode(rawVal, defaultDecoderConfig(&this, opts...))

	if err != nil {
		return err
	}

	return nil
}

// CopyTo copy the config into a Struct. Make sure that the tags
// on the fields of the structure are properly set.
func (this MapSI) CopyTo(dest interface{}, opts ...DecoderConfigOption) error {
	err := decode(this, defaultDecoderConfig(dest, opts...))

	if err != nil {
		return err
	}

	return nil
}

// defaultDecoderConfig returns default mapsstructure.DecoderConfig with suppot
// of time.Duration values & string slices
func defaultDecoderConfig(output interface{}, opts ...DecoderConfigOption) *mapstructure.DecoderConfig {
	c := &mapstructure.DecoderConfig{
		Metadata:         nil,
		Result:           output,
		WeaklyTypedInput: true,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// A wrapper around mapstructure.Decode that mimics the WeakDecode functionality
func decode(input interface{}, config *mapstructure.DecoderConfig) error {
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}
	return decoder.Decode(input)
}
