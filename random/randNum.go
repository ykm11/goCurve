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

