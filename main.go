package main

import (
	"fmt"
	"os"

	"github.com/cloudfoundry-incubator/receptor"
)

func main() {
	if len(os.Args) != 2 {
		PrintUsageAndExit()
	}

	domain := os.Args[1]

	receptorAddr := os.Getenv("RECEPTOR")
	if receptorAddr == "" {
		fmt.Println("No RECEPTOR set")
		PrintUsageAndExit()
	}

	client := receptor.NewClient(receptorAddr)
	tasks, err := client.TasksByDomain(domain)
	ExitIfError(err)
	desiredLRPs, err := client.DesiredLRPsByDomain(domain)
	ExitIfError(err)
	actualLRPs, err := client.ActualLRPsByDomain(domain)
	ExitIfError(err)

	report(tasks, desiredLRPs, actualLRPs)
}

func PrintUsageAndExit() {
	fmt.Println(`Usage:
troy DOMAIN

Set the receptor address with the RECEPTOR environment:
    export RECEPTOR=http://username:password@receptor.ketchup.cf-app.com

The address for a local Diego Edge box can be set via: 
    export RECEPTOR=http://receptor.192.168.11.11.xip.io
`)
	os.Exit(1)
}

func ExitIfError(err error) {
	if err != nil {
		fmt.Printf("Got an unexpected error:\n%s\n", err.Error())
		os.Exit(1)
	}
}
