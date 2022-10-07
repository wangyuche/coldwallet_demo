package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	usdt "github.com/wangyuche/coldwallet_demo/contract"
)

var pk string = "88d7a27b710e48c4ff80433a4d3f55f2036745e5caa713f4ed200615c5bc2219"

func main() {
	client, err := ethclient.Dial("http://localhost:8545")
	if err != nil {
		log.Fatal(err)
	}
	//deploy(client)
	//mint(client)
	//getbalance(client,"0x8afc4C8DA521433965bc033568234603B3dE753C")
	transfer(client,"0x08F3f690166202CC1838B5A46a4Df146e23A95EE")
	getbalance(client,"0x8afc4C8DA521433965bc033568234603B3dE753C")
	getbalance(client,"0x08F3f690166202CC1838B5A46a4Df146e23A95EE")
}

func deploy(client *ethclient.Client) {
	privateKey, err := crypto.HexToECDSA(pk)
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)      // in wei
	auth.GasLimit = uint64(6000000) // in units
	auth.GasPrice = gasPrice

	address, tx, instance, err := usdt.DeployUsdt(auth, client, "Tether USD", "USDT")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(address.Hex())   // 0x54d1D2D79f80163182190d1647c2AFE4421f7C43
	fmt.Println(tx.Hash().Hex()) // 0x014bf862a04db84f1ac4392e9d3dec05f5c25883a375ca64e49d50babd4be3b4
	_ = instance
}

func getbalance(client *ethclient.Client, playeraddr string) {
	privateKey, err := crypto.HexToECDSA(pk)
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)      // in wei
	auth.GasLimit = uint64(6000000) // in units
	auth.GasPrice = gasPrice

	contractaddress := common.HexToAddress("0x54d1D2D79f80163182190d1647c2AFE4421f7C43")
	instance, err := usdt.NewUsdt(contractaddress, client)
	if err != nil {
		log.Fatal(err)
	}
	player := common.HexToAddress(playeraddr)

	b, err := instance.BalanceOf(&bind.CallOpts{}, player)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("balance : %d", b)
}

func mint(client *ethclient.Client) {
	privateKey, err := crypto.HexToECDSA(pk)
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)      // in wei
	auth.GasLimit = uint64(6000000) // in units
	auth.GasPrice = gasPrice

	contractaddress := common.HexToAddress("0x54d1D2D79f80163182190d1647c2AFE4421f7C43")
	instance, err := usdt.NewUsdt(contractaddress, client)
	if err != nil {
		log.Fatal(err)
	}

	amount := big.NewInt(100000000)
	_, err = instance.Mint(auth, amount)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Mint")
}

func transfer(client *ethclient.Client, recipientadd string) {
	privateKey, err := crypto.HexToECDSA(pk)
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)      // in wei
	auth.GasLimit = uint64(6000000) // in units
	auth.GasPrice = gasPrice

	contractaddress := common.HexToAddress("0x54d1D2D79f80163182190d1647c2AFE4421f7C43")
	instance, err := usdt.NewUsdt(contractaddress, client)
	if err != nil {
		log.Fatal(err)
	}

	recipient := common.HexToAddress(recipientadd)
	amount := big.NewInt(100000000)
	_, err = instance.Transfer(auth, recipient, amount)
	if err != nil {
		log.Fatal(err)
	}
}
