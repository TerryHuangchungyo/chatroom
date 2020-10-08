var conn;
var userId;
var hubs = new Map();
var lastHubListItem;

window.onbeforeunload = function () {
    conn.close(1000);
};

$(document).ready(function(){

    $("#usernameSettingModal").modal("show");


    $("#usernameCommit").click( function(){
        $("#usernameSettingModal").hide();
        createUser( $("#usernameInput").val() )
        console.log( $("#usernameInput").val() );
    });
    
    $("#hubnameCommit").click( function(){
        $("#hubnameSettingModal").hide();
        createHub( $("#hubnameInput").val() )
    });

    $("#inviteCommit").click( function(){
        let message = { action: INVITE,
                        userId: parseInt($("#inviteInput").val()),
                        userName: "",
                        hubId: lastHubListItem.data("id"),
                        hubName: "",
                        content: "" };
        conn.send( JSON.stringify(message) );
        $("#inviteModal").modal("hide");
    });

    $("#msgForm").submit( function( event ) {
        if( !conn ) {
            alert("聊天室沒有連線")
            return false;
        }

        if( $("#msgInput").val() ) {
            let message = { action: 0,
                userId: userId,
                userName: "",
                hubId: lastHubListItem.data("id"),
                hubName: "",
                content: $("#msgInput").val() };
            conn.send( JSON.stringify( message ) );
            $("#msgInput").val("");
        }
        return false;
    });
});

function createConn( userId ) {
    if ( window["WebSocket"] ) {
        conn = new WebSocket( "ws://" + document.location.host + "/chat/" + userId );

        conn.onclose = function( event ) {
            appendMessage( SYSTEM, "系統訊息", msgCurrentTimeStr(), "聊天通道關閉" );
        };

        conn.onmessage = function( event ) {
            var message = JSON.parse(event.data);
            console.log( message );
            handleMessage( message );
        }

        appendMessage( SYSTEM, "系統訊息", msgCurrentTimeStr(), "聊天通道開啟" );
    } else {
        appendMessage( SYSTEM, "系統訊息", msgCurrentTimeStr(), "你的聊天室不支援websocket" );                         
    }
}

function handleMessage( message ) {
    switch( message.action ) {
        case REPLY:
            let type = OTHER;
            if( message.userId == userId ) {
                type = USER;
            }
            appendMessage( type, message.userName, msgCurrentTimeStr(), message.content );
            break;
        case INVITE:
            let replyMessage = { action: ANSWER,
                                userId: userId,
                                userName: "",
                                hubId: message.hubId,
                                hubName: "",
                                content: "0" };
            
            $.confirm({
                title: '聊天室邀請',
                content: `${message.userName} 邀請你加入 ${message.hubName} 聊天室`,
                buttons: {
                    confirm: {
                        text: '是',
                        btnClass: 'btn-green',
                        keys: ['enter'],
                        action: function() {
                            replyMessage.content = "1";
                            hubs.set( message.hubId, message.hubName);
                            updateHubList( message.hubId );
                            conn.send( JSON.stringify(replyMessage) );
                        }
                    },
                    cancel: {
                        text: '否',
                        action: function() {
                            conn.send( JSON.stringify(replyMessage) );
                        }
                    }
                }
            });
            break;
    }
}

function createUser( username ) {
    $.ajax({
        type: "POST",
        url: document.location.protocol + "//" + document.location.host + "/user",
        data: { username: username }
    }).then( function( data ){
        userId = data["id"];
        $("#username").text( data["username"] );
        createConn( userId )
        alert("創建使用者成功");
    }).fail( function( data ){
        alert("創建使用者失敗");
    })
    $("#usernameSettingModal").modal("hide");
}

function createHub( hubname ) {
    $.ajax({
        type: "POST",
        url: document.location.protocol + "//" + document.location.host + "/hub",
        data: { userId: userId, hubname: hubname }
    }).then( function( data ){
        let hubId = data["id"];
        hubs.set( hubId, data["hubname"])
        updateHubList( hubId )
        console.log("創建聊天室成功");
    }).fail( function( data ){
        console.log("創建聊天室失敗");
    })
    $("#hubnameInput").val("");
    $("#hubnameSettingModal").modal("hide");
}

function updateHubList( hubId = 0 ) {
    $("#hubList").empty()
    for( let [ id, name ] of hubs.entries() ) {
        let list = $("<a></a>").addClass( "nav-link" )
        list.data( "id", id )
        list.text( name )
        list.appendTo( "#hubList" )  
        if( hubId == id ) {
            list.addClass("active")
            lastHubListItem = list
            $("#hubName").text( lastHubListItem.text() )
            $("#hubId").text( String(lastHubListItem.data("id")).padStart( 12, "0") )
        }
        list.click( function(){
            lastHubListItem.removeClass( "active" )
            $(this).addClass("active")
            lastHubListItem = $(this)
            $("#hubName").text( lastHubListItem.text() )
            $("#hubId").text( String(lastHubListItem.data("id")).padStart( 12, "0") )
        })                
    }
}

    
function appendMessage( type, name, time, msg ) {
    let messageBox = $("<div></div>").addClass("bg-light");

    switch( type ) {
        case SYSTEM:
            messageBox.addClass("offset-4")
                    .addClass("col-4")
                    .addClass("rounded")
                    .addClass("mt-2");
            break;
        case OTHER:
            messageBox.addClass("col-4");
            break;
        case USER:
            messageBox.addClass("offset-8")
                    .addClass("col-4");
            break;
    }
    
    messageBox.html( `<small>${name} ${time}</small><div>${msg}</div>`);

    let wrapper = $("<div></div>").addClass("row")
                                .addClass("mt-2")
                                .append( messageBox );
    wrapper.appendTo($("#dialog"))
}

function msgCurrentTimeStr() {
    let str = (new Date()).toISOString().replaceAll("-","/").replace("T", " ");
    str = str.substr( 0, str.lastIndexOf( ":" ) );
    return str;
}