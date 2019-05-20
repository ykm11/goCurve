package mathUtils

import (
    "math/big"
)

func pollard_rho_f(alpha, beta, x, a, b, p, order *big.Int) (*big.Int, *big.Int, *big.Int) {
    if Mod(x, THREE).Cmp(ZERO) == 0 {
        return Exp(x, TWO, p), Mul(a, TWO, order), Mul(b, TWO, order)
    } else if Mod(x, THREE).Cmp(ONE) == 0 {
        return Mul(x, alpha, p), Add(a, ONE, order), b
    } else {
        return Mul(x, beta, p), a, Add(b, ONE, order)
    }
}

// when order is not prime, it doesn't work
func pollard_rho_DLP(alpha, beta, p, order *big.Int) *big.Int {
    // beta = alpha^x  (mod modulus)
    a, b, x := ZERO, ZERO, ONE
    A, B, X := ValCopy(a), ValCopy(b), ValCopy(x)

    for i := ONE; i.Cmp(order) == -1; i = Add(i, ONE, nil) {
        x, a, b = pollard_rho_f(alpha, beta, x, a, b, p, order)
        X, A, B = pollard_rho_f(alpha, beta, X, A, B, p, order)
        X, A, B = pollard_rho_f(alpha, beta, X, A, B, p, order)

        //fmt.Println(i, x, a, b, X, A, B)
        if x.Cmp(X) == 0 &&
        (A.Cmp(ZERO) != 0 && a.Cmp(ZERO) != 0 && B.Cmp(ZERO) != 0 && b.Cmp(ZERO) != 0) {
            r := Sub(B, b, order)
            if r.Cmp(ZERO) != 0 {
                return Mul(InvMod(r, order), Sub(a, A, order), order)
            }
        }
    }
    return nil
}
