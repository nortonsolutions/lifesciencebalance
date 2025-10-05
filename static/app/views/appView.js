/**
 * appView
 * @author Norton 2022
 */
class AppView {
  
  constructor() {}

  render = async function () {

    await this.preLoad("/cdn/handlebars.min.js");
    await this.preLoad("https://cdnjs.cloudflare.com/ajax/libs/font-awesome/5.15.4/js/all.min.js");
    await this.preLoad("/cdn/bootstrap.bundle.min.js");
    await this.preLoad("/cdn/sweetalert2.all.min.js");
    await this.preLoad("/cdn/jquery-3.6.0.min.js");
    await this.preLoad("/general.js");
    await this.preLoad("/public/common.js");

    let t = document.createElement("template");
    let text = await handleGet("app/views/templates/appHead.html");
    t.innerHTML = text;
    document.head.append(t.content);

  };

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

export { AppView };
