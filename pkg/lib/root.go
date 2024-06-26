package lib

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

type ClusterRequest struct {
	Command *flag.FlagSet
	Name    string
	Host    *string
}

func (c ClusterRequest) GetCommand() flag.FlagSet {
	return *c.Command
}

func (c ClusterRequest) GetName() string {
	return c.Name
}

func (c ClusterRequest) PrintDefaults() {
	c.Command.PrintDefaults()
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

func BuildClusterFlagSet(cmdMap map[string]CpCommand) ClusterRequest {

	name := "cluster"
	flagSet := flag.NewFlagSet(name, flag.ExitOnError)

	request := ClusterRequest{
		Command: flagSet,
		Name:    name,
		Host:    flagSet.String("host", getEnv("CP_CONNECT_HOST", "http://localhost:8083"), "cluster host"),
	}

	cmdMap[name] = request

	return request

}

func (request ClusterRequest) Execute() {
	//Cluster(request)
	//}
	//
	//func Cluster(request ClusterRequest) {

	path := "/"

	// build request
	req, err := http.NewRequest("GET", *request.Host+path, nil)
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
