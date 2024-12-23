package alipay

// https://opendocs.alipay.com/open/cd12c885_alipay.trade.app.pay?pathHash=ab686e33&ref=api&scene=20
func (c *Client) appPayV2()  {}
func (c *Client) h5PayV2()   {}
func (c *Client) pagePayV2() {}

// https://opendocs.alipay.com/open/6c0cdd7d_alipay.trade.refund?pathHash=4081e89c&ref=api&scene=common
// 退款
func (c *Client) refundV2() {}

// https://opendocs.alipay.com/open/ce0b4954_alipay.trade.close?pathHash=7b0fdae1&ref=api&scene=common
// 关闭订单
func (c *Client) closeV2() {}

// https://opendocs.alipay.com/open/82ea786a_alipay.trade.query?pathHash=0745ecea&ref=api&scene=23
// 订单查询
func (c *Client) queryV2() {}

// https://opendocs.alipay.com/open/8c776df6_alipay.trade.fastpay.refund.query?pathHash=fb6e1894&ref=api&scene=common
// 退款查询
func (c *Client) refundQueryV2() {}

// 支付完成通知
func (c *Client) notifyV2() {}
