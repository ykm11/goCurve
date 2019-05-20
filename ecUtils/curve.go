package ecUtils

import (
    "fmt"
    "math/big"
    "../random"
    . "../mathUtils"
)

var (
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

func Point2Str(p Point) string {
    str := fmt.Sprintf("(%d, %d; %d)", p.X, p.Y, Bool2int(!p.IsUnit))
    return str
}

func (ec ellipticCurve) PrintCurve() {
    fmt.Printf("[+] EC: y^2 = x^3 + %dx + %d  OVER Zmod(%d)\n",
        ec.A, ec.B, ec.Modulus)
}

func (ec ellipticCurve) Exist(X, Y *big.Int) bool {
    l := Exp(Y, TWO, ec.Modulus)
    r1 := Exp(X, THREE, ec.Modulus)
    r2 := Mul(ec.A, X, ec.Modulus)
    r := Add(Add(r1, r2, ec.Modulus), ec.B, ec.Modulus)
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

    lmd := Mul(Sub(Q.Y, P.Y, ec.Modulus),
                InvMod(Sub(Q.X, P.X, ec.Modulus), ec.Modulus), ec.Modulus)
    x3 := Sub(Sub(Exp(lmd, TWO, ec.Modulus), P.X, ec.Modulus), Q.X, ec.Modulus)
    y3 := Sub(Mul(lmd, Sub(P.X, x3, ec.Modulus), ec.Modulus), P.Y, ec.Modulus)

    return Point{x3, y3, false}
}

func (ec ellipticCurve) PointDoubling(P Point) Point {
    if P.Y.Cmp(ZERO) == 0 {
        return Origin
    }

    lmd := Mul(Add(Mul(
                    THREE, Exp(P.X, TWO, ec.Modulus), ec.Modulus),
                ec.A, ec.Modulus),
            InvMod(Mul(TWO, P.Y, ec.Modulus), ec.Modulus), ec.Modulus)
    x3 := Sub(Sub(Exp(lmd, TWO, ec.Modulus), P.X, ec.Modulus), P.X, ec.Modulus)
    y3 := Sub(Mul(lmd, Sub(P.X, x3, ec.Modulus), ec.Modulus), P.Y, ec.Modulus)

    return Point{x3, y3, false}
}

func (ec ellipticCurve) Point_xP(x *big.Int, p Point) Point {
    k := Origin
    n := ValCopy(x) // Not to change 'x' (address) value

    for ; n.Cmp(ZERO) != 0 ; {
        if Mod(n, TWO).Cmp(ONE) == 0 {
            k = ec.PointAdd(k, p)
        }
        p = ec.PointDoubling(p)
        n.Div(n, TWO)
    }
    return k
}

func (ec ellipticCurve) VerifySignature(r, s, e, n *big.Int, G, Q Point) bool {
    w := InvMod(s, n)
    u1 := Mul(e, w, n)
    u2 := Mul(r, w, n)
    V := ec.PointAdd(ec.Point_xP(u1, G), ec.Point_xP(u2, Q))
    //fmt.Println("[+] (x2, y2):", Point2Str(V))
    return V.X.Cmp(r) == 0
}

func (ec ellipticCurve) Sign(e, d, n *big.Int, G Point) (*big.Int, *big.Int) {
    k := random.Randint(nil, n)
    R := ec.Point_xP(k, G)
    fmt.Println("[+] (x1, y1) = [k]G:", Point2Str(R))
    r, _ := R.Xy()
    s := Mul(InvMod(k, n), Add(e, Mul(r, d, n), n), n)
    fmt.Printf("[+] Signature (r, s): (%d, %d)\n", r, s)
    return r, s
}
