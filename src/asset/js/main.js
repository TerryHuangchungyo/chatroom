var conn;
var userId;
var hubs = new Map();
var lastHubListItem;

window.onbeforeunload = function () {
    conn.close(1000);
};

$(document).ready(function(){
    let queryStr = window.location.search
    const urlParams = new URLSearchParams( queryStr )
    userId = urlParams.get( "userId" );
    $("#userId").text( userId )
    createConn( userId )

    $.ajax({
        url: document.location.protocol + "//" + document.location.host + "/hub/" + userId ,
        type: "GET",
        dataType: "json"
    }).then( function( hubLists, textStatus, xhr ) {
        for( hubInfo of hubLists ) {
            let hub = new Hub( hubInfo.hubId, hubInfo.hubName );
            hubs.set( hubInfo.hubId, hub);
        }

        if (hubLists)
            updateHubList( hubLists[0].hubId);
    }).fail( function( xhr, textStatus ) {
        console.log( xhr.status + ":" + textStatus )
    });

    $("#hubnameCommit").click( function(){
        $("#hubnameSettingModal").modal("hide")
        createHub( $("#hubnameInput").val() )
    });

    $("#inviteCommit").click( function(){
        let message = { action: INVITE,
                        userId: $("#inviteInput").val(),
                        userName: "",
                        hubId: lastHubListItem.data("id"),
                        hubName: "",
                        content: "" };
        conn.send( JSON.stringify(message) );
        $("#inviteModal").modal("hide");
    });

    $("#msgForm").submit( function( event ) {
        if( !conn ) {
            alert("聊天室沒有連線");
            return false;
        } else {
            if( $("#msgInput").val() ) {
                let message = { action: 0,
                    userId: userId,
                    userName: "",
                    hubId: lastHubListItem.data("id"),
                    hubName: hubs.get(lastHubListItem.data("id")).name,
                    content: $("#msgInput").val() };
                conn.send( JSON.stringify( message ) );
                $("#msgInput").val("");
            }
            return false;
        }
    });
});

function createConn( userId ) {
    if ( window["WebSocket"] ) {
        conn = new WebSocket( "ws://" + document.location.host + "/chat/" + userId  );

        conn.onclose = function( event ) {
            for( let [ id, hub ] of hubs.entries() ) {
                hub.appendMessage( SYSTEM, "系統訊息", msgCurrentTimeStr(), "聊天通道關閉" );         
            }
        };

        conn.onmessage = function( event ) {
            var message = JSON.parse(event.data);
            console.log( message );
            handleMessage( message );
        }

        // alert("已開啟websocket"); 
    } else {
        alert("你的瀏覽器不支援websocket");                          
    }
}

function handleMessage( message ) {
    switch( message.action ) {
        case MESSAGE:
            let type = OTHER;
            if( message.userId == userId ) {
                type = USER;
            }
            let hub = hubs.get( message.hubId );
            hub.appendMessage( type, message.userName, msgTimeStrToFormat(message.time), message.content );
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
                            let hub = new Hub( message.hubId, message.hubName );
                            hubs.set( message.hubId, hub);
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

function createHub( hubName ) {
    $.ajax({
        type: "POST",
        url: document.location.protocol + "//" + document.location.host + "/hub",
        data: { userId: userId, hubName: hubName }
    }).then( function( data, textStatus, xhr ){
        let hubId = data["hubId"];
        let hub = new Hub( data["hubId"], data["hubName"] );
        hubs.set( hubId, hub )
        updateHubList( hubId )
        console.log( xhr.status + ":" + data["msg"]);
    }).fail( function( xhr, textStatus ){
        console.log( xhr.status + ":" + textStatus);
    })
    $("#hubnameInput").val("");
    $("#hubnameSettingModal").modal("hide");
}

function updateHubList( hubId ) {
    $("#hubList").empty();
    for( let [ id, hub ] of hubs.entries() ) {
        console.log( hub );
        let list = $("<a></a>").addClass( "nav-link" );
        list.data( "id", id );
        list.text( hub.name );
        list.appendTo( "#hubList" );

        if( hubId == id ) {
            list.addClass("active");
            lastHubListItem = list;
            $("#hubName").text( lastHubListItem.text() );
            $("#hubId").text( String(lastHubListItem.data("id")).padStart( 12, "0") );
            $("#dialog-container").html( hub.dialog );
        }

        list.click( function(){
            lastHubListItem.removeClass( "active" );
            $(this).addClass("active");
            lastHubListItem = $(this);
            $("#hubName").text( lastHubListItem.text() );
            $("#hubId").text( String(lastHubListItem.data("id")).padStart( 12, "0") );
            let hub = hubs.get( lastHubListItem.data("id") );
            $("#dialog-container").html( hub.dialog );
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
    return msgTimeStrToFormat((new Date()).toISOString());
}

function msgTimeStrToFormat( timeStr ) {
    let str = timeStr.replaceAll("-","/").replace("T", " ");
    str = str.substr( 0, str.lastIndexOf( ":" ) );
    return str;
}