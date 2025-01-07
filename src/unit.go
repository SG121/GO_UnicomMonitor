package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math"
	"net/url"
	"os"
	"time"
)

// 全局变量
var (
	globalVar   []byte
	globalVoice []byte
	config      Config
)

// 配置文件
type Config struct {
	Size     int    `json:"size"`
	WsHost   string `json:"wsHost"`
	ParamMsg string `json:"paramMsg"`
}

// 读取配置文件
func ReadConfig() (Config, error) {
	var config Config
	filePath := "config.json"
	data, err := os.ReadFile(filePath)
	if err != nil {
		return config, err
	}
	err = json.Unmarshal(data, &config)
	return config, err
}

// 编码
func Encode(s string) string {
	t := int(math.Ceil(float64(len(s)) / 2))
	shiftedString := s[t:] + s[:t]
	escapedString := []byte(shiftedString)
	encodedString := base64.StdEncoding.EncodeToString(escapedString)
	result := "MTc2NDAxND" + encodedString
	return result
}

// 解码
func Decode(encodedString string) string {
	encodedString = encodedString[10:]
	decodedBytes, _ := base64.StdEncoding.DecodeString(encodedString)
	decodedString := string(decodedBytes)
	decodedString, _ = url.QueryUnescape(decodedString)
	t := int(math.Ceil(float64(len(decodedString)) / 2))
	originalString := decodedString[t:] + decodedString[:t]
	return originalString
}

// 保存数据
func SaveData() {
	if len(globalVar) > 1024*config.Size {
		fileName := fmt.Sprintf("%d_video.flv", GetNowTime())
		file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			FmtPrint("Failed to save voice data: %v", err)
			return
		}
		defer file.Close()
		file.Write(globalVar)
		globalVar = []byte{}
	}
}

// 保存声音
func SaveVoice() {
	if len(globalVoice) > 1024*config.Size {
		fileName := fmt.Sprintf("%d_voice.flv", GetNowTime())
		file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			FmtPrint("Failed to save voice data: %v", err)
			return
		}
		defer file.Close()
		file.Write(globalVoice)
		globalVoice = []byte{}
	}
}

// 定义内置的打印语句
func FmtPrint(data ...any) {
	date := time.Now().Format("2006-01-02 15:04:05")
	if len(data) == 1 {
		fmt.Println(date+": ", data[0])
	} else {
		fmt.Println(date+": ", data)
	}
}

// 获取当前时间
func GetNowTime() int64 {
	timestamp := time.Now().Unix()
	return timestamp
}
