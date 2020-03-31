# Web IM

[中文文档](README_ZH.md)

This sample is about using long polling and WebSocket to build a web-based chat room based on beego.

- [Documentation](http://beego.me/docs/examples/chat.md)

## Installation

```
cd $GOPATH/src/samples/WebIm
go get github.com/gorilla/websocket
go get github.com/beego/i18n
bee run
```

## Usage

enter chat room from 

```
http://127.0.0.1:8080 
```

## API 接口说明

```
 / 欢迎页
 /webim  聊天室主页
 /room/create  //创建聊天室接口
 /room/delete  //删除聊天室接口
 /ws           //websoket长链接
 /push/unicast   //单播推送接口
 /push/broadcast  //广播推送接口
 /monitor/status  //系统监控接口
```
