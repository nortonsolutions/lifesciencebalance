import { ContactView } from "./views/contactView.js";

class ContactController {
    constructor(callback) {

        this.contactView = new ContactView(() => {
            callback();
        });
    }
}

ContactController.prototype.load = function(request, callback) {

    if (request.parameters[0]) {
        let params = getParamsFromRequest(request);
    } 

    this.contactView.render(() => {
        if (callback) callback();
    });
}

export { ContactController };