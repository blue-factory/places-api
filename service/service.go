package service

import (
	"context"
	"fmt"
	"log"
	"os"

	places "github.com/microapis/places-api"
	"googlemaps.github.io/maps"

	// history "github.com/microapis/places-history-api"
	nats "github.com/nats-io/nats.go"
	// "github.com/kr/pretty"
)

// NewPlaces ...
func NewPlaces(conn *nats.EncodedConn) *Places {
	return &Places{
		Nats: conn,
	}
}

// Places ...
type Places struct {
	Nats *nats.EncodedConn
}

// ListByCoord ...
func (us *Places) ListByCoord(coord places.Coord, userID string) ([]*places.Place, error) {
	// get API_KEY env value
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Println("missing env variable API_KEY, using default value...")
		os.Exit(1)
	}

	// create new google maps client
	client, err := maps.NewClient(maps.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}

	// parse cood values
	loc, err := maps.ParseLatLng(coord.GetLatLngStr())
	if err != nil {
		return nil, err
	}

	// prepare request
	r := &maps.NearbySearchRequest{
		// TODO(ca): get radius value from req param
		Radius:   5000,
		Keyword:  "restaurantes",
		Language: "spanish",
		Location: &loc,
	}

	// fetch request to google maps api
	resp, err := client.NearbySearch(context.Background(), r)

	// prepare result in slice
	rr := make([]*places.Place, 0)

	// parse response
	for _, res := range resp.Results {
		c := places.Coord{
			Latitude:  res.Geometry.Location.Lat,
			Longitude: res.Geometry.Location.Lng,
		}

		place := &places.Place{
			ID:      res.PlaceID,
			Name:    res.Name,
			Address: res.Vicinity,
			Rating:  fmt.Sprintf("%f", res.Rating),
			Coord:   c,
		}

		if res.OpeningHours != nil {
			place.Open = *res.OpeningHours.OpenNow
		}

		if res.Photos != nil {
			place.PhotoReference = res.Photos[0].PhotoReference
		}

		rr = append(rr, place)
	}

	// emit event to history service
	// go func() {
	// 	he := &history.CreateHistoryEvent{
	// 		UserID:    userID,
	// 		Latitude:  coord.GetLatStr(),
	// 		Longitude: coord.GetLngStr(),
	// 	}

	// 	us.Nats.Publish("history.create", he)
	// 	log.Println("Published to History.create service")
	// }()

	return rr, nil
}
