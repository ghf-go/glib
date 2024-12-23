package alipay

// app支付
func (c *Client) AppPay(orderSn string, amount uint64, subject string) (string, error) {
	return c.appPayV2(orderSn, amount, subject)
}

// h5支付
func (c *Client) H5Pay(orderSn, return_url string, amount uint64, subject string) (string, error) {
	return c.h5PayV2(orderSn, return_url, amount, subject)
}

// pc支付
func (c *Client) PagePay(orderSn, return_url string, amount uint64, subject string) (string, error) {
	return c.pagePayV2(orderSn, return_url, amount, subject)
}

// 退款
func (c *Client) Refund(orderSn string, amount uint64) (*RefundV2Response, error) {
	return c.refundV2(orderSn, amount)
}

// 关闭订单
func (c *Client) Close(orderSn string) (*CloseV2Response, error) {
	return c.closeV2(orderSn)
}

// 订单查询
func (c *Client) Query(orderSn string) (*OrderQueryV2Response, error) {
	return c.queryV2(orderSn)
}

// 退款查询
func (c *Client) RefundQuery(orderSn string) (*RefundQueryV2Response, error) {
	return c.refundQueryV2(orderSn)
}

// 支付完成通知
// https://opendocs.alipay.com/open/00dn78?pathHash=fef00e6d#%E5%BC%82%E6%AD%A5%E8%BF%94%E5%9B%9E%E7%BB%93%E6%9E%9C%E9%AA%8C%E7%AD%BE
func (c *Client) Notify() {
	c.notifyV2()
}
