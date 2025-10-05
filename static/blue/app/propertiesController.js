/* Norton's rendition of the BlueSky Page (June 2019) */
import { PropertyStore } from "./propertyStore.js";
import { PropertiesView } from "./views/propertiesView.js";

class PropertiesController {
  constructor(callback) {
    this.propertyStore = new PropertyStore();
    this.propertiesView = new PropertiesView(() => {
      callback();
    });
  }
}

PropertiesController.prototype.load = function(request, callback) {
  let allProperties = this.propertyStore.getAll();
  let context = { properties: allProperties };
  context.types = getFromProperties(allProperties, "type");
  context.cities = getFromProperties(allProperties, "city");

  var params = {};

  if (request) {
    params = getParamsFromRequest(request);

    if (Object.keys(params).length > 0) {
      for (let paramName in params) {
        if (params[paramName] != "null") {
          context.properties = context.properties.filter(property => {
            return property[paramName] == params[paramName];
          });
        }
      }
    }
  }

  this.propertiesView.render(() => {
    if (callback) callback();
  }, context);
};

export { PropertiesController };
