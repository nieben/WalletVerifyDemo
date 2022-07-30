package main

import (
	"WalletVerifyDemo/types"
	"bytes"
	"crypto/ecdsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	testBaseUrl      = "http://127.0.0.1:2333"
	defaultHexPrvkey = "19ea515793137c8ef6402767d3226574808a2309903d47ba683c9ea30a92fa38"
	defaultAddress   = "0x44c9e15458514ee416365c4DbCF98E4A2259406e"
	//defaultAddress = "0x0000000000000000000000000000000000000000"  // invalidAddress
)

var (
	testOption int
)

func main() {
	var prvkey *ecdsa.PrivateKey
	var err error
	var address string

	fmt.Println("Please select a test option:")
	fmt.Println("1: use default key")
	fmt.Println("2: use random generated key")

	_, err = fmt.Scan(&testOption)
	if err != nil {
		panic(err)
	}
	if testOption != 1 && testOption != 2 {
		fmt.Println("unknown test option")
		return
	}

	if testOption == 1 { // default key
		prvkey, err = crypto.HexToECDSA(defaultHexPrvkey)
		if err != nil {
			panic(err)
		}
		address = defaultAddress
	} else { // random key
		prvkey, err = crypto.GenerateKey()
		if err != nil {
			panic(err)
		}
		address = crypto.PubkeyToAddress(prvkey.PublicKey).Hex()
	}

	// get random message
	msg, err := getMessage()
	if err != nil {
		panic(err)
	}

	fmt.Println("got message: ", msg)

	sig, err := secp256k1.Sign(msg, crypto.FromECDSA(prvkey))
	if err != nil {
		panic(err)
	}

	fmt.Println("sig message: ", sig)

	// verify
	verified, err := verify(address, msg, sig)
	if err != nil {
		panic(err)
	}

	fmt.Printf("verify result: [%v], address: %s", verified, address)
}

func getMessage() ([]byte, error) {
	req, err := http.NewRequest("GET", testBaseUrl+"/get_message", nil)
	if err != nil {
		return nil, err
	}
	req.Close = true

	body, err := doRequest(req)
	if err != nil {
		return nil, err
	}

	resp := &types.MessageResponse{}
	if err := json.Unmarshal(body, resp); err != nil {
		return nil, err
	}

	b, err := base64.StdEncoding.DecodeString(resp.Message)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func verify(address string, msg, sig []byte) (bool, error) {
	b, _ := json.Marshal(types.VerifyRequest{
		Address:       address,
		Message:       base64.StdEncoding.EncodeToString(msg),
		SignedMessage: base64.StdEncoding.EncodeToString(sig),
	})

	req, err := http.NewRequest("POST", testBaseUrl+"/verify", bytes.NewBuffer(b))
	if err != nil {
		return false, err
	}
	req.Close = true
	req.Header.Add("Content-Type", "application/json")

	body, err := doRequest(req)
	if err != nil {
		return false, err
	}

	resp := &types.VerifyResponse{}
	if err := json.Unmarshal(body, resp); err != nil {
		return false, err
	}

	return resp.Verified, nil
}

func doRequest(req *http.Request) ([]byte, error) {
	client := &http.Client{Timeout: time.Second * time.Duration(10)}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code:%d, body: %s", resp.StatusCode, string(body))
	}

	return body, nil
}
