/* Norton's rendition of the MVC Property Mgmt Page (June 2019) */
class AppView {

    constructor(callback) {
  
      this.preLoad("http://redrockcodecamp.org/educationMaterials/frameworks/js/handlebars/4.0.12/handlebars.min.js", () => {
          this.preLoad("general.js", () => {
  
          // After HB is loaded, get the header
          fetch("app/views/templates/appHead.html").then(response => response.text()).then(text => {
              let t = document.createElement('template');
              t.innerHTML = text;
              document.head.append(t.content);
              callback();
            });
          });
      });
    }
  }
  
  AppView.prototype.render = function(callback) {
  
      let messagesPlaceholder = document.createElement("div");
      messagesPlaceholder.id = "messagesPlaceholder"
  
      let mainPlaceholder = document.createElement("div");
      mainPlaceholder.id = "mainPlaceholder"
      
      document.body.appendChild(messagesPlaceholder);
      document.body.appendChild(mainPlaceholder);
      callback();
      
  };
  
  AppView.prototype.preLoad = function (scripting, callback) {
      var script = document.createElement("script");
      script.src = scripting;
      script.addEventListener("load", function() {
        callback();
      });
      document.body.appendChild(script);
  }
  
  export { AppView };
  