package alipay

// 支付
func (c *Client) AppPay() {
	if c.isV3 {
		c.appPayV3()
	} else {
		c.appPayV2()
	}
}
func (c *Client) H5Pay() {
	if c.isV3 {
		c.h5PayV3()
	} else {
		c.h5PayV2()
	}
}
func (c *Client) PagePay() {
	if c.isV3 {
		c.pagePayV3()
	} else {
		c.pagePayV2()
	}
}

// 退款
func (c *Client) Refund() {
	if c.isV3 {
		c.refundV3()
	} else {
		c.refundV2()
	}
}

// 关闭订单
func (c *Client) Close() {
	if c.isV3 {
		c.closeV3()
	} else {
		c.closeV2()
	}
}

// 订单查询
func (c *Client) Query() {
	if c.isV3 {
		c.queryV3()
	} else {
		c.queryV2()
	}
}

// 退款查询
func (c *Client) RefundQuery() {
	if c.isV3 {
		c.refundQueryV3()
	} else {
		c.refundQueryV2()
	}
}

// 支付完成通知
func (c *Client) Notify() {
	if c.isV3 {
		c.notifyV3()
	} else {
		c.notifyV2()
	}
}
