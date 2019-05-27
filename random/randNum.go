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

func GetPrime(nbits int64) *big.Int {
    nbits_big := big.NewInt(nbits)
    offset := Exp(TWO, Sub(nbits_big, ONE, nil), nil)
    end := Exp(TWO, nbits_big, nil)

    for ;; {
        k := Randint(offset, end)
        if Mod(k, TWO).Cmp(ZERO) == 0 {
            k = Add(k, ONE, nil)
        }
        if IsPrime(k) {
            return k
        }
    }
    return nil
}
