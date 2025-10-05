/**
 * ProductView - Render products using Blue template styling
 * @author Norton 2024
 */

export class ProductView {
    constructor(router, eventDepot) {
        this.router = router;
        this.eventDepot = eventDepot;
        this.getTemplates();
    }

    async getTemplates() {
        this.productsTemplate = await getTemplate('app/views/templates/products.hbs');
        this.productDetailTemplate = await getTemplate('app/views/templates/product-detail.hbs');
    }

    async render(context, dom = document.getElementById("main")) {
        await this.getTemplates();
        
        if (context.product) {
            // Render single product detail
            dom.innerHTML = this.productDetailTemplate(context);
            this.attachProductDetailEvents();
        } else {
            // Render products list
            dom.innerHTML = this.productsTemplate(context);
            this.attachProductListEvents();
        }
    }

    attachProductListEvents() {
        // Add to cart buttons
        const addToCartButtons = document.querySelectorAll('.add-to-cart-btn');
        addToCartButtons.forEach(button => {
            button.addEventListener('click', (e) => {
                e.preventDefault();
                const productId = button.getAttribute('data-product-id');
                this.eventDepot.trigger('cart:add', { productId: productId, quantity: 1 });
            });
        });
    }

    attachProductDetailEvents() {
        // Add to cart button in detail page
        const addToCartBtn = document.getElementById('addToCartBtn');
        if (addToCartBtn) {
            addToCartBtn.addEventListener('click', (e) => {
                e.preventDefault();
                const productId = addToCartBtn.getAttribute('data-product-id');
                const quantity = document.getElementById('quantity')?.value || 1;
                this.eventDepot.trigger('cart:add', { productId: productId, quantity: parseInt(quantity) });
            });
        }
    }
}
