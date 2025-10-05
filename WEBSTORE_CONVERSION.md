# Webstore Conversion Documentation

## Overview
This document details the conversion of the Life Science Balance application from a course/quiz delivery system to an e-commerce webstore platform, integrating the BlueSky frontend template styling.

## Architecture

### Backend (Go + Google Cloud Datastore)

#### New Models
1. **Product** (`models/ProductModel.go`)
   - Evolved from Course model
   - Fields: name, description, price, image_url, category, vendor_id, components, approved, in_stock, sku, tags
   - Represents items for sale in the webstore

2. **Customer** (`models/CustomerModel.go`)
   - Evolved from User with Student role context
   - Fields: user_id (reference), shipping_address, billing_address, phone, preferred_payment
   - Extends User model for customer-specific information

3. **Vendor** (`models/VendorModel.go`)
   - Evolved from User with Teacher role context
   - Fields: user_id (reference), business_name, business_email, phone, address, tax_id, payment_info, approved
   - Extends User model for vendor-specific information

4. **Order** (`models/OrderModel.go`)
   - New model for e-commerce transactions
   - Fields: customer_id, items (OrderItem array), total_amount, status, shipping/billing addresses, payment info
   - Manages the complete order lifecycle

5. **Cart** (`models/CartModel.go`)
   - New model for shopping cart functionality
   - Fields: customer_id, items (CartItem array), updated_at
   - Temporary storage for products before checkout

6. **ProductComponent** (`models/ProductComponentModel.go`)
   - Evolved from Module model
   - Can represent product page components like image galleries, reviews, specifications, or transaction handlers

#### New Controllers
1. **ProductController** - CRUD operations for products
2. **CustomerController** - Customer management
3. **VendorController** - Vendor management and approval
4. **CartController** - Shopping cart operations
5. **OrderController** - Order creation, checkout, and management

#### New API Endpoints

**Product Management:**
- `POST /product` - Create new product
- `GET /product` - Get all products
- `GET /product/{id}` - Get product by ID
- `PUT /product/{id}` - Update product
- `DELETE /product/{id}` - Delete product
- `GET /product/approved` - Get approved products
- `GET /product/category/{category}` - Filter by category
- `GET /product/vendor/{vendorId}` - Get vendor's products
- `PUT /product/{id}/approve` - Approve product (admin)

**Customer Management:**
- `POST /customer` - Create customer
- `GET /customer` - Get all customers
- `GET /customer/{id}` - Get customer by ID
- `GET /customer/user/{userId}` - Get customer by user ID
- `PUT /customer/{id}` - Update customer
- `DELETE /customer/{id}` - Delete customer

**Vendor Management:**
- `POST /vendor` - Create vendor
- `GET /vendor` - Get all vendors
- `GET /vendor/{id}` - Get vendor by ID
- `GET /vendor/user/{userId}` - Get vendor by user ID
- `PUT /vendor/{id}` - Update vendor
- `DELETE /vendor/{id}` - Delete vendor
- `GET /vendor/approved` - Get approved vendors
- `PUT /vendor/{id}/approve` - Approve vendor (admin)

**Shopping Cart:**
- `GET /cart/customer/{customerId}` - Get customer's cart
- `POST /cart/customer/{customerId}` - Add item to cart
- `PUT /cart/customer/{customerId}` - Update cart item
- `DELETE /cart/customer/{customerId}/product/{productId}` - Remove item
- `DELETE /cart/customer/{customerId}/clear` - Clear cart

**Orders:**
- `POST /order` - Create order
- `GET /order` - Get all orders
- `GET /order/{id}` - Get order by ID
- `PUT /order/{id}/status` - Update order status
- `GET /order/customer/{customerId}` - Get customer's orders
- `POST /order/customer/{customerId}/checkout` - Create order from cart

### Frontend (Vanilla JavaScript + Handlebars)

#### New Components

1. **ProductStore** (`static/app/productStore.js`)
   - Unified store combining PropertyStore (Blue) and Course logic
   - Manages product data with localStorage caching
   - API integration methods for fetching products

2. **CartStore** (`static/app/cartStore.js`)
   - Shopping cart state management
   - Local and API synchronization
   - Cart operations: add, update, remove, clear

3. **ProductController** (`static/app/productController.js`)
   - Handles product browsing and display
   - Routes: `/generic.html?page=products`, `/generic.html?page=product&detail={id}`

