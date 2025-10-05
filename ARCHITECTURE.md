# Webstore Architecture Overview

## System Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────────┐
│                           USER INTERFACE                             │
│                     (Browser - HTML/CSS/JS)                          │
└────────────────────────┬────────────────────────────────────────────┘
                         │
                         │ HTTP Requests
                         ▼
┌─────────────────────────────────────────────────────────────────────┐
│                      FRONTEND LAYER                                  │
│                   (Vanilla JavaScript + Handlebars)                  │
│                                                                       │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐              │
│  │   Product    │  │     Cart     │  │   Customer   │              │
│  │  Controller  │  │  Controller  │  │  Controller  │              │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘              │
│         │                  │                  │                       │
│         ▼                  ▼                  ▼                       │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐              │
│  │   Product    │  │     Cart     │  │   Generic    │              │
│  │     View     │  │     View     │  │     View     │              │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘              │
│         │                  │                  │                       │
│         ▼                  ▼                  ▼                       │
│  ┌──────────────┐  ┌──────────────┐                                 │
│  │   Product    │  │     Cart     │                                 │
│  │    Store     │  │    Store     │  (LocalStorage + API Sync)      │
│  └──────────────┘  └──────────────┘                                 │
│                                                                       │
│  Templates: products.hbs, product-detail.hbs, cart.hbs              │
└────────────────────────┬────────────────────────────────────────────┘
                         │
                         │ REST API Calls
                         ▼
┌─────────────────────────────────────────────────────────────────────┐
│                       BACKEND LAYER                                  │
│                   (Go + Gorilla Mux Router)                          │
│                                                                       │
│  ┌──────────────────────────────────────────────────────────────┐   │
│  │                     ROUTES LAYER                              │   │
│  │  /product, /cart, /order, /customer, /vendor                 │   │
│  └────────────────────┬─────────────────────────────────────────┘   │
│                       │                                               │
│                       ▼                                               │
│  ┌──────────────────────────────────────────────────────────────┐   │
│  │                  CONTROLLERS LAYER                            │   │
│  │                                                               │   │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐       │   │
│  │  │ Product  │ │   Cart   │ │  Order   │ │ Customer │       │   │
│  │  │Controller│ │Controller│ │Controller│ │Controller│       │   │
│  │  └────┬─────┘ └────┬─────┘ └────┬─────┘ └────┬─────┘       │   │
│  │       │            │            │            │               │   │
│  └───────┼────────────┼────────────┼────────────┼───────────────┘   │
│          │            │            │            │                    │
│          ▼            ▼            ▼            ▼                    │
│  ┌──────────────────────────────────────────────────────────────┐   │
│  │                 REPOSITORIES LAYER                            │   │
│  │                                                               │   │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐       │   │
│  │  │ Product  │ │   Cart   │ │  Order   │ │ Customer │       │   │
│  │  │   Repo   │ │   Repo   │ │   Repo   │ │   Repo   │       │   │
│  │  └────┬─────┘ └────┬─────┘ └────┬─────┘ └────┬─────┘       │   │
│  │       │            │            │            │               │   │
│  └───────┼────────────┼────────────┼────────────┼───────────────┘   │
│          │            │            │            │                    │
│          ▼            ▼            ▼            ▼                    │
│  ┌──────────────────────────────────────────────────────────────┐   │
│  │                     MODELS LAYER                              │   │
│  │                                                               │   │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐       │   │
│  │  │ Product  │ │   Cart   │ │  Order   │ │ Customer │       │   │
│  │  │  Model   │ │  Model   │ │  Model   │ │  Model   │       │   │
│  │  └──────────┘ └──────────┘ └──────────┘ └──────────┘       │   │
│  │                                                               │   │
│  └───────────────────────┬───────────────────────────────────────┘   │
│                          │                                            │
└──────────────────────────┼────────────────────────────────────────────┘
                           │
                           │ Datastore API
                           ▼
