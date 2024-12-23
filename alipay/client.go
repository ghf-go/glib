package alipay

// https://opendocs.alipay.com/open/cd12c885_alipay.trade.app.pay?pathHash=ab686e33&ref=api&scene=20
// 支付
func (c *Client) AppPayV2() {}

// https://opendocs.alipay.com/open/6c0cdd7d_alipay.trade.refund?pathHash=4081e89c&ref=api&scene=common
// 退款
func (c *Client) AppRefundV2() {}

// https://opendocs.alipay.com/open/ce0b4954_alipay.trade.close?pathHash=7b0fdae1&ref=api&scene=common
// 关闭订单
func (c *Client) AppCloseV2() {}

// https://opendocs.alipay.com/open/82ea786a_alipay.trade.query?pathHash=0745ecea&ref=api&scene=23
// 订单查询
func (c *Client) AppQueryV2() {}

// https://opendocs.alipay.com/open/8c776df6_alipay.trade.fastpay.refund.query?pathHash=fb6e1894&ref=api&scene=common
// 退款查询
func (c *Client) AppRefundQueryV2() {}

// https://opendocs.alipay.com/open/fca5d17e_alipay.trade.refund.depositback.completed?pathHash=47d46319&ref=api&scene=common
// 退款通知
func (c *Client) RefundNotifyV2() {}

// 支付完成通知
func (c *Client) PayNotifyV2() {}
