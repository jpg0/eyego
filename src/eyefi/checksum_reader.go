package eyefi

import (
	"io"
	"encoding/hex"
	"fmt"
	"crypto/md5"
)

type ChecksumReader struct {
	delegate io.Reader
	buf [/*512*/]byte
	ptr int
	checksums []uint16
	self *ChecksumReader
}

func NewChecksumReader(r io.Reader) ChecksumReader {
	cr := ChecksumReader {
		delegate: r,
		ptr: 0,
		buf: make([]byte, 512),
		checksums: make([]uint16, 0, 16)}
	cr.self = &cr
	return cr
}

func (cr ChecksumReader) Read(p []byte) (n int, err error){
	n, err = cr.delegate.Read(p)

	if err == nil {
		cr.appendBytes(p, n)
	}

	return n, nil
}

func (cr ChecksumReader) Checksum(uploadKey string) string {

	b, _ := hex.DecodeString(uploadKey)

	if len(b) % 2 != 0 { panic("Bad upload key")}

	h := md5.New()

	for i := 0; i < len(cr.checksums); i++ {
		h.Write([]byte{
			byte(cr.checksums[i]),
			byte(cr.checksums[i] >> 8)})
	}

	h.Write(b)

	return fmt.Sprintf("%x", h.Sum(nil))
}

func (cr ChecksumReader) appendBytes(b []byte, len int) {

	cr = *cr.self
	if cr.ptr + len >= 512 {
		copy(cr.buf[cr.ptr:512], b[0:512-cr.ptr]) //copy bytes to fill temp buffer
		cr.checksums = append(cr.checksums, tcp_checksum(cr.buf))
		cr.buf = cr.buf[:0]
		copy(cr.buf, b[512-cr.ptr:len]) //copy remaining bytes
		cr.ptr = len - (512 - cr.ptr)
	} else { //
		copy(cr.buf[cr.ptr:cr.ptr+len], b[0:len])
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
