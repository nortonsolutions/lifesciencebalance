/* Norton's rendition of the BlueSky Page (June 2019) */
import { PropertyStore } from "./propertyStore.js";
import { PropertyManagementView } from "./views/propertyManagementView.js";
import { PropertyMasterView } from "./views/propertyMasterView.js";
import { PropertyDetailsView } from "./views/propertyDetailsView.js";
import { EventDepot } from "./eventDepot.js";
import { NewsletterListStore } from "./newsletterListStore.js";

class PropertyManagementController {
    
    constructor(callback) {

        this.propertyObj = {};
        this.configObj = {};
        this.eventDepot = new EventDepot();
        this.addEventListeners();
        this.newPropertyInProcess = false;
        this.propertyStore = new PropertyStore();
        this.newsletterListStore = new NewsletterListStore();
        this.propertyMasterView = new PropertyMasterView(this.eventDepot);
        this.propertyDetailsView = new PropertyDetailsView(this.eventDepot);
        this.siteConfigDisplayed = false;
        this.propertyManagementView = new PropertyManagementView(() => {
            callback();
        }, this.eventDepot);

    }

    addEventListeners() {

        this.eventDepot.addListener('newProperty', e => { this.new(e) });
        this.eventDepot.addListener('deleteProperty', e => { this.delete(e) });
        this.eventDepot.addListener('saveProperty', e => { this.save(e) });
        this.eventDepot.addListener('cancel', e => { this.cancel(e) });
        this.eventDepot.addListener('deleteNewsletterListItem', e => { this.deleteEmail(e) });
    }

    deleteEmail(e) {
        this.newsletterListStore.remove(e.index);
        this.configObj.mailingList = this.newsletterListStore.getAll();
        this.propertyManagementView.reloadSiteConfiguration(this.configObj);
    }
}

PropertyManagementController.prototype.load = function(request, callback) {

    this.newsletterListStore = new NewsletterListStore();
    this.configObj.mailingList = this.newsletterListStore.getAll();

    if (request) {
        let params = getParamsFromRequest(request);
        this.propertyObj = this.propertyStore.get(params.propertyId);
    }

    this.propertyManagementView.render(() => {
        this.propertyMasterView.resetPropertyRows(this.propertyStore.getAll());
        
        if (this.propertyObj) {
            this.propertyDetailsView.populatePropertyDetails(this.propertyObj);
            this.propertyDetailsView.showInputForm();
        } else {
            this.propertyDetailsView.clearSlate();
            this.propertyDetailsView.hideInputForm();
        }

        if (callback) callback();

    }, this.propertyObj, this.configObj);
}

PropertyManagementController.prototype.new = function() {

    // throw new Error("Some error occurred.");
    let newId = this.propertyStore.uniqueId();
    this.propertyDetailsView.setCurrentPropertyId(newId);
    this.propertyDetailsView.clearSlate();
    this.propertyDetailsView.showInputForm();
    this.newPropertyInProcess = true;

}

PropertyManagementController.prototype.save = function(e) {

    e.preventDefault();

    // Check validity on the form
    if (this.propertyDetailsView.validateForm()) {

        let newProperty = {};
        newProperty.id = this.propertyDetailsView.getCurrentPropertyId();
    
        if (confirm (`Are you sure you want to save changes for property ${newProperty.id}?`)) {
    
            newProperty = this.propertyDetailsView.getFormInputValues();

            if (this.newPropertyInProcess) {
                this.propertyStore.add(newProperty);
            } else {
                this.propertyStore.update(newProperty);
            }
            
            this.newPropertyInProcess = false;
            this.load();
    
            this.propertyDetailsView.hideInputForm();
        }
    }
}

PropertyManagementController.prototype.cancel = function() {
    
    // cancel the edit and reload the same property or clear the 
    // Propertys Details view if the property was new

    if (this.newPropertyInProcess) {
        this.propertyDetailsView.clearSlate();
        this.propertyDetailsView.hideInputForm();
    } else {
        this.edit(null, this.propertyDetailsView.getCurrentPropertyId());
    }

}

PropertyManagementController.prototype.edit = function(e, id) {

    this.newPropertyInProcess = false;
    
    let propertyId;
    if (e) { propertyId = e.currentTarget.id.slice(1); } else propertyId = id;

    let currentProperty = this.propertyStore.get(propertyId);

    this.propertyDetailsView.populatePropertyDetails(currentProperty);
    this.propertyDetailsView.showInputForm();
}

PropertyManagementController.prototype.edit = function(e, id) {

    this.newPropertyInProcess = false;
    
    let propertyId;
    if (e) { propertyId = e.currentTarget.id.slice(1); } else propertyId = id;

    let currentProperty = this.propertyStore.get(propertyId);

    this.propertyDetailsView.populatePropertyDetails(currentProperty);
    this.propertyDetailsView.showInputForm();
}

PropertyManagementController.prototype.delete = function(e, id) {
    
    let propertyId;
    if (e) { propertyId = e.currentTarget.id.slice(1); } else propertyId = id;
    let propertyName = this.propertyStore.get(propertyId).name;

    if (confirm("Are you sure you want to delete " + propertyName + "?")) {
        this.propertyStore.remove(propertyId);    
        this.load();
    }
}

export { PropertyManagementController };