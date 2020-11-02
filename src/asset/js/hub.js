class Hub {
    constructor( id, name ) {
        this.id = id;
        this.name = name;
        this.aliveUser = new Map();
        this.aliveUserListUI = $("<ul class='list-group'></ul>");
        let dialog = $("<div style='overflow-y:auto;' class='p-2 rounded bg-white h-100'></div>");
        this.dialog = dialog;
    }

    appendMessage( type, name, time, msg ) {
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

        var urlReplaceStr = msg.replaceAll( urlPattern, function replacer( match ) {
            return `<a href="${match}" target="_blank">${match}</a>`;
        });

        messageBox.html( `<small>${name} ${time}</small><div>${urlReplaceStr}</div>`);
    
        let wrapper = $("<div></div>").addClass("row")
                                    .addClass("mt-2")
                                    .append( messageBox );
        wrapper.appendTo(this.dialog);

        this.dialog.prop( "scrollTop", this.dialog.prop('scrollHeight') - this.dialog.prop('clientHeight'))
    }

    updateAliveUserList( userList ) {
        for( let userInfo of userList ) {
            let userInfoUI = $("<li class='list-group-item'></li>");
            let badge = $("<div class='badge text-wrap'></div>").appendTo( userInfoUI );
            let userIdUI = $("<span class='bg-info'></span>").appendTo( userInfoUI );
            let userNameUI = $("<span></span>").appendTo( userInfoUI );

            userIdUI.text( userInfo.userId );
            userNameUI.text( userInfo.userName );

            this.aliveUser.set( userInfo.userId, userInfoUI );

            if ( userInfo.active ) {
                badge.addClass("badge-success").text("上線中");
                this.aliveUserListUI.prepend( userInfoUI );
            } else {
                badge.addClass("badge-secondary").text("已離線");
                this.aliveUserListUI.append( userInfoUI );
            }
        }
        
    }

    userOnline( userId, userName ) {
        let userInfoUI = this.aliveUser.get( userId );

        userInfoUI.empty();
        let badge = $("<div class='badge badge-success text-wrap'></div>").appendTo( userInfoUI );
        let userIdUI = $("<span class='bg-info'></span>").appendTo( userInfoUI );
        let userNameUI = $("<span></span>").appendTo( userInfoUI );
        
        badge.addClass("badge-success").text("上線中");
    
        userIdUI.text( userId );
        userNameUI.text( userName );

        this.aliveUserListUI.remove( userInfoUI );
        this.aliveUserListUI.prepend( userInfoUI );
    }

    userOffline( userId, userName ) {
        let userInfoUI = this.aliveUser.get( userId );

        userInfoUI.empty();
        let badge = $("<div class='badge badge-secondary text-wrap'></div>").appendTo( userInfoUI );
        let userIdUI = $("<span class='bg-info'></span>").appendTo( userInfoUI );
        let userNameUI = $("<span></span>").appendTo( userInfoUI );
        
        badge.addClass("badge-secondary").text("已離線");
    
        userIdUI.text( userId );
        userNameUI.text( userName );

        this.aliveUserListUI.remove( userInfoUI );
        this.aliveUserListUI.append( userInfoUI );
    }
}