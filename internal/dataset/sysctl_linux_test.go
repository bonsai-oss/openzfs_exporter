//go:build linux
// +build linux

package dataset

import (
	"reflect"
	"testing"
)

const testData = "test\n" +
	"38 1 0x01 7 2160 5326301430 8579547267790\n" +
	"name                            type data\n" +
	"dataset_name                    7    rpool/ROOT/ubuntu_2bfmnx/var/lib/NetworkManager\n" +
	"writes                          4    34\n" +
	"nwritten                        4    2054\n" +
	"reads                           4    5\n" +
	"nread                           4    1175\n" +
	"nunlinks                        4    34\n" +
	"nunlinked                       4    34\n"

func TestParseDatasetValues(t *testing.T) {
	expectedResult := Dataset{
		Name: "rpool/ROOT/ubuntu_2bfmnx/var/lib/NetworkManager",
		Parameter: map[string]string{
			"writes":    "34",
			"nwritten":  "2054",
			"reads":     "5",
			"nread":     "1175",
			"nunlinks":  "34",
			"nunlinked": "34",
		},
		ObjectID: "0x01",
	}

	testDataset := &Dataset{ObjectID: "0x01", Name: "rpool/ROOT/ubuntu_2bfmnx/var/lib/NetworkManager", Parameter: map[string]string{}}

	parseDatasetValues(testDataset, []byte(testData))
	if !reflect.DeepEqual(*testDataset, expectedResult) {
		t.Errorf("Expected %#v,\n   got %#v", expectedResult, *testDataset)
	}

}
