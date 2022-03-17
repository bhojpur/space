# Bhojpur Space - Client Library

The `client-side` library provide Maps API wrappers.

## Modules

- [X] Geocoding
- [X] Directions
- [X] Directions Matrix
- [X] Map Matching
- [ ] Styles
- [X] Maps
- [ ] Static
- [ ] Datasets

## Examples

### Initialisation

```go
// Import the core module (and any required APIs)
import (
    engine "github.com/bhojpur/space/pkg/client"
    "github.com/bhojpur/space/pkg/client/base"
)

// Fetch token from somewhere
token := os.Getenv("BHOJPUR_SPACE_MAPS_TOKEN")

// Create new MapEngine instance
engine := engine.NewMapEngine(token)

```

## Maps API

``` go
import (
    "github.com/bhojpur/space/pkg/client/maps"
)

img, err := engine.Maps.GetTiles(maps.MapIDSatellite, 1, 0, 0, maps.MapFormatJpg90, true)
```

### Geocoding

```go
import (
    "github.com/bhojpur/space/pkg/client/geocode"
)

// Forward Geocoding
var forwardOpts geocode.ForwardRequestOpts
forwardOpts.Limit = 1

place := "2 lincoln memorial circle nw"

forward, err := engine.Geocode.Forward(place, &forwardOpts)

// Reverse Geocoding
var reverseOpts geocode.ReverseRequestOpts
reverseOpts.Limit = 1

loc := &base.Location{72.438939, 34.074122}

reverse, err := engine.Geocode.Reverse(loc, &reverseOpts)
```

### Directions

```go
import (
    "github.com/bhojpur/space/pkg/client/directions"
)

var directionOpts directions.RequestOpts

locs := []base.Location{{-122.42, 37.78}, {-77.03, 38.91}}

directions, err := engine.Directions.GetDirections(locs, directions.RoutingCycling, &directionOpts)
```

## Layout

- [base](base/) contains a common base for Bhojpur Space API modules
- [maps](maps/) contains the Maps API module
- [directions](directions/) contains the Directions API module
- [geocode](geocode/) contains the Geocoding API module
