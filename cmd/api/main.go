package main

import (
	"context"
	"fmt"
	"goauth/internal/utils"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("api.env")
	checkErr(err)

	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	handler := slog.NewJSONHandler(os.Stderr, opts)

	logger := slog.New(handler)

	slog.SetDefault(logger)

	jwtConfig := &utils.JWTMiddlewareConfig{
		Issuer:     os.Getenv("AUTH_ISSUER"),
		Iss:        os.Getenv("AUTH_CLAIM_ISS"),
		Aud:        []string{os.Getenv("AUTH_CLAIM_AUD")},
		ContextKey: utils.TokenKey,
		Logger:     logger,
	}

	router := chi.NewRouter()

	router.Use(middleware.Recoverer)

	router.Get("/api/public", handlePublic)

	router.
		With(utils.JWTMiddleware(context.Background(), jwtConfig)).
		Get("/api/private", handlePrivate)

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
