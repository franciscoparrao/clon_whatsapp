package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"
	
	"whatsapp-clone/services"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	jwtSecret []byte
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type AuthResponse struct {
	Token     string `json:"token"`
	UserID    string `json:"userId"`
	Username  string `json:"username"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	ExpiresAt int64  `json:"expiresAt"`
}

type Claims struct {
	UserID   string `json:"userId"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

func NewAuthHandler(jwtSecret string) *AuthHandler {
	return &AuthHandler{
		jwtSecret: []byte(jwtSecret),
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Validate input
	if req.Username == "" || req.Password == "" || req.Email == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Username, email and password are required"})
		return
	}

	// Check if user already exists
	existingUser, _ := services.Store.GetUserByEmail(req.Email)
	if existingUser != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Email already registered"})
		return
	}

	existingUser, _ = services.Store.GetUserByUsername(req.Username)
	if existingUser != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Username already taken"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create new user
	userID := fmt.Sprintf("user_%d", time.Now().UnixNano())
	newUser := &services.User{
		ID:       userID,
		Username: req.Username,
		Email:    req.Email,
		Name:     req.Name,
		Password: string(hashedPassword),
		CreatedAt: time.Now(),
	}

	// Save user to store
	if err := services.Store.CreateUser(newUser); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Create initial chats with bots
	services.InitializeBotChats(userID)

	// Generate JWT token
	token, expiresAt, err := h.generateToken(userID, req.Username, req.Email)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	response := AuthResponse{
		Token:     token,
		UserID:    userID,
		Username:  req.Username,
		Name:      req.Name,
		Email:     req.Email,
		ExpiresAt: expiresAt,
	}

	c.JSON(http.StatusOK, response)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Validate input
	if req.Username == "" || req.Password == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Username and password are required"})
		return
	}

	// Fetch user from store
	user, err := services.Store.GetUserByUsername(req.Username)
	if err != nil || user == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	userID := user.ID
	email := user.Email
	name := user.Name

	// Generate JWT token
	token, expiresAt, err := h.generateToken(userID, req.Username, email)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	response := AuthResponse{
		Token:     token,
		UserID:    userID,
		Username:  req.Username,
		Name:      name,
		Email:     email,
		ExpiresAt: expiresAt,
	}

	c.JSON(http.StatusOK, response)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	// TODO: Implement token blacklisting if needed
	// For now, just return success
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	// Get token from Authorization header
	tokenString := extractTokenFromHeader(c)
	if tokenString == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization token"})
		return
	}

	// Parse and validate token
	claims, err := h.validateToken(tokenString)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// Generate new token
	token, expiresAt, err := h.generateToken(claims.UserID, claims.Username, claims.Email)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":     token,
		"expiresAt": expiresAt,
	})
}

func (h *AuthHandler) generateToken(userID, username, email string) (string, int64, error) {
	expiresAt := time.Now().Add(24 * time.Hour)
	
	claims := &Claims{
		UserID:   userID,
		Username: username,
		Email:    email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(h.jwtSecret)
	if err != nil {
		return "", 0, err
	}

	return tokenString, expiresAt.Unix(), nil
}

func (h *AuthHandler) validateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return h.jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	return claims, nil
}

func extractTokenFromHeader(c *gin.Context) string {
	bearerToken := c.GetHeader("Authorization")
	if len(bearerToken) > 7 && strings.HasPrefix(bearerToken, "Bearer ") {
		return bearerToken[7:]
	}
	return ""
}

func generateUserID() string {
	// TODO: Implement proper ID generation
	return "user_" + time.Now().Format("20060102150405")
}

func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := services.Store.GetUser(userID)
	if err != nil || user == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"name":     user.Name,
	})
}