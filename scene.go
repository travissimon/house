package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var lastSceneId = 0
var sceneDirectory = "data/scenes"

type Scene struct {
	Id                 int
	Name               string
	ActiveHold         time.Duration
	InactiveHold       time.Duration
	InactiveTransition time.Duration
	ActiveTransition   time.Duration
	ActiveScheme       int
	InactiveSchemes    []int
	SensorName         string
	isActive           bool
	timer              *time.Timer
	schemeIdx          int
}

func NewScene() *Scene {
	return &Scene{
		0,
		"",               // Name
		5 * time.Minute,  // Active hold
		5 * time.Minute,  // Inactive hold
		10 * time.Second, // Inactive transition
		2 * time.Second,  // Active transition
		1,                // active scheme
		make([]int, 0),   // inactive scheme
		"",               // sensor name
		false,            // isActive
		nil,              // private timer
		0,                // current scheme index
	}
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
	return ioutil.WriteFile(sceneDirectory+"/"+filename, data, 0664)
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

func LoadScenes() ([]*Scene, error) {
	files, err := ioutil.ReadDir(sceneDirectory)
	if err != nil {
		return nil, err
	}
	scenes := make([]*Scene, 0, len(files))
	for _, file := range files {
		scene, err := loadSceneByName(file.Name())
		if err != nil {
			return nil, err
		}
		scenes = append(scenes, scene)
	}
	return scenes, nil
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
	if !s.isActive {
		log.Printf("Transitioning to Active Scene")
		s.transitionToScheme(s.ActiveScheme, s.ActiveTransition)
	}
	s.isActive = true
	s.timer.Reset(s.ActiveHold + s.ActiveTransition)
}

func (s *Scene) Start() {
	log.Printf("Starting scene '%v'\n", s.Name)
	s.timer = time.AfterFunc(s.InactiveHold, s.durationComplete)
}

func (s *Scene) durationComplete() {
	s.isActive = false
	curId := s.InactiveSchemes[s.schemeIdx]
	s.transitionToScheme(curId, s.InactiveTransition)
	s.incrSchemeIdx()

	s.timer = time.AfterFunc(s.InactiveTransition+s.InactiveHold, s.durationComplete)
}

func (s *Scene) transitionToScheme(id int, duration time.Duration) {
	scheme, _ := LoadSchemeById(strconv.Itoa(id))
	for _, light := range scheme.Lights {
		l := getLightById(light.Id)
		l.SetColourFromHex(light.Hex)
		nanos := duration.Nanoseconds()
		tmp := nanos / (1000 * 1000 * 100)
		tenths := uint16(tmp)
		l.SetStateWithTransition(tenths)
	}
}

func (s *Scene) incrSchemeIdx() {
	s.schemeIdx++
	if s.schemeIdx >= len(s.InactiveSchemes) {
		s.schemeIdx = 0
	}
}

func (s *Scene) Stop() {
	if s == nil {
		return
	}
	log.Printf("Stopping scene '%v'\n", s.Name)
	if s.timer != nil {
		s.timer.Stop()
	}
}
