/**
 * ProductStore - Unified store for e-commerce products
 * Evolved from PropertyStore (Blue template) and Course system
 * @author Norton 2024
 */

class ProductStore {
    constructor() {
        this.products = [];
        
        if (localStorage.getItem('productStore')) {
            this.products = JSON.parse(localStorage.getItem('productStore'));
        } else {
            localStorage.setItem('productStore', JSON.stringify(this.products));
        }
    }

    uniqueId() {
        let id = 0;
        for (let product of this.products) {
            if (product.id >= id) {
                id = Number(product.id) + 1;
            }
        }
        return id;
    }

    add(product) {
        product.dateCreated = new Date().toISOString();
        product.dateUpdated = new Date().toISOString();
        this.products.push(product);
        this.syncStorage();
    }

    get(id) {
        return this.products.find(p => p.id == id);
    }

    getAll() {
        return this.products;
    }

    getByCategory(category) {
        return this.products.filter(p => p.category === category);
    }

    getByVendor(vendorId) {
        return this.products.filter(p => p.vendor_id == vendorId);
    }

    getApproved() {
        return this.products.filter(p => p.approved === true);
    }

    update(product) {
        product.dateUpdated = new Date().toISOString();
        this.remove(product.id);
        this.products.push(product);
        this.syncStorage();
    }

    remove(id) {
        this.products = this.products.filter(p => p.id != id);
        this.syncStorage();
    }

    syncStorage() {
        localStorage.setItem('productStore', JSON.stringify(this.products));
    }

    overwriteStorage(products) {
        localStorage.setItem('productStore', JSON.stringify(products));
    }

    // API integration methods
    async fetchFromAPI() {
        try {
            const response = await fetch('/product');
            if (response.ok) {
                const products = await response.json();
                this.products = products;
                this.syncStorage();
                return products;
            }
        } catch (error) {
            console.error('Error fetching products:', error);
        }
        return this.products;
    }

    async fetchApprovedFromAPI() {
        try {
            const response = await fetch('/product/approved');
            if (response.ok) {
                const products = await response.json();
                this.products = products;
                this.syncStorage();
                return products;
            }
        } catch (error) {
            console.error('Error fetching approved products:', error);
        }
        return this.getApproved();
    }
}

export { ProductStore };
