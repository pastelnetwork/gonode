package cli

import (
	"reflect"
	"time"

	"github.com/urfave/cli/v2"
)

// Flag is a wrapper of cli.Flag
type Flag struct {
	cli.Flag

	destination interface{}
}

// SetUsage sets 'Usage' field for the cli.Flag
func (s *Flag) SetUsage(val string) *Flag {
	s.setField("Usage", func(v reflect.Value) {
		v.SetString(val)
	})
	return s
}

// SetAliases sets 'Aliases' field for the cli.Flag
func (s *Flag) SetAliases(val []string) *Flag {
	s.setField("Aliases", func(v reflect.Value) {
		v.Set(reflect.ValueOf(val))
	})
	return s
}

// SetEnvVars sets 'EnvVars' field for the cli.Flag
func (s *Flag) SetEnvVars(val []string) *Flag {
	s.setField("EnvVars", func(v reflect.Value) {
		v.Set(reflect.ValueOf(val))
	})
	return s
}

// SetRequired sets 'Required' field for the cli.Flag
func (s *Flag) SetRequired() *Flag {
	s.setField("Required", func(v reflect.Value) {
		v.SetBool(true)
	})
	return s
}

// SetValue sets 'Value' field for the  cli.Flag
func (s *Flag) SetValue(val interface{}) *Flag {
	s.setField("Value", func(v reflect.Value) {
		v.Set(reflect.ValueOf(val))
	})
	return s
}

// SetHidden sets 'Hidden' field for the  cli.Flag
func (s *Flag) SetHidden() *Flag {
	s.setField("Hidden", func(v reflect.Value) {
		v.SetBool(true)
	})
	return s
}

// Assign a value given from `setValue` callback func to the field of `Flag` by the given name.
func (s *Flag) setField(name string, setValue func(v reflect.Value)) reflect.Value {
	v := reflect.ValueOf(&s.Flag).Elem()
	tmp := reflect.New(v.Elem().Type()).Elem()
	tmp.Set(v.Elem())
	setValue(tmp.FieldByName(name))
	v.Set(tmp)

	return tmp.FieldByName(name)
}

// NewFlag create a new instance of the Flag struct
func NewFlag(name string, destination interface{}) *Flag {
	var flag cli.Flag

	switch ptr := destination.(type) {
	case *[]string:
		flag = &cli.StringSliceFlag{
			Name:        name,
			Destination: cli.NewStringSlice(*ptr...),
		}
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
		Flag:        flag,
		destination: destination,
	}
}
