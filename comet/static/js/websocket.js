var socket;
var sid;
$(document).ready(function () {
    // Create a socket
    socket = new WebSocket('ws://' + window.location.host + '/ws?uname=' + $('#uname').text());
    socket.onopen = function() {
        console.log("建立长连接");
        var data = {"type":4,"data":{"room_id":"1", "device_token":"web_test"}}
        socket.send(JSON.stringify(data))
    };
    // Message received on the socket
    socket.onmessage = function (event) {
        var data = JSON.parse(event.data);
        var li = document.createElement('li');

        console.log(data);

        switch (data['type']) {
            case 2: // JOIN
                if (data.User == $('#uname').text()) {
                    li.innerText = 'You joined the chat room.';
                } else {
                    li.innerText = data.User + ' joined the chat room.';
                }
                break;
            case 3: // LEAVE
                li.innerText = data.User + ' left the chat room.';
                break;
            case 1: // MESSAGE
                var username = document.createElement('strong');
                var content = document.createElement('span');
                if (data['sid']) {
                    sid = data['sid']
                }
                username.innerText = data.data['user'] || "匿名用户";
                content.innerText = data.data['content'];

                li.appendChild(username);
                li.appendChild(document.createTextNode(': '));
                li.appendChild(content);

                break;
            case 99:
                var data = {
                    "type": 100,
                    "data": {
                        "sid": sid,
                        "uname": $('#uname').text(),
                        "content": ""
                    }
                }
                socket.send(JSON.stringify(data));
        }

        $('#chatbox li').first().before(li);
    };

    // Send messages.
    var postConecnt = function () {
        var content = $('#sendbox').val()
        if (!content){
            alert("发送的内容不能为空")
            return
        }
        var data = {
            "type": 1,
            "data": {
                "room_id": "1",
                "device_token": sid,
                "uname": $('#uname').text(),
                "content": $('#sendbox').val()
            }
        }
        socket.send(JSON.stringify(data));
        $('#sendbox').val('');
    }

    $('#sendbtn').click(function () {
        postConecnt();
    });
});
