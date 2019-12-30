# Places-API

Microservice implemented in Golang that get nearby places from Google Places API.

## GRPC Service

```go
service PlaceService {
  rpc ListByCoord(PlaceListByCoordRequest) returns (PlaceListByCoordResponse) {}
}

message Place {
  string id = 1;
  string name = 2;
  string rating = 3;
  string address = 4;
  bool open = 5;
  string photo_reference = 6;
  Coord coord = 7;
}

message PlaceListByCoordRequest {
  Coord coord = 1;
  string page_token = 2;
}

message MetaPlaceListByCoord {
  string page_token = 1;
}

message PlaceListByCoordResponse {
  repeated Place data = 1;
  MetaPlaceListByCoord meta = 2;
  Error error = 3;
}
```

## Commands (Development)

`make build`: build places service for osx.

`make linux`: build places service for linux os.

`make docker`: build docker.

`docker run -it -p 5030:5030 places-api`: run docker.

`PORT=<port> API_KEY=<api_key> make r`: run places service.
