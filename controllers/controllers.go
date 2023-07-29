package controllers

import (
	"crypto/sha1"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
	"wexin/models"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// 校验微信请求是否合法
func Wechat(c *gin.Context) {
	signature := c.Query("signature")
	timestamp := c.Query("timestamp")
	nonce := c.Query("nonce")
	echostr := c.Query("echostr")

	// 第一步：自然排序
	tmp := []string{viper.GetString("token"), timestamp, nonce}
	sort.Strings(tmp)

	// 第二步：sha1 加密
	sourceStr := strings.Join(tmp, "")
	h := sha1.New()
	h.Write([]byte(sourceStr))
	localSignature := fmt.Sprintf("%x", h.Sum(nil))

	// 第三步：验证签名
	if signature == localSignature {
		c.String(200, echostr)
	} else {
		c.String(401, "Unauthorized")
	}
}

// 接收用户发送的消息
func PostWecaht(c *gin.Context) {
	//判断用户是否是第一次关注，并做出相应的提示
	openid := c.Query("openid")
	go models.FirstAttention(openid)
	buf, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Println("读取消息体失败:", err)
		return
	}

	var msg models.Message
	err = xml.Unmarshal(buf, &msg)
	if err != nil {
		fmt.Println("解析消息体失败:", err)
		return
	}
	//当判断接受到的数据类型是text类型时进行判断
	switch msg.MsgType {
	case "text":
		models.WechatText(msg, msg.FromUserName)
		//case "images":
	}
	c.String(200, "success")
}
