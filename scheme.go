package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

var lastSchemeId = 0
var schemeDirectory = "data/schemes"

type Scheme struct {
	Id     int
	Name   string
	Lights []*SetLightArguments
}

func NewScheme() *Scheme {
	return &Scheme{}
}

func nextSchemeId() int {
	if lastSchemeId == 0 {
		initLastSchemeId()
	}
	lastSchemeId++
	return lastSchemeId
}

func initLastSchemeId() {
	files, _ := ioutil.ReadDir(schemeDirectory)
	log.Printf("init scheme id: found %v files\n", len(files))
	for _, f := range files {
		sepIndex := strings.Index(f.Name(), ".")
		idStr := f.Name()[:sepIndex]
		id64, _ := strconv.ParseInt(idStr, 10, 32)
		id := int(id64)
		if id > lastSchemeId {
			lastSchemeId = id
		}
	}
}

func (s *Scheme) Persist() error {
	if s.Lights == nil {
		return fmt.Errorf("Cannot persist scheme - nil lights\n")
	}
	if s.Name == "" {
		return fmt.Errorf("Cannot persist scheme - missing name\n")
	}
	if s.Id == 0 {
		s.Id = nextSchemeId()
	}
	data, _ := json.Marshal(s)
	filename := getSchemeFileName(strconv.Itoa(s.Id))
	return ioutil.WriteFile(schemeDirectory+"/"+filename, data, 0664)
}

func LoadSchemes() ([]*Scheme, error) {
	files, err := ioutil.ReadDir(schemeDirectory)
	if err != nil {
		return nil, err
	}
	schemes := make([]*Scheme, 0, len(files))
	for _, file := range files {
		scheme, err := loadSchemeByName(file.Name())
		if err != nil {
			return nil, err
		}
		schemes = append(schemes, scheme)
	}
	return schemes, nil
}

func LoadSchemeById(id string) (*Scheme, error) {
	return loadSchemeByName(getSchemeFileName(id))
}

func DeleteSchemeById(id string) {
	filename := getSchemeFileName(id)
	deleteSchemeByName(filename)
}

func loadSchemeByName(name string) (*Scheme, error) {
	fileContents, err := ioutil.ReadFile(schemeDirectory + "/" + name)
	if err != nil {
		return nil, err
	}
	var scheme = Scheme{}
	err = json.Unmarshal(fileContents, &scheme)
	if err != nil {
		return nil, err
	}
	return &scheme, nil
}

func getSchemeFileName(id string) string {
	maxDigits := 4
	idStr := strings.Repeat("0", maxDigits-len(id)) + id
	return idStr + ".scheme"
}

func deleteSchemeByName(name string) error {
	return os.Remove(schemeDirectory + "/" + name)
}
