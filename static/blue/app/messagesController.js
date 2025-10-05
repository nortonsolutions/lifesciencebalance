/* Norton's rendition of the MVC Property Mgmt Page (June 2019) */
import { Messages } from "./views/messagesView.js"

class MessagesController {

    constructor() {
        this.messages = new Messages(this);
    }
}

MessagesController.prototype.displayMessage = function(message) {
    this.messages.updateMessage(message);
}


export { MessagesController };