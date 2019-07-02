# `grafanactl` a simple CLI for grafana 

`grafanactl` is a simple CLI that allows to perform operation on folders and dashboards in Grafana.

# Installation

Download latest release:

```
curl --silent --location "https://github.com/overdrive3000/grafanactl/releases/download/0.1.2/grafanactl-darwin-amd64" -o grafanactl
sudo mv grafanactl-darwin-amd64 /usr/local/bin/grafanactl
sudo chmod +x /usr/local/bin/grafanactl
```

Check here the release for your OS https://github.com/overdrive3000/grafanactl/releases/tag/0.1.2

# Usage

In order to get `grafanactl` performing operations against your Grafana instance you must first create an API key that will be used to authorize `grafanactl`. To create the API key follow the instructions shown at [Create API Token](https://grafana.com/docs/http_api/auth/#create-api-token)

Once you get the API key you must configure `grafanactl` to use it. `grfanactl` by default will read a configuration file at `~/.grafanactl` using the below format:


```
url: <GRAFANA_URL> 
apiKey: <GRAFANA_API_KEY>
```

Example:

```
url: http://localhost:8080
apiKey: eyJrIjoicllsVW82Q2xENDR6UjZNTTVvazRRa2VybDZNd0Q5ZEciLCJuIjoiYXV0b21hdGlvbiIsImlkIjoxfQ==
```

This configuration can be also overrided via CLI flags as shown below:

```
grafanactl --url http://localhost:8080 --key eyJrIjoicllsVW82Q2xENDR6UjZNTTVvazRRa2VybDZNd0Q5ZEciLCJuIjoiYXV0b21hdGlvbiIsImlkIjoxfQ==
```

## Help

You can get a list of all available commands and flags by using running:

```
$ grafanactl --help

Usage:
  grafanactl [command]

Available Commands:
  dashboard   Grafana Dashboards
  folder      Grafana Folders
  help        Help about any command

Flags:
  -c, --config string      config file (default is $HOME/.grafanactl.yaml)
  -h, --help               help for grafanactl
      --key string         Grafana API Key
  -o, --output string      Output format (table, json) (default "table")
      --url string         Grafana URL (https://localhost:3000)
  -v, --verbosity string   Log level (debug, warn, error) (default "warning")
      --version            version for grafanactl

Use "grafanactl [command] --help" for more information about a command.
```

## Examples

1. Creating a Grafana folder
```
grafanactl folder create --title Test
```

2. List all folders
```
grafanactl folder list
```

3. Get information about an specific folder
```
grafanactl folder get --name 'Kubernetes Dashboards'

grafanactl folder get --name 'Kubernetes Dashboards' -o json
```

4. List all dashboards in a folder
```
grafanactl dashboard search --folder-id 2
```

5. Search folder by name
```
grafanactl dashboard search --name Kafka
```

6. Get whole dashboard information
```
grafanactl dashboard get --uid fa49a4706d07a042595b664c87fb33ea -o json 
```

7. Import all dashboards in a local folder into Grafana
```
for i in `ls -w1 mydashboards/`; do grafanactl dashboard create -f $i; done
```
