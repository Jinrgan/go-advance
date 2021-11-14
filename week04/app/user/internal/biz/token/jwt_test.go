package token

import (
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const privateKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEA0S/jrkNRQihPPU2m9G9k3KELp2NqQDNXfNG9Oo0LXmGztzgp
ZhhIo9btkmGXBpHgY9VOaCwoNQ6Y6BPgoKa2jAlT4HJudZ1R6mQZGfQrpUkby8pW
A+kFVM4L0OExUaqnVHGjU0gzcThJozPJAIOwxjb4upXgI51akOHFm8V53p7gIikt
4Xn8wtBZV9rUX3B2EKzezGGbC9Eg+WxpwLdJE8UynGV86M8KKtK7n74Ox2m+TfJn
59EAzGHbVKVAxeVl/HjYAZGOHOrg1BsHUm9/aCAJcnF+CyEldYdmLbkSVLyAGRKG
tbeIT0bV4Sim0TmjfMqF34dI50+KSkJfSTUEfQIDAQABAoIBAQCPFt/6CsghpesV
9vD2EOCAXpTXKhS616PHmLyEuYgGRnSlJoCC+qdtkw4s7B5fexdvkrAwZ8wVBugn
D7m+imsh/Rtn0z6lqzgmSdQ1waS9SfX+f2g5AoMIEG1asz+GKmKNS7I5vJCbqLIO
NdUPSgV4gI/BKdYI5pDVu+ns9La5RN0uaFiSPMh8LZgyTGGMKjjoCpqawO/LCEi6
5faS6Uw7H2s7pqiq3DJYoIm81ZoqhHvfwSNHCIapDPcUARnhPhv7R3KQqKhy/dQz
3yNN3e49UqQKyeQ5HNuRW/HI4mOo4twiJVwBP3P0J4G/ev33kVSsgx3tQd5leSmF
tFKsbWFNAoGBAO827pjMQE+QSgSeUui6E0JW8u6LzNKhpfttJGFikH5l64LmNuIa
aXLHYk2oC/cKlGVcVzysdrpEFm9qWeDKaLDp628k02OAd8khJ6tBgzUtXPVeyCVz
1yiumC91ibkiOJg5nf2r3KcIadl7nnzwxbrLIOo4cqXb7so7B5B0EzBzAoGBAN/d
kOu+VLk2cRXeb2fjXAracfD/kZlUwlSFsW+gmzxLaceWVDlov8cLV7NcSp3ptwQp
us+iDZEh2sA5/YpOquqG7oZYbtQgUpeuRY6Nq2N1HQ7fO6POHDaRJ6aS4eSBbUqt
bXnlcRE8kHzIPL+PPaPNZmorE4rdJcP9OCrgk2tPAoGBALI+a0DdiNoAOLcCReL1
q54V7cRD1SXpnyUeaSpLaEFWrksGQUTuyz3kRWJ54hh9AKAaU0J5e6pFS7ZPN5Nh
HuscEfrqJL0Sn671jnp0QVEhcQ/ARUBq9ZpxpiJO4YVac3MyE4BOTAcGJOER1MFi
IuORsf/0ebEOlPqJS5SeeSHzAoGACQ4yXYbecHuGSYcs5Hvq7jl14HTGE/i8v6SE
z+okPWUji0JGd+gH0epgis3R6t9YWt/BQcLhX5yJ97qgyeZyvXfl0CNloEkKbj5L
a//JDgXfvglDpVWiCIcInpFUd+TQYfPv+L1SPItBoPqMkocdzDFz0hmZ+cUGUQ4+
JmXdMLsCgYEA7VQmLzNtjh6ZrGst/OhDeQE86AjP1vkzIig2GlCn0aV7x1MsVW6k
u7QIXvgWxer5ZJrSNMQPrBOjOFQsySrC9AtIAOCnOLRCTibiRyqtl16beIRP7U98
wJHn9DcgplFs6luEoVO74Ju63MjHNPuPfQlZoGS3pHcQxRn4tk/fTg8=
-----END RSA PRIVATE KEY-----`

func TestJWTTokenGen_GenerateToken(t *testing.T) {
	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKey))
	if err != nil {
		t.Fatalf("cannot parse private key: %v", key)
	}

	g := NewJWTTokenGen("e-commerce/user", key)
	g.nowFunc = func() time.Time {
		return time.Unix(1516239022, 0)
	}

	tkn, err := g.GenerateToken("5", 2*time.Hour)
	if err != nil {
		t.Errorf("cannot generate token: %v", err)
	}

	want := "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTYyNDYyMjIsImlhdCI6MTUxNjIzOTAyMiwiaXNzIjoiZS1jb21tZXJjZS91c2VyIiwic3ViIjoiNSJ9.Da1OUSTN_GEwLNErwR8nfn3C3ZJnkjF8rNj578upuobuJfE9ZFa8rXpI2FoH678NWRf4fX5v8WnrRbf6pskVRWZh3TzNaukYTccGzKiPFCrEVOjqt6mHPLgecdgrOmz7dJ__Lt8AKPv3DgmjYnge0jlmIUsT3rx_w-SwVvpKJ0KCr9wZanNhivPAurBJ0r-Dw-7arixYfOsDrAlq7VOtkHz4BakI7pvSzxJhlzhPWcQjZGQxagIV9JX5-JlcMu2f5LdfUdQw_eQ2qRuHNAPmP9h1X0GfoEXKlOwAe8RcQlt8tLoL6h2fcQh9gHTju_hFoQuiJCziO-gejqBpY4r8Kg"
	if tkn != want {
		t.Errorf("wrong token generated. want %q; got %q", want, tkn)
	}
}
