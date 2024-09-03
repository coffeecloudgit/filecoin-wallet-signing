package msig

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/specs-actors/v8/actors/builtin/multisig"
	"github.com/spf13/cobra"

	"github.com/coffeecloudgit/filecoin-wallet-signing/chain/actors"
	"github.com/coffeecloudgit/filecoin-wallet-signing/chain/types"
	"github.com/coffeecloudgit/filecoin-wallet-signing/internal"
	"github.com/coffeecloudgit/filecoin-wallet-signing/pkg"
	init14 "github.com/filecoin-project/go-state-types/builtin/v14/init"
)

// proposeCmd represents the msigpropose command
var createCmd = &cobra.Command{
	Use:   "create a new  multisig Address",
	Short: "make a new  multisig",
	Run: func(cmd *cobra.Command, args []string) {
		create(cmd, args)
	},
}

func create(ccmd *cobra.Command, args []string) {
	if len(args) < 3 {
		_ = ccmd.Help()
		return
	}

	from, err := address.NewFromString(args[0])
	if err != nil {
		fmt.Println("decode from address failed:", err.Error())
		return
	}

	if from.Protocol() != address.Actor && from.Protocol() != address.ID {
		fmt.Println("please input a correct address")
		return
	}

	if from == address.Undef {
		fmt.Println("must provide source address")
		return
	}

	addresses := strings.Split(args[1], ",")
	var signers []address.Address
	for _, a := range addresses {
		addr, err := address.NewFromString(a)
		if err != nil {
			fmt.Println("please input correct address")
			return
		}
		signers = append(signers, addr)
	}

	//acceptAddr, err := address.NewFromString(args[1])
	//if err != nil {
	//	fmt.Println("decode miner address failed:", err.Error())
	//	return
	//}

	threshold, err := strconv.ParseUint(args[2], 10, 64)
	if err != nil {
		fmt.Println("The Threshold is wrong or the format is wrong:", err.Error())
		return
	}

	lenAddrs := uint64(len(signers))

	if lenAddrs < threshold {
		fmt.Println("cannot require signing of more addresses than provided for multisig")
		return
	}

	if threshold <= 0 {
		threshold = lenAddrs
	}
	//ud, err := mstate.UnlockDuration()
	d := abi.ChainEpoch(0)
	// Set up constructor parameters for multisig
	sigParams := &multisig.ConstructorParams{
		Signers:               signers,
		NumApprovalsThreshold: threshold,
		UnlockDuration:        d,
		StartEpoch:            abi.ChainEpoch(0),
	}

	enc, actErr := actors.SerializeParams(sigParams)
	if actErr != nil {
		fmt.Println(actErr)
		return
	}

	//code, ok := actors.GetActorCodeID(actorstypes.Version14, manifest.MultisigKey)
	//if !ok {
	//	fmt.Println("failed to get multisig code ID")
	//	return
	//}
	code, err := internal.StateActorManifestMultisigKeyCID()
	if err != nil {
		fmt.Println("failed to get multisig code ID")
		return
	}
	// new actors are created by invoking 'exec' on the init actor with the constructor params
	execParams := &init14.ExecParams{
		CodeCID:           code,
		ConstructorParams: enc,
	}

	enc, actErr = actors.SerializeParams(execParams)
	if actErr != nil {
		fmt.Println(actErr)
		return
	}

	msg := &types.Message{
		To:     actors.InitActorAddr,
		From:   from,
		Method: actors.MethodsInit.Exec,
		Params: enc,
		Value:  abi.NewTokenAmount(0),
	}

	key, err := pkg.ReadPrivteKey()
	if err != nil {
		fmt.Println("decode private key failed: ", err)
		return
	}

	err = internal.PushSignedMsg(msg, key.PrivateKey)
	if err != nil {
		fmt.Println(err)
		return
	}

	///fmt.Printf("send from %v to %v amount %v \n", mtsaddr.String(), acceptAddr.String(), pkg.ToFloat64(abi.TokenAmount(sfil)))
}

func CreateMessage(fromAddr, stringAddresses, stringThreshold string) (error, string) {
	from, err := address.NewFromString(fromAddr)
	if err != nil {
		fmt.Println("发送地址错误: ", err.Error())
		return fmt.Errorf("发送地址错误: %s", err.Error()), ""
	}

	if from == address.Undef {
		fmt.Println("must provide source address")
		return fmt.Errorf("发送地址不能为空"), ""
	}

	balance, err := GetWalletBalance(from)

	if err != nil {
		return err, ""
	}

	addresses := strings.Split(stringAddresses, ",")
	var signers []address.Address
	for _, a := range addresses {
		addr, err := address.NewFromString(a)
		if err != nil {
			fmt.Println("please input correct address")
			return fmt.Errorf("签名地址格式错误"), ""
		}
		signers = append(signers, addr)
	}

	threshold, err := strconv.ParseUint(stringThreshold, 10, 64)
	if err != nil {
		fmt.Println("The Threshold is wrong or the format is wrong:", err.Error())
		return fmt.Errorf("需要签名数量"), ""
	}

	lenAddrs := uint64(len(signers))

	if lenAddrs < threshold {
		fmt.Println("cannot require signing of more addresses than provided for multisig")
		return fmt.Errorf("需要的签名数量不能大于提供的签名地址总数"), ""
	}

	if threshold <= 0 {
		threshold = lenAddrs
	}
	//ud, err := mstate.UnlockDuration()
	d := abi.ChainEpoch(0)
	// Set up constructor parameters for multisig
	sigParams := &multisig.ConstructorParams{
		Signers:               signers,
		NumApprovalsThreshold: threshold,
		UnlockDuration:        d,
		StartEpoch:            abi.ChainEpoch(0),
	}

	enc, actErr := actors.SerializeParams(sigParams)
	if actErr != nil {
		fmt.Println(actErr)
		return actErr, ""
	}

	//code, ok := actors.GetActorCodeID(actorstypes.Version14, manifest.MultisigKey)
	//if !ok {
	//	fmt.Println("failed to get multisig code ID")
	//	return
	//}
	code, err := internal.StateActorManifestMultisigKeyCID()
	if err != nil {
		fmt.Println("failed to get multisig code ID")
		return fmt.Errorf("failed to get multisig code ID"), ""
	}

	fmt.Println(code)
	// new actors are created by invoking 'exec' on the init actor with the constructor params
	execParams := &init14.ExecParams{
		CodeCID:           code,
		ConstructorParams: enc,
	}

	enc, actErr = actors.SerializeParams(execParams)
	if actErr != nil {
		fmt.Println(actErr)
		return actErr, ""
	}

	msg := &types.Message{
		To:     actors.InitActorAddr,
		From:   from,
		Method: actors.MethodsInit.Exec,
		Params: enc,
		Value:  abi.NewTokenAmount(0),
	}

	err, str, msgWithGas := internal.GetUnSignedMsg(msg)

	if err != nil {
		return err, str
	}

	requireFunds := msgWithGas.RequiredFunds()

	if balance.Int.Cmp(requireFunds.Int) < 0 {
		err = fmt.Errorf("地址余额不足:%s, 余额：%v, 需要：%v", fromAddr, balance, requireFunds)
	}

	return err, str
}
