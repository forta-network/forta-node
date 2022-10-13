package content

import (
	"path"
	"strconv"
	"time"
)

// DefaultBasePath is the base path for all Forta storage content.
const DefaultBasePath = "/forta_storage"

// RepoDir constructs the repository dir path for a user.
func RepoDir(user string) string {
	return path.Join(DefaultBasePath, user)
}

// ContentDir constructs the dir path for a specific kind of content.
func ContentDir(user string, kind string) string {
	return path.Join(RepoDir(user), kind)
}

// NewContentPath creates the full path for a specific kind of content.
func NewContentPath(user string, kind string) string {
	ts := time.Now().UnixNano()
	tsStr := strconv.FormatInt(ts, 10)
	return path.Join(ContentDir(user, kind), tsStr)
}
