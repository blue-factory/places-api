package placessvc

import (
	"context"
	"fmt"
	"log"

	places "github.com/microapis/places-api"
	pb "github.com/microapis/places-api/proto"
	"github.com/microapis/places-api/service"

	"googlemaps.github.io/maps"

	nats "github.com/nats-io/nats.go"
)

var _ pb.PlaceServiceServer = (*Service)(nil)

// Service ...
type Service struct {
	placesSvc places.Service
}

// New ...
func New(conn *nats.EncodedConn) *Service {
	return &Service{
		placesSvc: service.NewPlaces(conn),
	}
}

// ListByCoord List nearby places by coord.
func (as *Service) ListByCoord(ctx context.Context, gr *pb.PlaceListByCoordRequest) (*pb.PlaceListByCoordResponse, error) {
	lat := gr.GetCoord().GetLatitude()
	lng := gr.GetCoord().GetLongitude()

	c := places.Coord{
		Latitude:  lat,
		Longitude: lng,
	}

	_, err := maps.ParseLatLng(c.GetLatLngStr())
	if err != nil {
		return &pb.PlaceListByCoordResponse{
			Data: nil,
			Error: &pb.PlaceError{
				Code:    500,
				Message: "invalid coord values",
			},
		}, nil
	}

	userID := gr.GetUserId()

	listedPlacess, err := as.placesSvc.ListByCoord(c, userID)
	if err != nil {
		log.Println(fmt.Sprintf("[GRPC][PlacessService][ListByCoord][Error] %v", err))
		return &pb.PlaceListByCoordResponse{
			Data: nil,
			Error: &pb.PlaceError{
				Code:    500,
				Message: err.Error(),
			},
		}, nil
	}

	data := make([]*pb.Place, 0)
	for _, place := range listedPlacess {
		data = append(data, place.ToProto())
	}

	res := &pb.PlaceListByCoordResponse{
		Data:  data,
		Error: nil,
	}

	log.Println(fmt.Sprintf("[GRPC][PlacessService][List][Response] Listed %v places", len(res.Data)))
	return res, nil
}
