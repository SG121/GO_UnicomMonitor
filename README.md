# 联通监控数据抓取

## 来源

[py 源码](back)是妖友的（妖火论坛）。本人转换成 golang 更方便使用

## 入口

联通看家：https://we.wo.cn/web/smart-club-pc-v2/?clientId=1001000001

## 说明

1. 从 [Releases](https://github.com/zgcwkjOpenProject/GO_UnicomMonitor/releases) 下载 **二进制程序** 和 **config.json** 文件
2. 修改配置文件 **config.json**，具体参考 [妖友源码说明](back)。
3. 启动程序，会立刻抓取数据。
4. 视频文件缺少文件头，可以用 ffmpeg 转换或补上缺少文件头。

## 开发

- [x] 摄像头录像
- [ ] 摄像头录音
- [ ] 补上缺少的视频文件头，使文件可以直接播放
- [ ] 设置摄像头录下存储上限
- [ ] 支持多个摄像头录像
