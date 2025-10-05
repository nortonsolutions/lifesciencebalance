import { NewsletterListStore } from "./newsletterListStore.js";

/* Norton's rendition of the MVC Property Mgmt Page (June 2019) */
class Router {

    constructor() {
        this.routingTable = {};
        this.navigateTo = this.navigateTo.bind(this);
        this.setRouteLinks = this.setRouteLinks.bind(this);

    }
}

/* Add(path, handler) - adds a path to the router list with a handler 
function that takes a request object.  The request object should have a 
parameters Prop with any parameters provided by the path.  All 
parameters will be parsed from query string, i.e. /property.html?propertyId=4*/
Router.prototype.add = function(path, handler) {
    this.routingTable[path] = handler;
}

/* navigateTo(path, query) - calls the handler method from the router list 
and passes a request object.  This function will need to parse out any 
query string parameters and build the request object */
Router.prototype.navigateTo = function(path, query, pop) {

    let handler = this.routingTable[path];

    // Handle history state push
    if (!pop) {
        var historyState = { path: path, query: query };
        history.pushState(historyState, path + query, path + query);
    }

    // Strip the initial '?' and create an array of params,
    // include these in a request object
    var query = query.substring(1);
    let request = { path: path, parameters: query.split('&') };

    handler(request);
}

/* setRouteLinks - finds all the links (anchors) that need to be handled 
by the router and attaches a method to call navigateTo

Note: the HTMLAnchorElement has the following properties:
pathname - path from the href attribute (no query string)
search - query string including leading ? from the href attribute */
Router.prototype.setRouteLinks = function() {

    let anchorElements = document.querySelectorAll('a[data-route-link]');
    anchorElements.forEach(el => {
        el.addEventListener('click', e => {
            this.navigateTo(el.pathname, el.search, false);
            e.preventDefault();
        });
    });

    // Searchbar form submit behavior
    if (document.querySelector('#searchbar') != null) {
        document.querySelector('#searchbar').addEventListener('submit', e => {
        
            var pathname = e.currentTarget.action;
            pathname = "/" + pathname.split("/")[3];
    
            // Create query string based on form-control items
            var querystring = "?";
            e.currentTarget.querySelectorAll('.form-control').forEach(fe => {
                querystring += fe.name + "=" + fe.value + "&";
            });
            
            querystring = querystring.slice(0,querystring.length - 1);
            
            this.navigateTo(pathname, querystring, false);
            e.preventDefault();
        });
    }

    // newsletterForm form submit behavior
    if (document.querySelector('#newsletterForm') != null) {
        document.querySelector('#newsletterForm').addEventListener('submit', e => {

            this.newsletterListStore = new NewsletterListStore();
            var emailAddress = e.currentTarget.querySelector('.form-control').value;
            let newEntry = { id: this.newsletterListStore.uniqueId(), email: emailAddress };
            if (this.newsletterListStore.add(newEntry)) {
                alert("Thank you!  You have been added to our mailing list.")
            };
            e.preventDefault();
        });
    }
}


export { Router };