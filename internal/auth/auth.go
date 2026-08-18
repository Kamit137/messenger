package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"
	"messenger/internal/jwt"
	
	"messenger/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

const cookieName = "messenger_token"

type authRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type userResponse struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}

type errorResponse struct {
	Error string `json:"error"`
}


func Register(w http.ResponseWriter, r *http.Request) {
	req, ok := decodeJSON(w, r)
	if !ok {
		return
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	if !validEmail(req.Email) {
		writeError(w, http.StatusBadRequest, "invalid email")
		return
	}
	if len(req.Password) < 6 {
		writeError(w, http.StatusBadRequest, "password too short")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "server error")
		return
	}

	id, err := storage.SaveUser(r.Context(), req.Email, hash)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			writeError(w, http.StatusConflict, "email already registered")
			return
		}
		writeError(w, http.StatusInternalServerError, "registration failed")
		return
	}

	if acceptsJSON(r) {
		writeJSON(w, http.StatusCreated, userResponse{ID: id, Email: req.Email})
		return
	}

	http.Redirect(w, r, "/login?registered=true", http.StatusFound)
}

func Login(w http.ResponseWriter, r *http.Request) {
	req, ok := decodeJSON(w, r)
	if !ok {
		return
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	user, err := storage.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			writeError(w, http.StatusUnauthorized, "invalid credentials")
			return
		}
		writeError(w, http.StatusInternalServerError, "login failed")
		return
	}

	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(req.Password)); err != nil {
		writeError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	token, err := jwt.GenerateToken(user.ID, user.Email)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "token generation failed")
		return
	}

	setCookie(w, token)

	if acceptsJSON(r) {
		writeJSON(w, http.StatusOK, userResponse{ID: user.ID, Email: user.Email})
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	if acceptsJSON(r) {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	http.Redirect(w, r, "/login", http.StatusFound)
}

func Me(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	claims, err := jwt.ValidateToken(cookie.Value)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "invalid token")
		return
	}

	user, err := storage.GetUserByID(r.Context(), claims.UserID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "user not found")
		return
	}

	writeJSON(w, http.StatusOK, userResponse{ID: user.ID, Email: user.Email})
}

// --- ВСПОМОГАТЕЛЬНЫЕ ФУНКЦИИ ---

func decodeJSON(w http.ResponseWriter, r *http.Request) (authRequest, bool) {
	defer r.Body.Close()

	var req authRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return authRequest{}, false
	}
	return req, true
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, errorResponse{Error: msg})
}

func setCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    token,
		Path:     "/",
		MaxAge:   86400,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

func acceptsJSON(r *http.Request) bool {
	return strings.Contains(r.Header.Get("Accept"), "application/json")
}

func validEmail(email string) bool {
	at := strings.IndexByte(email, '@')
	return at > 0 && at < len(email)-1 && strings.Contains(email[at+1:], ".")
}