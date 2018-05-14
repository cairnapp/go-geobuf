# Go Geobuf

A compact Protobuf representation of GeoJSON. Based on Mapbox's [geobuf](https://github.com/mapbox).

## Limitations

Due to the nature of Go being a statically typed language, custom properties are not currently supported.

Currently, [orb](https://github.com/paulmach/orb) is the underlying library supporting GeoJSON. 
Long term plans are to offer our own representation with orb as a supported extension.

Some orb properties may lose their types through encoding/decoding. For instance, `int8`s may become `uint`s
or just `int`s.

## Encoding/Decoding

A basic example shows how this library will infer the proper precision for encoding/decoding 

```go
import (
    "github.com/cairnapp/go-geobuf"
    "github.com/paulmach/orb"
    "github.com/paulmach/orb/geojson"
)

point := geojson.NewGeometry(orb.Point([2]float64{
    124.123, 
    234.456
}))
data := geobuf.Encode(point)
decoded_point := geobuf.Decode(point)
```

##
