package eyefi

import (
	"io"
	"fmt"
	"net/http"
	"mime/multipart"
	"bytes"
	"tar2"
	"os"
	"io/ioutil"
	"path"
)

func doPhotoUpload(r *http.Request) (err error) {
	multipartReader, err := r.MultipartReader()

	if err != nil {
		return
	}

	var soapString, integrityDigest, mediaFile, logFile string
	var mediaChecksummer func (string) string

	for {
		part, err2 := multipartReader.NextPart()
		if err2 == io.EOF {
			break
		} else if err2 != nil {
			return
		}

		switch part.FormName() {
		case "SOAPENVELOPE":
			soapString, err = readString(part)
		case "INTEGRITYDIGEST":
			integrityDigest, err = readString(part)
		case "FILENAME":
			mediaFile, mediaChecksummer, logFile, err = writeFiles(part)
		}

		if err != nil { return }
	}

	soap := new(UploadPhoto)
	ParseSoap(soapString, soap)
	card := GetCard(soap.MacAddress)

	if integrityDigest != mediaChecksummer(card.MacAddress) {
		panic("Bad integrity digest")
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

	if len(aps) > 0 {
		location, err := GPSCoordinates(aps)
		if err != nil {panic(err)}

		mediaFile = writeGeotag(mediaFile, location)
	}

	move(mediaFile)

	return "ok"
}

func writeGeotag(mediaFile string, location LocationResult) string {
	//os.StartProcess()
	return ""
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

	var header *tar2.Header
	var out *os.File

	checksumReader := NewChecksumReader(r)
	tarReader := tar2.NewReader(checksumReader)

	targetDir := "/tmp/eyego"
	err = os.Mkdir(targetDir, 0755)


	if err != nil {
		panic(err)
	}

	for {
		header, err = tarReader.Next()

		if err == io.EOF {
			break
		} else if err != nil {
			return
		}

		out, err = os.OpenFile(fmt.Sprintf("%s/%s", targetDir, header.Name), os.O_WRONLY | os.O_CREATE | os.O_EXCL, 0600)

		if err != nil {
			panic(err)
		}

		fmt.Println(out)
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
