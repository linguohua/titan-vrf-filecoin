package test

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"
	"titan-vrf/filrpc"
	"titan-vrf/trand"

	"github.com/filecoin-project/go-address"
)

var (
	chainHeight = int64(3290700)

	filPrivateKey = []byte{152, 112, 35, 145, 21, 31, 99, 206, 204, 113, 33, 99, 241, 180, 157, 194, 91, 224, 34, 186, 137, 10, 136, 38, 133, 32, 109, 255, 59, 81, 45, 26}
	// filecoin bls public key
	filPublicKey = []byte{146, 209, 52, 147, 166, 127, 130, 148, 172, 13, 162, 254, 17, 85, 254, 151, 93, 182, 28, 218, 103, 106, 200, 115, 178, 101, 156, 74, 25, 214, 220, 136, 167, 32, 147, 231, 40, 250, 149, 109, 229, 58, 7, 135, 214, 93, 55, 169}
)

func dumpbytes(vbytes []byte) string {
	var sb strings.Builder
	sb.WriteString("{")
	for _, v := range vbytes {
		sb.WriteString(fmt.Sprintf("%d,", v))
	}
	sb.WriteString("}")
	return sb.String()
}

func TestVRFGenV(t *testing.T) {
	nodeURL := "http://api.node.glif.io/rpc/v1"

	client := filrpc.New(
		filrpc.NodeURLOption(nodeURL),
	)

	tps, err := client.ChainGetTipSetByHeight(chainHeight)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("chain head tipset height:%d", tps.Height())

	ctx := context.Background()
	// privateKey, err := trand.BlsGenPrivateKey()
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// public, err := trand.BlsToPublic(privateKey)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	privateKey := filPrivateKey
	public := filPublicKey

	t.Logf("private key:%s", dumpbytes(privateKey))
	t.Logf("public key:%s", dumpbytes(public))

	var entropy []byte
	var gameRoundInfo = GameRoundInfo{
		GameID:    "abc-efg-hi",
		PlayerIDs: "a,b,c,d",
		RoundID:   "gogogogo1",
		ReplayID:  "bilibili",
	}

	buf := new(bytes.Buffer)
	err = gameRoundInfo.MarshalCBOR(buf)
	if err != nil {
		t.Fatal(err)
	}
	entropy = buf.Bytes()

	vrfout, err := trand.GenerateVRFByTipSet(ctx, trand.DomainSeparationTag_GameBasic, privateKey, tps, entropy)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("proof:%s", dumpbytes(vrfout.Proof))

	addr, err := address.NewBLSAddress(public)
	if err != nil {
		t.Fatal(err)
	}

	err = trand.VerifyVRFByTipSet(ctx, addr, trand.DomainSeparationTag_GameBasic, tps, entropy, vrfout)
	if err != nil {
		t.Fatal(err)
	}
}
