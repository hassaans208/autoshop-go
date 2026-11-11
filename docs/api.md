AutoWorkshop Pro — Backend API Specification
This document lists every REST endpoint the frontend expects from the Go backend. All endpoints return JSON. All timestamps are ISO-8601 strings (UTC). All money fields are numbers in PKR (no formatting). IDs are server-generated strings (UUID recommended).

0. Global Conventions
Base URL: /api/v1

Auth: All endpoints except POST /auth/login require Authorization: Bearer <jwt>.

Errors: { "error": "message", "code": "string", "details"?: any } with appropriate HTTP status (400, 401, 403, 404, 409, 422, 500).

Pagination (lists): query ?page=1&limit=50&q=search. Response: { data: [...], total: number, page, limit }.

Roles: owner | manager | mechanic | receptionist | accountant.

1. Authentication
POST
/auth/login
Authenticates a staff member by email + password. Returns JWT and the user profile.

Request Body
{
  "email": "owner@workshop.com",
  "password": "owner123"
}
Response (200/201)
{
  "token": "<jwt>",
  "user": {
    "id": "uuid",
    "email": "owner@workshop.com",
    "name": "Ahmed Khan",
    "role": "owner",
    "phone": "0300-1111111",
    "active": true
  }
}
POST
/auth/logout
Invalidates the current token (optional — frontend just drops it).

Response (200/201)
{ "ok": true }
GET
/auth/me
Returns the currently authenticated user. Used to bootstrap the app on page reload.

Response (200/201)
{
  "id": "uuid", "email": "...", "name": "...", "role": "owner",
  "phone": "...", "active": true
}
2. Workshop Settings
GET
/settings
Returns workshop-wide settings (used by header, prints, invoice prefixes, tax).

Response (200/201)
{
  "workshopName": "AutoWorkshop Pro",
  "logo": "data:image/png;base64,...",
  "tagline": "Quality Auto Repair & Service",
  "address": "...", "phone": "...", "email": "...", "whatsapp": "...",
  "ntn": "1234567-8",
  "defaultTaxPct": 0,
  "currency": "PKR",
  "terms": "string (newline-separated)",
  "warrantyDefault": "3 months ...",
  "footer": "Thank you ...",
  "estPrefix": "EST", "estStart": 1,
  "invPrefix": "INV", "invStart": 1
}
PUT
/settings
Replaces workshop settings. Owner only.

Request Body
{ /* full Settings object — same shape as GET */ }
Response (200/201)
{ /* updated Settings */ }
3. Lookups (Companies & Job Types)
GET
/companies
List of car brand names (for the New Car dropdown).

Response (200/201)
["Toyota","Honda","Suzuki", ...]
POST
/companies
Add a brand if it does not exist (case-insensitive).

Request Body
{ "name": "Lexus" }
Response (200/201)
{ "name": "Lexus" }
GET
/job-types
List of job-type tags.

Response (200/201)
["Mechanical","Electrical", ...]
PUT
/job-types
Replace the full list (manager+ only).

Request Body
["Mechanical","Electrical", ...]
Response (200/201)
["Mechanical", ...]
4. Staff / Users
GET
/staff
List all staff. Owner/Manager only. Passwords MUST never be returned.

Response (200/201)
[{ "id":"uuid","email":"...","name":"...","role":"mechanic","phone":"...","active":true }]
POST
/staff
Create a staff user. Owner only.

Request Body
{
  "email": "new@workshop.com",
  "password": "plain-text-on-create",
  "name": "...",
  "role": "mechanic",
  "phone": "0300-...",
  "active": true
}
Response (200/201)
{ "id":"uuid", /* user without password */ }
PUT
/staff/:id
Update user. password is optional; if present, rehash.

Request Body
{ "name":"...", "role":"manager", "phone":"...", "active":true, "password"?:"new" }
Response (200/201)
{ "id":"uuid", /* user */ }
DELETE
/staff/:id
Soft-delete (set active=false). Owner only.

Response (200/201)
{ "ok": true }
5. Customers
GET
/customers
Paginated list. Supports ?q= (search by name/phone).

Response (200/201)
{ "data": [{ "id":"uuid","name":"...","phone":"0300-...","email":"...","address":"...","notes":"...","createdAt":"ISO" }], "total":123 }
GET
/customers/:id
Single customer.

