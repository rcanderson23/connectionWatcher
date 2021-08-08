package connections

import "testing"

func TestIsEphemeralPort(t *testing.T) {
	type args struct {
		port uint16
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "60000",
			args: args{
				port: uint16(60000),
			},
			want: true,
		},
		{
			name: "80",
			args: args{
				port: uint16(80),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsEphemeralPort(tt.args.port); got != tt.want {
				t.Errorf("IsEphemeralPort() = %v, want %v", got, tt.want)
			}
		})
	}
}
