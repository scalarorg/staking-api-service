package services

import (
	"context"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/scalarorg/staking-api-service/internal/db"
	"github.com/scalarorg/staking-api-service/internal/types"
)

func (s *Services) TransitionToBurningState(
	ctx context.Context, vaultTxHashHex string,
) *types.Error {
	err := s.DbClient.TransitionToBurningState(ctx, vaultTxHashHex)
	if err != nil {
		if ok := db.IsNotFoundError(err); ok {
			log.Ctx(ctx).Warn().Str("vaultTxHashHex", vaultTxHashHex).Err(err).Msg("delegation not found or no longer eligible for burning")
			return types.NewErrorWithMsg(http.StatusForbidden, types.NotFound, "delegation not found or no longer eligible for burning")
		}
		log.Ctx(ctx).Error().Str("vaultTxHashHex", vaultTxHashHex).Err(err).Msg("failed to transition to burning state")
		return types.NewError(http.StatusInternalServerError, types.InternalServiceError, err)
	}
	return nil
}

func (s *Services) TransitionToSlashingOrLostKeyState(
	ctx context.Context, vaultTxHashHex string,
) *types.Error {
	err := s.DbClient.TransitionToSlashingOrLostKeyState(ctx, vaultTxHashHex)
	if err != nil {
		if ok := db.IsNotFoundError(err); ok {
			log.Ctx(ctx).Warn().Str("vaultTxHashHex", vaultTxHashHex).Err(err).Msg("delegation not found or no longer eligible for slashing or lost key")
			return types.NewErrorWithMsg(http.StatusForbidden, types.NotFound, "delegation not found or no longer eligible for slashing or lost key")
		}
		log.Ctx(ctx).Error().Str("vaultTxHashHex", vaultTxHashHex).Err(err).Msg("failed to transition to slashing or lost key state")
		return types.NewError(http.StatusInternalServerError, types.InternalServiceError, err)
	}
	return nil
}

func (s *Services) TransitionToBurnWithoutDAppState(
	ctx context.Context, vaultTxHashHex string,
) *types.Error {
	err := s.DbClient.TransitionToBurnWithoutDAppState(ctx, vaultTxHashHex)
	if err != nil {
		if ok := db.IsNotFoundError(err); ok {
			log.Ctx(ctx).Warn().Str("vaultTxHashHex", vaultTxHashHex).Err(err).Msg("delegation not found or no longer eligible for burn without dapp state")
			return types.NewErrorWithMsg(http.StatusForbidden, types.NotFound, "delegation not found or no longer eligible for burn without dapp state")
		}
		log.Ctx(ctx).Error().Str("vaultTxHashHex", vaultTxHashHex).Err(err).Msg("failed to transition to burn without dapp state")
		return types.NewError(http.StatusInternalServerError, types.InternalServiceError, err)
	}
	return nil
}
