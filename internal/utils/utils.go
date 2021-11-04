package utils

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/heikkilamarko/goutils"
)

type ctxKey string

var TokenKey = ctxKey("token")

func GetToken(r *http.Request) map[string]interface{} {
	return r.Context().Value(TokenKey).(map[string]interface{})
}

func GetUserName(token map[string]interface{}) string {
	switch v := token["name"].(type) {
	case string:
		return v
	default:
		return ""
	}
}

func GetRolesAzure(token map[string]interface{}) []string {
	var roles []string

	switch v := token["roles"].(type) {
	case []interface{}:
		for _, role := range v {
			roles = append(roles, role.(string))
		}
	}

	return roles
}

func GetRolesKeycloak(resource string, token map[string]interface{}) []string {
	var roles []string

	switch v1 := token["resource_access"].(type) {
	case map[string]interface{}:
		switch v2 := v1[resource].(type) {
		case map[string]interface{}:
			switch v3 := v2["roles"].(type) {
			case []interface{}:
				for _, role := range v3 {
					roles = append(roles, role.(string))
				}
			}
		}
	}

	return roles
}

func LogToken(r *http.Request) {
	token := goutils.TokenFromHeader(r)
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
