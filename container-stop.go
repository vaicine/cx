package main

import (
	"os"
	"text/tabwriter"
	"time"

	"github.com/cloud66/cloud66"

	"github.com/cloud66/cli"
)

func runContainerStop(c *cli.Context) {
	if len(c.Args()) == 0 {
		cli.ShowSubcommandHelp(c)
		os.Exit(2)
	}

	stack := mustStack(c)
	w := tabwriter.NewWriter(os.Stdout, 1, 2, 2, ' ', 0)
	defer w.Flush()

	containerUid := c.Args()[0]
	container, err := client.GetContainer(stack.Uid, containerUid)
	must(err)

	if container == nil {
		printFatal("Container with Id '" + containerUid + "' not found")
	}

	asyncId, err := startContainerStop(stack.Uid, containerUid)
	if err != nil {
		printFatal(err.Error())
	}
	genericRes, err := endServerSet(*asyncId, stack.Uid)
	if err != nil {
		printFatal(err.Error())
	}
	printGenericResponse(*genericRes)
	return
}

func startContainerStop(stackUid string, containerUid string) (*int, error) {
	asyncRes, err := client.StopContainer(stackUid, containerUid)
	if err != nil {
		return nil, err
	}
	return &asyncRes.Id, err
}

func endContainerStop(asyncId int, stackUid string) (*cloud66.GenericResponse, error) {
	return client.WaitStackAsyncAction(asyncId, stackUid, 3*time.Second, 20*time.Minute, true)
}
