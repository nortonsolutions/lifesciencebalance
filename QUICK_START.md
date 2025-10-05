# Quick Start Guide - Webstore

## For Developers

### 1. Setup Environment
```bash
# Clone repository
git clone https://github.com/nortonsolutions/lifesciencebalance.git
cd lifesciencebalance

# Copy environment template
cp .env_SAMPLE .env

# Edit .env with your credentials
# Required variables:
# - DATASTORE_PROJECT_ID
# - GOOGLE_APPLICATION_CREDENTIALS
# - PORT
# - SSL_ENABLED
```

### 2. Build & Run
```bash
# Install dependencies
go mod download

# Build application
go build -o app

# Run server
./app

# Server starts on port specified in .env (default: 8000)
```

### 3. Access Webstore
```
Homepage: http://localhost:8000/
Products: http://localhost:8000/generic.html?page=products
Cart: http://localhost:8000/generic.html?page=cart
```

---

## For Users

### Browse Products
1. Navigate to Products page (click "Products" in navigation)
2. Browse available products
3. Click on any product for details

### Add to Cart
1. On product detail page, select quantity
2. Click "Add to Cart" button
3. See cart badge update in navigation

### Manage Cart
1. Click shopping cart icon in navigation
2. Update quantities as needed
3. Remove items if desired
4. Click "Proceed to Checkout" when ready

---

## For Vendors

### Register as Vendor
1. Contact administrator to create vendor account
2. Wait for admin approval
3. Receive notification when approved

### Add Products
```bash
# Via API
curl -X POST http://localhost:8000/product \
  -H "Content-Type: application/json" \
  -d '{
    "name": "My Product",
    "description": "Product description",
    "price": 29.99,
    "category": "Category Name",
    "vendor_id": YOUR_VENDOR_ID,
    "in_stock": true,
    "sku": "PROD-001"
  }'
```

### Manage Products
- View your products: `GET /product/vendor/{vendorId}`
- Update product: `PUT /product/{id}`
- Delete product: `DELETE /product/{id}`

---

## For Administrators

### Approve Vendors
```bash
# Approve a vendor
curl -X PUT http://localhost:8000/vendor/{id}/approve \
  -H "Cookie: session_token=YOUR_SESSION"
```

### Approve Products
```bash
# Approve a product
curl -X PUT http://localhost:8000/product/{id}/approve \
  -H "Cookie: session_token=YOUR_SESSION"
```

### View Statistics
- Access admin panel: `/generic.html?page=admin`
- View system stats: `GET /admin/stats`

---

## API Quick Reference

### Products
```bash
# List all approved products
GET /product/approved

# Get specific product
GET /product/{id}

# Create product
POST /product

# Update product
PUT /product/{id}

# Filter by category
GET /product/category/{category}
```

### Shopping Cart
```bash
# Get cart
GET /cart/customer/{customerId}

# Add to cart
POST /cart/customer/{customerId}
{
  "product_id": 1,
  "quantity": 2
}

# Update cart item
PUT /cart/customer/{customerId}
{
  "product_id": 1,
  "quantity": 3
}

# Remove from cart
DELETE /cart/customer/{customerId}/product/{productId}

# Clear cart
DELETE /cart/customer/{customerId}/clear
```

### Orders
```bash
# Create order from cart
POST /order/customer/{customerId}/checkout

# Get customer orders
GET /order/customer/{customerId}

# Get order details
GET /order/{id}

# Update order status
PUT /order/{id}/status
{
  "status": "shipped"
}
```

---

## Common Tasks

### Seed Sample Data
```bash
# When database access is available
go run scripts/seed_sample_products.go
```

### Migrate Courses to Products
```bash
# Convert existing courses
go run scripts/migrate_courses_to_products.go
```

### Check System Health
```bash
# Verify server is running
curl http://localhost:8000/product/approved

# Should return JSON array of products
```

