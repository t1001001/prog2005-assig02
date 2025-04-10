# Country Dashboard Service

## Authors
- Tobias Nguyen - tobiasng@stud.ntnu.no
- Rayan Defoor - rayand@stud.ntnu.no

## Service
The Country Dashboard Service is a RESTful API that provides the client the ability to configure information dashboards that are dynamically populated when requested.

- **RestCountries API:** https://restcountries.com
- **Open-Meteo API:** https://open-meteo.com
- **Exchange Rate API:** https://www.exchangerate-api.com

*Please note that these URLs are the actual API endpoints.*  
*This Assignment strictly invokes a self-hosted version of the RestCountries API and Exchange Rate API endpoints.*

## Dependencies
Please install the 1.24.1 version of Go.  
To install it, please run the following commands in your terminal:

```bash
go install golang.org/dl/go1.24.1@latest
go1.24.1 download
```

## How to use the service?
To run the service, simply run the `main.go` file:

```bash
go run cmd/country-dashboard-service/main.go
```

## Available Endpoints
```
/dashboard/v1/registrations/
/dashboard/v1/dashboards/
/dashboard/v1/notifications/
/dashboard/v1/status/
```

---

## Endpoint `/dashboard/v1/registrations/` — Dashboard Configuration

### Register new dashboard configuration

**Method:** POST  
**Path:** `/dashboard/v1/registrations/`  
**Content-Type:** `application/json`

**Request Example:**
```json
{
  "country": "Norway",
  "isoCode": "NO",
  "features": {
    "temperature": true,
    "precipitation": true,
    "capital": true,
    "coordinates": true,
    "population": true,
    "area": true,
    "targetCurrencies": ["EUR", "USD", "SEK"]
  }
}
```

**Response Example:**
```json
{
  "id": "516dba7f015f2a68",
  "lastChange": "20250229 12:31"
}
```

---

### View a specific registered dashboard configuration

**Method:** GET  
**Path:** `/dashboard/v1/registrations/{id}`

**Response Example:**
```json
{
  "id": "516dba7f015f2a68",
  "country": "Norway",
  "isoCode": "NO",
  "features": {
    "temperature": true,
    "precipitation": true,
    "capital": true,
    "coordinates": true,
    "population": true,
    "area": false,
    "targetCurrencies": ["EUR", "USD", "SEK"]
  },
  "lastChange": "20250229 14:07"
}
```

---

### View all registered dashboard configurations

**Method:** GET  
**Path:** `/dashboard/v1/registrations/`

**Response Example:**
```json
[
  {
    "id": "516dba7f015f2a68",
    "country": "Norway",
    "isoCode": "NO",
    "features": {
      "temperature": true,
      "precipitation": true,
      "capital": true,
      "coordinates": true,
      "population": true,
      "area": false,
      "targetCurrencies": ["EUR", "USD", "SEK"]
    },
    "lastChange": "20250229 14:07"
  },
  {
    "id": "bc89adc23e27f42a",
    "country": "Denmark",
    "isoCode": "DK",
    "features": {
      "temperature": false,
      "precipitation": true,
      "capital": true,
      "coordinates": true,
      "population": false,
      "area": true,
      "targetCurrencies": ["NOK", "MYR", "JPY", "EUR"]
    },
    "lastChange": "20250224 08:27"
  }
]
```

> Advanced Task: Implement the HEAD method functionality (only return the header, not the body).

---

### Replace a specific registered dashboard configuration

**Method:** PUT  
**Path:** `/dashboard/v1/registrations/{id}`  
**Content-Type:** `application/json`

**Request Example:**
```json
{
  "country": "Norway",
  "isoCode": "NO",
  "features": {
    "temperature": false,
    "precipitation": true,
    "capital": true,
    "coordinates": true,
    "population": true,
    "area": false,
    "targetCurrencies": ["EUR", "SEK"]
  }
}
```

**Response:** Empty body with appropriate HTTP status code.

---

### Delete a specific registered dashboard configuration

**Method:** DELETE  
**Path:** `/dashboard/v1/registrations/{id}`

**Response:** Empty body with appropriate HTTP status code.

---

## Endpoint `/dashboard/v1/dashboards/{id}` — Retrieve Populated Dashboard

This endpoint can be used to retrieve a populated dashboard using the registered configuration identified by its ID.

**Method:** GET  
**Path:** `/dashboard/v1/dashboards/{id}`  
**Content-Type:** `application/json`

**Response Example:**
```json
{
  "country": "Norway",
  "isoCode": "NO",
  "features": {
    "temperature": -1.2,
    "precipitation": 0.80,
    "capital": "Oslo",
    "coordinates": {
      "latitude": 62.0,
      "longitude": 10.0
    },
    "population": 5379475,
    "area": 323802.0,
    "targetCurrencies": {
      "EUR": 0.087701435,
      "USD": 0.095184741,
      "SEK": 0.97827275
    }
  },
  "lastRetrieval": "20250229 18:15"
}
```

> Note: Only one dashboard can be retrieved at a time to avoid overloading external APIs.

---

## Endpoint `/dashboard/v1/notifications/` — Webhook Management

This endpoint manages webhook registrations triggered by events like configuration changes or dashboard invocations.

### Register a webhook

**Method:** POST  
**Path:** `/dashboard/v1/notifications/`  
**Content-Type:** `application/json`

**Request Example:**
```json
{
  "url": "https://localhost:8080/client/",
  "country": "NO",
  "event": "INVOKE"
}
```

**Response Example:**
```json
{
  "id": "OIdksUDwveiwe"
}
```

---

### Delete a webhook

**Method:** DELETE  
**Path:** `/dashboard/v1/notifications/{id}`

**Response:** Empty body with appropriate status code.

---

### View a specific webhook

**Method:** GET  
**Path:** `/dashboard/v1/notifications/{id}`

**Response Example:**
```json
{
  "id": "OIdksUDwveiwe",
  "url": "https://localhost:8080/client/",
  "country": "NO",
  "event": "INVOKE"
}
```

---

### View all registered webhooks

**Method:** GET  
**Path:** `/dashboard/v1/notifications/`

**Response Example:**
```json
[
  {
    "id": "OIdksUDwveiwe",
    "url": "https://localhost:8080/client/",
    "country": "NO",
    "event": "INVOKE"
  },
  {
    "id": "DiSoisivucios",
    "url": "https://localhost:8081/anotherClient/",
    "country": "",
    "event": "REGISTER"
  }
]
```

---

### Webhook Invocation Payload

When a webhook is triggered, the service sends:

**Method:** POST  
**Path:** `<url specified in the corresponding webhook registration>`  
**Content-Type:** `application/json`

**Body Example:**
```json
{
  "id": "OIdksUDwveiwe",
  "country": "NO",
  "event": "INVOKE",
  "time": "20240223 06:23"
}
```

---

## Endpoint `/dashboard/v1/status/` — Service Monitoring

Returns the availability and status of upstream services, number of registered webhooks, and the uptime.

**Method:** GET  
**Path:** `/dashboard/v1/status/`

**Response Example:**
```json
{
  "countries_api": 200,
  "meteo_api": 200,
  "currency_api": 200,
  "notification_db": 200,
  "webhooks": 4,
  "version": "v1",
  "uptime": 38290
}
```

---


