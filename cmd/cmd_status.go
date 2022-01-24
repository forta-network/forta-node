package cmd

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/fatih/color"
	"github.com/forta-protocol/forta-node/clients/health"
	"github.com/forta-protocol/forta-node/config"
	"github.com/spf13/cobra"
)

// status formats
const (
	StatusFormatPretty  = "pretty"
	StatusFormatOneline = "oneline"
	StatusFormatJSON    = "json"
	StatusFormatCSV     = "csv"
)

var ballPrefix = "⬤ "

func handleFortaStatus(cmd *cobra.Command, args []string) error {
	format, err := cmd.Flags().GetString("format")
	if err != nil {
		return err
	}

	noColor, err := cmd.Flags().GetBool("no-color")
	if err != nil {
		return err
	}
	if noColor {
		color.NoColor = true
		ballPrefix = ""
	}

	// call the runner health server on localhost
	reports := health.NewClient().CheckHealth("forta", config.DefaultHealthPort)
	sort.Slice(reports, func(i, j int) bool {
		return sort.StringsAreSorted([]string{reports[i].Name, reports[j].Name})
	})

	switch format {
	case StatusFormatPretty:
		formatReportsPretty(reports)

	case StatusFormatOneline:
		formatReportsOneline(reports)

	case StatusFormatJSON:
		color.NoColor = true
		return formatReportsJSON(reports)

	case StatusFormatCSV:
		color.NoColor = true
		return formatReportsCSV(reports)

	default:
		return fmt.Errorf("unknown format: %v", format)
	}

	return nil
}

func formatReportsPretty(reports health.Reports) {
	w := new(bytes.Buffer)
	for _, report := range reports {
		writeName(w, report.Name)
		fmt.Fprint(w, "\n")

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
}

func formatReportsOneline(reports health.Reports) {
	w := new(bytes.Buffer)
	for _, report := range reports {
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
		fmt.Fprint(w, " | ")
		writeName(w, report.Name)
		if len(report.Details) > 0 {
			fmt.Fprint(w, " | ")
			writeDetails(w, report.Status, report.Details)
		}

		fmt.Fprint(w, "\n")
	}

	fmt.Fprint(os.Stdout, w.String())
}

func formatReportsJSON(reports health.Reports) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(reports)
}

func formatReportsCSV(reports health.Reports) error {
	w := csv.NewWriter(os.Stdout)
	defer w.Flush()

	for _, report := range reports {
		if err := w.Write([]string{report.Name, string(report.Status), report.Details}); err != nil {
			return fmt.Errorf("failed to write csv record: %v", err)
		}
	}
	return nil
}

func writeStatus(w io.Writer, s string) {
	color.New(color.Bold).Fprint(w, s)
}

func writeName(w io.Writer, s string) {
	color.New(color.FgWhite, color.Bold).Fprint(w, s)
}

func writeDetails(w io.Writer, status health.Status, s string) {
	switch status {
	case health.StatusOK, health.StatusInfo:
		color.New(color.Faint).Fprint(w, s)
	default:
		color.New(color.FgYellow).Fprint(w, s)
	}
}

func writeColoredBall(w io.Writer, c color.Attribute) {
	if c == 0 {
		fmt.Fprintf(w, ballPrefix)
		return
	}
	color.New(c).Fprint(w, ballPrefix)
}
