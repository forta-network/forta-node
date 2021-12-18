package utils

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ipfs/go-cid"
)

// Image ref errors
var (
	ErrDiscoRefInvalid      = errors.New("invalid ipfs ref or digest in disco image ref")
	ErrDiscoRefNotIPFSCIDv1 = errors.New("image ref not an ipfs cidv1 ref")
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

// ValidateDiscoImageRef validates given image ref to be only a disco image ref.
func ValidateDiscoImageRef(discoHost, ref string) (string, error) {
	imageRef, digest := SplitImageRef(ref)
	if len(imageRef) == 0 || len(digest) == 0 {
		return "", ErrDiscoRefInvalid
	}

	// strip host from image ref
	if strings.Contains(imageRef, "/") {
		parts := strings.Split(imageRef, "/")
		imageRef = parts[1] // strip
	}

	if !isCidv1(imageRef) {
		return "", fmt.Errorf("%w: %s", ErrDiscoRefNotIPFSCIDv1, imageRef)
	}

	return fmt.Sprintf("%s/%s@sha256:%s", discoHost, imageRef, digest), nil
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
