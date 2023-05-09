package main

import (
	"fmt"
	"log"
	"os"

	"github.com/forta-network/forta-node/scripts/manifest/builder"
)

func main() {
	manifest, err := builder.BuildManifest(
		os.Getenv("VERSION"), os.Getenv("GITHUB_SHA"), os.Getenv("IMAGE_REF"), os.Getenv("RELEASE_NOTES"),
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(manifest)
}
