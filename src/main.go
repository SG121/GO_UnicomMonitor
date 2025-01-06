package main

import (
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/gorilla/websocket"
)

// 全局变量
var (
	globalVar   []byte
	globalVoice []byte
	config      Config
)

// 保存数据
func saveData() {
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
func saveVoice() {
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

// 发送消息
func sendMessage(wsHost string, paramMsg string) {
	//wss://vd-file-hnzz2-wcloud.wojiazongguan.cn:50443/h5player/live
	uri := url.URL{
		Scheme: "wss",
		Host:   wsHost,
		Path:   "/h5player/live",
	}
	conn, _, err := websocket.DefaultDialer.Dial(uri.String(), nil)
	if err != nil {
		log.Fatalf("Failed to connect to WebSocket server: %v", err)
	}
	defer conn.Close()
	// 发送消息
	message := "_paramStr_=" + paramMsg
	err = conn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		FmtPrint("Failed to send message: %v", err)
		return
	}
	// 接收消息
	_, response, err := conn.ReadMessage()
	if err != nil {
		FmtPrint("Failed to receive response: %v", err)
		return
	}
	FmtPrint("Received response:", Decode(string(response)))
	// 发送消息
	cmdMessage := `{"time":1243,"cmd":3}`
	err = conn.WriteMessage(websocket.TextMessage, []byte(cmdMessage))
	if err != nil {
		FmtPrint("Failed to send command: %v", err)
		return
	}
	// 继续接收消息并处理数据
	for {
		_, response, err := conn.ReadMessage()
		if err != nil {
			FmtPrint("Error receiving message: %v", err)
			break
		}
		// 检查特定条件
		if len(response) > 1 && response[1] == 0x63 {
			// 将数据附加到全局变量
			globalVar = append(globalVar, response[0x4e:]...)
			// 打印数据的长度
			FmtPrint("Data length:", len(globalVar))
			// 保存数据
			saveData()
		}
	}
}

// 主函数
func main() {
	FmtPrint("开始启动")
	// https://we.wo.cn/web/smart-club-pc-v2/?clientId=1001000001
	// 读取配置文件
	config, err := ReadConfig()
	if err != nil {
		FmtPrint("读取配置文件出错:", err)
		return
	}
	// 发送消息
	wsHost := config.WsHost
	paramMsg := config.ParamMsg
	sendMessage(wsHost, paramMsg)
	//
	FmtPrint("启动完成")
}
