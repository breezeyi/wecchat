package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type TemplateMessage struct {
	ToUser     string       `json:"touser"`
	TemplateID string       `json:"template_id"`
	Data       TemplateData `json:"data"`
}

type TemplateData struct {
	First    TemplateValue `json:"first"`
	Keyword1 TemplateValue `json:"keyword1"`
	Keyword2 TemplateValue `json:"keyword2"`
	Remark   TemplateValue `json:"remark"`
}

type TemplateValue struct {
	Value string `json:"value"`
}

type name struct {
	ToUser     string `json:"touser"`
	TemplateID string `json:"template_id"`
	Data       name01 `json:"data"`
}
type name01 struct {
	Data TemplateValue `json:"data"`
}

var (
	APPID     = "替换为自己的" //测试公众号ID
	APPSECRET = "替换为自己的" // 密码
)

// 定义结构体类型，用于解析从微信 API 返回的访问令牌信息
type token struct {
	// 访问令牌字符串
	AccessToken string `json:"access_token"`

	// 访问令牌的有效期，单位为秒
	ExpiresIn int `json:"expires_in"`
}

func main() {
	accessToken := getaccesstoken()
	apiUrl := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=%s", accessToken)

	/* // 构造模板消息数据
	message := TemplateMessage{
		ToUser:     "oyB_B6sc58sXV6OThHc4bzCWXZ3I",
		TemplateID: "xGtXOFXQ4-QzshrSGjqPxhCCNXHlxA_IbHBAsMOSVB0",
		Data: TemplateData{
			First:    TemplateValue{Value: "您好，您的订单已支付成功！"},
			Keyword1: TemplateValue{Value: "订单号123456"},
			Keyword2: TemplateValue{Value: "10元"},
			Remark:   TemplateValue{Value: "感谢您的购买！"},
		},
	} */
	message := name{
		ToUser:     "oyB_B6sc58sXV6OThHc4bzCWXZ3I",                //要发送给谁的用户ID
		TemplateID: "5OLjTkqytEPmdtsUUtUKXQu5TnuvCV0t18qGl9QGxaI", //消息模版ID
		Data: name01{
			Data: TemplateValue{Value: "你好！小的目前不懂你的心思~不如换个方式说说（每日一句||...的天气)"},
		},
	}
	// 将模板消息转换为JSON格式
	messageBytes, err := json.Marshal(message)
	if err != nil {
		fmt.Println("转换消息失败：", err)
		return
	}
	fmt.Println(string(messageBytes))
	// 发送POST请求
	resp, err := http.Post(apiUrl, "application/json", strings.NewReader(string(messageBytes)))
	if err != nil {
		fmt.Println("发送消息失败：", err)
		return
	}
	defer resp.Body.Close()

	// 解析响应
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取响应失败：", err)
		return
	}
	fmt.Println(string(respBody))
}

// 获取微信accesstoken
func getaccesstoken() string {
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%v&secret=%v", APPID, APPSECRET)
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

	token := token{}
	err = json.Unmarshal(body, &token)
	if err != nil {
		fmt.Println("微信token解析json失败", err)
		return ""
	}

	return token.AccessToken
}
