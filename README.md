# Sumo Logic Search Job CLI

## Install

Install Golang (Homebrew or https://go.dev/doc/install).

Build a binary:
```bash
git clone git@github.com:nhoag/sumo-search-job-cli.git
cd sumo-search-job-cli
cp .sumo-search-job-cli.yaml.dist ~/.sumo-search-job-cli.yaml
go build -o sumo
./sumo -h
```

Create a Sumo Logic Access Key [here](https://service.sumologic.com/ui/#/preferences).

Add credentials to `~/.sumo-search-job-cli.yaml`.

## Example Commands

Perform the full life-cycle of initiating a search job, polling for status, fetching results, and deleting the job:
```bash
./sumo jobProcessFull -J ./resources/jobDefinition.json
```
