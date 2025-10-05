/* Norton's rendition of the MVC Property Mgmt Page (June 2019) */
class Messages {

    constructor() { }

}

Messages.prototype.updateMessage = function(message) {

    let context = { 'hidden': false, 'message': message };

    fetch("app/views/templates/messages.hbs").then(response => response.text()).then(text => {
        let template = Handlebars.compile(text);
        document.getElementById('messagesPlaceholder').innerHTML = template(context);
    });
    
}

export { Messages };
