package main

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/ybbus/jsonrpc"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"
)

type JsonRpcRequest struct {
	Id      interface{}   `json:"id"`
	JsonRpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

type JsonRpcResponse struct {
	Id     interface{}  `json:"id"`
	Result *interface{} `json:"result"`
	Error  *Err         `json:"error"`
}

type AddressKeyPair struct {
	Address    string `json:"address"`
	PrivateKey string `json:"privateKey"`
	Encrypted  bool   `json:"encrypted"`
}

type GenerateAddressResponse struct {
	Id     interface{}       `json:"id"`
	Result *[]AddressKeyPair `json:"result"`
	Error  *Err              `json:"error"`
}

type SignTransactionResponse struct {
	Id     interface{} `json:"id"`
	Result *string     `json:"result"`
	Error  *Err        `json:"error"`
}

type MultiSigAddressRes struct {
	RedeemScript    string `json:"redeemScript"`
	MultiSigAddress string `json:"multiSigAddress"`
}

type GenerateMultiAddressResponse struct {
	Id     interface{}         `json:"id"`
	Result *MultiSigAddressRes `json:"result"`
	Error  *Err                `json:"error"`
}

type MultiSignTransactionResponse struct {
	Id     interface{} `json:"id"`
	Result *string     `json:"result"`
	Error  *Err        `json:"error"`
}

type ImportAddressesResponse struct {
	Id     interface{}  `json:"id"`
	Result *interface{} `json:"result"`
	Error  *Err         `json:"error"`
}

type UtxoRes struct {
	Address      string `json:"address"`
	Txid         string `json:"txid"`
	Vout         int    `json:"vout"`
	Amount       string `json:"amount"`
	ScriptPubKey string `json:"scriptPubKey"`
}

type QueryUtxosResponse struct {
	Id     interface{} `json:"id"`
	Result *[]UtxoRes  `json:"result"`
	Error  *Err        `json:"error"`
}

var app *iris.Application

func ReadJsonRpcBody(ctx iris.Context) (interface{}, string, []byte, error) {
	reader := ctx.Request().Body
	//bodyBytes := make([]byte, ctx.Request().ContentLength, ctx.Request().ContentLength)
	//readCount, err := reader.Read(bodyBytes)
	//if err.Error() != "EOF" || int64(readCount) != ctx.Request().ContentLength {
	//	return 0, "", nil, errors.New("read http jsonrpc body error")
	//}
	bodyBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return 0, "", nil, err
	}
	if int64(len(bodyBytes)) != ctx.Request().ContentLength {
		return 0, "", nil, errors.New("read http jsonrpc body error")
	}

	var jsonRpcRequest JsonRpcRequest
	err = json.Unmarshal(bodyBytes, &jsonRpcRequest)
	if err != nil {
		return 0, "", nil, err
	}
	return jsonRpcRequest.Id, jsonRpcRequest.Method, bodyBytes, nil
}

