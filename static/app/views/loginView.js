/**
 * LoginView is for non-dynamic pages like Contact and About;
 * just grabs the HBS template and places it in the main;
 * when rendering provide templateName and optional context params.
 */
class LoginView {
  constructor(router) {
    this.router = router;
  }

  async render(context) {

    this.template = await getTemplate(`app/views/templates/login.hbs`);

    this.dom = document.getElementById("main");
    this.dom.innerHTML = this.template(context);
    this.load();

    this.router.setRouteLinks(this.dom);
    document.title = `nortonApp - Login`;
    // document.documentElement.scrollTop = 0;
    window.scrollTo(0, 0);
  }

  async load() {
    localStorage.clear();

    // clear out the forms
    $("#localLogin").trigger("reset");
    $("#newUser").trigger("reset");

    window.addEventListener("error", (e) => {
      swalAlert(e);
    });

    document
      .getElementById("localLogin")
      .addEventListener("submit", function (e) {
        e.preventDefault();

        var body = {
          username: e.target.elements.username.value,
          password: e.target.elements.password.value,
        };

        fetch("/login", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify(body),
        })
          .then(function (response) {
            switch (response.status) {
              case 200: // success
                response.json().then((json) => {
                  window.localStorage.setItem(
                    "user",
                    JSON.stringify(json.user)
                  );
                  window.location.href = json.redirect_path;
                });
                break;
              case 401: // unauthorized
                throw "Invalid username or password";
              default: // error
                throw "Error: " + response.status;
            }
          })
          .catch(function (error) {
            window.dispatchEvent(
              new ErrorEvent("error", {
                error: error,
                message: error.message,
              })
            );
          });
      });

    document.getElementById("newUser").addEventListener("submit", function (e) {
      e.preventDefault();
      var body = {
        username: e.target.elements.username.value,
        email: e.target.elements.email.value,
        firstName: e.target.elements.firstName.value,
        lastName: e.target.elements.lastName.value,
        password: e.target.elements.password.value,
        confirmPassword: e.target.elements.confirmPassword.value,
      };
      fetch("/user", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(body),
      })
        .then(function (response) {
          switch (response.status) {
            case 200:
              response.json().then((json) => {
                localStorage.setItem("user", JSON.stringify(json.user));
                localStorage.setItem("userId", json.user.id);
                window.location.href = json.redirect_path;
              });
              break;
            case 400: // bad request
              throw "Error: " + response.status;
            default: // error
              throw "Error: " + response.status;
          }
        })
        .catch(function (error) {
          window.dispatchEvent(
            new ErrorEvent("error", {
              error: error,
              message: error.message,
            })
          );
        });
    });
  }
}

export { LoginView };
