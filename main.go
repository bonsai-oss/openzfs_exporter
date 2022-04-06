package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

var fields = []string{
	"nunlinked",
	"nunlinks",
	"nread",
	"reads",
	"nwritten",
	"writes",
}

var (
	zpool_stats = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "fsrv_bsd_userland_version",
		Help: "version of the FreeBSD userland",
	})
)

type Dataset struct {
	Name       string
	ObjectID   string
	ObjectPath []string
	Values     map[string]uint64
}

type State struct {
	Datasets []Dataset
}

func (ds *Dataset) getStringValue(key string) (string, error) {
	key = strings.Join(append(ds.ObjectPath, key), ".")
	out, err := exec.Command("/sbin/sysctl", "-nq", key).Output()
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

func (ds *Dataset) parseValues() {
	for _, field := range fields {
		ds.Values[field], _ = ds.getUint64Value(field)
	}
}

func detectDatasets(pool string) []Dataset {
	validator := regexp.MustCompile(`^^kstat\.zfs\.\w*\.dataset\.objset-\w*\.dataset_name\:\s\S*`)
	var outList []Dataset
	out, err := exec.Command("/sbin/sysctl", "-it", "kstat.zfs."+pool).Output()
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if validator.MatchString(line) {
			parts := strings.Split(line, ".")
			ds := Dataset{ObjectID: parts[4], ObjectPath: parts[:5]}
			ds.Name, _ = ds.getStringValue("dataset_name")
			ds.Values = make(map[string]uint64)
			outList = append(outList, ds)
		}
	}
	return outList
}

func (st *State) RefreshDatasets(pools ...string) {
	for _, pool := range pools {
		st.Datasets = detectDatasets(pool)
	}
}

func main() {
	s := new(State)
	s.RefreshDatasets("tank")
}