func GenerateAddressController(ctx iris.Context, jsonRpcBody []byte) {
	var req JsonRpcRequest
	_ = json.Unmarshal(jsonRpcBody, &req)

	var res GenerateAddressResponse
	res.Id = req.Id

	if len(req.Params) != 1 {
		res.Error = MakeError(-1, "invalid jsonrpc request params length")
		ctx.JSON(res)
		return
	}

	var count uint32
	typeStr := reflect.TypeOf(req.Params[0]).String()
	if typeStr == "float64" {
		count = uint32(req.Params[0].(float64))
	} else if typeStr == "string" {
		i, err := strconv.Atoi(req.Params[0].(string))
		if err != nil {
			res.Error = MakeError(-1, "invalid jsonrpc request params")
			ctx.JSON(res)
			return
		}
		count = uint32(i)
	} else {
		res.Error = MakeError(-1, "invalid jsonrpc request params")
		ctx.JSON(res)
		return
	}

	if count > 1000 {
		count = 1000
	}

	pairs := make([]AddressKeyPair, 0)
	addresses := make([]string, 0)
	res.Result = &pairs
	for i := uint32(0); i < count; i++ {
		_, privHex, _, addrStr, err := BTCGenerateNewAddress()
		if err != nil {
			res.Error = MakeError(-1, err.Error())
			ctx.JSON(res)
			return
		}
		cryptedBytes := AesEncrypt(privHex, SecurityPassStr)
		if len(cryptedBytes) == 0 {
			res.Error = MakeError(-1, "AesEncrypt fail")
			ctx.JSON(res)
			return
		}
		cryptedHex := hex.EncodeToString(cryptedBytes)
		pairs = append(pairs, AddressKeyPair{Address: addrStr, PrivateKey: cryptedHex, Encrypted: true})
		addresses = append(addresses, addrStr)
	}

	err := GlobalDBMgr.TblAddressMgr.AddNewAddresses(addresses)
	if err != nil {
		res.Error = MakeError(-1, err.Error())
		ctx.JSON(res)
		return
	}

	ctx.JSON(res)
	return
}

func SignTransactionController(ctx iris.Context, jsonRpcBody []byte) {
	var req JsonRpcRequest
	_ = json.Unmarshal(jsonRpcBody, &req)

	var res SignTransactionResponse
	res.Id = req.Id

	if len(req.Params) != 3 {
		res.Error = MakeError(-1, "invalid jsonrpc request params length")
		ctx.JSON(res)
		return
	}

	rawTrxStr, privKeyEncryptHexStr, utxosStr := "", "", ""
	typeStr := reflect.TypeOf(req.Params[0]).String()
	if typeStr == "string" {
		rawTrxStr = req.Params[0].(string)
	} else {
		res.Error = MakeError(-1, "invalid jsonrpc request params[0]")
		ctx.JSON(res)
		return
	}

	typeStr = reflect.TypeOf(req.Params[1]).String()
	if typeStr == "string" {
		privKeyEncryptHexStr = req.Params[1].(string)
	} else {
		res.Error = MakeError(-1, "invalid jsonrpc request params[1]")
		ctx.JSON(res)
		return
	}

	typeStr = reflect.TypeOf(req.Params[2]).String()
	if typeStr == "string" {
		utxosStr = req.Params[2].(string)
	} else {
		res.Error = MakeError(-1, "invalid jsonrpc request params[2]")
		ctx.JSON(res)
		return
	}

	privKeyEncryptBytes, err := hex.DecodeString(privKeyEncryptHexStr)
	if err != nil {
		res.Error = MakeError(-1, "invalid jsonrpc request params[1], privKeyEncryptHexStr not hex format string")
		ctx.JSON(res)
		return
	}

	privKeyHexStr := string(AesDecrypt(privKeyEncryptBytes, []byte(SecurityPassStr)))
	if len(privKeyHexStr) == 0 {
		res.Error = MakeError(-1, "AesDecrypt fail")
		ctx.JSON(res)
		return
	}

	var utxos UTXOsDetail
	//err = json.Unmarshal([]byte(utxosStr), &utxos)
	//if err != nil {
	//	res.Error = MakeError(-1, "invalid jsonrpc request params[2], Unmarshal fail")
	//	ctx.JSON(res)
	//	return
	//}
	_ = utxosStr

	trxSigStr, err := BTCSignRawTransaction(rawTrxStr, privKeyHexStr, utxos)
	if err != nil {
		res.Error = MakeError(-1, "sign raw transaction fail")
		ctx.JSON(res)
		return
	}

	res.Result = &trxSigStr
	ctx.JSON(res)
	return
}

