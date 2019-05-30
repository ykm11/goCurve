package mathUtils

import (
    "fmt"
    "math/big"
)

var (
    ZERO = big.NewInt(0)
    ONE = big.NewInt(1)
    TWO = big.NewInt(2)
    THREE = big.NewInt(3)
)

func Add(a, b, modulus *big.Int) *big.Int {
    r := new(big.Int).Add(a, b)
    if modulus != nil {
        r.Mod(r, modulus)
    }
    return r
}
func Mul(a, b, modulus *big.Int) *big.Int {
    r := new(big.Int).Mul(a, b)
    if modulus != nil {
        r.Mod(r, modulus)
    }
    return r
}
func Sub(a, b, modulus *big.Int) *big.Int {
    r := new(big.Int).Sub(a, b)
    if modulus != nil {
        r.Mod(r, modulus)
    }
    return r
}
func Exp(a, b, modulus *big.Int) *big.Int {
    r := new(big.Int).Exp(a, b, modulus)
    return r
}
func InvMod(a, n *big.Int) *big.Int {
    r := new(big.Int).ModInverse(a, n)
    return r
}
func Div(a, b *big.Int) *big.Int { // a % b == 0
    r := new(big.Int).Div(a, b)
    return r
}
func Mod(x, modulus *big.Int) *big.Int {
    return new(big.Int).Mod(x, modulus)
}
func GCD(x, y *big.Int) *big.Int {
    return new(big.Int).GCD(nil, nil, x, y)
}
func Ceil_sqrt(x *big.Int) *big.Int {
    r := new(big.Int).Sqrt(Sub(x, ONE, nil))
    return Add(r, ONE, nil)
}

func Extgcd(x, y *big.Int) (*big.Int, *big.Int, *big.Int) {
    s0, s1 := ONE, ZERO
    t0, t1 := ZERO, ONE

    m, n := ValCopy(x), ValCopy(y)
    for ; n.Cmp(ZERO) != 0; {
        //fmt.Printf("m, n : %d, %d\n", m, n)
        r := Mod(m, n)
        q := Div(Sub(m, r, nil), n)
        s1, s0 = Sub(s0, Mul(q, s1, nil), nil), s1
        t1, t0 = Sub(t0, Mul(q, t1, nil), nil), t1

        m, n = n, Mod(m, n)
    }

    if Add(Mul(x, s0, nil), Mul(y, t0, nil), nil).Cmp(m) != 0 {
        panic("ExtGCD panic..")
    }
    return m, s0, t0
}

func CRT(remainders, modulus []*big.Int) *big.Int {
    //fmt.Printf("[+] Modulus %d\n", modulus)
    //fmt.Printf("[+] Remainders %d\n", remainders)

    N := ONE
    for i := 0; i < len(modulus); i++ {
        N = Mul(N, modulus[i], nil)
    }
    x := ZERO
    for i := 0; i < len(modulus); i++ {
        g, r, s := Extgcd(modulus[i], Div(N, modulus[i]))
        if g.Cmp(ONE) != 0 {
            fmt.Printf("GCD(%d, %d) should be 1\n", r, s)
            panic("")
        }
        x = Add(x,
            Mul(Mul(remainders[i], s, nil), Div(N, modulus[i]), nil), nil)
    }
    return Mod(x, N)
}

func Legendre(x, p *big.Int) bool { // is x a quadratic-residue?
    r := Exp(x, Div(Sub(p, ONE, p), TWO), p)
    return r.Cmp(ONE) == 0
}

func GetQuadraticResidueRoot(x, modulus *big.Int) *big.Int {
    if !Legendre(x, modulus) {
        return nil
    }
    z := ONE
    for ; Exp(z, Div(Sub(modulus, ONE, nil), TWO), modulus).Cmp(Sub(modulus, ONE, nil)) != 0; {
        z = Add(z, ONE, nil)
    }
    q := Sub(modulus, ONE, nil)

    m := ZERO
    for ; Mod(q, TWO).Cmp(ZERO) == 0; {
        q = Div(q, TWO)
        m = Add(m, ONE, nil)
    }

    c := Exp(z, q, modulus)
    t := Exp(x, q, modulus)
    r := Exp(x, Div(Add(q, ONE, nil), TWO), modulus)
    for i := ValCopy(m); i.Cmp(ONE) != 0; i = Sub(i, ONE, nil) {
        tmp := Exp(t, Exp(TWO, Sub(i, TWO, nil), nil), modulus)
        if tmp.Cmp(ONE) != 0 {
            r = Mul(r, c, modulus)
            t = Mul(t, Exp(c, TWO, modulus), modulus)
        }
        c = Exp(c, TWO, modulus)
    }
    return r
}
