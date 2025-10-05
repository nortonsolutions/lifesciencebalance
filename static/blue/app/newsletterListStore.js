/* Norton's rendition of the BlueSky Page (June 2019) */

class NewsletterListStore {

    constructor() {
        this.newsletterListStore = [];

        if (localStorage.getItem('newsletterListStore')) {
            this.newsletterListStore = JSON.parse(localStorage.getItem('newsletterListStore'));
        } else {
            localStorage.setItem('newsletterListStore', JSON.stringify(this.newsletterListStore));
        }
    }
}

NewsletterListStore.prototype.uniqueId = function() {
    let id = 0;
    for (let property of this.newsletterListStore) 
        if (property.id >= id) 
            id = Number(property.id) + 1;
    return id;
}

NewsletterListStore.prototype.add = function(object) {
    object.dateCreated = new Date().toISOString();
    object.dateUpdated = new Date().toISOString();
    this.newsletterListStore.push(object);
    this.syncStorage();
    return true;
}

NewsletterListStore.prototype.get = function(id) {
    return this.newsletterListStore.find(e => e.id == id);
}

NewsletterListStore.prototype.getAll = function() {
    return this.newsletterListStore;
}

NewsletterListStore.prototype.update = function(object) {
    object.dateUpdated = new Date().toISOString();
    this.remove(object.id);
    this.newsletterListStore.push(object);
    this.syncStorage();
}

NewsletterListStore.prototype.remove = function(id) {
    this.newsletterListStore = this.newsletterListStore.filter(u => u.id != id );
    this.syncStorage();
}

NewsletterListStore.prototype.syncStorage = function() {
    localStorage.setItem('newsletterListStore', JSON.stringify(this.newsletterListStore));
}

export { NewsletterListStore };