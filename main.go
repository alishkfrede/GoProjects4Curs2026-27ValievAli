package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// ==================== Models ====================

type User struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
}

type Note struct {
	ID        string    `json:"id"`
	UserID    string    `json:"-"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateNoteRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

// ==================== JWT ====================

var jwtSecret = []byte("super_secret_key_change_in_production_123")

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func generateToken(user User) (string, error) {
	claims := Claims{
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ==================== Store ====================

type Store struct {
	mu    sync.RWMutex
	users map[string]User
	notes map[string]Note
}

func NewStore() *Store {
	return &Store{
		users: make(map[string]User),
		notes: make(map[string]Note),
	}
}

// ==================== Context ====================

type contextKey string

const userIDKey contextKey = "userID"

// ==================== Middleware ====================

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			respondError(w, http.StatusUnauthorized, "missing token")
			return
		}
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(authHeader[7:], claims, func(t *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			respondError(w, http.StatusUnauthorized, "invalid token")
			return
		}
		ctx := context.WithValue(r.Context(), userIDKey, claims.Subject)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// ==================== Handlers ====================

func (s *Store) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if req.Email == "" || req.Password == "" {
		respondError(w, http.StatusBadRequest, "email and password are required")
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.users[req.Email]; exists {
		respondError(w, http.StatusConflict, "user already exists")
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "hash error")
		return
	}
	user := User{ID: uuid.New().String(), Email: req.Email, PasswordHash: string(hash)}
	s.users[user.Email] = user
	respondJSON(w, http.StatusCreated, map[string]string{"message": "registered", "id": user.ID})
}

func (s *Store) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	s.mu.RLock()
	user, exists := s.users[req.Email]
	s.mu.RUnlock()
	if !exists || bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)) != nil {
		respondError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}
	token, err := generateToken(user)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "token error")
		return
	}
	respondJSON(w, http.StatusOK, LoginResponse{Token: token})
}

func (s *Store) ListNotes(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIDKey).(string)
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := []Note{}
	for _, n := range s.notes {
		if n.UserID == userID {
			result = append(result, n)
		}
	}
	respondJSON(w, http.StatusOK, result)
}

func (s *Store) CreateNote(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIDKey).(string)
	var req CreateNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if req.Title == "" {
		respondError(w, http.StatusBadRequest, "title is required")
		return
	}
	note := Note{
		ID: uuid.New().String(), UserID: userID,
		Title: req.Title, Content: req.Content, CreatedAt: time.Now(),
	}
	s.mu.Lock()
	s.notes[note.ID] = note
	s.mu.Unlock()
	respondJSON(w, http.StatusCreated, note)
}

func (s *Store) GetNote(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIDKey).(string)
	noteID := r.PathValue("id")
	s.mu.RLock()
	note, exists := s.notes[noteID]
	s.mu.RUnlock()
	if !exists || note.UserID != userID {
		respondError(w, http.StatusNotFound, "not found")
		return
	}
	respondJSON(w, http.StatusOK, note)
}

func (s *Store) DeleteNote(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIDKey).(string)
	noteID := r.PathValue("id")
	s.mu.Lock()
	defer s.mu.Unlock()
	note, exists := s.notes[noteID]
	if !exists || note.UserID != userID {
		respondError(w, http.StatusNotFound, "not found")
		return
	}
	delete(s.notes, noteID)
	w.WriteHeader(http.StatusNoContent)
}

// ==================== Helpers ====================

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if payload != nil {
		json.NewEncoder(w).Encode(payload)
	}
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}

// ==================== Main ====================

func main() {
	store := NewStore()
	mux := http.NewServeMux()

	// API
	mux.HandleFunc("POST /register", store.Register)
	mux.HandleFunc("POST /login", store.Login)
	mux.Handle("GET /notes", authMiddleware(http.HandlerFunc(store.ListNotes)))
	mux.Handle("POST /notes", authMiddleware(http.HandlerFunc(store.CreateNote)))
	mux.Handle("GET /notes/{id}", authMiddleware(http.HandlerFunc(store.GetNote)))
	mux.Handle("DELETE /notes/{id}", authMiddleware(http.HandlerFunc(store.DeleteNote)))

	// Статика (фронтенд)
	fs := http.FileServer(http.Dir("./frontend"))
	mux.Handle("/", fs)

	handler := corsMiddleware(mux)

	log.Println("🚀 Server running at http://localhost:8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}