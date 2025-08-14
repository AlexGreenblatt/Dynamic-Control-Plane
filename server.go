package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type RouteConfig struct {
	RouteName       string          `json:"routeName"`
	Method          string          `json:"method"`
	RequestSchema   json.RawMessage `json:"requestSchema"`
	ResponseSchema  json.RawMessage `json:"responseSchema"`
	PolicyFileNames []string        `json:"policies"`
}

type RequestBody struct {
	ServiceID any `json:"serviceId"`
	Verbose   any `json:"verbose"`
}

type ResponseSuccess struct {
	Message     string    `json:"message"`
	ServiceID   string    `json:"serviceId"`
	Name        string    `json:"name"`
	Status      string    `json:"status"`
	LastUpdated time.Time `json:"lastUpdated"`
}

type ResponseFailure struct {
	Violations []string `json:"violations"`
	Message    string   `json:"message"`
	Policies   []string `json:"policies"`
}

type PolicyResult struct {
	Violations []string `json:"violations"`
}

// serve is the main entry point for starting the server
func serve() {
	configs := loadRouteConfigs()
	router := gin.Default()
	prepareRoutes(router, configs)
	startServer(router)
}

// loadRouteConfigs loads route configurations from JSON
func loadRouteConfigs() []RouteConfig {
	configs, err := readRouteConfigJSON()
	if err != nil {
		log.Fatalf("Failed to load route config: %v", err)
	}
	return configs
}

// prepareRoutes sets up routes and handlers for the given router
func prepareRoutes(router *gin.Engine, configs []RouteConfig) {
	for _, cfg := range configs {
		registerRoute(router, cfg)
	}
}

// registerRoute registers a single route and its handler
func registerRoute(router *gin.Engine, cfg RouteConfig) {
	cfgCopy := cfg // avoid closure capture

	regoPolicies, err := readRegoPoliciesFromFiles(cfgCopy.PolicyFileNames)
	if err != nil {
		log.Fatalf("Failed to read policies: %v", err)
	}

	switch strings.ToUpper(cfgCopy.Method) {
	case http.MethodGet:
		router.GET(cfgCopy.RouteName, createGetHandler())
	case http.MethodPost:
		router.POST(cfgCopy.RouteName, createPostHandler(cfgCopy.PolicyFileNames, regoPolicies))
	}
}

func createGetHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		serviceID := ctx.Param("serviceId")

		var body map[string]any
		if err := ctx.ShouldBindJSON(&body); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
			return
		}

		ctx.JSON(http.StatusOK, ResponseSuccess{
			Message:     "OK",
			ServiceID:   serviceID,
			Name:        "Some name",
			Status:      "Active",
			LastUpdated: time.Now().UTC(),
		})
	}
}

func createPostHandler(policyFileNames []string, policies []string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		serviceID := ctx.Param("serviceId")

		var body map[string]any
		if err := ctx.ShouldBindJSON(&body); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
			return
		}

		violations := evaluateAllPolicies(body, policies)
		if len(violations) > 0 {
			ctx.JSON(http.StatusBadRequest, ResponseFailure{
				Violations: violations,
				Message:    "Validation failed",
				Policies:   policyFileNames,
			})
			return
		}

		ctx.JSON(http.StatusOK, ResponseSuccess{
			Message:     "OK",
			ServiceID:   serviceID,
			Name:        "Some name",
			Status:      "Active",
			LastUpdated: time.Now().UTC(),
		})
	}
}

// evaluateAllPolicies evaluates multiple Rego policies and returns all violations
// using map[string]any here instead of a struct to allow Rego to do the validation
func evaluateAllPolicies(body map[string]any, policies []string) []string {
	allViolations := []string{}
	for _, policy := range policies {
		violations, err := evaluatePolicy(body, policy)
		if err != nil {
			log.Printf("Policy evaluation error: %v", err)
			continue
		}
		allViolations = append(allViolations, violations...)
	}
	return allViolations
}

func startServer(router *gin.Engine) {
	router.Run()
}
