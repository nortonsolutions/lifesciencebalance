/**
 * ProductController - Main controller for product browsing and management
 * Evolved from PropertyController (Blue) and CourseController
 * @author Norton 2024
 */
import { ProductView } from "./views/productView.js";
import { ProductStore } from "./productStore.js";
import { EventDepot } from "./eventDepot.js";

export class ProductController {
    constructor(router) {
        this.eventDepot = new EventDepot();
        this.productStore = new ProductStore();
        this.productView = new ProductView(router, this.eventDepot);
    }

    async load(request) {
        let params = getParamsFromRequest(request);
        
        // Fetch products from API
        await this.productStore.fetchApprovedFromAPI();
        
        let context = {
            products: this.productStore.getApproved()
        };

        // If specific product ID is requested
        if (params?.productId) {
            const product = this.productStore.get(params.productId);
            context.product = product;
        }

        // If category filter is requested
        if (params?.category) {
            context.products = this.productStore.getByCategory(params.category);
            context.category = params.category;
        }

        this.productView.render(context);
    }
}
