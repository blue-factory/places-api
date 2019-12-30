package places

import (
	"strconv"

	pb "github.com/microapis/places-api/proto"
)

// Place ...
type Place struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Rating         string `json:"rating"`
	Address        string `json:"address"`
	Open           bool   `json:"open"`
	PhotoReference string `json:"photo_reference"`
	Coord          Coord  `json:"coord"`
}

// Service ...
type Service interface {
	ListByCoord(coord Coord, userID string) ([]*Place, error)
}

// ToProto ...
func (r *Place) ToProto() *pb.Place {
	c := &pb.Coord{
		Latitude:  r.Coord.Latitude,
		Longitude: r.Coord.Longitude,
	}

	return &pb.Place{
		Id:             r.ID,
		Name:           r.Name,
		Rating:         r.Rating,
		Address:        r.Address,
		Open:           r.Open,
		PhotoReference: r.PhotoReference,
		Coord:          c,
	}
}

// FromProto ...
func (r *Place) FromProto(rr *pb.Place) *Place {
	c := Coord{
		Latitude:  rr.GetCoord().GetLatitude(),
		Longitude: rr.GetCoord().GetLongitude(),
	}

	r.ID = rr.Id
	r.Name = rr.Name
	r.Rating = rr.Rating
	r.Address = rr.Address
	r.Open = rr.Open
	r.PhotoReference = rr.PhotoReference
	r.Coord = c

	return r
}

// Coord ...
type Coord struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// GetLatLngStr ...
func (c *Coord) GetLatLngStr() string {
	lat := strconv.FormatFloat(c.Latitude, 'f', -1, 64)
	lng := strconv.FormatFloat(c.Longitude, 'f', -1, 64)

	return lat + "," + lng
}

// GetLatStr ...
func (c *Coord) GetLatStr() string {
	return strconv.FormatFloat(c.Longitude, 'f', -1, 64)
}

// GetLngStr ...
func (c *Coord) GetLngStr() string {
	return strconv.FormatFloat(c.Longitude, 'f', -1, 64)
}
