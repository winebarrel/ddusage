package ddusage

import (
	"time"

	"github.com/araddon/dateparse"
)

var (
	defaultStartMonth time.Time
	defaultEndMonth   time.Time
)

func init() {
	defaultEndMonth = time.Now()
	defaultStartMonth = defaultEndMonth.AddDate(0, -6, 0)
}

type ClientOptions struct {
	APIKey string `env:"DD_API_KEY" required:"" help:"Datadog API key."`
	APPKey string `env:"DD_APP_KEY" required:"" help:"Datadog APP key."`
}

type PrintUsageSummaryOptions struct {
	IncludeOrgDetails bool   `short:"x" help:"Include usage summaries for each sub-org.."`
	Output            string `short:"o" enum:"table,tsv,json,csv" default:"table" help:"Formatting style for output (table, tsv, json, csv)."`
	StartMonth        string `short:"s" help:"Usage beginning this month."`
	EndMonth          string `short:"e" help:"Usage ending this month."`
	Humanize          bool   `short:"H" help:"Convert usage numbers to human-friendly strings."`
}

func (options *PrintUsageSummaryOptions) calcPeriod() (time.Time, time.Time, error) {
	timeStartMonth := defaultStartMonth
	timeEndMonth := defaultEndMonth

	if options.StartMonth != "" {
		t, err := dateparse.ParseAny(options.StartMonth)

		if err != nil {
			return timeStartMonth, timeEndMonth, err
		}

		timeStartMonth = t
	}

	if options.EndMonth != "" {
		t, err := dateparse.ParseAny(options.EndMonth)

		if err != nil {
			return timeStartMonth, timeEndMonth, err
		}

		timeEndMonth = t
	}

	return timeStartMonth, timeEndMonth, nil
}
