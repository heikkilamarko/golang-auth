package utils

import (
	"fmt"
	"net/http"
	"strings"
)

type ctxKey string

var TokenKey = ctxKey("token")

func GetToken(r *http.Request) map[string]any {
	return r.Context().Value(TokenKey).(map[string]any)
}

func GetUserName(token map[string]any) string {
	for _, name := range []string{"name", "preferred_username"} {
		switch v := token[name].(type) {
		case string:
			return v
		}
	}
	return ""
}

func GetRolesAzure(token map[string]any) []string {
	var roles []string

	switch v := token["roles"].(type) {
	case []any:
		for _, role := range v {
			roles = append(roles, role.(string))
		}
	}

	return roles
}

func GetRolesKeycloak(resource string, token map[string]any) []string {
	var roles []string

	switch v1 := token["resource_access"].(type) {
	case map[string]any:
		switch v2 := v1[resource].(type) {
		case map[string]any:
			switch v3 := v2["roles"].(type) {
			case []any:
				for _, role := range v3 {
					roles = append(roles, role.(string))
				}
			}
		}
	}

	return roles
}

func TokenFromHeader(r *http.Request) string {
	a := r.Header.Get("Authorization")
	if 7 < len(a) && strings.ToUpper(a[0:6]) == "BEARER" {
		return a[7:]
	}
	return ""
}

func LogToken(r *http.Request) {
	token := TokenFromHeader(r)
	if token == "" {
		token = "empty token"
	}
	logTokenDivider()
	fmt.Println(token)
	logTokenDivider()
}

func logTokenDivider() {
	fmt.Println(strings.Repeat("-", 20), "TOKEN", strings.Repeat("-", 20))
}
