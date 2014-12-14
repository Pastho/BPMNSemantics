package model

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
	rand "math/rand"
)

const STRING = "STRING"
const INT = "INT"
const FLOAT = "FLOAT"
const BOOL = "BOOL"
const NONE = "NONE"

type Simulator interface {
	Simulate(previousElement Simulator, waitGroup *sync.WaitGroup) string
	Result() (string, string)
}

type Element struct {
	id          int
	description string
}

/***************************************
 define the activities
***************************************/

type Activity struct {
	Element
	startTime, endTime time.Time
}

type DefaultActivityElement struct {
	Activity
}

func InitDefaultActivityElement(text string) *DefaultActivityElement {
	defaultActivity := new(DefaultActivityElement)

	defaultActivity.Activity.description = text

	return defaultActivity
}

func (defaultActivity *DefaultActivityElement) Simulate(previousElement Simulator, waitGroup *sync.WaitGroup) string {
	time.Sleep(1 * time.Second)
	waitGroup.Done()
	return " ---> " + defaultActivity.description
}

func (defaultActivity DefaultActivityElement) Result() (string, string) {
	return NONE, STRING
}

type ActiveActivityElement struct {
	Activity
	result, dataType string
	function         func() (string, string)
}

func InitActiveActivityElement(text string, function func() (string, string)) *ActiveActivityElement {
	activeActivityElement := new(ActiveActivityElement)

	activeActivityElement.Activity.description = text
	activeActivityElement.function = function

	return activeActivityElement
}

func (activeActivityElement *ActiveActivityElement) Simulate(previousElement Simulator, waitGroup *sync.WaitGroup) string {
	time.Sleep(time.Duration((rand.Intn(3) + 1)) * time.Second)
	activeActivityElement.result, activeActivityElement.dataType = activeActivityElement.function()
	waitGroup.Done()
	return " ---> " + activeActivityElement.description + " (" + activeActivityElement.result + ")"
}

func (activeActivityElement ActiveActivityElement) Result() (string, string) {

	return activeActivityElement.result, activeActivityElement.dataType
}

/***************************************
 define the events
***************************************/

type Event struct {
	Element
}

type StartEventElement struct {
	Event
}

func InitStartEventElement() *StartEventElement {
	start := new(StartEventElement)
	start.Event.description = "Start Event"
	return start
}

func (startEventElement *StartEventElement) Simulate(previousElement Simulator, waitGroup *sync.WaitGroup) string {
	fmt.Println("Start of business process simulation...")
	rand.Seed( time.Now().UTC().UnixNano())
	waitGroup.Done()
	return startEventElement.description
}

func (startEventElement StartEventElement) Result() (string, string) {
	return NONE, STRING
}

type EndEventElement struct {
	Event
}

func InitEndEventElement() *EndEventElement {
	end := new(EndEventElement)
	end.Event.description = "End Event"
	return end
}

func (endEventElement *EndEventElement) Simulate(previousElement Simulator, waitGroup *sync.WaitGroup) string {
	waitGroup.Done()
	return " ----> End event"
}

func (endEventElement EndEventElement) Result() (string, string) {
	return NONE, STRING
}

/***************************************
 define the business process element
***************************************/

type BusinessProcess struct {
	Element
	processElements []Simulator
}

func InitBusinessProcess(id int, description string) *BusinessProcess {
	businessProcess := new(BusinessProcess)
	businessProcess.Element.description = description
	return businessProcess
}

func (businessProcess *BusinessProcess) Simulate(previousElement Simulator, waitGroup *sync.WaitGroup) string {

	var subWaitGroup sync.WaitGroup

	for i := range businessProcess.processElements {
		subWaitGroup.Add(1)
		if i > 0 {
			fmt.Print(businessProcess.processElements[i].Simulate(businessProcess.processElements[i-1], &subWaitGroup))
		} else {
			fmt.Print(businessProcess.processElements[i].Simulate(businessProcess.processElements[i], &subWaitGroup))
		}
	}

	waitGroup.Done()

	return ""
}

func (businessProcess BusinessProcess) Result() (string, string) {
	return NONE, STRING
}

func AddStartEvent(businessProcess *BusinessProcess, startEventElement *StartEventElement) {
	startEventElement.id = len(businessProcess.processElements)
	businessProcess.processElements = append(businessProcess.processElements, startEventElement)
}

func AddDefaultActivity(businessProcess *BusinessProcess, defaultActivity *DefaultActivityElement) {
	defaultActivity.id = len(businessProcess.processElements)
	businessProcess.processElements = append(businessProcess.processElements, defaultActivity)
}

