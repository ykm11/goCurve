package ecUtils

import (
    "fmt"
    "math/big"
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
