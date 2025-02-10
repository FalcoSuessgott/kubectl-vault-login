package jwt

import (
	_ "github.com/golang-jwt/jwt/v5"
)

// func Decode(token string) error {
// 	token, err := jwt.ParseWithClaims(token, func(token *jwt.Token) (interface{}, error) {
// 		// Don't forget to validate the alg is what you expect:
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
// 		}

// 		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
// 		return hmacSampleSecret, nil
// 	})
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	if claims, ok := token.Claims.(jwt.MapClaims); ok {
// 		fmt.Println(claims["foo"], claims["nbf"])
// 	} else {
// 		fmt.Println(err)
// 	}
// }
