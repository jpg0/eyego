package eyefi

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"
)

type Card struct {
	UploadKey string `json:"upload_key"`
	MacAddress string `json:"mac_address"`
}

var cards = map[string]Card{}

func AddCardConfigs(configs []CardConfig) {
	for i := range configs {
		AddCard(Card{MacAddress: configs[i].MacAddress, UploadKey: configs[i].UploadKey})
	}
}

func GenerateSNonce() string {
	src := rand.NewSource(time.Now().UnixNano())
	p := make([]byte, 16)
	for i := range p {
		p[i] = byte(src.Int63() & 0xff)
	}
	return fmt.Sprintf("%x", p)
}

func AddCard(card Card) {
	cards[card.MacAddress] = card
}

func GetCard(mac_address string) Card {
	return cards[mac_address]
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
