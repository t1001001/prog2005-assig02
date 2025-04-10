package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strings"

	"cloud.google.com/go/firestore"
	c "github.com/t1001001/prog-assig02/internal/constants"
	h "github.com/t1001001/prog-assig02/internal/helpers"
)

// DashboardsHandler handles the GET request for the /dashboards/{id} endpoint
func DashboardsHandler(w http.ResponseWriter, r *http.Request, client *firestore.Client) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := context.Background()
	id := strings.TrimPrefix(r.URL.Path, c.ROOT+c.VERSION+c.DASHBOARDS_PATH)
	id = strings.Trim(id, "/")

	if id == "" {
		http.Error(w, "Missing ID in path", http.StatusBadRequest)
		return
	}

	// Retrieve the corresponding configuration from Firestore
	doc, err := client.Collection("registrations").Doc(id).Get(ctx)
	if err != nil {
		http.Error(w, "Configuration not found", http.StatusNotFound)
		return
	}

	var config c.Country
	if err := doc.DataTo(&config); err != nil {
		http.Error(w, "Failed to decode configuration", http.StatusInternalServerError)
		return
	}

	// Initialize the dashboard response structure
	response := c.DashboardResponse{
		Country:       config.Country,
		ISOCode:       config.ISOCode,
		LastRetrieval: h.GetCurrentTimestamp(),
		Features:      c.EnrichedFeatures{},
	}

	var countryData map[string]interface{}

	// Fetch data from RestCountries API
	if config.Features.Capital || config.Features.Coordinates || config.Features.Population || config.Features.Area || len(config.Features.TargetCurrencies) > 0 {
		restURL := fmt.Sprintf("%s/name/%s", c.RESTCOUNTRIES_API_URL, config.Country)
		if req, err := http.Get(restURL); err == nil && req.StatusCode == 200 {
			defer req.Body.Close()
			body, _ := ioutil.ReadAll(req.Body)
			var result []map[string]interface{}
			if err := json.Unmarshal(body, &result); err == nil && len(result) > 0 {
				countryData = result[0]
				// Extracting features based on the configuration
				// Extracting capital, population, area, and coordinates
				if config.Features.Capital {
					if capitals, ok := countryData["capital"].([]interface{}); ok && len(capitals) > 0 {
						response.Features.Capital = fmt.Sprintf("%v", capitals[0])
					}
				}

				if config.Features.Population {
					if pop, ok := countryData["population"].(float64); ok {
						response.Features.Population = int(pop)
					}
				}

				if config.Features.Area {
					if area, ok := countryData["area"].(float64); ok {
						response.Features.Area = area
					}
				}

				if config.Features.Coordinates {
					if latlng, ok := countryData["latlng"].([]interface{}); ok && len(latlng) >= 2 {
						lat, _ := latlng[0].(float64)
						lng, _ := latlng[1].(float64)
						response.Features.Coordinates = c.Coordinates{
							Latitude:  lat,
							Longitude: lng,
						}
					}
				}

				// Fetch currency rates from Currency API
				if config.Features.TargetCurrencies != nil && len(config.Features.TargetCurrencies) > 0 {
					if currencies, ok := countryData["currencies"].(map[string]interface{}); ok {
						for currencyCode := range currencies {
							currencyURL := fmt.Sprintf("%s/%s", c.CURRENCY_API_URL, strings.ToLower(currencyCode))
							if res, err := http.Get(currencyURL); err == nil && res.StatusCode == 200 {
								defer res.Body.Close()
								body, _ := ioutil.ReadAll(res.Body)
								var parsed map[string]interface{}
								if err := json.Unmarshal(body, &parsed); err == nil {
									if rates, ok := parsed["rates"].(map[string]interface{}); ok {
										targetMap := make(map[string]float64)
										for _, curr := range config.Features.TargetCurrencies {
											if val, exists := rates[curr]; exists {
												if rate, ok := val.(float64); ok {
													targetMap[curr] = rate
												}
											}
										}
										response.Features.TargetCurrencies = targetMap
										break
									}
								}
							}
						}
					}
				}
			}
		}

		// Fetch weather data from OpenMeteo if coordinates and weather features are enabled
		if config.Features.Coordinates && (config.Features.Temperature || config.Features.Precipitation) {
			lat := response.Features.Coordinates.Latitude
			lng := response.Features.Coordinates.Longitude
			weatherURL := fmt.Sprintf("%s/forecast?latitude=%.2f&longitude=%.2f&hourly=apparent_temperature,precipitation&forecast_days=1&timezone=auto", c.OPENMETEO_API_URL, lat, lng)

			if res, err := http.Get(weatherURL); err == nil && res.StatusCode == 200 {
				defer res.Body.Close()
				body, _ := ioutil.ReadAll(res.Body)
				var parsed map[string]interface{}
				if err := json.Unmarshal(body, &parsed); err == nil {
					if hourly, ok := parsed["hourly"].(map[string]interface{}); ok {
						// Extracting temperature and precipitation data
						if config.Features.Temperature {
							if temps, ok := hourly["apparent_temperature"].([]interface{}); ok && len(temps) > 0 {
								sum := 0.0
								count := 0
								for _, t := range temps {
									if v, ok := t.(float64); ok {
										sum += v
										count++
									}
								}
								if count > 0 {
									mean := sum / float64(count)
									response.Features.Temperature = math.Round(mean*100) / 100
								}
							}
						}
						if config.Features.Precipitation {
							if precs, ok := hourly["precipitation"].([]interface{}); ok && len(precs) > 0 {
								sum := 0.0
								count := 0
								for _, p := range precs {
									if v, ok := p.(float64); ok {
										sum += v
										count++
									}
								}
								if count > 0 {
									mean := sum / float64(count)
									response.Features.Precipitation = math.Round(mean*100) / 100
								}
							}
						}
					}
				}
			}
		}

		// Respond with the formatted dashboard data
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		out, err := json.MarshalIndent(response, "", "  ")
		if err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
		w.Write(out)

		// Trigger Webhook
		go TriggerWebhooks(client, config.ISOCode, c.EVENT_REGISTER)
	}
}
