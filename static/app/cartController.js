/**
 * CartController - Shopping cart management
 * @author Norton 2024
 */
import { CartView } from "./views/cartView.js";
import { CartStore } from "./cartStore.js";
import { ProductStore } from "./productStore.js";
import { EventDepot } from "./eventDepot.js";

export class CartController {
    constructor(router) {
        this.eventDepot = new EventDepot();
        this.cartStore = new CartStore();
        this.productStore = new ProductStore();
        this.cartView = new CartView(router, this.eventDepot);
        
        // Listen for cart events
        this.eventDepot.on('cart:add', this.handleAddToCart.bind(this));
        this.eventDepot.on('cart:remove', this.handleRemoveFromCart.bind(this));
        this.eventDepot.on('cart:update', this.handleUpdateCart.bind(this));
    }

    async handleAddToCart(data) {
        const { productId, quantity } = data;
        const customerId = localStorage.getItem('userId');
        
        if (customerId) {
            await this.cartStore.addToAPICart(customerId, productId, quantity);
        } else {
            this.cartStore.addItem(productId, quantity);
        }
        
        // Show notification
        this.showNotification('Product added to cart!');
        
        // Update cart badge if exists
        this.updateCartBadge();
    }

    handleRemoveFromCart(data) {
        const { productId } = data;
        this.cartStore.removeItem(productId);
        this.showNotification('Product removed from cart');
        this.updateCartBadge();
    }

    handleUpdateCart(data) {
        const { productId, quantity } = data;
        this.cartStore.updateItem(productId, quantity);
        this.updateCartBadge();
    }

    async load(request) {
        const customerId = localStorage.getItem('userId');
        
        if (customerId) {
            await this.cartStore.fetchFromAPI(customerId);
        }
        
        // Fetch product details for cart items
        await this.productStore.fetchApprovedFromAPI();
        const cartItems = this.cartStore.getItems();
        
        const itemsWithDetails = cartItems.map(item => {
            const product = this.productStore.get(item.product_id);
            return {
                ...item,
                product: product
            };
        });
        
        const context = {
            items: itemsWithDetails,
            itemCount: this.cartStore.getItemCount(),
            total: this.calculateTotal(itemsWithDetails)
        };
        
        this.cartView.render(context);
    }

    calculateTotal(items) {
        return items.reduce((total, item) => {
            return total + (item.product?.price || 0) * item.quantity;
        }, 0).toFixed(2);
    }

    updateCartBadge() {
        const badge = document.getElementById('cartBadge');
        if (badge) {
            const count = this.cartStore.getItemCount();
            badge.textContent = count;
            badge.style.display = count > 0 ? 'inline-block' : 'none';
        }
    }

    showNotification(message) {
        // Simple notification - can be enhanced with a toast library
        const notification = document.createElement('div');
        notification.className = 'alert alert-success alert-dismissible fade show position-fixed';
        notification.style.cssText = 'top: 20px; right: 20px; z-index: 9999; min-width: 250px;';
        notification.innerHTML = `
            ${message}
            <button type="button" class="close" data-dismiss="alert">
                <span>&times;</span>
            </button>
        `;
        document.body.appendChild(notification);
        
        setTimeout(() => {
            notification.remove();
        }, 3000);
    }
}
