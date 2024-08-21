package api

import (
	_ "github.com/babylonchain/staking-api-service/docs"
	"github.com/go-chi/chi"
	httpSwagger "github.com/swaggo/http-swagger"
)

func (a *Server) SetupRoutes(r *chi.Mux) {
	handlers := a.handlers
	r.Get("/healthcheck", registerHandler(handlers.HealthCheck))

	r.Get("/v1/staker/delegations", registerHandler(handlers.GetStakerDelegations))
	r.Post("/v1/unbonding", registerHandler(handlers.UnbondDelegation))
	r.Get("/v1/unbonding/eligibility", registerHandler(handlers.GetUnbondingEligibility))
	r.Get("/v1/global-params", registerHandler(handlers.GetBabylonGlobalParams))
	r.Get("/v1/finality-providers", registerHandler(handlers.GetFinalityProviders))
	r.Get("/v1/stats", registerHandler(handlers.GetOverallStats))
	r.Get("/v1/stats/staker", registerHandler(handlers.GetTopStakerStats))
	r.Get("/v1/staker/delegation/check", registerHandler(handlers.CheckStakerDelegationExist))
	r.Get("/v1/delegation", registerHandler(handlers.GetDelegationByTxHash))

	r.Get("/v1/dApp", registerHandler(handlers.GetDApp))
	r.Post("/v1/dApp", registerHandler(handlers.CreateDApp))
	r.Put("/v1/dApp", registerHandler(handlers.UpdateDApp))
	r.Patch("/v1/dApp", registerHandler(handlers.ToggleDApp))
	r.Delete("/v1/dApp", registerHandler(handlers.DeleteDApp))

	r.Get("/v1/gmp", registerHandler(handlers.GetGMPs))

	r.Get("/swagger/*", httpSwagger.WrapHandler)
}