4. **CartController** (`static/app/cartController.js`)
   - Shopping cart functionality
   - Event-driven architecture for cart updates
   - Route: `/generic.html?page=cart`

5. **ProductView** (`static/app/views/productView.js`)
   - Renders product listings and detail pages
   - Uses Blue template styling

6. **CartView** (`static/app/views/cartView.js`)
   - Renders shopping cart interface
   - Interactive quantity updates and item removal

#### Templates

1. **products.hbs** - Product grid/list view
2. **product-detail.hbs** - Individual product page with add-to-cart
3. **cart.hbs** - Shopping cart with order summary

#### Navigation Updates

- Added "Products" menu item in navbar
- Added shopping cart icon with badge showing item count
- Cart accessible from navigation bar

## Conversion Mapping

### Terminology Changes
| Old (Course System) | New (Webstore) |
|---------------------|----------------|
| Course              | Product        |
| Student             | Customer       |
| Teacher             | Vendor         |
| Module              | ProductComponent |
| Enrollment          | Order          |

### Role Evolution
- **Student role** → **Customer role**: Users who purchase products
- **Teacher role** → **Vendor role**: Users who sell products
- **Admin role** remains for platform management

### Data Migration Strategy

To migrate existing data:

1. **Courses → Products**
   ```go
   // Map Course fields to Product fields
   product.Name = course.Name
   product.Description = course.Description
   product.Price = 0.00 // Set default price
   product.Category = course.Department
   product.VendorID = course.OwnerID
   product.Approved = course.Approved
   ```

2. **Users with Student role → Customers**
   ```go
   // Create Customer record for each Student
   customer.UserID = user.KeyID
   // Collect additional customer information
   ```

3. **Users with Teacher role → Vendors**
   ```go
   // Create Vendor record for each Teacher
   vendor.UserID = user.KeyID
   vendor.BusinessName = user.Firstname + " " + user.Lastname
   // Collect additional vendor information
   ```

## Usage Guide

### For Customers
1. Browse products at `/generic.html?page=products`
2. View product details by clicking on a product
3. Add products to cart
4. View cart at `/generic.html?page=cart`
5. Proceed to checkout (to be implemented)

### For Vendors
1. Register as a vendor
2. Wait for admin approval
3. Create and manage products
4. View orders for their products

### For Admins
1. Approve/unapprove vendors
2. Approve/unapprove products
3. Manage orders and system settings

## Backward Compatibility

The original course system remains functional. Both systems coexist:
- Original course routes: `/course`, `/module`, etc.
- New e-commerce routes: `/product`, `/cart`, `/order`, etc.

This allows for gradual migration and testing without disrupting existing functionality.

## Styling Integration

The Blue template's styling has been integrated into the new product and cart views:
- Card-based layouts for product display
- Clean, modern design aesthetic
- Responsive grid layouts
- Bootstrap 5 components
- Font Awesome icons

## Future Enhancements

1. **Payment Integration**: Add Stripe/PayPal for actual transactions
2. **Product Reviews**: Allow customers to review products
3. **Wishlist**: Save products for later
4. **Product Search**: Enhanced search and filtering
5. **Order Tracking**: Real-time order status updates
6. **Vendor Dashboard**: Analytics and sales reports
7. **Email Notifications**: Order confirmations and updates
8. **Inventory Management**: Track stock levels
9. **Shipping Integration**: Calculate shipping costs
10. **Product Variants**: Size, color options, etc.

## Testing

### Manual Testing Steps
1. Build the Go application: `go build`
2. Start the server with appropriate environment variables
3. Navigate to products page
4. Test adding items to cart
5. Verify cart operations (update quantity, remove items)
6. Test product approval workflow (admin)
7. Test vendor approval workflow (admin)

### API Testing
Use tools like Postman or curl to test the new endpoints:

```bash
# Create a product
curl -X POST http://localhost:8000/product \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Sample Product",
    "description": "A great product",
    "price": 29.99,
    "category": "Electronics",
    "in_stock": true
  }'

# Get all products
curl http://localhost:8000/product

# Add to cart
curl -X POST http://localhost:8000/cart/customer/1 \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": 1,
    "quantity": 2
  }'
```

## Conclusion

This conversion successfully transforms the Life Science Balance platform from an educational course system into a fully-functional e-commerce webstore while:
- Maintaining backward compatibility
- Leveraging existing authentication and user management
- Integrating modern, clean UI from the Blue template
- Providing a solid foundation for future enhancements

The modular architecture allows for easy extension and customization to meet specific business needs.
