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
	"mime"
)

func teeMultipartReader(r *http.Request) (*multipart.Reader, *bytes.Buffer, error) {
    v := r.Header.Get("Content-Type")

    if v == "" {
        return nil, nil, http.ErrNotMultipart
    }
    d, params, err := mime.ParseMediaType(v)
    if err != nil || d != "multipart/form-data" {
        return nil, nil, http.ErrNotMultipart
    }
    boundary, ok := params["boundary"]
    if !ok {
        return nil, nil, http.ErrMissingBoundary
    }

		var buffer bytes.Buffer
		teeReader := io.TeeReader(r.Body, &buffer)

    return multipart.NewReader(teeReader, boundary),  &buffer, nil
}

func doPhotoUpload(r *http.Request) (s string, err error) {
	multipartReader, buffer, err := teeMultipartReader(r)

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
			Debug("%s", buffer.String())
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

	calculatedDigest := mediaChecksummer(card.UploadKey)

	if integrityDigest != calculatedDigest {
		panic(fmt.Sprintf("Bad integrity digest. Calculated %s, sent %s", calculatedDigest, integrityDigest))
	} else {
		Debug("Integrity digest verified ok")
	}


	return processUpload(mediaFile, logFile, *soap), nil
}

func processUpload(mediaFile string, logFile string, soap UploadPhoto) string {

	if Config().GoogleAPIKey != "" {
		mediaFile = geotag(mediaFile, logFile, path.Base(mediaFile))
	}

	target, err := move(mediaFile)

	if err != nil {
		panic(err)
	}

	Info("Photo %s archived to %s", soap.Filename, target)

	return CreateSoap(UploadPhotoResponse{Success:"true"})
}

func move(mediaFile string) (targetFile string, err error) {
	targetFile = path.Join(Config().TargetDir, path.Base(mediaFile))
	err = os.Rename(mediaFile, targetFile)

	if err != nil {
		Info("Failed to move file, falling back to copy [%s]", err)
		_, err = CopyFile(targetFile, mediaFile)
	}

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
			mediaFile = fmt.Sprintf("%s/%s", targetDir, header.Name)
		} else {
			logFile = fmt.Sprintf("%s/%s", targetDir, header.Name)
		}

		if err != nil {
			return
		}
	}

	ioutil.ReadAll(checksumReader)

	return mediaFile, func(s string) string {return checksumReader.Checksum(s)}, logFile, nil
}
