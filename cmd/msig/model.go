package msig

import (
	"github.com/coffeecloudgit/filecoin-wallet-signing/chain/types"
	addr "github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
)

type AccountInfo struct {
	Address addr.Address   `json:"address"`
	Id      string         `json:"id"`
	Height  abi.ChainEpoch `json:"height"`
	Balance types.BigInt   `json:"balance"`
}

type MultiAccountInfo struct {
	Signers               []addr.Address  `json:"signers"`
	NumApprovalsThreshold uint64          `json:"numApprovalsThreshold"`
	InitialBalance        abi.TokenAmount `json:"initialBalance"`
	StartEpoch            abi.ChainEpoch  `json:"startEpoch"`
	UnlockDuration        abi.ChainEpoch  `json:"unlockDuration"`
	MultiSignTxs          []MultiSignTx   `json:"multiSignTxs"`
}

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
