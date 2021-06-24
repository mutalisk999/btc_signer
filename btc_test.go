package main

import (
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcec"
	"github.com/mutalisk999/bitcoin-lib/src/base58"
	"github.com/ybbus/jsonrpc"
	"testing"
)

func TestBTCPrivKeyBytesToWIF(t *testing.T) {
	privKeyB58 := "KzutU4gAuMqf9qFayh57Xb6JkCZv6o4jKkZmuEpk5tpVfvvzKqUT"
	privKeyBytes, _ := base58.Decode(privKeyB58)
	privKeyBytes = privKeyBytes[1:33]
	privKeyHex := hex.EncodeToString(privKeyBytes)
	fmt.Println("privKeyHex:", privKeyHex)

	wifKey, _ := BTCPrivKeyBytesToWIF(privKeyBytes)
	fmt.Println("wifKey:", wifKey)
}

func TestBTCCalcAddressByPubKey(t *testing.T) {
	privKeyB58 := "KzutU4gAuMqf9qFayh57Xb6JkCZv6o4jKkZmuEpk5tpVfvvzKqUT"
	privKeyBytes, _ := base58.Decode(privKeyB58)
	privKeyBytes = privKeyBytes[1:33]

	privKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), privKeyBytes)
	pubKey := privKey.PubKey()
	pubKeyBytes := pubKey.SerializeUncompressed()[1:]
	fmt.Println("pubKeyStr:", hex.EncodeToString(pubKeyBytes))

	addrStr, _ := BTCCalcAddressByPubKey(hex.EncodeToString(pubKeyBytes))
	fmt.Println("addrStr:", addrStr)
}

func TestBTCGetP2PKHScriptPubKey(t *testing.T) {
	pubKeyHexStr := "f6613d8d57a0baa36ed7afd86513d6f6a492417022b0aa37f32a7ac429a9714deff765056b37a15bc3c8e4b45d9b58dde4a3ec92455a780774b08972ad3aca35"
	scriptPubKey, _ := BTCGetP2PKHScriptPubKey(pubKeyHexStr)
	fmt.Println("scriptPubKey:", hex.EncodeToString(scriptPubKey))
}

func TestBTCSignRawTransaction(t *testing.T) {
	privKeyB58 := "KzutU4gAuMqf9qFayh57Xb6JkCZv6o4jKkZmuEpk5tpVfvvzKqUT"
	privKeyBytes, _ := base58.Decode(privKeyB58)
	privKeyBytes = privKeyBytes[1:33]
	privKeyHex := hex.EncodeToString(privKeyBytes)

	// not necessary
	var utxos UTXOsDetail

	rawTrxStr := "0200000001b93156b79bc8c535d13a2af4dede3bbd7866aaf6badc38e135c08b0d0ec627930000000000ffffffff01606b042a0100000017a914d62ba8e2fb39688d4afcca101cc4ff7abf57f27f8700000000"
	rawTrxSignedStr, _ := BTCSignRawTransaction(rawTrxStr, privKeyHex, utxos)
	fmt.Println("rawTrxSignedStr:", rawTrxSignedStr)
}

//cryptedHexStr: 3ceed5da0cb6c7a84d6687e0409735172f14b2c85d319ac8b27d69f82f72e5372308590488b206fd6eae2cfe6259cb9699d7b7e915dfbf419c2c7fb48c9e1d4c2d2fede6a8001a1a2a5abb2e46426445
//pubKeyHexStr: 0303b98c2753cb48a456d88c89727936797d7fa890eb600dddf32940a1e835188b
//cryptedHexStr: daefb6933ccb396deb1006bb7700bd3e43765b58af38b3c14fa38e5f18d0f6ae2ece9f36f4880f32fdba9ab09454355e7a7212846ec7a82ddd5ef188dddaa8bd2d2fede6a8001a1a2a5abb2e46426445
//pubKeyHexStr: 02cd7c2fe2be798cf062de43783177fab7a3436af29a6aeb65c78399cbf25f84a9
//cryptedHexStr: 7664f5a2bfccd83c7701afbd45013f8865c9f040984fea7dc9a22aed571a097bad4ff8e7a4e459aa79b4cee8cb81782a710cc4f7d54d4bf82e59671dbe9065312d2fede6a8001a1a2a5abb2e46426445
//pubKeyHexStr: 0351519038c945c71a5268ae27729731f886b56b5e14b202d351530a92bdec8f59
//cryptedHexStr: 7042666190a4c0b50ef130ba05a795ef1dc424a29c3bb61dc6bc13e3794788e83d904444ac41bba14f35104dc4c83a155e32926bf6903e13312d1b0828dc1e802d2fede6a8001a1a2a5abb2e46426445
//pubKeyHexStr: 02ec30578e5647e00a20ad3ef98b08381cd57e28e00293ff5a27bf0981bac008b5
//cryptedHexStr: 3a21c428616d451409f4c78d714d3855a6b99797996fb705457f076cc82d877132322e79a21c96d85cc8eaeb6bc02b6fff0a288f031c2c2c8973a843ebbbb70d2d2fede6a8001a1a2a5abb2e46426445
//pubKeyHexStr: 036ff86d871899f06bd68f201c894cd872a19b15f4e284c2d86227176fbdc0a9bf
func TestBTCGenerateNewAddress(t *testing.T) {
	_, privKeyHexStr, pubKeyHexStr, _, _ := BTCGenerateNewAddress()
	cryptedBytes := AesEncrypt(privKeyHexStr, SecurityPassStr)
	cryptedHexStr := hex.EncodeToString(cryptedBytes)
	fmt.Println("cryptedHexStr:", cryptedHexStr)
	fmt.Println("pubKeyHexStr:", pubKeyHexStr)
}

