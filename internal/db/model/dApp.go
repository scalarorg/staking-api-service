package model

type DAppDocument struct {
	AddressHex   string `bson:"address_hex"`
	PublicKeyHex string `bson:"public_key`
	ChainName    string `bson:"chain_name"`
	State        bool   `bson:"state"`
}
