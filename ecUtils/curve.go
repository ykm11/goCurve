package ecUtils

import (
    "fmt"
    "math/big"
    "../random"
)

var (
    ZERO = big.NewInt(0)
    ONE = big.NewInt(1)
    TWO = big.NewInt(2)
    THREE = big.NewInt(3)

    Origin = Point{ZERO, ONE, true}
)

type Point struct {
    X *big.Int
    Y *big.Int
    IsUnit bool
}

type ellipticCurve struct {
    A *big.Int
    B *big.Int
    Modulus *big.Int
}

func NewCurve(a, b, modulus *big.Int) ellipticCurve {
    return ellipticCurve{a, b, modulus}
}

func (ec ellipticCurve) Point(x, y *big.Int) Point {
    if !ec.Exist(x, y) {
        panic("(x, y) doesn't exist")
    }
    return Point{x, y, false}
}

func (p Point) Xy() (*big.Int, *big.Int){
    return p.X, p.Y
}

func (ec ellipticCurve) PrintCurve() {
    fmt.Printf("[+] EC: y^2 = x^3 + %dx + %d  OVER Zmod(%d)\n",
        ec.A, ec.B, ec.Modulus)
}

func (ec ellipticCurve) Exist(X, Y *big.Int) bool {
    l := exp(Y, TWO, ec.Modulus)
    r1 := exp(X, THREE, ec.Modulus)
    r2 := mul(ec.A, X, ec.Modulus)
    r := add(add(r1, r2, ec.Modulus), ec.B, ec.Modulus)
    return l.Cmp(r) == 0
}

func cmpPoint(P, Q Point) bool {
    return (P.X.Cmp(Q.X) == 0) && (P.Y.Cmp(Q.Y) == 0)
}

func (ec ellipticCurve) PointAdd(P, Q Point) Point {
    if cmpPoint(P, Q) {
        return ec.PointDoubling(P)
    }
    if P.X.Cmp(Q.X) == 0 {
        return Origin
    }
    if P.IsUnit {
        return Q
    }
    if Q.IsUnit {
        return P
    }

    lmd := mul(sub(Q.Y, P.Y, ec.Modulus),
                invmod(sub(Q.X, P.X, ec.Modulus), ec.Modulus), ec.Modulus)
    x3 := sub(sub(exp(lmd, TWO, ec.Modulus), P.X, ec.Modulus), Q.X, ec.Modulus)
    y3 := sub(mul(lmd, sub(P.X, x3, ec.Modulus), ec.Modulus), P.Y, ec.Modulus)

    return Point{x3, y3, false}
}

func (ec ellipticCurve) PointDoubling(P Point) Point {
    if P.Y.Cmp(ZERO) == 0 {
        return Origin
    }

    lmd := mul(add(mul(
                    THREE, exp(P.X, TWO, ec.Modulus), ec.Modulus),
                ec.A, ec.Modulus),
            invmod(mul(TWO, P.Y, ec.Modulus), ec.Modulus), ec.Modulus)
    x3 := sub(sub(exp(lmd, TWO, ec.Modulus), P.X, ec.Modulus), P.X, ec.Modulus)
    y3 := sub(mul(lmd, sub(P.X, x3, ec.Modulus), ec.Modulus), P.Y, ec.Modulus)

    return Point{x3, y3, false}
}

func (ec ellipticCurve) Point_xP(x *big.Int, p Point) Point {
    k := Origin
    n := ValCopy(x) // Not to change 'x' (address) value

    /*
    def mul(a, n)
        k = Origin
        while n != 0
            if n % 2 == 1
                k = add(k, a)
            end
            a = add(a, a)
            n /= 2
        end
    */
    for ; n.Cmp(ZERO) != 0 ; {
        if new(big.Int).Mod(n, TWO).Cmp(ONE) == 0 {
            k = ec.PointAdd(k, p)
        }
        p = ec.PointDoubling(p)
        n.Div(n, TWO)
    }
    return k
}

func (ec ellipticCurve) VerifySignature(r, s, e, n *big.Int, G, Q Point) bool {
    w := invmod(s, n)
    u1 := mul(e, w, n)
    u2 := mul(r, w, n)
    V := ec.PointAdd(ec.Point_xP(u1, G), ec.Point_xP(u2, Q))
    //fmt.Println("[+] (x2, y2):", Point2Str(V))
    return V.X.Cmp(r) == 0
}

func (ec ellipticCurve) Sign(e, d, n *big.Int, G Point) (*big.Int, *big.Int) {
    k := random.Randint(n)
    R := ec.Point_xP(k, G)
    fmt.Println("[+] (x1, y1) = [k]G:", Point2Str(R))
    r, _ := R.Xy()
    s := mul(invmod(k, n), add(e, mul(r, d, n), n), n)
    fmt.Printf("[+] Signature (r, s): (%d, %d)\n", r, s)
    return r, s
}

func (ec ellipticCurve) Order(p Point, algorithm string) *big.Int {
    var cardinality *big.Int
    switch algorithm {
    case "exhausive":
        cardinality = exhaosiveSearchOrder(ec, p)
    default:
        cardinality = exhaosiveSearchOrder(ec, p)
    }
    return cardinality
}

func exhaosiveSearchOrder(ec ellipticCurve, p Point) *big.Int {
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
