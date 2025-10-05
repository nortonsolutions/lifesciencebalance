import { MessagesController } from "./app/messagesController.js";
import { AppController } from "./app/appController.js";
import { HomeController } from "./app/homeController.js";
import { NavController } from "/app/navController.js";
import { LoginController } from "./app/loginController.js";
import { GenericController } from "./app/genericController.js";
import { Router } from "./app/router.js";
class App {
  constructor() {
    this.router = new Router();
    this.addRouterLinks();

    this.messageController = new MessagesController();
    this.addWindowEventListeners();
  }
}

App.prototype.load = async function () {
  
  this.appController = new AppController();
  await this.appController.load();

  this.navController = new NavController(this.router);
  await this.navController.load({});
  await this.registerComponents();
  
  this.homeController = new HomeController(this.router);
  this.loginController = new LoginController(this.router);
  this.genericController = new GenericController(this.router);

  this.router.navigateTo(
    window.location.pathname,
    window.location.search,
    false
  );

};

App.prototype.addWindowEventListeners = function () {
  window.addEventListener("error", (e) => {
    // catches runtime errors, e.g. throw new Error(....);
    this.messageController.displayError(e);
  });

  window.addEventListener("popstate", (e) => {
    let { path, query } = e.state;
    this.router.navigateTo(path, query, true);
  });
};

App.prototype.addRouterLinks = function () {
  this.router.add("/", (request) => {
    this.loginController.load(request);
  });

  this.router.add("/app.html", (request) => {
    this.homeController.load(request);
  });

  // Add route for generic.html to handle various pages including admin
  this.router.add("/generic.html", (request) => {
    this.genericController.load(request);
  });
};

App.prototype.registerComponents = async function () {

  var componentsBase = "app/views/templates"

  var components = [
    "searchbar",
    "myCourses",
    "adSpace",
    "featuredCourse",
    "headerStuff",
  ]

  for (let component of components) {
    var template = await handleGet(
      componentsBase + "/" + component + ".hbs"
    );
    Handlebars.registerPartial(component, template);
  }

  registerHandlebarsHelpers();
};

export { App };
