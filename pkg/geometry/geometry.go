package geometry

type Geometry interface {
	private()
}

type Point []float64

func (p Point) private() {}

func (p Point) Equal(other Point) bool {
	if len(p) != len(other) {
		return false
	}

	for i, coord := range p {
		if other[i] != coord {
			return false
		}
	}

	return true
}

type MultiPoint []Point

func (m MultiPoint) private() {}

func (m MultiPoint) Equal(other MultiPoint) bool {
	if len(m) != len(other) {
		return false
	}

	for i, p := range m {
		if !other[i].Equal(p) {
			return false
		}
	}

	return true
}

type LineString []Point

func (m LineString) private() {}

func (ls LineString) Equal(lineString LineString) bool {
	return MultiPoint(ls).Equal(MultiPoint(lineString))
}

type MultiLineString []LineString

func (m MultiLineString) private() {}

func (m MultiLineString) Equal(other MultiLineString) bool {
	if len(m) != len(other) {
		return false
	}

	for i, p := range m {
		if !other[i].Equal(p) {
			return false
		}
	}

	return true
}

type Ring LineString

func (r Ring) private() {}

func (r Ring) Equal(ring Ring) bool {
	return MultiPoint(r).Equal(MultiPoint(ring))
}

type Polygon []Ring

func (p Polygon) private() {}

func (p Polygon) Equal(polygon Polygon) bool {
	if len(p) != len(polygon) {
		return false
	}

	for i := range p {
		if !p[i].Equal(polygon[i]) {
			return false
		}
	}

	return true
}

type MultiPolygon []Polygon

func (mp MultiPolygon) private() {}

func (mp MultiPolygon) Equal(multiPolygon MultiPolygon) bool {
	if len(mp) != len(multiPolygon) {
		return false
	}

	for i, p := range mp {
		if !p.Equal(multiPolygon[i]) {
			return false
		}
	}

	return true
}

type Collection []Geometry

func (c Collection) private() {}
