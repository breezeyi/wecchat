package models

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// 发送天气的格式
type Weather struct {
	City    string // 城市
	Week    string // 日期
	Weather string // 天气
	Temp    string // 温度
	Wind    string // 风向
	Power   string // 风力
	Prompt  string // 提示
}

// 处理有无用户特殊要求
// 接收的参数 ---用户发送的消息    ---用户请求天气的信息   ----发送信息的用户ID
func UserRequirements(content string, data Data, userid string) {

	//做出判断用户是否有特殊要求，没有则发送默认格式（根据当前时间来判断发送）
	switch {
	//判断是否有明天字段 ------1
	case strings.Contains(content, "明天"):
		dispose(1, content, data, userid)

	//判断是否有后天字段 ------2
	case strings.Contains(content, "后天"):
		dispose(2, content, data, userid)

	//表示默认当天  -------0
	default:
		dispose(0, content, data, userid)
	}
}

// 根据要求值做出相应的判断 ---要求值   -----用户发送的消息   -----请求天气的参数
// 这里只做简单处理只认可白天或者晚上
func dispose(number int, content string, data Data, userid string) {
	//做出判断用户是否有特殊要求，没有则发送默认格式（根据当前时间来判断发送）

	if strings.Contains(content, "晚上") {
		//赋值为晚上天气信息
		weather := Weather{
			City: data.Forecasts[0].City,
			//Week:    data.Forecasts[0].Casts[number].Week,
			Weather: data.Forecasts[0].Casts[number].Nightweather,
			Temp:    data.Forecasts[0].Casts[number].Nighttemp + "℃",
			Wind:    data.Forecasts[0].Casts[number].Nightwind,
			Power:   data.Forecasts[0].Casts[number].Nightpower + "级",
			Prompt:  data.Prompt,
		}
		//根据number值给日期添加后缀  今天 明天 后天
		if number == 0 {
			weather.Week = data.Forecasts[0].Casts[number].Week + "(今天)"
		} else if number == 1 {
			weather.Week = data.Forecasts[0].Casts[number].Week + "(明天)"
		} else {
			weather.Week = data.Forecasts[0].Casts[number].Week + "(后天)"
		}

		//传入赋值后的结构体和接受用户的ID
		packdata(weather, userid)
	} else {
		//赋值为白天天气信息
		//赋值为晚上天气信息
		weather := Weather{
			City: data.Forecasts[0].City,
			//Week:    data.Forecasts[0].Casts[number].Week,
			Weather: data.Forecasts[0].Casts[number].Dayweather,
			Temp:    data.Forecasts[0].Casts[number].Daytemp + "℃",
			Wind:    data.Forecasts[0].Casts[number].Daywind,
			Power:   data.Forecasts[0].Casts[number].Daypower + "级",
			Prompt:  data.Prompt,
		}
		//根据number值给日期添加后缀  今天 明天 后天
		if number == 0 {
			weather.Week = data.Forecasts[0].Casts[number].Week + "(今天)"
		} else if number == 1 {
			weather.Week = data.Forecasts[0].Casts[number].Week + "(明天)"
		} else {
			weather.Week = data.Forecasts[0].Casts[number].Week + "(后天)"
		}

		//传入赋值后的结构体和接受用户的ID
		packdata(weather, userid)
	}
}

// 将赋值好的消息结构体打包
func packdata(data Weather, userid string) {
	reqdata := "{\"city\":{\"value\":\"" +
		data.City + "\"}, \"week\":{\"value\":\"" +
		data.Week + "\"}, \"weather\":{\"value\":\"" +
		data.Weather + "\"}, \"temp\":{\"value\":\"" +
		data.Temp + "\"}, \"wind\":{\"value\":\"" +
		data.Wind + "\"}, \"power\":{\"value\":\"" +
		data.Power + "\"}, \"prompt\":{\"value\":\"" +
		data.Prompt + "\"}}"
	fmt.Println(reqdata)
	Templatepost(Getaccesstoken(), reqdata, "", viper.GetString("WeatTemplateID"), userid)
}

// 判断用户请求的时间是否早上还是晚上 白天----true   晚上-----false
func DetermineTheCurrentTime() bool {
	// 获取当前时间
	now := time.Now()

	// 获取当天日出和日落的时间
	sunrise := time.Date(now.Year(), now.Month(), now.Day(), 6, 0, 0, 0, time.Local)
	sunset := time.Date(now.Year(), now.Month(), now.Day(), 18, 0, 0, 0, time.Local)

	// 判断当前时间是否在日出和日落之间
	if now.After(sunrise) && now.Before(sunset) {
		log.Println("现在是白天")
		return true
	} else {
		log.Println("现在是晚上")
		return false
	}
}