### Debug Issues
```bash
# Check server logs
tail -f /path/to/logs

# Test API endpoints
curl -v http://localhost:8000/product/approved

# Verify environment variables
env | grep DATASTORE
```

---

## File Locations

### Backend
```
models/         - Data models
repositories/   - Database access
controllers/    - Business logic
routes/         - API endpoints
```

### Frontend
```
static/app/                  - Application code
static/app/views/            - View components
static/app/views/templates/  - Handlebars templates
```

### Documentation
```
WEBSTORE_CONVERSION.md      - Complete technical docs
IMPLEMENTATION_SUMMARY.md   - Executive summary
ARCHITECTURE.md             - System architecture
QUICK_START.md             - This file
```

---

## Troubleshooting

### Build Fails
```bash
# Update dependencies
go mod tidy
go mod download

# Clean and rebuild
go clean
go build
```

### Cannot Connect to Database
```bash
# Check environment variables
echo $DATASTORE_PROJECT_ID
echo $GOOGLE_APPLICATION_CREDENTIALS

# Verify credentials file exists
ls -l $GOOGLE_APPLICATION_CREDENTIALS

# Test datastore connection
gcloud datastore operations list --project=$DATASTORE_PROJECT_ID
```

### Frontend Not Loading
```bash
# Check static file directory
ls -l static/

# Verify STATIC_DIR in .env
grep STATIC_DIR .env

# Check server logs for errors
```

### Cart Not Updating
```bash
# Clear browser cache and localStorage
# In browser console:
localStorage.clear()

# Reload page
location.reload()
```

---

## Environment Variables

### Required
```bash
DATASTORE_PROJECT_ID=your-project-id
GOOGLE_APPLICATION_CREDENTIALS=/path/to/credentials.json
PORT=8000
```

### Optional
```bash
SSL_ENABLED=false
CERT_FILE=/path/to/cert.pem
KEY_FILE=/path/to/key.pem
STATIC_DIR=./static
```

---

## Development Workflow

### Making Changes

1. **Backend Changes**
   ```bash
   # Edit Go files
   # Rebuild
   go build
   # Restart server
   ./app
   ```

2. **Frontend Changes**
   ```bash
   # Edit JS/HTML files
   # No rebuild needed
   # Just refresh browser
   ```

3. **Testing**
   ```bash
   # Manual testing
   # Follow test checklist in IMPLEMENTATION_SUMMARY.md
   ```

4. **Commit**
   ```bash
   git add .
   git commit -m "Description of changes"
   git push
   ```

---

## Support

### Documentation
- Technical Details: `WEBSTORE_CONVERSION.md`
- Architecture: `ARCHITECTURE.md`
- Full Summary: `IMPLEMENTATION_SUMMARY.md`

### Resources
- Go Documentation: https://golang.org/doc/
- Gorilla Mux: https://github.com/gorilla/mux
- Google Cloud Datastore: https://cloud.google.com/datastore/docs
- Handlebars: https://handlebarsjs.com/
- Bootstrap 5: https://getbootstrap.com/docs/5.0/

---

## Next Steps

1. ✅ Set up environment
2. ✅ Build and run application
3. ✅ Explore webstore features
4. ✅ Read documentation
5. ⏭️ Add your products
6. ⏭️ Test complete workflow
7. ⏭️ Deploy to production
8. ⏭️ Implement Phase 2 features

---

## Quick Tips

- **Start Simple**: Browse products first, then explore cart functionality
- **Use API Docs**: Reference `WEBSTORE_CONVERSION.md` for all endpoints
- **Check Logs**: Server logs show request/response details
- **Browser Console**: Check for JavaScript errors
- **LocalStorage**: Cart is cached locally, clear if needed
- **Session Management**: Login to access protected features

---

**Need Help?** Check the documentation files or contact your administrator.

**Ready to Go?** Start with: `go build && ./app`

🚀 **Happy Shopping!**
