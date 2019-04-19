package ecUtils

import (
    "math/big"
)

func add(a, b, modulus *big.Int) *big.Int {
    r := new(big.Int).Add(a, b)
    r.Mod(r, modulus)
    return r
}
func mul(a, b, modulus *big.Int) *big.Int {
    r := new(big.Int).Mul(a, b)
    r.Mod(r, modulus)
    return r
}
func sub(a, b, modulus *big.Int) *big.Int {
    r := new(big.Int).Sub(a, b)
    r.Mod(r, modulus)
    return r
}
func expMod(a, b, modulus *big.Int) *big.Int {
    r := new(big.Int).Exp(a, b, modulus)
    return r
}
func invmod(a, n *big.Int) *big.Int {
    r := new(big.Int).ModInverse(a, n)
    return r
}
