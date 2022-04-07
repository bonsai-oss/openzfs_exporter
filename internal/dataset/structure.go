package dataset

var fields = []string{
	"nunlinked",
	"nunlinks",
	"nread",
	"reads",
	"nwritten",
	"writes",
}

type Dataset struct {
	Name       string
	ObjectID   string
	ObjectPath []string
	Parameter  map[string]uint64
}
