package storage

import (
	"path"
	"strconv"
	"time"
)

// Content kinds
const (
	BloomLimit = 10000

	KindBatchReceipt = "batchReceipt"
)

// ContentLimit returns the limit for the doc count of a specific type.
func ContentLimit(kind string) int {
	switch kind {
	case KindBatchReceipt:
		return 10000

	default:
		return 10000
	}
}

// DefaultBasePath is the base path for all Forta storage content.
const DefaultBasePath = "/forta"

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

// BloomPath constructs the bloom filter path for a user.
func BloomPath(user string) string {
	return path.Join(RepoDir(user), "bloom")
}
