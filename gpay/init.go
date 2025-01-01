package gpay

// 通知内容
type PayNotifyParam struct {
	OrderSn   string
	TradeId   string
	IsSuccess bool
	Status    string
}

// 通知回掉
type PayNotifyHandle func(PayNotifyParam) bool

type Pay interface {
	CreateOrderApp(ordersn string, amount uint64, orderDesc ...string) map[string]any
	CreateOrderH5(ordersn string, amount uint64, orderDesc ...string) map[string]any
	CreateOrderNative(ordersn string, amount uint64, orderDesc ...string) map[string]any
	CreateOrderJsapi(ordersn string, amount uint64, orderDesc ...string) map[string]any
	CancelOrder(ordersn string) bool
	RefundOrder(ordersn string, amount uint64)
	NotifyOrder(call PayNotifyHandle) bool
}
