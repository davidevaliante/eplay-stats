package handlers

import (
	"encoding/json"
	"eplay-reports/env"
	"eplay-reports/models"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// Define the external API response format
type APIResult struct {
	ID                 string      `json:"id"`
	Date               string      `json:"date"`
	Campaign           string      `json:"campaign"`
	SignUps            int         `json:"signups"`
	Ftd                *int        `json:"ftd"`
	Cpa                *int        `json:"cpa"`
	Deposits           interface{} `json:"deposits"` // Can be string or float64
	Ggr                interface{} `json:"ggr"`      // Can be string or float64
	Bet                interface{} `json:"bet"`      // Can be string or float64
	Win                interface{} `json:"win"`      // Can be string or float64
	Bonus              interface{} `json:"bonus"`    // Can be string or float64
	Depo               interface{} `json:"depo"`     // Can be string or float64
	Withd              interface{} `json:"withd"`    // Can be string or float64
	Netrev             interface{} `json:"netrev"`   // Can be string or float64
	RevShareCommission float64     `json:"revShareCommission"`
}

func RootGet(c *gin.Context) {
	partnerID := env.Env.PartnerId
	subscriptionKey := env.Env.SubriptionKey
	serviceToken := env.Env.ServiceToken

	startDateStr := c.Query("startdate")
	endDateStr := c.Query("enddate")

	log.Printf("Received request with startdate: %s, enddate: %s", startDateStr, endDateStr)

	if startDateStr == "" || endDateStr == "" {
		log.Println("Error: Missing startdate or enddate in query parameters")
		c.JSON(http.StatusBadRequest, gin.H{"error": "startdate and enddate query parameters are required"})
		return
	}

	const layout = "2006-01-02"
	startDate, err := time.Parse(layout, startDateStr)
	if err != nil {
		log.Printf("Error parsing startdate: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid format for startdate, expected YYYY-MM-DD"})
		return
	}

	endDate, err := time.Parse(layout, endDateStr)
	if err != nil {
		log.Printf("Error parsing enddate: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid format for enddate, expected YYYY-MM-DD"})
		return
	}

	if startDate.After(endDate) {
		log.Println("Error: startdate cannot be after enddate")
		c.JSON(http.StatusBadRequest, gin.H{"error": "startdate cannot be after enddate"})
		return
	}

	apiURL := fmt.Sprintf(
		"https://api-eplay24.azure-api.net/bi/stats_gb/GetBySkin?id_conto_partner=%s&token=%s&data_inizio=%s&data_fine=%s",
		partnerID, serviceToken, startDateStr, endDateStr,
	)

	log.Printf("Making GET request to external API: %s", apiURL)

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	req.Header.Set("Ocp-Apim-Subscription-Key", subscriptionKey)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error making external API request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to make external API request"})
		return
	}
	defer resp.Body.Close()

	log.Printf("Received response with status code: %d", resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response body"})
		return
	}

	// Define the structure of the response from the API
	type APIResponse struct {
		Results []APIResult `json:"results"`
	}

	// Unmarshal the body into the APIResponse structure
	var apiResponse APIResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		log.Printf("Error unmarshalling API response: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse external API response"})
		return
	}

	// Transform the APIResult into Campaign struct
	var campaigns []models.Campaign
	for _, result := range apiResponse.Results {
		campaign := models.Campaign{
			Date:               result.Date,
			Campaign:           result.Campaign,
			SignUps:            result.SignUps,
			Ftd:                getOrDefault(result.Ftd, 0),
			Cpa:                getOrDefault(result.Cpa, 0),
			Deposits:           parseToCents(result.Deposits),
			Clicks:             0,                         // Assuming Clicks is 0 as it's not provided in the API response
			CpaCommission:      0.0,                       // Assuming CpaCommission is 0 as it's not provided in the API response
			RevShareCommission: result.RevShareCommission, // Assuming RevShareCommission is 0 as it's not provided in the API response
			TotalCommission:    0.0,                       // Assuming TotalCommission is 0 as it's not provided in the API response
			Ggr:                parseToCents(result.Ggr),
			Bet:                parseToCents(result.Bet),
			Win:                parseToCents(result.Win),
			Bonus:              parseToCents(result.Bonus),
			Depo:               parseToCents(result.Depo),
			Withd:              parseToCents(result.Withd),
			Netrev:             parseToCents(result.Netrev),
		}
		campaigns = append(campaigns, campaign)
	}

	// Return the transformed data as JSON with the "data" field
	c.JSON(http.StatusOK, gin.H{"data": campaigns})
}

// Helper function to get a default value for an int pointer
func getOrDefault(val *int, defaultVal int) int {
	if val != nil {
		return *val
	}
	return defaultVal
}

// Helper function to parse an interface (string or float64) into cents
func parseToCents(value interface{}) float64 {
	var floatValue float64
	var err error

	switch v := value.(type) {
	case string:
		// Try parsing the string to a float
		floatValue, err = strconv.ParseFloat(v, 64)
		if err != nil {
			log.Printf("Error parsing string to float: %v", err)
			return 0.0
		}
	case float64:
		floatValue = v
	default:
		return 0.0
	}

	// Multiply by 100 to convert to cents
	return floatValue
}
