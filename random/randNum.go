package random

import (
    "math"
    "math/big"
    . "../mathUtils"
    "crypto/rand"
    mrand "math/rand"
)

func Setup() error {
    seed, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
    mrand.Seed(seed.Int64())
    return err
}

func Randint(offset, n *big.Int) *big.Int {
    if offset != nil { // [offset, n)
        randNum, _ := rand.Int(rand.Reader, Sub(n, offset, nil))
        return Add(randNum, offset, nil)
    } else { // [0, n)
        randNum, _ := rand.Int(rand.Reader, n)
        return randNum
    }
}

func IsPrime(x *big.Int) bool {
    return x.ProbablyPrime(20)
}

func GetPrime(nbits uint) *big.Int {
    offset := new(big.Int).Lsh(ONE, nbits-1)
    end := new(big.Int).Lsh(ONE, nbits)

    for ;; {
        k := Randint(offset, end)
        k_lsb := new(big.Int).And(k, ONE)
        k = Add(k, k_lsb.Xor(k_lsb, ONE), nil)
        //if IsPrime(k) {
        if k.ProbablyPrime(20) {
            return k
        }
    }
    return nil
}
