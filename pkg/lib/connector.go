package lib

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func BuildConnectorFlagSet(cmdMap map[string]CpCommand) ConnectorRequest {

	name := "connector"
	connectorCmd := flag.NewFlagSet(name, flag.ExitOnError)

	request := ConnectorRequest{
		Name:    name,
		Command: connectorCmd,
		Host:    connectorCmd.String("host", getEnv("CP_CONNECT_HOST", "http://localhost:8083"), "cluster host"),
		Json:    connectorCmd.String("json", "", "json configuration file"),
		Action:  connectorCmd.String("action", "Config", "action to perform (Config | Validate | Restart | RestartAll | Put | Status | Pause | Resume | Delete)"),
	}
	cmdMap[name] = request
	return request

}

type ConnectorRequest struct {
	Name    string
	Command *flag.FlagSet
	Host    *string
	Json    *string
	Action  *string
}

func (request ConnectorRequest) GetCommand() flag.FlagSet {
	return *request.Command
}

func (request ConnectorRequest) GetName() string {
	return request.Name
}

func (request ConnectorRequest) PrintDefaults() {
	request.Command.PrintDefaults()
}

func (request ConnectorRequest) Execute() {
	//	Connector(request)
	//}
	//
	//func Connector(request ConnectorRequest) {

	logger := log.Default()

	if len(*request.Json) == 0 {
		err := filepath.Walk(".", func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && filepath.Ext(path) == ".json" {
				*request.Json = path
			}
			return nil
		})
		if err != nil {
			log.Fatalf("unable to find json configuration\n")
		}
	}

	logger.Printf("host:            %s\n", *request.Host)
	logger.Printf("json file:       %s\n", *request.Json)
	logger.Printf("action:          %s\n", *request.Action)

	jsonFile, err := os.Open(*request.Json)
	if err != nil {
		log.Fatalf("unable to read json configuration %s\n", request.Json)
	}
	defer jsonFile.Close()
	jsonBytes, _ := io.ReadAll(jsonFile)
	//logger.Printf(string(jsonBytes))
	var config JsonConfig
	json.Unmarshal(jsonBytes, &config)

	logger.Printf("connector name:  %s\n", config.Name)
	logger.Printf("connector class: %s\n", config.ConnectorClass)

	var uri string

	switch *request.Action {
	case "Put":
		//	PUT /connectors/(string:name)/config
		uri = fmt.Sprintf("/connectors/%s/config", config.Name)
		req, _ := http.NewRequest("PUT", *request.Host+uri, bytes.NewReader(jsonBytes))
		req.Header.Add("Content-Type", "application/json")
		client := http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Print(err.Error())
		}
		defer resp.Body.Close()
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Print(err.Error())
		}
		dst := &bytes.Buffer{}
		if err := json.Indent(dst, bodyBytes, "", "  "); err != nil {
			panic(err)
		}
		fmt.Println(dst.String())
	case "Validate":
		uri = fmt.Sprintf("/connector-plugins/%s/config/validate", config.ConnectorClass)
		req, _ := http.NewRequest("PUT", *request.Host+uri, bytes.NewReader(jsonBytes))
		req.Header.Add("Content-Type", "application/json")
		client := http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Print(err.Error())
		}
		defer resp.Body.Close()
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Print(err.Error())
		}
		dst := &bytes.Buffer{}
		if err := json.Indent(dst, bodyBytes, "", "  "); err != nil {
			panic(err)
		}
		fmt.Println(dst.String())
	case "Pause":

		uri = fmt.Sprintf("/connectors/%s/pause", config.Name)
		req, _ := http.NewRequest("PUT", *request.Host+uri, nil)

		resp := executeRequestWithoutResponseBody(req)

		fmt.Println("response status: ", resp.Status)

	case "Resume":

		uri = fmt.Sprintf("/connectors/%s/resume", config.Name)
		req, _ := http.NewRequest("PUT", *request.Host+uri, nil)

		resp := executeRequestWithoutResponseBody(req)

		fmt.Println("response status: ", resp.Status)

	case "Restart":

		uri = fmt.Sprintf("/connectors/%s/restart?includeTasks=true&onlyFailed=true", config.Name)
		req, _ := http.NewRequest("POST", *request.Host+uri, nil)

		resp := executeRequestWithoutResponseBody(req)

		fmt.Println("response status: ", resp.Status)

	case "RestartAll":

		uri = fmt.Sprintf("/connectors/%s/restart?includeTasks=true&onlyFailed=false", config.Name)
		req, _ := http.NewRequest("POST", *request.Host+uri, nil)

		resp := executeRequestWithoutResponseBody(req)

		fmt.Println("response status: ", resp.Status)

	case "Delete":

		// DELETE /connectors/(string:name)/
		uri = fmt.Sprintf("/connectors/%s", config.Name)
		req, _ := http.NewRequest("DELETE", *request.Host+uri, nil)

		resp := executeRequestWithoutResponseBody(req)

		fmt.Println("response status: ", resp.Status)

	case "Status":
		uri = fmt.Sprintf("/connectors/%s/status", config.Name)
		req, _ := http.NewRequest("GET", *request.Host+uri, nil)
		req.Header.Add("Content-Type", "application/json")
		client := http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Print(err.Error())
		}
		defer resp.Body.Close()
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Print(err.Error())
		}
		dst := &bytes.Buffer{}
		if err := json.Indent(dst, bodyBytes, "", "  "); err != nil {
			panic(err)
		}
		fmt.Println(dst.String())
	default:
		// Config is the default behavior
		uri = fmt.Sprintf("/connectors/%s/config", config.Name)
		req, _ := http.NewRequest("GET", *request.Host+uri, nil)
		req.Header.Add("Content-Type", "application/json")
		client := http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Print(err.Error())
		}
		defer resp.Body.Close()
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Print(err.Error())
		}
		dst := &bytes.Buffer{}
		if err := json.Indent(dst, bodyBytes, "", "  "); err != nil {
			panic(err)
		}
		fmt.Println(dst.String())
	}
	logger.Printf("uri:             %s\n", uri)

}

type JsonConfig struct {
	Name           string `json:"name"`
	ConnectorClass string `json:"connector.class"`
}

func executeRequestWithoutResponseBody(req *http.Request) *http.Response {
	req.Header.Add("Content-Type", "application/json")
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Print(err.Error())
	}
	if err != nil {
		fmt.Print(err.Error())
	}
	return resp
}
