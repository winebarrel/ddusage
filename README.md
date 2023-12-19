# ddusage

[![CI](https://github.com/winebarrel/ddusage/actions/workflows/ci.yml/badge.svg)](https://github.com/winebarrel/ddusage/actions/workflows/ci.yml)

A tool that shows a breakdown of Datadog usages in a table.

## Usage

```
Usage: ddusage --api-key=STRING --app-key=STRING

Flags:
  -h, --help                   Show context-sensitive help.
      --api-key=STRING         Datadog API key ($DD_API_KEY).
      --app-key=STRING         Datadog APP key ($DD_APP_KEY).
  -x, --include-org-details    Include usage summaries for each sub-org..
  -o, --output="table"         Formatting style for output (table, tsv, json, csv).
  -s, --start-month=STRING     Usage beginning this month.
  -e, --end-month=STRING       Usage ending this month.
  -H, --humanize               Convert usage numbers to to human-friendly strings.
      --version
```

```
$ export DD_API_KEY=...
$ export DD_APP_KEY=...
$ ddusage -x -H
       ORG       |               PRODUCT               | 2022-12 | 2023-01 | 2023-02 | 2023-03 | 2023-04
-----------------+-------------------------------------+---------+---------+---------+---------+----------
  organization1  | apm_host_incl_usm_top99p            | 1       | 1       | 1       | 1       | 1
  organization2  | logs_indexed_logs_usage_sum_15_day  | 1M      | 1M      | 1M      | 1M      | 1M
```

## Installation

```
brew install winebarrel/ddusage/ddusage
```

## Related Links

- https://github.com/winebarrel/ddcost
