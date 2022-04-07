package pool

import (
	"os/exec"
	"strings"
)

type Pool struct {
	Name string
}

func Discover() (pools []Pool, err error) {
	raw, err := exec.Command("/sbin/zpool", "list", "-o", "name", "-H").Output()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(strings.TrimSuffix(string(raw), "\n"), "\n")
	for _, line := range lines {
		pools = append(pools, Pool{Name: line})
	}
	return pools, err
}
