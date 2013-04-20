package eyefi

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

type AccessPoint struct {
	MacAddress string
	Strength string
	Data int32
}

type PoweredCycle struct {
	Photos map[string] []NewPhoto
	AccessPoints map[string] []AccessPoint
}

type PowerOn struct {

}

func ParseLog(log io.Reader) (p ParsedLog, err error) {
	lines, err := EachLine(log, readLine)

	p = ParsedLog{
		Cycles: make([]PoweredCycle,0)}

	cycle := PoweredCycle{
		Photos: make(map[string] []NewPhoto),
		AccessPoints: make(map[string] []AccessPoint)}

	for i := range lines {
		switch t := lines[i].(type){
		case AccessPoint:
			if cycle.AccessPoints[t.MacAddress] == nil {
				cycle.AccessPoints[t.MacAddress] = make([]AccessPoint, 0, 1)
			}

			cycle.AccessPoints[t.MacAddress] = append(cycle.AccessPoints[t.MacAddress], t)
		case NewPhoto:
			if cycle.Photos[t.Filename] == nil {
				cycle.Photos[t.Filename] = make([]NewPhoto, 0, 1)
			}

			cycle.Photos[t.Filename] = append(cycle.Photos[t.Filename], t)
		case PowerOn:
			p.Cycles = append(p.Cycles, cycle)
			cycle = PoweredCycle{
				Photos: make(map[string] []NewPhoto),
				AccessPoints: make(map[string] []AccessPoint)}
		}
	}

	p.Cycles = append(p.Cycles, cycle)
	return
}

func readLine(line string) interface {} {

	elements := strings.Split(strings.Trim(line, " "), ",")

	power_secs := atoi(elements[0])
	secs := atoi(elements[1])
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
		return AccessPoint{
			MacAddress: mac,
			Strength: strength,
			Data: int32(data)}
	case "NEWPHOTO":
		filename := args[0]
		size := atoi(args[1])
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

func atoi(s string) (i int) {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return
}

func (l *ParsedLog) newAccessPoint() {

}
