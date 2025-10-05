/**
 * HomeController
 * @author Norton 2022
 */
import { HomeView } from "./views/homeView.js";
import { GenericView } from "./views/genericView.js";
import { EventDepot } from "./eventDepot.js";
export class HomeController {
  constructor(router) {
    this.eventDepot = new EventDepot();
    this.homeView = new HomeView(router, this.eventDepot);
    this.genericView = new GenericView(router, this.eventDepot);
    this.user = {};
  }

  async load(request) {
    let params = getParamsFromRequest(request);

    if (params?.user) {
      // decode the user using base-64 standard decoding
      this.user = JSON.parse(atob(params.user));
      localStorage.setItem("user", JSON.stringify(this.user));
      localStorage.setItem("userId", this.user.id);
    } else {
      // retrieve the user from localStorage and parse
      this.user = JSON.parse(localStorage.getItem("user"));
      
    }

    if (params?.page) {
        this.genericView.render(params.page, params.detail);
    } else {
      // temp sample featuredCourse data
      let featuredCourse = {
        title: "Learn to Code",
        description: "Learn to code in just a few weeks",
        imageUrl:
          "https://images.unsplash.com/photo-1518791841217-8f162f1e1131?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=800&q=60",
        url: "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
      };

      // temp sample ad data
      let ad = {
        title: "Buy something!",
        description: "Random advertisement",
        imageUrl: "https://source.unsplash.com/random",
        url: "https://random-ize.com/random-youtube/",
      };

      let context = {
        ad,
        featuredCourse,
        ...this.user,
        title: "Home",
        description: "Welcome to the home page",
      };

      this.homeView.render(context);
    }
  }
}
