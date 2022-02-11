# Sumo Logic Search Job CLI

## Install

```bash
git clone git@github.com:nhoag/sumo-search-job-cli.git
cd sumo-search-job-cli
cp .sumo-search-job-cli.yaml.dist ~/.sumo-search-job-cli.yaml
go build -o sumo
./sumo -h
```

Create an Access Key [here](https://service.sumologic.com/ui/#/preferences).

Add credentials to `~/.sumo-search-job-cli.yaml`.
