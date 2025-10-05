/* Norton's rendition of the MVC Property Mgmt Page (June 2019) */
class PropertyManagementView {

    constructor(callback, eventDepot) {

        this.eventDepot = eventDepot;
        this.siteConfigDisplayed = false;
        this.dom = document.getElementById('mainPlaceholder');

        fetch("app/views/templates/propertyManage.hbs").then(response => response.text()).then(text => {
            this.template = Handlebars.compile(text);

            fetch("app/views/templates/siteConfiguration.hbs").then(response => response.text()).then(text => {
                this.configTemplate = Handlebars.compile(text);

                callback();
            });
        });
    }

    render(callback, propertyObj, configObj) {

        if (!propertyObj) var propertyObj = {};
        propertyObj.indexActive = '';
        propertyObj.propertyActive = '';
        this.dom.innerHTML = this.template(propertyObj);

        if (!configObj) var configObj = {};

        let domSiteConfig = document.getElementById('siteConfiguration');
        domSiteConfig.innerHTML = this.configTemplate(configObj);
        
        this.addEventListeners();
        
        callback();
    }

    toggleView() {
        
        if (! this.siteConfigDisplayed) {
    
            document.querySelector('#siteConfiguration').classList.remove('inactive');
            document.querySelector('#propertyDetailPlaceholder').classList.add('inactive');
            document.querySelector('#propertyMasterPlaceholder').classList.add('inactive');

        } else {

            document.querySelector('#siteConfiguration').classList.add('inactive');
            document.querySelector('#propertyDetailPlaceholder').classList.remove('inactive');
            document.querySelector('#propertyMasterPlaceholder').classList.remove('inactive');

        }

        this.siteConfigDisplayed = !this.siteConfigDisplayed;
    }

    addEventListeners() {

        document.querySelector('#siteConfigToggle').addEventListener('click', e => {
            this.toggleView();
        });

        document.querySelector('#addToMailingList').addEventListener('click', e => {
            
        });

        document.querySelectorAll('.deleteEmail').forEach(el => { 
            el.addEventListener('click', e => {
                let index = e.currentTarget.id.split('-')[1];
                this.eventDepot.fire('deleteNewsletterListItem', { index: index });
            });
        });
    }

    reloadSiteConfiguration(configObj) {
        let domSiteConfig = document.getElementById('siteConfiguration');
        domSiteConfig.innerHTML = this.configTemplate(configObj);
        this.addEventListeners();
    }
}

export { PropertyManagementView };