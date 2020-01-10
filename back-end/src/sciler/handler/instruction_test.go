package handler

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"sciler/config"
	"testing"
	"time"
)

////////////////////////////// Instruction tests //////////////////////////////
func TestInstructionSetUp(t *testing.T) {
	communicatorMock := new(CommunicatorMock)
	workingConfig := config.ReadFile("../../../resources/testing/test_setup.json")
	handler := Handler{
		Config:       workingConfig,
		ConfigFile:   "../../../resources/testing/test_setup.json",
		Communicator: communicatorMock,
	}
	instructionMsg := Message{
		DeviceID: "front-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "instruction",
		Contents: []map[string]interface{}{
			{
				"instruction": "send setup",
			},
		},
	}
	timerGeneralMessage, _ := json.Marshal(Message{
		DeviceID: "back-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "time",
		Contents: map[string]interface{}{
			"state":    "stateIdle",
			"duration": 60000,
			"id":       "general",
		},
	})
	statusInstructionMsg, _ := json.Marshal(Message{
		DeviceID: "back-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "instruction",
		Contents: []map[string]interface{}{
			{
				"instruction": "status update",
			},
		},
	})
	statusMessage, _ := json.Marshal(Message{
		DeviceID: "back-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "status",
		Contents: map[string]interface{}{
			"id":         "telephone",
			"connection": false,
			"status":     map[string]interface{}{},
		},
	})
	returnMessage, _ := json.Marshal(Message{
		DeviceID: "back-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "setup",
		Contents: map[string]interface{}{
			"name": "Escape X",
			"hints": map[string][]string{
				"Telefoon puzzels": {"De knop verzend jouw volgorde", "Heb je al even gewacht?"},
				"Control puzzel":   {"Zet de schuiven nauwkeurig"},
			},
			"events": map[string]string{
				"correctSequence": "De juiste volgorde van cijfers moet gedraaid worden.",
			},
			"cameras": []map[string]string{
				{"link": "https://raccoon.games", "name": "camera1"},
				{"link": "https://debrouwerij.io", "name": "camera2"},
			},
		},
	})
	statusMessageFrontEnd, _ := json.Marshal(Message{
		DeviceID: "back-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "status",
		Contents: map[string]interface{}{
			"id":         "front-end",
			"connection": false,
			"status": map[string]interface{}{
				"start": 0,
				"stop":  0},
		},
	})
	messageEventStatus, _ := json.Marshal(Message{
		DeviceID: "back-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "event status",
		Contents: []map[string]interface{}{{
			"id":     "correctSequence",
			"status": false},
		},
	})
	communicatorMock.On("Publish", "front-end", string(returnMessage), 3)
	communicatorMock.On("Publish", "front-end", string(messageEventStatus), 3)
	communicatorMock.On("Publish", "telephone", string(statusInstructionMsg), 3)
	communicatorMock.On("Publish", "front-end", string(statusInstructionMsg), 3)
	communicatorMock.On("Publish", "front-end", string(statusMessageFrontEnd), 3)
	communicatorMock.On("Publish", "front-end", string(timerGeneralMessage), 3)
	communicatorMock.On("Publish", "front-end", string(statusMessage), 3)
	handler.msgMapper(instructionMsg)
	communicatorMock.AssertNumberOfCalls(t, "Publish", 7)
}

func TestOnInstructionMsgSendStatus(t *testing.T) {
	msg := Message{
		DeviceID: "front-end",
		TimeSent: "05-12-2019 09:42:10",
		Type:     "instruction",
		Contents: []map[string]interface{}{
			{"instruction": "send status"},
		},
	}
	communicatorMock := new(CommunicatorMock)
	handler := Handler{
		Config:       config.ReadFile("../../../resources/testing/test_instruction_status.json"),
		Communicator: communicatorMock,
	}
	statusMessageDisplay, _ := json.Marshal(Message{
		DeviceID: "back-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "status",
		Contents: map[string]interface{}{
			"id":         "display",
			"connection": false,
			"status":     map[string]interface{}{},
		},
	})
	statusMessageFrontEnd, _ := json.Marshal(Message{
		DeviceID: "back-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "status",
		Contents: map[string]interface{}{
			"id":         "front-end",
			"connection": false,
			"status": map[string]interface{}{
				"start": 0,
				"stop":  0},
		},
	})
	timerMessage, _ := json.Marshal(Message{
		DeviceID: "back-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "time",
		Contents: map[string]interface{}{
			"state":    "stateIdle",
			"duration": 10000,
			"id":       "timer1",
		},
	})
	timerGeneralMessage, _ := json.Marshal(Message{
		DeviceID: "back-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "time",
		Contents: map[string]interface{}{
			"state":    "stateIdle",
			"duration": 60000,
			"id":       "general",
		},
	})
	eventMessage, _ := json.Marshal(Message{
		DeviceID: "back-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "event status",
		Contents: []map[string]interface{}{{
			"id":     "rule",
			"status": true},
		},
	})
	communicatorMock.On("Publish", "front-end", string(eventMessage), 3)
	communicatorMock.On("Publish", "front-end", string(timerMessage), 3)
	communicatorMock.On("Publish", "front-end", string(timerGeneralMessage), 3)
	communicatorMock.On("Publish", "front-end", string(statusMessageDisplay), 3)
	communicatorMock.On("Publish", "front-end", string(statusMessageFrontEnd), 3)
	handler.onInstructionMsg(msg)
	communicatorMock.AssertNumberOfCalls(t, "Publish", 5)
}

