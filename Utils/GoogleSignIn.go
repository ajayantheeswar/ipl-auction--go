package utils

import (
	"google.golang.org/api/idtoken"
	"ipl/firebase"
	"context"
)

/*
func getTokenInfo(idToken string) (*oauth2.Tokeninfo, error) {
	oauth2Service, err := oauth2.New(&http.Client{})
	if err != nil {
		return nil, err
	}
	tokenInfoCall := oauth2Service.Tokeninfo()
	tokenInfoCall.IdToken(idToken)
	return tokenInfoCall.Do()
}*/



func VerifyIDToken(idToken string) (*idtoken.Payload, error) {
	
    // this comes from your web or mobile app maybe
	googleClientID := "GoogleCLIENTID"  // from credentials in the Google dev console
	tokenValidator, err := idtoken.NewValidator(context.Background(),firebase.ClientOpts)
	if err != nil {
		return nil , err
	}

	payload, err := tokenValidator.Validate(context.Background(), idToken, googleClientID)
	if err != nil {
		return nil , err
	}

	return payload , nil
	/*
	email := payload.Claims["email"]
	name  := payload.Claims["name"]
	*/
}
