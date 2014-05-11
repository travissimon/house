package main

import (
	"encoding/json"
	"github.com/travissimon/huego"
	"io/ioutil"
	"log"
)

// Utilities for dealing with Huego bases

var baseDirectory = "data/bases"

type BaseData struct {
	Id         string `json:"id"`
	InternalIp string `json:"internalipaddress"`
	Username   string
}

func GetHueBase() (*huego.Base, error) {
	bases, err := loadBasesFromFilesystem()
	if err != nil {
		return nil, err
	}
	if bases != nil && len(bases) > 0 {
		b := bases[0]
		log.Printf("Loaded from filessystem: %v", b)
		b.GetLights()
		log.Printf("After get lights: %v", b)
		return bases[0], nil
	}
	baseInstances, err := huego.DiscoverBases()
	if err != nil {
		return nil, err
	}
	if baseInstances != nil && len(baseInstances) > 0 {
		base := &baseInstances[0]
		persistBase(base)
		return base, nil
	}
	return nil, nil
}

func getBaseFileName(id string) string {
	return id + ".base"
}

func loadBasesFromFilesystem() ([]*huego.Base, error) {
	files, err := ioutil.ReadDir(baseDirectory)
	if err != nil {
		return nil, err
	}
	bases := make([]*huego.Base, 0, len(files))
	for _, file := range files {
		base, err := loadBaseByName(file.Name())
		if err != nil {
			return nil, err
		}
		bases = append(bases, base)
	}
	return bases, nil
}

func loadBaseByName(name string) (*huego.Base, error) {
	fileContents, err := ioutil.ReadFile(baseDirectory + "/" + name)
	if err != nil {
		return nil, err
	}
	var cfgData BaseData
	err = json.Unmarshal(fileContents, &cfgData)
	if err != nil {
		return nil, err
	}
	b := huego.Base{}
	b.Id = cfgData.Id
	b.InternalIp = cfgData.InternalIp
	b.Username = cfgData.Username
	// initialise base
	_, _ = b.GetLights()
	return &b, nil
}

func persistBase(b *huego.Base) error {
	cfgData := BaseData{b.Id, b.InternalIp, b.Username}

	data, err := json.Marshal(cfgData)
	if err != nil {
		return err
	}
	filename := getBaseFileName(b.Id)
	return ioutil.WriteFile(baseDirectory+"/"+filename, data, 0664)
}
