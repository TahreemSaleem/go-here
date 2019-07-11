package here

import (
	"net/http"

	"github.com/dghubble/sling"
)

// GeocodingService provides for HERE Geocoding api.
type GeocodingService struct {
	sling *sling.Sling
}

// Parameters by search text for Geocoding Service.
type SearchTextParameters struct {
	SearchText string `url:"searchtext"`
}

type GeocodingResponse struct {
	Response struct {
		MetaInfo struct {
			Timestamp string `json:"Timestamp"`
		} `json:"MetaInfo"`
		View []struct {
			Type   string `json:"_type"`
			ViewID int    `json:"ViewId"`
			Result []struct {
				Relevance    int    `json:"Relevance"`
				MatchLevel   string `json:"MatchLevel"`
				MatchQuality struct {
					State       int       `json:"State"`
					City        int       `json:"City"`
					Street      []float64 `json:"Street"`
					HouseNumber int       `json:"HouseNumber"`
				} `json:"MatchQuality"`
				MatchType string `json:"MatchType"`
				Location  struct {
					LocationID      string `json:"LocationId"`
					LocationType    string `json:"LocationType"`
					DisplayPosition struct {
						Latitude  float64 `json:"Latitude"`
						Longitude float64 `json:"Longitude"`
					} `json:"DisplayPosition"`
					NavigationPosition []struct {
						Latitude  float64 `json:"Latitude"`
						Longitude float64 `json:"Longitude"`
					} `json:"NavigationPosition"`
					MapView struct {
						TopLeft struct {
							Latitude  float64 `json:"Latitude"`
							Longitude float64 `json:"Longitude"`
						} `json:"TopLeft"`
						BottomRight struct {
							Latitude  float64 `json:"Latitude"`
							Longitude float64 `json:"Longitude"`
						} `json:"BottomRight"`
					} `json:"MapView"`
					Address struct {
						Label          string `json:"Label"`
						Country        string `json:"Country"`
						State          string `json:"State"`
						County         string `json:"County"`
						City           string `json:"City"`
						District       string `json:"District"`
						Street         string `json:"Street"`
						HouseNumber    string `json:"HouseNumber"`
						PostalCode     string `json:"PostalCode"`
						AdditionalData []struct {
							Value string `json:"value"`
							Key   string `json:"key"`
						} `json:"AdditionalData"`
					} `json:"Address"`
				} `json:"Location"`
			} `json:"Result"`
		} `json:"View"`
	} `json:"Response"`
}

// newGeocodingService returns a new GeocodingService.
func newGeocodingService(sling *sling.Sling) *GeocodingService {
	return &GeocodingService{
		sling: sling,
	}
}

// Geocode by search text.
func (s *GeocodingService) Search(text string) (*GeocodingResponse, *http.Response, error) {
	searchTextParams := &SearchTextParameters{SearchText: text}
	geocodingResponse := new(GeocodingResponse)
	resp, err := s.sling.New().Get("geocode.json").QueryStruct(searchTextParams).ReceiveSuccess(geocodingResponse)
	return geocodingResponse, resp, err
}