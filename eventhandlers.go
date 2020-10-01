package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	keptn "github.com/keptn/go-utils/pkg/lib"

	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
	"github.com/google/uuid"
)

type PACResultFile struct {
	Results []PACResult `json:"results"`
}

type PACResult struct {
	ID   string             `json:"id"`
	URL  string             `json:"url"`
	Date string             `json:"date"`
	Data map[string]float64 `json:"data"`
}

/**
* Here are all the handler functions for the individual event
  See https://github.com/keptn/spec/blob/0.1.3/cloudevents.md for details on the payload

  -> "sh.keptn.event.configuration.change"
  -> "sh.keptn.events.deployment-finished"
  -> "sh.keptn.events.tests-finished"
  -> "sh.keptn.event.start-evaluation"
  -> "sh.keptn.events.evaluation-done"
  -> "sh.keptn.event.problem.open"
	-> "sh.keptn.events.problem"
	-> "sh.keptn.event.action.triggered"
*/

// Handles ConfigureMonitoringEventType = "sh.keptn.event.monitoring.configure"
func HandleConfigureMonitoringEvent(myKeptn *keptn.Keptn, incomingEvent cloudevents.Event, data *keptn.ConfigureMonitoringEventData) error {
	log.Printf("Handling Configure Monitoring Event: %s", incomingEvent.Context.GetID())

	return nil
}

//
// Handles ConfigurationChangeEventType = "sh.keptn.event.configuration.change"
// TODO: add in your handler code
//
func HandleConfigurationChangeEvent(myKeptn *keptn.Keptn, incomingEvent cloudevents.Event, data *keptn.ConfigurationChangeEventData) error {
	log.Printf("Handling Configuration Changed Event: %s", incomingEvent.Context.GetID())

	return nil
}

//
// Handles DeploymentFinishedEventType = "sh.keptn.events.deployment-finished"
// TODO: add in your handler code
//
func HandleDeploymentFinishedEvent(myKeptn *keptn.Keptn, incomingEvent cloudevents.Event, data *keptn.DeploymentFinishedEventData) error {
	log.Printf("Handling Deployment Finished Event: %s", incomingEvent.Context.GetID())

	// capture start time for tests
	// startTime := time.Now()

	// run tests
	// ToDo: Implement your tests here

	// Send Test Finished Event
	// return myKeptn.SendTestsFinishedEvent(&incomingEvent, "", "", startTime, "pass", nil, "pac-sliprovider")
	return nil
}

//
// Handles TestsFinishedEventType = "sh.keptn.events.tests-finished"
// TODO: add in your handler code
//
func HandleTestsFinishedEvent(myKeptn *keptn.Keptn, incomingEvent cloudevents.Event, data *keptn.TestsFinishedEventData) error {
	log.Printf("Handling Tests Finished Event: %s", incomingEvent.Context.GetID())

	return nil
}

//
// Handles EvaluationDoneEventType = "sh.keptn.events.evaluation-done"
// TODO: add in your handler code
//
func HandleStartEvaluationEvent(myKeptn *keptn.Keptn, incomingEvent cloudevents.Event, data *keptn.StartEvaluationEventData) error {
	log.Printf("Handling Start Evaluation Event: %s", incomingEvent.Context.GetID())

	return nil
}

//
// Handles DeploymentFinishedEventType = "sh.keptn.events.deployment-finished"
// TODO: add in your handler code
//
func HandleEvaluationDoneEvent(myKeptn *keptn.Keptn, incomingEvent cloudevents.Event, data *keptn.EvaluationDoneEventData) error {
	log.Printf("Handling Evaluation Done Event: %s", incomingEvent.Context.GetID())

	return nil
}

