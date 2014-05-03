package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

var lastSceneId = 0
var sceneDirectory = "data/scenes"

type Scene struct {
	isRunning        bool
	Id               int
	ActiveHold       time.Duration
	InactiveHoldHold time.Duration
	TransitionTime   time.Duration
	ActiveScheme     int
	InactiveSchemes  []int
}

func (s *Scene) Persist() error {
	if s.Id == 0 {
		s.Id = nextSceneId()
	}

	data, err := json.Marshal(s)
	if err != nil {
		return err
	}
	filename := getSceneFilename(strconv.Itoa(s.Id))
	return ioutil.WriteFile(filename, data, 0664)
}

func nextSceneId() int {
	if lastSceneId == 0 {
		initLastSceneId()
	}
	lastSceneId++
	return lastSceneId
}

func initLastSceneId() {
	files, _ := ioutil.ReadDir(sceneDirectory)
	for _, f := range files {
		sepIndex := strings.Index(f.Name(), ".")
		idStr := f.Name()[:sepIndex]
		id64, _ := strconv.ParseInt(idStr, 10, 32)
		id := int(id64)
		if id > lastSceneId {
			lastSceneId = id
		}
	}
}

func getSceneFilename(id string) string {
	maxDigits := 4
	idStr := strings.Repeat("0", maxDigits-len(id)) + id
	return idStr + ".scene"
}

func LoadSceneById(id string) (*Scene, error) {
	filename := getSceneFilename(id)
	return loadSceneByName(filename)
}

func DeleteSceneById(id string) error {
	filename := getSceneFilename(id)
	return deleteSceneByName(filename)
}

func loadSceneByName(name string) (*Scene, error) {
	fileContents, err := ioutil.ReadFile(sceneDirectory + "/" + name)
	if err != nil {
		return nil, err
	}

	var scene = Scene{}
	err = json.Unmarshal(fileContents, &scene)
	if err != nil {
		return nil, err
	}

	return &scene, nil
}

func deleteSceneByName(name string) error {
	return os.Remove(sceneDirectory + "/" + name)
}

func (s *Scene) NotifyOfMovement(sensorName string) {
}

func (s *Scene) Start() {
	s.isRunning = true
}

func (s *Scene) Stop() {
	s.isRunning = false
}
