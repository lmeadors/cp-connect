package lib

import (
	"flag"
	"fmt"
	"os"
)

func HandleArguments() {

	var commands []flag.FlagSet

	clusterRequest := BuildClusterFlagSet()
	commands = append(commands, *clusterRequest.Command)

	connectorsRequest := BuildConnectorsFlagSet()
	commands = append(commands, *connectorsRequest.Command)

	healthRequest := BuildHealthFlagSet()
	commands = append(commands, *healthRequest.Command)

	connectorCmd := flag.NewFlagSet("connector", flag.ExitOnError)
	connectorHost := connectorCmd.String("host", "http://localhost:8083", "cluster host")
	connectorJson := connectorCmd.String("json", "", "json configuration file")
	connectorAction := connectorCmd.String("action", "Config", "action to perform (Config | Validate | Put | Status | Pause | Resume | Delete)")

	commands = append(commands, *connectorCmd)

	if len(os.Args) < 2 {
		fmt.Println("expected one for the subcommands:")
		for i := range commands {
			println(commands[i].Name())
			commands[i].PrintDefaults()
			//println(commands[i].Args())
		}
		os.Exit(1)
	}

	switch os.Args[1] {
	case "cluster":
		clusterRequest.Command.Parse(os.Args[2:])
		Cluster(clusterRequest)
	case "health-check":
		healthRequest.Command.Parse(os.Args[2:])
		Health(healthRequest)
	case "connector":
		connectorCmd.Parse(os.Args[2:])
		Connector(*connectorHost, *connectorJson, *connectorAction)
	case "connectors":
		connectorsRequest.Command.Parse(os.Args[2:])
		Connectors(connectorsRequest)
	}

}
