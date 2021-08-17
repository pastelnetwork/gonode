package database

import (
	"fmt"
	"reflect"
	"testing"
)

func Test_safeStringQueryStruct(t *testing.T) {
	type testStruct struct {
		String_      string
		Float64_     float64
		Float64Star_ *float64
		Int_         int
		IntStart_    *int
	}
	hackString := "abc213ef ' OR 1==1 --"
	f10 := float64(10.0)
	i10 := int(10)

	processedHackString := "abc213ef '' OR 1==1 --"
	tests := []struct {
		v    interface{}
		want interface{}
	}{
		{
			v: &testStruct{
				String_:      hackString,
				Float64_:     f10,
				Float64Star_: &f10,
				Int_:         i10,
				IntStart_:    &i10,
			},
			want: &testStruct{
				String_:      processedHackString,
				Float64_:     f10,
				Float64Star_: &f10,
				Int_:         i10,
				IntStart_:    &i10,
			},
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("Test_safeStringQueryStruct-%d", i), func(t *testing.T) {
			if got := safeStringQueryStruct(tt.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("safeStringQueryStruct() = %v, want %v", got, tt.want)
			}
		})
	}
}
