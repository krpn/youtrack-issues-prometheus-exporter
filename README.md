# youtrack-issues-prometheus-exporter

[![Build Status](https://travis-ci.org/krpn/youtrack-issues-prometheus-exporter.svg?branch=master)](https://travis-ci.org/krpn/youtrack-issues-prometheus-exporter) [![Quality Gate](https://sonarcloud.io/api/project_badges/measure?project=krpn_youtrack-issues-prometheus-exporter&metric=alert_status)](https://sonarcloud.io/dashboard?id=krpn_youtrack-issues-prometheus-exporter) [![Coverage Status](https://sonarcloud.io/api/project_badges/measure?project=krpn_youtrack-issues-prometheus-exporter&metric=coverage)](https://sonarcloud.io/component_measures?id=krpn_youtrack-issues-prometheus-exporter&metric=coverage) [![Technical Debt](https://sonarcloud.io/api/project_badges/measure?project=krpn_youtrack-issues-prometheus-exporter&metric=sqale_index)](https://sonarcloud.io/component_measures?id=krpn_youtrack-issues-prometheus-exporter&metric=sqale_index) [![License](https://img.shields.io/github/license/krpn/youtrack-issues-prometheus-exporter.svg)](https://github.com/krpn/youtrack-issues-prometheus-exporter/blob/master/LICENSE)

youtrack-issues-prometheus-exporter exports YouTrack issues to Prometheus for any search queries

# Table of Contents
* [Features](#features)
* [Quick Start](#quick-start)
* [Configuration](#configuration)
* [Exposed Prometheus Metrics](#exposed-prometheus-metrics)
* [Command-Line Flags](#command-line-flags)
* [Contribute](#contribute)

# Features

* Export issues for any search query from config
* [!] Works only with YouTrack 2018.3 and above because uses "new" REST API
* A docker image available on [Docker Hub](https://hub.docker.com/r/krpn/youtrack-issues-prometheus-exporter/)

[(back to top)](#youtrack-issues-prometheus-exporter)

# Quick Start

1. Prepare config.json file based on [example](https://github.com/krpn/youtrack-issues-prometheus-exporter/blob/master/example/config.json) (details in [configuration](#configuration))

2. Run container with command ([cli flags](#command-line-flags)):

    `docker run -d -p <port>:8080 -v <path to config.json dir>:/config --name youtrack-exporter krpn/youtrack-issues-prometheus-exporter`

3. Checkout logs (will be empty if ok):

    `docker logs youtrack-exporter`
    
4. Add youtrack-exporter instance to [Prometheus scrape targets](https://prometheus.io/docs/prometheus/latest/configuration/configuration/#%3Cscrape_config%3E) (port from `docker run` command)

5. Add [alerting rules](https://prometheus.io/docs/prometheus/latest/configuration/alerting_rules/) based on [metrics](#exposed-prometheus-metrics) if needed, for example:

    ```yaml
    groups:
    - name: youtrack.rules
      rules:
      
      - alert: YouTrackShowStopper
        expr: youtrack_issues{query="showstopper"} == 1
        for: 1m
        annotations:
          description: 'Show-Stopper {{ $labels.id }} {{ $labels.title }}: https://youtrack.company.com/issue/{{ $labels.id }}'
      
      - alert: YouTrackExporterError
        expr: sum(increase(youtrack_errors[1m])) by (error) > 0
        for: 1m
        annotations:
          description: 'YouTrack exporter got error: {{ $labels.error }}'
    ```

[(back to top)](#youtrack-issues-prometheus-exporter)

# Configuration

Configuration file based on JSON format. Example:

```json
{
  "endpoint": "https://youtrack.company.com/",
  "token": "perm:YWxleGtydXBpbg==.QWxleGFuZGVy.9nvYkHL4aHy0zHaEGIXmjcGjVNx6Kr",
  "queries": {
    "showstopper": "Show-Stopper #Unresolved #Unassigned",
    "unresolved": "#Unresolved State: Submitted"
  },
  "refresh_delay_seconds": 10,
  "request_timeout_seconds": 10,
  "listen_port": 8080
}
```
| Setting                   | Type      | Description                                                                                                                              | Example                                                                                                 |
|---------------------------|:---------:|------------------------------------------------------------------------------------------------------------------------------------------|---------------------------------------------------------------------------------------------------------|
| `endpoint`                | `string`  | YouTrack URL without path                                                                                                                | `https://youtrack.company.com/`                                                                         |
| `token`                   | `string`  | [YouTrack API permanent token](https://www.jetbrains.com/help/youtrack/standalone/authentication-with-permanent-token.html)              | `perm:YWxleGtydXBpbg==.QWxleGFuZGVy.9nvYkHL4aHy0zHaEGIXmjcGjVNx6Kr`                                     |
| `queries`                 | `object`  | Map of search queries where key is search query name and value is search query string. Query name will be passed to metric label `query` | `{"showstopper": "Show-Stopper #Unresolved #Unassigned", "unresolved": "#Unresolved State: Submitted"}` |
| `refresh_delay_seconds`   | `integer` | (optional, default: 10) Refresh metrics delay seconds. Metrics automatically refreshes in background                                     | `60`                                                                                                    |
| `request_timeout_seconds` | `integer` | (optional, default: 10) Request timeout seconds for YouTrack REST API HTTP request                                                       | `30`                                                                                                    |
| `listen_port`             | `integer` | (optional, default: 8080) HTTP port to listen on                                                                                         | `80`                                                                                                    |

[(back to top)](#youtrack-issues-prometheus-exporter)

# Exposed Prometheus Metrics

| Name              | Description                                                                                              | Labels               |
|-------------------|----------------------------------------------------------------------------------------------------------|----------------------|
| `youtrack_issues` | Query issues. Equals `1` if task for this query is found. Equals `0` if not found (but was found before) | `query` `id` `title` |
| `youtrack_errors` | Errors counter. Increments when error is occurred                                                        | `query` `error`      |

[(back to top)](#youtrack-issues-prometheus-exporter)

# Command-Line Flags

Usage: `youtrack-issues-prometheus-exporter [<flags>]`

| Flag                 | Type     | Description         | Default              |
|----------------------|:--------:|---------------------|----------------------|
| `-c` or `--config`   | `string` | Path to config file | `config/config.json` |
| `--help`             |          | Show help           |                      |

[(back to top)](#youtrack-issues-prometheus-exporter)

# Contribute

Please feel free to send me [pull requests](https://github.com/krpn/youtrack-issues-prometheus-exporter/pulls).

[(back to top)](#youtrack-issues-prometheus-exporter)