func GenerateMultiAddressController(ctx iris.Context, jsonRpcBody []byte) {
	var req JsonRpcRequest
	_ = json.Unmarshal(jsonRpcBody, &req)

	var res GenerateMultiAddressResponse
	res.Id = req.Id

	if len(req.Params) != 2 {
		res.Error = MakeError(-1, "invalid jsonrpc request params length")
		ctx.JSON(res)
		return
	}

	var need uint32
	typeStr := reflect.TypeOf(req.Params[0]).String()
	if typeStr == "float64" {
		need = uint32(req.Params[0].(float64))
	} else if typeStr == "string" {
		i, err := strconv.Atoi(req.Params[0].(string))
		if err != nil {
			res.Error = MakeError(-1, "invalid jsonrpc request params[0]")
			ctx.JSON(res)
			return
		}
		need = uint32(i)
	} else {
		res.Error = MakeError(-1, "invalid jsonrpc request params[0]")
		ctx.JSON(res)
		return
	}

	//pubKeyHexStrs := make([]string, 0)
	//typeStr = reflect.TypeOf(req.Params[1]).String()
	//if typeStr[0] == '[' {
	//	argArray := reflect.ValueOf(req.Params[1]).Convert(reflect.TypeOf(req.Params[1]))
	//	for i := 0; i < argArray.Len(); i++ {
	//		typeStr = reflect.TypeOf(argArray.Index(i)).String()
	//		if typeStr != "string" {
	//			res.Error = MakeError(-1, "invalid jsonrpc request params[1], param type err")
	//			ctx.JSON(res)
	//			return
	//		}
	//		pubKeyHexStrs = append(pubKeyHexStrs, argArray.Index(i).String())
	//	}
	//}

	multiPubKeyHexStr := ""
	typeStr = reflect.TypeOf(req.Params[1]).String()
	if typeStr == "string" {
		multiPubKeyHexStr = req.Params[1].(string)
	} else {
		res.Error = MakeError(-1, "invalid jsonrpc request params[1]")
		ctx.JSON(res)
		return
	}

	pubKeyHexStrs := make([]string, 0)
	l := strings.Split(multiPubKeyHexStr, ",")
	for _, e := range l {
		pubKeyHexStrs = append(pubKeyHexStrs, e)
	}

	redeemScript, err := BTCGetRedeemScriptByPubKeys(int(need), pubKeyHexStrs)
	if err != nil {
		res.Error = MakeError(-1, err.Error())
		ctx.JSON(res)
		return
	}

	multiSigAddr, err := BTCGetMultiSignAddressByRedeemScript(redeemScript)
	if err != nil {
		res.Error = MakeError(-1, err.Error())
		ctx.JSON(res)
		return
	}

	res.Result = new(MultiSigAddressRes)
	res.Result.RedeemScript = redeemScript
	res.Result.MultiSigAddress = multiSigAddr

	ctx.JSON(res)
	return
}

