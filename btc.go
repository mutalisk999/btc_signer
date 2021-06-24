package main

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/btcec"
	"github.com/mutalisk999/bitcoin-lib/src/base58"
	"github.com/mutalisk999/bitcoin-lib/src/blob"
	"github.com/mutalisk999/bitcoin-lib/src/keyid"
	"github.com/mutalisk999/bitcoin-lib/src/pubkey"
	"github.com/mutalisk999/bitcoin-lib/src/script"
	"github.com/mutalisk999/bitcoin-lib/src/serialize"
	"github.com/mutalisk999/bitcoin-lib/src/transaction"
	"github.com/mutalisk999/bitcoin-lib/src/utility"
	"io"
)

type UTXODetail struct {
	TxId          string `json:"txid"`
	Vout          int    `json:"vout"`
	Address       string `json:"address"`
	Account       string `json:"account"`
	ScriptPubKey  string `json:"scriptPubKey"`
	RedeemScript  string `json:"redeemScript"`
	Amount        int64  `json:"amount"`
	Confirmations int    `json:"confirmations"`
	Spendable     bool   `json:"spendable"`
	Solvable      bool   `json:"solvable"`
}
type UTXOsDetail []UTXODetail

func BTCPrivKeyBytesToWIF(privKeyBytes []byte) (string, error) {
	if len(privKeyBytes) != 32 {
		return "", errors.New("invalid privKeyBytes size")
	}

	privkeyPaddingBytes := make([]byte, 38, 38)
	// mainnet version
	privkeyPaddingBytes[0] = 0x80
	copy(privkeyPaddingBytes[1:], privKeyBytes[0:32])
	// compress privkey
	privkeyPaddingBytes[33] = 0x1

	bytes := utility.Sha256(utility.Sha256(privkeyPaddingBytes[0:34]))
	copy(privkeyPaddingBytes[34:], bytes[0:4])

	wifKeyStr := base58.Encode(privkeyPaddingBytes)
	return wifKeyStr, nil
}

func BTCGetCompressPubKey(pubKeyBytes []byte) ([]byte, error) {
	if len(pubKeyBytes) != 64 {
		return nil, errors.New("invalid pubKeyBytes size")
	}

	pubkeyCompress := make([]byte, 33, 33)
	if pubKeyBytes[63]%2 == 0 {
		pubkeyCompress[0] = 0x2
	} else {
		pubkeyCompress[0] = 0x3
	}
	copy(pubkeyCompress[1:], pubKeyBytes[0:32])
	return pubkeyCompress, nil
}

func BTCCalcAddressByPubKey(pubKeyStr string) (string, error) {
	pubKeyBytes, err := hex.DecodeString(pubKeyStr)
	if err != nil {
		return "", err
	}

	pubkeyCompress, err := BTCGetCompressPubKey(pubKeyBytes)
	if err != nil {
		return "", err
	}

	pubKey := new(pubkey.PubKey)
	pubKey.SetPubKeyData(pubkeyCompress)

	keyIdBytes, err := pubKey.CalcKeyIDBytes()
	if err != nil {
		return "", err
	}
	keyId := new(keyid.KeyID)
	keyId.SetKeyIDData(keyIdBytes)

	var version byte
	version = 0

	addrStr, err := keyId.ToBase58Address(version)
	if err != nil {
		return "", err
	}

	return addrStr, nil
}

