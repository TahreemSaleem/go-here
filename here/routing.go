package here

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Jeffail/gabs"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dghubble/sling"
	"github.com/ompluscator/dynamic-struct"
)

// RoutingService provides for HERE routing api.
type RoutingService struct {
	sling *sling.Sling
}

// RoutingParams parameters for Routing Service.
type RoutingParams struct {
	Waypoint0 string `url:"waypoint0"`
	Waypoint1 string `url:"waypoint1"`
	APIKey    string `url:"apikey"`
	Modes     string `url:"mode"`
	Departure string `url:"departure"`
}

// RoutingResponse model for routing service.
type RoutingResponse struct {
	Response struct {
		MetaInfo struct {
			Timestamp           time.Time `json:"timestamp"`
			MapVersion          string    `json:"mapVersion"`
			ModuleVersion       string    `json:"moduleVersion"`
			InterfaceVersion    string    `json:"interfaceVersion"`
			AvailableMapVersion []string  `json:"availableMapVersion"`
		} `json:"metaInfo"`
		Route []struct {
			Waypoint []struct {
				LinkID         string `json:"linkId"`
				MappedPosition struct {
					Latitude  float64 `json:"latitude"`
					Longitude float64 `json:"longitude"`
				} `json:"mappedPosition"`
				OriginalPosition struct {
					Latitude  float64 `json:"latitude"`
					Longitude float64 `json:"longitude"`
				} `json:"originalPosition"`
				Type           string  `json:"type"`
				Spot           float64 `json:"spot"`
				SideOfStreet   string  `json:"sideOfStreet"`
				MappedRoadName string  `json:"mappedRoadName"`
				Label          string  `json:"label"`
				ShapeIndex     int     `json:"shapeIndex"`
				Source         string  `json:"source"`
			} `json:"waypoint"`
			Mode struct {
				Type           string        `json:"type"`
				TransportModes []string      `json:"transportModes"`
				TrafficMode    string        `json:"trafficMode"`
				Feature        []interface{} `json:"feature"`
			} `json:"mode"`
			Leg []struct {
				Start struct {
					LinkID         string `json:"linkId"`
					MappedPosition struct {
						Latitude  float64 `json:"latitude"`
						Longitude float64 `json:"longitude"`
					} `json:"mappedPosition"`
					OriginalPosition struct {
						Latitude  float64 `json:"latitude"`
						Longitude float64 `json:"longitude"`
					} `json:"originalPosition"`
					Type           string  `json:"type"`
					Spot           float64 `json:"spot"`
					SideOfStreet   string  `json:"sideOfStreet"`
					MappedRoadName string  `json:"mappedRoadName"`
					Label          string  `json:"label"`
					ShapeIndex     int     `json:"shapeIndex"`
					Source         string  `json:"source"`
				} `json:"start"`
				End struct {
					LinkID         string `json:"linkId"`
					MappedPosition struct {
						Latitude  float64 `json:"latitude"`
						Longitude float64 `json:"longitude"`
					} `json:"mappedPosition"`
					OriginalPosition struct {
						Latitude  float64 `json:"latitude"`
						Longitude float64 `json:"longitude"`
					} `json:"originalPosition"`
					Type           string  `json:"type"`
					Spot           float64 `json:"spot"`
					SideOfStreet   string  `json:"sideOfStreet"`
					MappedRoadName string  `json:"mappedRoadName"`
					Label          string  `json:"label"`
					ShapeIndex     int     `json:"shapeIndex"`
					Source         string  `json:"source"`
				} `json:"end"`
				Length     int `json:"length"`
				TravelTime int `json:"travelTime"`
				Maneuver   []struct {
					Position struct {
						Latitude  float64 `json:"latitude"`
						Longitude float64 `json:"longitude"`
					} `json:"position"`
					Instruction string `json:"instruction"`
					TravelTime  int    `json:"travelTime"`
					Length      int    `json:"length"`
					ID          string `json:"id"`
					Type        string `json:"_type"`
				} `json:"maneuver"`
			} `json:"leg"`
			Summary struct {
				Distance    int      `json:"distance"`
				TrafficTime int      `json:"trafficTime"`
				BaseTime    int      `json:"baseTime"`
				Flags       []string `json:"flags"`
				Text        string   `json:"text"`
				TravelTime  int      `json:"travelTime"`
				Type        string   `json:"_type"`
			} `json:"summary"`
		} `json:"route"`
		Language string `json:"language"`
	} `json:"response"`
}

// newRoutingService returns a new RoutingService.
func newRoutingService(sling *sling.Sling) *RoutingService {
	return &RoutingService{
		sling: sling,
	}
}

// Returns waypoints as a formatted string.
func createWaypoint(waypoint [2]float32) string {
	waypoints := fmt.Sprintf("%f,%f", waypoint[0], waypoint[1])
	return waypoints
}

// CreateRoutingParams creates routing parameters struct.
func (s *RoutingService) CreateRoutingParams(origin [2]float32, destination [2]float32, waypoints [][2]float32, apiKey string, modes []Enum) interface{} {
	stringOrigin := createWaypoint(origin)
	stringDestination := createWaypoint(destination)
	var buffer bytes.Buffer
	for _, routeMode := range modes {
		mode := Enum.ValueOfRouteMode(routeMode)
		buffer.WriteString(mode + ";")
	}
	routeModes := buffer.String()
	routeModes = routeModes[:len(routeModes)-1]

	if len(waypoints) <= 0 {
		routingParams := RoutingParams{
			Waypoint0: stringOrigin,
			Waypoint1: stringDestination,
			APIKey:    apiKey,
			Modes:     routeModes,
			Departure: "now",
		}
		return routingParams
	}

	var builder dynamicstruct.Builder
	var name string
	var index int
	jsonObj := gabs.New()
	extendStruct := dynamicstruct.ExtendStruct(RoutingParams{})
	for index = 0; index <= len(waypoints); index++ {
		name = "Waypoint" + strconv.Itoa(index+1)
		builder = extendStruct.AddField(name, "", `url:"`+strings.ToLower(name)+`"`)
		if index < len(waypoints) {
			jsonObj.Set(createWaypoint(waypoints[index]), name)
		}
	}
	routingParams := builder.Build().New()

	jsonObj.Set(apiKey, "APIKey")
	jsonObj.Set(routeModes, "Modes")
	jsonObj.Set("now", "Departure")
	jsonObj.Set(stringOrigin, "Waypoint0")
	jsonObj.Set(stringDestination, "Waypoint"+strconv.Itoa(index))

	err := json.Unmarshal(jsonObj.Bytes(), &routingParams)

	if err != nil {
		log.Fatal(err)
	}
	return routingParams
}

// Route with given parameters.
func (s *RoutingService) Route(params interface{}) (*RoutingResponse, *http.Response, error) {
	routes := new(RoutingResponse)
	apiError := new(APIError)
	resp, err := s.sling.New().Get("calculateroute.json").QueryStruct(params).Receive(routes, apiError)
	return routes, resp, relevantError(err, *apiError)
}
