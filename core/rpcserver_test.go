package core

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestCreateHTTPClient(t *testing.T) {
	assert.IsType(t, &http.Client{}, createHTTPClient())
}

func TestGetErrorResponseBytes(t *testing.T) {
	bts := getErrorResponseBytes(1, "test msg")

	assert.IsType(t, "{\"error\":{\"code\":-32602,\"message\":\"test msg\"},\"id\":1,\"jsonrpc\":\"2.0\"}", string(bts))
}

func TestBuildRunningConfigFromConfigNAIVE(t *testing.T) {
	var testConfigStr1 = `{
		"_upstreams": "support http, https, ws, wss",
		"upstreams": [
		  "https://ropsten.infura.io/v3/83438c4dcf834ceb8944162688749707"
		],
		
		"_strategy": "support NAIVE, RACE, FALLBACK",
		"strategy": "NAIVE",
	  
		"_methodLimitationEnabled": "limit or not",
		"methodLimitationEnabled": true,
	  
		"_allowedMethods": "can be ignored when set methodLimitationEnabled false",
		"allowedMethods": ["eth_blockNumber"],
	  
		"_contractWhitelist": "can be ignored when set methodLimitationEnabled false",
		"contractWhitelist": []
	  }`

	config := &Config{}

	err := json.Unmarshal([]byte(testConfigStr1), config)

	rcfg1, err := buildRunningConfigFromConfig(context.Background(), config)

	if err != nil {
		logrus.Fatal(err)
	}

	assert.Equal(t, true, rcfg1.MethodLimitationEnabled)

	var testConfigStr2 = `{
		"_upstreams": "support http, https, ws, wss",
		"upstreams": [
		  "https://ropsten.infura.io/v3/83438c4dcf834ceb8944162688749707",
		  "https://test1.com"
		],
		
		"_strategy": "support NAIVE, RACE, FALLBACK",
		"strategy": "NAIVE",
	  
		"_methodLimitationEnabled": "limit or not",
		"methodLimitationEnabled": true,
	  
		"_allowedMethods": "can be ignored when set methodLimitationEnabled false",
		"allowedMethods": ["eth_blockNumber"],
	  
		"_contractWhitelist": "can be ignored when set methodLimitationEnabled false",
		"contractWhitelist": []
	  }`

	err = json.Unmarshal([]byte(testConfigStr2), config)

	assert.Panics(t, func() { buildRunningConfigFromConfig(context.Background(), config) })
}

func TestBuildRunningConfigFromConfigRACE(t *testing.T) {
	var testConfigStr1 = `{
		"_upstreams": "support http, https, ws, wss",
		"upstreams": [
		  "https://ropsten.infura.io/v3/83438c4dcf834ceb8944162688749707",
		  "https://test1.com"
		],
		
		"_strategy": "support NAIVE, RACE, FALLBACK",
		"strategy": "RACE",
	  
		"_methodLimitationEnabled": "limit or not",
		"methodLimitationEnabled": true,
	  
		"_allowedMethods": "can be ignored when set methodLimitationEnabled false",
		"allowedMethods": ["eth_blockNumber"],
	  
		"_contractWhitelist": "can be ignored when set methodLimitationEnabled false",
		"contractWhitelist": []
	  }`

	config := &Config{}

	err := json.Unmarshal([]byte(testConfigStr1), config)

	rcfg1, err := buildRunningConfigFromConfig(context.Background(), config)

	if err != nil {
		logrus.Fatal(err)
	}

	assert.Equal(t, true, rcfg1.MethodLimitationEnabled)

	var testConfigStr2 = `{
		"_upstreams": "support http, https, ws, wss",
		"upstreams": [
		  "https://ropsten.infura.io/v3/83438c4dcf834ceb8944162688749707"
		],
		
		"_strategy": "support NAIVE, RACE, FALLBACK",
		"strategy": "RACE",
	  
		"_methodLimitationEnabled": "limit or not",
		"methodLimitationEnabled": true,
	  
		"_allowedMethods": "can be ignored when set methodLimitationEnabled false",
		"allowedMethods": ["eth_blockNumber"],
	  
		"_contractWhitelist": "can be ignored when set methodLimitationEnabled false",
		"contractWhitelist": []
	  }`

	err = json.Unmarshal([]byte(testConfigStr2), config)

	assert.Panics(t, func() { buildRunningConfigFromConfig(context.Background(), config) })
}

func TestBuildRunningConfigFromConfigFALLBACK(t *testing.T) {
	var testConfigStr1 = `{
		"_upstreams": "support http, https, ws, wss",
		"upstreams": [
		  "https://ropsten.infura.io/v3/83438c4dcf834ceb8944162688749707",
		  "https://test1.com"
		],
		
		"_strategy": "support NAIVE, RACE, FALLBACK",
		"strategy": "FALLBACK",
	  
		"_methodLimitationEnabled": "limit or not",
		"methodLimitationEnabled": true,
	  
		"_allowedMethods": "can be ignored when set methodLimitationEnabled false",
		"allowedMethods": ["eth_blockNumber"],
	  
		"_contractWhitelist": "can be ignored when set methodLimitationEnabled false",
		"contractWhitelist": []
	  }`

	config := &Config{}

	err := json.Unmarshal([]byte(testConfigStr1), config)

	rcfg1, err := buildRunningConfigFromConfig(context.Background(), config)

	if err != nil {
		logrus.Fatal(err)
	}

	assert.Equal(t, true, rcfg1.MethodLimitationEnabled)

	var testConfigStr2 = `{
		"_upstreams": "support http, https, ws, wss",
		"upstreams": [
		  "https://ropsten.infura.io/v3/83438c4dcf834ceb8944162688749707"
		],
		
		"_strategy": "support NAIVE, RACE, FALLBACK",
		"strategy": "FALLBACK",
	  
		"_methodLimitationEnabled": "limit or not",
		"methodLimitationEnabled": true,
	  
		"_allowedMethods": "can be ignored when set methodLimitationEnabled false",
		"allowedMethods": ["eth_blockNumber"],
	  
		"_contractWhitelist": "can be ignored when set methodLimitationEnabled false",
		"contractWhitelist": []
	  }`

	err = json.Unmarshal([]byte(testConfigStr2), config)

	assert.Panics(t, func() { buildRunningConfigFromConfig(context.Background(), config) })
}
