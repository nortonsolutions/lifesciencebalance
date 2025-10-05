class PropertyView {
    constructor(callback) {
        
        this.dom = document.getElementById('mainPlaceholder');
        
        fetch("app/views/templates/property.hbs").then(response => response.text()).then(text => {
            this.template = Handlebars.compile(text);
            callback();
        });
         
    }
}

PropertyView.prototype.render = function(callback, params) {
    
    if (!params) var params = {};
    params.indexActive = '';
    params.propertiesActive = 'active';
    this.dom.innerHTML = this.template(params);
    document.title = `Norton BlueSky - ${params.title}`;
    document.documentElement.scrollTop = 0;
    window.scrollTo(0, 500);
    callback();
}

export { PropertyView };