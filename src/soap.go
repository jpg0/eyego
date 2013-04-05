package soap

import (
	"encoding/xml"
	"bytes"
	"fmt"
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
}

type SoapFault struct {
	Faultstring string
	Detail      string
}

type SoapBody struct {
	Fault           SoapFault
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

func CreateSoap(body SoapBody) SoapEnvelope {
	return SoapEnvelope{
		Body: body,
		XMLName: xml.Name{
			Space:"http://schemas.xmlsoap.org/soap/envelope/",
			Local:"Envelope"}}
}
