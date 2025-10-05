/* Norton's rendition of the MVC Property Mgmt Page (June 2019) */
class HomeView {

    constructor() {
        this.dom = document.getElementById('mainPlaceholder');
        this.templateRetrieved = false;

        fetch("app/views/templates/home.hbs")
        .then(response => response.text())
        .then(text => {
          this.template = Handlebars.compile(text);
          this.templateRetrieved = true;
        });
    }
}


HomeView.prototype.render = async function(callback, context) {

    if (!context) var context = {};
    while (! this.templateRetrieved) {
        await delay();
    }

    context.indexActive = 'active';
    context.propertyActive = '';
    this.dom.innerHTML = this.template(context);

    document.querySelector('#seeMore').addEventListener('click', () => {
        document.querySelector('#more').classList.remove("inactive");
        document.querySelector('#seeMore').classList.add("inactive");
    });

    document.title = `Norton BlueSky Home`;
    window.scrollTo(0,0);
    callback();

}

export { HomeView };