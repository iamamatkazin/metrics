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

func TestCompareIP(t *testing.T) {
	type args struct {
		ipClient  string
		ipTrusted string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "ipTrusted is empty",
			args: args{ipTrusted: ""},
			want: true,
		},
		{
			name: "ipClient is bad",
			args: args{ipTrusted: "192.168.0/24", ipClient: "qqq"},
			want: false,
		},
		{
			name: "ipClient To4 == nil",
			args: args{ipTrusted: "192.168.0/24", ipClient: "::ffff:192.168.1.228"},
			want: false,
		},
		{
			name: "ipNet.Contains == false",
			args: args{ipTrusted: "192.168.2.0/24", ipClient: "192.168.1.228"},
			want: false,
		},
		{
			name: "ipNet.Contains == true",
			args: args{ipTrusted: "192.168.1.0/24", ipClient: "192.168.1.228"},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CompareIP(tt.args.ipClient, tt.args.ipTrusted); got != tt.want {
				t.Errorf("CompareIP() = %v, want %v", got, tt.want)
			}
		})
	}
}
