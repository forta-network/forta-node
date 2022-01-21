package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/fatih/color"
	"github.com/forta-protocol/forta-node/clients/health"
	"github.com/forta-protocol/forta-node/config"
	"github.com/spf13/cobra"
)

func handleFortaStatus(cmd *cobra.Command, args []string) error {
	// call the runner health server on localhost
	reports := health.NewClient().CheckHealth("forta", config.DefaultHealthPort)
	sort.Slice(reports, func(i, j int) bool {
		return sort.StringsAreSorted([]string{reports[i].Name, reports[j].Name})
	})

	w := new(bytes.Buffer)
	for _, report := range reports {
		writeName(w, report.Name)

		switch report.Status {
		case health.StatusOK:
			writeColoredBall(w, color.FgGreen)
		case health.StatusDown:
			writeColoredBall(w, color.FgRed)
		case health.StatusFailing:
			writeColoredBall(w, color.FgYellow)
		case health.StatusInfo:
			writeColoredBall(w, color.FgBlue)
		}
		writeStatus(w, string(report.Status))
		if len(report.Details) > 0 {
			writeStatus(w, ": ") // put colon at the end of the status
			writeDetails(w, report.Status, report.Details)
		}

		fmt.Fprint(w, "\n")
	}

	fmt.Fprint(os.Stdout, w.String())
	return nil
}

func writeStatus(w io.Writer, s string) {
	color.New(color.Bold).Fprint(w, s)
}

func writeName(w io.Writer, s string) {
	color.New(color.FgWhite, color.Bold).Fprintln(w, s)
}

func writeDetails(w io.Writer, status health.Status, s string) {
	switch status {
	case health.StatusOK, health.StatusInfo:
		color.New(color.Faint).Fprintln(w, s)
	default:
		color.New(color.FgYellow).Fprintln(w, s)
	}
}

const ballPrefix = "⬤ "

func writeColoredBall(w io.Writer, c color.Attribute) {
	if c == 0 {
		fmt.Fprintf(w, ballPrefix)
		return
	}
	color.New(c).Fprintf(w, ballPrefix)
}