Response (200/201)
{ "id":"uuid","name":"...","phone":"...","email":"...","address":"...","notes":"...","createdAt":"ISO" }
GET
/customers/by-phone/:phone
Used by the New Job flow to look up existing customers. Returns 404 if not found.

Response (200/201)
{ /* Customer */ }
POST
/customers
Create. Must 409 if phone already exists.

Request Body
{ "name":"...", "phone":"0300-...", "email"?:"...", "address"?:"...", "notes"?:"..." }
Response (200/201)
{ /* Customer with id + createdAt */ }
PUT
/customers/:id
Partial update accepted.

Request Body
{ "name"?:"...", "phone"?:"...", "email"?:"...", "address"?:"...", "notes"?:"..." }
Response (200/201)
{ /* Customer */ }
DELETE
/customers/:id
Hard delete only if no cars/invoices exist; else 409.

Response (200/201)
{ "ok": true }
6. Cars / Vehicles
GET
/cars
Paginated list. Supports ?q= (plate/company/model) and ?customerId=.

Response (200/201)
{ "data":[{ "id":"uuid","customerId":"uuid","plate":"LEA-1234","chassis":"...","company":"Toyota","model":"Corolla","year":2020,"color":"#ffffff","colorLabel":"White","mileage":45000,"notes":"...","createdAt":"ISO" }], "total":50 }
GET
/cars/:id
Single car.

Response (200/201)
{ /* Car */ }
GET
/customers/:id/cars
All cars belonging to a customer.

Response (200/201)
[{ /* Car */ }]
POST
/cars
Create. Plate is stored uppercase. Auto-adds company to /companies if new.

Request Body
{
  "customerId":"uuid",
  "plate":"lea-1234",
  "chassis"?:"...", "company":"Toyota", "model"?:"Corolla",
  "year"?:2020, "color"?:"#ffffff", "colorLabel"?:"White",
  "mileage"?:45000, "notes"?:"..."
}
Response (200/201)
{ /* Car */ }
PUT
/cars/:id
Partial update. Plate normalised to uppercase.

Request Body
{ /* any Car fields */ }
Response (200/201)
{ /* Car */ }
DELETE
/cars/:id
Block if invoices exist; else delete.

Response (200/201)
{ "ok": true }
7. Shared Document Types (Estimate & Invoice)
LineItem {
  id: string
  name: string
  description?: string
  qty: number
  unitPrice: number          // PKR
  discountPct: number        // 0-100
}

ImageItem {
  id: string
  label: string
  dataUrl: string            // base64 data URL OR replace with "url": string after backend stores files
}

Payment {
  id: string
  date: string               // ISO
  amount: number             // PKR (can EXCEED invoice total — change-returned is computed client side)
  method: "Cash" | "Bank Transfer" | "Cheque" | "Card" | "Split"
  note?: string
}

Totals (computed server-side AND echoed in responses, for safety):
  subtotal     = sum(qty * unitPrice * (1 - discountPct/100))
  afterDisc    = max(0, subtotal - discountAmount)
  tax          = afterDisc * taxPct/100
  grand        = afterDisc + tax
  paid         = sum(payments.amount)
  balance      = max(0, grand - paid)
  change       = max(0, paid - grand)
Numbering: EST-YYYY-00001 and INV-YYYY-00001. The backend owns the counter — frontend never sends number on create.

8. Estimates
GET
/estimates
Paginated list. Filters: ?status=&customerId=&carId=&from=ISO&to=ISO&q=number.

Response (200/201)
{ "data":[{ /* Estimate (see POST response) */ }], "total":100 }
GET
/estimates/:id
Single estimate (full doc).

Response (200/201)
{ /* Estimate */ }
GET
/cars/:id/estimates
Estimates for a car.

Response (200/201)
[{ /* Estimate */ }]
POST
/estimates
Create. Backend assigns id + number.

Request Body
{
  "date": "ISO",
  "customerId": "uuid",
  "carId": "uuid",
  "jobTypes": ["Mechanical","AC & Cooling"],
  "mechanicId"?: "uuid",
  "status": "Draft" | "Sent" | "Approved" | "Rejected" | "Converted",
  "notes"?: "string",
  "items": [LineItem],
  "images": [ImageItem],
  "discountAmount": 0,
  "taxPct": 0
}
Response (200/201)
{
  "id":"uuid",
  "number":"EST-2026-00001",
  "date":"ISO",
  "customerId":"uuid","carId":"uuid",
  "jobTypes":[...], "mechanicId":"uuid|null",
  "status":"Draft",
  "notes":"...", "items":[...], "images":[...],
  "discountAmount":0, "taxPct":0,
  "convertedInvoiceId": null,
  "totals": { "subtotal":0, "tax":0, "grand":0 }
}
PUT
/estimates/:id
Partial update. Any field from POST body is accepted.

