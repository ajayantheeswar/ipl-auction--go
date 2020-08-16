package utils

import(
	"github.com/dgrijalva/jwt-go"
	"errors"
)



func GetAuthClaims (token string) (string,error){   
	claims := jwt.MapClaims{}

	Parsedtoken, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("jdnfksdmfksd"), nil
	})
	// ... error handling
	if err != nil{
		return "" ,errors.New("Invalid Token")
	}
	if Parsedtoken.Valid {
		return claims["user_id"].(string) , nil
	}else {
		return "" , errors.New("Invlalid Token")
	}
}


func CreateToken(Id string) (string, error) {
	var err error
	
	Claims := jwt.MapClaims{}
	Claims["user_id"] = Id

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims)
	token, err := at.SignedString([]byte("jdnfksdmfksd"))
	if err != nil {
	   return "", err
	}
	return token, nil
  }