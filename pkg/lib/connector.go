package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func Connector(host string, jsonFilename string, action string) {

	logger := log.Default()

	if len(jsonFilename) == 0 {
		err := filepath.Walk(".", func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && filepath.Ext(path) == ".json" {
				jsonFilename = path
			}
			return nil
		})
		if err != nil {
			log.Fatalf("unable to find json configuration\n")
		}
	}

	logger.Printf("host:            %s\n", host)
	logger.Printf("json file:       %s\n", jsonFilename)
	logger.Printf("action:          %s\n", action)

	jsonFile, err := os.Open(jsonFilename)
	if err != nil {
		log.Fatalf("unable to read json configuration %s\n", jsonFilename)
	}
	defer jsonFile.Close()
	jsonBytes, _ := io.ReadAll(jsonFile)
	var config JsonConfig
	json.Unmarshal(jsonBytes, &config)

	logger.Printf("connector name:  %s\n", config.Name)
	logger.Printf("connector class: %s\n", config.ConnectorClass)

	var uri string
	//var req *http.Request
	switch action {
	case "Put":
		//	PUT /connectors/(string:name)/config
		uri = fmt.Sprintf("/connectors/%s/config", config.Name)
		req, _ := http.NewRequest("PUT", host+uri, bytes.NewReader(jsonBytes))
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
		req, _ := http.NewRequest("PUT", host+uri, bytes.NewReader(jsonBytes))
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
		req, _ := http.NewRequest("PUT", host+uri, nil)

		resp := executeRequestWithoutResponseBody(req)

		fmt.Println("response status: ", resp.Status)

	case "Resume":

		uri = fmt.Sprintf("/connectors/%s/resume", config.Name)
		req, _ := http.NewRequest("PUT", host+uri, nil)

		resp := executeRequestWithoutResponseBody(req)

		fmt.Println("response status: ", resp.Status)

	case "Delete":

		// DELETE /connectors/(string:name)/
		uri = fmt.Sprintf("/connectors/%s", config.Name)
		req, _ := http.NewRequest("DELETE", host+uri, nil)

		resp := executeRequestWithoutResponseBody(req)

		fmt.Println("response status: ", resp.Status)

	case "Status":
		uri = fmt.Sprintf("/connectors/%s/status", config.Name)
		req, _ := http.NewRequest("GET", host+uri, nil)
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
		req, _ := http.NewRequest("GET", host+uri, nil)
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
