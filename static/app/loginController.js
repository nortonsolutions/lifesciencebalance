/**
 * LoginController
 * @author Norton 2022
 */
import { LoginView } from "./views/loginView.js";
import { EventDepot } from "./eventDepot.js";
export class LoginController {
  constructor(router) {
    this.eventDepot = new EventDepot();
    this.loginView = new LoginView(router);
  }

  load(request) {
    var referral;
    let params = getParamsFromRequest(request);

    if (params?.referral) {
        referral = params.referral;
    }

    let context = {
      referral,
      title: "Login",
      description: "Welcome to the login page",
    };

    this.loginView.render(context);
  }
}
