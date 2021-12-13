package main

import (
	"errors"
	"fmt"
	"net/http"
	"shopping_system/common"
	"shopping_system/encrypt"
)

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
	//1、过滤器
	filter := common.NewFilter()
	filter.RegisterFilterUri("/check", Auth)
	//2、启动服务
	http.HandleFunc("/check", filter.Handle(Check))
	http.ListenAndServe(":8083", nil)
}
