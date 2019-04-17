var socket;
var sid;
$(document).ready(function () {
    // Create a socket
    socket = new WebSocket('ws://' + window.location.host + '/ws?uname=' + $('#uname').text());
    socket.onopen = function() {
        console.log("建立长连接");
        var data = {"type": 11, "data": JSON.stringify({"device_id": $('#uname').text()})}
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
                alert("device_token:" + data['data']['device_token'])
                var tMsg = {"type":4,"data": JSON.stringify({"room_id":"1"})}
                socket.send(JSON.stringify(tMsg))
                break;
            case 2: // JOIN
                alert("加入房间")
                if (data.User == $('#uname').text()) {
                    li.innerText = 'You joined the chat room.';
                } else {
                    li.innerText = data.User + ' joined the chat room.';
                }
                break;
            case 3: // LEAVE
                alert("退出聊天室")
                li.innerText = data.User + ' left the chat room.';
                break;
            case 1: // MESSAGE
                var li = document.createElement('li');
                var username = document.createElement('strong');
                var content = document.createElement('span');
                if (data['sid']) {
                    sid = data['sid']
                }
                username.innerText = data.data['uname'] || "匿名用户";
                content.innerText = data.data['content'];

                li.appendChild(username);
                li.appendChild(document.createTextNode(': '));
                li.appendChild(content);
                $('#chatbox li').first().before(li);
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
                        var li = document.createElement('li');
                        var content = document.createElement('span');
                        var username = document.createElement('strong');

                        if (data['sid']) {
                            sid = data['sid']
                        }
                        username.innerText = data.data["chat_history"][i]['uname'] || "匿名用户";
                        content.innerText = data.data["chat_history"][i]['content'];

                        li.appendChild(username);
                        li.appendChild(document.createTextNode(': '));
                        li.appendChild(content);
                        $('#chatbox li').first().before(li);
                    }
                }
        }
    };

    // Send messages.
    var postConecnt = function () {
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

    $('#sendbtn').click(function () {
        postConecnt();
    });
});