Request Body
{ /* partial */ }
Response (200/201)
{ /* Estimate */ }
DELETE
/estimates/:id
Hard delete. 409 if status='Converted'.

Response (200/201)
{ "ok": true }
POST
/estimates/:id/convert
Converts an estimate into an Invoice. Copies items/images/discount/tax. Sets estimate.status='Converted' and links convertedInvoiceId. Returns the new invoice.

Request Body
{ "dueDate"?:"ISO", "deliveryDate"?:"ISO", "finalMileage"?:0, "warranty"?:"..." }
Response (200/201)
{ /* Invoice (see §9 POST response) with fromEstimateId set */ }
9. Invoices
GET
/invoices
Paginated list. Filters: ?status=&customerId=&carId=&from=&to=&q=number.

Response (200/201)
{ "data":[{ /* Invoice */ }], "total":200 }
GET
/invoices/:id
Single invoice (full doc including payments).

Response (200/201)
{ /* Invoice */ }
GET
/customers/:id/invoices
Invoices for a customer.

Response (200/201)
[{ /* Invoice */ }]
GET
/cars/:id/invoices
Invoices for a car.

Response (200/201)
[{ /* Invoice */ }]
POST
/invoices
Create. Backend assigns id + number. payments starts as [].

Request Body
{
  "date":"ISO",
  "dueDate"?:"ISO",
  "customerId":"uuid","carId":"uuid",
  "jobTypes":[...], "mechanicId"?:"uuid",
  "status":"Unpaid",
  "notes"?:"...",
  "items":[LineItem],
  "images":[ImageItem],
  "discountAmount":0, "taxPct":0,
  "deliveryDate"?:"ISO",
  "finalMileage"?:0,
  "warranty"?:"...",
  "fromEstimateId"?:"uuid"
}
Response (200/201)
{
  "id":"uuid", "number":"INV-2026-00001",
  "date":"ISO", "dueDate":"ISO|null",
  "customerId":"uuid","carId":"uuid",
  "jobTypes":[...], "mechanicId":"uuid|null",
  "status":"Unpaid" | "Partially Paid" | "Paid" | "Overdue",
  "notes":"...", "items":[...], "images":[...],
  "discountAmount":0, "taxPct":0,
  "deliveryDate":"ISO|null", "finalMileage":0, "warranty":"...",
  "payments":[],
  "fromEstimateId":"uuid|null",
  "totals": { "subtotal":0, "tax":0, "grand":0, "paid":0, "balance":0, "change":0 }
}
PUT
/invoices/:id
Partial update (any field from POST). Recompute status from payments vs grand.

Request Body
{ /* partial */ }
Response (200/201)
{ /* Invoice */ }
DELETE
/invoices/:id
Hard delete. Owner/Manager only.

Response (200/201)
{ "ok": true }
10. Payments (sub-resource of Invoice)
Payments may exceed the invoice grand total — the frontend displays the excess as Change Returned. Backend should still store the full amount as received and set status='Paid' when paid >= grand.

POST
/invoices/:id/payments
Record a payment.

Request Body
{
  "date":"ISO",
  "amount": 5000,
  "method": "Cash" | "Bank Transfer" | "Cheque" | "Card" | "Split",
  "note"?: "Change returned: Rs. 200"
}
Response (200/201)
{
  "payment": { "id":"uuid","date":"ISO","amount":5000,"method":"Cash","note":"..." },
  "invoice": { /* full updated Invoice with recomputed status + totals */ }
}
DELETE
/invoices/:id/payments/:paymentId
Reverse a payment. Recompute status.

Response (200/201)
{ "invoice": { /* Invoice */ } }
11. File / Image Upload (optional upgrade)
Today the frontend embeds images as base64 inside the document. If you want a real file store, expose this endpoint and the frontend will switch ImageItem.dataUrl to a URL field.

POST
/uploads
Multipart upload. field name: file. Returns a public URL.

Request Body
multipart/form-data:
  file: <binary>
  purpose?: "estimate" | "invoice" | "logo"
Response (200/201)
{ "url": "https://cdn.example.com/uploads/xyz.jpg", "id": "uuid" }
12. Reports / Dashboard
GET
/reports/dashboard
Numbers shown on the dashboard home.

