package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/tidwall/gjson"
)

type Sentence struct {
	Content     string `json:"content"`     //英文字段
	Note        string `json:"note"`        //中文字段
	Translation string `json:"translation"` //标题来自来哪里
}

// 获取每日一句
func Getsen() (Sentence, string) {
	resp, err := http.Get("http://open.iciba.com/dsapi/?date")
	sent := Sentence{}
	if err != nil {
		fmt.Println("获取每日一句失败", err)
		return sent, ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取内容失败", err)
		return sent, ""
	}

	err = json.Unmarshal(body, &sent)
	if err != nil {
		fmt.Println("每日一句解析json失败")
		return sent, ""
	}
	fenxiangurl := gjson.Get(string(body), "fenxiang_img").String()
	// fmt.Println(sent)
	// fmt.Println(fenxiangurl)
	return sent, fenxiangurl //返回的是中文、英文、来自哪里、地址
}

// 封装返回的信息
func EcapsulationData() (string, string) {
	req, fxurl := Getsen()
	if req.Content == "" {
		return "", ""
	}

	reqdata := "{\"content\":{\"value\":\"" +
		req.Content + "\"}, \"note\":{\"value\":\"" +
		req.Note + "\"}, \"translation\":{\"value\":\"" +
		req.Translation + "\"}}"
	return reqdata, fxurl
}