func MultiSignTransactionController(ctx iris.Context, jsonRpcBody []byte) {
	var req JsonRpcRequest
	_ = json.Unmarshal(jsonRpcBody, &req)

	var res MultiSignTransactionResponse
	res.Id = req.Id

	if len(req.Params) != 4 {
		res.Error = MakeError(-1, "invalid jsonrpc request params length")
		ctx.JSON(res)
		return
	}

	rawTrxStr, multiPrivKeyEncryptHexStr, redeemScriptStr, utxosStr := "", "", "", ""
	typeStr := reflect.TypeOf(req.Params[0]).String()
	if typeStr == "string" {
		rawTrxStr = req.Params[0].(string)
	} else {
		res.Error = MakeError(-1, "invalid jsonrpc request params[0]")
		ctx.JSON(res)
		return
	}

	typeStr = reflect.TypeOf(req.Params[1]).String()
	if typeStr == "string" {
		multiPrivKeyEncryptHexStr = req.Params[1].(string)
	} else {
		res.Error = MakeError(-1, "invalid jsonrpc request params[1]")
		ctx.JSON(res)
		return
	}

	typeStr = reflect.TypeOf(req.Params[2]).String()
	if typeStr == "string" {
		redeemScriptStr = req.Params[2].(string)
	} else {
		res.Error = MakeError(-1, "invalid jsonrpc request params[2]")
		ctx.JSON(res)
		return
	}

	typeStr = reflect.TypeOf(req.Params[3]).String()
	if typeStr == "string" {
		utxosStr = req.Params[3].(string)
	} else {
		res.Error = MakeError(-1, "invalid jsonrpc request params[3]")
		ctx.JSON(res)
		return
	}

	privKeyHexStrList := make([]string, 0)
	privKeyHexStrSet := make(map[string]struct{})
	l := strings.Split(multiPrivKeyEncryptHexStr, ",")
	for _, e := range l {
		privKeyEncryptBytes, err := hex.DecodeString(e)
		if err != nil {
			res.Error = MakeError(-1, "invalid jsonrpc request params[1], privKeyEncryptHexStr not hex format string")
			ctx.JSON(res)
			return
		}

		privKeyHexStr := string(AesDecrypt(privKeyEncryptBytes, []byte(SecurityPassStr)))
		if len(privKeyHexStr) == 0 {
			res.Error = MakeError(-1, "AesDecrypt fail")
			ctx.JSON(res)
			return
		}

		privKeyHexStrList = append(privKeyHexStrList, privKeyHexStr)
		privKeyHexStrSet[privKeyHexStr] = struct{}{}
	}

	if len(privKeyHexStrList) != len(privKeyHexStrSet) {
		res.Error = MakeError(-1, "Duplicated private key input")
		ctx.JSON(res)
		return
	}

	var utxos UTXOsDetail
	//err := json.Unmarshal([]byte(utxosStr), &utxos)
	//if err != nil {
	//	res.Error = MakeError(-1, "invalid jsonrpc request params[3], Unmarshal fail")
	//	ctx.JSON(res)
	//	return
	//}
	_ = utxosStr

	trxSigStrList := make([]string, 0)
	for _, key := range privKeyHexStrList {
		trxSigStr, err := BTCMultiSignRawTransaction(rawTrxStr, redeemScriptStr, key, utxos)
		if err != nil {
			res.Error = MakeError(-1, fmt.Sprintf("multi sign raw transaction fail: %s", err.Error()))
			ctx.JSON(res)
			return
		}
		trxSigStrList = append(trxSigStrList, trxSigStr)
	}

	reqParams := make([]interface{}, 0)
	reqParams = append(reqParams, trxSigStrList)

	Info.Println(fmt.Sprintf("rpcClient to [%s] for combining raw transaction", GlobalConfig.ServerUrl))
	rpcClient := jsonrpc.NewClient(GlobalConfig.ServerUrl)
	rpcResponse, err := rpcClient.Call("combinerawtransaction", reqParams)
	if err != nil {
		res.Error = MakeError(-1, fmt.Sprintf("rpc combinerawtransaction fail: %s", err.Error()))
		ctx.JSON(res)
		return
	}

	if rpcResponse.Error != nil {
		Error.Println("combinerawtransaction rpcResponse: %s", rpcResponse.Error.Error())
		res.Error = MakeError(-1, "rpc combinerawtransaction fail: rpcResponse.Error not nil")
		ctx.JSON(res)
		return
	}

	if rpcResponse.Result == nil {
		Error.Println("combinerawtransaction rpcResponse is nil")
		res.Error = MakeError(-1, "rpc combinerawtransaction fail: rpcResponse.Result is nil")
		ctx.JSON(res)
		return
	}

	trxSigStr := rpcResponse.Result.(string)
	Info.Println("rawTrxCombinedStr:", trxSigStr)

	// set trx utxos state to pending
	trx, err := BTCUnPackRawTransaction(rawTrxStr)
	if err != nil {
		Error.Println("BTCUnPackRawTransaction rawTrx fail:", err.Error())
		res.Error = MakeError(-1, "BTCUnPackRawTransaction rawTrx fail")
		ctx.JSON(res)
		return
	}
	for _, vin := range trx.Vin {
		txId := vin.PrevOut.Hash.GetHex()
		vout := vin.PrevOut.N
		err = GlobalDBMgr.TblUtxoMgr.UpdateUtxoPendingState(txId, int(vout), 1)
		if err != nil {
			Error.Printf("UpdateUtxoPendingState [%s/%d] fail: %s", txId, int(vout), err.Error())
			res.Error = MakeError(-1, "UpdateUtxoPendingState fail")
			ctx.JSON(res)
			return
		}
	}

	res.Result = &trxSigStr
	ctx.JSON(res)
	return
}

