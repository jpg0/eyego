package eyego

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"
)

func GenerateSNonce() string {
	src := rand.NewSource(time.Now().UnixNano())
	p := make([]byte, 16)
	for i := range p {
		p[i] = byte(src.Int63() & 0xff)
	}
	return fmt.Sprintf("%x", p)
}

func (c Card) Credential(cnonce string) string {
	if len(cnonce) == 0 {
		panic("no cnonce")
	}

	h := md5.New()
	binary, _ := hex.DecodeString(c.MacAddress + cnonce + c.UploadKey)
	h.Write(binary)
	return fmt.Sprintf("%x", h.Sum(nil))
}
