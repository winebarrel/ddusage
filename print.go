package ddusage

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
	"github.com/dustin/go-humanize"
	"github.com/olekukonko/tablewriter"
	"github.com/winebarrel/ddusage/internal/util"
)

type Usage float64

func (c Usage) String() string {
	return fmt.Sprintf("%.0f", c)
}

func (c Usage) Float64() float64 {
	return float64(c)
}

// month/usage
type UsageByMonth map[string]Usage

// product_name/month/usage
type UsageByProduct map[string]UsageByMonth

// org_name/product_name/month/cost
type UsageBreakdown map[string]UsageByProduct

func breakdownUsage(resp *datadogV1.UsageSummaryResponse) (UsageBreakdown, []string) {
	ubd := UsageBreakdown{}
	monthSet := map[string]struct{}{}

	for _, u := range resp.Usage {
		month := u.Date.Format("2006-01")
		monthSet[month] = struct{}{}

		breakdownUsageByOrg := func(org string, props map[string]any) {
			byProduct := util.MapValueOrDefault(ubd, org, UsageByProduct{})

			for product, usage := range props {
				if usage == nil {
					continue
				}

				if v, ok := usage.(float64); ok {
					byMonth := util.MapValueOrDefault(byProduct, product, UsageByMonth{})
					byMonth[month] = Usage(v)
				}
			}
		}

		if len(u.Orgs) == 0 {
			breakdownUsageByOrg("-", u.AdditionalProperties)
		} else {
			for _, org := range u.Orgs {
				breakdownUsageByOrg(*org.Name, org.AdditionalProperties)
			}
		}
	}

	return ubd, util.MapSortKeys(monthSet)
}

func printTable(resp *datadogV1.UsageSummaryResponse, out io.Writer, h bool) {
	ubd, months := breakdownUsage(resp)

	table := tablewriter.NewWriter(out)
	table.SetBorder(false)

	if h {
		table.SetAlignment(tablewriter.ALIGN_LEFT)
	}

	header := []string{"org", "product"}
	header = append(header, months...)
	table.SetHeader(header)

	printTable0(ubd, months, out, h, func(row []string) {
		table.Append(row)
	})

	table.Render()
}

func printTSV(resp *datadogV1.UsageSummaryResponse, out io.Writer, sep string) {
	ubd, months := breakdownUsage(resp)

	header := []string{"org", "product"}
	header = append(header, months...)
	fmt.Fprintln(out, strings.Join(header, sep))

	printTable0(ubd, months, out, false, func(row []string) {
		if strings.Join(row, "") != "" {
			fmt.Fprintln(out, strings.Join(row, sep))
		} else {
			fmt.Fprintln(out)
		}
	})
}

func printTable0(ubd UsageBreakdown, months []string, out io.Writer, h bool, procRow func([]string)) {
	emptyLine := make([]string, len(months)+2)
	idxOrg := 0

	for _, org := range util.MapSortKeys(ubd) {
		usageByProduct := ubd[org]
		rows := [][]string{}

		for _, product := range util.MapSortKeys(usageByProduct) {
			usageByMonth := usageByProduct[product]
			row := []string{"", product}

			for _, month := range util.MapSortKeys(usageByMonth) {
				usage := usageByMonth[month]
				column := ""

				if h {
					if usage != 0 {
						v, unit := humanize.ComputeSI(usage.Float64())
						column = humanize.FtoaWithDigits(v, 2) + unit
					}
				} else {
					column = usage.String()
				}

				row = append(row, column)
			}

			if h && strings.Join(row[2:], "") == "" {
				continue
			}

			rows = append(rows, row)
		}

		if len(rows) > 0 {
			if idxOrg != 0 {
				procRow(emptyLine)
			}

			rows[0][0] = org

			for _, row := range rows {
				procRow(row)
			}

			idxOrg++
		}
	}
}

func printJSON(resp *datadogV1.UsageSummaryResponse, out io.Writer) {
	ubd, _ := breakdownUsage(resp)
	m, _ := json.MarshalIndent(ubd, "", "  ")
	fmt.Fprintln(out, string(m))
}
