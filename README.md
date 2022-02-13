# Sumo Logic Search Job CLI

## Install

Install Golang (Homebrew or https://go.dev/doc/install).

Build a binary:
```bash
git clone git@github.com:nhoag/sumo-search-job-cli.git
cd sumo-search-job-cli
cp .sumo-search-job-cli.yaml.dist ~/.sumo-search-job-cli.yaml
go build -o sumo
# Optional: Move the binary into a PATH directory
./sumo -h
```

Create a Sumo Logic Access Key [here](https://service.sumologic.com/ui/#/preferences).

Add credentials to `~/.sumo-search-job-cli.yaml`.

## Example Commands

Perform the full life-cycle of initiating a search job, polling for status, fetching results, and deleting the job:
```bash
sumo jobProcessFull -J ./resources/jobDefinition.json
```

Create a search job:
```bash
sumo jobCreate -J ./resources/jobDefinition.json
```

Get search job status, and poll until complete:
```bash
sumo jobStatusCheck JOB_ID -p
```

Keep the search job alive for 1h (default lifetime is 5m after last activity):
```bash
sumo jobKeepAlive JOB_ID -k60
```

Fetch search job results after a job has completed:
```bash
sumo jobResultsGet JOB_ID -a -p
```

Delete search job:
```bash
sumo jobDelete JOB_ID
```
