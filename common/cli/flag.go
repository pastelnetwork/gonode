package cli

import (
	"reflect"
	"time"

	"github.com/urfave/cli/v2"
)

// Flag is a wrapper of cli.Flag
type Flag struct {
	cli.Flag
}

// SetUsage assigns 'Usage' field for the cli.Flag
func (s *Flag) SetUsage(val string) {
	s.setField("Usage", func(v reflect.Value) {
		v.SetString(val)
	})
}

// SetAliases assigns 'Aliases' field for the cli.Flag
func (s *Flag) SetAliases(val ...string) {
	s.setField("Aliases", func(v reflect.Value) {
		v.Set(reflect.ValueOf(val))
	})
}

// SetEnvVars assigns 'EnvVars' field for the cli.Flag
func (s *Flag) SetEnvVars(val ...string) {
	s.setField("EnvVars", func(v reflect.Value) {
		v.Set(reflect.ValueOf(val))
	})
}

// SetRequired assigns 'Required' field for the cli.Flag
func (s *Flag) SetRequired() {
	s.setField("Required", func(v reflect.Value) {
		v.SetBool(true)
	})
}

// SetValue assigns 'Value' field for the cli.Flag
func (s *Flag) SetValue(val interface{}) {
	s.setField("Value", func(v reflect.Value) {
		v.Set(reflect.ValueOf(val))
	})
}

// SetDefaultText assigns 'DefaultText' field for the cli.Flag
func (s *Flag) SetDefaultText(val interface{}) {
	s.setField("DefaultText", func(v reflect.Value) {
		v.Set(reflect.ValueOf(val))
	})
}

// SetHidden assigns 'Hidden' field for the cli.Flag
func (s *Flag) SetHidden() {
	s.setField("Hidden", func(v reflect.Value) {
		v.SetBool(true)
	})
}

// Assign a value using callback func `setValue` to the field of `Flag` by the given name.
func (s *Flag) setField(name string, setValue func(v reflect.Value)) reflect.Value {
	val := reflect.ValueOf(&s.Flag).Elem()
	tmp := reflect.New(val.Elem().Type()).Elem()
	tmp.Set(val.Elem())
	setValue(tmp.Elem().FieldByName(name))
	val.Set(tmp)

	return tmp.Elem().FieldByName(name)
}

// NewFlag returns a new Flag instance. Where the `name` is a flag name and
// `destination` is a pointer to which the flag value will be assigned.
func NewFlag(name string, destination interface{}) *Flag {
	var flag cli.Flag

	switch ptr := destination.(type) {
	case *string:
		flag = &cli.StringFlag{
			Name:        name,
			Destination: ptr,
		}
	case *bool:
		flag = &cli.BoolFlag{
			Name:        name,
			Destination: ptr,
		}
	case *time.Duration:
		flag = &cli.DurationFlag{
			Name:        name,
			Destination: ptr,
		}
	case *int:
		flag = &cli.IntFlag{
			Name:        name,
			Destination: ptr,
		}
	case *uint:
		flag = &cli.UintFlag{
			Name:        name,
			Destination: ptr,
		}
	case *int64:
		flag = &cli.Int64Flag{
			Name:        name,
			Destination: ptr,
		}
	case *uint64:
		flag = &cli.Uint64Flag{
			Name:        name,
			Destination: ptr,
		}
	}

	return &Flag{
		Flag: flag,
	}
}
