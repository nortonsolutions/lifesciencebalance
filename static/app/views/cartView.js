/**
 * CartView - Shopping cart display
 * @author Norton 2024
 */

export class CartView {
    constructor(router, eventDepot) {
        this.router = router;
        this.eventDepot = eventDepot;
        this.getTemplates();
    }

    async getTemplates() {
        this.template = await getTemplate('app/views/templates/cart.hbs');
    }

    async render(context, dom = document.getElementById("main")) {
        await this.getTemplates();
        dom.innerHTML = this.template(context);
        this.attachEvents();
    }

    attachEvents() {
        // Remove item buttons
        const removeButtons = document.querySelectorAll('.remove-item-btn');
        removeButtons.forEach(button => {
            button.addEventListener('click', (e) => {
                e.preventDefault();
                const productId = button.getAttribute('data-product-id');
                this.eventDepot.trigger('cart:remove', { productId: productId });
                // Reload cart view
                setTimeout(() => {
                    window.location.reload();
                }, 500);
            });
        });

        // Update quantity inputs
        const quantityInputs = document.querySelectorAll('.quantity-input');
        quantityInputs.forEach(input => {
            input.addEventListener('change', (e) => {
                const productId = input.getAttribute('data-product-id');
                const quantity = parseInt(e.target.value);
                if (quantity > 0) {
                    this.eventDepot.trigger('cart:update', { productId: productId, quantity: quantity });
                    // Reload to update totals
                    setTimeout(() => {
                        window.location.reload();
                    }, 500);
                }
            });
        });

        // Checkout button
        const checkoutBtn = document.getElementById('checkoutBtn');
        if (checkoutBtn) {
            checkoutBtn.addEventListener('click', (e) => {
                e.preventDefault();
                // Navigate to checkout page
                window.location.href = '/app.html?page=checkout';
            });
        }
    }
}
