package alipay

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/google/go-querystring/query"
)

const (
	mediaTypeXWWWFormURLEncodedCharsetUTF8 = "application/x-www-form-urlencoded;charset=utf-8"

	GatewayDevelopment = "https://openapi.alipaydev.com/gateway.do"
	GatewayProduction  = "https://openapi.alipay.com/gateway.do"

	format          = "JSON"
	charset         = "utf-8"
	signType        = "RSA2"
	version         = "1.0"
	timestampLayout = "2006-01-02 15:04:05"

	successCode = "10000"
)

type ResponseCommon struct {
	Code    string `json:"code"`
	Msg     string `json:"msg"`
	SubCode string `json:"sub_code"`
	SubMsg  string `json:"sub_msg"`
}

func (r ResponseCommon) IsSuccess() bool {
	return successCode == r.Code
}

type Request interface {
	Method() string
}

type Client struct {
	gateway         string
	appID           string
	appPrivateKey   *rsa.PrivateKey
	alipayPublicKey *rsa.PublicKey
}

type Options struct {
	Gateway         string
	AppID           string
	AppPrivateKey   string
	AlipayPublicKey string
}

func NewClient(options Options) (client *Client, err error) {
	client = &Client{
		gateway: options.Gateway,
		appID:   options.AppID,
	}
	client.appPrivateKey, err = parseRSAPrivateKey(options.AppPrivateKey)
	if err != nil {
		return
	}
	client.alipayPublicKey, err = parseRSAPublicKey(options.AlipayPublicKey)
	return
}

func (c *Client) VerifyNotification(request *http.Request) (err error) {
	err = request.ParseForm()
	if err != nil {
		return
	}
	appID := request.Form.Get("app_id")
	if appID != c.appID {
		err = errors.New("alipay: app_id does not match")
		return
	}
	params := make(url.Values)
	for key := range request.Form {
		params.Set(key, request.Form.Get(key))
	}
	params.Del("sign")
	params.Del("sign_type")

	data := c.buildStringToSign(params)
	sign := request.Form.Get("sign")
	err = c.verify(data, sign)
	return
}

func (c *Client) Execute(request Request, response interface{}) (err error) {
	params, err := c.buildRequestParams(request)
	if err != nil {
		return
	}
	resp, err := http.Post(c.gateway, mediaTypeXWWWFormURLEncodedCharsetUTF8, strings.NewReader(params.Encode()))
	if err != nil {
		return
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	if !json.Valid(body) {
		err = errors.New("alipay: body is not a valid JSON encoding")
		return
	}
	rawMap := make(map[string]json.RawMessage)
	err = json.Unmarshal(body, &rawMap)
	if err != nil {
		return
	}
	respKey := strings.ReplaceAll(request.Method(), ".", "_") + "_response"
	rawResp := rawMap[respKey]
	rawSign := rawMap["sign"]

	if !json.Valid(rawResp) {
		err = fmt.Errorf("alipay: %s is not a valid JSON encoding", respKey)
		return
	}
	if response != nil {
		err = json.Unmarshal(rawResp, response)
		if err != nil {
			return
		}
	}
	if !json.Valid(rawSign) {
		err = errors.New("alipay: sign is not a valid JSON encoding")
		return
	}
	var sign string
	err = json.Unmarshal(rawSign, &sign)
	if err != nil {
		return
	}
	err = c.verify(string(rawResp), sign)
	return
}

func (c *Client) SDKExecute(request Request) (query string, err error) {
	params, err := c.buildRequestParams(request)
	if err != nil {
		return
	}
	query = params.Encode()
	return
}

func (c *Client) PageExecute(request Request) (url string, err error) {
	params, err := c.buildRequestParams(request)
	if err != nil {
		return
	}
	url = c.gateway + "?" + params.Encode()
	return
}

func (c *Client) buildRequestParams(request Request) (params url.Values, err error) {
	params, err = query.Values(request)
	if err != nil {
		return
	}
	params.Set("app_id", c.appID)
	params.Set("format", format)
	params.Set("charset", charset)
	params.Set("sign_type", signType)
	params.Set("version", version)

	method := request.Method()
	params.Set("method", method)

	timestamp := time.Now().Format(timestampLayout)
	params.Set("timestamp", timestamp)

	bizContent, err := c.encodeBizContent(request)
	if err != nil {
		return
	}
	params.Set("biz_content", bizContent)

	params = c.cleanRequestParams(params)
	stringToSign := c.buildStringToSign(params)
	sign, err := c.sign(stringToSign)
	if err != nil {
		return
	}
	params.Set("sign", sign)
	return
}

func (c *Client) buildStringToSign(params url.Values) string {
	keys := make(sort.StringSlice, 0, len(params))
	for key := range params {
		keys = append(keys, key)
	}
	keys.Sort()

	var buf strings.Builder
	for i, key := range keys {
		if i > 0 {
			buf.WriteByte('&')
		}
		buf.WriteString(key)
		buf.WriteByte('=')
		buf.WriteString(params.Get(key))
	}
	return buf.String()
}

func (c *Client) cleanRequestParams(params url.Values) url.Values {
	cleanedParams := make(url.Values)
	for oldKey := range params {
		oldVal := params.Get(oldKey)
		newKey := strings.TrimSpace(oldKey)
		newVal := strings.TrimSpace(oldVal)
		if newKey != "sign" && newKey != "" && newVal != "" {
			cleanedParams.Set(newKey, newVal)
		}
	}
	return cleanedParams
}

func (c *Client) encodeBizContent(request Request) (bizContent string, err error) {
	var buf strings.Builder
	err = json.NewEncoder(&buf).Encode(request)
	bizContent = buf.String()
	return
}

func (c *Client) verify(data, sign string) (err error) {
	sig, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return
	}
	sum := sha256.Sum256([]byte(data))
	err = rsa.VerifyPKCS1v15(c.alipayPublicKey, crypto.SHA256, sum[:], sig)
	return
}

func (c *Client) sign(data string) (sign string, err error) {
	sum := sha256.Sum256([]byte(data))
	sig, err := rsa.SignPKCS1v15(rand.Reader, c.appPrivateKey, crypto.SHA256, sum[:])
	if err != nil {
		return
	}
	sign = base64.StdEncoding.EncodeToString(sig)
	return
}
