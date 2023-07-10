package lib

import (
	"flag"
	"fmt"
	"os"
)

type CpCommand interface {
	GetCommand() flag.FlagSet
	GetName() string
	PrintDefaults()
	Execute()
}

func HandleArguments() {

	var cmdMap = make(map[string]CpCommand)

	BuildClusterFlagSet(cmdMap)
	BuildConnectorsFlagSet(cmdMap)
	BuildHealthFlagSet(cmdMap)
	BuildConnectorFlagSet(cmdMap)
	BuildPluginsRequest(cmdMap)

	if len(os.Args) < 2 {
		fmt.Println("expected one of the subcommands below:")
		showHelp(cmdMap)
		os.Exit(1)
	}

	r, ok := cmdMap[os.Args[1]]
	if ok {
		command := r.GetCommand()
		command.Parse(os.Args[2:])
		r.Execute()
	} else {
		showHelp(cmdMap)
	}

}

func showHelp(cmdMap map[string]CpCommand) {
	for _, cmd := range cmdMap {
		println(cmd.GetName())
		cmd.PrintDefaults()
	}
}
