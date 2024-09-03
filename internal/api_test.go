package internal

import (
	"fmt"
	"github.com/filecoin-project/go-address"
	abinetwork "github.com/filecoin-project/go-state-types/network"
	"testing"
)

func TestChangeTks(t *testing.T) {
	ChangeTks()

	fmt.Println(CurrentTsk)

	result, _ := Lapi.StateActorCodeCIDs(Ctx, abinetwork.Version23)

	fmt.Println(result)
}

func TestStateActorCodeCIDs(t *testing.T) {
	ChangeTks()

	//fmt.Println(CurrentTsk)

	result, _ := Lapi.StateActorCodeCIDs(Ctx, NetworkVersion)

	fmt.Println(result)
}

func TestStateNetworkVersion(t *testing.T) {
	ChangeTks()

	result, _ := Lapi.StateNetworkVersion(Ctx, *CurrentTsk)

	fmt.Println(result)
}

func TestStateActorManifestCID(t *testing.T) {
	ChangeTks()

	//fmt.Println(CurrentTsk)

	result, _ := Lapi.StateActorManifestCID(Ctx, NetworkVersion)

	fmt.Println(result)
}

func TestStateGetActor(t *testing.T) {
	ChangeTks()
	account := "f14qmuid2b6ne4342m5dk56f4rcr7y5sz4sg5fiwy"
	addr, err := address.NewFromString(account)
	if err != nil {
		return
	}
	//result, _ := Lapi.StateGetActor(Ctx, addr, *CurrentTsk)

	result, _ := Lapi.StateLookupID(Ctx, addr, *CurrentTsk)

	fmt.Println(result)
}
