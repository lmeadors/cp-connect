package lib

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
)

type PluginsRequest struct {
	Command *flag.FlagSet
	Host    *string
}

func BuildPluginsRequest() PluginsRequest {

	flagSet := flag.NewFlagSet("plugins", flag.ExitOnError)

	return PluginsRequest{
		Command: flagSet,
		Host:    flagSet.String("host", "http://localhost:8083", "cluster host"),
	}

}

func Plugins(request PluginsRequest) {

	logger := log.Default()

	path := "/connector-plugins/"

	// build request
	req, err := http.NewRequest("GET", *request.Host+path, nil)
	if err != nil {
		fmt.Print(err.Error())
	}
	req.Header.Add("Content-Type", "application/json")

	logger.Printf("url: %s\n", req.URL.String())
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
	var cpResponse PluginsResponse
	err2 := json.Unmarshal(bodyBytes, &cpResponse)
	if err2 != nil {
		fmt.Print(err.Error())
	}

	// print it
	cpResponse.Show()

}

type PluginResponse struct {
	Class   string `json:"class"`
	Type    string `json:"type"`
	Version string `json:"version"`
}

type PluginsResponse []PluginResponse

func (responseObject PluginsResponse) Show() {
	for _, plugin := range responseObject {
		fmt.Printf("%-80s %-25s %s\n", plugin.Class, plugin.Version, plugin.Type)
	}
}
