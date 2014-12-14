package main

import (
	"main/model"
	"sync"
	"fmt"
)

const STRING = "STRING"
const INT = "INT"
const FLOAT = "FLOAT"
const BOOL = "BOOL"
const NONE = "NONE"

func main() {

	// initialize the business process elements
	businessProcess := model.InitBusinessProcess(1, "Default")
	startEvent := model.InitStartEventElement()
	defaultActivity01 := model.InitDefaultActivityElement("Activity 0-1")
	defaultActivity02 := model.InitDefaultActivityElement("Activity 0-2")
	defaultActivity03 := model.InitDefaultActivityElement("Activity 0-3")
	defaultActivity04 := model.InitDefaultActivityElement("Activity 0-4")

	activeActivity01 := model.InitActiveActivityElement("Active Activity 0-1", func() (string, string) {
		var result = 4 * 6;
		
		return fmt.Sprintf("%d", result), INT
	})

	inclusiveGateway1 := model.InitInclusiveGateway(1, "Inclusive Gateway")
	endEvent := model.InitEndEventElement()

	// ================================
	// initialize the first sub process
	subProcess1 := model.InitBusinessProcess(2, "Sub Process 1")
	defaultActivity11 := model.InitDefaultActivityElement("Activity 1-1")
	defaultActivity12 := model.InitDefaultActivityElement("Activity 1-2")
	defaultActivity13 := model.InitDefaultActivityElement("Activity 1-3")
	defaultActivity14 := model.InitDefaultActivityElement("Activity 1-4")

	model.AddDefaultActivity(subProcess1, defaultActivity11)
	model.AddDefaultActivity(subProcess1, defaultActivity12)
	model.AddDefaultActivity(subProcess1, defaultActivity13)
	model.AddDefaultActivity(subProcess1, defaultActivity14)

	// =================================
	// initialize the second sub process
	subProcess2 := model.InitBusinessProcess(3, "Sub Process 2")
	defaultActivity21 := model.InitDefaultActivityElement("Activity 2-1")
	defaultActivity22 := model.InitDefaultActivityElement("Activity 2-2")
	defaultActivity23 := model.InitDefaultActivityElement("Activity 2-3")
	//defaultActivity24 := model.InitDefaultActivityElement("Activity 2-4")

	model.AddDefaultActivity(subProcess2, defaultActivity21)
	model.AddDefaultActivity(subProcess2, defaultActivity22)
	model.AddDefaultActivity(subProcess2, defaultActivity23)
	//model.AddDefaultActivity(subProcess2, defaultActivity24)

	// =================================
	// create the business process
	model.AddStartEvent(businessProcess, startEvent)
	model.AddDefaultActivity(businessProcess, defaultActivity01)
	model.AddDefaultActivity(businessProcess, defaultActivity02)
	model.AddActiveActivity(businessProcess, activeActivity01)

	// Add the inclusive gateway and the sub processes
	model.AddInclusiveGateway(businessProcess, inclusiveGateway1)
	model.AddSubProcessInclusiveGateway(inclusiveGateway1, ">;30", subProcess1)
	model.AddSubProcessInclusiveGateway(inclusiveGateway1, "<;25", subProcess2)

	model.AddDefaultActivity(businessProcess, defaultActivity03)
	model.AddDefaultActivity(businessProcess, defaultActivity04)
	model.AddEndEvent(businessProcess, endEvent)

	var waitGroup sync.WaitGroup
	waitGroup.Add(1)
	businessProcess.Simulate(businessProcess, &waitGroup)
	waitGroup.Wait()
}
