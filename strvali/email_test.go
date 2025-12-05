package strvali

import (
	"testing"
)

func TestIsValidEmail(t *testing.T) {
	type args struct {
		email string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "standard qq email",
			args: args{
				"xxxxx@qq.com",
			},
			want: true,
		},
		{
			name: "long TLD (valid)",
			args: args{
				"user@example.northwesternmutual", // 18 chars
			},
			want: true,
		},
		{
			name: "underscore in domain (valid)",
			args: args{
				"user@sub_domain.example.com",
			},
			want: true,
		},
		{
			name: "consecutive dots in local part (invalid)",
			args: args{
				"user..name@example.com",
			},
			want: false,
		},
		{
			name: "start with dot (invalid)",
			args: args{
				".user@example.com",
			},
			want: false,
		},
		{
			name: "end with dot in local part (invalid)",
			args: args{
				"user.@example.com",
			},
			want: false,
		},
		{
			name: "plus sign in local part (valid)",
			args: args{
				"user+tag@example.com",
			},
			want: true,
		},
		{
			name: "domain with hyphen (valid)",
			args: args{
				"user@my-domain.com",
			},
			want: true,
		},
		{
			name: "domain ends with hyphen (invalid)",
			args: args{
				"user@example-.com",
			},
			want: false,
		},
		{
			name: "domain ends with underscore (invalid)",
			args: args{
				"user@example_.com",
			},
			want: false,
		},
		{
			name: "email too long (invalid)",
			args: args{
				"user@" + string(make([]byte, 250)) + ".com",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidEmail(tt.args.email); got != tt.want {
				t.Errorf("IsValidEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}
