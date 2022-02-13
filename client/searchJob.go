package client

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"

	openapi "github.com/nhoag/sumologic-search-job-client-go"
	"github.com/spf13/viper"
)

func getClient() *openapi.APIClient {
	configuration := openapi.NewConfiguration()
	// @todo: Update endpoint resolution to be more robust
	endpoint := getEndpoint(configuration.Servers)
	if len(endpoint) == 0 {
		endpoint = viper.GetString("host")
	}
	configuration.Host = endpoint
	return openapi.NewAPIClient(configuration)
}

func getEndpoint(servers openapi.ServerConfigurations) string {
	serverConfig := getServerConfig(servers)
	m := regexp.MustCompile(`[^http(s)?//:][a-z0-9.-]+[^/api]`)
	return m.FindString(serverConfig.URL)
}

func getServerConfig(servers openapi.ServerConfigurations) openapi.ServerConfiguration {
	deployment := viper.GetString("deployment")
	for _, server := range servers {
		if strings.Contains(server.Description, strings.ToUpper(deployment)) {
			return server
		}
	}
	return openapi.ServerConfiguration{}
}

func getContext() context.Context {
	return context.WithValue(context.Background(), openapi.ContextBasicAuth, openapi.BasicAuth{
		UserName: viper.GetString("accessId"),
		Password: viper.GetString("accessKey"),
	})
}

func CreateSearchJob(searchJob openapi.SearchJobDefinition) (*url.URL, string) {
	request := getClient().DefaultApi.CreateSearchJob(getContext()).SearchJobDefinition(searchJob)
	resp, err := request.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.CreateSearchJob``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", resp)
		os.Exit(1)
	}
	location, err := resp.Location()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when retrieving Location header: %v\n", err)
		os.Exit(1)
	}

	locationArray := strings.Split(location.String(), "/")
	jobId := locationArray[len(locationArray)-1]
	return location, jobId
}

func DeleteSearchJob(jobId string) {
	request := getClient().DefaultApi.DeleteSearchJob(getContext(), jobId)
	resp, err := request.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.DeleteSearchJob``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", resp)
		os.Exit(1)
	}
}

func GetSearchJobStatus(jobId string) *openapi.SearchJobState {
	request := getClient().DefaultApi.GetSearchJobStatus(getContext(), jobId)
	status, resp, err := request.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.GetSearchJobStatus``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", resp)
		os.Exit(1)
	}
	return status
}

func GetSearchJobMessages(jobId string, limit int32, offset int32) *openapi.SearchJobMessages {
	request := getClient().DefaultApi.GetSearchJobMessages(getContext(), jobId).Offset(offset).Limit(limit)
	messages, resp, err := request.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.DeleteSearchJob``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", resp)
		os.Exit(1)
	}
	return messages
}

func GetSearchJobRecords(jobId string, limit int32, offset int32) *openapi.SearchJobRecords {
	request := getClient().DefaultApi.GetSearchJobRecords(getContext(), jobId).Offset(offset).Limit(limit)
	records, resp, err := request.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.DeleteSearchJob``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", resp)
		os.Exit(1)
	}
	return records
}
