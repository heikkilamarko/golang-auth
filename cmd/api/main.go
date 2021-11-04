package main

import (
	"context"
	"fmt"
	"goauth/internal/utils"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/heikkilamarko/goutils"
	"github.com/heikkilamarko/goutils/middleware"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

func main() {
	err := godotenv.Load("api.env")
	checkErr(err)

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	logger := zerolog.New(os.Stderr).
		With().
		Timestamp().
		Logger()

	jwtConfig := &middleware.JWTConfig{
		Issuer:   os.Getenv("AUTH_ISSUER"),
		Iss:      os.Getenv("AUTH_CLAIM_ISS"),
		Aud:      []string{os.Getenv("AUTH_CLAIM_AUD")},
		TokenKey: utils.TokenKey,
	}

	r := mux.NewRouter()

	r.Use(
		middleware.Logger(&logger),
		middleware.RequestLogger(),
		middleware.ErrorRecovery(),
	)

	r.HandleFunc("/api/public", handlePublic).Methods("GET")

	s := r.PathPrefix("/api/private").Subrouter()
	s.Use(middleware.JWT(context.Background(), jwtConfig))
	s.HandleFunc("", handlePrivate).Methods("GET")

	log.Fatal(http.ListenAndServe(":8090", r))
}

func handlePublic(w http.ResponseWriter, r *http.Request) {
	goutils.WriteOK(w, "Hello from the public endpoint!", nil)
}

func handlePrivate(w http.ResponseWriter, r *http.Request) {
	utils.LogToken(r)

	token := utils.GetToken(r)

	userName := utils.GetUserName(token)

	var roles []string
	ap := os.Getenv("AUTH_PROVIDER")
	switch ap {
	case "azure":
		roles = utils.GetRolesAzure(token)
	case "keycloak":
		roles = utils.GetRolesKeycloak("todo-api", token)
	default:
		fmt.Printf("invalid auth provider: %s", ap)
	}

	fmt.Printf("username: '%s'\n", userName)
	fmt.Printf("roles (%d):\n", len(roles))
	for _, role := range roles {
		fmt.Printf("  - %s\n", role)
	}

	goutils.WriteOK(w, token, nil)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
