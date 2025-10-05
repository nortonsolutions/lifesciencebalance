/**
 * GenericController to handle generic page requests
 * @author Norton 2022
 */
import { GenericView } from "./views/genericView.js";
import { EventDepot } from "./eventDepot.js";
import { defaultServices } from "./database.js";
import { initCourseSubmissionForm } from "./courseSubmissionForm.js";
import { loadAdminDashboard } from "./adminFunctions.js";
import { ProductController } from "./productController.js";
import { CartController } from "./cartController.js";

export class GenericController {
  constructor(router) {
    this.eventDepot = new EventDepot();
    this.genericView = new GenericView(router, this.eventDepot);
    this.router = router;
    
    // Initialize e-commerce controllers
    this.productController = new ProductController(router);
    this.cartController = new CartController(router);
  }

  async load(request) {
    // Parse the parameters from the request
    let params = getParamsFromRequest(request);
    
    if (!params || !params.page) {
      // Default to home if no page specified
      this.router.navigateTo("/app.html", "", false);
      return;
    }
    
    // Handle different page types
    switch (params.page) {
      case 'services':
        // Render services overview page
        this.genericView.render('pages/services', { 
          services: defaultServices,
          title: "Our Services"
        });
        break;
        
      case 'serviceDetail':
        // Render specific service detail page
        if (params.serviceId) {
          const serviceId = parseInt(params.serviceId);
          const service = defaultServices.find(s => s.id === serviceId);
          
          if (service) {
            switch (service.title) {
              case 'Students':
                this.genericView.render('pages/students', {
                  service: service,
                  title: "For Students"
                });
                break;
              case 'Teachers':
                this.genericView.render('pages/teachers', {
                  service: service,
                  title: "For Teachers"
                });
                break;
              case 'Prospective Teachers':
                this.genericView.render('pages/prospective-teachers', {
                  service: service,
                  title: "Become a Teacher"
                });
                // Initialize the course submission form
                setTimeout(() => {
                  initCourseSubmissionForm();
                }, 100);
                break;
              default:
                // Generic service detail page
                this.genericView.render('serviceDetail', service);
                break;
            }
          } else {
            // Service not found
            this.router.navigateTo("/generic.html", "?page=services", false);
          }
        }
        break;

      case 'admin':
        // Render admin page if the user has the right role
        const userRole = localStorage.getItem('userRole') || 'user';
        if (userRole === 'admin') {
          this.genericView.render('pages/admin', { title: "Administration" });
          // Initialize the admin dashboard
          setTimeout(() => {
            loadAdminDashboard();
          }, 100);
        } else {
          // Unauthorized - redirect to home
          this.router.navigateTo("/app.html", "", false);
        }
        break;
      
      case 'products':
        // Render products page (e-commerce)
        await this.productController.load(request);
        break;
      
      case 'product':
        // Render product detail page (e-commerce)
        await this.productController.load(request);
        break;
      
      case 'cart':
        // Render shopping cart page (e-commerce)
        await this.cartController.load(request);
        break;
        
      default:
        // For any other page, just render the template with that name
        this.genericView.render(`pages/${params.page}`, {
          title: params.page.charAt(0).toUpperCase() + params.page.slice(1)
        });
        break;
    }
  }
}

// Helper function to parse parameters from the request
function getParamsFromRequest(request) {
  if (!request || !request.parameters || !request.parameters.length) {
    return null;
  }
  
  const params = {};
  request.parameters.forEach(param => {
    const [key, value] = param.split('=');
    if (key && value) {
      params[key] = decodeURIComponent(value);
    }
  });
  
  return params;
}