//
// Handles InternalGetSLIEventType = "sh.keptn.internal.event.get-sli"
// TODO: add in your handler code
//
func HandleInternalGetSLIEvent(myKeptn *keptn.Keptn, incomingEvent cloudevents.Event, data *keptn.InternalGetSLIEventData) error {
	log.Printf("Handling Internal Get SLI Event: %s", incomingEvent.Context.GetID())

	incomingGetSLIEventData := data

	// Step 1 - Do we need to do something?
	// Lets make sure we are only processing an event that really belongs to our SLI Provider
	if incomingGetSLIEventData.SLIProvider != "pac-sliprovider" {
		return fmt.Errorf("Not handling event because its not for pac-sliprovider. Its for: " + incomingGetSLIEventData.SLIProvider)
	}

	// Step 2 - prep-work
	// Get the incoming GetSLIEvent for accessing incoming data
	// Labels: get the incoming labels for potential config data and use it to pass more labels on result, e.g: links
	// Indicators: this is the list of indicators as requested in the SLO.yaml
	// SLIResult: this is the array that will receive the results
	labels := incomingGetSLIEventData.Labels
	if labels == nil {
		labels = make(map[string]string)
	}

	// Step 3 - query the data entry that matches what the user wants
	indicators := incomingGetSLIEventData.Indicators
	sliResults := []*keptn.SLIResult{}
	pacResult, resultFile, err := loadPACData(myKeptn, incomingGetSLIEventData, labels)

	// if we dont have any data return the error
	if err != nil || pacResult == nil {
		SendInternalGetSLIDoneEvent(myKeptn, incomingGetSLIEventData, indicators, sliResults, labels, err, "pac-sliprovider")
		return err
	}

	// Step 4 - now lets return those result properties
	for _, indicatorName := range indicators {

		pacValue, pacValueExists := pacResult.Data[indicatorName]
		valueMessage := ""
		if !pacValueExists {
			valueMessage = "Couldnt find data for SLI " + indicatorName
		}

		sliResult := &keptn.SLIResult{
			Metric:  indicatorName,
			Value:   pacValue,
			Success: pacValueExists,
			Message: valueMessage,
		}

		sliResults = append(sliResults, sliResult)
	}

	// Step 5 - add the result file link to the labels so it shows up in the bridage
	labels["PAC Data Source"] = resultFile
	labels["Link to "+pacResult.ID] = pacResult.URL

	sliResultAsText, err := json.Marshal(sliResults)
	log.Printf(string(sliResultAsText))

	return SendInternalGetSLIDoneEvent(myKeptn, incomingGetSLIEventData, indicators, sliResults, labels, nil, "pac-sliprovider")
}

//
// Handles ProblemOpenEventType = "sh.keptn.event.problem.open"
// Handles ProblemEventType = "sh.keptn.events.problem"
// TODO: add in your handler code
//
func HandleProblemEvent(myKeptn *keptn.Keptn, incomingEvent cloudevents.Event, data *keptn.ProblemEventData) error {
	log.Printf("Handling Problem Event: %s", incomingEvent.Context.GetID())

	// Deprecated since Keptn 0.7.0 - use the HandleActionTriggeredEvent instead

	return nil
}

//
// Handles ActionTriggeredEventType = "sh.keptn.event.action.triggered"
// TODO: add in your handler code
//
func HandleActionTriggeredEvent(myKeptn *keptn.Keptn, incomingEvent cloudevents.Event, data *keptn.ActionTriggeredEventData) error {
	log.Printf("Handling Action Triggered Event: %s", incomingEvent.Context.GetID())

	// check if action is supported
	if data.Action.Action == "action-xyz" {
		//myKeptn.SendActionStartedEvent()

		// Implement your remediation action here

		//myKeptn.SendActionFinishedEvent()
	}
	return nil
}

