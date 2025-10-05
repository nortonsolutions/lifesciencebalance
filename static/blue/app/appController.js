/* Norton's rendition of the MVC Property Mgmt Page (June 2019) */
import { AppView } from './views/appView.js'

class AppController {

    constructor(callback) {
        this.appView = new AppView(() => { callback() });
    }
}

AppController.prototype.load = function(callback) {
    this.appView.render(() => {
        callback();
    });
}

export { AppController };