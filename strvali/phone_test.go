package strvali

import (
	"testing"
)

func TestIsValidPhoneNum(t *testing.T) {
	type args struct {
		area  string
		phone string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "cn",
			args: args{
				area:  "+86",
				phone: "13890101000",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidPhoneNum(tt.args.area, tt.args.phone); got != tt.want {
				t.Errorf("IsValidPhoneNum() = %v, want %v", got, tt.want)
			}
		})
	}
}
