package db

import (
	"context"

	"github.com/babylonchain/staking-api-service/internal/db/model"
	"go.mongodb.org/mongo-driver/bson"
)

func (db *Database) IsDAppExist(ctx context.Context, addressHex, publicKeyHex, chainName string) error {
	dApps := db.Client.Database(db.DbName).Collection(model.DAppCollection)
	dAppsFilter := bson.M{
		"address":    addressHex,
		"public_key": publicKeyHex,
		"chain_name": chainName,
	}
	// Check if the dApp already exists
	err := dApps.FindOne(ctx, dAppsFilter).Err()
	// If the dApp already exists, return an error
	if err != nil {
		return err
	}
	return nil
}
