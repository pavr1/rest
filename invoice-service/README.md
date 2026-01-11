# Invoice Service

The Invoice Service manages purchase invoices, invoice details, and stock existences for the Bar-Restaurant system. This service implements **Approach 1** for inventory cost tracking where invoice input automatically creates stock items with calculated unit costs.

## Features

### Purchase Invoices
- Create and manage supplier purchase invoices
- Track invoice details (supplier, dates, amounts, status)
- Upload scanned invoice images
- Status tracking (pending, received, paid, cancelled)

### Invoice Details
- Add line items to purchase invoices
- Specify quantities, units, and pricing
- Link to existing or new stock items
- Automatic total price calculation

### Stock Existences
- Create stock batches from invoice details
- Automatic unit cost calculation
- Weighted average cost updates for stock items
- Track expiry dates and batch numbers
- Inventory quantity management

## Cost Calculation Logic

### Unit Cost Calculation
```
If items_per_unit specified (e.g., 12 breadsticks per "dozen"):
  units_purchased = quantity × items_per_unit
  cost_per_unit = unit_price ÷ items_per_unit

Example: 2 dozen breadsticks at $24/dozen
  units_purchased = 2 × 12 = 24 breadsticks
  cost_per_unit = $24 ÷ 12 = $2 per breadstick
```

### Weighted Average Cost
When new stock arrives, the system calculates weighted average costs:
```
current_value = current_stock × current_unit_cost
new_value = new_units × new_unit_cost
total_value = current_value + new_value
new_average_cost = total_value ÷ total_stock
```

## API Endpoints

### Purchase Invoices
- `GET /api/v1/invoices/purchase` - List invoices
- `POST /api/v1/invoices/purchase` - Create invoice
- `GET /api/v1/invoices/purchase/{id}` - Get invoice
- `PUT /api/v1/invoices/purchase/{id}` - Update invoice
- `DELETE /api/v1/invoices/purchase/{id}` - Delete invoice

### Invoice Details
- `GET /api/v1/invoices/purchase/{invoiceId}/details` - List details
- `POST /api/v1/invoices/purchase/{invoiceId}/details` - Create detail
- `GET /api/v1/invoices/purchase/{invoiceId}/details/{id}` - Get detail
- `PUT /api/v1/invoices/purchase/{invoiceId}/details/{id}` - Update detail
- `DELETE /api/v1/invoices/purchase/{invoiceId}/details/{id}` - Delete detail

### Existences
- `GET /api/v1/invoices/existences` - List existences
- `POST /api/v1/invoices/existences` - Create existence
- `GET /api/v1/invoices/existences/{id}` - Get existence
- `PUT /api/v1/invoices/existences/{id}` - Update existence
- `DELETE /api/v1/invoices/existences/{id}` - Delete existence

## Data Flow

1. **Create Purchase Invoice** - Supplier information and basic details
2. **Add Invoice Details** - Line items with stock item references
3. **Create Existences** - Convert invoice details to stock batches with calculated costs
4. **Update Stock Items** - Automatic inventory updates with weighted average costs

## Usage Example

```bash
# Create a purchase invoice
curl -X POST http://localhost:8092/api/v1/invoices/purchase \
  -H "Content-Type: application/json" \
  -d '{
    "invoice_number": "INV-001",
    "supplier_name": "Bakery Supply Co",
    "invoice_date": "2024-01-15",
    "total_amount": 240.00,
    "status": "pending"
  }'

# Add invoice detail (2 dozen breadsticks at $24/dozen)
curl -X POST http://localhost:8092/api/v1/invoices/purchase/{invoice_id}/details \
  -H "Content-Type: application/json" \
  -d '{
    "description": "Breadsticks",
    "quantity": 2,
    "unit_of_measure": "dozen",
    "items_per_unit": 12,
    "unit_price": 24.00
  }'

# Create existence (this automatically updates stock)
curl -X POST http://localhost:8092/api/v1/invoices/existences \
  -H "Content-Type: application/json" \
  -d '{
    "invoice_detail_id": "{detail_id}"
  }'
```

## Integration with Menu Costs

The calculated unit costs from existences are used by the menu service to calculate dish costs:

```
Menu Item Cost = Σ(ingredient_quantity × stock_unit_cost)
Selling Price = Menu Item Cost × markup_percentage
```

This provides accurate profitability analysis per menu item.