#!/bin/bash
# ============================================================
#  AutoWorkshop API — cURL Test Suite
#  Usage: bash curls.sh
#  Set TOKEN after login and update IDs as you go
# ============================================================

BASE="https://api.avinertech.com/api/v1"
TOKEN=""          # filled after login
CUSTOMER_ID=""
CAR_ID=""
ESTIMATE_ID=""
INVOICE_ID=""
PAYMENT_ID=""
STAFF_ID=""

# ── colour helpers ────────────────────────────────────────────
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

section() { echo -e "\n${BLUE}══════════════════════════════════════${NC}"; echo -e "${YELLOW}  $1${NC}"; echo -e "${BLUE}══════════════════════════════════════${NC}"; }
label()   { echo -e "\n${GREEN}▶ $1${NC}"; }


# ============================================================
# 1. AUTH
# ============================================================
section "1. AUTH"

label "POST /auth/login (valid)"
curl -s -X POST "$BASE/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"owner@workshop.com","password":"owner123"}' | jq .

label "POST /auth/login (wrong password)"
curl -s -X POST "$BASE/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"owner@workshop.com","password":"wrongpass"}' | jq .

label "POST /auth/login (missing fields)"
curl -s -X POST "$BASE/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"owner@workshop.com"}' | jq .

# ── Grab token automatically ──────────────────────────────────
echo -e "\n${YELLOW}  → Grabbing token...${NC}"
TOKEN=$(curl -s -X POST "$BASE/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"owner@workshop.com","password":"owner123"}' | jq -r '.token')
echo "  TOKEN=$TOKEN"

label "GET /auth/me"
curl -s "$BASE/auth/me" \
  -H "Authorization: Bearer $TOKEN" | jq .

label "POST /auth/logout"
curl -s -X POST "$BASE/auth/logout" \
  -H "Authorization: Bearer $TOKEN" | jq .

label "GET /auth/me (no token — expect 401)"
curl -s "$BASE/auth/me" | jq .


# ============================================================
# 2. SETTINGS
# ============================================================
section "2. SETTINGS"

label "GET /settings"
curl -s "$BASE/settings" \
  -H "Authorization: Bearer $TOKEN" | jq .

label "PUT /settings (owner only)"
curl -s -X PUT "$BASE/settings" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "workshopName": "AutoWorkshop Pro",
    "tagline": "Quality Auto Repair & Service",
    "address": "Plot 12, Johar Town, Lahore",
    "phone": "0300-1234567",
    "email": "info@autoworkshop.pk",
    "whatsapp": "0300-1234567",
    "ntn": "1234567-8",
    "defaultTaxPct": 0,
    "currency": "PKR",
    "terms": "Payment due within 7 days.\nWarranty void if tampered.",
    "warrantyDefault": "3 months or 5000 km",
    "footer": "Thank you for your business!",
    "estPrefix": "EST",
    "estStart": 1,
    "invPrefix": "INV",
    "invStart": 1
  }' | jq .


# ============================================================
# 3. STAFF
# ============================================================
section "3. STAFF"

label "GET /staff"
curl -s "$BASE/staff" \
  -H "Authorization: Bearer $TOKEN" | jq .

label "POST /staff (create mechanic)"
STAFF_RESP=$(curl -s -X POST "$BASE/staff" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "mechanic@workshop.com",
    "password": "mech1234",
    "name": "Usman Ali",
    "role": "mechanic",
    "phone": "0311-2222222",
    "active": true
  }')
echo $STAFF_RESP | jq .
STAFF_ID=$(echo $STAFF_RESP | jq -r '.id')
echo "  STAFF_ID=$STAFF_ID"

label "POST /staff (create receptionist)"
curl -s -X POST "$BASE/staff" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "reception@workshop.com",
    "password": "recep1234",
    "name": "Sara Malik",
    "role": "receptionist",
    "phone": "0322-3333333",
    "active": true
  }' | jq .

label "POST /staff (duplicate email — expect 409)"
curl -s -X POST "$BASE/staff" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "mechanic@workshop.com",
    "password": "mech1234",
    "name": "Duplicate",
    "role": "mechanic"
  }' | jq .

label "PUT /staff/:id (update phone)"
curl -s -X PUT "$BASE/staff/$STAFF_ID" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"phone": "0311-9999999"}' | jq .

label "PUT /staff/:id (change password)"
curl -s -X PUT "$BASE/staff/$STAFF_ID" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"password": "newpass456"}' | jq .