func BTCGenerateNewAddress() (string, string, string, string, error) {
	privkey, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		return "", "", "", "", err
	}
	privkeyBytes := privkey.Serialize()
	privkeyBlob := blob.Byteblob{}
	privkeyBlob.SetData(privkeyBytes)

	privkeyWif, err := BTCPrivKeyBytesToWIF(privkeyBytes)
	if err != nil {
		return "", "", "", "", err
	}

	pubkey := privkey.PubKey()
	pubkeyUncompressedBytes := pubkey.SerializeUncompressed()
	pubkeyOrigBytes := pubkeyUncompressedBytes[1:]
	pubkeyOrigBlob := blob.Byteblob{}
	pubkeyOrigBlob.SetData(pubkeyOrigBytes)

	pubkeyCompressedBytes, err := BTCGetCompressPubKey(pubkeyOrigBytes)
	if err != nil {
		return "", "", "", "", err
	}
	pubkeyCompressedBlob := blob.Byteblob{}
	pubkeyCompressedBlob.SetData(pubkeyCompressedBytes)

	addrStr, err := BTCCalcAddressByPubKey(pubkeyOrigBlob.GetHex())
	if err != nil {
		return "", "", "", "", err
	}

	return privkeyWif, privkeyBlob.GetHex(), pubkeyCompressedBlob.GetHex(), addrStr, nil
}

func BTCUnPackRawTransaction(rawTrx string) (*transaction.Transaction, error) {
	Blob := new(blob.Byteblob)
	err := Blob.SetHex(rawTrx)
	if err != nil {
		return nil, err
	}
	bytesBuf := bytes.NewBuffer(Blob.GetData())
	bufReader := io.Reader(bytesBuf)
	trx := new(transaction.Transaction)
	err = trx.UnPack(bufReader)
	if err != nil {
		return nil, err
	}
	return trx, nil
}

func BTCPackRawTransaction(trxSig transaction.Transaction) (string, error) {
	bytesBuf := bytes.NewBuffer([]byte{})
	bufWriter := io.Writer(bytesBuf)
	err := trxSig.Pack(bufWriter)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytesBuf.Bytes()), nil
}

func BTCCombineSignatureAndPubKey(signature []byte, pubKey []byte) []byte {
	scriptSig := make([]byte, 0, 1+len(signature)+1+len(pubKey))
	scriptSig = append(scriptSig, byte(len(signature)))
	scriptSig = append(scriptSig, signature...)
	scriptSig = append(scriptSig, byte(len(pubKey)))
	scriptSig = append(scriptSig, pubKey...)
	//Info.Println("scriptSig:", hex.EncodeToString(scriptSig))
	return scriptSig
}

func SerializeDerEncoding(rBytes []byte, sBytes []byte) ([]byte, error) {
	if len(rBytes) != 32 {
		return nil, errors.New("invalid rBytes len")
	}
	if len(sBytes) != 32 {
		return nil, errors.New("invalid sBytes len")
	}

	var r []byte
	r = append(r, 0)
	r = append(r, rBytes...)
	var s []byte
	s = append(s, 0)
	s = append(s, sBytes...)
	for {
		if len(r) > 1 && r[0] == 0 && r[1] < 0x80 {
			r = r[1:]
		} else {
			break
		}
	}
	for {
		if len(s) > 1 && s[0] == 0 && s[1] < 0x80 {
			s = s[1:]
		} else {
			break
		}
	}

	size := 6 + len(r) + len(s)
	signedData := make([]byte, size, size)
	signedData[0] = 0x30
	signedData[1] = 4 + byte(len(r)) + byte(len(s))
	signedData[2] = 0x2
	signedData[3] = byte(len(r))
	copy(signedData[4:4+len(r)], r)
	signedData[4+len(r)] = 0x2
	signedData[5+len(r)] = byte(len(s))
	copy(signedData[6+len(r):6+len(r)+len(s)], s)

	return signedData, nil
}

func BTCCoinSignTrx(privKeyBytes []byte, signData []byte) ([]byte, error) {
	privKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), privKeyBytes)
	signature, err := privKey.Sign(signData)
	if err != nil {
		return nil, err
	}
	signedData := signature.Serialize()
	return signedData, nil
}

func BTCCoinVerifyTrx(pubKeyBytes []byte, signData []byte, signedData []byte) (bool, error) {
	pubKey, err := btcec.ParsePubKey(pubKeyBytes, btcec.S256())
	if err != nil {
		return false, err
	}
	signature, err := btcec.ParseSignature(signedData, btcec.S256())
	if err != nil {
		return false, err
	}
	verified := signature.Verify(signData, pubKey)
	if !verified {
		return false, nil
	}
	return true, nil
}

