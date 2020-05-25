package version

import "fmt"

var (
	Version = ""
	GitHash = ""
	Created = ""
	Release = ""
)

func init() {
	Release = fmt.Sprintf("%s-%s-%s", Version, GitHash, Created)
}
