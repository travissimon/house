package main

import (
	"github.com/travissimon/huego"
)

type LightProxy struct {
	Id   string
	Name string
	On   bool
	Hex  string
}

func NewLightProxyFromLight(light *huego.Light) *LightProxy {
	hex := light.ToHex()
	return &LightProxy{light.Id, light.Name, light.State.On, hex}
}

func (lp *LightProxy) String() string {
	return lp.Hex
}
