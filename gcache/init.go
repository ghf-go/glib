package gcache

type Cache interface {
	Get(key string, defval string) string
	GetAll(key ...string) map[string]string
	GetObj(key string, out any) error
	GetAllObj(data map[string]any)
	Set(key, val string, timeOut ...int) error
	SetObj(key string, obj any, timeOut ...int) error
	SetNx(key, val string, timeOut ...int) error
	SetObjNx(key string, obj any, timeOut ...int) error
	Incr(key string, step ...int) int
	Decr(key string, step ...int) int
	Del(key ...string) error
	Flush() error
	Lock(key string, callfunc func(), timeOut ...int) error
}
