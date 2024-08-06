package services

import (
	"context"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/scalarorg/staking-api-service/internal/db"
	"github.com/scalarorg/staking-api-service/internal/db/model"
	"github.com/scalarorg/staking-api-service/internal/types"
)

func (s *Services) CreateDApp(ctx context.Context, chainName, btcAddress, publicKeyHex string) *types.Error {
	err := s.DbClient.SaveDApp(ctx, chainName, btcAddress, publicKeyHex)
	if err != nil {
		return types.NewError(http.StatusInternalServerError, types.InternalServiceError, err)
	}
	return nil

}

func (s *Services) GetDApp(ctx context.Context, page string) ([]*model.DAppDocument, string, *types.Error) {
	dApps, err := s.DbClient.GetDApp(ctx)
	if err != nil {
		return nil, "", types.NewError(http.StatusInternalServerError, types.InternalServiceError, err)
	}
	if len(dApps) == 0 {
		log.Ctx(ctx).Info().Msg("No dApps found")
		return nil, "", types.NewErrorWithMsg(http.StatusInternalServerError, types.InternalServiceError, "No dApp found")
	}
	dAppsMap := make(map[string]*model.DAppDocument)
	for _, d := range dAppsMap {
		dAppsMap[d.PublicKeyHex] = d
	}

	resultMap, err := s.DbClient.FindDAppStats(ctx, page)
	if err != nil {
		if db.IsInvalidPaginationTokenError(err) {
			log.Ctx(ctx).Warn().Err(err).Msg("Invalid pagination token when fetching dApp stats")
			return nil, "", types.NewError(http.StatusBadRequest, types.BadRequest, err)
		}
		// We don't want to return an error here in case of DB error.
		// we will continue the process with the data we have from global params as a fallback.
		// TODO: Add metric for this error and alerting
		log.Ctx(ctx).Error().Err(err).Msg("Error while fetching dApp from DB")
		// Return the dApp from global params as a fallback
		return buildFallbackDAppDetailsPublic(dApps), "", nil
	}

	/*
		If no dApp are found in the database and no pagination token
		is provided (indicating this is the first page), return the dApp
		from the global parameters as a fallback. This fallback is only necessary when
		launching the service for the first time and no dApp are found in the database.
	*/
	if (len(resultMap.Data) == 0) && (page == "") {
		return buildFallbackDAppDetailsPublic(dApps), "", nil
	}

	var dAppsDetails []*model.DAppDocument
	for _, d := range resultMap.Data {
		detail := &model.DAppDocument{
			ID:                d.ID,
			Description:       d.Description,
			ChainName:         d.ChainName,
			BTCAddress:        d.BTCAddress,
			PublicKeyHex:      d.PublicKeyHex,
			State:             d.State,
			ActiveTvl:         d.ActiveTvl,
			TotalTvl:          d.TotalTvl,
			ActiveDelegations: d.ActiveDelegations,
			TotalDelegations:  d.TotalDelegations,
		}
		dAppsDetails = append(dAppsDetails, detail)
	}
	return dAppsDetails, resultMap.PaginationToken, nil
}

func (s *Services) UpdateDApp(ctx context.Context, ID, chainName, btcAddress, publicKeyHex string) *types.Error {
	err := s.DbClient.UpdateDApp(ctx, ID, chainName, btcAddress, publicKeyHex)
	if err != nil {
		return types.NewError(http.StatusInternalServerError, types.InternalServiceError, err)
	}
	return nil
}

func (s *Services) ToggleDApp(ctx context.Context, ID string) *types.Error {
	err := s.DbClient.ToggleDApp(ctx, ID)
	if err != nil {
		return types.NewError(http.StatusInternalServerError, types.InternalServiceError, err)
	}
	return nil
}

func (s *Services) DeleteDApp(ctx context.Context, ID string) *types.Error {
	err := s.DbClient.DeleteDApp(ctx, ID)
	if err != nil {
		return types.NewError(http.StatusInternalServerError, types.InternalServiceError, err)
	}
	return nil
}

func buildFallbackDAppDetailsPublic(dApps []*model.DAppDocument) []*model.DAppDocument {
	var dAppsDetails []*model.DAppDocument
	for _, d := range dApps {
		detail := &model.DAppDocument{
			Description:       d.Description,
			ChainName:         d.ChainName,
			BTCAddress:        d.BTCAddress,
			PublicKeyHex:      d.PublicKeyHex,
			State:             d.State,
			ActiveTvl:         0,
			TotalTvl:          0,
			ActiveDelegations: 0,
			TotalDelegations:  0,
		}
		dAppsDetails = append(dAppsDetails, detail)
	}
	return dAppsDetails
}
