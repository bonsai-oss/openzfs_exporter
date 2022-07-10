//go:build linux
// +build linux

package dataset

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const sysFSBasePath = "/proc/spl/kstat/zfs/"

var matcher = regexp.MustCompile(`^((?P<Key>\w+)\s+\d\s+((?P<Value>\d+))|(\w+\s+\d\s+(?P<Name>\S+)))$`)

func DetectDatasets(pool string) ([]*Dataset, error) {
	var datasets []*Dataset
	filepath.Walk(sysFSBasePath+pool, func(path string, info os.FileInfo, err error) error {
		if !strings.HasPrefix(info.Name(), "objset-") {
			return nil
		}
		content, _ := os.ReadFile(path)
		ds := &Dataset{
			ObjectID:  info.Name()[len("objset-"):],
			Parameter: make(map[string]string),
		}
		for _, line := range strings.Split(string(content), "\n") {
			matchGroupContent := matcher.FindStringSubmatch(line)
			// skip if regex is not matching
			if matchGroupContent == nil {
				continue
			}
			key := matchGroupContent[matcher.SubexpIndex("Key")]
			value := matchGroupContent[matcher.SubexpIndex("Value")]
			name := matchGroupContent[matcher.SubexpIndex("Name")]

			if name != "" {
				ds.Name = name
				continue // cancel if dataset name is found. No need for further parameter/value checking
			}
			ds.Parameter[key] = value
		}
		datasets = append(datasets, ds)
		return nil
	})

	return datasets, nil
}
