package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func getSessionSecret() []byte {
	secret := os.Getenv("SESSION_SECRET")
	if secret == "" {
		secret = "clave-secreta-ootd-tp6-cambiar-en-prod"
	}
	return []byte(secret)
}

func firmarID(idStr string) string {
	mac := hmac.New(sha256.New, getSessionSecret())
	mac.Write([]byte(idStr))
	return hex.EncodeToString(mac.Sum(nil))
}

func establecerCookieSesion(w http.ResponseWriter, usuarioID int32) {
	idStr := strconv.Itoa(int(usuarioID))
	http.SetCookie(w, &http.Cookie{
		Name:     "session_user_id",
		Value:    idStr + ":" + firmarID(idStr),
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

func borrarCookieSesion(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_user_id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

func obtenerUsuarioDeSesion(r *http.Request) (int32, bool) {
	cookie, err := r.Cookie("session_user_id")
	if err != nil || cookie.Value == "" {
		return 0, false
	}

	partes := strings.Split(cookie.Value, ":")
	if len(partes) != 2 {
		return 0, false
	}

	firmaEsperada := firmarID(partes[0])
	if !hmac.Equal([]byte(partes[1]), []byte(firmaEsperada)) {
		return 0, false
	}

	id, err := strconv.Atoi(partes[0])
	if err != nil || id <= 0 {
		return 0, false
	}

	return int32(id), true
}
