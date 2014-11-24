package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/cloudfoundry-incubator/receptor"
	"github.com/pivotal-cf-experimental/veritas/say"
)

func report(tasks []receptor.TaskResponse, desiredLRPs []receptor.DesiredLRPResponse, actualLRPs []receptor.ActualLRPResponse) {
	reportOnTasks(tasks)
	reportOnLRPs(desiredLRPs, actualLRPs)
}

func reportOnTasks(tasks []receptor.TaskResponse) {
	say.Println(0, say.Green("Tasks"))
	if len(tasks) == 0 {
		say.Println(1, say.Red("No Tasks"))
	}

	for _, task := range tasks {
		state := fmt.Sprintf("[%s]", task.State)
		switch task.State {
		case receptor.TaskStatePending:
			state = say.LightGray(state)
		case receptor.TaskStateClaimed:
			state = say.Yellow(state)
		case receptor.TaskStateRunning:
			state = say.Cyan(state)
		case receptor.TaskStateCompleted:
			if task.Failed {
				state = say.Red(state)
			} else {
				state = say.Green(state)
			}
		case receptor.TaskStateResolving:
			state = say.Yellow(state)
		}
		say.Println(1,
			"%s %s on %s",
			state,
			task.TaskGuid,
			task.CellID,
		)
		if task.State == receptor.TaskStateCompleted {
			if task.Failed {
				say.Println(2, say.Red("Failed"))
				say.Println(2, say.Red(task.FailureReason))
			} else {
				say.Println(2, say.Green("Success"))
				say.Println(2, say.Green(task.Result))
			}
		}
	}
}

type LRP struct {
	DesiredLRP        receptor.DesiredLRPResponse
	ActualLRPsByIndex map[int][]receptor.ActualLRPResponse
}

func (l LRP) OrderedActualLRPIndices() []int {
	indices := []int{}
	for index := range l.ActualLRPsByIndex {
		indices = append(indices, index)
	}

	sort.Ints(indices)
	return indices
}

func reportOnLRPs(desiredLRPs []receptor.DesiredLRPResponse, actualLRPs []receptor.ActualLRPResponse) {
	say.Println(0, say.Green("LRPs"))

	processGuids := []string{}
	lrps := map[string]LRP{}
	for _, desiredLRP := range desiredLRPs {
		processGuids = append(processGuids, desiredLRP.ProcessGuid)
		lrps[desiredLRP.ProcessGuid] = LRP{
			DesiredLRP:        desiredLRP,
			ActualLRPsByIndex: map[int][]receptor.ActualLRPResponse{},
		}
	}

	for _, actualLRP := range actualLRPs {
		processGuid := actualLRP.ProcessGuid
		if _, ok := lrps[processGuid]; !ok {
			processGuids = append(processGuids, actualLRP.ProcessGuid)
			lrps[actualLRP.ProcessGuid] = LRP{
				ActualLRPsByIndex: map[int][]receptor.ActualLRPResponse{},
			}
		}
		lrps[processGuid].ActualLRPsByIndex[actualLRP.Index] = append(lrps[processGuid].ActualLRPsByIndex[actualLRP.Index], actualLRP)
	}

	sort.Strings(processGuids)

	if len(processGuids) == 0 {
		say.Println(1, say.Red("No LRPs"))
	}

	for _, processGuid := range processGuids {
		lrp := lrps[processGuid]
		if lrp.DesiredLRP.ProcessGuid == "" {
			say.Println(1, say.Red("%s - Undesired", processGuid))
		} else {
			routes := ""
			if len(lrp.DesiredLRP.Routes) > 0 {
				routes = say.Green(strings.Join(lrp.DesiredLRP.Routes, ","))
			}
			say.Println(1,
				"%s [%d] %s",
				say.Green(processGuid), lrp.DesiredLRP.Instances, routes)
		}
		indices := lrp.OrderedActualLRPIndices()
		for _, index := range indices {
			actuals := lrp.ActualLRPsByIndex[index]
			for _, actual := range actuals {
				state := actual.State
				switch actual.State {
				case receptor.ActualLRPStateStarting:
					state = say.Yellow(state)
				case receptor.ActualLRPStateRunning:
					state = say.Green(state)
				}
				ports := []string{}
				for _, port := range actual.Ports {
					ports = append(ports, fmt.Sprintf("%d:%d", port.HostPort, port.ContainerPort))
				}
				say.Println(2,
					"[%d] %s on %s - %s (%s)",
					actual.Index,
					state,
					actual.CellID,
					actual.Host,
					strings.Join(ports, ","),
				)
			}
		}
	}

}
