package msig

import (
	"fmt"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	builtintypes "github.com/filecoin-project/go-state-types/builtin"
	"github.com/filecoin-project/specs-actors/v8/actors/builtin/miner"
	"github.com/filecoin-project/specs-actors/v8/actors/builtin/multisig"
	"github.com/spf13/cobra"

	"github.com/coffeecloudgit/filecoin-wallet-signing/chain/actors"
	"github.com/coffeecloudgit/filecoin-wallet-signing/chain/types"
	"github.com/coffeecloudgit/filecoin-wallet-signing/internal"
	"github.com/coffeecloudgit/filecoin-wallet-signing/pkg"
)

// proposeCmd represents the msigpropose command
var proposeWhithdrawCmd = &cobra.Command{
	Use:   "withdraw <multisigAddress> <minerAddress> <amount> ",
	Short: "propose withdraw from miner ",
	Run: func(cmd *cobra.Command, args []string) {
		proposeWithdraw(cmd, args)
	},
}

func proposeWithdraw(cmd *cobra.Command, args []string) {
	if len(args) < 3 {
		_ = cmd.Help()
		return
	}

	mtsaddr, err := address.NewFromString(args[0])
	if err != nil {
		fmt.Println("decode address failed:", err.Error())
		return
	}

	if mtsaddr.Protocol() != address.Actor && mtsaddr.Protocol() != address.ID {
		fmt.Println("please input a correct multisigAddress")
		return
	}

	mnersaddr, err := address.NewFromString(args[1])
	if err != nil {
		fmt.Println("decode miner address failed:", err.Error())
		return
	}

	if mnersaddr.Protocol() != address.Actor && mnersaddr.Protocol() != address.ID {
		fmt.Println("please input a correct miner address")
		return
	}

	wdfil, err := types.ParseFIL(args[2])
	if err != nil {
		fmt.Println("The withdrawal amount is wrong or the format is wrong: ", err.Error())
		return
	}

	key, err := pkg.ReadPrivteKey()
	if err != nil {
		fmt.Println("decode private key failed: ", err)
		return
	}

	withdrawBalanceParams, err := actors.SerializeParams(&miner.WithdrawBalanceParams{
		AmountRequested: abi.TokenAmount(wdfil), // Default to attempting to withdraw all the extra funds in the miner actor
	})
	if err != nil {
		fmt.Println("actors.SerializeParams &miner2.WithdrawBalanceParams failed: ", err)
		return
	}

	proposeParams, err := actors.SerializeParams(&multisig.ProposeParams{
		To:     mnersaddr,
		Method: builtintypes.MethodsMiner.WithdrawBalance,
		Value:  abi.NewTokenAmount(0),
		Params: withdrawBalanceParams,
	})
	if err != nil {
		fmt.Println("actors.SerializeParams &miner2.WithdrawBalanceParams failed: ", err)
		return
	}

	msg := types.Message{
		From:   key.Address,
		To:     mtsaddr,
		Value:  abi.NewTokenAmount(0),
		Method: builtintypes.MethodsMultisig.Propose,
		Params: proposeParams,
	}

	err = internal.PushSignedMsg(&msg, key.PrivateKey)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("withdraw %v FIL from %v \n", pkg.ToFloat64(abi.TokenAmount(wdfil)), mtsaddr.String())

}

func ProposeWithdrawMessage(fromAddr, mts, minerAddr, fil string) (error, string) {
	from, err := address.NewFromString(fromAddr)
	if err != nil {
		fmt.Println("decode fromAddr failed:: ", err.Error())
		return fmt.Errorf("发送地址错误: %s", err.Error()), ""
	}
	mtsaddr, err := address.NewFromString(mts)
	if err != nil {
		fmt.Println("decode multisigAddress failed:", err.Error())
		return fmt.Errorf("多签地址错误: %s", err.Error()), ""
	}

	if mtsaddr.Protocol() != address.Actor && mtsaddr.Protocol() != address.ID {
		fmt.Println("please input a correct multisigAddress")
		return fmt.Errorf("多签地址错误"), ""
	}

	mnersaddr, err := address.NewFromString(minerAddr)
	if err != nil {
		fmt.Println("decode miner address failed:", err.Error())
		return fmt.Errorf("miner地址错误: %s", err.Error()), ""
	}

	if mnersaddr.Protocol() != address.Actor && mnersaddr.Protocol() != address.ID {
		fmt.Println("please input a correct miner address")
		return fmt.Errorf("miner地址错误"), ""
	}

	sfil, err := types.ParseFIL(fil)
	if err != nil {
		fmt.Println("The withdrawal amount is wrong or the format is wrong:", err.Error())
		return err, ""
	}
	withdrawBalanceParams, err := actors.SerializeParams(&miner.WithdrawBalanceParams{
		AmountRequested: abi.TokenAmount(sfil), // Default to attempting to withdraw all the extra funds in the miner actor
	})

	if err != nil {
		fmt.Println("actors.SerializeParams &miner2.WithdrawBalanceParams failed: ", err)
		return err, ""
	}
	proposeParams, err := actors.SerializeParams(&multisig.ProposeParams{
		To:     mnersaddr,
		Method: builtintypes.MethodsMiner.WithdrawBalance,
		Value:  abi.NewTokenAmount(0),
		Params: withdrawBalanceParams,
	})

	if err != nil {
		fmt.Println("actors.SerializeParams &miner2.WithdrawBalanceParams failed: ", err)
		return err, ""
	}

	msg := types.Message{
		From:   from,
		To:     mtsaddr,
		Value:  abi.NewTokenAmount(0),
		Method: builtintypes.MethodsMultisig.Propose,
		Params: proposeParams,
	}

	err, str, msgWithGas := internal.GetUnSignedMsg(&msg)

	if err != nil {
		return err, str
	}

	requireFunds := msgWithGas.RequiredFunds()
	balance, err := GetWalletBalance(from)

	if err != nil {
		return err, ""
	}

	if balance.Int.Cmp(requireFunds.Int) < 0 {
		err = fmt.Errorf("地址余额不足:%s, 余额：%v, 需要：%v", fromAddr, balance, requireFunds)
	}

	return err, str
}
