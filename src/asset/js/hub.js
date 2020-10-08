class Hub {
    constructor( id, name ) {
        this.id = id;
        this.name = name;
        let dialog = $("<div style='overflow-y:auto;' class='p-2 rounded bg-white h-100'></div>");
        this.dialog = dialog;
    }

    appendMessage( type, name, time, msg ) {
        let messageBox = $("<div></div>").addClass("bg-light");
    
        switch( type ) {
            case BROADCAST:
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
        wrapper.appendTo(this.dialog);

        this.dialog.prop( "scrollTop", this.dialog.prop('scrollHeight') - this.dialog.prop('clientHeight'))
    }
}