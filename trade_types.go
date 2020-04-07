package alipay

// 交易查询响应
type TradeQueryResponse struct {
	ResponseCommon
	TradeNo      string `json:"trade_no"`
	OutTradeNo   string `json:"out_trade_no"`
	BuyerLogonID string `json:"buyer_logon_id"`
	TradeStatus  string `json:"trade_status"`
	TotalAmount  string `json:"total_amount"`
	BuyerUserID  string `json:"buyer_user_id"`
}

// 交易查询请求
type TradeQueryRequest struct {
	OutTradeNo string `url:"-" json:"out_trade_no,omitempty"`
	TradeNo    string `url:"-" json:"trade_no,omitempty"`
}

func (TradeQueryRequest) Method() string {
	return "alipay.trade.query"
}

// 条码支付响应
type TradePayResponse struct {
	ResponseCommon
	TradeNo       string `json:"trade_no"`
	OutTradeNo    string `json:"out_trade_no"`
	BuyerLogonID  string `json:"buyer_logon_id"`
	TotalAmount   string `json:"total_amount"`
	ReceiptAmount string `json:"receipt_amount"`
	GMTPayment    string `json:"gmt_payment"`
	BuyerUserID   string `json:"buyer_user_id"`
}

// 条码支付请求
type TradePayRequest struct {
	NotifyURL   string `json:"-" url:"notify_url,omitempty"`
	OutTradeNo  string `url:"-" json:"out_trade_no,omitempty"`
	Scene       string `url:"-" json:"scene,omitempty"`
	AuthCode    string `url:"-" json:"auth_code,omitempty"`
	Subject     string `url:"-" json:"subject,omitempty"`
	TotalAmount string `url:"-" json:"total_amount,omitempty"`
}

func (TradePayRequest) Method() string {
	return "alipay.trade.pay"
}

// 扫码支付响应
type TradePrecreateResponse struct {
	ResponseCommon
	OutTradeNo string `json:"out_trade_no"`
	QRCode     string `json:"qr_code"`
}

// 扫码支付请求
type TradePrecreateRequest struct {
	NotifyURL   string `json:"-" url:"notify_url,omitempty"`
	OutTradeNo  string `url:"-" json:"out_trade_no,omitempty"`
	TotalAmount string `url:"-" json:"total_amount,omitempty"`
	Subject     string `url:"-" json:"subject,omitempty"`
}

func (TradePrecreateRequest) Method() string {
	return "alipay.trade.precreate"
}

// App 支付请求
type TradeAppPayRequest struct {
	ReturnURL   string `json:"-" url:"return_url,omitempty"`
	NotifyURL   string `json:"-" url:"notify_url,omitempty"`
	TotalAmount string `url:"-" json:"total_amount,omitempty"`
	Subject     string `url:"-" json:"subject,omitempty"`
	OutTradeNo  string `url:"-" json:"out_trade_no,omitempty"`
}

func (TradeAppPayRequest) Method() string {
	return "alipay.trade.app.pay"
}

// 手机网站支付请求
type TradeWapPayRequest struct {
	ReturnURL   string `json:"-" url:"return_url,omitempty"`
	NotifyURL   string `json:"-" url:"notify_url,omitempty"`
	Subject     string `url:"-" json:"subject,omitempty"`
	OutTradeNo  string `url:"-" json:"out_trade_no,omitempty"`
	TotalAmount string `url:"-" json:"total_amount,omitempty"`
	QuitURL     string `url:"-" json:"quit_url,omitempty"`
	ProductCode string `url:"-" json:"product_code,omitempty"`
}

func (TradeWapPayRequest) Method() string {
	return "alipay.trade.wap.pay"
}

// 电脑网站支付请求
type TradePagePayRequest struct {
	ReturnURL   string `json:"-" url:"return_url,omitempty"`
	NotifyURL   string `json:"-" url:"notify_url,omitempty"`
	OutTradeNo  string `url:"-" json:"out_trade_no,omitempty"`
	ProductCode string `url:"-" json:"product_code,omitempty"`
	TotalAmount string `url:"-" json:"total_amount,omitempty"`
	Subject     string `url:"-" json:"subject,omitempty"`
}

func (TradePagePayRequest) Method() string {
	return "alipay.trade.page.pay"
}
