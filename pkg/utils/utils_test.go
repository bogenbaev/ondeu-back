package utils

import "testing"

func TestParseUint(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    uint
		wantErr bool
	}{
		{
			name: "should parse uint",
			args: args{
				s: "1",
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "should return error",
			args: args{
				s: "a",
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "should parse large uint",
			args: args{
				s: "18446744073709551615",
			},
			want:    18446744073709551615,
			wantErr: false,
		},
		{
			name: "should parse small uint",
			args: args{
				s: "0",
			},
			want:    0,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseUint(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseUint() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseUint() got = %v, want %v", got, tt.want)
			}
		})
	}
}
