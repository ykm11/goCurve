package mathUtils

import (
    "crypto/sha256"
    "math/big"
)

func Str2Int(str string, base int) *big.Int {
    n, ok := new(big.Int).SetString(str, base)
    if !ok {
        panic("UWAAAAAAAAAAA")
    }
    return n
}

func Byte2Int(b []byte) *big.Int {
    return new(big.Int).SetBytes(b)
}

func Bool2int(b bool) int {
    if b {
        return 1
    } else {
        return 0
    }
}

func ValCopy(x *big.Int) *big.Int {
    r := new(big.Int).SetBytes(x.Bytes())
    return r
}

func GetLeftmostNBits(x *big.Int, n_bits, x_bits uint) *big.Int {
    return new(big.Int).Rsh(x, x_bits-n_bits)
}
func GetHash(message string) []byte {
    hash := sha256.Sum256([]byte(message))
    return hash[:]
}

