package db

import (
	"context"

	"github.com/scalarorg/staking-api-service/internal/types"
	"github.com/scalarorg/staking-api-service/internal/utils"
)

func (db *Database) TransitionToBurningState(ctx context.Context, txHashHex string) error {
	err := db.transitionState(
		ctx, txHashHex, types.Burned.ToString(),
		utils.QualifiedStatesToBurning(), nil,
	)
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) TransitionToSlashingOrLostKeyState(ctx context.Context, txHashHex string) error {
	err := db.transitionState(
		ctx, txHashHex, types.Slashed.ToString(),
		utils.QualifiedStatesToSlashingOrLostKey(), nil,
	)
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) TransitionToBurnWithoutDAppState(ctx context.Context, txHashHex string) error {
	err := db.transitionState(
		ctx, txHashHex, types.WithoutDAppBurned.ToString(),
		utils.QualifiedStatesToBurnWithoutDApp(), nil,
	)
	if err != nil {
		return err
	}
	return nil
}
