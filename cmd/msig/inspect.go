package msig

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/coffeecloudgit/filecoin-wallet-signing/chain/types"
	multisig2 "github.com/filecoin-project/specs-actors/v2/actors/builtin/multisig"
	"io"
	"net/http"
	"time"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/specs-actors/v8/actors/builtin/miner"
	"github.com/filecoin-project/specs-actors/v8/actors/builtin/multisig"
	"github.com/filecoin-project/specs-actors/v8/actors/util/adt"
	cbor "github.com/ipfs/go-ipld-cbor"
	"github.com/spf13/cobra"
	"golang.org/x/xerrors"

	"github.com/coffeecloudgit/filecoin-wallet-signing/chain/blockstore"
	"github.com/coffeecloudgit/filecoin-wallet-signing/internal"
	"github.com/coffeecloudgit/filecoin-wallet-signing/pkg"
)

// inspectCmd represents the msiginspect command
var inspectCmd = &cobra.Command{
	Use:   "inspect <multisigAddress> ",
	Short: "inspect multisigAddress ",

	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 1 {
			_ = cmd.Help()
			return
		}

		//mtsaddr, err := address.NewFromString("t2i35vaqpkqpx3rcmqpttayaa3k4b7qm2fgrqiq3q")
		mtsaddr, err := address.NewFromString(args[0])
		if err != nil {
			fmt.Println("decode multisigAddress failed:", err.Error())
			return
		}

		if mtsaddr.Protocol() != address.Actor && mtsaddr.Protocol() != address.ID {
			fmt.Println("please input a correct multisigAddress")
			return
		}

		multisigID, err := internal.Lapi.StateLookupID(internal.Ctx, mtsaddr, *internal.CurrentTsk)
		if err != nil {
			fmt.Println("get address ID failed:", err.Error())
			return
		}

		fmt.Printf("Address: %s, ID: %s \n", mtsaddr.String(), multisigID.String())

		a, err := internal.Lapi.StateGetActor(internal.Ctx, mtsaddr, *internal.CurrentTsk)
		if err != nil {
			fmt.Println("Failed to get the address information:", err.Error())
			return
		}

		hd, err := internal.Lapi.ChainReadObj(internal.Ctx, a.Head)
		if err != nil {
			fmt.Println("Failed to get the address HEAD:", err.Error())
			return
		}

		var mstate multisig.State

		err = mstate.UnmarshalCBOR(bytes.NewReader(hd))
		if err != nil {
			fmt.Println("unmarshal address state failed:", err.Error())
			return
		}

		fmt.Printf("Number of signatories %v threshold  %v \n", len(mstate.Signers), mstate.NumApprovalsThreshold)
		for _, signer := range mstate.Signers {
			signerAddr, err := internal.Lapi.StateAccountKey(internal.Ctx, signer, *internal.CurrentTsk)
			if err != nil {
				fmt.Println("get singer of multisigAddress failed : ", err.Error())
				return
			}
			fmt.Printf("%s : %s \n", signer.String(), signerAddr.String())
		}

		store := adt.WrapStore(internal.Ctx, cbor.NewCborStore(blockstore.NewAPIBlockstore(internal.Lapi)))

		arr, err := adt.AsMap(store, mstate.PendingTxns, 5)
		if err != nil {
			fmt.Println("map address pending transaction failed:", err.Error())
			return
		}
		ks, err := arr.CollectKeys()
		if err != nil {
			fmt.Println("Collect address pending transaction failed:", err.Error())
			return
		}
		if len(ks) == 0 {
			fmt.Println("No pending transactions")
			return
		}
		fmt.Println("Pending transaction: ")
		var out multisig.Transaction
		err = arr.ForEach(&out, func(key string) error {
			txid, n := binary.Varint([]byte(key))
			if n <= 0 {
				return xerrors.Errorf("invalid pending transaction key: %v", key)
			}
			p := ""
			msg := ""
			var mwdp miner.WithdrawBalanceParams
			msg = "send out"
			if out.Method == 16 {
				err = mwdp.UnmarshalCBOR(bytes.NewReader(out.Params))
				if err != nil {
					fmt.Println("Parameter parsing failed:", err.Error())
					return nil
				}
				b, _ := json.Marshal(mwdp)
				p = string(b)
				msg = fmt.Sprintf("withdraw from miner  %v FIL", pkg.ToFloat64(mwdp.AmountRequested))
			}
			if out.Method == 23 {
				addr := address.Address{}
				err = addr.UnmarshalCBOR(bytes.NewReader(out.Params))
				if err != nil {
					fmt.Println("Parameter parsing failed:", err.Error())
					return nil
				}

				msg = fmt.Sprintf("change miner %v owner is %v ", out.To.String(), addr.String())
			}
			fmt.Printf("pending id: %v , to : %v , method: %v , amount: %v FIL, Params: %s, approved %v, ps: %s \n",
				txid, out.To, out.Method, pkg.ToFloat64(out.Value), p, out.Approved, msg)
			return nil
		})
		if err != nil {
			fmt.Println("get address pinding transation failed:", err.Error())
			return
		}
	},
}