label "DELETE /staff/:id (soft delete)"
curl -s -X DELETE "$BASE/staff/$STAFF_ID" \
  -H "Authorization: Bearer $TOKEN" | jq .


# ============================================================
# 4. COMPANIES & JOB TYPES
# ============================================================
section "4. COMPANIES & JOB TYPES"

label "GET /companies"
curl -s "$BASE/companies" \
  -H "Authorization: Bearer $TOKEN" | jq .

label "POST /companies (Toyota)"
curl -s -X POST "$BASE/companies" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Toyota"}' | jq .

label "POST /companies (Honda)"
curl -s -X POST "$BASE/companies" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Honda"}' | jq .

label "POST /companies (Suzuki)"
curl -s -X POST "$BASE/companies" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Suzuki"}' | jq .

label "GET /job-types"
curl -s "$BASE/job-types" \
  -H "Authorization: Bearer $TOKEN" | jq .

label "PUT /job-types (replace list)"
curl -s -X PUT "$BASE/job-types" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '["Mechanical","Electrical","AC & Cooling","Body Work","Tyre & Wheel","Oil & Fluids","Detailing"]' | jq .


# ============================================================
# 5. CUSTOMERS
# ============================================================
section "5. CUSTOMERS"

label "POST /customers (create — Ahmed)"
CUST_RESP=$(curl -s -X POST "$BASE/customers" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Ahmed Raza",
    "phone": "0300-1111111",
    "email": "ahmed@example.com",
    "address": "House 5, Model Town, Lahore",
    "notes": "VIP customer"
  }')
echo $CUST_RESP | jq .
CUSTOMER_ID=$(echo $CUST_RESP | jq -r '.id')
echo "  CUSTOMER_ID=$CUSTOMER_ID"

label "POST /customers (create — Bilal)"
curl -s -X POST "$BASE/customers" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Bilal Hussain",
    "phone": "0311-2222222",
    "email": "bilal@example.com"
  }' | jq .

label "POST /customers (duplicate phone — expect 409)"
curl -s -X POST "$BASE/customers" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Dup","phone":"0300-1111111"}' | jq .

label "GET /customers (paginated)"
curl -s "$BASE/customers?page=1&limit=10" \
  -H "Authorization: Bearer $TOKEN" | jq .

label "GET /customers?q=ahmed (search)"
curl -s "$BASE/customers?q=ahmed" \
  -H "Authorization: Bearer $TOKEN" | jq .

label "GET /customers/:id"
curl -s "$BASE/customers/$CUSTOMER_ID" \
  -H "Authorization: Bearer $TOKEN" | jq .

label "GET /customers/by-phone/:phone"
curl -s "$BASE/customers/by-phone/03001111111" \
  -H "Authorization: Bearer $TOKEN" | jq .

label "GET /customers/by-phone/:phone (not found — expect 404)"
curl -s "$BASE/customers/by-phone/0399-0000000" \
  -H "Authorization: Bearer $TOKEN" | jq .

label "PUT /customers/:id (update address)"
curl -s -X PUT "$BASE/customers/$CUSTOMER_ID" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"address":"House 10, DHA Phase 5, Lahore"}' | jq .

label "DELETE /customers/:id (has cars — expect 409 after car is added)"
# Run this after adding a car to see the 409
# curl -s -X DELETE "$BASE/customers/$CUSTOMER_ID" \
#   -H "Authorization: Bearer $TOKEN" | jq .


# ============================================================
# 6. CARS / VEHICLES
# ============================================================
section "6. CARS"

label "POST /cars (create)"
CAR_RESP=$(curl -s -X POST "$BASE/cars" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"customerId\": \"$CUSTOMER_ID\",
    \"plate\": \"lea-1234\",
    \"chassis\": \"JT123456789\",
    \"company\": \"Toyota\",
    \"model\": \"Corolla\",
    \"year\": 2020,
    \"color\": \"#ffffff\",
    \"colorLabel\": \"White\",
    \"mileage\": 45000,
    \"notes\": \"Owner's daily driver\"
  }")
echo $CAR_RESP | jq .
CAR_ID=$(echo $CAR_RESP | jq -r '.id')
echo "  CAR_ID=$CAR_ID"

label "POST /cars (plate auto-uppercase check)"
curl -s -X POST "$BASE/cars" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"customerId\": \"$CUSTOMER_ID\",
    \"plate\": \"abc-0001\",
    \"company\": \"Honda\",
    \"model\": \"Civic\",
    \"year\": 2022,
    \"mileage\": 12000
  }" | jq .

