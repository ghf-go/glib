# 支付宝支付

## 概述

支付宝支付库
## 安装

```bash
go get -u github.com/ghf-go/glib/alipay
```

## 文档
### H5下单
```go
func Pay(amount int, subject string, outTradeNo string, notifyUrl string, returnUrl string) (string, error)
```
### App下单
```go
func AppPay(amount int, subject string, outTradeNo string, notifyUrl string, returnUrl string) (string, error)
### Page下单
```go
func PagePay(amount int, subject string, outTradeNo string, notifyUrl string, returnUrl string) (string, error)
```
### 验签
```go
func VerifySign(params map[string]string) bool
```
### 退款
```go
func Refund(tradeNo string, refundAmount int, outRequestNo string) (string, error)
```
### 退款查询
```go
func RefundQuery(tradeNo string, outRequestNo string) (string, error)
```
### 关闭订单
```go
func Close(tradeNo string) (string, error)
```
### 查询订单
```go
func Query(tradeNo string) (string, error)
```
## 示例
```go
package main

import (

)   
func main() {
}