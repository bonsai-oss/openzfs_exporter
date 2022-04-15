package dataset

import (
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

const (
	SYSCTL = "/sbin/sysctl"
)

func (ds *Dataset) getStringValue(key string) (string, error) {
	key = strings.Join(append(ds.ObjectPath, key), ".")
	out, err := exec.Command(SYSCTL, "-nq", key).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(string(out), "\n"), nil
}

func (ds *Dataset) getUint64Value(key string) (uint64, error) {
	rawValue, err := ds.getStringValue(key)
	if err != nil {
		return 0, err
	}
	number, err := strconv.ParseUint(rawValue, 10, 64)
	if err != nil {
		return 0, err
	}
	return number, nil
}

func (ds *Dataset) ParseParameters() {
	for _, field := range fields {
		ds.Parameter[field], _ = ds.getUint64Value(field)
	}
}

func DetectDatasets(pool string) []Dataset {
	validator := regexp.MustCompile(`^^kstat\.zfs\.\w*\.dataset\.objset-\w*\.dataset_name\:\s\S*`)
	var outList []Dataset
	out, err := exec.Command(SYSCTL, "-it", "kstat.zfs."+pool).Output()
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if validator.MatchString(line) {
			parts := strings.Split(line, ".")
			ds := Dataset{ObjectID: parts[4], ObjectPath: parts[:5]}
			ds.Name, _ = ds.getStringValue("dataset_name")
			ds.Parameter = make(map[string]uint64)
			outList = append(outList, ds)
		}
	}
	return outList
}