label "GET /cars (list)"
curl -s "$BASE/cars?page=1&limit=10" \
  -H "Authorization: Bearer $TOKEN" | jq .

label "GET /cars?q=corolla (search)"
curl -s "$BASE/cars?q=corolla" \
  -H "Authorization: Bearer $TOKEN" | jq .

label "GET /cars?customerId=... (filter)"
curl -s "$BASE/cars?customerId=$CUSTOMER_ID" \
  -H "Authorization: Bearer $TOKEN" | jq .

label "GET /cars/:id"
curl -s "$BASE/cars/$CAR_ID" \
  -H "Authorization: Bearer $TOKEN" | jq .

label "GET /customers/:id/cars"
curl -s "$BASE/customers/$CUSTOMER_ID/cars" \
  -H "Authorization: Bearer $TOKEN" | jq .

label "PUT /cars/:id (update mileage)"
curl -s -X PUT "$BASE/cars/$CAR_ID" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"mileage": 47500}' | jq .

label "DELETE /cars/:id (has invoices — expect 409 after invoice added)"
# curl -s -X DELETE "$BASE/cars/$CAR_ID" \
#   -H "Authorization: Bearer $TOKEN" | jq .


# ============================================================
# 7. ESTIMATES
# ============================================================
section "7. ESTIMATES"

label "POST /estimates (create)"
EST_RESP=$(curl -s -X POST "$BASE/estimates" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"date\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\",
    \"customerId\": \"$CUSTOMER_ID\",
    \"carId\": \"$CAR_ID\",
    \"jobTypes\": [\"Mechanical\", \"Oil & Fluids\"],
    \"status\": \"Draft\",
    \"notes\": \"Full service + oil change\",
    \"discountAmount\": 500,
    \"taxPct\": 0,
    \"items\": [
      {
        \"name\": \"Engine Oil (4L)\",
        \"description\": \"Castrol 5W-30\",
        \"qty\": 1,
        \"unitPrice\": 4500,
        \"discountPct\": 0
      },
      {
        \"name\": \"Oil Filter\",
        \"qty\": 1,
        \"unitPrice\": 800,
        \"discountPct\": 0
      },
      {
        \"name\": \"Labour\",
        \"qty\": 1,
        \"unitPrice\": 1000,
        \"discountPct\": 10
      }
    ],
    \"images\": []
  }")
echo $EST_RESP | jq .
ESTIMATE_ID=$(echo $EST_RESP | jq -r '.id')
echo "  ESTIMATE_ID=$ESTIMATE_ID"

label "GET /estimates (list)"
curl -s "$BASE/estimates?page=1&limit=10" \
  -H "Authorization: Bearer $TOKEN" | jq .

label "GET /estimates?status=Draft"
curl -s "$BASE/estimates?status=Draft" \
  -H "Authorization: Bearer $TOKEN" | jq .

label "GET /estimates/:id"
curl -s "$BASE/estimates/$ESTIMATE_ID" \
  -H "Authorization: Bearer $TOKEN" | jq .

label "GET /cars/:id/estimates"
curl -s "$BASE/cars/$CAR_ID/estimates" \
  -H "Authorization: Bearer $TOKEN" | jq .

label "PUT /estimates/:id (update status to Sent)"
curl -s -X PUT "$BASE/estimates/$ESTIMATE_ID" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"status":"Sent","notes":"Sent to customer via WhatsApp"}' | jq .

label "PUT /estimates/:id (update items)"
curl -s -X PUT "$BASE/estimates/$ESTIMATE_ID" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "Approved",
    "discountAmount": 200,
    "items": [
      {
        "name": "Engine Oil (4L)",
        "description": "Castrol 5W-30",
        "qty": 1,
        "unitPrice": 4500,
        "discountPct": 0
      },
      {
        "name": "Oil Filter",
        "qty": 1,
        "unitPrice": 800,
        "discountPct": 0
      },
      {
        "name": "Labour",
        "qty": 1,
        "unitPrice": 1000,
        "discountPct": 0
      }
    ]
  }' | jq .

label "POST /estimates/:id/convert (to invoice)"
INV_RESP=$(curl -s -X POST "$BASE/estimates/$ESTIMATE_ID/convert" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"dueDate\": \"$(date -u -d '+7 days' +%Y-%m-%dT%H:%M:%SZ 2>/dev/null || date -u -v+7d +%Y-%m-%dT%H:%M:%SZ)\",
    \"deliveryDate\": null,
    \"finalMileage\": 47600,
    \"warranty\": \"3 months or 5000 km\"
  }")
