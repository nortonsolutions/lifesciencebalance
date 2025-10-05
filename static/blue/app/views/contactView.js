class ContactView {
    constructor(callback) {
        
        this.dom = document.getElementById('mainPlaceholder');

        fetch("app/views/templates/contact.hbs").then(response => response.text()).then(text => {
            this.template = Handlebars.compile(text);
            callback();
        });
         
    }
}

ContactView.prototype.render = function(callback, params) {
    
    if (!params) var params = {};
    params.indexActive = '';
    params.propertiesActive = '';
    params.contactActive = 'active';
    this.dom.innerHTML = this.template(params);
    document.title = `Norton BlueSky - Contact`;
    document.documentElement.scrollTop = 0;
    window.scrollTo(0, 420);
    // document.documentElement.scrollTop = 200;
    callback();
}

export { ContactView };