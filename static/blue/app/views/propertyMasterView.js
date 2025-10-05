/* Norton's rendition of the MVC Property Mgmt Page (June 2019) */
class PropertyMasterView {

    constructor(eventDepot) {

        this.eventDepot = eventDepot;
        this.dom = document.getElementById('propertyMaster');
        this.templateRetrieved = false;

        fetch("app/views/templates/propertyMaster.hbs").then(response => response.text()).then(text => {
            this.template = Handlebars.compile(text);
            this.templateRetrieved = true;
        });
    }
}


PropertyMasterView.prototype.resetPropertyRows = async function(rows) {

    while (! this.templateRetrieved) {
        await delay();
    }

    let context = { 'properties': rows };    
    document.getElementById('propertyMasterPlaceholder').innerHTML = this.template(context);

    document.getElementById("btnNew").addEventListener("click", e => { this.eventDepot.fire('newProperty', e); });

    document.querySelectorAll(".deleteButton").forEach(el => {
        el.addEventListener("click", e => { this.eventDepot.fire('deleteProperty', e); });
    });

    document.title = "Property Master View";
}



export { PropertyMasterView };