echo $INV_RESP | jq .
INVOICE_ID=$(echo $INV_RESP | jq -r '.id')
echo "  INVOICE_ID=$INVOICE_ID"

label "POST /estimates/:id/convert (already converted — expect 409)"
curl -s -X POST "$BASE/estimates/$ESTIMATE_ID/convert" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{}' | jq .

label "DELETE /estimates/:id (converted — expect 409)"
curl -s -X DELETE "$BASE/estimates/$ESTIMATE_ID" \
  -H "Authorization: Bearer $TOKEN" | jq .

label "POST /estimates (another one to test delete)"
DEL_EST=$(curl -s -X POST "$BASE/estimates" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"date\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\",
    \"customerId\": \"$CUSTOMER_ID\",
    \"carId\": \"$CAR_ID\",
    \"status\": \"Draft\",
    \"discountAmount\": 0,
    \"taxPct\": 0,
    \"items\": [{\"name\":\"Brake Pads\",\"qty\":1,\"unitPrice\":3500,\"discountPct\":0}],
    \"images\": []
  }")
DEL_EST_ID=$(echo $DEL_EST | jq -r '.id')

label "DELETE /estimates/:id (draft — should succeed)"
curl -s -X DELETE "$BASE/estimates/$DEL_EST_ID" \
  -H "Authorization: Bearer $TOKEN" | jq .


# ============================================================
# 8. INVOICES
# ============================================================
section "8. INVOICES"

# Use invoice from conversion above, or create fresh
label "POST /invoices (create directly)"
INV_DIRECT=$(curl -s -X POST "$BASE/invoices" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"date\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\",
    \"customerId\": \"$CUSTOMER_ID\",
    \"carId\": \"$CAR_ID\",
    \"jobTypes\": [\"Electrical\"],
    \"status\": \"Unpaid\",
    \"notes\": \"AC compressor repair\",
    \"discountAmount\": 0,
    \"taxPct\": 0,
    \"finalMileage\": 47800,
    \"warranty\": \"1 month\",
    \"items\": [
      {
        \"name\": \"AC Compressor\",
        \"qty\": 1,
        \"unitPrice\": 25000,
        \"discountPct\": 0
      },
      {
        \"name\": \"Labour\",
        \"qty\": 1,
        \"unitPrice\": 3000,
        \"discountPct\": 0
      }
    ],
    \"images\": []
  }")
echo $INV_DIRECT | jq .
INVOICE_ID2=$(echo $INV_DIRECT | jq -r '.id')
echo "  INVOICE_ID2=$INVOICE_ID2"

label "GET /invoices (list)"
curl -s "$BASE/invoices?page=1&limit=10" \
  -H "Authorization: Bearer $TOKEN" | jq .

label "GET /invoices?status=Unpaid"
curl -s "$BASE/invoices?status=Unpaid" \
  -H "Authorization: Bearer $TOKEN" | jq .

label "GET /invoices/:id"
curl -s "$BASE/invoices/$INVOICE_ID2" \
  -H "Authorization: Bearer $TOKEN" | jq .

label "GET /customers/:id/invoices"
curl -s "$BASE/customers/$CUSTOMER_ID/invoices" \
  -H "Authorization: Bearer $TOKEN" | jq .

label "GET /cars/:id/invoices"
curl -s "$BASE/cars/$CAR_ID/invoices" \
  -H "Authorization: Bearer $TOKEN" | jq .

label "PUT /invoices/:id (update notes + warranty)"
curl -s -X PUT "$BASE/invoices/$INVOICE_ID2" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"notes":"AC compressor replaced. Tested OK.","warranty":"2 months"}' | jq .


# ============================================================
# 9. PAYMENTS
# ============================================================
section "9. PAYMENTS"

label "POST /invoices/:id/payments (partial — 10000)"
PAY_RESP=$(curl -s -X POST "$BASE/invoices/$INVOICE_ID2/payments" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"date\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\",
    \"amount\": 10000,
    \"method\": \"Cash\",
    \"note\": \"Advance payment\"
  }")
echo $PAY_RESP | jq .
PAYMENT_ID=$(echo $PAY_RESP | jq -r '.payment.id')
echo "  PAYMENT_ID=$PAYMENT_ID"
echo "  Invoice status should be: Partially Paid"

