package cmd

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/forta-protocol/forta-core-go/clients/health"
	"github.com/forta-protocol/forta-node/config"
	"github.com/spf13/cobra"
)

// status formats
const (
	StatusFormatPretty  = "pretty"
	StatusFormatOneline = "oneline"
	StatusFormatJSON    = "json"
	StatusFormatCSV     = "csv"

	StatusShowSummary   = "summary"
	StatusShowImportant = "important"
	StatusShowAll       = "all"
)

var ballPrefix = "⬤ "

func handleFortaStatus(cmd *cobra.Command, args []string) error {
	format, err := cmd.Flags().GetString("format")
	if err != nil {
		return err
	}

	show, err := cmd.Flags().GetString("show")
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
	allReports := health.NewClient().CheckHealth("forta", config.DefaultHealthPort)
	sort.Slice(allReports, func(i, j int) bool {
		return sort.StringsAreSorted([]string{allReports[i].Name, allReports[j].Name})
	})

	var reports health.Reports
	for _, report := range allReports {
		var shouldInclude bool
		switch show {
		case StatusShowSummary:
			shouldInclude = strings.Contains(report.Name, "summary")

		case StatusShowImportant:
			shouldInclude = report.Status != health.StatusInfo

		case StatusShowAll:
			shouldInclude = true
		}
		if shouldInclude {
			reports = append(reports, report)
		}
	}

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

		writeStatusBall(w, report.Status)

		writeStatus(w, string(report.Status))
		if len(report.Details) > 0 {
			writeStatus(w, ": ") // put colon at the end of the status
			writeDetails(w, report.Status, report.Details)
		}

		fmt.Fprint(w, "\n\n")
	}
	fmt.Fprint(os.Stdout, w.String())
}

func formatReportsOneline(reports health.Reports) {
	w := new(bytes.Buffer)
	for _, report := range reports {
		writeStatusBall(w, report.Status)

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

func writeStatusBall(w io.Writer, status health.Status) {
	switch status {
	case health.StatusOK:
		writeColoredBall(w, color.FgGreen)
	case health.StatusDown:
		writeColoredBall(w, color.FgRed)
	case health.StatusFailing:
		writeColoredBall(w, color.FgYellow)
	case health.StatusLagging:
		writeColoredBall(w, color.FgYellow)
	case health.StatusInfo:
		writeColoredBall(w, color.FgBlue)
	case health.StatusUnknown:
		writeColoredBall(w, color.Faint)
	}
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
