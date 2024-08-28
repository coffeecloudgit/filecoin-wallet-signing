package msig

import (
	addr "github.com/filecoin-project/go-address"
)

type MultiSignTx struct {
	Id       int64          `json:"id"`
	To       string         `json:"to"`
	Method   string         `json:"method"`
	Mount    float64        `json:"mount"`
	Params   string         `json:"params"`
	Approved []addr.Address `json:"approved"`
	Ps       string         `json:"ps"`
	TxId     string         `json:"txId"`
}
