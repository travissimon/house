package main

import (
	"fmt"
)

type SetSchemeArguments struct {
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

func (m *SocketMessage) GetSetLightArguments() SetLightArguments {
	args := m.Request.Arguments
	id := args["Id"].(string)
	val := args["Value"].(string)
	return SetLightArguments{id, val}
}

func (m *SocketMessage) GetSetSchemeArguments() SetSchemeArguments {
	args := m.Request.Arguments
	primaryColour := args["PrimaryColour"].(string)
	strategy := args["Strategy"].(string)
	angle := args["Angle"].(float64)
	tint := args["Tint"].(float64)
	shade := args["Shade"].(float64)
	return SetSchemeArguments{primaryColour, strategy, angle, tint, shade}
}

func (m *SocketMessage) String() string {
	return fmt.Sprintf("Request: %v, Response: %v", m.Request, m.Response)
}
