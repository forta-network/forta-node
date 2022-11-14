package storage

import (
	"path"
	"strconv"
	"time"
)

// Content constants
const (
	BloomLimit             = 10000
	BloomFalsePositiveRate = 0.0001

	KindBatchReceipt = "batchReceipt"
)

// Globals
var (
	HistorySupport = time.Hour * 48
	BucketInterval = time.Minute * 30
	MaxBuckets     = int(HistorySupport / BucketInterval)
)

// DefaultBasePath is the base path for all Forta storage content.
const DefaultBasePath = "/forta"

// RepoDir constructs the repository dir path for a user.
func RepoDir(user string) string {
	return path.Join(DefaultBasePath, user)
}

// ContentDir constructs the dir path for a specific kind of content.
func ContentDir(user, kind string) string {
	return path.Join(RepoDir(user), kind)
}

// BucketDir constructs the dir path for a specific kind of content bucket.
func BucketDir(user, kind, bucket string) string {
	return path.Join(RepoDir(user), kind, bucket)
}

// NewContentPath creates the full path for a specific kind of content.
func NewContentPath(user, kind string) (contentPath, bucketDir string) {
	now := time.Now()

	bucketTs := now.Truncate(BucketInterval).UnixNano()
	bucketTsStr := strconv.FormatInt(bucketTs, 10)
	bucketDir = BucketDir(user, kind, bucketTsStr)

	contentTs := now.UnixNano()
	contentTsStr := strconv.FormatInt(contentTs, 10)
	contentPath = path.Join(bucketDir, contentTsStr)

	return
}

// BloomPath constructs the bloom filter path for a user.
func BloomPath(user string) string {
	return path.Join(RepoDir(user), "bloom")
}
