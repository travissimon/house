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
	Id    string
	Value string
}

type CreateSceneArguments struct {
	Name   string
	Lights []*SetLightArguments
}

type SetSchemeArguments struct {
	Id string
}

func (m *SocketMessage) GetSetLightArguments() SetLightArguments {
	args := m.Request.Arguments
	id := args["Id"].(string)
	val := args["Value"].(string)
	return SetLightArguments{id, val}
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

func (m *SocketMessage) GetCreateSceneArguments() CreateSceneArguments {
	args := m.Request.Arguments
	name := args["Name"].(string)
	lightArray := args["Lights"].([]interface{})
	lights := make([]*SetLightArguments, 0, len(lights))
	for _, lightMap := range lightArray {
		lm := lightMap.(map[string]interface{})
		id := int32(lm["Id"].(float64))
		hex := lm["Hex"].(string)
		l := SetLightArguments{string(id), hex}
		lights = append(lights, &l)
	}

	return CreateSceneArguments{name, lights}
}

func (m *SocketMessage) GetSetSchemeArguments() SetSchemeArguments {
	args := m.Request.Arguments
	id := args["Id"].(string)
	return SetSchemeArguments{id}
}

func (m *SocketMessage) String() string {
	return fmt.Sprintf("Request: %v, Response: %v", m.Request, m.Response)
}
