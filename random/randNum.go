package random

import (
    "math"
    "math/big"
    "crypto/rand"
    mrand "math/rand"
)

func Setup() error {
    seed, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
    mrand.Seed(seed.Int64())
    return err
}

func Randint(n *big.Int) *big.Int {
    randNum, _ := rand.Int(rand.Reader, n)
    return randNum
}
