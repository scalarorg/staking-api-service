package handlers

import (
	"net/http"

	"github.com/scalarorg/staking-api-service/internal/types"
)

// CreateDApp creates a new dApp.
// @Summary Create dApp
// @Description Creates a new dApp with the provided chain name, BTC address hex, and public key hex.
// @Produce json
// @Param chain_name body string true "Chain name"
// @Param btc_address_hex body string true "BTC address hex"
// @Param public_key_hex body string true "Public key hex"
// @Success 200 {object} PublicResponse[CreateDAppRequestPayload] "The dApp has been created successfully"
// @Router /v1/dApp [post]
func (h *Handler) CreateDApp(request *http.Request) (*Result, *types.Error) {
	payload, err := parseCreateDAppPayload(request, h.config.Server.BTCNetParam)
	if err != nil {
		return nil, err
	}

	err = h.services.CreateDApp(request.Context(), payload.ChainName, payload.BTCAddress, payload.PublicKeyHex)
	if err != nil {
		return nil, err
	}

	return NewResult(payload), nil
}

// GetDApp gets active dApp sorted by ActiveTvl.
// @Summary Get Active dApp
// @Description Fetches details of all active dApp sorted by their active total value locked (ActiveTvl) in descending order.
// @Produce json
// @Param pagination_key query string false "Pagination key to fetch the next page of finality providers"
// @Success 200 {object} PublicResponse[[]services.FpDetailsPublic] "A list of dApp sorted by ActiveTvl in descending order"
// @Router /v1/dApp [get]
func (h *Handler) GetDApp(request *http.Request) (*Result, *types.Error) {
	paginationKey, err := parsePaginationQuery(request)
	if err != nil {
		return nil, err
	}
	dApps, paginationToken, err := h.services.GetDApp(request.Context(), paginationKey)
	if err != nil {
		return nil, err
	}
	return NewResultWithPagination(dApps, paginationToken), nil
}

// UpdateDApp updates an existing dApp.
// @Summary Update dApp
// @Description Updates an existing dApp with the provided chain name, BTC address hex, and public key hex.
// @Produce json
// @Param id body string true "ID of the dApp to update"
// @Param chain_name body string true "Chain name"
// @Param btc_address_hex body string true "BTC address hex"
// @Param public_key_hex body string true "Public key hex"
// @Success 200 {object} PublicResponse[UpdateDAppRequestPayload] "The dApp has been updated successfully"
// @Router /v1/dApp [put]
func (h *Handler) UpdateDApp(request *http.Request) (*Result, *types.Error) {
	payload, err := parseUpdateDAppPayload(request)
	if err != nil {
		return nil, err
	}
	err = h.services.UpdateDApp(request.Context(), payload.ID, payload.ChainName, payload.BTCAddressHex, payload.PublicKeyHex)
	if err != nil {
		return nil, err
	}

	return NewResult(payload), nil
}

// ToggleDApp toggles the state of an existing dApp.
// @Summary Toggle dApp
// @Description Toggles the state of an existing dApp between active and inactive.
// @Produce json
// @Param id body string true "ID of the dApp to toggle"
// @Success 200 {object} PublicResponse[IdRequestPayload] "The dApp has been toggled successfully"
// @Router /v1/dApp [patch]
func (h *Handler) ToggleDApp(request *http.Request) (*Result, *types.Error) {
	payload, err := parseIdDAppPayload(request)
	if err != nil {
		return nil, err
	}
	err = h.services.ToggleDApp(request.Context(), payload.ID)
	if err != nil {
		return nil, err
	}
	return NewResult(payload), nil
}

// DeleteDApp deletes an existing dApp.
// @Summary Delete dApp
// @Description Deletes an existing dApp.
// @Produce json
// @Param id body string true "ID of the dApp to delete"
// @Success 200 {object} PublicResponse[IdRequestPayload] "The dApp has been deleted successfully"
// @Router /v1/dApp [delete]
func (h *Handler) DeleteDApp(request *http.Request) (*Result, *types.Error) {
	payload, err := parseIdDAppPayload(request)
	if err != nil {
		return nil, err
	}
	err = h.services.DeleteDApp(request.Context(), payload.ID)
	if err != nil {
		return nil, err
	}

	return NewResult("Delete successfully"), nil
}
