package eyefi

import "io"
import (
	"encoding/hex"
	"fmt"
//	"encoding/binary"
//	"bytes"
)

type ChecksumReader struct {
	delegate io.Reader
	buf [/*512*/]byte
	ptr int
	checksums []uint16
}

func NewChecksumReader(r io.Reader) ChecksumReader {
	return ChecksumReader {
		delegate: r,
		ptr: 0,
		checksums: make([]uint16, 16)}
}

func (cr ChecksumReader) Read(p []byte) (n int, err error){
	n, err = cr.Read(p)

	if err != nil {
		cr.appendBytes(p,n)
	}

	return n, nil
}

func (cr ChecksumReader) Checksum(uploadKey string) {
	b, _ := hex.DecodeString(uploadKey)

	if len(b) % 2 != 0 { panic("Bad upload key")}

	cs := make([]uint16, len(cr.checksums))
	copy(cr.checksums, cs)

//	for c := 0; c < len(b); c = c + 2 {
//		cs = append(cs, int16(b[c:c+1]))
//	}

}

func (cr ChecksumReader) appendBytes(b []byte, len int) {

	if cr.ptr + len >= 512 {
		copy(b[0:512-cr.ptr], cr.buf[cr.ptr:512])
		cr.checksums = append(cr.checksums, tcp_checksum(cr.buf))
		copy(b[512-cr.ptr:len], cr.buf[0:len-(512-cr.ptr)])
		cr.ptr = len - (512 - cr.ptr)
	} else { //
		copy(b[0:len], cr.buf[cr.ptr:cr.ptr+len])
		cr.ptr += len
	}
}

func tcp_checksum(b []byte) uint16 {
	if len(b) % 2 != 0 { panic(fmt.Sprintf("tcp checksum bad length: %d", len(b))) }

	var sum uint32 = 0
	var tmp uint16

	for c := 0; c < len(b); c = c + 2 {
		tmp = uint16(b[c]) | uint16(b[c+1]) << 8
		sum += uint32(tmp)
	}

	sum = (sum >> 16) + (sum & 0xffff)
	sum += (sum >> 16)
	return uint16(^sum)
}