//redeemScript: 53210303b98c2753cb48a456d88c89727936797d7fa890eb600dddf32940a1e835188b2102cd7c2fe2be798cf062de43783177fab7a3436af29a6aeb65c78399cbf25f84a9210351519038c945c71a5268ae27729731f886b56b5e14b202d351530a92bdec8f592102ec30578e5647e00a20ad3ef98b08381cd57e28e00293ff5a27bf0981bac008b521036ff86d871899f06bd68f201c894cd872a19b15f4e284c2d86227176fbdc0a9bf55ae
func TestBTCGetRedeemScriptByPubKeys(t *testing.T) {
	pubKeyHexStr1 := "0303b98c2753cb48a456d88c89727936797d7fa890eb600dddf32940a1e835188b"
	pubKeyHexStr2 := "02cd7c2fe2be798cf062de43783177fab7a3436af29a6aeb65c78399cbf25f84a9"
	pubKeyHexStr3 := "0351519038c945c71a5268ae27729731f886b56b5e14b202d351530a92bdec8f59"
	pubKeyHexStr4 := "02ec30578e5647e00a20ad3ef98b08381cd57e28e00293ff5a27bf0981bac008b5"
	pubKeyHexStr5 := "036ff86d871899f06bd68f201c894cd872a19b15f4e284c2d86227176fbdc0a9bf"
	redeemScript, _ := BTCGetRedeemScriptByPubKeys(3, []string{pubKeyHexStr1, pubKeyHexStr2, pubKeyHexStr3, pubKeyHexStr4, pubKeyHexStr5})
	fmt.Println("redeemScript:", redeemScript)
}

//multiSigAddr: 3MDSq8EZGz71f9BCLy1tpndHjvbXH8Wj4V
func TestBTCGetMultiSignAddressByRedeemScript(t *testing.T) {
	redeemScript := "53210303b98c2753cb48a456d88c89727936797d7fa890eb600dddf32940a1e835188b2102cd7c2fe2be798cf062de43783177fab7a3436af29a6aeb65c78399cbf25f84a9210351519038c945c71a5268ae27729731f886b56b5e14b202d351530a92bdec8f592102ec30578e5647e00a20ad3ef98b08381cd57e28e00293ff5a27bf0981bac008b521036ff86d871899f06bd68f201c894cd872a19b15f4e284c2d86227176fbdc0a9bf55ae"
	multiSigAddr, _ := BTCGetMultiSignAddressByRedeemScript(redeemScript)
	fmt.Println("multiSigAddr:", multiSigAddr)
}

func TestBTCMultiSignRawTransaction(t *testing.T) {
	redeemScript := "53210303b98c2753cb48a456d88c89727936797d7fa890eb600dddf32940a1e835188b2102cd7c2fe2be798cf062de43783177fab7a3436af29a6aeb65c78399cbf25f84a9210351519038c945c71a5268ae27729731f886b56b5e14b202d351530a92bdec8f592102ec30578e5647e00a20ad3ef98b08381cd57e28e00293ff5a27bf0981bac008b521036ff86d871899f06bd68f201c894cd872a19b15f4e284c2d86227176fbdc0a9bf55ae"
	rawTrxStr := "02000000012df0530a929aff47252c24e9b371fc39acd02ace2f79189a3e55922c807956de0000000000ffffffff01c0e4022a0100000017a914d62ba8e2fb39688d4afcca101cc4ff7abf57f27f8700000000"

	// not necessary
	var utxos UTXOsDetail

	privKeyEncryptHexStr1 := "3ceed5da0cb6c7a84d6687e0409735172f14b2c85d319ac8b27d69f82f72e5372308590488b206fd6eae2cfe6259cb9699d7b7e915dfbf419c2c7fb48c9e1d4c2d2fede6a8001a1a2a5abb2e46426445"
	privKeyEncryptHexStr2 := "daefb6933ccb396deb1006bb7700bd3e43765b58af38b3c14fa38e5f18d0f6ae2ece9f36f4880f32fdba9ab09454355e7a7212846ec7a82ddd5ef188dddaa8bd2d2fede6a8001a1a2a5abb2e46426445"
	privKeyEncryptHexStr3 := "7664f5a2bfccd83c7701afbd45013f8865c9f040984fea7dc9a22aed571a097bad4ff8e7a4e459aa79b4cee8cb81782a710cc4f7d54d4bf82e59671dbe9065312d2fede6a8001a1a2a5abb2e46426445"

	privKeyEncryptBytes1, _ := hex.DecodeString(privKeyEncryptHexStr1)
	privKeyHexStr1 := string(AesDecrypt(privKeyEncryptBytes1, []byte(SecurityPassStr)))

	trxSignedData1, _ := BTCMultiSignRawTransaction(rawTrxStr, redeemScript, privKeyHexStr1, utxos)
	fmt.Println("trxSignedData1:", trxSignedData1)

	privKeyEncryptBytes2, _ := hex.DecodeString(privKeyEncryptHexStr2)
	privKeyHexStr2 := string(AesDecrypt(privKeyEncryptBytes2, []byte(SecurityPassStr)))

	trxSignedData2, _ := BTCMultiSignRawTransaction(rawTrxStr, redeemScript, privKeyHexStr2, utxos)
	fmt.Println("trxSignedData2:", trxSignedData2)

	privKeyEncryptBytes3, _ := hex.DecodeString(privKeyEncryptHexStr3)
	privKeyHexStr3 := string(AesDecrypt(privKeyEncryptBytes3, []byte(SecurityPassStr)))

	trxSignedData3, _ := BTCMultiSignRawTransaction(rawTrxStr, redeemScript, privKeyHexStr3, utxos)
	fmt.Println("trxSignedData3:", trxSignedData3)
}