func TestOnInstructionMsgResetAll(t *testing.T) {
	communicatorMock := new(CommunicatorMock)
	workingConfig := config.ReadFile("../../../resources/testing/test_config.json")
	handler := Handler{
		Config:       workingConfig,
		ConfigFile:   "../../../resources/testing/test_config.json",
		Communicator: communicatorMock,
	}
	instructionMsg := Message{
		DeviceID: "front-end",
		TimeSent: "05-12-2019 09:42:10",
		Type:     "instruction",
		Contents: []map[string]interface{}{
			{"instruction": "reset all"},
		},
	}
	responseMsg := Message{
		DeviceID: "back-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "instruction",
		Contents: []map[string]interface{}{{
			"instruction":   "reset",
			"instructed_by": "front-end"},
		},
	}
	statusMsg := Message{
		DeviceID: "back-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "time",
		Contents: map[string]interface{}{
			"id":       "general",
			"duration": 1800000,
			"state":    "stateIdle",
		},
	}

	jsonInstructionMsg, _ := json.Marshal(&responseMsg)
	jsonStatusMsg, _ := json.Marshal(&statusMsg)

	communicatorMock.On("Publish", "client-computers", string(jsonInstructionMsg), 3)
	communicatorMock.On("Publish", "front-end", string(jsonInstructionMsg), 3)
	communicatorMock.On("Publish", "front-end", string(jsonStatusMsg), 3)
	handler.msgMapper(instructionMsg)
	communicatorMock.AssertNumberOfCalls(t, "Publish", 3)
}

func TestOnInstructionMsgTestAll(t *testing.T) {
	communicatorMock := new(CommunicatorMock)
	workingConfig := config.ReadFile("../../../resources/testing/test_config.json")
	handler := Handler{
		Config:       workingConfig,
		Communicator: communicatorMock,
	}
	instructionMsg := Message{
		DeviceID: "front-end",
		TimeSent: "05-12-2019 09:42:10",
		Type:     "instruction",
		Contents: []map[string]interface{}{
			{"instruction": "test all"},
		},
	}
	responseMsg, _ := json.Marshal(Message{
		DeviceID: "back-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "instruction",
		Contents: []map[string]interface{}{{
			"instruction":   "test",
			"instructed_by": "front-end"},
		},
	})
	communicatorMock.On("Publish", "client-computers", string(string(responseMsg)), 3)
	handler.onInstructionMsg(instructionMsg)
	communicatorMock.AssertNumberOfCalls(t, "Publish", 1)
}

func TestOnInstructionMsgTestDevice(t *testing.T) {
	msg := Message{
		DeviceID: "front-end",
		TimeSent: "05-12-2019 09:42:10",
		Type:     "instruction",
		Contents: []map[string]interface{}{{
			"instruction": "test device",
			"device":      "display"},
		},
	}
	communicatorMock := new(CommunicatorMock)
	handler := Handler{
		Config:       config.ReadFile("../../../resources/testing/test_instruction.json"),
		Communicator: communicatorMock,
	}
	returnMessage, _ := json.Marshal(Message{
		DeviceID: "back-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "instruction",
		Contents: []map[string]interface{}{{
			"instruction":   "test",
			"instructed_by": "front-end"},
		},
	})
	communicatorMock.On("Publish", "display", string(returnMessage), 3)
	handler.onInstructionMsg(msg)
	communicatorMock.AssertNumberOfCalls(t, "Publish", 1)
}

func TestOnInstructionMsgFinishRule(t *testing.T) {
	msg := Message{
		DeviceID: "front-end",
		TimeSent: "05-12-2019 09:42:10",
		Type:     "instruction",
		Contents: []map[string]interface{}{{
			"instruction": "finish rule",
			"rule":        "rule"},
		},
	}
	communicatorMock := new(CommunicatorMock)
	handler := Handler{
		Config:       config.ReadFile("../../../resources/testing/test_instruction.json"),
		Communicator: communicatorMock,
	}
	returnMessage, _ := json.Marshal(Message{
		DeviceID: "back-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "event status",
		Contents: []map[string]interface{}{{
			"id":     "rule",
			"status": true},
		},
	})
	communicatorMock.On("Publish", "front-end", string(returnMessage), 3)
	handler.onInstructionMsg(msg)
	communicatorMock.AssertNumberOfCalls(t, "Publish", 1)
}

