package ecUtils

import (
    "math/big"
    . "../mathUtils"
)

func (ec EllipticCurve) Order(p Point, algorithm string) *big.Int {
	var cardinality *big.Int
	switch algorithm {
	case "exhaustive":
		cardinality = exhaostiveSearchOrder(ec, p)
	default:
		cardinality = exhaostiveSearchOrder(ec, p)
	}
	return cardinality
}

func exhaostiveSearchOrder(ec EllipticCurve, p Point) *big.Int {
	if p.IsUnit {
		return ONE
	}
	p2 := ec.PointDoubling(p)
	cardinality := TWO
	for !p2.IsUnit {
		p2 = ec.PointAdd(p2, p)
		cardinality = Add(cardinality, ONE, nil)
	}
	return cardinality
}

