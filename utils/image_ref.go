package utils

import (
	"fmt"
	"strings"

	"github.com/ipfs/go-cid"
)

// ValidateImageRef validates an image reference. Supported formats:
//  - IPFS CID v1: bafy...@sha256:<digest>
//  - Any reference with digest: registry.hub.docker.com/library/busybox@sha256:<digest>
func ValidateImageRef(defaultRegistry, ref string) (string, bool) {
	imageRef, digest := SplitImageRef(ref)
	if len(imageRef) == 0 || len(digest) == 0 {
		return "", false
	}

	if isCidv1(imageRef) {
		return fmt.Sprintf("%s/%s", defaultRegistry, ref), true
	}

	return ref, true
}

// SplitImageRef splits the full image ref to the actual <host>/<repo> and <digest>.
func SplitImageRef(ref string) (string, string) {
	parts := strings.Split(ref, "@sha256:")
	if len(parts) == 1 {
		return "", ""
	}
	imageRef, digest := parts[0], parts[1]
	if len(digest) != 64 {
		return "", ""
	}
	return imageRef, digest
}

func isCidv1(fileCid string) bool {
	parsed, err := cid.Parse(fileCid)
	if err != nil {
		return false
	}
	return parsed.Version() == 1
}