func TestOnInstructionMsgHint(t *testing.T) {
	msg := Message{
		DeviceID: "front-end",
		TimeSent: "05-12-2019 09:42:10",
		Type:     "instruction",
		Contents: []map[string]interface{}{{
			"instruction": "hint",
			"value":       "This is my hint"},
		},
	}
	communicatorMock := new(CommunicatorMock)
	handler := Handler{
		Config:       config.ReadFile("../../../resources/testing/test_instruction.json"),
		Communicator: communicatorMock,
	}
	returnMessage, _ := json.Marshal(Message{
		DeviceID: "back-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "instruction",
		Contents: []map[string]interface{}{{
			"instruction":   "hint",
			"value":         "This is my hint",
			"instructed_by": "front-end"},
		},
	})
	communicatorMock.On("Publish", "hint", string(returnMessage), 3)
	handler.onInstructionMsg(msg)
	communicatorMock.AssertNumberOfCalls(t, "Publish", 1)
}

func TestOnInstructionMsgMapper(t *testing.T) {
	handler := getTestHandler()
	msg := Message{
		DeviceID: "TestDevice",
		TimeSent: "05-12-2019 09:42:10",
		Type:     "instruction",
		Contents: []map[string]interface{}{
			{"instruction": "name"},
		},
	}

	before := handler.Config
	handler.msgMapper(msg)
	assert.Equal(t, before, handler.Config,
		"Nothing should have been changed after an instruction message type")
}

// TODO fix test to pass full config in json message
func TestInstructionCheckConfigNoErrors(t *testing.T) {
	communicatorMock := new(CommunicatorMock)
	fileName := "../../../resources/testing/test_config.json"
	workingConfig := config.ReadFile(fileName)
	handler := Handler{
		Config:       workingConfig,
		ConfigFile:   fileName,
		Communicator: communicatorMock,
	}
	jsonFile, _ := ioutil.ReadFile(fileName)
	configToTest := make(map[string]interface{})
	if err := json.Unmarshal(jsonFile, &configToTest); err != nil {
		assert.FailNow(t, "cannot create instruction message")
	}
	instructionMsg := Message{
		DeviceID: "front-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "instruction",
		Contents: []map[string]interface{}{
			{
				"instruction": "check config",
				"config":      configToTest,
			},
		},
	}
	returnMsg := Message{
		DeviceID: "back-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "config",
		Contents: map[string][]string{
			"errors": {},
		},
	}
	jsonMessage, _ := json.Marshal(&returnMsg)
	communicatorMock.On("Publish", "front-end", string(jsonMessage), 3)
	handler.msgMapper(instructionMsg)
	communicatorMock.AssertNumberOfCalls(t, "Publish", 1)
}

func TestInstructionCheckConfigWithErrors(t *testing.T) {
	communicatorMock := new(CommunicatorMock)
	fileName := "../../../resources/testing/test_config_errors.json"
	workingConfig := config.ReadFile(fileName)
	handler := Handler{
		Config:       workingConfig,
		ConfigFile:   fileName,
		Communicator: communicatorMock,
	}
	jsonFile, _ := ioutil.ReadFile(fileName)
	configToTest := make(map[string]interface{})
	if err := json.Unmarshal(jsonFile, &configToTest); err != nil {
		assert.FailNow(t, "cannot create instruction message")
	}
	instructionMsg := Message{
		DeviceID: "front-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "instruction",
		Contents: []map[string]interface{}{
			{"instruction": "check config", "config": configToTest},
		},
	}
	returnMsg := Message{
		DeviceID: "back-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "config",
		Contents: map[string][]string{
			"errors": {
				"json: cannot unmarshal number into Go struct field General.general.duration of type string",
			},
		},
	}
	jsonMessage, _ := json.Marshal(&returnMsg)
	communicatorMock.On("Publish", "front-end", string(jsonMessage), 3)
	handler.msgMapper(instructionMsg)
	communicatorMock.AssertNumberOfCalls(t, "Publish", 1)
}

func TestInstructionUseConfig(t *testing.T) {
	communicatorMock := new(CommunicatorMock)
	fileName := "../../../resources/testing/test_config.json"
	workingConfig := config.ReadFile(fileName)
	handler := Handler{
		Config:       workingConfig,
		ConfigFile:   fileName,
		Communicator: communicatorMock,
	}
	jsonFile, _ := ioutil.ReadFile(fileName)
	configToTest := make(map[string]interface{})
	if err := json.Unmarshal(jsonFile, &configToTest); err != nil {
		assert.FailNow(t, "cannot create instruction message")
	}
	instructionMsg := Message{
		DeviceID: "front-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "instruction",
		Contents: []map[string]interface{}{
			{"instruction": "use config", "config": configToTest},
		},
	}
	returnMsg := Message{
		DeviceID: "back-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "config",
		Contents: map[string]interface{}{},
	}
	jsonMessage, _ := json.Marshal(&returnMsg)
	communicatorMock.On("Publish", "front-end", string(jsonMessage), 3)
	handler.msgMapper(instructionMsg)
	communicatorMock.AssertNumberOfCalls(t, "Publish", 1)
}
