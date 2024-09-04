package msig

import (
	"fmt"
	"strconv"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/specs-actors/actors/builtin"
	"github.com/filecoin-project/specs-actors/v8/actors/builtin/multisig"
	"github.com/spf13/cobra"

	"github.com/coffeecloudgit/filecoin-wallet-signing/chain/actors"
	"github.com/coffeecloudgit/filecoin-wallet-signing/chain/types"
	"github.com/coffeecloudgit/filecoin-wallet-signing/internal"
	"github.com/coffeecloudgit/filecoin-wallet-signing/pkg"
)

// approveCmd represents the Approve command
var approveCmd = &cobra.Command{
	Use:   "approve <multisigAddress> <TxID>",
	Short: "approve  transaction of multisigAddress",

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			_ = cmd.Help()
			return
		}

		maddr, err := address.NewFromString(args[0])
		if err != nil {
			fmt.Println("decode multisigAddress failed:: ", err.Error())
			return
		}

		if maddr.Protocol() != address.Actor && maddr.Protocol() != address.ID {
			fmt.Println("please input a correct multisigAddress")
			return
		}

		txid, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("Transaction ID failed: ", err.Error())
			return
		}
		key, err := pkg.ReadPrivteKey()
		if err != nil {
			fmt.Println("decode private key failed: ", err)
			return
		}

		params, err := actors.SerializeParams(&multisig.TxnIDParams{ID: multisig.TxnID(txid)})
		if err != nil {
			fmt.Println("actors.SerializeParams &miner2.WithdrawBalanceParams failed: ", err)
			return
		}

		msg := types.Message{
			From:   key.Address,
			To:     maddr,
			Value:  abi.NewTokenAmount(0),
			Method: builtin.MethodsMultisig.Approve,
			Params: params,
		}
		err = internal.PushSignedMsg(&msg, key.PrivateKey)
		if err != nil {
			fmt.Println(err)
		}
	},
}

func GetMessage(fromAddr, multiAddr, txId string) (error, string) {
	from, err := address.NewFromString(fromAddr)
	if err != nil {
		fmt.Println("decode fromAddr failed:: ", err.Error())
		return err, ""
	}

	balance, err := GetWalletBalance(from)

	if err != nil {
		return err, ""
	}

	maddr, err := address.NewFromString(multiAddr)
	if err != nil {
		fmt.Println("decode multisigAddress failed:: ", err.Error())
		return err, ""
	}

	if maddr.Protocol() != address.Actor && maddr.Protocol() != address.ID {
		fmt.Println("please input a correct multisigAddress")
		return fmt.Errorf("please input a correct multisigAddress"), ""
	}

	txid, err := strconv.Atoi(txId)
	if err != nil {
		fmt.Println("Transaction ID failed: ", err.Error())
		return err, ""
	}

	params, err := actors.SerializeParams(&multisig.TxnIDParams{ID: multisig.TxnID(txid)})
	if err != nil {
		fmt.Println("actors.SerializeParams &miner2.WithdrawBalanceParams failed: ", err)
		return err, ""
	}

	msg := types.Message{
		From:   from,
		To:     maddr,
		Value:  abi.NewTokenAmount(0),
		Method: builtin.MethodsMultisig.Approve,
		Params: params,
	}
	err, str, msgWithGas := internal.GetUnSignedMsg(&msg)

	if err != nil {
		return err, str
	}

	requireFunds := msgWithGas.RequiredFunds()

	if balance.Int.Cmp(requireFunds.Int) < 0 {
		err = fmt.Errorf("地址余额不足:%s, 余额：%v, 需要：%v", fromAddr, balance, requireFunds)
	}

	return err, str

}

func GetCancelMessage(fromAddr, multiAddr, txId string) (error, string) {
	from, err := address.NewFromString(fromAddr)
	if err != nil {
		fmt.Println("decode fromAddr failed:: ", err.Error())
		return err, ""
	}

	balance, err := GetWalletBalance(from)

	if err != nil {
		return err, ""
	}

	maddr, err := address.NewFromString(multiAddr)
	if err != nil {
		fmt.Println("decode multisigAddress failed:: ", err.Error())
		return err, ""
	}

	if maddr.Protocol() != address.Actor && maddr.Protocol() != address.ID {
		fmt.Println("please input a correct multisigAddress")
		return fmt.Errorf("please input a correct multisigAddress"), ""
	}

	txid, err := strconv.Atoi(txId)
	if err != nil {
		fmt.Println("Transaction ID failed: ", err.Error())
		return err, ""
	}

	params, err := actors.SerializeParams(&multisig.TxnIDParams{ID: multisig.TxnID(txid)})
	if err != nil {
		fmt.Println("actors.SerializeParams &miner2.WithdrawBalanceParams failed: ", err)
		return err, ""
	}

	msg := types.Message{
		From:   from,
		To:     maddr,
		Value:  abi.NewTokenAmount(0),
		Method: builtin.MethodsMultisig.Cancel,
		Params: params,
	}
	err, str, msgWithGas := internal.GetUnSignedMsg(&msg)

	if err != nil {
		return err, str
	}

	requireFunds := msgWithGas.RequiredFunds()

	if balance.Int.Cmp(requireFunds.Int) < 0 {
		err = fmt.Errorf("地址余额不足:%s, 余额：%v, 需要：%v", fromAddr, balance, requireFunds)
	}

	return err, str

}

func PushTx(message, signature string) (error, string) {
	return internal.PushMsg(message, signature)
}
