package ecUtils

import (
    "math/big"
)

func (ec ellipticCurve) Order(p Point, algorithm string) *big.Int {
    var cardinality *big.Int
    switch algorithm {
    case "exhaustive":
        cardinality = exhaostiveSearchOrder(ec, p)
    default:
        cardinality = exhaostiveSearchOrder(ec, p)
    }
    return cardinality
}

func exhaostiveSearchOrder(ec ellipticCurve, p Point) *big.Int {
    if p.IsUnit {
        return ONE
    }
    p2 := ec.PointDoubling(p)
    cardinality := TWO
    for ; !p2.IsUnit ; {
        p2 = ec.PointAdd(p2, p)
        cardinality = add(cardinality, ONE, nil)
    }
    return cardinality
}


func BsGs(ec ellipticCurve, P, Q Point) *big.Int {
    // Find d; Q = [d]P, otherwise nil
    m := ceil_sqrt(ec.Order(P, "exhaustive"))

    // Baby Step
    b := Q
    minus_P := P
    minus_P.Y = sub(ZERO, minus_P.Y, ec.Modulus)
    baby := []Point{}
    for j := ZERO; j.Cmp(m) == -1; j = add(j, ONE, nil) { // [0, m)
        baby = append(baby, b)
        b = ec.PointAdd(b, minus_P)
    }

    // Giant Step
    mP := ec.Point_xP(m, P)
    Giant := mP
    for i := ONE; i.Cmp(m) == -1; i = add(i, ONE, nil) { // [1, m)
        index := ExistInPointArray(Giant, baby)
        if index != nil {
            return add(mul(i, m, nil), index, nil)
        } else {
            Giant = ec.PointAdd(Giant, mP)
        }
    }

    return nil
}

func ExistInPointArray(P Point, PointArray []Point) *big.Int {
    // is P in Array? if so, return index, otherwise nil
    j := ZERO
    for _, Q := range PointArray { // [0, array_length)
        if cmpPoint(P, Q) {
            return j
        }
        j = add(j, ONE, nil)
    }
    return nil
}
