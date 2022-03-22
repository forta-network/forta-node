package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"

	"github.com/forta-protocol/forta-core-go/encoding"
	"github.com/forta-protocol/forta-core-go/protocol"
	"github.com/forta-protocol/forta-core-go/security"
	"github.com/ipfs/go-cid"
	"github.com/spf13/cobra"
)

func handleFortaBatchDecode(cmd *cobra.Command, args []string) error {
	batchCid, err := cmd.Flags().GetString("cid")
	if err != nil {
		return err
	}
	_, err = cid.Parse(batchCid)
	if err != nil {
		return fmt.Errorf("invalid cid")
	}
	fileName, err := cmd.Flags().GetString("o")
	if err != nil {
		return err
	}
	// rest of the info lines go to stderr - useful when piping stdout to e.g. jq
	printToStdout, err := cmd.Flags().GetBool("stdout")
	if err != nil {
		return err
	}

	cmd.PrintErrln("Downloading...")

	batchResp, err := http.Get(fmt.Sprintf("%s/ipfs/%s", cfg.Publish.IPFS.GatewayURL, batchCid))
	if err != nil {
		return fmt.Errorf("failed to get batch: %v", err)
	}
	if batchResp.StatusCode != http.StatusOK {
		return fmt.Errorf("request to get batch failed with status %d", batchResp.StatusCode)
	}
	defer batchResp.Body.Close()

	cmd.PrintErrln("Successfully downloaded the batch.")

	var signedBatch protocol.SignedPayload
	if err := json.NewDecoder(batchResp.Body).Decode(&signedBatch); err != nil {
		return fmt.Errorf("failed to decode batch json: %v", err)
	}

	if err := security.VerifySignedPayload(&signedBatch); err != nil {
		yellowBold("Invalid batch signature: %v\n", err)
	} else {
		cmd.PrintErrf("Valid batch signature found - scanner: %s\n", signedBatch.Signature.Signer)
	}
	// continue decoding in any case

	var alertBatch protocol.AlertBatch
	if err := encoding.DecodeGzippedProto(signedBatch.Encoded, &alertBatch); err != nil {
		redBold("Invalid batch encoding!\n")
		return fmt.Errorf("failed to decode: %v", err)
	}

	// indent by two spaces
	b, _ := json.MarshalIndent(&alertBatch, "", "  ")

	if printToStdout {
		fmt.Println(string(b))
		return nil
	}

	dir, _ := os.Getwd()
	filePath := path.Join(dir, fileName)
	if err := os.WriteFile(filePath, b, 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %v", filePath, err)
	}
	greenBold("Successfully wrote the decoded batch to %s\n", filePath)
	return nil
}
