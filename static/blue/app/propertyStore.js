/* Norton's rendition of the BlueSky Page (June 2019) */

class PropertyStore {

    constructor() {
        this.propertyStore = [];

        if (localStorage.getItem('propertyStore')) {
            this.propertyStore = JSON.parse(localStorage.getItem('propertyStore'));
        } else {
            localStorage.setItem('propertyStore', JSON.stringify(this.propertyStore));
        }
    }
}

PropertyStore.prototype.uniqueId = function() {
    let id = 0;
    for (let property of this.propertyStore) 
        if (property.id >= id) 
            id = Number(property.id) + 1;
    return id;
}

PropertyStore.prototype.add = function(object) {
    object.dateCreated = new Date().toISOString();
    object.dateUpdated = new Date().toISOString();
    this.propertyStore.push(object);
    this.syncStorage();
}

PropertyStore.prototype.get = function(id) {
    return this.propertyStore.find(e => e.id == id);
}

PropertyStore.prototype.getAll = function() {
    return this.propertyStore;
}

PropertyStore.prototype.update = function(object) {
    object.dateUpdated = new Date().toISOString();
    this.remove(object.id);
    this.propertyStore.push(object);
    this.syncStorage();
}

PropertyStore.prototype.remove = function(id) {
    this.propertyStore = this.propertyStore.filter(u => u.id != id );
    this.syncStorage();
}

PropertyStore.prototype.syncStorage = function() {
    localStorage.setItem('propertyStore', JSON.stringify(this.propertyStore));
}

PropertyStore.prototype.overwriteStorage = function(properties) {
    localStorage.setItem('propertyStore', JSON.stringify(properties));
}

export { PropertyStore };