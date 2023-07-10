package lib

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
)

func BuildHealthFlagSet(cmdMap map[string]CpCommand) HealthRequest {
	name := "health-check"
	healthCmd := flag.NewFlagSet(name, flag.ExitOnError)
	healthHost := healthCmd.String("host", "http://localhost:8083", "cluster host")
	request := HealthRequest{
		Command: healthCmd,
		Name:    name,
		Host:    healthHost,
	}
	cmdMap[name] = request
	return request
}

type HealthRequest struct {
	Command *flag.FlagSet
	Name    string
	Host    *string
}

func (request HealthRequest) GetCommand() flag.FlagSet {
	return *request.Command
}

func (request HealthRequest) GetName() string {
	return request.Name
}

func (request HealthRequest) PrintDefaults() {
	request.Command.PrintDefaults()
}

func (request HealthRequest) Execute() {
	//Health(request)
	//}
	//
	//func Health(request HealthRequest) {

	//logger := log.Default()
	//expand := Status

	//logger.Printf("host:   %s\n", *request.Host)
	//logger.Printf("expand: %s\n", expand.Name())
	//logger.Printf("uri:    %s\n", url)

	url := *request.Host + "/connectors?expand=status"

	// build request
	req, err := http.NewRequest("GET", url, nil)
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
	var statusMap map[string]ConnectorStatus
	err2 := json.Unmarshal(bodyBytes, &statusMap)
	if err2 != nil {
		fmt.Print(err.Error())
	}

	var colorMap = map[string]Color{
		"RUNNING":    ColorGreen,
		"UNASSIGNED": ColorCyan,
		"PAUSED":     ColorYellow,
		"FAILED":     ColorRed,
	}

	for connector, status := range statusMap {
		var states StateTracker
		states.StateMap = make(map[string]float64)
		fmt.Printf("%s: %s%s%s\n", connector, colorMap[status.Status.Connector.State], status.Status.Connector.State, ColorReset)
		for _, task := range status.Status.Tasks {
			value := states.StateMap[task.State]
			states.Total++
			value += 1
			states.StateMap[task.State] = value
			//fmt.Printf("- %d : %s\n", task.Id, task.State)
		}
		//fmt.Printf("states: %f\n", states.Total)
		for key, f := range states.StateMap {
			fmt.Printf("  - %s%s:%s %.2f%s (%.0f of %.0f tasks) \n", colorMap[key], key, ColorReset, 100*(f/states.Total), "%", f, states.Total)
		}
	}

}

type StateTracker struct {
	StateMap map[string]float64
	Total    float64
}

type ConnectorStatus struct {
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
