/**
 * CartStore - Shopping cart management
 * @author Norton 2024
 */

class CartStore {
    constructor() {
        this.cart = {
            items: [],
            customerId: null,
            updatedAt: null
        };
        
        if (localStorage.getItem('cart')) {
            this.cart = JSON.parse(localStorage.getItem('cart'));
        } else {
            localStorage.setItem('cart', JSON.stringify(this.cart));
        }
    }

    setCustomerId(customerId) {
        this.cart.customerId = customerId;
        this.syncStorage();
    }

    addItem(productId, quantity = 1) {
        const existingItem = this.cart.items.find(item => item.product_id == productId);
        
        if (existingItem) {
            existingItem.quantity += quantity;
        } else {
            this.cart.items.push({
                product_id: productId,
                quantity: quantity
            });
        }
        
        this.cart.updatedAt = new Date().toISOString();
        this.syncStorage();
    }

    updateItem(productId, quantity) {
        const item = this.cart.items.find(item => item.product_id == productId);
        
        if (item) {
            if (quantity <= 0) {
                this.removeItem(productId);
            } else {
                item.quantity = quantity;
                this.cart.updatedAt = new Date().toISOString();
                this.syncStorage();
            }
        }
    }

    removeItem(productId) {
        this.cart.items = this.cart.items.filter(item => item.product_id != productId);
        this.cart.updatedAt = new Date().toISOString();
        this.syncStorage();
    }

    getItems() {
        return this.cart.items;
    }

    getItemCount() {
        return this.cart.items.reduce((total, item) => total + item.quantity, 0);
    }

    clear() {
        this.cart.items = [];
        this.cart.updatedAt = new Date().toISOString();
        this.syncStorage();
    }

    syncStorage() {
        localStorage.setItem('cart', JSON.stringify(this.cart));
    }

    // API integration methods
    async fetchFromAPI(customerId) {
        try {
            const response = await fetch(`/cart/customer/${customerId}`);
            if (response.ok) {
                const cart = await response.json();
                this.cart = {
                    items: cart.items || [],
                    customerId: cart.customer_id,
                    updatedAt: cart.updated_at
                };
                this.syncStorage();
                return this.cart;
            }
        } catch (error) {
            console.error('Error fetching cart:', error);
        }
        return this.cart;
    }

    async addToAPICart(customerId, productId, quantity = 1) {
        try {
            const response = await fetch(`/cart/customer/${customerId}`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    product_id: productId,
                    quantity: quantity
                })
            });
            
            if (response.ok) {
                const cart = await response.json();
                this.cart = {
                    items: cart.items || [],
                    customerId: cart.customer_id,
                    updatedAt: cart.updated_at
                };
                this.syncStorage();
                return this.cart;
            }
        } catch (error) {
            console.error('Error adding to cart:', error);
        }
        
        // Fallback to local storage
        this.addItem(productId, quantity);
        return this.cart;
    }
}

export { CartStore };
