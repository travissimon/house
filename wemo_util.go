package main

import (
	"github.com/savaki/go.wemo"
	"log"
	"time"
)

type WemoUtil struct {
	wemoApi       *wemo.Wemo
	motionSensors []*wemo.Device
	sensorNames   []string
	poller        *time.Ticker
	stopChan      chan bool
}

func NewWemoUtil() *WemoUtil {
	util := new(WemoUtil)
	util.wemoApi, _ = wemo.NewByInterface("wlan0")
	util.sensorNames = make([]string, 3)

	devices, _ := util.wemoApi.DiscoverAll(1 * time.Second)
	for _, device := range devices {
		info, _ := device.FetchDeviceInfo()
		if info.DeviceType[len(info.DeviceType)-1:] == "1" {
			util.motionSensors = append(util.motionSensors, device)
			util.sensorNames = append(util.sensorNames, info.FriendlyName)
		}
	}

	return util
}

type MotionDetected func(string)

func (util *WemoUtil) Start(onMotion MotionDetected) {
	log.Printf("Starting Wemo util")
	util.stopChan = make(chan bool)
	util.poller = time.NewTicker(100 * time.Millisecond)
	go func() {
		for {
			select {
			case <-util.poller.C:
				for idx, d := range util.motionSensors {
					isActive := d.GetBinaryState()
					if isActive == 1 {
						onMotion(util.sensorNames[idx])
					}
				}
			case <-util.stopChan:
				return
			}
		}
	}()
}

func (util *WemoUtil) Stop() {
	if util.stopChan != nil {
		log.Printf("Requesting Wemo stop")
		util.stopChan <- true
		close(util.stopChan)
		util.stopChan = nil
	}
}