//
// This actually loads the data from the results file
//
func loadPACData(myKeptn *keptn.Keptn, incomingGetSLIEventData *keptn.InternalGetSLIEventData, labels map[string]string) (*PACResult, string, error) {

	// Step 1 - figure out where to load our data from
	resultFile := ""

	// Step 1a - option 1 could be e.g: from customFilters which is data from the SLO.yaml
	if incomingGetSLIEventData.CustomFilters != nil {
		for _, customFilter := range incomingGetSLIEventData.CustomFilters {
			if customFilter.Key == "resultfile" {
				resultFile = customFilter.Value
			}
		}
	}

	// Step 1b - just to give you a different example - loading some configuration from an environment variable on the containers if it wasnt set through filters
	if resultFile == "" {
		resultFile = os.Getenv("RESULTFILE")
	}

	// Step 1c - another option is loading it from a sli-provider specific configuration file stored in the Keptn Repo
	// This could be a sli.yaml, myconfig.conf, ... - whatever makes sense for your SLI provider
	// In this case we load it from a file in a sli-provider specific subfolder that is stored for this keptn project, stage & service
	sliConfigFile, err := myKeptn.GetKeptnResource("pac-sliprovider/sli.yaml")
	if err != nil {
		log.Printf("TODO: Downloaded SLI.yaml but not doing anything with it yet in this sample SLI Provider: " + sliConfigFile)
		// resultFile = keptnResourceContent
	}

	// Step 1d - if no configuration was set through customFilters, Envs or Config File lets go with a default!
	if resultFile == "" {
		resultFile = "https://raw.githubusercontent.com/grabnerandi/pac-sliprovider/master/resultfiles/results.json"
	}

	// Step 2 - Lets actually download the result file from the URL we have:-)
	resp, err := http.Get(resultFile)
	if err != nil {
		return nil, resultFile, err
	}
	defer resp.Body.Close()
	resultFileContent, _ := ioutil.ReadAll(resp.Body)
	resultJSON := &PACResultFile{}
	err = json.Unmarshal(resultFileContent, &resultJSON)
	if err != nil {
		return nil, resultFile, err
	}
	if resultJSON == nil || len(resultJSON.Results) == 0 {
		return nil, resultFile, fmt.Errorf("No results in pacresult.json")
	}

	// Step 3 - query the data the user actually wants
	// Depending on the SLI service we could use the timeframe (start/end) to e.g: query timeranges, we could also use information passed via labels, e.g: a test run id
	// In our example we simply look for a label with the name "pacId" - if that is present we return the data we have for that pacId. Otherwise we just return the first result
	pacID := labels["pacId"]
	var matchingPacResult *PACResult
	if pacID == "" {
		matchingPacResult = &resultJSON.Results[0]
	} else {
		for _, pacResult := range resultJSON.Results {
			if pacResult.ID == pacID {
				matchingPacResult = &pacResult
			}
		}
	}

	if matchingPacResult == nil {
		return nil, resultFile, fmt.Errorf("Couldn't find matching results")
	}

	return matchingPacResult, resultFile, nil
}

/**
 * Sends the SLI Done Event. If err != nil it will send an error message
 * THIS SHOULD MAKE IT INTO go-utils!
 */
func SendInternalGetSLIDoneEvent(myKeptn *keptn.Keptn, incomingGetSLIEventData *keptn.InternalGetSLIEventData, indicators []string, indicatorValues []*keptn.SLIResult, labels map[string]string, err error, eventSource string) error {

	source, _ := url.Parse(eventSource)
	contentType := "application/json"

	// if an error was set - the indicators will be set to failed and error message is set to each
	if err != nil {
		errMessage := err.Error()

		if (indicatorValues == nil) || (len(indicatorValues) == 0) {
			if indicators == nil || len(indicators) == 0 {
				indicators = []string{"no metric"}
			}

			for _, indicatorName := range indicators {
				indicatorValues = []*keptn.SLIResult{
					{
						Metric: indicatorName,
						Value:  0.0,
					},
				}
			}
		}

		for _, indicator := range indicatorValues {
			indicator.Success = false
			indicator.Message = errMessage
		}
	}

	sliDoneEvent := keptn.InternalGetSLIDoneEventData{}

	// reuse data from the incoming GetSLIEventData
	if incomingGetSLIEventData != nil {
		sliDoneEvent.Project = incomingGetSLIEventData.Project
		sliDoneEvent.Stage = incomingGetSLIEventData.Stage
		sliDoneEvent.Service = incomingGetSLIEventData.Service

		sliDoneEvent.Start = incomingGetSLIEventData.Start
		sliDoneEvent.End = incomingGetSLIEventData.End
		sliDoneEvent.TestStrategy = incomingGetSLIEventData.TestStrategy
		sliDoneEvent.DeploymentStrategy = incomingGetSLIEventData.DeploymentStrategy
		sliDoneEvent.Deployment = incomingGetSLIEventData.Deployment
	}

	if labels != nil {
		sliDoneEvent.Labels = labels
	}
	if indicatorValues != nil {
		sliDoneEvent.IndicatorValues = indicatorValues
	}

	event := cloudevents.Event{
		Context: cloudevents.EventContextV02{
			ID:          uuid.New().String(),
			Time:        &types.Timestamp{Time: time.Now()},
			Type:        keptn.InternalGetSLIDoneEventType,
			Source:      types.URLRef{URL: *source},
			ContentType: &contentType,
			Extensions:  map[string]interface{}{"shkeptncontext": myKeptn.KeptnContext},
		}.AsV02(),
		Data: sliDoneEvent,
	}

	log.Println(fmt.Printf("%s", event))

	return myKeptn.SendCloudEvent(event)
}
