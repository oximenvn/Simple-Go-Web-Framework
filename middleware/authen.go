package middleware

import (
	"context"
	//"encoding/json"
	"log"
	"net/http"
)

func AuthHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		authToken := r.Header.Get("Authorization")
		user, err := getUser(authToken)

		if err != nil {
			http.Error(w, http.StatusText(401), 401)
			return
		}

		// add value to Context
		requestContext := r.Context()
		log.Printf("%+v\n", requestContext)
		for key, value := range user {
			requestContext = context.WithValue(requestContext, key, value)
		}
		log.Printf("%+v\n", requestContext)
		r = r.WithContext(requestContext)

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func getUser(token string) (map[string]string, error) {
	result := make(map[string]string)
	result["id"] = "id1"
	return result, nil
}
