/* Norton's rendition of the BlueSky Page (June 2019) */
import { PropertyStore } from "./propertyStore.js";
import { PropertyView } from "./views/propertyView.js";

class PropertyController {
    constructor(callback) {

        this.propertyStore = new PropertyStore();
        this.propertyView = new PropertyView(() => {
            callback();
        });
    }
}

PropertyController.prototype.load = function(request, callback) {

    let allProperties = this.propertyStore.getAll()

    var property = { };

    if (request.parameters[0]) {

        let params = getParamsFromRequest(request);
        property = this.propertyStore.get(params.propertyId);
    
    } 

    property.types = getFromProperties(allProperties, "type");
    property.cities = getFromProperties(allProperties, "city");

    this.propertyView.render(() => {

        if (callback) callback();
    
    }, property);
}

export { PropertyController };