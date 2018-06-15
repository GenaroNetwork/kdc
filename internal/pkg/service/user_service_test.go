package service

import (
	"testing"
	"bytes"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"fmt"
	"encoding/hex"
)
var (
	testmsg     = hexutil.MustDecode("0xce0677bb30baa8cf067c88db9811f4333d131bf8bcf12fe7065d211dce971008")
	testsig     = hexutil.MustDecode("0x90f27b8b488db00b00606796d2987f6a5f59ae62ea05effe84fef5b8b0e549984a691139ad57a3f0b906637673aa2f63d1f55cb1a69199d4009eea23ceaddc9301")
	testpubkey  = hexutil.MustDecode("0x04e32df42865e97135acfb65f3bae71bdc86f4d49150ad6a440b6f15878109880a0a2b2667f7e725ceea70c673093bf67663e0312623c8e091b13cf2c0f11ef652")
	testpubkeyc = hexutil.MustDecode("0x02e32df42865e97135acfb65f3bae71bdc86f4d49150ad6a440b6f15878109880a")
)

func TestRunService(t *testing.T) {
	RunService()
}

func TestEcrecover(t *testing.T) {
	pubkey, err := crypto.Ecrecover(testmsg, testsig)
	if err != nil {
		t.Fatalf("recover error: %s", err)
	}
	if !bytes.Equal(pubkey, testpubkey) {
		t.Errorf("pubkey mismatch: want: %x have: %x", testpubkey, pubkey)
	}
}

func TestFullSignature(t *testing.T) {
	prik, _ := crypto.GenerateKey()
	pubk := &prik.PublicKey
	prikStr := hex.EncodeToString(crypto.FromECDSA(prik))
	pubkStr := hex.EncodeToString(crypto.FromECDSAPub(pubk))

	addr := crypto.PubkeyToAddress(prik.PublicKey).Hex()
	fmt.Printf("private key: %s\n", prikStr)
	fmt.Printf("public  key: %s\n", pubkStr)
	fmt.Printf("addr is    : %s\n", addr)

	msg := "a good example"
	shaMsg := crypto.Keccak256([]byte(msg))
	fmt.Printf("msg    : %s\n", msg)
	fmt.Printf("msg sha: %s\n", shaMsg)

	// signing
	sig, _ := crypto.Sign(shaMsg, prik)
	fmt.Printf("signature : %s\n", sig)

	//recover pub
	recoveredPub, _ := crypto.Ecrecover(shaMsg, sig)
	pubkStr2 := hex.EncodeToString(recoveredPub)
	fmt.Printf("recovered public key : %s\n", pubkStr2)
	//recover pub2
	recoveredPub2, _ := crypto.SigToPub(shaMsg, sig)
	pubkStr3 := hex.EncodeToString(crypto.FromECDSAPub(recoveredPub2))
	fmt.Printf("recovered public key2: %s\n", pubkStr3)

	// validate
	finalAddr := crypto.PubkeyToAddress(*recoveredPub2).Hex()
	fmt.Printf("calculated addr is %s\n", finalAddr)
}

/**
	{
		"jsonrpc": "2.0",
		"method": "subtract",
		"id": 1,
		"params": {
			"fileId": "fhdusihfdisuhdihui",
			"data": "0x00dB21164B6510a4A0c6BC7C48178e31Cd8B6145",
			"amount": "0x91a",
			"signature": "1f02df499f15c8757754c11251a6e5238296f56b17f7229202fce6ccd7289e224c49c32eaf77d5905e2b4d8a8a5ddcc215c51ce45c207ef0f038328200578d1bee"
		}
	}
private key: 7b27e6bc6bf67b99b26f2e00c6acc1102e187a5ff469735f1f65805354809598
public  key: 04dad3c1205c456ca05996285725177c700b901a85f1cc2d02aa0729de49754d2e8db26834f51a3d72cd3e03cc28a6128dc8d35445764fec423ded742d4c2865d2
addr is    : 0x3056D12c1Df12B5299E2eC559eE60a781661Cc48
 */
func TestSignReq(t *testing.T) {
	prik, _ := crypto.GenerateKey()
	pubk := &prik.PublicKey
	prikStr := hex.EncodeToString(crypto.FromECDSA(prik))
	pubkStr := hex.EncodeToString(crypto.FromECDSAPub(pubk))

	addr := crypto.PubkeyToAddress(prik.PublicKey).Hex()
	fmt.Printf("private key: %s\n", prikStr)
	fmt.Printf("public  key: %s\n", pubkStr)
	fmt.Printf("addr is    : %s\n", addr)

	a := "2.0" + "subtract" + "1" + "fhdusihfdisuhdihui" + addr + "0x91a"
	shaMsg := crypto.Keccak256([]byte(a))
	sig, _ := crypto.Sign(shaMsg, prik)
	sigStr := hex.EncodeToString(sig)
	fmt.Printf("signature : %s\n", sigStr)
}
func TestSignReadReq(t *testing.T) {
	prik, _ := crypto.GenerateKey()
	pubk := &prik.PublicKey
	prikStr := hex.EncodeToString(crypto.FromECDSA(prik))
	pubkStr := hex.EncodeToString(crypto.FromECDSAPub(pubk))

	addr := crypto.PubkeyToAddress(prik.PublicKey).Hex()
	fmt.Printf("private key: %s\n", prikStr)
	fmt.Printf("public  key: %s\n", pubkStr)
	fmt.Printf("addr is    : %s\n", addr)

	a := "2.0" + "read" + "1" + "fhdusihfdisuhdihui" + addr
	shaMsg := crypto.Keccak256([]byte(a))
	sig, _ := crypto.Sign(shaMsg, prik)
	sigStr := hex.EncodeToString(sig)
	fmt.Printf("signature : %s\n", sigStr)
}

func TestSignTerminateReq(t *testing.T) {
	prik, _ := crypto.GenerateKey()
	pubk := &prik.PublicKey
	prikStr := hex.EncodeToString(crypto.FromECDSA(prik))
	pubkStr := hex.EncodeToString(crypto.FromECDSAPub(pubk))

	addr := crypto.PubkeyToAddress(prik.PublicKey).Hex()
	fmt.Printf("private key: %s\n", prikStr)
	fmt.Printf("public  key: %s\n", pubkStr)
	fmt.Printf("addr is    : %s\n", addr) //0x5406183353Ec94eA26d9e9c170a4Bd8f26f2b2DD

	a := "2.0" + "terminate" + "1" + "fhdusihfdisuhdihui"
	shaMsg := crypto.Keccak256([]byte(a))
	sig, _ := crypto.Sign(shaMsg, prik)
	sigStr := hex.EncodeToString(sig)
	fmt.Printf("signature : %s\n", sigStr)
}

