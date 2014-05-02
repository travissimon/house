package main

import (
	"github.com/travissimon/huego"
)

type LightProxy struct {
	Id   string
	Name string
	Hex  string
}

func NewLightProxyFromLight(light *huego.Light) *LightProxy {
	hex := light.ToHex()
	return &LightProxy{light.Id, light.Name, hex}
}

func (lp *LightProxy) String() string {
	return lp.Hex
}