func AddActiveActivity(businessProcess *BusinessProcess, activeActivity *ActiveActivityElement) {
	activeActivity.id = len(businessProcess.processElements)
	businessProcess.processElements = append(businessProcess.processElements, activeActivity)
}

func AddEndEvent(businessProcess *BusinessProcess, endEventElement *EndEventElement) {
	endEventElement.id = len(businessProcess.processElements)
	businessProcess.processElements = append(businessProcess.processElements, endEventElement)
}

func AddInclusiveGateway(businessProcess *BusinessProcess, inclusiveGateway *InclusiveGateway) {
	inclusiveGateway.id = len(businessProcess.processElements)
	businessProcess.processElements = append(businessProcess.processElements, inclusiveGateway)
}

/***************************************
 define the gateways
***************************************/

type Gateway struct {
	Element
	conditions   []string
	subProcesses []Simulator
}

type ExclusiveGateway struct {
	Gateway
}

type ParallelGateway struct {
	Gateway
}

func InitParallelGateway(id int, description string) *ParallelGateway {
	parallelGateway := new(ParallelGateway)
	parallelGateway.Gateway.id = id
	parallelGateway.Gateway.description = description
	return parallelGateway
}

func AddSubProcessParallelGateway(parallelGateway *ParallelGateway, businessProcess *BusinessProcess) {
	parallelGateway.subProcesses = append(parallelGateway.subProcesses, businessProcess)
}

func (parallelGateway *ParallelGateway) Simulate(previousElement Simulator, waitGroup *sync.WaitGroup) string {

	fmt.Print("\n\n\t\t\t\t")

	var subWaitGroup sync.WaitGroup

	for i := range parallelGateway.subProcesses {
		subWaitGroup.Add(1)
		go parallelGateway.subProcesses[i].Simulate(parallelGateway, &subWaitGroup)
	}
	
	subWaitGroup.Wait()
	waitGroup.Done()
	fmt.Print("\n\n")

	return ""
}

func (parallelGateway ParallelGateway) Result() (string, string) {
	return NONE, STRING
}

type InclusiveGateway struct {
	Gateway
}

func InitInclusiveGateway(id int, description string) *InclusiveGateway {
	inclusiveGateway := new(InclusiveGateway)
	inclusiveGateway.Gateway.id = id
	inclusiveGateway.Gateway.description = description
	return inclusiveGateway
}

func AddSubProcessInclusiveGateway(inclusiveGateway *InclusiveGateway, condition string, businessProcess *BusinessProcess) {
	inclusiveGateway.conditions = append(inclusiveGateway.conditions, condition)
	inclusiveGateway.subProcesses = append(inclusiveGateway.subProcesses, businessProcess)
}

func (inclusiveGateway *InclusiveGateway) Simulate(previousElement Simulator, waitGroup *sync.WaitGroup) string {

	//get result of previous element
	var previousResult, dataType = previousElement.Result()

	var result, err interface{}

	// Parse result
	switch dataType {
	case INT:
		result, err = strconv.ParseInt(previousResult, 10, 32)
		break
	case BOOL:
		result, err = strconv.ParseInt(previousResult, 10, 32)
		break
	case FLOAT:
		result, err = strconv.ParseInt(previousResult, 10, 32)
		break
	case STRING:
		result = previousResult
		break
	default:
		result = ""
		break
	}

	if result != nil && err == nil {

		fmt.Print("\n\n\t\t\t\t")

		var subWaitGroup sync.WaitGroup

		for i := range inclusiveGateway.conditions {
			var condition = strings.Split(inclusiveGateway.conditions[i], ";")

			var conditionResult, conditionERR interface{}
			conditionResult, conditionERR = strconv.ParseInt(condition[1], 10, 32)

			if conditionResult != nil && conditionERR == nil {
				switch condition[0] {
				case "<":
					if result.(int64) < conditionResult.(int64) {
						subWaitGroup.Add(1)
						go inclusiveGateway.subProcesses[i].Simulate(inclusiveGateway, &subWaitGroup)
					}
					break
				case ">":
					if result.(int64) > conditionResult.(int64) {
						subWaitGroup.Add(1)
						go inclusiveGateway.subProcesses[i].Simulate(inclusiveGateway, &subWaitGroup)
					}
					break
				case "=":
					if result == conditionResult {
						subWaitGroup.Add(1)
						go inclusiveGateway.subProcesses[i].Simulate(inclusiveGateway, &subWaitGroup)
					}
					break
				}
			}
			
		}
		
		subWaitGroup.Wait()
		waitGroup.Done()
		fmt.Print("\n\n")
	}

	return ""
}

func (inclusiveGateway InclusiveGateway) Result() (string, string) {
	return NONE, STRING
}
