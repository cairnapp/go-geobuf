# Go Geobuf

A compact Protobuf representation of GeoJSON. Based on Mapbox's [geobuf](https://github.com/mapbox).

## Limitations

Due to the nature of Go being a statically typed language, custom properties are not currently supported.

Some properties may lose their types through encoding/decoding. For instance, `int8`s may become `uint`s
or just `int`s.

## Encoding/Decoding

A basic example shows how this library will infer the proper precision for encoding/decoding 

```go
import (
    "github.com/cairnapp/go-geobuf"
    "github.com/cairnapp/go-geobuf/pkg/geometry"
    "github.com/cairnapp/go-geobuf/pkg/geojson"
)

point := geojson.NewGeometry(geometry.Point([]float64{
    124.123, 
    234.456
}))
data := geobuf.Encode(point)
decoded_point := geobuf.Decode(point)
```

##
