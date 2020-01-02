package core

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rlp"
	"strings"
)

type RequestData struct {
	JsonRpc string        `json:"jsonrpc"`
	ID      int64         `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

//eth_call
//eth_estimateGas
//eth_getLogs
//eth_getBalance
//eth_getCode
//eth_getStorageAt
//eth_getTransactionCount

var DecodeError = fmt.Errorf("decode error")
var Denied = fmt.Errorf("not allowed method")

func isAllowedMethod(method string) bool {
	return currentRunningConfig.allowedMethods[method]
}

func inWhitelist(contractAddress string) bool {
	return currentRunningConfig.allowedCallContracts[strings.ToLower(contractAddress)]
}

func isValidCall(req *RequestData) (err error) {
	defer func() {
		if er := recover(); er != nil {
			err = DecodeError
		}
	}()

	if !isAllowedMethod(req.Method) {
		return Denied
	}

	if req.Method == "eth_getBalance" ||
		req.Method == "eth_getTransactionReceipt" {
		return nil
	}

	if req.Method == "eth_call" || req.Method == "eth_estimateGas" {
		to := req.Params[0].(map[string]interface{})["to"].(string)

		if !inWhitelist(to) {
			return Denied
		}

		return nil
	}

	if req.Method == "eth_sendRawTransaction" {
		// 0. nonce
		// 1. gasPrice
		// 2. gasLimit
		// 3. to
		// 4. value
		// 5. data
		// 6. signature
		var fields []interface{}

		data := req.Params[0].(string)
		bts, _ := hexutil.Decode(data)
		err = rlp.DecodeBytes(bts, &fields)

		if err != nil {
			return DecodeError
		}

		if !inWhitelist(fields[3].(string)) {
			return Denied
		}

		return nil
	}

	return Denied
}
