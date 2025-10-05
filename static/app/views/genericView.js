/**
 * GenericView is for non-dynamic pages like Contact and About;
 * just grabs the HBS template and places it in the main;
 * when rendering provide templateName and optional context params.
 */

import { ModuleView } from "./moduleView.js";
class GenericView {
  constructor(router, eventDepot) {
    this.router = router;
    this.moduleView = new ModuleView(router, eventDepot);
  }

  async render(templateName, detail, dom = document.getElementById("main")) {
    this.template = await getTemplate(
      `app/views/templates/${templateName}.hbs`
    );

    var context;
    switch (templateName) {
      case "course":
        localStorage.setItem("courseId", detail);
        context = JSON.parse(await handleGet(`/course/${detail}`));
        dom.innerHTML = this.template(context);

        // recurse for modules and instructors (sub pages to be rendered)
        let domModules = document.getElementById("modules");
        await this.render("modules", detail, domModules);

        let domInstructors = document.getElementById("instructors");
        await this.render("instructors", detail, domInstructors);
        window.scrollTo(0, 0);
        break;

      case "modules":
        context = JSON.parse(await handleGet(`/course/${detail}/module`));
        dom.innerHTML = this.template(context);
        break;

      case "instructors":
        context = JSON.parse(await handleGet(`/course/${detail}/instructor`));
        dom.innerHTML = this.template(context);
        break;

      case "module":
        // pass off to module view
        await this.moduleView.render(detail, dom);
        break;

      default:
        dom.innerHTML = this.template(context);
        break;
    }

    this.router.setRouteLinks(dom);
    // document.title = `nortonApp - ${templateName}`;
    // document.documentElement.scrollTop = 0;
  }

  preLoad = function (scripting) {
    return new Promise((resolve, reject) => {
      var script = document.createElement("script");
      script.src = scripting;
      script.addEventListener("load", resolve);
      script.addEventListener("error", reject);
      document.body.appendChild(script);
    });
  };
}

export { GenericView };