func BTCGetP2PKHScriptPubKey(pubKeyStr string) ([]byte, error) {
	pubKeyBytes, err := hex.DecodeString(pubKeyStr)
	if err != nil {
		return nil, err
	}
	pubkeyCompress, err := BTCGetCompressPubKey(pubKeyBytes)
	if err != nil {
		return nil, err
	}

	pubKey := new(pubkey.PubKey)
	pubKey.SetPubKeyData(pubkeyCompress)

	keyIdBytes, err := pubKey.CalcKeyIDBytes()
	if err != nil {
		return nil, err
	}
	keyId := new(keyid.KeyID)
	keyId.SetKeyIDData(keyIdBytes)
	keyIdData, _ := keyId.GetKeyIDData()

	bufBytes := make([]byte, 0)
	bufBytes = append(bufBytes, script.OP_DUP, script.OP_HASH160, byte(keyid.KEY_ID_SIZE))
	bufBytes = append(bufBytes, keyIdData...)
	bufBytes = append(bufBytes, script.OP_EQUALVERIFY, script.OP_CHECKSIG)
	return bufBytes, nil
}

func BTCSignRawTransaction(rawTrx string, privKeyStr string, utxos []UTXODetail) (string, error) {
	privKeyBytes, err := hex.DecodeString(privKeyStr)
	if err != nil {
		return "", err
	}

	_, pubKey := btcec.PrivKeyFromBytes(btcec.S256(), privKeyBytes)
	pubKeyBytes := pubKey.SerializeUncompressed()[1:]

	pubkeyCompress, err := BTCGetCompressPubKey(pubKeyBytes)
	if err != nil {
		return "", err
	}

	Info.Println("rawTrxStr:", rawTrx)

	trx, err := BTCUnPackRawTransaction(rawTrx)
	if err != nil {
		return "", err
	}

	signedDataList := make([][]byte, len(trx.Vin))

	// add scriptPubKey
	for i := 0; i < len(trx.Vin); i++ {
		trxTemp, err := BTCUnPackRawTransaction(rawTrx)
		if err != nil {
			return "", err
		}

		p2pkhScriptPubKey, err := BTCGetP2PKHScriptPubKey(hex.EncodeToString(pubKeyBytes))
		if err != nil {
			return "", err
		}
		trxTemp.Vin[i].ScriptSig.SetScriptBytes(p2pkhScriptPubKey)

		rawTrxWithScript, err := BTCPackRawTransaction(*trxTemp)
		if err != nil {
			return "", err
		}

		rawTrxBytes, err := hex.DecodeString(rawTrxWithScript)
		if err != nil {
			return "", err
		}
		// append SIGHASH_ALL
		rawTrxBytes = append(rawTrxBytes, []byte{0x1, 0x0, 0x0, 0x0}...)
		hashBytes := utility.Sha256(utility.Sha256(rawTrxBytes))

		//Info.Println("rawTrxBytes:", hex.EncodeToString(rawTrxBytes))
		//Info.Println("hashBytes:", hex.EncodeToString(hashBytes))

		// signature
		signedData, err := BTCCoinSignTrx(privKeyBytes, hashBytes)
		if err != nil {
			return "", err
		}

		verifyOk, err := BTCCoinVerifyTrx(pubkeyCompress, hashBytes, signedData)
		if err != nil {
			return "", err
		}
		if !verifyOk {
			return "", errors.New("verify signature error")
		}

		Info.Println("signedDataStr:", hex.EncodeToString(signedData))

		// append SIGHASH_ALL
		signedData = append(signedData, 0x1)

		scriptSig := BTCCombineSignatureAndPubKey(signedData, pubkeyCompress)

		signedDataList[i] = scriptSig
	}

	for i := 0; i < len(trx.Vin); i++ {
		trx.Vin[i].ScriptSig.SetScriptBytes(signedDataList[i])
	}

	trxSigStr, err := BTCPackRawTransaction(*trx)
	if err != nil {
		return "", err
	}

	Info.Println("rawTrxSignedStr:", trxSigStr)

	return trxSigStr, nil
}

