package Chain

import (
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"math"
	"math/big"
	"strconv"
	"strings"
	coin "vgo/core/eth"
	"vgo/core/global"
	"vgo/core/response"
)

var (
	rpcTestUrl = "https://data-seed-prebsc-1-s1.binance.org:8545"
	rpcUrl     = "https://bsc-dataseed.binance.org"
)

// Transfer erc20 代币转账
func Transfer(ctx *gin.Context) {
	url, err1 := getUrl(ctx)
	if err1 != nil {
		response.Fail(ctx, err1.Error(), nil, nil)
		return
	}
	wallet := coin.NewEthChain()
	_, err := wallet.InitRemote(coin.UrlParam{RpcUrl: url})
	if err != nil {
		response.Fail(ctx, "请求失败！！", nil, nil)
		return
	}
	fromAddress := ctx.PostForm("from_address")
	if fromAddress == "" {
		response.Fail(ctx, "缺少参数【from_address】-转出地址", nil, nil)
	}
	fromPrivateKey := ctx.PostForm("from_private")
	if fromPrivateKey == "" {
		response.Fail(ctx, "缺少参数【fromPrivateKey】-钱包私钥", nil, nil)
		return
	}
	itemContract := ctx.PostForm("item_contract")
	if itemContract == "" {
		response.Fail(ctx, "缺少参数【itemContract】-代币合约地址", nil, nil)
		return
	}
	toAddress := ctx.PostForm("to_address")
	if toAddress == "" {
		response.Fail(ctx, "缺少参数【to_address】-接受者钱包地址", nil, nil)
		return
	}
	price := ctx.PostForm("price")
	if price == "" {
		response.Fail(ctx, "缺少参数【price】-金额", nil, nil)
		return
	}
	nonce, _ := wallet.Nonce(fromAddress)
	callMethodOpts := &coin.CallMethodOpts{
		Nonce: nonce,
	}
	f, err := strconv.ParseFloat(price, 64)
	if err != nil {
		fmt.Println("转换失败：", err)
		response.Fail(ctx, "传值错误！！", nil, nil)
		return
	}
	amount := big.NewFloat(f)
	tenDecimal := big.NewFloat(math.Pow(10, float64(18)))
	convertAmount, _ := new(big.Float).Mul(tenDecimal, amount).Int(&big.Int{})
	buildTxResult, err := wallet.BuildCallMethodTx(
		fromPrivateKey,
		itemContract,
		coin.Erc20AbiStr,
		"transfer",
		callMethodOpts,
		common.HexToAddress(toAddress), // 目标账户
		convertAmount,
	)
	if err != nil {
		fmt.Printf("构建调用方法 tx 错误: %v\n", err)
		response.Fail(ctx, "构建调用方法 tx 错误！！", nil, nil)
		return
	}
	sendTxResult, err := wallet.SendRawTransaction(buildTxResult.TxHex)
	if err != nil {
		fmt.Printf("发送原始交易错误: %v\n", err)
		response.Fail(ctx, "发送原始交易错误 tx 错误！！", nil, nil)
		return
	}
	fmt.Printf("发送结果: %v\n", sendTxResult)
	response.Success(ctx, "成功", map[string]interface{}{
		"hash": sendTxResult,
	}, nil)
}

// BalanceOther  获取代余额
func BalanceOther(ctx *gin.Context) {
	url, err1 := getUrl(ctx)
	if err1 != nil {
		response.Fail(ctx, err1.Error(), nil, nil)
		return
	}
	wallet := coin.NewEthChain()
	_, err := wallet.InitRemote(coin.UrlParam{RpcUrl: url})
	if err != nil {
		response.Fail(ctx, "请求失败！！", nil, nil)
		return
	}
	itemContract := ctx.PostForm("item_contract")
	if itemContract == "" {
		response.Fail(ctx, "缺少参数【item_contract】-代币合约地址", nil, nil)
		return
	}
	address := ctx.PostForm("address")
	if address == "" {
		response.Fail(ctx, "缺少参数【address】-钱包地址", nil, nil)
		return
	}
	itemBalance, _ := wallet.TokenBalance(itemContract, address)
	a := new(big.Float)
	a.SetString(itemBalance.String())
	itemValue := new(big.Float).Quo(a, big.NewFloat(math.Pow10(18)))
	//fmt.Printf("代币余额11: %v\n", itemValue)
	response.Success(ctx, "成功", map[string]interface{}{
		"balance": itemValue.String(),
	}, nil)
}