func ImportAddressesController(ctx iris.Context, jsonRpcBody []byte) {
	var req JsonRpcRequest
	_ = json.Unmarshal(jsonRpcBody, &req)

	var res ImportAddressesResponse
	res.Id = req.Id

	addresses := make([]string, 0)
	for _, param := range req.Params {
		typeStr := reflect.TypeOf(param).String()
		if typeStr == "string" {
			addresses = append(addresses, param.(string))
		} else {
			res.Error = MakeError(-1, "invalid jsonrpc request params")
			ctx.JSON(res)
			return
		}
	}

	err := GlobalDBMgr.TblAddressMgr.AddNewAddresses(addresses)
	if err != nil {
		res.Error = MakeError(-1, err.Error())
		ctx.JSON(res)
		return
	}

	res.Result = nil
	ctx.JSON(res)
	return
}

func QueryUtxosController(ctx iris.Context, jsonRpcBody []byte) {
	var req JsonRpcRequest
	_ = json.Unmarshal(jsonRpcBody, &req)

	var res QueryUtxosResponse
	res.Id = req.Id

	if len(req.Params) != 1 {
		res.Error = MakeError(-1, "invalid jsonrpc request params length")
		ctx.JSON(res)
		return
	}

	addr := ""
	typeStr := reflect.TypeOf(req.Params[0]).String()
	if typeStr == "string" {
		addr = req.Params[0].(string)
	} else {
		res.Error = MakeError(-1, "invalid jsonrpc request params[0]")
		ctx.JSON(res)
		return
	}

	utxos, err := GlobalDBMgr.TblUtxoMgr.ListAddrUtxos(addr)
	if err != nil {
		res.Error = MakeError(-1, err.Error())
		ctx.JSON(res)
		return
	}

	utxosRes := make([]UtxoRes, 0)
	for _, utxo := range utxos {
		amountFloat, _ := strconv.ParseFloat(utxo.Amount, 64)
		utxoRes := UtxoRes{Txid: utxo.Txid,
			Address:      utxo.Address,
			Amount:       fmt.Sprintf("%.08f", amountFloat),
			ScriptPubKey: utxo.Scriptpubkey,
			Vout:         utxo.Vout}
		utxosRes = append(utxosRes, utxoRes)
	}

	res.Result = &utxosRes
	ctx.JSON(res)
	return
}

func Controller(ctx iris.Context) {
	id, funcName, jsonRpcBody, err := ReadJsonRpcBody(ctx)
	if err != nil {
		Info.Println("Internal Error:", err.Error())
		SetInternalError(ctx, err.Error())
		return
	}

	if funcName == "generate_address" {
		GenerateAddressController(ctx, jsonRpcBody)
	} else if funcName == "sign_transaction" {
		SignTransactionController(ctx, jsonRpcBody)
	} else if funcName == "generate_multi_address" {
		GenerateMultiAddressController(ctx, jsonRpcBody)
	} else if funcName == "multi_sign_transaction" {
		MultiSignTransactionController(ctx, jsonRpcBody)
	} else if funcName == "import_addresses" {
		ImportAddressesController(ctx, jsonRpcBody)
	} else if funcName == "query_utxos" {
		QueryUtxosController(ctx, jsonRpcBody)
	} else {
		var res JsonRpcResponse
		res.Id = id
		res.Result = nil
		res.Error = MakeError(-1, "invalid jsonrpc request method")
		ctx.JSON(res)
	}
}
