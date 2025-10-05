/**
 * AppView
 * @author Norton 2022
 */
import { AppView } from './views/appView.js'

export class AppController {

    constructor() {
        this.appView = new AppView();
    }

    async load() {
        await this.appView.render();
    }
}
