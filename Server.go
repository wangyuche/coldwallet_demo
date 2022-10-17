package main

import (
	"context"
	"crypto/ecdsa"
	"crypto/tls"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	usdt "github.com/wangyuche/coldwallet_demo/contract"
	"gitlab.edufgs.com/gomod/Mysql"
)

var (
	commitid  = ""
	buildtime = ""
	version   = ""
)

var pk string = "88d7a27b710e48c4ff80433a4d3f55f2036745e5caa713f4ed200615c5bc2219"
var client *ethclient.Client

func Setup() *fiber.App {
	engine := html.New("./view", ".html")
	engine.AddFunc(
		// add unescape function
		"unescape", func(s string) template.HTML {
			return template.HTML(s)
		},
	)
	engine.Reload(true)
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{})
	})
	app.Post("/bindfb", bindfb)
	app.Post("/balance", balance)
	app.Static("/", "./view/")
	return app
}

func main() {
	Mysql.SetWriteConnectionInfo("test", os.Getenv("test"), 5, 5)
	app := Setup()
	log.Print("Listen" + os.Getenv("Port"))
	cer, err := tls.LoadX509KeyPair("ca.crt", "ca.key")
	if err != nil {
		log.Fatal(err)
	}

	config := &tls.Config{Certificates: []tls.Certificate{cer}}
	ln, err := tls.Listen("tcp", os.Getenv("Port"), config)
	if err != nil {
		panic(err)
	}
	client, err = ethclient.Dial("https://cloudflare-eth.com")
	if err != nil {
		log.Fatal(err)
	}
	app.Listener(ln)
	return
	//deploy(client)
	//mint(client)
	//getbalance(client,"0x8afc4C8DA521433965bc033568234603B3dE753C")
	transfer(client, "0x08F3f690166202CC1838B5A46a4Df146e23A95EE")
	getbalance(client, "0x8afc4C8DA521433965bc033568234603B3dE753C")
	getbalance(client, "0x08F3f690166202CC1838B5A46a4Df146e23A95EE")
}

type StateCode struct {
	Code int `json:"Code"`
}

type BindCS struct {
	Fbid string `json:"fbid"`
}

type BindSC struct {
	Address string    `json:"Address"`
	State   StateCode `json:"State"`
}

func bindfb(c *fiber.Ctx) error {
	cs := new(BindCS)
	data := new(BindSC)
	state := new(StateCode)
	if err := c.BodyParser(cs); err != nil {
		log.Print(err)
		state.Code = 9999
		data.State = *state
		return c.JSON(data)
	}
	conn, err := Mysql.GetWriteConnection("test").Begin()
	if err != nil {
		log.Print(err)
		state.Code = 9999
		data.State = *state
		return c.JSON(data)
	}
	defer conn.Commit()
	var _pk, address string
	err = conn.QueryRow("SELECT pk,address FROM user where fbid = ?", cs.Fbid).Scan(&_pk, &address)
	if err != nil {
		if err == sql.ErrNoRows {
			_pk, address := createpk()
			_, err = conn.Exec("INSERT INTO user (fbid,pk,address) VALUES (?,?,?)", cs.Fbid, _pk, address)
			if err != nil {
				log.Print(err)
				state.Code = 9999
				data.State = *state
				return c.JSON(data)
			}
			state.Code = 10000
			data.State = *state
			data.Address = address
			return c.JSON(data)
		}
		log.Print(err)
		state.Code = 9999
		data.State = *state
		return c.JSON(data)
	}
	state.Code = 10000
	data.State = *state
	data.Address = address
	return c.JSON(data)
}

type BalanceCS struct {
	Fbid string `json:"fbid"`
}

type BalanceSC struct {
	Balance int       `json:"Balance"`
	State   StateCode `json:"State"`
}

func balance(c *fiber.Ctx) error {
	cs := new(BalanceCS)
	data := new(BalanceSC)
	state := new(StateCode)
	if err := c.BodyParser(cs); err != nil {
		log.Print(err)
		state.Code = 9999
		data.State = *state
		return c.JSON(data)
	}
	conn, err := Mysql.GetWriteConnection("test").Begin()
	if err != nil {
		log.Print(err)
		state.Code = 9999
		data.State = *state
		return c.JSON(data)
	}
	defer conn.Commit()
	var pk, address string
	err = conn.QueryRow("SELECT pk,address FROM user where fbid = ?", cs.Fbid).Scan(&pk, &address)
	if err != nil {
		log.Print(err)
		state.Code = 9999
		data.State = *state
		return c.JSON(data)
	}
	b := getbalance(client, address)
	data.Balance = b
	state.Code = 10000
	data.State = *state
	return c.JSON(data)
}

func createpk() (string, string) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)
	_pk := hexutil.Encode(privateKeyBytes)[2:]
	fmt.Print(_pk)

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	fmt.Println(hexutil.Encode(publicKeyBytes)[4:]) // 9a7df67f79246283fdc93af76d4f8cdd62c4886e8cd870944e817dd0b97934fdd7719d0810951e03418205868a5c1b40b192451367f28e0088dd75e15de40c05

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	fmt.Print(address) // 0x96216849c49358B10257cb55b28eA603c874b05E
	return _pk, address
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

func getbalance(client *ethclient.Client, playeraddr string) int {
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

	contractaddress := common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7")
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
	return int(b.Int64())
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
