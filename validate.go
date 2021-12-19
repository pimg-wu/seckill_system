package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"shopping_system/common"
	datamodels "shopping_system/dataModels"
	"shopping_system/encrypt"
	"shopping_system/rabbitmq"
	"strconv"
	"sync"
)

//设置集群地址，最好内外ip
var hostArray = []string{"127.0.0.1", "127.0.0.1"}

var localHost = ""

//数量控制接口服务器内网ip，或者getone的slb内网ip
var GetOneIp = "127.0.0.1"

var GetOnePort = "8084"

var port = "8083"

var hashConsistent *common.Consistent

var rabbitmqValidate *rabbitmq.RabbitMQ

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
	hostUrl := "http://" + host + ":" + port + "/checkRight"
	response, body, err := GetCurl(hostUrl, request)
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

//模拟请求
func GetCurl(hostUrl string, request *http.Request) (response *http.Response, body []byte, err error) {
	uidPre, err := request.Cookie("uid")
	if err != nil {
		return
	}
	uidSign, err := request.Cookie("sign")
	if err != nil {
		return
	}

	//模拟接口访问
	client := &http.Client{}
	req, err := http.NewRequest("GET", hostUrl, nil)
	if err != nil {
		return
	}

	//手动指定，排查多余cookies
	cookieUid := &http.Cookie{Name: "uid", Value: uidPre.Value, Path: "/"}
	cookieSign := &http.Cookie{Name: "sign", Value: uidSign.Value, Path: "/"}
	//添加cookie到模拟的请求中
	req.AddCookie(cookieUid)
	req.AddCookie(cookieSign)

	//获取返回结果
	response, err = client.Do(req)
	defer response.Body.Close()
	if err != nil {
		return
	}
	body, err = ioutil.ReadAll(response.Body)
	return
}

func CheckRight(w http.ResponseWriter, r *http.Request) {
	right := accessControl.GetDistributedRight(r)
	if !right {
		w.Write([]byte("false"))
		return
	}
	w.Write([]byte("true"))
	return
}

//执行正常业务逻辑
func Check(w http.ResponseWriter, r *http.Request) {
	fmt.Println("执行check！")
	queryForm, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil || len(queryForm["productID"]) <= 0 {
		w.Write([]byte("false"))
		return
	}

	productString := queryForm["productID"][0]
	fmt.Println(productString)

	userCookie, err := r.Cookie("uid")
	if err != nil {
		w.Write([]byte("false"))
		return
	}

	//分布式权限验证
	right := accessControl.GetDistributedRight(r)
	if right == false {
		w.Write([]byte("false"))
		return
	}
	//获取数量控制权限，防止秒杀出现超卖现象
	hostUrl := "http://" + GetOneIp + ":" + GetOnePort + "/getOne"
	responseValidate, validateBody, err := GetCurl(hostUrl, r)
	if err != nil {
		w.Write([]byte("false"))
		return
	}
	//判断数量控制接口请求状态
	if responseValidate.StatusCode == 200 {
		if string(validateBody) == "true" {
			//整合下单
			productID, err := strconv.ParseInt(productString, 10, 64)
			if err != nil {
				w.Write([]byte("false"))
				return
			}

			userID, err := strconv.ParseInt(userCookie.Value, 10, 64)
			if err != nil {
				w.Write([]byte("false"))
				return
			}

			//创建消息体
			message := datamodels.NewMessage(userID, productID)
			byteMessage, err := json.Marshal(message)
			if err != nil {
				w.Write([]byte("false"))
				return
			}
			//生产消息
			err = rabbitmqValidate.PublishSimple(string(byteMessage))
			if err != nil {
				w.Write([]byte("false"))
				return
			}
			w.Write([]byte("true"))
			return
		}
	}
	w.Write([]byte("false"))
	return
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

	localIp, err := common.GetIntranceIp()
	if err != nil {
		fmt.Println(err)
	}
	localHost = localIp
	fmt.Println(localHost)

	rabbitmqValidate = rabbitmq.NewRabbitMQSimple("product")
	defer rabbitmqValidate.Destory()

	//1、过滤器
	filter := common.NewFilter()
	filter.RegisterFilterUri("/check", Auth)
	//2、启动服务
	http.HandleFunc("/check", filter.Handle(Check))
	http.ListenAndServe(":8083", nil)
}
