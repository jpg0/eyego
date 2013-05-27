package eyego

import (
	"strings"
	"fmt"
	"os"
	"os/exec"
)

func geotag(mediaFile string, logFile string, originalFilename string) string {
	f, err := os.OpenFile(logFile, os.O_RDONLY, 0)
	if err != nil {panic(err)}
	defer f.Close()

	p, err := ParseLog(f)
	if err != nil {panic(err)}

	aps := p.AccessPoints(originalFilename)

	if strings.HasSuffix(strings.ToLower(mediaFile), ".jpg") {
		if len(aps) > 1 {
			location, err := GPSCoordinates(aps)
			if err != nil {panic(err)}

			mediaFile = WriteGeotag(mediaFile, location)

			Info("Photo %s geotagged %v:%v", originalFilename, location.Location.Latitude, location.Location.Longitude)
		} else {
			Info("Insufficient Access Points logged for %s, skipping geotag.", originalFilename)
		}
	}

	return mediaFile
}

func WriteGeotag(mediaFile string, location LocationResult) string {
	args := []string{
		"-GPSMapDatum=WGS-84",
		"-GPSMeasureMode=2",
		"-GPSVersionID=2 0 0 0"}


	lat := location.Location.Latitude
	if lat != 0 {
		var latRef string
		if lat > 0 {
			latRef = "N"
		} else {
			latRef = "S"
			lat = -lat
		}
		args = append(args,
			fmt.Sprintf("-GPSLatitudeRef=%s", latRef),
			fmt.Sprintf("-GPSLatitude=%v", lat))

	}

	lon := location.Location.Longitude
	if lon != 0 {
		var lonRef string
		if lon > 0 {
			lonRef = "E"
		} else {
			lonRef = "W"
			lon = -lon
		}
		args = append(args,
			fmt.Sprintf("-GPSLongitudeRef=%s", lonRef),
			fmt.Sprintf("-GPSLongitude=%v", lon))

	}

	args = append(args, mediaFile)

	exiftool_location, err := exec.LookPath("exiftool")

	if err != nil {
		panic("exiftool not found")
	}

	Trace("Running %s, %v", exiftool_location, args)
	cmd := exec.Command(exiftool_location, args...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		panic(fmt.Sprintf("Failed to run exiftool: %s", output))
	}

	return mediaFile
}
