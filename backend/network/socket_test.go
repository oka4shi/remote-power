package network

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestListPhisicalIfaces(t *testing.T) {
	cases := map[string]struct {
		in   []string
		want []string
	}{
		"物理デバイスが認識できるか確認": {
			in: []string{
				"/sys/devices/pci0000:00/0000:00:02.1/0000:03:00.0/0000:04:0b.0/0000:0d:00.0/net/wlp13s0",
				"/sys/devices/pci0000:00/0000:00:02.1/0000:03:00.0/0000:04:0a.0/0000:0c:00.0/net/enp12s0",
				"/sys/devices/platform/fe2a0000.ethernet/net/eth0",
				"/sys/devices/platform/3c0000000.pcie/pci0000:00/0000:00:00.0/0000:01:00.0/net/eth1",
			},
			want: []string{
				"enp12s0",
				"eth0",
				"eth1",
				"wlp13s0",
			},
		},
		"仮想デバイスが除外されているか確認": {
			in: []string{
				"/sys/devices/virtual/net/br0",
				"/sys/devices/virtual/net/tailscale0",
			},
			want: []string{},
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			d := t.TempDir()
			t.Logf("TempDir: %s", d)

			err := makeDirs(t, d, tt.in)
			if err != nil {
			}

			got, err := ListPhisicalIfaces(d)
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(tt.want, got) {
				fmt.Println("want: ", tt.want)
				fmt.Println("data: ", got)
				t.Fatal("Unexpected data")
			}
		})
	}
}

func makeDirs(t *testing.T, root string, paths []string) error {
	t.Helper()

	for _, path := range paths {
		if err := os.MkdirAll(filepath.Join(root, path), 0755); err != nil {
			t.Error("Failed to arrange a test: ", err)
		}
	}
	return nil
}
