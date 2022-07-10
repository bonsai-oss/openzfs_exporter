//go:build linux
// +build linux

package dataset

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	sysFSBasePath = "/proc/spl/kstat/zfs/"
	objsetPrefix  = "objset-"

	matchGroupKey   = "Key"
	matchGroupValue = "Value"
	matchGroupName  = "Name"
)

var matcher = regexp.MustCompile(`^((?P<Key>\w+)\s+\d\s+((?P<Value>\d+))|(\w+\s+\d\s+(?P<Name>\S+)))$`)

func DetectDatasets(pool string) (datasets []*Dataset, err error) {
	return datasets, filepath.Walk(sysFSBasePath+pool, func(path string, info os.FileInfo, err error) error {
		if !strings.HasPrefix(info.Name(), objsetPrefix) {
			return nil
		}
		content, fileReadError := os.ReadFile(path)
		if fileReadError != nil {
			return nil
		}
		ds := &Dataset{
			ObjectID:  strings.TrimPrefix(info.Name(), objsetPrefix),
			Parameter: make(map[string]string),
		}
		parseDatasetValues(ds, content)
		datasets = append(datasets, ds)
		return nil
	})
}

func parseDatasetValues(ds *Dataset, content []byte) {
	for _, line := range strings.Split(string(content), "\n") {
		matchGroupContent := matcher.FindStringSubmatch(line)
		// skip if regex is not matching
		if matchGroupContent == nil {
			continue
		}
		key := matchGroupContent[matcher.SubexpIndex(matchGroupKey)]
		value := matchGroupContent[matcher.SubexpIndex(matchGroupValue)]
		name := matchGroupContent[matcher.SubexpIndex(matchGroupName)]

		if name != "" {
			ds.Name = name
			continue // cancel if dataset name is found. No need for further parameter/value checking
		}
		ds.Parameter[key] = value
	}
}
