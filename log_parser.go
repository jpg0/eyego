package eyego

import (
	"io"
	"strconv"
	"fmt"
	"strings"
)

type ParsedLog struct {
	Cycles []PoweredCycle
}

type NewPhoto struct {
	Filename string
	PowerSecs int
	Secs int
	Size int
}

type AccessPointSighting struct {
	MacAddress string
	SNR int
	Data int32
	Secs int
	PowerSecs int
}

type PoweredCycle struct {
	Photos map[string] []NewPhoto
	AccessPoints map[string] []AccessPointSighting
}

type PowerOn struct {}

type AccessPointSightingInfo struct {
	MacAddress string `json:"macAddress"`
	Age int `json:"age"`
	SNR int `json:"signalToNoiseRatio"`
}

func ParseLog(log io.Reader) (p ParsedLog, err error) {
	lines, err := EachLine(log, readLine)

	p = ParsedLog{
		Cycles: make([]PoweredCycle,0)}

	cycle := PoweredCycle{
		Photos: make(map[string] []NewPhoto),
		AccessPoints: make(map[string] []AccessPointSighting)}

	for i := range lines {
		switch t := lines[i].(type){
		case AccessPointSighting:
			Trace("Found AccessPointSighting")
			if cycle.AccessPoints[t.MacAddress] == nil {
				cycle.AccessPoints[t.MacAddress] = make([]AccessPointSighting, 0, 1)
			}

			cycle.AccessPoints[t.MacAddress] = append(cycle.AccessPoints[t.MacAddress], t)
		case NewPhoto:
			Trace("Found NewPhoto")
			if cycle.Photos[t.Filename] == nil {
				cycle.Photos[t.Filename] = make([]NewPhoto, 0, 1)
			}

			cycle.Photos[t.Filename] = append(cycle.Photos[t.Filename], t)
		case PowerOn:
			Trace("Found PowerOn")
			p.Cycles = append(p.Cycles, cycle)
			cycle = PoweredCycle{
				Photos: make(map[string] []NewPhoto),
				AccessPoints: make(map[string] []AccessPointSighting)}
		}
	}

	p.Cycles = append(p.Cycles, cycle)
	return
}

func (p ParsedLog) AccessPoints(photoName string) []AccessPointSightingInfo {
	for i := range p.Cycles {
		cycle := p.Cycles[i]
		if cycle.Photos[photoName] != nil {
			return cycle.accessPoints(photoName)
		}
	}

	return make([]AccessPointSightingInfo,0)
}

func (c PoweredCycle) accessPoints(photoName string) []AccessPointSightingInfo {
	m := make(map[string]AccessPointSightingInfo)
	timeTaken := c.Photos[photoName][0].PowerSecs

	for i := range c.AccessPoints {
		aps := c.AccessPoints[i]

		for j := range aps {
			ap := aps[j]
			if Abs(ap.PowerSecs - timeTaken) < 300 {
				m[ap.MacAddress] =  AccessPointSightingInfo{
					MacAddress: formatMac(ap.MacAddress),
					Age: Abs(ap.PowerSecs - timeTaken) * 1000,
					SNR: ap.SNR}
			}
		}
	}

	rv := make([]AccessPointSightingInfo, 0, len(m))

	for  _, value := range m {
		rv = append(rv, value)
	}

	return rv
}

func readLine(line string) interface {} {

	elements := strings.Split(strings.Trim(line, " "), ",")

	Trace("Parsing line: %s", elements)

	power_secs := Atoi(elements[0])
	secs := Atoi(elements[1])
	args := elements[3:]

	switch (elements[2]){
	case "POWERON":
		return PowerOn{}
	case "AP","NEWAP":
		mac, strength := args[0], args[1]
		var data int64
		var err error
		if len(args) > 2 {
			data, err = strconv.ParseInt(args[2], 16, 32)
			if err != nil { panic(fmt.Sprintf("Cannot parse signal strength: %v", err)) }
		} else {
			data = 0
		}
		return AccessPointSighting{
			MacAddress: mac,
			SNR: Atoi(strength),
			Data: int32(data),
			PowerSecs: power_secs,
			Secs: secs}
	case "NEWPHOTO":
		filename := args[0]
		size := Atoi(args[1])
		return NewPhoto{
			Filename: filename,
			PowerSecs: power_secs,
			Secs: secs,
			Size: size}
	default:
		panic(fmt.Sprintf("Unknown event type: '%s'", elements[2]))
	}

	panic("unreachable")
}

func formatMac(mac string) string {
	parts := []string{mac[0:2],mac[2:4],mac[4:6],mac[6:8],mac[8:10],mac[10:12]}
	return strings.Join(parts, ":")
}
