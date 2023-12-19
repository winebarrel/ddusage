package ddusage

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadog"
	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
)

type Client struct {
	options *ClientOptions
	api     *datadogV1.UsageMeteringApi
}

func NewClient(options *ClientOptions) *Client {
	configuration := datadog.NewConfiguration()
	apiClient := datadog.NewAPIClient(configuration)
	api := datadogV1.NewUsageMeteringApi(apiClient)

	client := &Client{
		options: options,
		api:     api,
	}

	return client
}

func (client *Client) withAPIKey(ctx context.Context) context.Context {
	ctx = context.WithValue(
		ctx,
		datadog.ContextAPIKeys,
		map[string]datadog.APIKey{
			"apiKeyAuth": {
				Key: client.options.APIKey,
			},
			"appKeyAuth": {
				Key: client.options.APPKey,
			},
		},
	)

	return ctx
}

func (client *Client) PrintUsageSummary(out io.Writer, options *PrintUsageSummaryOptions) error {
	timeStartMonth, timeEndMonth, err := options.calcPeriod()

	if err != nil {
		return err
	}

	ctx := client.withAPIKey(context.Background())
	resp, _, err := client.api.GetUsageSummary(
		ctx,
		timeStartMonth,
		*datadogV1.NewGetUsageSummaryOptionalParameters().
			WithEndMonth(timeEndMonth).
			WithIncludeOrgDetails(options.IncludeOrgDetails),
	)

	if err != nil {
		var dderr datadog.GenericOpenAPIError

		if errors.As(err, &dderr) {
			err = fmt.Errorf("%w: %s", err, dderr.ErrorBody)
		}

		return err
	}

	switch options.Output {
	case "table":
		printTable(&resp, out, options.Humanize)
	case "tsv":
		printTSV(&resp, out, "\t")
	case "json":
		printJSON(&resp, out)
	case "csv":
		printTSV(&resp, out, ",")
	}

	return nil
}
