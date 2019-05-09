package ecUtils

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

func add(a, b, modulus *big.Int) *big.Int {
    r := new(big.Int).Add(a, b)
    if modulus != nil {
        r.Mod(r, modulus)
    }
    return r
}
func mul(a, b, modulus *big.Int) *big.Int {
    r := new(big.Int).Mul(a, b)
    if modulus != nil {
        r.Mod(r, modulus)
    }
    return r
}
func sub(a, b, modulus *big.Int) *big.Int {
    r := new(big.Int).Sub(a, b)
    if modulus != nil {
        r.Mod(r, modulus)
    }
    return r
}
func exp(a, b, modulus *big.Int) *big.Int {
    r := new(big.Int).Exp(a, b, modulus)
    return r
}
func invmod(a, n *big.Int) *big.Int {
    r := new(big.Int).ModInverse(a, n)
    return r
}
func div(a, b *big.Int) *big.Int { // a % b == 0
    r := new(big.Int).Div(a, b)
    return r
}
func mod(x, modulus *big.Int) *big.Int {
    return new(big.Int).Mod(x, modulus)
}
func ceil_sqrt(x *big.Int) *big.Int {
    r := new(big.Int).Sqrt(sub(x, ONE, nil))
    return add(r, ONE, nil)
}

func Extgcd(x, y *big.Int) (*big.Int, *big.Int, *big.Int) {
    s0, s1 := ONE, ZERO
    t0, t1 := ZERO, ONE

    m, n := ValCopy(x), ValCopy(y)
    for ; n.Cmp(ZERO) != 0; {
        //fmt.Printf("m, n : %d, %d\n", m, n)
        r := mod(m, n)
        q := div(sub(m, r, nil), n)
        s1, s0 = sub(s0, mul(q, s1, nil), nil), s1
        t1, t0 = sub(t0, mul(q, t1, nil), nil), t1

        m, n = n, mod(m, n)
    }

    if add(mul(x, s0, nil), mul(y, t0, nil), nil).Cmp(m) != 0 {
        panic("ExtGCD panic..")
    }
    return m, s0, t0
}

func CRT(remainders, modulus []*big.Int) *big.Int {
    //fmt.Printf("[+] Modulus %d\n", modulus)
    //fmt.Printf("[+] Remainders %d\n", remainders)

    N := ONE
    for i := 0; i < len(modulus); i++ {
        N = mul(N, modulus[i], nil)
    }
    x := ZERO
    for i := 0; i < len(modulus); i++ {
        g, r, s := Extgcd(modulus[i], div(N, modulus[i]))
        if g.Cmp(ONE) != 0 {
            fmt.Printf("GCD(%d, %d) should be 1\n", r, s)
            panic("")
        }
        x = add(x,
            mul(mul(remainders[i], s, nil), div(N, modulus[i]), nil), nil)
    }
    return mod(x, N)
}

func Legendre(x, p *big.Int) bool { // is x a quadratic-residue?
    r := exp(x, div(sub(p, ONE, p), TWO), p)
    return r.Cmp(ONE) == 0
}

func GetQuadraticResidueRoot(x, modulus *big.Int) *big.Int {
    if !Legendre(x, modulus) {
        return nil
    }
    z := ONE
    for ; exp(z, div(sub(modulus, ONE, nil), TWO), modulus).Cmp(sub(modulus, ONE, nil)) != 0; {
        z = add(z, ONE, nil)
    }
    q := sub(modulus, ONE, nil)

    m := ZERO
    for ; mod(q, TWO).Cmp(ZERO) == 0; {
        q = div(q, TWO)
        m = add(m, ONE, nil)
    }

    c := exp(z, q, modulus)
    t := exp(x, q, modulus)
    r := exp(x, div(add(q, ONE, nil), TWO), modulus)
    for i := ValCopy(m); i.Cmp(ONE) != 0; i = sub(i, ONE, nil) {
        tmp := exp(t, exp(TWO, sub(i, TWO, nil), nil), modulus)
        if tmp.Cmp(ONE) != 0 {
            r = mul(r, c, modulus)
            t = mul(t, exp(c, TWO, modulus), modulus)
        }
        c = exp(c, TWO, modulus)
    }
    return r
}
