package builder

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const testTimestamp = "2020-01-01T12:00:00Z"

const testResultWithDefaultValues = `{
  "release": {
    "timestamp": "2020-01-01T12:00:00Z",
    "repository": "https://github.com/forta-network/forta-node",
    "version": "v0.5.0",
    "commit": "abc",
    "services": {
      "updater": "disco.forta.network/a@sha256:b",
      "supervisor": "disco.forta.network/a@sha256:b"
    },
    "config": {
      "autoUpdateInHours": 24,
      "deprecationPolicy": {
        "supportedVersions": [
          "v0.5.0"
        ],
        "activatesInHours": 168
      }
    }
  }
}`

const testReleaseNotes = `
# Improvements

Some description about the release.

## Foo title
Foo feature

## Bar feature
Bar feature

## Other stuff
- Baz1
- Baz2

` +

	"```yaml" +

	beginReleaseConfiguration +

	`
autoUpdateInHours: 4
deprecationPolicy:
  supportedVersions:
    - v0.4.0
    - v0.3.0
  activatesInHours: 72
` +

	endReleaseConfiguration +

	"```" +

	`
# What's Changed
- Some PR reference https://github.com/forta-network/forta-node/etc/123
`

const testResultWithCustomValues = `{
  "release": {
    "timestamp": "2020-01-01T12:00:00Z",
    "repository": "https://github.com/forta-network/forta-node",
    "version": "v0.5.0",
    "commit": "abc",
    "services": {
      "updater": "disco.forta.network/a@sha256:b",
      "supervisor": "disco.forta.network/a@sha256:b"
    },
    "config": {
      "autoUpdateInHours": 4,
      "deprecationPolicy": {
        "supportedVersions": [
          "v0.4.0",
          "v0.3.0"
        ],
        "activatesInHours": 72
      }
    }
  }
}`

func TestBuildManifest_WithDefaultValues(t *testing.T) {
	r := require.New(t)

	testTs, err := time.Parse(time.RFC3339, testTimestamp)
	r.NoError(err)

	result, err := BuildManifestWithTimestamp(testTs, "v0.5.0", "abc", "disco.forta.network/a@sha256:b", "")
	r.NoError(err)

	r.Equal(testResultWithDefaultValues, result)
}

func TestBuildManifest_WithCustomValues(t *testing.T) {
	r := require.New(t)

	testTs, err := time.Parse(time.RFC3339, testTimestamp)
	r.NoError(err)

	result, err := BuildManifestWithTimestamp(testTs, "v0.5.0", "abc", "disco.forta.network/a@sha256:b", testReleaseNotes)
	r.NoError(err)

	r.Equal(testResultWithCustomValues, result)
}
