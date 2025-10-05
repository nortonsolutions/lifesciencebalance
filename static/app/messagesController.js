/**
 * MessagesView
 * @author Norton 2022
 */

export class MessagesController {
  constructor() {}

  displayError({ error, message }) {
    switch (error) {
      case "Session expired":
        this.swal(error, message + "<br> <b></b>", "error", 2000, () => {
            window.location.href = "/index.html";
        });
        break;

      default:
        break;
    }
  }

  swal(title, html, icon, timer=0, callback = () => {}) {
    let timerInterval;
    Swal.fire({
      title,
      html,
      timer: timer > 0 ? timer : null,
      icon,
      timerProgressBar: true,
      didOpen: () => {
        Swal.showLoading();
        const b = Swal.getHtmlContainer().querySelector("b");
        timerInterval = setInterval(() => {
          b.textContent = Swal.getTimerLeft();
        }, 100);
      },
      willClose: () => {
        clearInterval(timerInterval);
      },
    }).then((result) => {
      if (result.dismiss === Swal.DismissReason.timer) {
        callback();
      }
    });
  }
}
