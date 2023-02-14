package lib

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
)

func BuildHealthFlagSet() HealthRequest {
	healthCmd := flag.NewFlagSet("health-check", flag.ExitOnError)
	healthHost := healthCmd.String("host", "http://localhost:8083", "cluster host")
	return HealthRequest{
		Command: healthCmd,
		Host:    healthHost,
	}
}

type HealthRequest struct {
	Command *flag.FlagSet
	Host    *string
}

func Health(request HealthRequest) {

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

	for connector, status := range statusMap {
		var states StateTracker
		states.StateMap = make(map[string]float64)
		fmt.Printf("%s: %s\n", connector, status.Status.Connector.State)
		for _, task := range status.Status.Tasks {
			value := states.StateMap[task.State]
			states.Total++
			value += 1
			states.StateMap[task.State] = value
			//fmt.Printf("- %d : %s\n", task.Id, task.State)
		}
		//fmt.Printf("states: %f\n", states.Total)
		for key, f := range states.StateMap {
			fmt.Printf("\t- %s: %.2f%s (%.0f of %.0f tasks) \n", key, 100*(f/states.Total), "%", f, states.Total)
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
