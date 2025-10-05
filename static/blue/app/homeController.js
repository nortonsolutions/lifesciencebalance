/* Norton's rendition of the MVC Property Mgmt Page (June 2019) */
import { HomeView } from "./views/homeView.js";
import { PropertyStore } from "./propertyStore.js";
import { EventDepot } from "./eventDepot.js";

class HomeController {

    constructor() {
        this.propertyStore = new PropertyStore();
        this.eventDepot = new EventDepot();
        this.homeView = new HomeView(this.eventDepot);
    }
}

HomeController.prototype.load = function(callback) {

    let allProperties = this.propertyStore.getAll()

    if (Object.keys(allProperties).length == 0) {
        allProperties = loadDefaultProperties();
        this.propertyStore.overwriteStorage(allProperties);
    }

    allProperties = shuffle(allProperties);

    let context = {};
    // Featured (offer) property
    if (allProperties.find(el => el.status == 'offer') != undefined) 
        context = allProperties.find(el => el.status == 'offer');

    // Sort the allProperties array in reverse by dateUpdated; take the top three
    allProperties = allProperties.sort((a, b) => a.dateUpdated > b.dateUpdated);
    let someProperties = allProperties.filter((el, index) => index < 3);

    if (allProperties.length > 0) {
        context.types = getFromProperties(allProperties, "type");
        context.cities = getFromProperties(allProperties, "city");

        let citiesWithCodes = [];

        for (let city of context.cities) {
            citiesWithCodes.push({city: city, cityUrlSafe: city.replace(/ /g,"+"), code: cityMatrix[city]});
        }

        context.properties = someProperties;
        context.citiesWithCodes = citiesWithCodes;
        // context.testimonials = testimonials;
    }

    this.homeView.render(() => {

        if (callback) callback();
    
    }, context);

}



export { HomeController };