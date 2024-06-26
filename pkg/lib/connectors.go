package lib

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
)

func BuildConnectorsFlagSet(cmdMap map[string]CpCommand) ConnectorsRequest {

	name := "connectors"
	flagSet := flag.NewFlagSet(name, flag.ExitOnError)

	request := ConnectorsRequest{
		Name:     name,
		Command:  flagSet,
		Host:     flagSet.String("host", getEnv("CP_CONNECT_HOST", "http://localhost:8083"), "cluster host"),
		Expand:   flagSet.String("expand", None.Name(), "expanded info (None | Status | Info)"),
		ConnName: flagSet.String("name", "", "connector name"),
		Config:   flagSet.Bool("config", false, "configuration only"),
		Status:   flagSet.Bool("status", false, "status only"),
		Tasks:    flagSet.Bool("tasks", false, "tasks only"),
	}
	cmdMap[name] = request

	return request

}

type ConnectorsRequest struct {
	Name     string
	Command  *flag.FlagSet
	Host     *string
	Expand   *string
	ConnName *string
	Config   *bool
	Status   *bool
	Tasks    *bool
}

func (request ConnectorsRequest) GetCommand() flag.FlagSet {
	return *request.Command
}

func (request ConnectorsRequest) GetName() string {
	return request.Name
}

func (request ConnectorsRequest) PrintDefaults() {
	request.Command.PrintDefaults()
}

func (request ConnectorsRequest) Execute() {
	//	Connectors(request)
	//}
	//
	//func Connectors(request ConnectorsRequest) {

	logger := log.Default()

	expand := ExpandFromString(*request.Expand)

	path := "/connectors"
	// build request

	var uri string
	if len(*request.ConnName) > 0 {
		uri = *request.Host + path + "/" + *request.ConnName
		if *request.Config {
			uri += "/config"
		} else if *request.Status {
			uri += "/status"
		} else if *request.Tasks {
			uri += "/tasks"
		}
	} else {
		uri = *request.Host + path + expand.UrlParam()
	}

	logger.Printf("host:   %s\n", *request.Host)
	logger.Printf("name:   %s\n", *request.ConnName)
	logger.Printf("expand: %s\n", expand.Name())
	logger.Printf("uri:    %s\n", uri)

	req, err := http.NewRequest("GET", uri, nil)

	if err != nil {
		fmt.Print(err.Error())
	}

	//req.Header.Add("Accept", "application/json")
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

	// convert to our response object

	//var cpResponse CPResponse
	//err2 := json.Unmarshal(bodyBytes, &cpResponse)
	//if err2 != nil {
	//	fmt.Print(err.Error())
	//}

	// print it
	//cpResponse.Show()
	dst := &bytes.Buffer{}
	if err := json.Indent(dst, bodyBytes, "", "  "); err != nil {
		panic(err)
	}
	fmt.Println(dst.String())

}

type ConnectorsResponse struct {
	Connectors map[string]interface{}
}

type ConnectorsResponseInfo struct {
	Info struct {
		Name   string            `json:"name"`
		Config map[string]string `json:"config"`
		Tasks  []struct {
			Connector string `json:"connector"`
			Task      int    `json:"task"`
		} `json:"tasks"`
		Type string `json:"type"`
	} `json:"info"`
}

type ConnectorsResponseNamesOnly struct {
	// nothing else, just names
}

type ConnectorsResponseStatus struct {
	Status struct {
		Name      string `json:"name"`
		Connector struct {
			State    string `json:"state"`
			WorkerId string `json:"worker_id"`
		} `json:"connector"`
		Tasks []struct {
			Id       int    `json:"id"`
			State    string `json:"state"`
			WorkerId string `json:"worker_id"`
		} `json:"tasks"`
		Type string `json:"type"`
	} `json:"status"`
}

type Expand int8

const (
	None Expand = iota
	Status
	Info
)

func (e Expand) Name() string {
	switch e {
	case Status:
		return "Status"
	case Info:
		return "Info"
	}
	return "None"

}

func (e Expand) UrlParam() string {
	switch e {
	case Status:
		return "?expand=status"
	case Info:
		return "?expand=info"
	}
	return "" // nothing
}

func ExpandFromString(name string) Expand {
	switch name {
	case "Status":
		return Status
	case "Info":
		return Info
	}
	return None
}

func (e Expand) Unmarshal(bytes []byte) CPResponse {
	var response CPResponse
	switch e {
	case Status:
	case Info:
	default:
	}
	return response
}
