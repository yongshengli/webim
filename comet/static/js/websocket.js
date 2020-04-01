var socket;
var sid;
$(document).ready(function () {
    // Create a socket
    socket = new WebSocket('ws://' + window.location.host + '/ws');
    socket.onopen = function() {
        console.log("建立长连接");
        var data = {"type": 11, "data": JSON.stringify({"device_id": $('input[name="device_id"]').val()})}
        socket.send(JSON.stringify(data))
    };
    // Message received on the socket
    socket.onmessage = function (event) {
        var data = JSON.parse(event.data);
        if (data.data!="") {
            data.data = JSON.parse(data.data)
        }
        console.log(data);

        switch (data['type']) {
            case 0:
                alert("单播推送:"+data.data['content']);
                break;
            case 10:
                var nHtml = "<li>通知:"+data.data['content']+"</li>"
                $('#noticebox').append(nHtml);
                break;
            case 11:
                if (data.device_token!="undefined" && data.device_token.length>5){
                    showMainBox()
                    var tMsg = {"type":4,"data": JSON.stringify({"room_id":"1"})}
                    socket.send(JSON.stringify(tMsg))
                }
                break;
            case 2: // JOIN
                alert("加入房间")
                if (data.User == $('#uname').text()) {
                    // li.innerText = 'You joined the chat room.';
                } else {
                    // li.innerText = data.User + ' joined the chat room.';
                }
                break;
            case 3: // LEAVE
                alert("退出聊天室")
                // li.innerText = data.User + ' left the chat room.';
                break;
            case 1: // MESSAGE
                addMsg(data.data)
                break;
            case 99:
                var tMsg = {
                    "type": 100,
                    "data": JSON.stringify({
                        "sid": sid,
                        "uname": $('#uname').text(),
                        "content": ""
                    })
                }
                socket.send(JSON.stringify(tMsg));
            case 4:
                if (typeof data.data["chat_history"]!="undefined"){
                    for (var i=0; i<data.data["chat_history"].length; i++){
                        addMsg(data.data["chat_history"][i])
                    }
                }
            case 5:
                if(data.data['code']!="undefined" && data.data['code']==0){
                    showMainBox()
                    var tMsg = {"type":4,"data": JSON.stringify({"room_id":"1"})}
                    socket.send(JSON.stringify(tMsg))
                }else{
                    alert(data.data["content"])
                } 
        }
    };
    function showMainBox(){
        $('#login-box').hide()
        $('#main-box').show()
    }
    // Send messages.
    var postContent = function () {
        var content = $('#sendbox').val()
        if (!content){
            alert("发送的内容不能为空")
            return
        }
        var tmpMsg = {
            "type": 1,
            "data": JSON.stringify({
                "room_id": "1",
                "device_token": sid,
                "uname": $('#uname').text(),
                "content": $('#sendbox').val()
            })
        }
        socket.send(JSON.stringify(tmpMsg));
        $('#sendbox').val('');
    }
    Date.prototype.Format = function(fmt) { //author: meizz
        var o = {
            "M+" : this.getMonth()+1,                 //月份
            "d+" : this.getDate(),                    //日
            "h+" : this.getHours(),                   //小时
            "m+" : this.getMinutes(),                 //分
            "s+" : this.getSeconds(),                 //秒
            "q+" : Math.floor((this.getMonth()+3)/3), //季度
            "S"  : this.getMilliseconds()             //毫秒
        };
        if(/(y+)/.test(fmt))
            fmt=fmt.replace(RegExp.$1, (this.getFullYear()+"").substr(4 - RegExp.$1.length));
        for(var k in o)
            if(new RegExp("("+ k +")").test(fmt))
                fmt = fmt.replace(RegExp.$1, (RegExp.$1.length==1) ? (o[k]) : (("00"+ o[k]).substr((""+ o[k]).length)));
        return fmt;
    }
    function addMsg(data){
        var li = document.createElement('div');
        var face = document.createElement('img');
        var content = document.createElement('div');
        face.className = 'direct-chat-img'

        content.className = 'direct-chat-text'
        if (data['sid']) {
            sid = data['sid']
        }
        face.alt = data['uname'] || "匿名用户";
        content.innerText = data['content'];
        var info = document.createElement('div')
        info.className = "direct-chat-info clearfix"

        var uname = document.createElement("span")
        var ctime = document.createElement("span")

        if (data['uname']==$('#uname').text()) {
            li.className = "direct-chat-msg right";
            face.src = "/static/img/user3-128x128.jpg"
            uname.className = 'direct-chat-name pull-right'
            ctime.className = 'direct-chat-timestamp pull-left'
        }else{
            li.className = "direct-chat-msg";
            face.src = "/static/img/user1-128x128.jpg"
            uname.className = 'direct-chat-name pull-left'
            ctime.className = 'direct-chat-timestamp pull-right'
        }
        uname.innerText = data['uname'] || "匿名用户";
        ctime.innerText = new Date(data['c_t']*1000).Format("yyyy-M-d h:m:s")
        info.appendChild(uname)
        info.appendChild(ctime)
        li.appendChild(info);
        li.appendChild(face);
        li.appendChild(content);
        $('#chatbox').append(li);
        var ele = document.getElementById('chatbox');
        ele.scrollTop = ele.scrollHeight;
    }
    window.postContent = postContent
    $('#sendbtn').click(function () {
        postContent();
    });
    $('#login-btn').click(function () {
        uname = $('input[name="uname"]').val()
        deviceId = $('input[name="device_id"]').val()
        password = $('input[name="password"]').val()
        if (uname=="" || password=="" || deviceId==""){
            alert("用户名或者密码或者设备id为空")
            return
        }
        var tmpMsg = {
            "type": 5,
            "data": JSON.stringify({
                "device_id": deviceId,
                "username": uname,
                "password": password
            })
        }
        //console.log(tmpMsg)
        socket.send(JSON.stringify(tmpMsg))
    })
});
