package connections

import (
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestGetConnections(t *testing.T) {
	tcp1, err := os.Open("../test/tcp2")
	if err != nil {
		t.Errorf("failed to open file: %v", err)
	}
	defer tcp1.Close()

	empty, err := os.Open("../test/tcpEmpty")
	if err != nil {
		t.Errorf("failed to open file: %v", err)
	}
	defer empty.Close()

	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]Connection
		wantErr bool
	}{
		{
			name: "tcp2",
			args: args{
				r: tcp1,
			},
			want: map[string]Connection{
				"10.192.1.18:6443:10.192.1.21:55468": {
					LocalIP:    net.ParseIP("10.192.1.18"),
					LocalPort:  6443,
					RemoteIP:   net.ParseIP("10.192.1.21"),
					RemotePort: 55468,
				},
				"10.192.1.18:6443:10.192.1.24:60024": {
					LocalIP:    net.ParseIP("10.192.1.18"),
					LocalPort:  6443,
					RemoteIP:   net.ParseIP("10.192.1.24"),
					RemotePort: 60024,
				},
			},
			wantErr: false,
		},
		{
			name: "empty",
			args: args{
				r: empty,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getConnections(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("getConnections() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getConnections() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func benchmarkLogicLoop(path string, b *testing.B) {
	log.SetOutput(ioutil.Discard)
	blocker := NewIPBlocker()
	cw := NewConnectionWatcher(blocker)
	TTL := int64(60)

	for n := 0; n < b.N; n++ {
		t := time.Now().Unix()
		cw.Observe(path, t)
		cw.Blocker.RemoveOldConnections(t, TTL)
		hosts := cw.Blocker.HostsToBlock()
		cw.Blocker.BlockHosts(hosts)
	}
}

func BenchmarkLogicLoop100000(b *testing.B) { benchmarkLogicLoop("../test/tcp100000", b) }

func benchmarkconnectionwatcherObserve(path string, b *testing.B) {
	log.SetOutput(ioutil.Discard)
	blocker := NewIPBlocker()
	cw := NewConnectionWatcher(blocker)
	t := time.Now().Unix()
	for n := 0; n < b.N; n++ {
		cw.Observe(path, t)
	}
}

func BenchmarkConnectionWatcher_Observe2(b *testing.B) {
	benchmarkconnectionwatcherObserve("../test/tcp2", b)
}
func BenchmarkConnectionWatcher_Observe100000(b *testing.B) {
	benchmarkconnectionwatcherObserve("../test/tcp100000", b)
}
