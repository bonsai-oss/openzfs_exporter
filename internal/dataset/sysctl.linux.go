//go:build linux
// +build linux

package dataset

const sysFSBasePath = "/proc/spl/kstat/zfs/"

func DetectDatasets(pool string) ([]*Dataset, error) {
	return nil, nil
}
