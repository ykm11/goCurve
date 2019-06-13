package main

import (
    "fmt"
    . "./ecUtils"
    . "./mathUtils"
    "./random"
)

func main() { // Eaxmple
    random.Setup()

    A := Str2Int("1461501637330902918203684832716283019653785059324", 10)
    B := Str2Int("163235791306168110546604919403271579530548345413", 10)
    modulus := Str2Int("1461501637330902918203684832716283019653785059327", 10)

    n := Str2Int("1461501637330902918203687197606826779884643492439", 10)

    EC := NewCurve(A, B, modulus)
    EC.PrintCurve()
    //fmt.Printf("[+] Order n: %d\n", n)

    gx := Str2Int("598833001378563909320556562387727035658124457364", 10)
    gy := Str2Int("456273172676936625440583883939668862699127599796", 10)
    G := EC.Point(gx, gy)
    fmt.Println("[+] Base Point G:", Point2Str(G))

    m := "This is a message"
    e := GetLeftmostNBits(Bytes2Int(GetHash(m)), 126, 256)
    fmt.Printf("[+] Message: %s\n", m)
    fmt.Printf("[+] Message Hash: %x\n", e)
    d := random.Randint(nil, n)
    Q := EC.Point_xP(d, G)
    fmt.Println("[+] Q:", Point2Str(Q))

    // Signature
    r, s := EC.Sign(e, d, n, G)

    // Verify
    result := EC.VerifySignature(r, s, e, n, G, Q)
    fmt.Println("[+] Result of Validation:", result)

}
