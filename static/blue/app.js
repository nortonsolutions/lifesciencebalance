/* Norton's rendition of the BlueSky website (June 2019) */

import { MessagesController } from "./app/messagesController.js";
import { AppController } from "./app/appController.js";
import { Router } from "./app/router.js";
import { HomeController } from "./app/homeController.js";
import { PropertyManagementController } from "./app/propertyManagementController.js";
import { PropertiesController } from "./app/propertiesController.js";
import { PropertyController } from "./app/propertyController.js";
import { ContactController } from "./app/contactController.js";

class App {
  constructor() {}
}

App.prototype.load = function() {
  this.router = new Router();

  this.appController = new AppController(() => {
    this.appController.load(() => {
      this.registerPartials(() => {
        this.messageController = new MessagesController();
        this.homeController = new HomeController();
        this.propertyManagementController = new PropertyManagementController(
          () => {
            this.propertiesController = new PropertiesController(() => {
              this.propertyController = new PropertyController(() => {
                this.contactController = new ContactController(() => {
                  this.addRouterLinks();
                  this.addWindowEventListeners();
                  this.router.navigateTo(
                    window.location.pathname,
                    window.location.search,
                    false
                  );
                });
              });
            });
          }
        );
      });
    });
  });
};

App.prototype.addWindowEventListeners = function() {
  window.addEventListener("error", e => {
    // catches runtime errors, e.g. throw new Error(....);
    this.messageController.displayMessage(e.error);
  });

  window.addEventListener("popstate", e => {
    let { path, query } = e.state;
    this.router.navigateTo(path, query, true);
  });

  window.addEventListener('scroll', () => {
    if (document.body.scrollTop > 1 || document.documentElement.scrollTop > 1) {
      document.getElementsByTagName("nav")[0].style.maxHeight = "100px";
      document.getElementsByTagName("nav")[0].style.padding = "10px";
    } else {
      document.getElementsByTagName("nav")[0] .style.maxHeight = "120px";
      document.getElementsByTagName("nav")[0].style.padding = "20px";
    }
  });
};

App.prototype.addRouterLinks = function() {
  this.router.add("/", request => {
    this.homeController.load(() => {
      this.router.setRouteLinks();
    });
  });

  this.router.add("/home.html", request => {
    this.homeController.load(() => {
      this.router.setRouteLinks();
    });
  });

  this.router.add("/index.html", request => {
    this.homeController.load(() => {
      this.router.setRouteLinks();
    });
  });

  this.router.add("/educationMaterials/workbenchProjects/review/bluesky-website-36a8b4d0fda447acb4971b9d/index.html", request => {
    this.homeController.load(() => {
      this.router.setRouteLinks();
    });
  });

  this.router.add("/property-management.html", request => {
    this.propertyManagementController.load(request, () => {
      this.router.setRouteLinks();
    });
  });

  this.router.add("/properties.html", request => {
    this.propertiesController.load(request, () => {
      this.router.setRouteLinks();
    });
  });

  this.router.add("/property.html", request => {
    this.propertyController.load(request, () => {
      this.router.setRouteLinks();
    });
  });

  this.router.add("/contact.html", request => {
    this.contactController.load(request, () => {
      this.router.setRouteLinks();
    });
  });
};

App.prototype.registerPartials = function(callback) {
  fetch("app/views/templates/navbar.hbs")
    .then(r => r.text())
    .then(text => {
      Handlebars.registerPartial("navbar", text);

      fetch("app/views/templates/searchbar.hbs")
        .then(r => r.text())
        .then(text => {
          Handlebars.registerPartial("searchbar", text);
          
          fetch("app/views/templates/footer.hbs")
          .then(r => r.text())
          .then(text => {
            Handlebars.registerPartial("footer", text);
            callback();
          });
        });
    });
};

export { App };
