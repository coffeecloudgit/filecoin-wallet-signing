package msig

import (
	"fmt"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	multisig0 "github.com/filecoin-project/specs-actors/actors/builtin/multisig"
	multisig2 "github.com/filecoin-project/specs-actors/v2/actors/builtin/multisig"
	"github.com/filecoin-project/specs-actors/v8/actors/builtin"
	"github.com/filecoin-project/specs-actors/v8/actors/builtin/multisig"
	"github.com/spf13/cobra"

	"github.com/coffeecloudgit/filecoin-wallet-signing/chain/actors"
	"github.com/coffeecloudgit/filecoin-wallet-signing/chain/types"
	"github.com/coffeecloudgit/filecoin-wallet-signing/internal"
	"github.com/coffeecloudgit/filecoin-wallet-signing/pkg"
)

// proposeCmd represents the msigpropose command
var proposeCmd = &cobra.Command{
	Use:   "propose  multisigAddr toAddr amount",
	Short: "make a proposal",
	Run: func(cmd *cobra.Command, args []string) {
		propose(cmd, args)
	},
}

func propose(ccmd *cobra.Command, args []string) {
	if len(args) < 3 {
		_ = ccmd.Help()
		return
	}

	mtsaddr, err := address.NewFromString(args[0])
	if err != nil {
		fmt.Println("decode multisigAddress failed:", err.Error())
		return
	}

	if mtsaddr.Protocol() != address.Actor && mtsaddr.Protocol() != address.ID {
		fmt.Println("please input a correct multisigAddress")
		return
	}

	acceptAddr, err := address.NewFromString(args[1])
	if err != nil {
		fmt.Println("decode miner address failed:", err.Error())
		return
	}

	sfil, err := types.ParseFIL(args[2])
	if err != nil {
		fmt.Println("The withdrawal amount is wrong or the format is wrong:", err.Error())
		return
	}

	proposeParams, err := actors.SerializeParams(&multisig.ProposeParams{
		To:     acceptAddr,
		Method: builtin.MethodSend,
		Value:  abi.TokenAmount(sfil),
		Params: []byte{},
	})
	if err != nil {
		fmt.Println("actors.SerializeParams &miner2.WithdrawBalanceParams failed: ", err)
		return
	}

	key, err := pkg.ReadPrivteKey()
	if err != nil {
		fmt.Println("decode private key failed: ", err)
		return
	}

	msg := types.Message{
		From:   key.Address,
		To:     mtsaddr,
		Value:  abi.NewTokenAmount(0),
		Method: builtin.MethodsMultisig.Propose,
		Params: proposeParams,
	}

	err = internal.PushSignedMsg(&msg, key.PrivateKey)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("send from %v to %v amount %v \n", mtsaddr.String(), acceptAddr.String(), pkg.ToFloat64(abi.TokenAmount(sfil)))
}

