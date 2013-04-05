package eyefi

import (
"crypto/md5"
	"encoding/hex"
	"fmt"
)

type Card struct {
	upload_key string
	mac_address string
}

var cards = map[string]Card{}

func AddCard(card Card) {
	cards[card.mac_address] = card
}

func GetCard(mac_address string) Card {
	return cards[mac_address]
}

func (c Card) Credential(cnonce string) string {
	h := md5.New()
	binary, _ := hex.DecodeString(c.mac_address + cnonce + c.upload_key)
	h.Write(binary)
	return fmt.Sprintf("%x", h.Sum(nil))
}
