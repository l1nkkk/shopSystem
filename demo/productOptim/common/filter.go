package common

import (
	"net/http"
)

// l1nkkk: 有点类似与中间件，只不过把原本集成到后端服务中的中间件服务，独立了出来

// 声明一个新的数据类型（函数类型）
type FilterHandle func(rw http.ResponseWriter, req *http.Request) error

// Filter 拦截器结构体
type Filter struct {
	// filterMap 用来存储需要拦截的URI
	filterMap map[string]FilterHandle
}

// Filter 初始化函数
func NewFilter() *Filter {
	return &Filter{filterMap: make(map[string]FilterHandle)}
}

// RegisterFilterUri 注册拦截器
func (f *Filter) RegisterFilterUri(uri string, handler FilterHandle) {
	f.filterMap[uri] = handler
}

// GetFilterHandle 根据Uri获取对应的handle
func (f *Filter) GetFilterHandle(uri string) FilterHandle {
	return f.filterMap[uri]
}

// 声明一个新的数据类型
type WebHandle func(rw http.ResponseWriter, req *http.Request)

// Handle 将拦截器与正常业务绑定，返回新的处理逻辑
func (f *Filter) Handle(webHandle WebHandle) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		// 1.执行拦截业务逻辑
		for path, handle := range f.filterMap {
			if path == r.RequestURI {
				err := handle(rw, r)
				if err != nil {
					// 进入到这说明被拦截了
					rw.Write([]byte(err.Error()))
					return
				}
				break
			}
		}
		// 2.执行正常注册的函数(正常业务逻辑）
		webHandle(rw, r)
	}
}