Response (200/201)
{
  "todayRevenue": 0,
  "monthRevenue": 0,
  "openInvoices": 0,
  "overdueInvoices": 0,
  "estimatesPending": 0,
  "activeJobs": 0,
  "recentInvoices": [{ /* Invoice summary */ }],
  "recentEstimates": [{ /* Estimate summary */ }]
}
GET
/reports/revenue
?from=ISO&to=ISO&groupBy=day|month. Used by Reports page chart.

Response (200/201)
{
  "series": [{ "bucket":"2026-06-01","revenue":12500,"paid":10000,"invoices":4 }],
  "totals": { "revenue":350000, "paid":280000, "outstanding":70000, "invoices":120 }
}
GET
/reports/top-customers
?limit=10&from=&to=.

Response (200/201)
[{ "customerId":"uuid","name":"...","invoices":5,"revenue":120000 }]
GET
/reports/mechanic-performance
?from=&to=.

Response (200/201)
[{ "mechanicId":"uuid","name":"...","jobs":12,"revenue":80000 }]
13. Implementation Notes for the Go Backend
All write endpoints must validate role permissions. Suggested matrix: owner=all, manager=all except staff CRUD & settings, receptionist=customers/cars/estimates/invoices/payments (no delete), mechanic=read-only + update items on assigned jobs, accountant=read all + payments.
Always recompute totals server-side from items+discountAmount+taxPct — never trust client.
Invoice status transitions: payment recorded → recompute (paid>=grand = Paid, paid>0 = Partially Paid, else Unpaid); a nightly job flips past-due Unpaid/Partially Paid to Overdue.
Number generation must be atomic (DB sequence per prefix per year).
Phone search: strip non-digits before comparing.
Plates: store uppercase; index for fast lookup.
CORS: allow the deployed Lovable frontend origin + Authorization header.
Return ISO-8601 with timezone (e.g. 2026-06-17T10:30:00Z).
11. Vendors (NEW)
A vendor is either a company (e.g. Toyota Pakistan, manufacturer) or a normal vendor (parts shop). Used to attribute purchase cost of parts.

GET
/vendors
List vendors. Supports ?q=search&kind=company|vendor.

Response (200/201)
{ "data": [ { "id": "uuid", "name": "Toyota Pakistan", "kind": "company", "phone": "+92...", "email": "...", "notes": "..." } ], "total": 1 }
POST
/vendors
Create vendor or company.

Request Body
{ "name": "Local Parts Bazaar", "kind": "vendor", "phone": "+92 300 ...", "email": "?", "notes": "?" }
Response (200/201)
{ "id": "uuid", "name": "...", "kind": "vendor", "phone": "...", "email": "...", "notes": "..." }
GET
/vendors/:id
Get one vendor.

Response (200/201)
{ "id": "...", "name": "...", "kind": "company|vendor", ... }
PUT
/vendors/:id
Update vendor.

Request Body
{ "name": "...", "kind": "vendor", "phone": "...", "email": "...", "notes": "..." }
Response (200/201)
{ "id": "...", ... }
DELETE
/vendors/:id
Delete vendor. 409 if any Part still references it.

Response (200/201)
{ "ok": true }
12. Parts Catalog (NEW)
Reusable parts catalog. Each part can optionally link to a vendor/company that supplies it. Default purchase & sale prices auto-fill into wizard line items.

GET
/parts
List parts. Supports ?q=search&vendorId=uuid.

Response (200/201)
{ "data": [ { "id": "uuid", "name": "Engine Oil 5L", "vendorId": "uuid|null", "defaultPurchasePrice": 4500, "defaultSalePrice": 5500, "notes": "?" } ], "total": 1 }
POST
/parts
Add a part.

Request Body
{ "name": "Brake Pads", "vendorId": "uuid|null", "defaultPurchasePrice": 3000, "defaultSalePrice": 4200, "notes": "?" }
Response (200/201)
{ "id": "uuid", "name": "...", "vendorId": "...", "defaultPurchasePrice": 3000, "defaultSalePrice": 4200 }
GET
/parts/:id
Get one part.

Response (200/201)
{ "id": "...", ... }
PUT
/parts/:id
Update part.

Request Body
{ "name": "...", "vendorId": "uuid|null", "defaultPurchasePrice": 0, "defaultSalePrice": 0, "notes": "?" }
Response (200/201)
{ "id": "...", ... }
DELETE
/parts/:id
Delete part.

