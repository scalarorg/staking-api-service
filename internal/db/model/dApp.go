package model

type DAppDocument struct {
	ID                string                `bson:"_id,omitempty"`
	Description       dAppDescriptionPublic `bson:"description"`
	ChainName         string                `bson:"chain_name"`
	BTCAddress        string                `bson:"btc_address"`
	PublicKeyHex      string                `bson:"public_key_hex"`
	State             bool                  `bson:"state"`
	ActiveTvl         int64                 `json:"active_tvl"`
	TotalTvl          int64                 `json:"total_tvl"`
	ActiveDelegations int64                 `json:"active_delegations"`
	TotalDelegations  int64                 `json:"total_delegations"`
}
type dAppDescriptionPublic struct {
	Moniker         string `json:"moniker"`
	Identity        string `json:"identity"`
	Website         string `json:"website"`
	SecurityContact string `json:"security_contact"`
	Details         string `json:"details"`
}

type DAppStatsPagination struct {
	PublicKeyHex string `json:"public_key_hex"`
	ActiveTvl    int64  `json:"active_tvl"`
}

func BuildDAppStatsPaginationToken(d *DAppDocument) (string, error) {
	page := DAppStatsPagination{
		PublicKeyHex: d.PublicKeyHex,
		ActiveTvl:    d.ActiveTvl,
	}
	token, err := GetPaginationToken(page)
	if err != nil {
		return "", err
	}
	return token, nil
}
