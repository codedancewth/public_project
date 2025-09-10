package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// RandString 生成随机字符串
func RandString(codeLen int) string {
	// 1. 定义原始字符串
	timeDo := time.Now().Unix()
	str := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_%v_"
	rawStr := fmt.Sprintf(str, timeDo)
	// 2. 定义一个buf，并且将buf交给bytes往buf中写数据
	buf := make([]byte, 0, codeLen)
	b := bytes.NewBuffer(buf)
	// 随机从中获取
	rand.Seed(time.Now().UnixNano())
	for rawStrLen := len(rawStr); codeLen > 0; codeLen-- {
		randNum := rand.Intn(rawStrLen)
		b.WriteByte(rawStr[randNum])
	}
	return b.String()
}

// GetTimeId 获取时间id
func GetTimeId() int64 {
	var (
		now    = time.Now()
		timeId = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Unix() // 设置为当天的 00:00:00
	)
	return timeId
}

// ChunkSlice chunk函数
func ChunkSlice(slice []int64, chunkSize int) [][]int64 {
	var chunks [][]int64

	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize

		if end > len(slice) {
			end = len(slice)
		}

		chunks = append(chunks, slice[i:end])
	}

	return chunks
}

// If 三目运算符
func If(isTrue bool, a, b interface{}) interface{} {
	if isTrue {
		return a
	}
	return b
}

// GetDayTTl 获取每日的剩余ttl
func GetDayTTl() int64 {
	var (
		now   = time.Now()
		today = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()) // 设置为当天的 00:00:00
		ttl   = today.Unix() + 86400 - now.Unix()
	)

	return ttl
}

// GetHourTTl 获取每小时的剩余ttl
func GetHourTTl() int64 {
	var (
		now   = time.Now()
		today = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()) // 设置为当天的 00:00:00
		dis   = now.Unix() - today.Unix()
		ttl   = 3600 - dis + (dis / 3600 * 3600)
	)

	return ttl
}

const (
	DINGDING_MSG_TYPE_TEXT = "text"
)

type DingDingMsgSt struct {
	MsgType string `json:"msgtype"`
	Text    struct {
		Content string `json:"content"`
	} `json:"text"`
	Link struct {
		Text       string `json:"text"`
		Title      string `json:"title"`
		PicURL     string `json:"picUrl"`
		MessageURL string `json:"messageUrl"`
	} `json:"link"`
	Markdown struct {
		Title string `json:"title"`
		Text  string `json:"text"`
	} `json:"markdown"`
	At struct {
		AtMobiles []string `json:"atMobiles"`
		IsAtAll   bool     `json:"isAtAll"`
	} `json:"at"`
}

func DingDingSendMsg(param *DingDingMsgSt, URL string) (err error) {
	buf, err := json.Marshal(param)
	if err != nil {
		return
	}
	req, err := http.NewRequest("POST", URL, bytes.NewReader(buf))
	if err != nil {
		return
	}
	req.Header.Add("Content-Type", "application/json;charset=utf-8")

	client := &http.Client{
		Timeout: time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	var respData struct {
		ErrCode int64  `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}
	err = json.Unmarshal(body, &respData)
	if err != nil {
		return
	}
	if respData.ErrCode != 0 {
		err = fmt.Errorf("dingding resp: %s", body)
	}
	return
}

func GetMiliSecondString() string {
	mili := time.Now().UnixNano() / 1e3
	return strconv.Itoa(int(mili))
}

// UtilJsonMarshal 直接转化为string
func UtilJsonMarshal(any interface{}) string {
	marshal, _ := json.Marshal(any)
	return string(marshal)
}
