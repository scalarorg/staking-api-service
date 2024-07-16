package services

import (
	"context"
	"net/http"

	"github.com/babylonchain/staking-api-service/internal/types"
)

func (s *Services) CreateDApp(ctx context.Context, chainName, addressHex, publicKeyHex string) *types.Error {
	err := s.DbClient.IsDAppExist(ctx, addressHex, publicKeyHex, chainName)
	if err != nil {
		return types.NewError(http.StatusInternalServerError, types.InternalServiceError, err)
	}
	// err2 := s.DbClient.SaveDApp(ctx, chainName, addressHex, publicKeyHex)
	// if err2 != nil {
	// 	return types.NewError(http.StatusInternalServerError, types.InternalServiceError, err)
	// }
	return nil

}
