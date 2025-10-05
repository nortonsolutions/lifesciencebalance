class PropertiesView {
    constructor(callback) {
        
        this.dom = document.getElementById('mainPlaceholder');
        this.templateRetrieved = false;

        fetch("app/views/templates/properties.hbs").then(response => response.text()).then(text => {
            this.template = Handlebars.compile(text);
            this.templateRetrieved = true;
            callback();
        });
         
    }
}

PropertiesView.prototype.render = function(callback, params) {
    
    if (!params) var params = {};
    params.indexActive = '';
    params.propertiesActive = 'active';
    this.dom.innerHTML = this.template(params);
    document.title = "Norton BlueSky Properties";
    document.documentElement.scrollTop = 0;
    window.scrollTo(0, 220);
    callback();
}

export { PropertiesView };