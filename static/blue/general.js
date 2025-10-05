/* Norton's rendition of the MVC Property Mgmt Page (June 2019) */

/* Simple 10ms delay, useful for initial loop when 
awaiting an asynchronous load of some resource.
Call from within an async function, i.e., 
var result = await delay(); */
function delay() {
  return new Promise(resolve => {
    setTimeout(() => {
      resolve("resolved");
    }, 10);
  });
}

function getParamsFromRequest(request) {

  let { parameters } = request;

  let paramObj = {};

  if (request.parameters[0]) {
    for (let el of parameters) {
      let [ name, value ] = el.split("=");

      value = value.replace(/\+/g," ");
      if (value != "null") paramObj[name] = value;
    }

  }
  return paramObj;

}

function getFromProperties(allProperties, attrName) {

  let attrArray = [];
  for (let p of allProperties) {
      if (!attrArray.includes(p[attrName])) attrArray.push(p[attrName]);
  }

  return attrArray;

}

// Borrowed from D3 code :-)
let shuffle = function(array, i0, i1) {
  if ((m = arguments.length) < 3) {
    i1 = array.length;
    if (m < 2) i0 = 0;
  }
  var m = i1 - i0, t, i;
  while (m) {
    i = Math.random() * m-- | 0;
    t = array[m + i0], array[m + i0] = array[i + i0], array[i + i0] = t;
  }
  return array;
};



/* Register Handlebars helpers for views */
function registerHandlebarsHelpers() {
  Handlebars.registerHelper("maleCheck", function(data) {
    let str = "";
    if (data === "male") str = "checked";
    return new Handlebars.SafeString(str);
  });

  Handlebars.registerHelper("femaleCheck", function(data) {
    let str = "";
    if (data === "female") str = "checked";
    return new Handlebars.SafeString(str);
  });

  Handlebars.registerHelper("roleCompare", function(role, compareRole) {
    let str = "";
    if (role === compareRole) str = "selected";
    return new Handlebars.SafeString(str);
  });

  Handlebars.registerHelper("verifiedHelper", function(verified) {
    let str = "";
    if (verified) str = "checked";
    return new Handlebars.SafeString(str);
  });

  Handlebars.registerHelper("iconize", function(data) {
    let str;

    switch (data) {
      case "male":
        str = "<i class='fas fa-male'></i>";
        break;
      case "female":
        str = "<i class='fas fa-female'></i>";
        break;
      default:
        str = "?";
    }
    return new Handlebars.SafeString(str);
  });
}

function loadDefaultProperties() {
    return defaultProperties;
}

registerHandlebarsHelpers();



let defaultProperties = [{"id":0,"title":"Lovely cottage","description":"Quaint and awesome....","city":"Salt Lake City","price":"2341333","type":"Townhouse","status":"standard","rentable":"on","bedrooms":"3","bathrooms":"3","area":"2344","patios":"2","active":"on","imageFilename":"property_1.jpg","dateCreated":"2019-06-10T21:09:19.280Z","dateUpdated":"2019-06-10T21:09:19.280Z","":""},{"id":2,"title":"Wowzers","description":"Draw-dropping beauty","city":"Compton","price":"3453123","type":"Townhouse","status":"offer","rentable":"on","bedrooms":"3","bathrooms":"3","area":"3132","patios":"2","active":"on","imageFilename":"property_4.jpg","dateCreated":"2019-06-10T21:14:43.167Z","dateUpdated":"2019-06-10T21:14:43.167Z","":""},{"id":3,"title":"City dreams","description":"High-rise views of heaven","city":"Paris","price":"2341533","type":"City Property","status":"offer","rentable":"on","bedrooms":"5","bathrooms":"3","area":"4441","patios":"2","active":"on","imageFilename":"property_5.jpg","dateCreated":"2019-06-10T21:15:44.247Z","dateUpdated":"2019-06-10T21:15:44.247Z","":""},{"id":1,"title":"Amazing Backcountry","description":"Fit for kings and queens","city":"Petaluma","price":"3452135","type":"Vacation Home","status":"offer","rentable":"on","bedrooms":"3","bathrooms":"4","area":"3215","patios":"2","active":"on","imageFilename":"property_3.jpg","dateCreated":"2019-06-10T21:14:09.271Z","dateUpdated":"2019-06-10T21:16:05.967Z","":""}];

let cityMatrix = {

  "Los Angeles": 1,
  "Salt Lake City": 2,
  "Chicago": 3,
  "London": 4,
  "Petaluma": 5,
  "Ogden": 6,
  "Outback": 7,
  "Venice": 8,
  "Compton": 2,
  "Paris": 3

};

let testimonials = [

  { 
    title: "Amazing home for me", 
    testimonial: "Lorem ipsum dolor sit amet consectetur adipisicing elit. Animi veritatis exercitationem sapiente temporibus aut!", 
    image: "testimonial_1.jpg",
    name: "Diane Smith",
    role: "client"
  },
  { 
    title: "Friendly realtors", 
    testimonial: "Lorem ipsum dolor sit amet consectetur adipisicing elit. Animi veritatis exercitationem sapiente temporibus aut!", 
    image: "testimonial_2.jpg",
    name: "Michael Duncan",
    role: "client"
  },
  { 
    title: "Very good communication", 
    testimonial: "Lorem ipsum dolor sit amet consectetur adipisicing elit. Animi veritatis exercitationem sapiente temporibus aut!", 
    image: "testimonial_3.jpg",
    name: "Shawn Gaines",
    role: "client"
  }

];