package db

import (
	"context"

	"github.com/scalarorg/staking-api-service/internal/db/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (db *Database) SaveDApp(ctx context.Context, chainName, address, publicKeyHex string) error {
	dApps := db.Client.Database(db.DbName).Collection(model.DAppCollection)
	dApp := model.DAppDocument{
		ChainName:    chainName,
		BTCAddress:   address,
		PublicKeyHex: publicKeyHex,
		State:        true,
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

func (db *Database) UpdateDApp(ctx context.Context, ID, chainName, address, publicKeyHex string) error {
	dApps := db.Client.Database(db.DbName).Collection(model.DAppCollection)
	// convert ID to objectID
	_id, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": _id}
	update := bson.M{"$set": bson.M{"chain_name": chainName, "btc_address": address, "public_key_hex": publicKeyHex}}
	updateResult := dApps.FindOneAndUpdate(ctx, filter, update)
	return updateResult.Err()
}

func (db *Database) ToggleDApp(ctx context.Context, ID string) error {
	dApps := db.Client.Database(db.DbName).Collection(model.DAppCollection)
	// convert ID to objectID
	_id, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": _id}
	var result model.DAppDocument
	err = dApps.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return err
	}
	update := bson.M{"$set": bson.M{"state": !result.State}}
	updateResult := dApps.FindOneAndUpdate(ctx, filter, update)
	return updateResult.Err()
}

func (db *Database) DeleteDApp(ctx context.Context, ID string) error {
	dApps := db.Client.Database(db.DbName).Collection(model.DAppCollection)
	// convert ID to objectID
	_id, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": _id}
	_, err = dApps.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

// FindDAppStats fetches the finality provider stats from the database
func (db *Database) FindDAppStats(ctx context.Context, paginationToken string) (*DbResultMap[*model.DAppDocument], error) {
	client := db.Client.Database(db.DbName).Collection(model.DAppStatsCollection)
	options := options.Find().SetSort(bson.D{{Key: "active_tvl", Value: -1}}) // Sorting in descending order
	var filter bson.M

	// Decode the pagination token first if it exist
	if paginationToken != "" {
		decodedToken, err := model.DecodePaginationToken[model.DAppStatsPagination](paginationToken)
		if err != nil {
			return nil, &InvalidPaginationTokenError{
				Message: "Invalid pagination token",
			}
		}
		filter = bson.M{
			"$or": []bson.M{
				{"active_tvl": bson.M{"$lt": decodedToken.ActiveTvl}},
				{"active_tvl": decodedToken.ActiveTvl, "_id": bson.M{"$lt": decodedToken.PublicKeyHex}},
			},
		}
	}

	return findWithPagination(
		ctx, client, filter, options, db.cfg.MaxPaginationLimit,
		model.BuildDAppStatsPaginationToken,
	)
}
