package alipay

import (
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	client     *Client
	outTradeNo string
)

func init() {
	var err error
	client, err = newTestClient()
	if err != nil {
		panic(err)
	}
	outTradeNo = strconv.FormatInt(time.Now().UnixNano(), 10)
}

func newTestClient() (*Client, error) {
	options := Options{
		Gateway:         GatewayDevelopment,
		AppID:           os.Getenv("APP_ID"),
		AppPrivateKey:   os.Getenv("APP_PRIVATE_KEY"),
		AlipayPublicKey: os.Getenv("ALIPAY_PUBLIC_KEY"),
	}
	return NewClient(options)
}

func TestClient_Execute(t *testing.T) {
	request := TradePrecreateRequest{
		OutTradeNo:  outTradeNo,
		TotalAmount: "0.01",
		Subject:     "fresh meat",
	}
	var response TradePrecreateResponse
	err := client.Execute(request, &response)
	assert.NoError(t, err)
	assert.True(t, response.IsSuccess())

	t.Logf("%#v", response)
}

func TestClient_SDKExecute(t *testing.T) {
	request := TradeAppPayRequest{
		TotalAmount: "0.01",
		Subject:     "fresh meat",
		OutTradeNo:  outTradeNo,
	}
	query, err := client.SDKExecute(request)
	assert.NoError(t, err)
	assert.NotEmpty(t, query)

	t.Log(query)
}

func TestClient_PageExecute(t *testing.T) {
	request := TradePagePayRequest{
		OutTradeNo:  outTradeNo,
		ProductCode: "FAST_INSTANT_TRADE_PAY",
		TotalAmount: "0.01",
		Subject:     "fresh meat",
	}
	url, err := client.PageExecute(request)
	assert.NoError(t, err)
	assert.NotEmpty(t, url)

	t.Log(url)
}
