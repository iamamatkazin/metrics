package common

import (
	"reflect"
	"testing"
)

func TestCalcSign(t *testing.T) {
	type args struct {
		key []byte
		src []byte
	}
	tests := []struct {
		name string
		want string
		args args
	}{
		{
			name: "Test 1",
			args: args{
				key: []byte("key"),
				src: []byte("src"),
			},
			want: "40ef57b8ed738f93cbd214984959d50ee2a22c8dcdc3b746e3f1eb9bdf9279a4",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CalcSign(tt.args.key, tt.args.src); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CalcSign() = %v, want %v", got, tt.want)
			}
		})
	}
}
