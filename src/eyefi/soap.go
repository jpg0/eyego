package eyefi

import (
	"encoding/xml"
	"bytes"
	"io"
)

const (
	soapStart = `<?xml version="1.0" encoding="UTF-8"?><SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/" xmlns:ns1="EyeFi/SOAP/EyeFilm"><SOAP-ENV:Body>`
	soapEnd = `</SOAP-ENV:Body></SOAP-ENV:Envelope>`
)

type SoapStartSession struct {
	MacAddress string `xml:"macaddress"`
	CNonce     string `xml:"cnonce"`
	TransferMode          string `xml:"transfermode"`
	TransferModeTimestamp string `xml:"transfermodetimestamp"`
}

type SoapStartSessionResponse struct {
	Credential string `xml:"credential"`
	SNonce     string `xml:"snonce"`
	TransferMode          string `xml:"transfermode"`
	TransferModeTimestamp string `xml:"transfermodetimestamp"`
	UpSyncAllowed string `xml:"upsyncallowed"`
	XMLName xml.Name `xml:"ns1:StartSessionResponse"`
}

type GetPhotoStatus struct {
	Credential string `xml:"credential"`
	MacAddress string `xml:"macaddress"`
	Filename string `xml:"filename"`
	FileSize string `xml:"filesize"`
	FileSignature string `xml:"filesignature"`
	Flags string `xml:"flags"`
	XMLName xml.Name `xml:"ns1:GetPhotoStatus"`
}

type GetPhotoStatusResponse struct {
	FileID string `xml:"fileid"`
	Offset string `xml:"offset"`
	XMLName xml.Name `xml:"ns1:GetPhotoStatusResponse"`
}

func ParseSoap(s string, target interface {}) {
	if s[0:len(soapStart)] != soapStart || s[len(s)-len(soapEnd):len(s)] != soapEnd {
		panic("Unknown soap request:\n" + s)
	}

	body := s[len(soapStart):len(s)-len(soapEnd)]

	parser := xml.NewDecoder(bytes.NewBufferString(body))
	parser.Decode(target)
}

func CreateSoap(body interface{}) string {
	buffer := bytes.NewBuffer(make([]byte, 0))
	WriteSoap(body, buffer)
	return string(buffer.Bytes())
}

func WriteSoap(body interface{}, writer io.Writer) {
	encoder := xml.NewEncoder(writer)
	writer.Write([]byte(soapStart))
	encoder.Encode(body)
	writer.Write([]byte(soapEnd))
}
