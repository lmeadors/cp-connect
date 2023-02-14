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

	connectorRequest := BuildConnectorFlagSet()
	commands = append(commands, *connectorRequest.Command)

	if len(os.Args) < 2 {
		showHelp(commands)
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
		connectorRequest.Command.Parse(os.Args[2:])
		Connector(connectorRequest)
	case "connectors":
		connectorsRequest.Command.Parse(os.Args[2:])
		Connectors(connectorsRequest)
	default:
		showHelp(commands)
	}

}

func showHelp(commands []flag.FlagSet) {
	fmt.Println("expected one for the subcommands:")
	for i := range commands {
		println(commands[i].Name())
		commands[i].PrintDefaults()
		//println(commands[i].Args())
	}
}
