/**
 * ServicesController
 */

import { ServicesView } from "./views/servicesView.js";
import { EventDepot } from "./eventDepot.js";
import { defaultServices } from "./database.js";
import { initCourseSubmissionForm } from "./courseSubmissionForm.js";

export class ServicesController {
  constructor(router) {
    this.router = router;
    this.servicesView = new ServicesView(router);
    this.eventDepot = EventDepot;
  }

  async load(page) {
    // Handle the services overview page
    if (!page || page === "services") {
      this.servicesView.render({ 
        page: "services", 
        title: "Services Overview",
        services: defaultServices
      });
    }
    
    // Handle the students page
    if (page === "students") {
      this.servicesView.render({ 
        page: "students",
        title: "For Students"
      });
    }
    
    // Handle the teachers page
    if (page === "teachers") {
      this.servicesView.render({ 
        page: "teachers",
        title: "For Teachers"
      });
    }
    
    // Handle the prospective teachers page
    if (page === "prospective-teachers") {
      this.servicesView.render({ 
        page: "prospective-teachers",
        title: "Become a Teacher"
      });
      
      // Initialize the course submission form after rendering
      setTimeout(() => {
        initCourseSubmissionForm();
      }, 100);
    }
  }
}
