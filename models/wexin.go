package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
)

// 定义结构体类型，用于解析从微信 API 返回的访问令牌信息
type Token struct {
	// 访问令牌字符串
	AccessToken string `json:"access_token"`

	// 访问令牌的有效期，单位为秒
	ExpiresIn int `json:"expires_in"`
}

// 定义默认回复结构体
type Defaultreply struct {
	ToUser     string       `json:"touser"`      //接收者ID
	TemplateID string       `json:"template_id"` //发送的消息模版ID
	Data       TemplateData `json:"data"`        //发送消息的结构题
}
type TemplateData struct {
	Data TemplateValue `json:"data"`
}
type TemplateValue struct {
	Value string `json:"value"`
}

// 获取微信accesstoken
func Getaccesstoken() string {
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%v&secret=%v", viper.GetString("APPID"), viper.GetString("APPSECRET"))
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("获取微信token失败", err)
		return ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("微信token读取失败", err)
		return ""
	}

	token := Token{}
	err = json.Unmarshal(body, &token)
	if err != nil {
		fmt.Println("微信token解析json失败", err)
		return ""
	}

	return token.AccessToken
}

// 获取关注人列表
func Getflist(access_token string) []gjson.Result {
	url := "https://api.weixin.qq.com/cgi-bin/user/get?access_token=" + access_token + "&next_openid="
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("获取关注列表失败", err)
		return nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取内容失败", err)
		return nil
	}
	flist := gjson.Get(string(body), "data.openid").Array()
	return flist
}

// 发送模板消息
//
//	--------------获取微信的token 	发送消息的结构字段 	  详情跳转的地址  发送消息的模版ID 	 接受消息的用户ID
func Templatepost(access_token string, reqdata string, fxurl string, templateid string, openid string) {
	url := "https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=" + access_token

	reqbody := "{\"touser\":\"" + openid +
		"\", \"template_id\":\"" + templateid +
		"\", \"url\":\"" + fxurl +
		"\", \"data\": " + reqdata + "}"

	resp, err := http.Post(url,
		"application/x-www-form-urlencoded",
		strings.NewReader(string(reqbody)))
	if err != nil {
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(body))
}

// 处理接受到的数据
func WechatText(msg Message, FromUserName string) {
	log.Printf("%v发送的消息:%v", msg.FromUserName, msg.Content)
	log.Println(msg)
	switch {
	//用户发送你好
	case msg.Content == "你好":
		Templatepost(Getaccesstoken(), Hello(), "", viper.GetString("HelloTemplateID"), FromUserName)

		//用户发送每日一句
	case msg.Content == "每日一句", msg.Content == "每日一句.", msg.Content == "每日一句。":
		data, fxurl := EcapsulationData()
		Templatepost(Getaccesstoken(), data, fxurl, viper.GetString("SentTemplateID"), FromUserName)

		//用户想要查询天气时
	case strings.Contains(msg.Content, "天气"):
		Filtration(msg)
	}
}

// 自动回复你好字段
func Hello() string {

	reqdata := "{\"data\":{\"value\":\"" +
		viper.GetString("Reply.Hello") + "\"}}"
	return reqdata

}

type storage struct {
	UserID map[string]bool
	s      sync.Mutex
}

var UserStorage = storage{
	UserID: make(map[string]bool),
}

// 处理用户第一次关注-------（正式公众号不需要，因为可以自己自动设置）
func FirstAttention(userid string) {
	//判断是否是一个空请求
	if userid != "" {
		//判断用户ID是否已经存在,不存在则存入，否则不做处理
		UserStorage.s.Lock()
		if !UserStorage.UserID[userid] {
			//用户不存在，存入发送关于第一次提示信息
			defer UserStorage.s.Unlock()
			UserStorage.UserID[userid] = true
			//存入用户ID后发送提示信息
			Templatepost(Getaccesstoken(), OneConcern(), "", viper.GetString("HelloTemplateID"), userid)
		}
	}
}

// 第一次关注
func OneConcern() string {

	reqdata := "{\"data\":{\"value\":\"" +
		viper.GetString("Reply.OneConcern") + "\"}}"
	return reqdata

}
