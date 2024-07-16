package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/babylonchain/staking-api-service/internal/types"
	"github.com/babylonchain/staking-api-service/internal/utils"
)

type CreateDAppRequestPayload struct {
	ChainName     string `json:"chain_name"`
	BTCAddressHex string `json:"BTC_address_hex"`
	PublicKeyHex  string `json:"public_key_hex"`
}

func parseCreateDAppPayload(request *http.Request) (*CreateDAppRequestPayload, *types.Error) {
	payload := &CreateDAppRequestPayload{}
	err := json.NewDecoder(request.Body).Decode(payload)
	if err != nil {
		return nil, types.NewErrorWithMsg(http.StatusBadRequest, types.BadRequest, "invalid request payload")
	}
	// Validate the payload fields
	if !utils.IsValidChainName(payload.ChainName) {
		return nil, types.NewErrorWithMsg(
			http.StatusBadRequest, types.BadRequest, "invalid chain name",
		)
	}
	if !utils.IsValidAddressHex(payload.BTCAddressHex) {
		return nil, types.NewErrorWithMsg(
			http.StatusBadRequest, types.BadRequest, "invalid address hex",
		)
	}
	if !utils.IsValidPublickeyHex(payload.PublicKeyHex) {
		return nil, types.NewErrorWithMsg(
			http.StatusBadRequest, types.BadRequest, "invalid public key hex",
		)
	}

	return payload, nil
}

func (h *Handler) CreateDApp(request *http.Request) (*Result, *types.Error) {
	payload, err := parseCreateDAppPayload(request)
	if err != nil {
		return nil, err
	}

	err = h.services.CreateDApp(request.Context(), payload.ChainName, payload.BTCAddressHex, payload.PublicKeyHex)
	if err != nil {
		return nil, err
	}

	return NewResult(payload), nil
}

func (h *Handler) GetDApp(request *http.Request) (*Result, *types.Error) {
	// FUTURE WORK: Implement pagination
	// paginationKey, err := parsePaginationQuery(request)
	// if err != nil {
	// 	return nil, err
	// }
	dApps, err := h.services.GetDApp(request.Context())
	if err != nil {
		return nil, err
	}
	return NewResult(dApps), nil
}
