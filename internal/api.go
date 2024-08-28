package internal

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-jsonrpc"

	"github.com/coffeecloudgit/filecoin-wallet-signing/chain/api"
	"github.com/coffeecloudgit/filecoin-wallet-signing/chain/types"
)

const mainnetStartTimestamp = 1598306400

var Lapi *api.FullNodeStruct
var Ctx = context.Background()
var CurrentTsk *types.TipSetKey

func init() {
	err := ConnectLotus()
	if err != nil {
		panic(err)
	}
}

func ConnectLotus() error {
	Lapi = new(api.FullNodeStruct)
	lotusAPI := os.Getenv("LOTUS_API")
	token := "Bearer " + os.Getenv("LOTUS_API_TOKEN")
	if lotusAPI == "" {
		//lotusAPI = "https://api.node.glif.io/rpc/v1"
		lotusAPI = "https://node.filutils.com/rpc/v1"
		token = "Bearer " + "I6505IXYTJWI9LG3EZXQCO7FBZ846QPWXQX"
		//lotusAPI = "https://fil.coin.bit99.com:1443/rpc/v1"
	}

	fmt.Println("LOTUS_API : ", lotusAPI)
	fmt.Println("LOTUS_API_TOKEN : ", token)
	var headers http.Header
	if len(token) > 7 {
		headers = http.Header{
			"Authorization": []string{token},
			"content-type":  []string{"application/json"},
		}
	} else {
		headers = http.Header{
			"content-type": []string{"application/json"},
		}
	}

	closer, err := jsonrpc.NewMergeClient(context.Background(), lotusAPI, "Filecoin", []interface{}{&Lapi.Internal, &Lapi.CommonStruct.Internal}, headers)
	if err != nil {
		return fmt.Errorf("connecting with lotus failed: %s", err)
	}
	defer closer()

	gts, err := Lapi.ChainGetGenesis(Ctx)
	if err != nil {
		return fmt.Errorf("get genesis failed: %s", err.Error())
	}
	address.CurrentNetwork = address.Mainnet

	if gts.Blocks()[0].Timestamp != mainnetStartTimestamp {
		address.CurrentNetwork = address.Testnet
	}

	hts, _ := Lapi.ChainHead(Ctx)
	Tsk := hts.Key()
	CurrentTsk = &Tsk

	return nil
}
