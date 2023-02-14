package lib

import "flag"

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

	// look for any connectors in any state other than "running"

}
