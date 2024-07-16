package db

import (
	"context"

	"github.com/babylonchain/staking-api-service/internal/db/model"
	"go.mongodb.org/mongo-driver/bson"
)

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

func (db *Database) GetDApp(ctx context.Context) ([]*model.DAppDocument, error) {
	dApps := db.Client.Database(db.DbName).Collection(model.DAppCollection)
	// get the cursor to iterator over the dApps
	cursor, err := dApps.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var dAppDocuments []*model.DAppDocument
	for cursor.Next(ctx) {
		var dApp model.DAppDocument
		err := cursor.Decode(&dApp)
		if err != nil {
			return nil, err
		}
		dAppDocuments = append(dAppDocuments, &dApp)
	}
	return dAppDocuments, nil
}
