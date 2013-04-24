package eyego

import (
	"io"
	"encoding/hex"
	"fmt"
	"crypto/md5"
)

type ChecksumReader struct {
	delegate io.Reader
	state *checksumState
}

type checksumState struct {
	buf [/*512*/]byte
	ptr int
	checksums []uint16
}

func NewChecksumReader(r io.Reader) ChecksumReader {
	cr := ChecksumReader {
		delegate: r,
		state: &checksumState{
			buf: make([]byte, 512),
			checksums: make([]uint16, 0, 16)}}

	return cr
}

func (cr ChecksumReader) Read(p []byte) (n int, err error){
	n, err = cr.delegate.Read(p)

	if err == nil {
		cr.appendBytes(p, n)
	}

	return n, err
}

func (cr ChecksumReader) Checksum(uploadKey string) string {

	b, _ := hex.DecodeString(uploadKey)

	if len(b) % 2 != 0 { panic("Bad upload key")}

	h := md5.New()

	for i := 0; i < len(cr.state.checksums); i++ {
		h.Write([]byte{
			byte(cr.state.checksums[i]),
			byte(cr.state.checksums[i] >> 8)})
	}

	h.Write(b)

	return fmt.Sprintf("%x", h.Sum(nil))
}

func (cr ChecksumReader) appendBytes(b []byte, len int) {

	if cr.state.ptr + len >= 512 {
		Trace("Added block")
		copy(cr.state.buf[cr.state.ptr:512], b[0:512-cr.state.ptr]) //copy bytes to fill temp buffer
		cr.state.checksums = append(cr.state.checksums, tcp_checksum(cr.state.buf))
		cr.state.buf = cr.state.buf[:0]
		copy(cr.state.buf, b[512-cr.state.ptr:len]) //copy remaining bytes
		cr.state.ptr = len - (512 - cr.state.ptr)
	} else { //
		Trace("Added %d bytes to (%d byte) buffer", len, cr.state.ptr)
		copy(cr.state.buf[cr.state.ptr:cr.state.ptr+len], b[0:len])
		cr.state.ptr += len
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
