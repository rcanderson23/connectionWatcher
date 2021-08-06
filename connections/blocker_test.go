package connections

import (
	"reflect"
	"testing"
)

// shortcut: Test is flaky due to no guarantee on ordering in the slice. Helper functions to order the slices
// of RemoteHost and Ports should fix this and allow for accurate testing with reflect.DeepEqual.
//func TestIPBlocker_HostsToBlock(t *testing.T) {
//	type fields struct {
//		IPPortTime map[string]map[uint16]int64
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		want   []RemoteHost
//	}{
//		{
//			name: "block 192.168.1.1",
//			fields: fields{IPPortTime: map[string]map[uint16]int64{
//				"192.168.1.1": {80: 10, 81: 20, 82: 30},
//			}},
//			want: []RemoteHost{{
//				RemoteIP: net.ParseIP("192.168.1.1"),
//				Ports: []uint16{80, 81, 82},
//				},
//			},
//		},
//		{
//			name: "block 192.168.1.1 and 192.168.1.2",
//			fields: fields{IPPortTime: map[string]map[uint16]int64{
//				"192.168.1.1": {80: 10, 81: 20, 82: 30},
//				"192.168.1.2": {80: 10, 81: 20, 82: 30},
//			}},
//			want: []RemoteHost{
//				{
//					RemoteIP: net.ParseIP("192.168.1.1"),
//					Ports: []uint16{80, 81, 82},
//				},
//				{
//					RemoteIP: net.ParseIP("192.168.1.2"),
//					Ports: []uint16{80, 81, 82},
//				},
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			ipb := &IPBlocker{
//				IPPortTime: tt.fields.IPPortTime,
//			}
//			if got := ipb.HostsToBlock(); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("HostsToBlock() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

func TestIPBlocker_RemoveOldConnections(t *testing.T) {
	type fields struct {
		IPPortTime map[string]map[uint16]int64
	}
	type args struct {
		now int64
		ttl int64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []uint16
	}{
		{
			name: "remove all",
			fields: fields{IPPortTime: map[string]map[uint16]int64{
				"192.168.1.1": {80: 40, 81: 50, 82: 60},
			}},
			args: args{
				now: 90,
				ttl: 30,
			},
			want: []uint16{80, 81, 82},
		},
		{
			name: "remove 80",
			fields: fields{IPPortTime: map[string]map[uint16]int64{
				"192.168.1.1": {80: 40, 81: 50, 82: 60},
			}},
			args: args{
				now: 70,
				ttl: 30,
			},
			want: []uint16{80},
		},
		{
			name: "remove none",
			fields: fields{IPPortTime: map[string]map[uint16]int64{
				"192.168.1.1": {80: 40, 81: 50, 82: 60},
			}},
			args: args{
				now: 60,
				ttl: 30,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ipb := &IPBlocker{
				IPPortTime: tt.fields.IPPortTime,
			}

			if got := ipb.RemoveOldConnections(tt.args.now, tt.args.ttl); !containsSamePorts(got, tt.want) {
				t.Errorf("RemoveOldConnections() = %v, want %v", got, tt.want)
			}
		})
	}
}

func containsSamePorts(a, b []uint16) bool {
	if len(a) != len(b) {
		return false
	}

	present := func(item uint16, slice []uint16) bool {
		for _, s := range slice {
			if s == item {
				return true
			}
		}
		return false
	}

	for _, i := range b {
		if !present(i, a) {
			return false
		}
	}
	return true
}

func TestNewIPBlocker(t *testing.T) {
	tests := []struct {
		name string
		want *IPBlocker
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewIPBlocker(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewIPBlocker() = %v, want %v", got, tt.want)
			}
		})
	}
}
