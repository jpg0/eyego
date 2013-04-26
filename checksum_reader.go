package eyego

import (
	"io"
	"encoding/hex"
	"fmt"
	"crypto/md5"
	"encoding/base64"
	"bytes"
//	"io/ioutil"
)

type ChecksumReader struct {
	delegate  io.Reader
	state *checksumState
}

type checksumState struct {
	checksums []uint16
}

func NewChecksumReader(r io.Reader) ChecksumReader {

	cr := ChecksumReader {
		state:&checksumState{
		checksums: make([]uint16,0, 16)}}

	bnr := NewBlockNotifyingReader(r, 512, func(b []byte) {
			cr.appendBlock(b)
		})

	cr.delegate = bnr

	return cr
}

func (r ChecksumReader) Read(p []byte) (n int, err error){
	return r.delegate.Read(p)
}

func (cr ChecksumReader) Checksum(uploadKey string) string {

	b, _ := hex.DecodeString(uploadKey)

	if len(b)%2 != 0 { panic("Bad upload key")}

	h := md5.New()

	Trace("Beginning checksum")

	for i := 0; i < len(cr.state.checksums); i++ {
		Trace("Adding %v to hash", cr.state.checksums[i])
		h.Write([]byte{
			byte(cr.state.checksums[i]),
			byte(cr.state.checksums[i]>>8)})
	}

	Trace("Adding %v to hash", b)
	h.Write(b)

	return fmt.Sprintf("%x", h.Sum(nil))
}

func (cr *ChecksumReader) appendBlock(b []byte) {
	cr.state.checksums = append(cr.state.checksums, tcp_checksum(b))
}

func tcp_checksum(b []byte) uint16 {
	if len(b) %2 != 0 { panic(fmt.Sprintf("tcp checksum bad length: %d", len(b))) }

	buf := bytes.NewBuffer(make([]byte, 0))
	enc := base64.NewEncoder(base64.StdEncoding, buf)
	enc.Write(b)
//	b64, _ := ioutil.ReadAll(buf)

//	Trace("To Checksum: %v", string(b64))

	var sum uint32 = 0
	var tmp uint16

	for c := 0; c < len(b); c = c + 2 {
		tmp = uint16(b[c]) | uint16(b[c + 1])<<8
		sum += uint32(tmp)
	}

	sum = (sum>>16) + (sum & 0xffff)
	sum += (sum>>16)
	Trace("tcp checksum: %d", uint16(^sum))
	return uint16(^sum)
}
