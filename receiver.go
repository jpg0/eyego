package eyego

import (
	"io"
	"fmt"
	"net/http"
	"mime/multipart"
	"bytes"
	"os"
	"io/ioutil"
	"path"
	"os/exec"
	"strings"
)

func doPhotoUpload(r *http.Request) (err error) {
	multipartReader, err := r.MultipartReader()

	if err != nil {
		return
	}

	var soapString, integrityDigest, mediaFile, logFile string
	var mediaChecksummer func (string) string

	for {
		Debug("Parsing MIME part...")
		part, err2 := multipartReader.NextPart()
		if err2 == io.EOF {
			Debug("Completed MIME parsing.")
			break
		} else if err2 != nil {
			LogError("Failed MIME parsing: %s", err2)
			return
		}

		switch part.FormName() {
		case "SOAPENVELOPE":
			soapString, err = readString(part)
			Debug("SOAPENVELOPE: %s", soapString)
		case "INTEGRITYDIGEST":
			integrityDigest, err = readString(part)
			Debug("INTEGRITYDIGEST: %s", integrityDigest)
		case "FILENAME":
			mediaFile, mediaChecksummer, logFile, err = writeFiles(part)
			Debug("FILENAME: %s", mediaFile)
		}

		err = part.Close()

		if err != nil {
			LogError("Upload failed: %s", err)
			return
		}
	}

	soap := new(UploadPhoto)
	ParseSoap(soapString, soap)

	Info("Uploading %s", soap.Filename)

	card := GetCard(soap.MacAddress)

	if integrityDigest != mediaChecksummer(card.MacAddress) {
		panic("Bad integrity digest")
	} else {
		Debug("Integrity digest verified ok")
	}


	processUpload(mediaFile, logFile, *soap)
	return nil
}

func processUpload(mediaFile string, logFile string, soap UploadPhoto) {
	geotag(mediaFile, logFile, soap)
	move(mediaFile)
}

func geotag(mediaFile string, logFile string, soap UploadPhoto) string {
	f, err := os.OpenFile(logFile, os.O_RDONLY, 0)
	if err != nil {panic(err)}
	defer f.Close()

	p, err := ParseLog(f)
	if err != nil {panic(err)}

	aps := p.AccessPoints(soap.Filename)

	if len(aps) > 0 && strings.HasSuffix(strings.ToLower(soap.Filename), ".jpg") {
		location, err := GPSCoordinates(aps)
		if err != nil {panic(err)}

		mediaFile = WriteGeotag(mediaFile, location)

		Info("Photo %s geotagged %v:%v", soap.Filename, location.Location.Latitude, location.Location.Longitude)
	}

	move(mediaFile)

	Info("Photo %s archived", soap.Filename)

	return CreateSoap(UploadPhotoResponse{Success:"true"})
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

	cmd := exec.Command(exiftool_location, args...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		panic(fmt.Sprintf("Failed to run exiftool: %s", output))
	}

	return mediaFile
}

func move(mediaFile string) (targetFile string, err error) {
	targetFile = path.Join(Config().TargetDir, path.Base(mediaFile))
	err = os.Rename(mediaFile, targetFile)
	return
}

func readString(p *multipart.Part) (s string, err error) {
	buffer := bytes.NewBuffer(make([]byte, 0))
	_, err = io.Copy(buffer, p)
	if err != nil {
		return
	}

	return string(buffer.Bytes()), nil
}

func writeFiles(r io.Reader) (mediaFile string, mediaChecksum func (string) string, logFile string, err error) {

	var header *Header
	var out *os.File

	Debug("Unpacking and verifying tar")

	checksumReader := NewChecksumReader(r)
	tarReader := NewTarReader(checksumReader)

	targetDir := "/tmp/eyego"
	if _, err := os.Stat(targetDir); err != nil {
		err = os.Mkdir(targetDir, 0755)

		if err != nil {
			panic(err)
		}
	}

	for {
		header, err = tarReader.Next()

		if err == io.EOF {
			break
		} else if err != nil {
			return
		}

		out, err = os.OpenFile(fmt.Sprintf("%s/%s", targetDir, header.Name), os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0600)

		if err != nil {
			panic(err)
		}

		io.Copy(out, tarReader)
		err = out.Close()

		if mediaFile == "" {
			mediaFile = header.Name
		} else {
			logFile = header.Name
		}

		if err != nil {
			return
		}
	}

	ioutil.ReadAll(checksumReader)

	return mediaFile, func(s string) string {return checksumReader.Checksum(s)}, logFile, nil
}

func wrap()
