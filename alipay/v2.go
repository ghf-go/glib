package alipay

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"

	gcrypto "github.com/ghf-go/glib/gcrypto"
	"github.com/ghf-go/glib/gnet"
)

// https://opendocs.alipay.com/open/cd12c885_alipay.trade.app.pay?pathHash=ab686e33&ref=api&scene=20
func (c *Client) appPayV2(orderSn string, amount uint64, subject string) (string, error) {
	data := c.newReqDataV2("alipay.trade.app.pay", map[string]any{
		"out_trade_no": orderSn,
		"total_amount": fmt.Sprintf("%0.2f", float64(amount)/100),
		"subject":      subject,
		"notify_url":   c.notifyUrl,
	})
	return c.buildPostData(data)
}

// h5支付
func (c *Client) h5PayV2(orderSn, return_url string, amount uint64, subject string) (string, error) {
	data := c.newReqDataV2("alipay.trade.wap.pay", map[string]any{
		"out_trade_no": orderSn,
		"total_amount": fmt.Sprintf("%0.2f", float64(amount)/100),
		"subject":      subject,
		"product_code": "QUICK_WAP_WAY",
		"notify_url":   c.notifyUrl,
	})
	data["return_url"] = return_url
	q, e := c.buildPostData(data)
	if e != nil {
		return "", e
	}
	return c.getApiHost() + "?" + q, nil
}

// pc支持
func (c *Client) pagePayV2(orderSn, return_url string, amount uint64, subject string) (string, error) {
	data := c.newReqDataV2("alipay.trade.page.pay", map[string]any{
		"out_trade_no": orderSn,
		"total_amount": fmt.Sprintf("%0.2f", float64(amount)/100),
		"product_code": "FAST_INSTANT_TRADE_PAY",
		"subject":      subject,
		"notify_url":   c.notifyUrl,
	})
	data["return_url"] = return_url
	q, e := c.buildPostData(data)
	if e != nil {
		return "", e
	}
	return c.getApiHost() + "?" + q, nil
}

// https://opendocs.alipay.com/open/6c0cdd7d_alipay.trade.refund?pathHash=4081e89c&ref=api&scene=common
// 退款
func (c *Client) refundV2(orderSn string, amount uint64) (*RefundV2Response, error) {
	data := c.newReqDataV2("alipay.trade.refund", map[string]any{
		"out_trade_no":  orderSn,
		"refund_amount": float64(amount) / 100,
	})
	ret := &V2Response{}
	if e := c.sendV2(data, ret); e != nil {
		return nil, e
	}
	return ret.RefundResp, nil
}

// https://opendocs.alipay.com/open/ce0b4954_alipay.trade.close?pathHash=7b0fdae1&ref=api&scene=common
// 关闭订单
func (c *Client) closeV2(orderSn string) (*CloseV2Response, error) {
	data := c.newReqDataV2("alipay.trade.close", map[string]any{
		"out_trade_no": orderSn,
	})
	ret := &V2Response{}
	if e := c.sendV2(data, ret); e != nil {
		return nil, e
	}
	return ret.CloseResp, nil
}

// https://opendocs.alipay.com/open/82ea786a_alipay.trade.query?pathHash=0745ecea&ref=api&scene=23
// 订单查询
func (c *Client) queryV2(orderSn string) (*OrderQueryV2Response, error) {
	data := c.newReqDataV2("alipay.trade.query", map[string]any{
		"out_trade_no": orderSn,
	})
	ret := &V2Response{}
	if e := c.sendV2(data, ret); e != nil {
		return nil, e
	}
	return ret.OrderQueryResp, nil
}

// https://opendocs.alipay.com/open/8c776df6_alipay.trade.fastpay.refund.query?pathHash=fb6e1894&ref=api&scene=common
// 退款查询
func (c *Client) refundQueryV2(orderSn string) (*RefundQueryV2Response, error) {
	data := c.newReqDataV2("alipay.trade.fastpay.refund.query", map[string]any{
		"out_trade_no": orderSn,
	})
	ret := &V2Response{}
	if e := c.sendV2(data, ret); e != nil {
		return nil, e
	}
	return ret.RefundQueryResp, nil
}

// 支付完成通知
func (c *Client) notifyV2() {}

// 签名
func (c *Client) signV2(data map[string]string) error {
	if c.appCertSn != "" { //密码验签
		data["app_cert_sn"] = c.appCertSn
		data["alipay_root_cert_sn"] = c.alipayRootCertSn
	}
	keys := []string{}
	for k, _ := range data {
		keys = append(keys, k)
	}
	ks := sort.StringSlice(keys)
	ks.Sort()
	kvs := []string{}
	for _, k := range ks {
		if v, _ := data[k]; v != "" {
			kvs = append(kvs, fmt.Sprintf("%s=%s", k, v))
		}
	}
	str := strings.Join(kvs, "&")
	sign, e := gcrypto.RsaSign(c.appPrivateKey, []byte(str))
	if e != nil {
		return fmt.Errorf("sign is error")
	}
	data["sign"] = base64.RawStdEncoding.EncodeToString(sign)
	return nil
}

// 生产提交的内容
func (c *Client) buildPostData(data map[string]string) (string, error) {
	if e := c.signV2(data); e != nil {
		return "", e
	}
	kvs := []string{}
	for k, v := range data {
		kvs = append(kvs, fmt.Sprintf("%s=%s", url.QueryEscape(k), url.QueryEscape(v)))
	}
	return strings.Join(kvs, "&"), nil
}

// 构建请求数据
func (c *Client) newReqDataV2(method string, data map[string]any) map[string]string {
	bz, _ := json.Marshal(data)
	ret := map[string]string{
		"method":    method,
		"app_id":    c.appID,
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
		"format":    "json",
		"version":   "1.0",
		"charset":   "UTF-8",
		"sign_type": "RSA2",
		// "notify_url":  c.notifyUrl,
		"biz_content": string(bz),
	}
	if c.appCertSn != "" {
		ret["app_cert_sn"] = c.appCertSn
		ret["alipay_root_cert_sn"] = c.alipayRootCertSn
	}
	return ret
}

// 发送请求
func (c *Client) sendV2(data map[string]string, out any) error {
	pdata, e := c.buildPostData(data)
	if e != nil {
		return e
	}
	return gnet.PostJson(c.getApiHost(), "", pdata, out)
}
