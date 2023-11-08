package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/openpubkey/openpubkey/client"
	"github.com/openpubkey/openpubkey/client/cosigner/mfa"
	"github.com/openpubkey/openpubkey/client/providers"
	"github.com/openpubkey/openpubkey/examples/mfa/webauthn"
	"github.com/openpubkey/openpubkey/util"
)

// Variables for building our google provider
var (
	clientID = "184968138938-g1fddl5tglo7mnlbdak8hbsqhhf79f32.apps.googleusercontent.com"
	// The clientSecret was intentionally checked in for the purposes of this example,. It holds no power. Do not report as a security issue
	clientSecret = "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" // Google requires a ClientSecret even if this a public OIDC App
	issuer       = "https://accounts.google.com"
	scopes       = []string{"openid profile email"}
	redirURIPort = "3000"
	callbackPath = "/login-callback"
	redirectURI  = fmt.Sprintf("http://localhost:%v%v", redirURIPort, callbackPath)
)

func main() {
	authenticator, err := webauthn.New()
	if err != nil {
		fmt.Println("error instantiating mfa:", err.Error())
		return
	}

	fmt.Println("MFA ready, now initializing cosigner")

	cosigner, err := mfa.NewCosigner(authenticator)
	if err != nil {
		fmt.Println("error instantiating cosigner:", err.Error())
		return
	}

	signer, err := util.GenKeyPair(jwa.ES256)
	if err != nil {
		fmt.Println("error generating key pair:", err.Error())
		return
	}

	client := &client.OpkClient{
		Op: &providers.GoogleOp{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Issuer:       issuer,
			Scopes:       scopes,
			RedirURIPort: redirURIPort,
			CallbackPath: callbackPath,
			RedirectURI:  redirectURI,
		},
	}

	pkt, err := client.OidcAuth(context.TODO(), signer, jwa.ES256, map[string]any{"extra": "yes"}, false)
	if err != nil {
		fmt.Println("error generating key pair: ", err.Error())
		return
	}

	if err := cosigner.Cosign(pkt); err != nil {
		fmt.Println("error cosigning:", err.Error())
		return
	}

	pktJson, _ := json.MarshalIndent(pkt, "", "  ")
	fmt.Println(string(pktJson))

	select {}
}