func BTCGetRedeemScriptByPubKeys(needCount int, pubKeyStrList []string) (string, error) {
	if needCount <= 0 || needCount > 16 {
		return "", errors.New("BTCGetRedeemScriptByPubKeys error: invalid needCount")
	}
	if len(pubKeyStrList) == 0 || len(pubKeyStrList) > 16 {
		return "", errors.New("BTCGetRedeemScriptByPubKeys error: invalid pubKeyStrList size")
	}
	if needCount > len(pubKeyStrList) {
		return "", errors.New("BTCGetRedeemScriptByPubKeys error: needCount greater than pubKeyStrList size")
	}

	bytesBuf := bytes.NewBuffer([]byte{})
	bufWriter := io.Writer(bytesBuf)
	err := serialize.PackUint8(bufWriter, uint8(needCount+0x50))
	if err != nil {
		return "", err
	}
	for _, pubKeyStr := range pubKeyStrList {
		pubKeyBytes, err := hex.DecodeString(pubKeyStr)
		if err != nil {
			return "", err
		}

		var pubKeyCpsBytes []byte
		if len(pubKeyBytes) == 33 && (pubKeyBytes[0] == 0x2 || pubKeyBytes[0] == 0x3) {
			pubKeyCpsBytes = pubKeyBytes
		} else {
			if len(pubKeyBytes) == 65 && pubKeyBytes[0] == 0x4 {
				pubKeyBytes = pubKeyBytes[1:]
			}
			pubKeyCpsBytes, err = BTCGetCompressPubKey(pubKeyBytes)
			if err != nil {
				return "", err
			}
		}

		pubKey := new(pubkey.PubKey)
		pubKey.SetPubKeyData(pubKeyCpsBytes)

		err = pubKey.Pack(bufWriter)
		if err != nil {
			return "", err
		}
	}
	err = serialize.PackUint8(bufWriter, uint8(len(pubKeyStrList)+0x50))
	if err != nil {
		return "", err
	}
	err = serialize.PackUint8(bufWriter, uint8(script.OP_CHECKMULTISIG))
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytesBuf.Bytes()), nil
}

func BTCGetMultiSignAddressByRedeemScript(redeemScriptStr string) (string, error) {
	redeemScript, err := hex.DecodeString(redeemScriptStr)
	if err != nil {
		return "", err
	}

	scriptIdBytes := utility.Hash160(redeemScript)
	keyId := new(keyid.KeyID)
	keyId.SetKeyIDData(scriptIdBytes)

	var version byte
	version = 5

	addrStr, err := keyId.ToBase58Address(version)
	if err != nil {
		return "", err
	}
	return addrStr, nil
}

func BTCCombineSignatureAndRedeemScript(signature []byte, redeemScriptBytes []byte) ([]byte, error) {
	bytesBuf := bytes.NewBuffer([]byte{})
	bufWriter := io.Writer(bytesBuf)
	err := serialize.PackUint8(bufWriter, script.OP_0)
	if err != nil {
		return nil, err
	}
	signatureScript := new(script.Script)
	signatureScript.SetScriptBytes(signature)
	err = signatureScript.Pack(bufWriter)
	if err != nil {
		return nil, err
	}

	if len(redeemScriptBytes) < int(script.OP_PUSHDATA1) {
	} else {
		opPushData := uint8(0)
		if len(redeemScriptBytes) <= 0xff {
			opPushData = script.OP_PUSHDATA1
		} else if len(redeemScriptBytes) <= 0xffff {
			opPushData = script.OP_PUSHDATA2
		} else {
			opPushData = script.OP_PUSHDATA4
		}
		err = serialize.PackUint8(bufWriter, opPushData)
		if err != nil {
			return nil, err
		}
	}

	redeemScrip := new(script.Script)
	redeemScrip.SetScriptBytes(redeemScriptBytes)
	err = redeemScrip.Pack(bufWriter)
	if err != nil {
		return nil, err
	}
	return bytesBuf.Bytes(), nil
}