label "POST /invoices/:id/payments (remaining — 18000)"
curl -s -X POST "$BASE/invoices/$INVOICE_ID2/payments" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"date\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\",
    \"amount\": 18000,
    \"method\": \"Bank Transfer\",
    \"note\": \"Final payment\"
  }" | jq .
echo "  Invoice status should be: Paid"

label "POST /invoices/:id/payments (overpayment — change scenario)"
# Create a new invoice to test overpayment
OVR_INV=$(curl -s -X POST "$BASE/invoices" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"date\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\",
    \"customerId\": \"$CUSTOMER_ID\",
    \"carId\": \"$CAR_ID\",
    \"discountAmount\": 0,
    \"taxPct\": 0,
    \"items\": [{\"name\":\"Wiper Blades\",\"qty\":2,\"unitPrice\":800,\"discountPct\":0}],
    \"images\": []
  }")
OVR_INV_ID=$(echo $OVR_INV | jq -r '.id')
curl -s -X POST "$BASE/invoices/$OVR_INV_ID/payments" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"date\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\",
    \"amount\": 2000,
    \"method\": \"Cash\",
    \"note\": \"Customer paid 2000, change 400\"
  }" | jq .
echo "  Grand=1600, Paid=2000 → balance=0, status=Paid, change=400 (frontend shows change)"

label "DELETE /invoices/:id/payments/:paymentId (reverse payment)"
curl -s -X DELETE "$BASE/invoices/$INVOICE_ID2/payments/$PAYMENT_ID" \
  -H "Authorization: Bearer $TOKEN" | jq .
echo "  Invoice status should revert to Partially Paid"

label "DELETE /invoices/:id (owner only)"
curl -s -X DELETE "$BASE/invoices/$OVR_INV_ID" \
  -H "Authorization: Bearer $TOKEN" | jq .


# ============================================================
# 10. REPORTS
# ============================================================
section "10. REPORTS"

label "GET /reports/dashboard"
curl -s "$BASE/reports/dashboard" \
  -H "Authorization: Bearer $TOKEN" | jq .

label "GET /reports/revenue?groupBy=day"
curl -s "$BASE/reports/revenue?from=2026-01-01T00:00:00Z&to=2026-12-31T23:59:59Z&groupBy=day" \
  -H "Authorization: Bearer $TOKEN" | jq .

label "GET /reports/revenue?groupBy=month"
curl -s "$BASE/reports/revenue?groupBy=month&from=2026-01-01T00:00:00Z&to=2026-12-31T23:59:59Z" \
  -H "Authorization: Bearer $TOKEN" | jq .

label "GET /reports/top-customers"
curl -s "$BASE/reports/top-customers?limit=5" \
  -H "Authorization: Bearer $TOKEN" | jq .

label "GET /reports/mechanic-performance"
curl -s "$BASE/reports/mechanic-performance?from=2026-01-01T00:00:00Z&to=2026-12-31T23:59:59Z" \
  -H "Authorization: Bearer $TOKEN" | jq .


# ============================================================
# 11. ROLE / AUTH GUARD TESTS
# ============================================================
section "11. ROLE & GUARD TESTS"

label "Login as mechanic"
MECH_TOKEN=$(curl -s -X POST "$BASE/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"mechanic@workshop.com","password":"mech1234"}' | jq -r '.token')
echo "  MECH_TOKEN=$MECH_TOKEN"

label "Mechanic tries GET /staff (expect 403)"
curl -s "$BASE/staff" \
  -H "Authorization: Bearer $MECH_TOKEN" | jq .

label "Mechanic tries DELETE /invoices/:id (expect 403)"
curl -s -X DELETE "$BASE/invoices/$INVOICE_ID2" \
  -H "Authorization: Bearer $MECH_TOKEN" | jq .

label "Mechanic tries PUT /settings (expect 403)"
curl -s -X PUT "$BASE/settings" \
  -H "Authorization: Bearer $MECH_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"workshopName":"Hacked"}' | jq .

label "No token → GET /customers (expect 401)"
curl -s "$BASE/customers" | jq .

label "Garbage token (expect 401)"
curl -s "$BASE/customers" \
  -H "Authorization: Bearer thisisnotatoken" | jq .


echo -e "\n${BLUE}══════════════════════════════════════${NC}"
echo -e "${GREEN}  ✅ All curls done${NC}"
echo -e "${BLUE}══════════════════════════════════════${NC}\n"
