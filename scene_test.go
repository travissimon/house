package main

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
	"time"
)

func Test_ColourHookup(t *testing.T) {
	log.Println("Hookup Succeeded, testing scenes")
}

func Test_ScenDataSerialisation(t *testing.T) {
	sd := new(SceneData)
	sd.SceneHold = 5 * time.Second
	sd.TransitionTime = 2 * time.Second
	sd.ActiveScheme = 3
	sd.InactiveSchemes = make([]int, 0, 3)
	sd.InactiveSchemes = append(sd.InactiveSchemes, 1)
	sd.InactiveSchemes = append(sd.InactiveSchemes, 2)
	sd.InactiveSchemes = append(sd.InactiveSchemes, 3)

	data, _ := json.Marshal(sd)
	fmt.Printf("%s\n", data)
}