func TestBTCCombineSignedTrx(t *testing.T) {
	trxSigStrList := make([]string, 0)
	trxSigStrList = append(trxSigStrList,
		"0200000001a0262971a6196ddb554140f12aae68b738852247121f6769dcb6142f9cf6ede300000000f800473044022057e93001dab90716c333df85e5a591352c3468820cd02ceb4e2cd947991440c402201617fa72af1d4a39b487a19c1a32cf5c246c76f8aee8f08715fb54cae91ae7b8014cad532102feade17d70e308af54fcce9baee1c3d34066f100798c9839efee9b1d281abd8921030c6ed4af9836f9772e2b1813e4cd9e30c49f7b9924d59bfca6d4f54e39aaeb162102657d332743056fe81c72be70748e71ad5248524caaab5551c365326baac5279d2103eb860f625fc71dd1710f56a6c3d0082c28ea42e542966146d776296bf833fba0210229fafc185334bb65bc23fba054cb693b65a521fdb545bd73135da3bc5ebc2fca55aeffffffff0210270000000000001976a914451328751fbb4d981aea377f84511aede5c2b9e788ac701101000000000017a91427368ea17968c43f8dd6b5e944457300e4b539208700000000",
		"0200000001a0262971a6196ddb554140f12aae68b738852247121f6769dcb6142f9cf6ede300000000f90048304502210096519aa83b9e64822f3aa9302bf455c9b423b6e85a3ce96d62dd2cdfede9341c02201006e0fe2728dac2ffa6adb86efa955d4c3bb9a84fd21db05a70e1c14c86e32c014cad532102feade17d70e308af54fcce9baee1c3d34066f100798c9839efee9b1d281abd8921030c6ed4af9836f9772e2b1813e4cd9e30c49f7b9924d59bfca6d4f54e39aaeb162102657d332743056fe81c72be70748e71ad5248524caaab5551c365326baac5279d2103eb860f625fc71dd1710f56a6c3d0082c28ea42e542966146d776296bf833fba0210229fafc185334bb65bc23fba054cb693b65a521fdb545bd73135da3bc5ebc2fca55aeffffffff0210270000000000001976a914451328751fbb4d981aea377f84511aede5c2b9e788ac701101000000000017a91427368ea17968c43f8dd6b5e944457300e4b539208700000000",
		"0200000001a0262971a6196ddb554140f12aae68b738852247121f6769dcb6142f9cf6ede300000000f900483045022100eec9598e486c4582ebe03c1735eed918d3ee9bf55416044ab0cff8f305663648022046dc8d04d5bfe1abe832d1ef83b5cc5c0b98d9a42bd252e282869a58dc7c2dae014cad532102feade17d70e308af54fcce9baee1c3d34066f100798c9839efee9b1d281abd8921030c6ed4af9836f9772e2b1813e4cd9e30c49f7b9924d59bfca6d4f54e39aaeb162102657d332743056fe81c72be70748e71ad5248524caaab5551c365326baac5279d2103eb860f625fc71dd1710f56a6c3d0082c28ea42e542966146d776296bf833fba0210229fafc185334bb65bc23fba054cb693b65a521fdb545bd73135da3bc5ebc2fca55aeffffffff0210270000000000001976a914451328751fbb4d981aea377f84511aede5c2b9e788ac701101000000000017a91427368ea17968c43f8dd6b5e944457300e4b539208700000000")

	reqParams := make([]interface{}, 0)
	reqParams = append(reqParams, trxSigStrList)

	rpcClient := jsonrpc.NewClient("http://a:b@192.168.1.160:5100")
	rpcResponse, err := rpcClient.Call("combinerawtransaction", reqParams)
	if err != nil {
		fmt.Println("err:", err)
		return
	}

	trxSigStr := rpcResponse.Result
	fmt.Println("trxSigStr:", trxSigStr)
}
