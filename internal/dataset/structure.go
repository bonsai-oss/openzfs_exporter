package dataset

import "sync"

type Dataset struct {
	Name      string
	ObjectID  string
	Parameter map[string]string
	Mutex     sync.Mutex
}