// BalanceBnb 获取 BNB 余额
func BalanceBnb(ctx *gin.Context) {
	url, err1 := getUrl(ctx)
	if err1 != nil {
		response.Fail(ctx, err1.Error(), nil, nil)
		return
	}
	wallet := coin.NewEthChain()
	_, err := wallet.InitRemote(coin.UrlParam{RpcUrl: url})
	if err != nil {
		response.Fail(ctx, "请求失败！！", nil, nil)
		return
	}
	address := ctx.PostForm("address")
	if address == "" {
		response.Fail(ctx, "缺少参数【address】-钱包地址", nil, nil)
		return
	}
	balance, _ := wallet.Balance(address)
	fBalance := new(big.Float)
	fBalance.SetString(balance.String())
	ethValue := new(big.Float).Quo(fBalance, big.NewFloat(math.Pow10(18)))
	//fmt.Printf("币安币余额11: %v\n", ethValue)
	response.Success(ctx, "成功", map[string]interface{}{
		"balance": ethValue.String(),
	}, nil)
}

// TransactionQuery 通过hash 获取交易状态  合约地址  转出钱包地址  目标钱包地址
func TransactionQuery(ctx *gin.Context) {
	url, err1 := getUrl(ctx)
	if err1 != nil {
		response.Fail(ctx, err1.Error(), nil, nil)
		return
	}
	hash := ctx.PostForm("hash")
	if hash == "" {
		response.Fail(ctx, "缺少参数【hash】-交易hash", nil, nil)
		return
	}
	wallet := coin.NewEthChain()
	_, err := wallet.InitRemote(coin.UrlParam{RpcUrl: url})
	if err != nil {
		response.Fail(ctx, "请求失败！！", nil, nil)
		return
	}
	receipt, _ := wallet.TransactionReceiptByHash(hash)
	status := receipt.Status
	var arr []string
	var fromAddress []string
	for _, log := range receipt.Logs {
		fromAddress = append(fromAddress, log.Address.Hex())
		//fmt.Println("合约地址:", log.Address.Hex())
		for _, topic := range log.Topics {
			arr = append(arr, topic.Hex())
		}
	}
	fmt.Println(arr)

	chainQueryFromAddress := arr[1]
	chainQueryFromAddress = "0x" + strings.ToLower(chainQueryFromAddress[len(chainQueryFromAddress)-40:])
	//fmt.Println("转出钱包地址:", chainQueryFromAddress)
	chainQueryTo := arr[2]
	chainQueryTo = "0x" + strings.ToLower(chainQueryTo[len(chainQueryTo)-40:])
	//fmt.Println("目标钱包地址:", chainQueryTo)
	//fmt.Println("交易状态:", status)
	response.Success(ctx, "成功", map[string]interface{}{
		"fromContract": fromAddress[0],
		"fromAddress":  chainQueryFromAddress,
		"chainQueryTo": chainQueryTo,
		"status":       status,
	}, nil)
}

func getUrl(ctx *gin.Context) (string, error) {
	appConf := global.App.Config.App
	secret := ctx.PostForm("secret")
	if secret == "" {
		return "", errors.New("缺少参数【secret】-秘钥")
	}
	if secret != appConf.Secret {
		return "", errors.New("秘钥错误")
	}
	var url string
	debug := ctx.PostForm("debug")
	if debug == "" {
		return "", errors.New("缺少参数【debug】-调试开关")
	}
	if debug == "1" {
		url = rpcTestUrl
	} else {
		url = rpcUrl
	}
	return url, nil
}
