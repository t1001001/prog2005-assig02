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

`go install golang.org/dl/go1.24.1@latest`

`go1.24.1 download`

## How to use the service?
To run the service, simply run the `main.go` file

`go run cmd/country-dashboard-service/main.go`

# Endpoints
The service provides four endpoints:
```
/dashboard/v1/registrations/
/dashboard/v1/dashboards/
/dashboard/v1/notifications/
/dashboard/v1/status/
```