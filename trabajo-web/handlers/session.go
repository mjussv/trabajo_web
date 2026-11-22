package handlers

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
	"trabajo-web/logic"

	"github.com/a-h/templ"
)

const sessionDuration = 7 * 24 * time.Hour

var (
	sessionSecretOnce sync.Once
	sessionSecretKey  []byte
)

func getSessionSecret() []byte {
	sessionSecretOnce.Do(func() {
		secret := os.Getenv("SESSION_SECRET")
		if secret != "" {
			sessionSecretKey = []byte(secret)
			return
		}
		log.Println("ADVERTENCIA: SESSION_SECRET no configurado; generando clave aleatoria en memoria.")
		buf := make([]byte, 32)
		if _, err := rand.Read(buf); err != nil {
			log.Fatalf("No se pudo generar clave de sesión segura: %v", err)
		}
		sessionSecretKey = buf
	})
	return sessionSecretKey
}

func firmarPayload(payload string) string {
	mac := hmac.New(sha256.New, getSessionSecret())
	mac.Write([]byte(payload))
	return hex.EncodeToString(mac.Sum(nil))
}

func esConexionSegura(r *http.Request) bool {
	return r.TLS != nil ||
		strings.EqualFold(r.Header.Get("X-Forwarded-Proto"), "https") ||
		strings.EqualFold(os.Getenv("COOKIE_SECURE"), "true")
}

func establecerCookieSesion(w http.ResponseWriter, r *http.Request, usuarioID int32) {
	idStr := strconv.Itoa(int(usuarioID))
	expStr := strconv.FormatInt(time.Now().Add(sessionDuration).Unix(), 10)
	payload := idStr + ":" + expStr

	http.SetCookie(w, &http.Cookie{
		Name:     "session_user_id",
		Value:    payload + ":" + firmarPayload(payload),
		Path:     "/",
		MaxAge:   int(sessionDuration.Seconds()),
		HttpOnly: true,
		Secure:   esConexionSegura(r),
		SameSite: http.SameSiteLaxMode,
	})
}

func borrarCookieSesion(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_user_id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   esConexionSegura(r),
		SameSite: http.SameSiteLaxMode,
	})
}

func obtenerUsuarioDeSesion(r *http.Request) (int32, bool) {
	cookie, err := r.Cookie("session_user_id")
	if err != nil || cookie.Value == "" {
		return 0, false
	}

	partes := strings.Split(cookie.Value, ":")
	if len(partes) != 3 {
		return 0, false
	}

	payload := partes[0] + ":" + partes[1]
	firmaEsperada := firmarPayload(payload)
	if !hmac.Equal([]byte(partes[2]), []byte(firmaEsperada)) {
		return 0, false
	}

	expUnix, err := strconv.ParseInt(partes[1], 10, 64)
	if err != nil || time.Now().Unix() > expUnix {
		return 0, false
	}

	id64, err := strconv.ParseInt(partes[0], 10, 32)
	if err != nil || id64 <= 0 {
		return 0, false
	}

	return int32(id64), true
}

// renderHTML renderiza todos los componentes Templ en un buffer antes de escribir headers o status HTTP.
func renderHTML(w http.ResponseWriter, r *http.Request, status int, components ...templ.Component) {
	var buf bytes.Buffer
	for _, c := range components {
		if err := c.Render(r.Context(), &buf); err != nil {
			log.Printf("Error renderizando vista HTML: %v", err)
			http.Error(w, "Error interno al generar la vista", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	if _, err := buf.WriteTo(w); err != nil {
		log.Printf("Error escribiendo respuesta HTTP: %v", err)
	}
}

func responderNoAutenticado(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/")
		w.WriteHeader(http.StatusOK)
		return
	}
	http.Error(w, "No autenticado", http.StatusUnauthorized)
}

func validarUsuarioActivo(w http.ResponseWriter, r *http.Request, usuariosSvc *logic.UsuariosService, usuarioID int32) bool {
	if _, err := usuariosSvc.GetUsuariosByID(r.Context(), usuarioID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			borrarCookieSesion(w, r)
			responderNoAutenticado(w, r)
			return false
		}
		log.Printf("Error verificando usuario %d: %v", usuarioID, err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return false
	}
	return true
}