┌─────────────────────────────────────────────────────────────────────┐
│                    DATA PERSISTENCE LAYER                            │
│                  Google Cloud Datastore                              │
│                                                                       │
│  Entities: Product, Cart, Order, Customer, Vendor, User             │
└─────────────────────────────────────────────────────────────────────┘
```

## Request Flow Examples

### 1. Browsing Products

```
User → GET /generic.html?page=products
  ↓
Frontend: ProductController.load()
  ↓
ProductStore.fetchApprovedFromAPI()
  ↓
Backend: GET /product/approved
  ↓
ProductController.GetApprovedProducts()
  ↓
ProductRepository.GetApprovedProducts()
  ↓
Google Cloud Datastore Query
  ↓
Products returned as JSON
  ↓
ProductStore caches to localStorage
  ↓
ProductView renders products.hbs
  ↓
User sees product grid
```

### 2. Adding to Cart

```
User clicks "Add to Cart"
  ↓
ProductView fires event
  ↓
CartController.handleAddToCart()
  ↓
CartStore.addToAPICart()
  ↓
Backend: POST /cart/customer/{id}
  ↓
CartController.AddToCart()
  ↓
CartRepository.UpdateCart()
  ↓
Google Cloud Datastore Save
  ↓
Updated cart returned
  ↓
CartStore updates localStorage
  ↓
Cart badge updates with new count
  ↓
User sees notification
```

### 3. Checkout Process

```
User clicks "Proceed to Checkout"
  ↓
Navigate to /generic.html?page=cart
  ↓
CartController.load()
  ↓
Fetch cart items and product details
  ↓
CartView renders cart.hbs
  ↓
User clicks "Checkout"
  ↓
Backend: POST /order/customer/{id}/checkout
  ↓
OrderController.CreateOrderFromCart()
  ↓
- Validate cart items
  - Calculate totals
  - Create Order entity
  - Clear cart
  ↓
Order confirmation returned
  ↓
User sees order summary
```

## Data Flow Patterns

### Frontend State Management

```
┌──────────────┐
│ User Action  │
└──────┬───────┘
       │
       ▼
┌──────────────┐      ┌──────────────┐
│  Controller  │─────▶│  Event Bus   │
└──────┬───────┘      └──────────────┘
       │
       ▼
┌──────────────┐
│    Store     │◄─────┐
└──────┬───────┘      │
       │              │
       ├──────────────┤
       │              │
       ▼              │
┌──────────────┐      │
│ LocalStorage │      │
└──────────────┘      │
       │              │
       ▼              │
┌──────────────┐      │
│   API Call   │──────┘
└──────┬───────┘
       │
       ▼
┌──────────────┐
│     View     │
└──────────────┘
```

### Backend Data Access Pattern

```
HTTP Request
     │
     ▼
┌──────────────┐
│  Controller  │ (Request handling, validation)
└──────┬───────┘
       │
       ▼
┌──────────────┐
│  Repository  │ (Data access logic)
└──────┬───────┘
       │
       ▼
┌──────────────┐
│    Model     │ (Data structure)
└──────┬───────┘
       │
       ▼
┌──────────────┐
│  Datastore   │ (Persistence)
└──────────────┘
```

## Component Relationships

### E-Commerce Domain Model

```
┌──────────┐       ┌──────────┐       ┌──────────┐
│   User   │──────▶│ Customer │──────▶│  Order   │
└────┬─────┘       └────┬─────┘       └────┬─────┘
     │                  │                   │
     │                  │                   │
     │                  ▼                   │
     │            ┌──────────┐              │
     │            │   Cart   │              │
     │            └────┬─────┘              │
     │                 │                    │
     ▼                 ▼                    ▼
┌──────────┐       ┌──────────┐       ┌──────────┐
│  Vendor  │──────▶│ Product  │◄──────│OrderItem │
└──────────┘       └────┬─────┘       └──────────┘
                        │
                        ▼
                   ┌──────────┐
                   │Component │
                   └──────────┘

