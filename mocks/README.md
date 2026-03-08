# DETRAN Mock API

The **DETRAN Mock API** simulates a vehicle registry service used during development and testing.

It provides vehicle metadata such as:

• registration status  
• outstanding fines  
• model code  
• estimated market price  

This service allows the **CRE workflow** to fetch vehicle data without relying on external government APIs.

---

# Purpose

During development we cannot rely on real registry systems.

This mock API simulates the behavior of:

• Brazilian **DETRAN vehicle registry**  
• **FIPE vehicle price index**

The CRE workflow queries this service to retrieve vehicle information before minting the NFT.

---

# API Endpoint

The service exposes the following endpoint:

```
GET /detran/{plate}
```

Example request:

```
GET http://localhost:8080/detran/ABC1234
```

---

# Example Response

```json
{
  "plate": "ABC1234",
  "status": "clear",
  "fines": 0,
  "model_code": "005456-9",
  "price": 35000
}
```

Fields:

| Field | Description |
|------|-------------|
| plate | Vehicle license plate |
| status | Registry status (`clear`, `blocked`, `stolen`) |
| fines | Outstanding fines |
| model_code | Vehicle model identifier |
| price | Estimated vehicle value |

---

# Example Mock Rules

The API includes simple mock rules for demonstration:

| Plate | Price |
|------|------|
| ABC1234 | 35000 |
| XYZ9999 | 120000 |
| any other | 75000 |

---

# Running the Mock API

Start the server:

```bash
go run main.go
```

The API will run on:

```
http://localhost:8080
```

---

# Production Implementation

In a real production environment, this service should be replaced with integrations to official data providers.

Examples:

• **DETRAN vehicle registry APIs**  
• **FIPE price index services**  
• insurance or financing databases  

These data sources would provide **verified vehicle ownership and market valuation** used by the oracle workflow.

---

# Summary

The DETRAN Mock API is a lightweight development tool that simulates external vehicle data services.

It enables the CRE workflow to validate vehicle information during testing without relying on real-world government or financial systems.