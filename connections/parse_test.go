package connections

import (
	"net"
	"reflect"
	"testing"
)

func Test_ParseEndpoint(t *testing.T) {
	type args struct {
		ep string
	}
	tests := []struct {
		name    string
		args    args
		want    net.IP
		want1   uint16
		wantErr bool
	}{
		{
			name: "10.192.1.18",
			args: args{
				ep: "1201C00A:192B",
			},
			want:    net.IPv4(10, 192, 1, 18),
			want1:   6443,
			wantErr: false,
		},
		{
			name: "10.192.1.24",
			args: args{
				ep: "1801C00A:EA78",
			},
			want:    net.IPv4(10, 192, 1, 24),
			want1:   60024,
			wantErr: false,
		},
		{
			name: "10.192.1.21",
			args: args{
				ep: "1501C00A:D8AC",
			},
			want:    net.IPv4(10, 192, 1, 21),
			want1:   55468,
			wantErr: false,
		},
		{
			name: "invalid hex",
			args: args{
				ep: "x12:D8AC",
			},
			want:    nil,
			want1:   0,
			wantErr: true,
		},
		{
			name: "empty string",
			args: args{
				ep: "",
			},
			want:    nil,
			want1:   0,
			wantErr: true,
		},
		{
			name: "too many :",
			args: args{
				ep: "AAAA:AAAA:AAAA",
			},
			want:    nil,
			want1:   0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := parseEndpoint(tt.args.ep)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseEndpoint() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseEndpoint() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("parseEndpoint() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
