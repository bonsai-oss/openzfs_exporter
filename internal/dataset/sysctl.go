package dataset

import (
	"os/exec"
	"regexp"
	"strings"
)

const (
	SYSCTL = "/sbin/sysctl"

	Parameter = "Parameter"
	Object    = "Object"
	Value     = "Value"
	Name      = "Name"
)

func DetectDatasets(pool string) (datasets []*Dataset, err error) {
	out, err := exec.Command(SYSCTL, "kstat.zfs."+pool).Output()
	if err != nil {
		return nil, err
	}

	for _, dataSetDetail := range findObjectDetails(string(out)) {
		datasets = append(datasets, dataSetDetail)
	}
	return datasets, nil
}

func findObjectDetails(input string) map[string]*Dataset {
	matcher := regexp.MustCompile(`^kstat\.zfs\.\w*\.dataset.objset-(?P<Object>\w*).((?P<Parameter>\w*): (?P<Value>\d*)|dataset_name: (?P<Name>(\w*(/)?)+))$`)
	list := make(map[string]*Dataset)

	for _, line := range strings.Split(input, "\n") {
		matchGroupContent := matcher.FindStringSubmatch(line)
		// skip if regex is not matching
		if matchGroupContent == nil {
			continue
		}
		object := matchGroupContent[matcher.SubexpIndex(Object)]
		datasetName := matchGroupContent[matcher.SubexpIndex(Name)]
		parameter := matchGroupContent[matcher.SubexpIndex(Parameter)]
		value := matchGroupContent[matcher.SubexpIndex(Value)]

		// initialize object list entry if not already existing
		if list[object] == nil {
			list[object] = &Dataset{
				ObjectID:  object,
				Parameter: make(map[string]string),
			}
		}

		// assign dataset name
		if datasetName != "" {
			list[object].Name = datasetName
			continue // cancel if dataset name is found. No need for further parameter/value checking
		}
		// assign parameters/values to dataset of not empty
		if value != "" && parameter != "" {
			list[object].Parameter[parameter] = value
		}
	}
	return list
}
