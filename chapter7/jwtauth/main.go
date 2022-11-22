package main

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	jwt "github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

//using asymmetric crypto/RSA keys
//location of the files used for signing and verification
const (
	privKeyPath = "keys/app.rsa"     //openssl genrsa -out app.rsa 1024
	pubKeyPath  = "keys/app.rsa.pub" //openssl rsa -in app.rsa -pubout > app.rsa.pub
)

//verify key and sign key
var (
	verifyParsedKey *rsa.PublicKey
	signParsedKey   *rsa.PrivateKey
)

//User struct for parsing login credentials
type User struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

//read the key files before starting http handlers
func init() {
	var err error

	signKey, err := ioutil.ReadFile(privKeyPath)
	if err != nil {
		log.Fatal("Error reading private key")
		return
	}

	signParsedKey, err = jwt.ParseRSAPrivateKeyFromPEM(signKey)
	if err != nil {
		log.Fatal("Error parsing private key")
	}

	verifyKey, err := ioutil.ReadFile(pubKeyPath)
	if err != nil {
		log.Fatal("Error reading public key")
		return
	}

	verifyParsedKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyKey)
	if err != nil {
		log.Fatal("Error parsing verify key")
	}
}

//reads the login credentials, checks them and creates JWT the token
func loginHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	//decode into User string
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error in request body")
		return
	}
	//validate user credentials
	if user.UserName != "mihailmihaylov" && user.Password != "pass" {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintln(w, "Wrong info")
		return
	}

	//create signer for rsa 256
	t := jwt.New(jwt.SigningMethodRS256)

	//assign claims
	claims := t.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(20 * time.Minute).Unix()
	claims["iss"] = "admin"
	claims["user"] = "mihailmihaylov"
	claims["authorized"] = true

	tokenString, err := t.SignedString(signParsedKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Sorry, error while Signing Token!")
		log.Printf("token Signing error: %v\n", err)
		return
	}
	response := Token{tokenString}
	jsonResponse(response, w)
}

// only accessible with a valid token
func authHandler(w http.ResponseWriter, r *http.Request) {
	// validate the token
	token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
		// since we only use one private key to sign the tokens,
		// we also only use its public counterpart to verify
		_, ok := token.Method.(*jwt.SigningMethodRSA)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			_, err := w.Write([]byte("You're Unauthorized"))
			if err != nil {
				return nil, err
			}
		}
		return verifyParsedKey, nil
	})

	if err != nil {
		switch err.(type) {

		case *jwt.ValidationError: // something was wrong during the validation
			vErr := err.(*jwt.ValidationError)

			switch vErr.Errors {
			case jwt.ValidationErrorExpired:
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprintln(w, "Token Expired, get a new one.")
				return

			default:
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintln(w, "Error while Parsing Token!")
				log.Printf("ValidationError error: %+v\n", vErr.Errors)
				return
			}

		default: // something else went wrong
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "Error while Parsing Token!")
			log.Printf("Token parse error: %v\n", err)
			return
		}

	}
	if token.Valid {
		response := Response{"Authorized to the system"}
		jsonResponse(response, w)
	} else {
		response := Response{"Invalid token"}
		jsonResponse(response, w)
	}

}

type Response struct {
	Text string `json:"text"`
}
type Token struct {
	Token string `json:"token"`
}

func jsonResponse(response interface{}, w http.ResponseWriter) {
	json, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

//Entry point of the program
func main() {
	r := mux.NewRouter()
	r.HandleFunc("/login", loginHandler).Methods("POST")
	r.HandleFunc("/auth", authHandler).Methods("POST")

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	log.Println("Listening...")
	server.ListenAndServe()
}
