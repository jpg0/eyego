package eyefi

import (
	"io"
	"fmt"
	"promise"
	"net/http"
	"mime/multipart"
	"bytes"
	"tar2"
	"os"
)

func doPhotoUpload(r *http.Request) (rv promise.Promise, err error) {
	multipartReader, err := r.MultipartReader()

	if err != nil {
		return
	}

	var soapString,integrityDigest, mediaFile, logFile string

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
			mediaFile, logFile, err = writeFiles(part)
		}

		if err != nil { return }
	}

	return processUpload(mediaFile, logFile, soapString, integrityDigest), nil
}

func processUpload(mediaFile string, logFile string, soapString string, integrityDigest string) promise.Promise {
	return nil
}

func readString(p *multipart.Part) (s string, err error) {
	buffer := bytes.NewBuffer(make([]byte, 0))
	_, err = io.Copy(buffer, p)
	if err != nil {
		return
	}

	return string(buffer.Bytes()), nil
}

func writeFiles(r io.Reader) (mediaFile string, logFile string, err error) {

	var header *tar2.Header
	var out *os.File

	tarReader := tar2.NewReader(r)

	for {
		header, err = tarReader.Next()

		if err == io.EOF {
			break
		} else if err != nil {
			return
		}

		out, err = os.OpenFile(fmt.Sprintf("/tmp/%s", header.Name), os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
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

	return mediaFile, logFile, nil
}
