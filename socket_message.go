package main

import (
	"fmt"
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

type SetLightArguments struct {
	Id   string
	Name string
	Hex  string
}

func (s *SetLightArguments) String() string {
	return s.Hex
}

type SaveSceneArguments struct {
	Id     string
	Name   string
	Lights []*SetLightArguments
}

type SetSchemeArguments struct {
	Id string
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

	return SaveSceneArguments{id, name, lights}
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

func (m *SocketMessage) String() string {
	return fmt.Sprintf("Request: %v, Response: %v", m.Request, m.Response)
}
