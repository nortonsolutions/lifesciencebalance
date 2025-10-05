# Migration and Setup Scripts

This directory contains utility scripts for the webstore conversion.

## Scripts

### migrate_courses_to_products.go
Converts existing Course records to Product records in the database.

**Usage:**
```bash
go run scripts/migrate_courses_to_products.go
```

**What it does:**
- Fetches all courses from the database
- Creates corresponding Product records
- Maps course fields to product fields
- Sets default values for new fields (price=0, in_stock=true)

**Note:** After migration, manually update product prices and images.

### seed_sample_products.go
Seeds the database with 10 sample products for testing and demonstration.

**Usage:**
```bash
go run scripts/seed_sample_products.go
```

**What it does:**
- Creates 10 diverse sample products across different categories
- Includes realistic product data (names, descriptions, prices)
- Demonstrates various product states (in stock, out of stock)
- Useful for testing the webstore functionality

## Prerequisites

Both scripts require:
- Valid `.env` file with `DATASTORE_PROJECT_ID` and `GOOGLE_APPLICATION_CREDENTIALS`
- Go 1.18 or higher
- Access to Google Cloud Datastore

## After Running Scripts

After successfully running these scripts:
1. Navigate to `/generic.html?page=products` to see products
2. Test adding products to cart
3. Verify cart functionality at `/generic.html?page=cart`
4. Update product images and prices as needed through the admin panel
