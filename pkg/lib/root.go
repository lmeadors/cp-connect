package lib

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
)

func BuildClusterFlagSet() ClusterRequest {

	flagSet := flag.NewFlagSet("cluster", flag.ExitOnError)

	return ClusterRequest{
		Command: flagSet,
		Host:    flagSet.String("host", "http://localhost:8083", "cluster host"),
	}

}

type ClusterRequest struct {
	Command *flag.FlagSet
	Host    *string
}

func Cluster(request ClusterRequest) {

	path := "/"
	host := *request.Host
	// build request
	req, err := http.NewRequest("GET", host+path, nil)
	if err != nil {
		fmt.Print(err.Error())
	}
	req.Header.Add("Content-Type", "application/json")

	// execute request
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Print(err.Error())
	}

	// schedule the close of the response
	defer resp.Body.Close()

	// read the response
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err.Error())
	}

	// convert to our response
	var cpResponse ClusterResponse
	err2 := json.Unmarshal(bodyBytes, &cpResponse)
	if err2 != nil {
		fmt.Print(err.Error())
	}

	// print it
	cpResponse.Show()

}

type ClusterResponse struct {
	Version        string `json:"version"`
	Commit         string `json:"commit"`
	KafkaClusterId string `json:"kafka_cluster_id"`
}

func (responseObject ClusterResponse) Show() {
	fmt.Printf("Version:         %s\n", responseObject.Version)
	fmt.Printf("Commit:          %s\n", responseObject.Commit)
	fmt.Printf("KafakaClusterId: %s\n", responseObject.KafkaClusterId)
}

type CPResponse interface {
	Show()
}
