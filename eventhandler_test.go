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

func initializeTestObjects(eventFileName string) (*keptn.Keptn, *cloudevents.Event, error) {
	// load sample event
	eventFile, err := ioutil.ReadFile(eventFileName)
	if err != nil {
		return nil, nil, fmt.Errorf("Cant load test-events/get-sli.http: " + err.Error())
	}

	incomingEvent := &cloudevents.Event{}
	err = json.Unmarshal(eventFile, incomingEvent)
	if err != nil {
		return nil, nil, fmt.Errorf("Error parsing: " + err.Error())
	}

	var keptnOptions = keptn.KeptnOpts{}
	keptnOptions.UseLocalFileSystem = true
	myKeptn, err := keptn.NewKeptn(incomingEvent, keptnOptions)

	return myKeptn, incomingEvent, err
}

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