func GetMultiAccountInfo(account string) (*MultiAccountInfo, error) {

	//mtsaddr, err := address.NewFromString("t2i35vaqpkqpx3rcmqpttayaa3k4b7qm2fgrqiq3q")
	mtsaddr, err := address.NewFromString(account)
	if err != nil {
		return nil, err
	}

	if mtsaddr.Protocol() != address.Actor && mtsaddr.Protocol() != address.ID {
		return nil, err
	}

	//multisigID, err := internal.Lapi.StateLookupID(internal.Ctx, mtsaddr, *internal.CurrentTsk)
	//if err != nil {
	//	return nil, err
	//}
	//
	//fmt.Printf("Address: %s, ID: %s \n", mtsaddr.String(), multisigID.String())

	internal.ChangeTks()

	a, err := internal.Lapi.StateGetActor(internal.Ctx, mtsaddr, *internal.CurrentTsk)
	if err != nil {
		fmt.Println("Failed to get the address information:", err.Error())
		return nil, err
	}

	hd, err := internal.Lapi.ChainReadObj(internal.Ctx, a.Head)
	if err != nil {
		fmt.Println("Failed to get the address HEAD:", err.Error())
		return nil, err
	}

	var mstate multisig.State

	err = mstate.UnmarshalCBOR(bytes.NewReader(hd))
	if err != nil {
		fmt.Println("unmarshal address state failed:", err.Error())
		return nil, err
	}

	fmt.Printf("Number of signatories %v threshold  %v \n", len(mstate.Signers), mstate.NumApprovalsThreshold)
	//for _, signer := range mstate.Signers {
	//	signerAddr, err := internal.Lapi.StateAccountKey(internal.Ctx, signer, *internal.CurrentTsk)
	//	if err != nil {
	//		fmt.Println("get singer of multisigAddress failed : ", err.Error())
	//		return nil, err
	//	}
	//	fmt.Printf("%s : %s \n", signer.String(), signerAddr.String())
	//}

	store := adt.WrapStore(internal.Ctx, cbor.NewCborStore(blockstore.NewAPIBlockstore(internal.Lapi)))

	arr, err := adt.AsMap(store, mstate.PendingTxns, 5)
	if err != nil {
		fmt.Println("map address pending transaction failed:", err.Error())
		return nil, err
	}
	//ks, err := arr.CollectKeys()
	//if err != nil {
	//	fmt.Println("Collect address pending transaction failed:", err.Error())
	//	return nil, err
	//}
	//if len(ks) == 0 {
	//	fmt.Println("No pending transactions")
	//	return nil, err
	//}
	//fmt.Println("Pending transaction: ")
	var multiSigPendingTxs = make([]MultiSignTx, 0)
	var out multisig.Transaction
	err = arr.ForEach(&out, func(key string) error {
		txid, n := binary.Varint([]byte(key))
		if n <= 0 {
			return xerrors.Errorf("invalid pending transaction key: %v", key)
		}
		p := ""
		msg := ""
		var mwBp miner.WithdrawBalanceParams
		var addSp multisig2.AddSignerParams
		var removeSp multisig2.AddSignerParams
		var cnatP multisig2.ChangeNumApprovalsThresholdParams
		msg = "send out"

		if out.Method == 5 {
			err = addSp.UnmarshalCBOR(bytes.NewReader(out.Params))
			if err != nil {
				fmt.Println("Parameter parsing failed:", err.Error())
				return nil
			}
			b, _ := json.Marshal(addSp)
			p = string(b)
			msg = fmt.Sprintf("Add Signer  %v", addSp.Signer.String())
		}

		if out.Method == 6 {
			err = removeSp.UnmarshalCBOR(bytes.NewReader(out.Params))
			if err != nil {
				fmt.Println("Parameter parsing failed:", err.Error())
				return nil
			}
			b, _ := json.Marshal(removeSp)
			p = string(b)
			msg = fmt.Sprintf("Remove Signer  %v", removeSp.Signer.String())
		}

		if out.Method == 8 {
			err = cnatP.UnmarshalCBOR(bytes.NewReader(out.Params))
			if err != nil {
				fmt.Println("Parameter parsing failed:", err.Error())
				return nil
			}
			b, _ := json.Marshal(cnatP)
			p = string(b)
			msg = fmt.Sprintf("Change threshold to  %v", cnatP.NewThreshold)
		}

		if out.Method == 16 {
			err = mwBp.UnmarshalCBOR(bytes.NewReader(out.Params))
			if err != nil {
				fmt.Println("Parameter parsing failed:", err.Error())
				return nil
			}
			b, _ := json.Marshal(mwBp)
			p = string(b)
			msg = fmt.Sprintf("withdraw from miner  %v FIL", pkg.ToFloat64(mwBp.AmountRequested))
		}
		if out.Method == 23 {
			addr := address.Address{}
			err = addr.UnmarshalCBOR(bytes.NewReader(out.Params))
			if err != nil {
				fmt.Println("Parameter parsing failed:", err.Error())
				return nil
			}

			msg = fmt.Sprintf("change miner %v owner is %v ", out.To.String(), addr.String())
		}
		//fmt.Printf("pending id: %v , to : %v , method: %v , amount: %v FIL, Params: %s, approved %v, ps: %s \n",
		//txid, out.To, out.Method, pkg.ToFloat64(out.Value), p, out.Approved, msg)

		multiSigPendingTxs = append(multiSigPendingTxs,
			MultiSignTx{Id: txid, To: out.To.String(), Method: out.Method.String(),
				Mount: pkg.ToFloat64(out.Value), Params: p, Approved: out.Approved, Ps: msg})
		return nil
	})
	multiSignInfo := MultiAccountInfo{Signers: mstate.Signers, NumApprovalsThreshold: mstate.NumApprovalsThreshold,
		InitialBalance: mstate.InitialBalance, StartEpoch: mstate.StartEpoch, UnlockDuration: mstate.UnlockDuration, MultiSignTxs: multiSigPendingTxs}
	return &multiSignInfo, err
}

func GetAccountInfo(account string) (*AccountInfo, error) {
	accountAddr, err := address.NewFromString(account)
	if err != nil {
		return nil, err
	}
	tipSet, err := internal.Lapi.ChainHead(internal.Ctx)

	if err != nil {
		return nil, err
	}

	Tsk := tipSet.Key()

	balance, err := internal.Lapi.WalletBalance(internal.Ctx, accountAddr)

	if err != nil {
		return nil, err
	}

	result, err := internal.Lapi.StateLookupID(internal.Ctx, accountAddr, Tsk)

	if err != nil {
		return nil, err
	}

	return &AccountInfo{Address: accountAddr, Id: result.String(), Height: tipSet.Height(), Balance: balance}, nil
}

func GetWalletBalance(accountAddr address.Address) (types.BigInt, error) {

	return internal.Lapi.WalletBalance(internal.Ctx, accountAddr)

}

func GetActorAddress(address string) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	client := http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(fmt.Sprintf("https://api.filutils.com/api/v2/actor/%s", address))

	if err != nil {
		return data, err
	}

	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return data, err
	}

	// 解析JSON数据
	if err := json.Unmarshal(body, &data); err != nil {
		return data, err
	}
	return data, nil
}
