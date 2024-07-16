package services

import (
	"context"
	"net/http"

	"github.com/babylonchain/staking-api-service/internal/types"
)

func (s *Services) CreateDApp(ctx context.Context, chainName, btcAddressHex, publicKeyHex string) *types.Error {
	err := s.DbClient.IsDAppExist(ctx, chainName, btcAddressHex, publicKeyHex)
	if err != nil {
		return types.NewError(http.StatusInternalServerError, types.InternalServiceError, err)
	}
	err = s.DbClient.SaveDApp(ctx, chainName, btcAddressHex, publicKeyHex)
	if err != nil {
		return types.NewError(http.StatusInternalServerError, types.InternalServiceError, err)
	}
	return nil

}
