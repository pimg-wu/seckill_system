package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"shopping_system/common"
	"shopping_system/encrypt"
	"strconv"
	"sync"
)

//设置集群地址，最好内外ip
var hostArray = []string{"127.0.0.1", "127.0.0.1"}

var localHost = "127.0.0.1"

var port = "8081"

var hashConsistent *common.Consistent

type AccessControl struct {
	sourcesArray map[int]interface{} //存放用户想要存放的信息
	sync.RWMutex
}

var accessControl = &AccessControl{sourcesArray: make(map[int]interface{})}

//获取制定的数据
func (m *AccessControl) GetNewRecord(uid int) interface{} {
	m.RWMutex.RLock()
	defer m.RWMutex.RUnlock()
	data := m.sourcesArray[uid]
	return data
}

//设置记录
func (m *AccessControl) SetNewRecord(uid int) {
	m.RWMutex.Lock()
	m.sourcesArray[uid] = "hello world!"
	m.RWMutex.Unlock()
}

func (m *AccessControl) GetDistributedRight(req *http.Request) bool {
	uid, err := req.Cookie("uid")
	if err != nil {
		return false
	}

	//采用一致性hash算法，根据用户id判断获取具体机器
	hostRequest, err := hashConsistent.Get(uid.Value)
	if err != nil {
		return false
	}

	//判断是否为本机
	if hostRequest == localHost {
		return m.GetDataFromMap(uid.Value)
	} else {
		return GetDataFromOtherMap(hostRequest, req)
	}
}

//获取本机map，并处理业务逻辑
func (m *AccessControl) GetDataFromMap(uid string) (isOK bool) {
	uidInt, err := strconv.Atoi(uid)
	if err != nil {
		return false
	}
	data := m.GetNewRecord(uidInt)

	if data != nil {
		return true
	}
	return
}

//获取其他节点处理结果
func GetDataFromOtherMap(host string, request *http.Request) bool {
	uidPre, err := request.Cookie("uid")
	if err != nil {
		return false
	}

	uidSign, err := request.Cookie("sign")
	if err != nil {
		return false
	}

	//模拟接口访问
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://"+host+":"+port+"/check", nil)
	if err != nil {
		return false
	}

	cookieUid := &http.Cookie{Name: "uid", Value: uidPre.Value, Path: "/"}
	cookieSign := &http.Cookie{Name: "sign", Value: uidSign.Value, Path: "/"}
	req.AddCookie(cookieUid)
	req.AddCookie(cookieSign)

	response, err := client.Do(req)
	if err != nil {
		return false
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return false
	}

	if response.StatusCode == 200 {
		if string(body) == "true" {
			return true
		} else {
			return false
		}
	}
	return false
}

//执行正常业务逻辑
func Check(rw http.ResponseWriter, r *http.Request) {
	fmt.Println("执行check！")
}

//统一验证拦截器，每个接口都需要提前验证
func Auth(rw http.ResponseWriter, r *http.Request) error {
	fmt.Println("执行验证！")
	//添加基于Cookie的权限验证
	err := CheckUserInfo(r)
	if err != nil {
		return err
	}
	return nil
}

//身份校验函数
func CheckUserInfo(r *http.Request) error {
	uidCookie, err := r.Cookie("uid")
	if err != nil {
		return errors.New("用户uid Cookie 获取失败！")
	}
	signCookie, err := r.Cookie("sign")
	if err != nil {
		return errors.New("用户加密串 Cookie 获取失败！")
	}

	signByte, err := encrypt.DePwdCode(signCookie.Value)
	if err != nil {
		return errors.New("加密串已被篡改！")
	}

	fmt.Println("结果比对")
	println("用户id：" + uidCookie.Value)
	println("解密后用户id：" + string(signByte))
	if checkInfo(uidCookie.Value, string(signByte)) {
		return nil
	}
	return errors.New("校验失败！")
}

//
func checkInfo(checkStr string, signStr string) bool {
	if checkStr == signStr {
		return true
	}
	return false
}

func main() {
	//负载均衡设置
	//采用一致性哈希算法
	hashConsistent = common.NewConsistent()
	for _, v := range hostArray {
		hashConsistent.Add(v)
	}

	//1、过滤器
	filter := common.NewFilter()
	filter.RegisterFilterUri("/check", Auth)
	//2、启动服务
	http.HandleFunc("/check", filter.Handle(Check))
	http.ListenAndServe(":8083", nil)
}
