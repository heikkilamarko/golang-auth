package main

import (
	"context"
	"fmt"
	"goauth/internal/utils"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

	jwtConfig := &utils.JWTConfig{
		Issuer:   os.Getenv("AUTH_ISSUER"),
		Iss:      os.Getenv("AUTH_CLAIM_ISS"),
		Aud:      []string{os.Getenv("AUTH_CLAIM_AUD")},
		TokenKey: utils.TokenKey,
		Logger:   &logger,
	}

	router := chi.NewRouter()

	router.Use(middleware.Recoverer)

	router.Get("/api/public", handlePublic)

	router.Route("/api/private", func(r chi.Router) {
		r.Use(utils.JWT(context.Background(), jwtConfig))
		r.Get("/", handlePrivate)
	})

	log.Fatal(http.ListenAndServe(":8090", router))
}

func handlePublic(w http.ResponseWriter, r *http.Request) {
	utils.WriteResponse(w, http.StatusOK, utils.NewDataResponse("Hello from the public endpoint!", nil))
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
		roles = utils.GetRolesKeycloak(os.Getenv("AUTH_CLAIM_AUD"), token)
	default:
		fmt.Printf("invalid auth provider: %s", ap)
	}

	fmt.Printf("username: '%s'\n", userName)
	fmt.Printf("roles (%d):\n", len(roles))
	for _, role := range roles {
		fmt.Printf("  - %s\n", role)
	}

	utils.WriteResponse(w, http.StatusOK, utils.NewDataResponse(token, nil))
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
