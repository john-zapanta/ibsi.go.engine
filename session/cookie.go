package session

import (
	"fmt"
	// "log"
	"time"
	"context"
	"math/rand"
	"net/http"
	"encoding/base64"
	// "github.com/gorilla/sessions"
)

// var store = sessions.NewCookieStore([]byte("something-very-secret"))
type TNewSession func(r *http.Request, sessionId string) (int64)
type TLogin func(r *http.Request, session *Session) (bool)
type TLogout func(r *http.Request, session *Session) (bool)

var OnSessionStart TNewSession = nil
var OnLogin TLogin = nil
var OnLogout TLogout = nil

type Session struct {
	SessionId string
	Authenticated bool
	VisitorId int64
	Home string
	Vars map[string]interface{}
}

var Sessions = make(map[string]*Session)
	
func NewSession(r *http.Request, sessionId string) *Session {
	return &Session{
		SessionId: sessionId,
		Authenticated: false,
		VisitorId: NewVisit(r, sessionId),
		Home: "/",
		// Vars: make(map[string]interface{}),
		// Vars: make(map[string]interface{}, {"login_attempts":0}),
		Vars: map[string]interface{}{"login_attempts":0},
	}
}
	
func GetSession(r *http.Request) *Session {
	if cookie, err := r.Cookie("session-id"); err == nil && cookie != nil {
		return Sessions[cookie.Value]
	}
	
	return nil
}

// all requests will be capure by CheckCookieHandler to verify if a cookie was passed by the browser, if not create a new cookie
func VerifyCookieHandler(next http.Handler) http.Handler {
	// fmt.Println("3. here")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if cookie, err := r.Cookie("session-id"); err != nil {
			sessionId := RandomSessionID(24)
			CreateSession(w, r, sessionId, time.Now().Add(365 * 24 * time.Hour), 0)
			// log.Println("New session created: " + sessionId, r.RequestURI)
			// fmt.Println("**************************************************************************************************************")
			// fmt.Println("** New Session: " + sessionId, r.RequestURI)
			// fmt.Println("**************************************************************************************************************")
			next.ServeHTTP(w, r)
		} else if cookie != nil {
			if session, ok := Sessions[cookie.Value]; ok {
				ctx := context.WithValue(r.Context(), "session", session)
				next.ServeHTTP(w, r.WithContext(ctx))
				// fmt.Println("**************************************************************************************************************")
				// fmt.Println("** Session found ", Sessions[cookie.Value], r.RequestURI)
				// fmt.Println("**************************************************************************************************************")
			} else {
				// cookie was passed, but for some reason the session was not found, probably server was restartedm
				// reuse the session
				Sessions[cookie.Value] = NewSession(r, cookie.Value)
				ctx := context.WithValue(r.Context(), "session", Sessions[cookie.Value])
				next.ServeHTTP(w, r.WithContext(ctx))
				// log.Println("** Cookie present but session not found, reset", Sessions[cookie.Value], r.RequestURI)
				// log.Println("** Cookie present but session not found, reset", Sessions[cookie.Value].SessionId, r.RequestURI)
				// fmt.Println("**************************************************************************************************************")
				// fmt.Println("** Cookie present but session not found, reset", Sessions[cookie.Value], r.RequestURI)
				// fmt.Println("**************************************************************************************************************")
			}
		}
	})
}

func NewVisit(r *http.Request, sessionId string) int64 {
	if OnSessionStart != nil {
		return OnSessionStart(r, sessionId)
	}
	
	return 0
}

func RandomSessionID(size int) string {
    buf := make([]byte, size)

    if _, err := rand.Read(buf); err != nil {
		fmt.Println(err.Error())
        // log.Println(err)
        // return "", errors.New("Couldn't generate random string")
    }

    return base64.URLEncoding.EncodeToString(buf)[:size]
	
	// return string(buf[:]) //"..."
}

// type Cookie struct {
        // Name       string
        // Value      string
        // Path       string
        // Domain     string
        // Expires    time.Time
        // RawExpires string

    // MaxAge=0 means no 'Max-Age' attribute specified.
    // MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'
    // MaxAge>0 means Max-Age attribute present and given in seconds
        // MaxAge   int
        // Secure   bool
        // HttpOnly bool
        // Raw      string
        // Unparsed []string // Raw text of unparsed attribute-value pairs
    // }
	
func CreateSession(w http.ResponseWriter, r *http.Request, sessionId string, expiration time.Time, maxAge int) {
	// expiration := time.Now().Add(365 * 24 * time.Hour)
	
	cookie := http.Cookie{
		Name: "session-id",
		Value: sessionId,
		Path: "/",
		HttpOnly: true,
		Expires: expiration,
		MaxAge: maxAge,
	}
	
	http.SetCookie(w, &cookie)
	
	Sessions[sessionId] = NewSession(r, sessionId)
	// Sessions[sessionId] = Session{
		// SessionId: sessionId,
		// Authenticated: false,
		// VisitorId: NewVisit(r, sessionId),
		// Vars: make(map[string]interface{}),
	// }
}
