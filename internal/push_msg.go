package internal

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/coffeecloudgit/filecoin-wallet-signing/chain/types"
	"github.com/coffeecloudgit/filecoin-wallet-signing/signer"
	"github.com/filecoin-project/go-state-types/crypto"
	"github.com/filecoin-project/go-state-types/manifest"
	"github.com/ipfs/go-cid"
)

func PushSignedMsg(msg *types.Message, privateKey []byte) error {
	nonce, err := Lapi.MpoolGetNonce(Ctx, msg.From)
	if err != nil {
		fmt.Println("Mpool GetNonce failed: ", err)
		return err
	}
	msg.Nonce = nonce
	msgWithGas, err := Lapi.GasEstimateMessageGas(Ctx, msg, nil, *CurrentTsk)
	if err != nil {
		fmt.Println("GasEstimateMessageGas failed: ", err)
		return err
	}

	blk, err := msgWithGas.ToStorageBlock()
	if err != nil {
		fmt.Println("msg.ToStorageBlock() failed: ", err)
		return err
	}

	sigType := signer.AddressSigType(msg.From)
	signed, err := signer.Sign(sigType, privateKey, blk.Cid().Bytes())
	if err != nil {
		fmt.Println("sign failed: ", err.Error())
		return err
	}
	signedMsg := types.SignedMessage{
		Message:   *msgWithGas,
		Signature: *signed,
	}
	b, _ := json.MarshalIndent(signedMsg, " ", " ")
	fmt.Println("Signed message: ", string(b))

	msgCid, err := Lapi.MpoolPush(Ctx, &signedMsg)
	if err != nil {
		fmt.Println("push message failed: ", err.Error())
		return err
	}

	fmt.Println("message CID:", msgCid.String())
	return nil
}

func bytesToHexStr(byteArray []byte) string {
	hexStr := fmt.Sprintf("%x", byteArray)
	return hexStr
}

func GetUnSignedMsg(msg *types.Message) (error, string, *types.Message) {
	nonce, err := Lapi.MpoolGetNonce(Ctx, msg.From)
	if err != nil {
		fmt.Println("Mpool GetNonce failed: ", err)
		return err, "", nil
	}
	msg.Nonce = nonce
	fmt.Println("CurrentTsk", CurrentTsk)
	fmt.Println("msg", msg)
	msgWithGas, err := Lapi.GasEstimateMessageGas(Ctx, msg, nil, *CurrentTsk)
	if err != nil {
		fmt.Println("GasEstimateMessageGas failed: ", err)
		return err, "", nil
	}

	//blk, err := msgWithGas.ToStorageBlock()
	//if err != nil {
	//	fmt.Println("msg.ToStorageBlock() failed: ", err)
	//	return err, ""
	//}
	serialize, err := msgWithGas.Serialize()

	if err != nil {
		fmt.Println("serialize() failed: ", err)
		return err, "", nil
	}
	return nil, bytesToHexStr(serialize), msgWithGas
	//sigType := signer.AddressSigType(msg.From)
	//signed, err := signer.Sign(sigType, privateKey, blk.Cid().Bytes())
	//if err != nil {
	//	fmt.Println("sign failed: ", err.Error())
	//	return err
	//}
	//signedMsg := types.SignedMessage{
	//	Message:   *msgWithGas,
	//	Signature: *signed,
	//}
	//b, _ := json.MarshalIndent(signedMsg, " ", " ")
	//fmt.Println("Signed message: ", string(b))
	//
	//msgCid, err := Lapi.MpoolPush(Ctx, &signedMsg)
	//if err != nil {
	//	fmt.Println("push message failed: ", err.Error())
	//	return err
	//}
	//
	//fmt.Println("message CID:", msgCid.String())
	//return nil
}

func PushMsg(message, signature string) (error, string) {
	msgBytes, err := hex.DecodeString(message)
	if err != nil {
		fmt.Println("Error:", err)
		return err, ""
	}

	// Base64解码
	signatureBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		fmt.Println("Decode error:", err)
		return err, ""
	}

	//signatureBytes, err := hex.DecodeString(signature)
	//if err != nil {
	//	fmt.Println("Error:", err)
	//	return err, ""
	//}

	decodeMessage, err := types.DecodeMessage(msgBytes)
	if err != nil {
		return err, ""
	}

	//blk, err := decodeMessage.ToStorageBlock()
	//if err != nil {
	//	fmt.Println("msg.ToStorageBlock() failed: ", err)
	//	return err
	//}
	//
	sigType := signer.AddressSigType(decodeMessage.From)
	//signed, err := signer.Sign(sigType, privateKey, blk.Cid().Bytes())
	//if err != nil {
	//	fmt.Println("sign failed: ", err.Error())
	//	return err
	//}
	signed := &crypto.Signature{
		Type: sigType,
		Data: signatureBytes,
	}
	signedMsg := types.SignedMessage{
		Message:   *decodeMessage,
		Signature: *signed,
	}
	b, _ := json.MarshalIndent(signedMsg, " ", " ")
	fmt.Println("Signed message: ", string(b))

	msgCid, err := Lapi.MpoolPush(Ctx, &signedMsg)
	if err != nil {
		fmt.Println("push message failed: ", err.Error())
		return err, ""
	}

	fmt.Println("message CID:", msgCid.String())

	return nil, msgCid.String()
}

func StateActorManifestMultisigKeyCID() (cid.Cid, error) {
	result, _ := Lapi.StateActorCodeCIDs(Ctx, NetworkVersion)
	fmt.Println(result)
	return result[manifest.MultisigKey], nil
}