Relationships:
- User → Customer (one-to-one)
- User → Vendor (one-to-one)
- Customer → Cart (one-to-one)
- Customer → Order (one-to-many)
- Vendor → Product (one-to-many)
- Cart → CartItem (one-to-many)
- Order → OrderItem (one-to-many)
- Product → ProductComponent (one-to-many)
```

## Technology Stack

### Backend Stack
```
┌─────────────────────────────┐
│      Go 1.18+               │
├─────────────────────────────┤
│  - Gorilla Mux (routing)    │
│  - godotenv (config)        │
│  - Google Cloud SDK         │
└─────────────────────────────┘
          │
          ▼
┌─────────────────────────────┐
│  Google Cloud Datastore     │
│  (NoSQL Document DB)        │
└─────────────────────────────┘
```

### Frontend Stack
```
┌─────────────────────────────┐
│   Vanilla JavaScript        │
├─────────────────────────────┤
│  - Handlebars (templating)  │
│  - Bootstrap 5 (CSS)        │
│  - Font Awesome (icons)     │
│  - LocalStorage (state)     │
└─────────────────────────────┘
```

## Security Architecture

```
┌──────────────┐
│   Browser    │
└──────┬───────┘
       │ HTTPS (when enabled)
       ▼
┌──────────────┐
│   Gorilla    │
│   Router     │◄─────┐
└──────┬───────┘      │
       │              │
       ▼              │
┌──────────────┐      │
│ Permissions  │      │
│  Middleware  │──────┘ (Route protection)
└──────┬───────┘
       │
       ▼
┌──────────────┐
│   Session    │
│ Validation   │
└──────┬───────┘
       │
       ▼
┌──────────────┐
│  Controller  │
│  (Business   │
│    Logic)    │
└──────────────┘
```

## Deployment Architecture

```
┌─────────────────────────────────────────┐
│         Production Environment          │
│                                         │
│  ┌───────────────────────────────────┐ │
│  │      Load Balancer (Optional)     │ │
│  └─────────────┬─────────────────────┘ │
│                │                        │
│    ┌───────────┴───────────┐           │
│    │                       │           │
│    ▼                       ▼           │
│  ┌─────┐                 ┌─────┐      │
│  │ App │                 │ App │      │
│  │ Go  │                 │ Go  │      │
│  └──┬──┘                 └──┬──┘      │
│     │                       │          │
│     └───────────┬───────────┘          │
│                 │                      │
│                 ▼                      │
│  ┌──────────────────────────────────┐ │
│  │   Google Cloud Datastore         │ │
│  └──────────────────────────────────┘ │
│                                         │
│  ┌──────────────────────────────────┐ │
│  │   Static Files (CDN Optional)    │ │
│  └──────────────────────────────────┘ │
└─────────────────────────────────────────┘
```

## Scalability Considerations

### Horizontal Scaling
- Stateless Go application servers
- Session data in centralized store
- Load balancer for traffic distribution

### Vertical Scaling
- Google Cloud Datastore auto-scales
- No database server to manage
- Pay-per-use pricing model

### Caching Strategy
```
Browser Cache (LocalStorage)
       ↓
API Response Cache (Future)
       ↓
Datastore (Source of Truth)
```

## Monitoring & Observability

```
Application Logs
       ↓
┌──────────────┐
│   Logging    │ → Cloud Logging
└──────────────┘
       │
       ▼
┌──────────────┐
│   Metrics    │ → Cloud Monitoring
└──────────────┘
       │
       ▼
┌──────────────┐
│   Alerts     │ → Notification Channels
└──────────────┘
```

## Summary

This architecture provides:
- ✅ **Separation of Concerns**: Clear layer boundaries
- ✅ **Scalability**: Stateless servers, managed database
- ✅ **Maintainability**: Repository pattern, MVC structure
- ✅ **Security**: Session validation, role-based access
- ✅ **Performance**: LocalStorage caching, efficient queries
- ✅ **Flexibility**: Easy to extend with new features

The system is designed to be:
- **Modular**: Add features without breaking existing code
- **Testable**: Interface-based design allows mocking
- **Observable**: Logs and metrics for debugging
- **Resilient**: Graceful error handling and fallbacks
