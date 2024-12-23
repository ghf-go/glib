package alipay

type baseV2Response struct {
	Code    string `json:"code"`
	Msg     string `json:"msg"`
	SubCode string `json:"sub_code"`
	SubMsg  string `json:"sub_msg"`
}

// 是否操作成功
func (c baseV2Response) IsSuccess() bool {
	return c.Code == "10000"
}

// 关闭订单的返回
type CloseV2Response struct {
	baseV2Response
	TradeNo    string `json:"trade_no"`
	OutTradeNo string `json:"out_trade_no"`
}
type OrderQueryV2Response struct {
	baseV2Response
	TradeNo      string  `json:"trade_no"`
	OutTradeNo   string  `json:"out_trade_no"`
	BuyerLogonID string  `json:"buyer_logon_id"`
	TradeStatus  string  `json:"trade_status"`
	TotalAmount  float64 `json:"total_amount"`
	SendPayDate  string  `json:"send_pay_date"`
}
type RefundV2Response struct {
	baseV2Response
	TradeNo      string  `json:"trade_no"`
	OutTradeNo   string  `json:"out_trade_no"`
	BuyerLogonID string  `json:"buyer_logon_id"`
	RefundFee    float64 `json:"refund_fee"`
}
type RefundQueryV2Response struct {
	baseV2Response
	TradeNo      string  `json:"trade_no"`
	OutTradeNo   string  `json:"out_trade_no"`
	BuyerLogonID string  `json:"buyer_logon_id"`
	RefundStatus string  `json:"refund_status"`
	TotalAmount  float64 `json:"total_amount"`
	RefundAmount float64 `json:"refund_amount"`
}

// 接口返回
type V2Response struct {
	CloseResp       *CloseV2Response       `json:"alipay_trade_close_response"`                //关闭订单返回
	OrderQueryResp  *OrderQueryV2Response  `json:"alipay_trade_query_response"`                //查询订单返回
	RefundResp      *RefundV2Response      `json:"alipay_trade_refund_response"`               //退款返回
	RefundQueryResp *RefundQueryV2Response `json:"alipay_trade_fastpay_refund_query_response"` //查询订单返回

}
