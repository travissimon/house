package main

import (
	"fmt"
	"time"
)

type SetGeneratorArguments struct {
	PrimaryColour string
	Strategy      string
	Angle         float64
	Tint          float64
	Shade         float64
}

type SocketMessage struct {
	SenderId int
	Request  ClientRequest
	Response []*LightProxy
}

type ClientRequest struct {
	Action    string
	Arguments map[string]interface{}
}

type SetPowerArguments struct {
	Id     string
	TurnOn bool
}

type SetLightArguments struct {
	Id   string
	Name string
	Hex  string
}

func (s *SetLightArguments) String() string {
	return s.Hex
}

type SaveSceneArguments struct {
	Id                 int
	Name               string
	ActiveTransition   time.Duration
	ActiveHold         time.Duration
	InactiveTransition time.Duration
	InactiveHold       time.Duration
	ActiveScheme       int
	InactiveSchemes    []int
}

type SetSceneArguments struct {
	Id string
}

type DeleteSceneArguments struct {
	Id string
}

type SetSchemeArguments struct {
	Id string
}

type SaveSchemeArguments struct {
	Id     string
	Name   string
	Lights []*SetLightArguments
}

type DeleteSchemeArguments struct {
	Id string
}

func (m *SocketMessage) GetSetLightArguments() SetLightArguments {
	args := m.Request.Arguments
	id := args["Id"].(string)
	name := args["Name"].(string)
	val := args["Hex"].(string)
	return SetLightArguments{id, name, val}
}

func (m *SocketMessage) GetSetGeneratorArguments() SetGeneratorArguments {
	args := m.Request.Arguments
	primaryColour := args["PrimaryColour"].(string)
	strategy := args["Strategy"].(string)
	angle := args["Angle"].(float64)
	tint := args["Tint"].(float64)
	shade := args["Shade"].(float64)
	return SetGeneratorArguments{primaryColour, strategy, angle, tint, shade}
}

func (m *SocketMessage) GetSaveSceneArguments() SaveSceneArguments {
	args := m.Request.Arguments
	id := 0
	if args["id"] != nil {
		id = int(args["id"].(float64))
	}
	name := args["name"].(string)
	activeTransition, _ := time.ParseDuration(args["activeTransition"].(string))
	activeHold, _ := time.ParseDuration(args["activeHold"].(string))
	inactiveTransition, _ := time.ParseDuration(args["inactiveTransition"].(string))
	inactiveHold, _ := time.ParseDuration(args["inactiveHold"].(string))
	activeScheme := int(args["activeScheme"].(float64))
	inactiveSchemes := args["inactiveSchemes"].([]interface{})

	schs := make([]int, 0, len(inactiveSchemes))
	for _, inact := range inactiveSchemes {
		schs = append(schs, int(inact.(float64)))
	}

	return SaveSceneArguments{id, name, activeTransition, activeHold, inactiveTransition, inactiveHold, activeScheme, schs}
}

func (m *SocketMessage) GetSetSceneArguments() SetSceneArguments {
	args := m.Request.Arguments
	id := args["Id"].(string)
	return SetSceneArguments{id}
}

func (m *SocketMessage) GetDeleteSceneArguments() DeleteSceneArguments {
	args := m.Request.Arguments
	id := args["Id"].(string)
	return DeleteSceneArguments{id}
}

func (m *SocketMessage) GetSetSchemeArguments() SetSchemeArguments {
	args := m.Request.Arguments
	id := args["Id"].(string)
	return SetSchemeArguments{id}
}

func (m *SocketMessage) GetDeleteSchemeArguments() DeleteSchemeArguments {
	args := m.Request.Arguments
	id := args["Id"].(string)
	return DeleteSchemeArguments{id}
}

func (m *SocketMessage) GetSaveSchemeArguments() SaveSchemeArguments {
	args := m.Request.Arguments
	id := args["Id"].(string)
	name := args["Name"].(string)
	lightArray := args["Lights"].([]interface{})
	lights := make([]*SetLightArguments, 0, len(lightArray))
	for _, lightMap := range lightArray {
		lm := lightMap.(map[string]interface{})
		id := lm["Id"].(string)
		name := lm["Name"].(string)
		hex := lm["Hex"].(string)
		l := SetLightArguments{string(id), name, hex}
		lights = append(lights, &l)
	}
	return SaveSchemeArguments{id, name, lights}
}

func (m *SocketMessage) GetSetPowerArguments() SetPowerArguments {
	args := m.Request.Arguments
	id := args["Id"].(string)
	turnOn := args["TurnOn"].(bool)
	return SetPowerArguments{id, turnOn}
}

func (m *SocketMessage) String() string {
	return fmt.Sprintf("Request: %v, Response: %v", m.Request, m.Response)
}
