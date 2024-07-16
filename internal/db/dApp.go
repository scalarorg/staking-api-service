package db

import (
	"context"

	"github.com/babylonchain/staking-api-service/internal/db/model"
	"go.mongodb.org/mongo-driver/bson"
)

func (db *Database) IsDAppExist(ctx context.Context, chainName, addressHex, publicKeyHex string) error {
	dApps := db.Client.Database(db.DbName).Collection(model.DAppCollection)
	dAppsFilter := bson.M{
		"chain_name":      chainName,
		"btc_address_hex": addressHex,
		"public_key":      publicKeyHex,
	}
	// Check if the dApp already exists
	err := dApps.FindOne(ctx, dAppsFilter).Err()
	// If the dApp already exists, return an error
	if err != nil {
		return err
	}
	return nil
}
func (db *Database) SaveDApp(ctx context.Context, chainName, addressHex, publicKeyHex string) error {
	dApps := db.Client.Database(db.DbName).Collection(model.DAppCollection)
	dApp := model.DAppDocument{
		ChainName:     chainName,
		BTCAddressHex: addressHex,
		PublicKeyHex:  publicKeyHex,
		State:         true,
	}
	// insert unique dApp
	_, err := dApps.InsertOne(ctx, dApp)
	return err
}
