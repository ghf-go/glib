package gapp

import "strings"

type WebRouter struct {
	groups     map[string]*WebRouter
	methods    map[string][]Handle
	middlefunc []Handle
	err404     Handle
}

func newRootRouter(args ...Handle) *WebRouter {
	args = append([]Handle{func(c *Content) {
		c.Next()
		c.Flush()
	}}, args...)
	return &WebRouter{
		groups:     map[string]*WebRouter{},
		methods:    map[string][]Handle{},
		middlefunc: args,
		err404: func(c *Content) {
			c.FailJson(404, "接口不存在")
		},
	}
}

func (r *WebRouter) Error404(err404 Handle) {
	r.err404 = err404
}
func (r *WebRouter) Group(path string, args ...Handle) *WebRouter {
	for strings.HasPrefix(path, "/") {
		path = path[1:]
	}
	path = strings.ToLower(path)
	nr := &WebRouter{
		groups:     map[string]*WebRouter{},
		methods:    map[string][]Handle{},
		middlefunc: append(r.middlefunc, args...),
		err404:     r.err404,
	}
	r.groups[path] = nr
	return nr
}
func (r *WebRouter) Get(path string, hand Handle, args ...Handle) {
	r.addmethods("get", path, hand, args...)
}
func (r *WebRouter) Post(path string, hand Handle, args ...Handle) {
	r.addmethods("post", path, hand, args...)
}
func (r *WebRouter) Delete(path string, hand Handle, args ...Handle) {
	r.addmethods("delete", path, hand, args...)
}
func (r *WebRouter) Put(path string, hand Handle, args ...Handle) {
	r.addmethods("put", path, hand, args...)
}
func (r *WebRouter) Options(path string, hand Handle, args ...Handle) {
	r.addmethods("options", path, hand, args...)
}
func (r *WebRouter) Head(path string, hand Handle, args ...Handle) {
	r.addmethods("head", path, hand, args...)
}
func (r *WebRouter) Any(path string, hand Handle, args ...Handle) {
	r.addmethods("any", path, hand, args...)
}
func (r *WebRouter) addmethods(methodname, path string, hand Handle, args ...Handle) {
	for strings.HasPrefix(path, "/") {
		path = path[1:]
	}
	path = strings.ToLower(path)
	if item, ok := r.groups[path]; ok {
		item.methods[strings.ToUpper(methodname)] = append([]Handle{hand}, args...)
	} else {
		nr := &WebRouter{
			groups:     map[string]*WebRouter{},
			methods:    map[string][]Handle{strings.ToUpper(methodname): append([]Handle{hand}, args...)},
			middlefunc: r.middlefunc,
		}
		r.groups[path] = nr
	}

}
func (r *WebRouter) FindHandle(methodname, path string) (bool, []Handle) {
	for strings.HasPrefix(path, "/") {
		path = path[1:]
	}
	methodname = strings.ToUpper(methodname)
	path = strings.ToLower(path)
	ps := strings.Split(path, "/")
	// fmt.Printf("路由信息 %s - %v\n", path, ps)
	var rt *WebRouter
	error404 := r.err404
	funcs := r.middlefunc
	for _, k := range ps {
		if rt == nil {
			rt = r.findPath(k)
			if rt == nil {
				break
			}
		} else {
			funcs = rt.middlefunc
			error404 = rt.err404

			rt = rt.findPath(k)
			if rt == nil {
				break
			}
		}
	}
	if rt == nil {
		// fmt.Printf("没有找打路由 %d\n", len(append(funcs, error404)))
		return false, append(funcs, error404)
	}
	if ffs, on := rt.methods[methodname]; on {
		return true, append(funcs, ffs...)
	} else if ffs, on := rt.methods["ANY"]; on {
		return true, append(funcs, ffs...)
	}
	// fmt.Printf("没有找打路由 %d", len(append(funcs, error404)))
	return false, append(funcs, error404)
}
func (r *WebRouter) findPath(name string) *WebRouter {
	// fmt.Printf("查找信息 %s -> %v\n", name, r.groups)
	if rr, ok := r.groups[name]; ok {

		return rr
	}
	return nil
}
