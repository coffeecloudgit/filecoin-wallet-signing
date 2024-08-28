package internal

import (
	"encoding/hex"
	"fmt"
	"github.com/coffeecloudgit/filecoin-wallet-signing/chain/types"
	"testing"
)

func TestPushMsg(t *testing.T) {
	msg := "8a0058310396a1a3e4ea7a14d49985e661b22401d44fed402d1d0925b243c923589c0fbc7e32cd04e29ed78d15d37d3aaa3fe6da3358310386b454258c589475f7d16f5aac018a79f6c1169d20fc33921dd8b5ce1cac6c348f90a3603624f6aeb91b64518c2e80950144000186a01961a8430009c44200000040"

	msgBytes, err := hex.DecodeString(msg)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(string(msgBytes))

	decodeMessage, err := types.DecodeMessage(msgBytes)

	fmt.Println(decodeMessage)
}