func ProposeTransferMessage(fromAddr, mts, accept, fil string) (error, string) {
	from, err := address.NewFromString(fromAddr)
	if err != nil {
		fmt.Println("发送地址错误: ", err.Error())
		return fmt.Errorf("发送地址错误: %s", err.Error()), ""
	}
	mtsaddr, err := address.NewFromString(mts)
	if err != nil {
		fmt.Println("多签地址错误:", err.Error())
		return fmt.Errorf("多签地址错误: %s", err.Error()), ""
	}

	if mtsaddr.Protocol() != address.Actor && mtsaddr.Protocol() != address.ID {
		fmt.Println("please input a correct multisigAddress")
		return fmt.Errorf("多签地址错误"), ""
	}

	acceptAddr, err := address.NewFromString(accept)
	if err != nil {
		fmt.Println("收币地址错误:", err.Error())
		return fmt.Errorf("收币地址错误: %s", err.Error()), ""
	}

	sfil, err := types.ParseFIL(fil)
	if err != nil {
		fmt.Println("The withdrawal amount is wrong or the format is wrong:", err.Error())
		return fmt.Errorf("转账金额错误: %s", err.Error()), ""
	}

	proposeParams, err := actors.SerializeParams(&multisig.ProposeParams{
		To:     acceptAddr,
		Method: builtin.MethodSend,
		Value:  abi.TokenAmount(sfil),
		Params: []byte{},
	})
	if err != nil {
		fmt.Println("actors.SerializeParams &miner2.WithdrawBalanceParams failed: ", err)
		return err, ""
	}

	//key, err := pkg.ReadPrivteKey()
	//if err != nil {
	//	fmt.Println("decode private key failed: ", err)
	//	return
	//}

	msg := types.Message{
		From:   from,
		To:     mtsaddr,
		Value:  abi.NewTokenAmount(0),
		Method: builtin.MethodsMultisig.Propose,
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
	fmt.Printf("send from %v to %v amount %v \n", mtsaddr.String(), acceptAddr.String(), pkg.ToFloat64(abi.TokenAmount(sfil)))

	return err, str
}

func ProposeAddSignerMessage(fromAddr, mts, editeAddr string, inc bool) (error, string) {
	from, err := address.NewFromString(fromAddr)
	if err != nil {
		fmt.Println("发送地址错误: ", err.Error())
		return fmt.Errorf("发送地址错误: %s", err.Error()), ""
	}
	mtsaddr, err := address.NewFromString(mts)
	if err != nil {
		fmt.Println("多签地址错误:", err.Error())
		return fmt.Errorf("多签地址错误: %s", err.Error()), ""
	}

	if mtsaddr.Protocol() != address.Actor && mtsaddr.Protocol() != address.ID {
		fmt.Println("please input a correct multisigAddress")
		return fmt.Errorf("多签地址错误"), ""
	}

	editeAddress, err := address.NewFromString(editeAddr)
	if err != nil {
		fmt.Println("收币地址错误:", err.Error())
		return fmt.Errorf("收币地址错误: %s", err.Error()), ""
	}

	enc, actErr := actors.SerializeParams(&multisig2.AddSignerParams{
		Signer:   editeAddress,
		Increase: inc,
	})

	if actErr != nil {
		return actErr, ""
	}
	amt := abi.NewTokenAmount(0)
	enc2, actErr2 := actors.SerializeParams(&multisig0.ProposeParams{
		To:     mtsaddr,
		Value:  amt,
		Method: builtin.MethodsMultisig.AddSigner,
		Params: enc,
	})
	if actErr2 != nil {
		return fmt.Errorf("failed to serialize parameters: %w", actErr2), ""
	}

	msg := types.Message{
		From:   from,
		To:     mtsaddr,
		Value:  amt,
		Method: builtin.MethodsMultisig.Propose,
		Params: enc2,
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
	fmt.Printf("send from %v to %v amount %v \n", mtsaddr.String(), editeAddress.String())

	return err, str
}

func ProposeRemoveSignerMessage(fromAddr, mts, editeAddr string, dec bool) (error, string) {
	from, err := address.NewFromString(fromAddr)
	if err != nil {
		fmt.Println("发送地址错误: ", err.Error())
		return fmt.Errorf("发送地址错误: %s", err.Error()), ""
	}
	mtsaddr, err := address.NewFromString(mts)
	if err != nil {
		fmt.Println("多签地址错误:", err.Error())
		return fmt.Errorf("多签地址错误: %s", err.Error()), ""
	}

	if mtsaddr.Protocol() != address.Actor && mtsaddr.Protocol() != address.ID {
		fmt.Println("please input a correct multisigAddress")
		return fmt.Errorf("多签地址错误"), ""
	}

	editeAddress, err := address.NewFromString(editeAddr)
	if err != nil {
		fmt.Println("收币地址错误:", err.Error())
		return fmt.Errorf("收币地址错误: %s", err.Error()), ""
	}

	enc, actErr := actors.SerializeParams(&multisig2.RemoveSignerParams{
		Signer:   editeAddress,
		Decrease: dec,
	})

	if actErr != nil {
		return actErr, ""
	}
	amt := abi.NewTokenAmount(0)
	enc2, actErr2 := actors.SerializeParams(&multisig0.ProposeParams{
		To:     mtsaddr,
		Value:  amt,
		Method: builtin.MethodsMultisig.RemoveSigner,
		Params: enc,
	})
	if actErr2 != nil {
		return fmt.Errorf("failed to serialize parameters: %w", actErr2), ""
	}

	msg := types.Message{
		From:   from,
		To:     mtsaddr,
		Value:  amt,
		Method: builtin.MethodsMultisig.Propose,
		Params: enc2,
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
	fmt.Printf("send from %v to %v amount %v \n", mtsaddr.String(), editeAddress.String())

	return err, str
}

func ProposeChangeThresholdMessage(fromAddr, mts string, threshold uint64) (error, string) {
	from, err := address.NewFromString(fromAddr)
	if err != nil {
		fmt.Println("发送地址错误: ", err.Error())
		return fmt.Errorf("发送地址错误: %s", err.Error()), ""
	}
	mtsaddr, err := address.NewFromString(mts)
	if err != nil {
		fmt.Println("多签地址错误:", err.Error())
		return fmt.Errorf("多签地址错误: %s", err.Error()), ""
	}

	if mtsaddr.Protocol() != address.Actor && mtsaddr.Protocol() != address.ID {
		fmt.Println("please input a correct multisigAddress")
		return fmt.Errorf("多签地址错误"), ""
	}

	enc, actErr := actors.SerializeParams(&multisig2.ChangeNumApprovalsThresholdParams{
		NewThreshold: threshold,
	})

	if actErr != nil {
		return actErr, ""
	}
	amt := abi.NewTokenAmount(0)
	enc2, actErr2 := actors.SerializeParams(&multisig0.ProposeParams{
		To:     mtsaddr,
		Value:  amt,
		Method: builtin.MethodsMultisig.ChangeNumApprovalsThreshold,
		Params: enc,
	})
	if actErr2 != nil {
		return fmt.Errorf("failed to serialize parameters: %w", actErr2), ""
	}

	msg := types.Message{
		From:   from,
		To:     mtsaddr,
		Value:  amt,
		Method: builtin.MethodsMultisig.Propose,
		Params: enc2,
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
	fmt.Printf("change threshold from %v to %v \n", mtsaddr.String(), threshold)

	return err, str
}
