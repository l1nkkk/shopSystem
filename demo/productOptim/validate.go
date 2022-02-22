package main

import (
	"fmt"
	"github.com/l1nkkk/shopSystem/demo/productOptim/common"
	"github.com/l1nkkk/shopSystem/demo/productOptim/encrypt"
	"net/http"
	"sync"
	"strconv"
	"io/ioutil"
)

// l1nkkk: 这里连个IP都一样，在add virtual node hash 的适合不会碰撞?
// 设置集群地址，最好内外IP
var hostArray= []string{"127.0.0.1","127.0.0.1"}

var localHost = "127.0.0.1"

var port = "8081"

var hashConsistent *common.Consistent

// AccessControl 分布式存储控制器
type AccessControl struct {
	// 用来存放用户想要存放的信息
	sourcesArray map[int]interface{}
	sync.RWMutex
}

var accessControl = &AccessControl{sourcesArray:make(map[int]interface{})}

// GetNewRecord获取指定的数据
func (m *AccessControl) GetNewRecord(uid int) interface{} {
	m.RWMutex.RLock()
	defer m.RWMutex.RUnlock()
	data:=m.sourcesArray[uid]
	return data
}

// SetNewRecord设置记录
func (m *AccessControl) SetNewRecord(uid int) {
	m.RWMutex.Lock()
	// 注：这里是写死的，只是为了测试
	m.sourcesArray[uid]="hello imooc"
	m.RWMutex.Unlock()
}

// GetDistributedRight 分布式的方式获取uid对应的数据。
// 可能从本机获取，可能充当代理，直接到定位到的节点获取
func (m *AccessControl) GetDistributedRight(req *http.Request) bool {
	// 1. 获取用户UID
	uid ,err := req.Cookie("uid")
	if err !=nil {
		return false
	}

	// 2. 采用一致性hash算法，根据用户ID，判断获取具体机器
	hostRequest,err:=hashConsistent.Get(uid.Value)
	if err !=nil {
		return false
	}

	// 3.判断是否为本机
	if hostRequest == localHost {
		//执行本机数据读取和校验
		return m.GetDataFromMap(uid.Value)
	} else {
		//不是本机充当代理访问数据返回结果
		return GetDataFromOtherMap(hostRequest,req)
	}

}

// GetDataFromMap 获取本机map，并且处理业务逻辑，返回的结果类型为bool类型
func (m *AccessControl) GetDataFromMap(uid string) (isOk bool) {
	uidInt,err := strconv.Atoi(uid)
	if err !=nil {
		return false
	}
	data:=m.GetNewRecord(uidInt)

	//执行逻辑判断
	if data !=nil {
		return true
	}
	return
}

// GetDataFromOtherMap 获取其它节点处理结果
func GetDataFromOtherMap(host string,request *http.Request) bool  {
	// 1. 获取Uid
	uidPre,err := request.Cookie("uid")
	if err !=nil {
		return false
	}
	// 2. 获取sign
	uidSign,err:=request.Cookie("sign")
	if err !=nil {
		return  false
	}

	// 3. 模拟接口访问，
	client :=&http.Client{}
	req,err:= http.NewRequest("GET","http://"+host+":"+port+"/check",nil)
	if err !=nil {
		return false
	}

	// 4. 手动指定，排查多余cookies
	cookieUid :=&http.Cookie{Name:"uid",Value:uidPre.Value,Path:"/"}
	cookieSign :=&http.Cookie{Name:"sign",Value:uidSign.Value,Path:"/"}
	// 5. 添加cookie到模拟的请求中
	req.AddCookie(cookieUid)
	req.AddCookie(cookieSign)

	// 6. 获取返回结果
	response,err :=client.Do(req)
	if err !=nil {
		return false
	}
	body,err:=ioutil.ReadAll(response.Body)
	if err !=nil {
		return false
	}

	// 7. 判断状态
	if response.StatusCode == 200 {
		if string(body) == "true" {
			return true
		} else {
			return false
		}
	}
	return false
}

// Check 执行正常业务逻辑，这个函数命名就很有歧义了，其实是正常的逻辑，而不是验证。
// 假设该函数是接口本应的正常处理逻辑
func Check(w http.ResponseWriter, r *http.Request) {
	//执行正常业务逻辑
	fmt.Println("执行正常的逻辑check！")
}

// Auth 对用户请求进行统一验证的拦截器function，
// 每个接口都需要提前经过该函数验证
func Auth(w http.ResponseWriter, r *http.Request) error {
	fmt.Println("执行验证！")
	//添加基于cookie的权限验证
	err := CheckUserInfo(r)
	if err != nil {
		return err
	}
	return nil
}

// CheckUserInfo 身份校验函数
func CheckUserInfo(r *http.Request) error {
	// 1. 获取cookie中的uid，user ID
	uidCookie, err := r.Cookie("uid")
	if err != nil {
		//return errors.New("用户UID Cookie 获取失败！")
	}
	// 2. 获取cookie中的sign（加密的用户信息）
	signCookie, err := r.Cookie("sign")
	if err != nil {
		//return errors.New("用户加密串 Cookie 获取失败！")
	}

	// 3. 对信息进行解密
	signByte, err := encrypt.DePwdCode(signCookie.Value)
	if err != nil {
		//return errors.New("加密串已被篡改！")
	}

	// test code
	fmt.Println("结果比对")
	fmt.Println("用户ID：" + uidCookie.Value)
	fmt.Println("解密后用户ID：" + string(signByte))

	// 结果比对，只是简单的比对连个是否相等
	if checkInfo(uidCookie.Value, string(signByte)) {
		//return nil
	}
	//return errors.New("身份校验失败！")
	return nil
}

// checkInfo 自定义逻辑判断，这里只是进行简单的等值判断
func checkInfo(checkStr string, signStr string) bool {
	if checkStr == signStr {
		return true
	}
	return false
}

func main() {
	// 请求到来的适合，如果发现不是本机的IP，则本服务器充当代理

	//负载均衡器设置
	//采用一致性哈希算法
	hashConsistent =common.NewConsistent()
	//采用一致性hash算法，添加节点
	for _,v :=range hostArray {
		hashConsistent.Add(v)
	}


	// l1nkkk: 这里没用框架，使用裸的http库

	// 1、拦截器生成与注册
	filter := common.NewFilter()
	filter.RegisterFilterUri("/check", Auth)

	// 2、启动服务
	// Check 的命名非常不好，这里的 Check 其实是模拟的正常的接口处理逻辑。
	// /check 是模拟的正常请求接口，容易造成混淆，
	// 主要思路就是在原本的处理逻辑基础上套上一层处理函数，对用户请求该接口进行验证，
	// 当没有问题了，才执行 Check 的处理逻辑
	http.HandleFunc("/check", filter.Handle(Check))
	http.ListenAndServe(":8083", nil)
}
