package middleware

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type AuthMiddleware struct {
	jwtSecret []byte
	skipPaths []string
}

type Claims struct {
	UserID   string `json:"userId"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

// NewAuthMiddleware creates a new auth middleware instance
func NewAuthMiddleware(jwtSecret string, skipPaths []string) *AuthMiddleware {
	return &AuthMiddleware{
		jwtSecret: []byte(jwtSecret),
		skipPaths: skipPaths,
	}
}

// Authenticate is the middleware function that validates JWT tokens
func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the path should skip authentication
		for _, path := range m.skipPaths {
			if strings.HasPrefix(r.URL.Path, path) {
				next.ServeHTTP(w, r)
				return
			}
		}

		// Extract token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Printf("Auth middleware: No authorization header for %s", r.URL.Path)
			http.Error(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}
		log.Printf("Auth middleware: Authorization header = %s", authHeader)

		// Check if the header starts with "Bearer "
		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		tokenString := bearerToken[1]

		// Parse and validate the token
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// Verify the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return m.jwtSecret, nil
		})

		if err != nil {
			log.Printf("Token validation error: %v", err)
			if err == jwt.ErrSignatureInvalid {
				http.Error(w, "Invalid token signature", http.StatusUnauthorized)
				return
			}
			if errors.Is(err, jwt.ErrTokenExpired) {
				http.Error(w, "Token has expired", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			log.Printf("Token is not valid")
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Add user information to request context
		ctx := context.WithValue(r.Context(), "userID", claims.UserID)
		ctx = context.WithValue(ctx, "username", claims.Username)
		ctx = context.WithValue(ctx, "email", claims.Email)
		ctx = context.WithValue(ctx, "claims", claims)

		// Call the next handler with the updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireAuth is a simpler middleware that just checks if the user is authenticated
func (m *AuthMiddleware) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("userID")
		if userID == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
}

// ExtractUserID is a helper function to extract user ID from context
func ExtractUserID(r *http.Request) (string, bool) {
	userID, ok := r.Context().Value("userID").(string)
	return userID, ok
}

// ExtractUsername is a helper function to extract username from context
func ExtractUsername(r *http.Request) (string, bool) {
	username, ok := r.Context().Value("username").(string)
	return username, ok
}

// ExtractEmail is a helper function to extract email from context
func ExtractEmail(r *http.Request) (string, bool) {
	email, ok := r.Context().Value("email").(string)
	return email, ok
}

// ExtractClaims is a helper function to extract full claims from context
func ExtractClaims(r *http.Request) (*Claims, bool) {
	claims, ok := r.Context().Value("claims").(*Claims)
	return claims, ok
}