package eyefi

import (
	"encoding/xml"
	"bytes"
	"fmt"
	"io"
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
	XMLName xml.Name
}

type SoapFault struct {
	Faultstring string
	Detail      string
}

type SoapBody struct {
	Fault           *SoapFault
	//possible bodies
	StartSession    *SoapStartSession
	StartSessionResponse *SoapStartSessionResponse
}
type SoapEnvelope struct {
	XMLName xml.Name
	Body    SoapBody
}

func ParseSoap(s string) SoapBody {
	parser := xml.NewDecoder(bytes.NewBufferString(s))
	envelope := new(SoapEnvelope)
	parser.DecodeElement(&envelope, nil)
	fmt.Print(envelope.XMLName)
	return envelope.Body
}

func CreateSoap(body SoapBody) string {
	buffer := bytes.NewBuffer(make([]byte, 0))
	WriteSoap(body, buffer)
	return string(buffer.Bytes())
}

func WriteSoap(body SoapBody, writer io.Writer) {
	encoder := xml.NewEncoder(writer)
	writer.Write([]byte(xml.Header))
	encoder.Encode(SoapEnvelope{
		Body: body,
		XMLName: xml.Name{
			Space:"http://schemas.xmlsoap.org/soap/envelope/",
			Local:"Envelope"}})
}