func BTCMultiSignRawTransaction(rawTrx string, redeemScriptStr string, privKeyStr string, utxos []UTXODetail) (string, error) {
	privKeyBytes, err := hex.DecodeString(privKeyStr)
	if err != nil {
		return "", err
	}

	_, pubKey := btcec.PrivKeyFromBytes(btcec.S256(), privKeyBytes)
	pubKeyBytes := pubKey.SerializeUncompressed()[1:]

	pubkeyCompress, err := BTCGetCompressPubKey(pubKeyBytes)
	if err != nil {
		Error.Println("BTCGetCompressPubKey fail:", err.Error())
		return "", err
	}

	Info.Println("rawTrxStr:", rawTrx)

	redeemScriptBytes, err := hex.DecodeString(redeemScriptStr)
	if err != nil {
		Error.Println("DecodeString redeemScriptStr fail:", err.Error())
		return "", err
	}

	trx, err := BTCUnPackRawTransaction(rawTrx)
	if err != nil {
		Error.Println("BTCUnPackRawTransaction rawTrx fail:", err.Error())
		return "", err
	}

	signedDataList := make([][]byte, len(trx.Vin))

	// add scriptPubKey
	for i := 0; i < len(trx.Vin); i++ {
		trxTemp, err := BTCUnPackRawTransaction(rawTrx)
		if err != nil {
			Error.Println("BTCUnPackRawTransaction rawTrx fail:", err.Error())
			return "", err
		}

		trxTemp.Vin[i].ScriptSig.SetScriptBytes(redeemScriptBytes)

		rawTrxWithScript, err := BTCPackRawTransaction(*trxTemp)
		if err != nil {
			Error.Println("BTCPackRawTransaction fail:", err.Error())
			return "", err
		}

		rawTrxBytes, err := hex.DecodeString(rawTrxWithScript)
		if err != nil {
			Error.Println("DecodeString rawTrxWithScript fail:", err.Error())
			return "", err
		}
		// append SIGHASH_ALL
		rawTrxBytes = append(rawTrxBytes, []byte{0x1, 0x0, 0x0, 0x0}...)
		hashBytes := utility.Sha256(utility.Sha256(rawTrxBytes))

		//Info.Println("rawTrxBytes:", hex.EncodeToString(rawTrxBytes))
		//Info.Println("hashBytes:", hex.EncodeToString(hashBytes))

		// signature
		signedData, err := BTCCoinSignTrx(privKeyBytes, hashBytes)
		if err != nil {
			Error.Println("BTCCoinSignTrx fail:", err.Error())
			return "", err
		}

		verifyOk, err := BTCCoinVerifyTrx(pubkeyCompress, hashBytes, signedData)
		if err != nil {
			Error.Println("BTCCoinVerifyTrx fail:", err.Error())
			return "", err
		}
		if !verifyOk {
			Error.Println("verify signature error")
			return "", errors.New("verify signature error")
		}

		Info.Println("signedDataStr:", hex.EncodeToString(signedData))

		// append SIGHASH_ALL
		signedData = append(signedData, 0x1)

		scriptSig, err := BTCCombineSignatureAndRedeemScript(signedData, redeemScriptBytes)
		if err != nil {
			Error.Println("BTCCombineSignatureAndRedeemScript fail:", err.Error())
			return "", err
		}
		signedDataList[i] = scriptSig
	}

	for i := 0; i < len(trx.Vin); i++ {
		trx.Vin[i].ScriptSig.SetScriptBytes(signedDataList[i])
	}

	trxSigStr, err := BTCPackRawTransaction(*trx)
	if err != nil {
		Error.Println("BTCPackRawTransaction fail:", err.Error())
		return "", err
	}

	fmt.Println("rawTrxSignedStr:", trxSigStr)

	return trxSigStr, nil
}
