package ecUtils

import (
    "math/big"
    . "../mathUtils"
)

func BsGs(ec EllipticCurve, P, Q Point, cardinality *big.Int) *big.Int {
	// Find d; Q = [d]P, otherwise nil

	//cardinality := ec.Order(P, "exhaustive")
	m := Ceil_sqrt(cardinality)

	// Baby Step
	b := Q
	minus_P := P
	minus_P.Y = Sub(ZERO, minus_P.Y, ec.Modulus)
	baby := []Point{}
	for j := ZERO; j.Cmp(m) == -1; j = Add(j, ONE, nil) { // [0, m)
		baby = append(baby, b)
		b = ec.PointAdd(b, minus_P)
	}

	// Giant Step
	mP := ec.Point_xP(m, P)
	Giant := mP
	for i := ONE; i.Cmp(m) == -1; i = Add(i, ONE, nil) { // [1, m)
		index := ExistInPointArray(Giant, baby)
		if index != nil {
			d := Add(Mul(i, m, nil), index, nil)
			return Mod(d, cardinality)
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
		j = Add(j, ONE, nil)
	}
	return nil
}

func Pollard_rho_f(ec EllipticCurve, alpha, beta, x Point, a, b, order *big.Int) (Point, *big.Int, *big.Int){
    if Mod(x.X, THREE).Cmp(ZERO) == 0 {
        return ec.PointAdd(beta, x), a, Add(b, ONE, order)
    } else if Mod(x.X, THREE).Cmp(ONE) == 0 {
        return ec.PointDoubling(x), Mul(a, TWO, order), Mul(b, TWO, order)
    } else {
        return ec.PointAdd(alpha, x), Add(a, ONE, order), b
    }
}

// Solver doesn's work when "order" is not a prime num. You need to factorize.
func Pollard_rho_ECDLP(P, Q Point, ec EllipticCurve, order *big.Int) *big.Int {
    // Q = [d] * P; d < order
    a, b, x := ZERO, ZERO, Origin
    A, B, X := ValCopy(a), ValCopy(b), x

    for i := ZERO ; i.Cmp(order) == -1; i = Add(i, ONE, nil) {
        x, a, b = Pollard_rho_f(ec, P, Q, x, a, b, order)
        X, A, B = Pollard_rho_f(ec, P, Q, X, A, B, order)
        X, A, B = Pollard_rho_f(ec, P, Q, X, A, B, order)
        if cmpPoint(x, X) &&
        (A.Cmp(ZERO) != 0 && a.Cmp(ZERO) != 0 && B.Cmp(ZERO) != 0 && b.Cmp(ZERO) != 0) {
            r := Sub(B, b, order)
            if r.Cmp(ZERO) != 0 {
                return Mul(InvMod(r, order), Sub(a, A, order), order)
            }
        }
    }
    return nil
}

func Solve_ECDLP(ec EllipticCurve, P, Q Point, order *big.Int) *big.Int {
    if P.IsUnit && Q.IsUnit { // Already Q = P = O, so Q = [0]P
        return ZERO
    }

    if order.Cmp(big.NewInt(100)) == -1 {
        tmp_P := P
        for i := ONE; i.Cmp(big.NewInt(100)) == -1; i = Add(i, ONE, nil) {
            if cmpPoint(tmp_P, Q) {
                return i
            }
            tmp_P = ec.PointAdd(tmp_P, P)
        }
    } else if order.Cmp(big.NewInt(10000)) == -1 {
        return BsGs(ec, P, Q, order)
    }
    return Pollard_rho_ECDLP(P, Q, ec, order)
}
