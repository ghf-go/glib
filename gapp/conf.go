package gapp

import (
	"strings"
)

type conf struct {
	App     *appConfig        `yaml:"app"`
	Dbs     map[string]string `yaml:"dbs"`
	Caches  map[string]string `yaml:"cache"`
	Log     string            `yaml:"LogConfig"`
	Stmp    map[string]string `yaml:"smtp"`
	Payment PaymentConfig     `yaml:"payment"`
	Lang    langConf          `yaml:"lang"`
	Storage map[string]string `yaml:"storage"`
	Oauth   oauthConf         `yaml:"oauth"`
	Meta    MetaConf          `yaml:"meta"`
}
type appConfig struct {
	Port  int  `yaml:"port"`
	Debug bool `yaml:"debug"`
}

type oauthConf struct {
	WeChat oauthWeChatConf `yaml:"wechat"`
}
type oauthWeChatConf struct {
	App           *oauthWeChatItemConf `yaml:"app"`
	MiniProgram   *oauthWeChatItemConf `yaml:"mini_program"`
	PublicAccount *oauthWeChatItemConf `yaml:"public_account"`
}

type oauthWeChatItemConf struct {
	AppID     string `yaml:"app_id"`
	AppSecret string `yaml:"app_secret"`
}
type PaymentConfig struct {
	WexPay *WeChatPaymentConfig `yaml:"wechat"`
	AliPay *AliPaymentConfig    `yaml:"alipay"`
}

// 微信支付
type WeChatPaymentConfig struct {
	AppID           string `yaml:"app_id"`
	MchID           string `yaml:"mch_id"`
	MchIdNumber     string `yaml:"mch_id_num"`
	MchApiV3Key     string `yaml:"mch_api_v3_key"`
	PriviteKeyPem   string `yaml:"private_key_pem"`
	NotifyURL       string `yaml:"notify_url"`
	RefundNotifyURL string `yaml:"refund_notify_url"`
}

// 支付宝支付
type AliPaymentConfig struct {
	AppID         string `yaml:"app_id"`
	NotifyURL     string `yaml:"notify_url"`
	AppPublicPem  string `yaml:"app_public_pem"`
	AppPrivatePem string `yaml:"app_private_pem"`
	AliPublicPem  string `yaml:"ali_public_pem"`
	RootPem       string `yaml:"root_pem"`
	AliGateWay    string `yaml:"gateway"`
	AliPrivateKey string `yaml:"alipay_private_key"`
}

type MetaConf map[string]any

func (m MetaConf) Get(key string) any {
	if r, ok := m[key]; ok {
		return r
	}
	return nil
}

type langConf map[string]map[string]string

// 生成多语言
func (c langConf) Lang(key, lang string) string {
	lang = strings.ToLower(lang)
	if r, ok := c[lang]; ok {
		if ret, isok := r[key]; isok {
			return ret
		}
		return ""
	}
	return c.Lang(key, "us-en")
}
