var conn;
var userId;
var hubs = new Map();
var lastHubListItem;
var invitesCount = 0;
var invitesList = [];

window.onbeforeunload = function () {
    conn.close(1000);
};

$(document).ready(function(){
    let queryStr = window.location.search
    const urlParams = new URLSearchParams( queryStr )
    userId = urlParams.get( "userId" );
    $("#userId").text( userId )

    $.ajax({
        url: document.location.protocol + "//" + document.location.host + "/hub/" + userId,
        type: "GET",
        dataType: "json"
    }).then( function( hubLists, textStatus, xhr ) {
        for( hubInfo of hubLists ) {
            let hub = new Hub( hubInfo.hubId, hubInfo.hubName );
            hubs.set( hubInfo.hubId, hub);
            refreshHubHistoryMsg( hubInfo.hubId, hub );
            setTimeout(refreshOnlineUser, 1000, hubInfo.hubId, hub );
        }

        if (hubLists)
            updateHubList( hubLists[0].hubId);
    }).fail( function( xhr, textStatus ) {
        console.log( xhr.status + ":" + textStatus );
    });

    updateInviteBox();
    createConn( userId );

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

function updateInviteBox() {
    $.ajax({
        url: document.location.protocol + "//" + document.location.host + "/invite/" + userId,
        type: "GET",
        dataType: "json"
    }).then( function( inviteMsgs, textStatus, xhr){
        if (inviteMsgs == null) {
            return;
        }

        $("#inviteList").empty();
        invitesCount = 0;

        for( inviteInfo of inviteMsgs ) {
            console.log( inviteInfo )
            let dataRow = $("<tr></tr>");
            $("<td></td>").text( msgTimeStrToFormat(inviteInfo.time) ).appendTo( dataRow );
            $("<td></td>").text( `${inviteInfo.userName} 邀請你加入 ${inviteInfo.hubName} 聊天室` ).appendTo( dataRow );

            let accept = $("<a class='text-success'>接受</a>");
            let reject = $("<a class='text-danger'>拒絕</a>")
            $("<td class='text-right'></td>").append( accept )
                                            .append( "|" )
                                            .append( reject )
                                            .appendTo( dataRow );
            
            let acceptHandler = function() {
                let index = $("#inviteList tr").index($(this).closest("tr"));
                $(this).closest("tr").remove();

                let  inviteMsg = invitesList[ index ];
                if ( !conn )
                    return;
                
                let answer = {
                    action: ANSWER,
                    hubId: inviteMsg.hubId,
                    hubName: inviteMsg.hubName,
                    userId: userId,
                    userName: "",
                    content: "1"
                }

                answer.content = "1";
                let hubId = inviteMsg.hubId;
                let hub = new Hub( hubId, inviteMsg.hubName );
                conn.send( JSON.stringify( answer));
                refreshHubHistoryMsg( hubId, hub );
                updateHubList( hubId );
                updateInviteBox();
            }

            let rejectHandler = function () {
                let index = $("#inviteList tr").index($(this).closest("tr"));
                $(this).closest("tr").remove();
                let  inviteMsg = invitesList[ index ];

                if ( !conn )
                    return;

                let answer = {
                    action: ANSWER,
                    hubId: inviteMsg.hubId,
                    hubName: inviteMsg.hubName,
                    userId: userId,
                    userName: "",
                    content: "0"
                }

                conn.send( JSON.stringify( answer));
                updateInviteBox();
            }

            accept.click( acceptHandler );
            reject.click( rejectHandler ); 

            $("#inviteList").append( dataRow );
            invitesList.push( inviteInfo )
            invitesCount++;
        }
        $("#inviteCount").text( invitesCount );
    }).fail( function( xhr, textStatus) {
        console.log( xhr.status + ":" + textStatus )
    })
}

function handleMessage( message ) {
    let hub;
    let type;
    switch( message.action ) {
        case MESSAGE:
            type = OTHER;
            if( message.userId == userId ) {
                type = USER;
            }
            hub = hubs.get( message.hubId );
            hub.appendMessage( type, message.userName, msgTimeStrToFormat(message.time), message.content );
            break;
        case INVITE:
            updateInviteBox();
            break;
        case USER_ONLINE:
            hub = hubs.get( message.hubId );
            hub.userOnline( message.userId, message.userName );
            break;
        case USER_OFFLINE:
            hub = hubs.get( message.hubId );
            hub.userOffline( message.userId, message.userName );
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
        setTimeout( refreshOnlineUser, 3000, hubId, hub );
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
            $("#userListUIContainer").html( hub.aliveUserListUI );
        }

        list.click( function(){
            lastHubListItem.removeClass( "active" );
            $(this).addClass("active");
            lastHubListItem = $(this);
            $("#hubName").text( lastHubListItem.text() );
            $("#hubId").text( String(lastHubListItem.data("id")).padStart( 12, "0") );
            let hub = hubs.get( lastHubListItem.data("id") );
            $("#dialog-container").html( hub.dialog );
            $("#userListUIContainer").html( hub.aliveUserListUI );
        })                
    }
}

function refreshHubHistoryMsg( hubId, hub ) {
    $.ajax({
        type: "GET",
        url: document.location.protocol + "//" + document.location.host + "/history/" + hubId,
    }).then( function( historyMessages, textStatus, xhr){
        for( message of historyMessages ) {
            let type = USER;
            if( message.userId != userId )
                type = OTHER;
            hub.appendMessage( type, message.userName, msgTimeStrToFormat( message.time), message.content);
        }
    }).fail( function( xhr, textStatus ){
        console.log( `Load Hub:${hubId} messages failed in function refreshHubHistoryMsg.`);
    })
}

function refreshOnlineUser( hubId, hub ) {
    $.ajax({
        type: "GET",
        url: document.location.protocol + "//" + document.location.host + "/member/" + hubId,
    }).then( function( userList, textStatus, xhr){
        console.log( userList );
        hub.updateAliveUserList( userList );
    }).fail( function( xhr, textStatus ){
        console.log( `Load Hub:${hubId} userList failed in function refreshOnlineUser.`);
    })
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
    
    var urlPattern = /(http|https):\/\/(\w+\.)+\w+\/?/g;
    var urlReplaceStr = str.replaceAll( urlPattern, function replacer( match ) {
        return `<a href="${match}">${match}</a>`;
    });

    messageBox.html( `<small>${name} ${time}</small><div>${urlReplaceStr}</div>`);

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