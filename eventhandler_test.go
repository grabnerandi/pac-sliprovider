package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	_ "github.com/keptn/go-utils/pkg/lib"
	keptn "github.com/keptn/go-utils/pkg/lib"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
)

/**
 * loads a cloud event from the passed test json file and initializes a keptn object with it
 */
func initializeTestObjects(eventFileName string) (*keptn.Keptn, *cloudevents.Event, error) {
	// load sample event
	eventFile, err := ioutil.ReadFile(eventFileName)
	if err != nil {
		return nil, nil, fmt.Errorf("Cant load %s: %s", eventFileName, err.Error())
	}

	incomingEvent := &cloudevents.Event{}
	err = json.Unmarshal(eventFile, incomingEvent)
	if err != nil {
		return nil, nil, fmt.Errorf("Error parsing: %s", err.Error())
	}

	var keptnOptions = keptn.KeptnOpts{}
	keptnOptions.UseLocalFileSystem = true
	myKeptn, err := keptn.NewKeptn(incomingEvent, keptnOptions)

	return myKeptn, incomingEvent, err
}

// Tests the InternalGetSLIEvent Handler
func TestHandleInternalGetSLIEvent(t *testing.T) {

	myKeptn, incomingEvent, err := initializeTestObjects("test-events/get-sli.json")
	if err != nil {
		t.Error(err)
		return
	}

	getSLIEventData := &keptn.InternalGetSLIEventData{}
	err = incomingEvent.DataAs(getSLIEventData)
	if err != nil {
		t.Errorf("Error getting keptn event data")
	}

	err = HandleInternalGetSLIEvent(myKeptn, *incomingEvent, getSLIEventData)

	if err != nil {
		t.Errorf("Error: " + err.Error())
	}
}