Response (200/201)
{ "ok": true }
13. Job Types with Subtypes (EXTENDED)
Job types now carry a list of subtypes. Wizard step 1 picks job types; selected subtypes become labour rows in step 4.

GET
/job-types
List job types with subtypes.

Response (200/201)
{ "data": [ { "name": "Engine", "subtypes": ["Oil Change", "Tune-up", "Overhaul"] }, { "name": "Brakes", "subtypes": ["Pad Replacement"] } ] }
PUT
/job-types
Replace the full list.

Request Body
{ "list": [ { "name": "Engine", "subtypes": ["Oil Change", "Overhaul"] } ] }
Response (200/201)
{ "data": [ ... ] }
14. Staff Salary (EXTENDED)
Each staff member now has a monthly salary (PKR). Dashboard uses it for obligation / remaining-to-cover / profit.

GET
/staff
List staff including salary.

Response (200/201)
{ "data": [ { "id": "uuid", "name": "...", "email": "...", "phone": "...", "role": "mechanic", "active": true, "salary": 50000 } ] }
POST
/staff
Create staff with salary.

Request Body
{ "name": "...", "email": "...", "phone": "...", "password": "...", "role": "mechanic", "salary": 50000, "active": true }
Response (200/201)
{ "id": "uuid", "name": "...", "role": "mechanic", "salary": 50000, "active": true }
PUT
/staff/:id
Update staff (salary, role, active, password).

Request Body
{ "name": "?", "email": "?", "phone": "?", "role": "?", "salary": 55000, "active": true, "password": "?" }
Response (200/201)
{ "id": "...", ... }
15. Dashboard Financials (NEW)
Single endpoint that returns everything the dashboard needs — no client-side aggregation.

GET
/dashboard/financials
Current-month financials. Optional ?month=YYYY-MM.

Response (200/201)
{
  "month": "2026-06",
  "salaryObligation": 360000,       // sum of active staff salaries
  "revenue": 480000,                // sum of PAID amounts on invoices this month
  "remainingToCoverSalaries": 0,    // max(0, salaryObligation - revenue)
  "partsPurchaseCost": 95000,       // sum(items.purchasePrice * qty) on invoices this month
  "materialCost": 30000,            // sum(items.materialPrice) on labour items this month
  "profit": 0,                      // max(0, revenue - salaryObligation - partsPurchaseCost - materialCost)
  "counts": {
    "customers": 142, "cars": 198,
    "openEstimates": 12, "unpaidInvoices": 7, "overdueInvoices": 2
  }
}
16. Extended LineItem shape (BREAKING for estimates/invoices)
Existing items[] on Estimates & Invoices grows new optional fields driven by the wizard. Backend must accept & persist them as-is.

LineItem (extended):
{
  "id": "uuid",
  "name": "Engine Oil 5L",
  "description": "?",
  "qty": 1,
  "unitPrice": 5500,           // sale price per unit
  "discountPct": 0,            // legacy %; wizard uses discountAmount
  "discountAmount": 0,         // NEW — flat PKR discount on this row

  // NEW — part rows (kind = "part")
  "kind": "part" | "labour",
  "purchasePrice": 4500,       // INTERNAL — never on public URLs
  "vendorId": "uuid|null",     // INTERNAL — never on public URLs

  // NEW — labour rows (kind = "labour")
  "jobType": "Engine",
  "subtype": "Oil Change",
  "assignees": ["staffId1","staffId2"],   // INTERNAL
  "materialPrice": 200,                    // INTERNAL (cost)
  "servicePrice": 1500                     // shown to customer as labour
}
Totals are still server-recomputed. Tax is now a flat amount (taxAmount); keep taxPct for backward compat but prefer taxAmount when present.

17. Estimates & Invoices — additions
Existing endpoints keep working. Additions to payload/response:

publicSlug — server-generated 8-char URL-safe string returned on create. Used for the public link.
taxAmount — number PKR. Replaces taxPct when present.
POST
/estimates
Create estimate from wizard.

Request Body
{
  "date": "2026-06-17",
  "customerId": "uuid",
  "carId": "uuid",
  "jobTypes": ["Engine"],
  "subtypes": ["Oil Change","Tune-up"],
  "notes": "internal notes",
  "items": [ /* extended LineItem[] — see section 16 */ ],
  "images": [ { "id":"...", "label":"Before", "dataUrl":"data:image/jpeg;base64,..." } ],
  "discountAmount": 500,
  "taxAmount": 0,
  "status": "Draft"
}
Response (200/201)
{ "id":"uuid", "number":"EST-1001", "publicSlug":"a1b2c3d4", "status":"Draft", ... }
POST
/invoices
Create invoice from wizard.

