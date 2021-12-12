package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

//高级加密标准AES

//16,24,32位字符串的话，分别对应AES-128，AES-192, AES-256 加密方法
//key不能泄露
var PwdKey = []byte("DIS**#KKKDJJSKDI")

//PKCS7 填充模式
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	//Repeat()函数的功能是把切片[]byte{byte(padding)}复制padding个，然后合并成新的字节切片返回
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

//实现加密
func AesEcrypt(origData []byte, key []byte) ([]byte, error) {
	//创建加密算法实例
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	//获取块的大小
	blockSize := block.BlockSize()
	//对数据进行填充，让数据长度满足要求
	origData = PKCS7Padding(origData, blockSize)
	//采用AES加密方法中CBC加密模式
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	//执行加密
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

//加密base
func EnPwdCode(pwd []byte) (string, error) {
	result, err := AesEcrypt(pwd, PwdKey)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(result), err
}

//填充的反向操作，删除填充字符串
func PKCS7UnPadding