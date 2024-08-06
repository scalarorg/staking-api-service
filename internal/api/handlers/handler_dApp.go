package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/scalarorg/staking-api-service/internal/types"
	"github.com/scalarorg/staking-api-service/internal/utils"
)

type CreateDAppRequestPayload struct {
	ChainName    string `json:"chain_name"`
	BTCAddress   string `json:"btc_address"`
	PublicKeyHex string `json:"public_key_hex"`
}
type UpdateDAppRequestPayload struct {
	ID            string `json:"id"`
	ChainName     string `json:"chain_name"`
	BTCAddressHex string `json:"btc_address_hex"`
	PublicKeyHex  string `json:"public_key_hex"`
}

type IdRequestPayload struct {
	ID string `json:"id"`
}

func parseCreateDAppPayload(request *http.Request, netParam *chaincfg.Params) (*CreateDAppRequestPayload, *types.Error) {
	payload := &CreateDAppRequestPayload{}
	err := json.NewDecoder(request.Body).Decode(payload)
	if err != nil {
		return nil, types.NewErrorWithMsg(http.StatusBadRequest, types.BadRequest, "invalid request payload")
	}
	if !utils.IsValidChainName(payload.ChainName) {
		return nil, types.NewErrorWithMsg(
			http.StatusBadRequest, types.BadRequest, "invalid chain name",
		)
	}
	err = utils.IsValidBtcAddress(payload.BTCAddress, netParam)
	if err != nil {
		return nil, types.NewErrorWithMsg(
			http.StatusBadRequest, types.BadRequest, "invalid btc address hex",
		)
	}
	if !utils.IsValidPublicKeyHex(payload.PublicKeyHex) {
		return nil, types.NewErrorWithMsg(
			http.StatusBadRequest, types.BadRequest, "invalid public key hex",
		)
	}
	return payload, nil
}

func parseUpdateDAppPayload(request *http.Request) (*UpdateDAppRequestPayload, *types.Error) {
	payload := &UpdateDAppRequestPayload{}
	err := json.NewDecoder(request.Body).Decode(payload)
	if err != nil {
		return nil, types.NewErrorWithMsg(http.StatusBadRequest, types.BadRequest, "invalid request payload")
	}
	return payload, nil
}

func parseIdDAppPayload(request *http.Request) (*IdRequestPayload, *types.Error) {
	payload := &IdRequestPayload{}
	err := json.NewDecoder(request.Body).Decode(payload)
	if err != nil {
		return nil, types.NewErrorWithMsg(http.StatusBadRequest, types.BadRequest, "invalid request payload")
	}
	return payload, nil
}