Request Body
{
  "date":"...", "dueDate":"?", "deliveryDate":"?", "finalMileage":0, "warranty":"?",
  "customerId":"uuid", "carId":"uuid",
  "jobTypes":[...], "subtypes":[...], "items":[...], "images":[...],
  "discountAmount":0, "taxAmount":0,
  "status":"Unpaid",
  "fromEstimateId":"uuid|null"
}
Response (200/201)
{ "id":"uuid", "number":"INV-5001", "publicSlug":"e5f6g7h8", ... }
POST
/estimates/:id/convert
Convert estimate → invoice.

Request Body
{ "dueDate":"?", "deliveryDate":"?", "finalMileage":0, "warranty":"?" }
Response (200/201)
{ "invoiceId":"uuid", "number":"INV-5002", "publicSlug":"..." }
18. Public (unauthenticated) URLs (NEW)
Customer-facing read-only views. No auth header. Backend MUST strip every INTERNAL field (purchasePrice, vendorId, materialPrice, assignees, internal notes) before returning. Frontend builds the WhatsApp link as https://wa.me/?text=<encoded url>.

GET
/estimates/slug/:slug
Fetch estimate by publicSlug. NO auth header. Sanitized payload. (Frontend calls GET /api/v1/estimates/slug/{slug}.)

Response (200/201)
{
  "id": "...", "number": "EST-1001", "status": "Sent",
  "date": "2026-06-17",
  "publicSlug": "na626u0v",
  "customerId": "...", "carId": "...",
  "jobTypes": ["Engine"], "subtypes": ["Engine / Oil Change"],
  "items": [
    { "id":"...", "kind":"part",   "name":"Engine Oil 5L", "qty":1, "unitPrice":5500, "discountAmount":0 },
    { "id":"...", "kind":"labour", "name":"Engine / Oil Change", "subtype":"Engine / Oil Change", "qty":1, "unitPrice":1500, "servicePrice":1500 }
  ],
  "images": [],
  "discountAmount": 0, "taxPct": 0, "taxAmount": 0,
  "notes": "(omit internal notes)",
  "workshop": { "name":"...", "logo":"...", "phone":"...", "address":"...", "ntn":"...", "terms":"...", "warrantyDefault":"...", "footer":"..." },
  "customer": { "name":"...", "phone":"..." },
  "car": { "plate":"...", "company":"...", "model":"...", "year":2020 }
}
MUST strip purchasePrice, vendorId, materialPrice, assignees, and any internal notes from items before returning. Frontend currently consumes this via apiEstimates.bySlug().

GET
/invoices/slug/:slug
Same shape + paid/balance + due/delivery/warranty. NO auth header, sanitized.

Response (200/201)
{
  "number": "INV-5001", "status": "Paid|Partially Paid|Unpaid|Overdue",
  "date":"...", "dueDate":"?", "deliveryDate":"?", "warranty":"?", "finalMileage": 0,
  "workshop": {...}, "customer": {...}, "car": {...},
  "items": [ /* sanitized */ ],
  "subtotal": 7000, "discountAmount": 0, "taxAmount": 0, "grandTotal": 7000,
  "paidAmount": 7000, "balanceDue": 0
}
19. Implementation notes for new endpoints
publicSlug: 8-char URL-safe random (crypto). Unique per document type, indexed.
Public endpoints: rate-limit by IP (e.g. 60 req/min). No cookies.
Internal-field stripping: enforce in a single serializer for public routes — never rely on the frontend.
Dashboard financials: one SQL view, not N+1.
Salary obligation: only sum staff where active = true.
Profit formula: max(0, revenue − salaryObligation − partsPurchaseCost − materialCost).
Parts catalog vs invoice items: catalog is a template; values are copied into the line item — editing the catalog later must NOT mutate historical documents.
RBAC additions: vendors & parts CRUD = owner/manager only; staff salary edit = owner only; dashboard financials = owner/manager/accountant.
Frontend contract version 1 — generated for the Go backend integration. The mock layer lives in src/services/api.ts; swap each method's body with fetch() calls to the endpoints above without changing any UI code.