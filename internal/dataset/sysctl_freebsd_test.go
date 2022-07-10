//go:build freebsd
// +build freebsd

package dataset

import (
	"reflect"
	"testing"
)

var testData = "test\n" +
	"kstat.zfs.zroot.dataset.objset-0x84a1.nunlinked: 304598\n" +
	"kstat.zfs.zroot.dataset.objset-0x84a1.nunlinks: 304598\n" +
	"kstat.zfs.zroot.dataset.objset-0x84a1.nread: 278766165\n" +
	"kstat.zfs.zroot.dataset.objset-0x84a1.reads: 1606921\n" +
	"kstat.zfs.zroot.dataset.objset-0x84a1.nwritten: 52392378153\n" +
	"kstat.zfs.zroot.dataset.objset-0x84a1.writes: 1060365\n" +
	"kstat.zfs.zroot.dataset.objset-0x84a1.dataset_name: zroot/s3\n" +
	"kstat.zfs.zroot.dataset.objset-0x30b.nunlinked: 30459\n" +
	"kstat.zfs.zroot.dataset.objset-0x30b.nunlinks: 30459\n" +
	"kstat.zfs.zroot.dataset.objset-0x30b.nread: 27876616\n" +
	"kstat.zfs.zroot.dataset.objset-0x30b.reads: 160692\n" +
	"kstat.zfs.zroot.dataset.objset-0x30b.nwritten: 5239237815\n" +
	"kstat.zfs.zroot.dataset.objset-0x30b.writes: 106036\n" +
	"kstat.zfs.zroot.dataset.objset-0x30b.dataset_name: zroot/var/tmp"

func BenchmarkFindObjectDetails(b *testing.B) {
	for n := 0; n < b.N; n++ {
		findObjectDetails(testData)
	}
}

func TestFindObjectDetails(t *testing.T) {
	expectedResult := map[string]Dataset{
		"0x84a1": {
			Name: "zroot/s3",
			Parameter: map[string]string{
				"nunlinks":  "304598",
				"nunlinked": "304598",
				"nread":     "278766165",
				"reads":     "1606921",
				"nwritten":  "52392378153",
				"writes":    "1060365",
			},
		},
		"0x30b": {
			Name: "zroot/var/tmp",
			Parameter: map[string]string{
				"nunlinks":  "30459",
				"nunlinked": "30459",
				"nread":     "27876616",
				"reads":     "160692",
				"nwritten":  "5239237815",
				"writes":    "106036",
			},
		},
	}

	result := findObjectDetails(testData)
	if len(result) != len(expectedResult) {
		t.Errorf("number of result entries are not equal")
	}
	if !reflect.DeepEqual(result["0x30b"].Parameter, expectedResult["0x30b"].Parameter) {
		t.Errorf("result fields are not equal")
	}
}
