package system

import (
	"fmt"
	"io/ioutil"
	"os"
	"log"
	"encoding/json"
	// "regexp"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	_ "github.com/denisenkom/go-mssqldb"
	// "ibsi/system"
	"ibsi/dbase"
	// "ibsi/crud"
	"ibsi/bundle"
	"ibsi/template"
	"ibsi/session"
)

type DatabaseConnection struct {
    Name string `json:"name"`
    Server string `json:"server"`
    Port int `json:"port"`
    Database string `json:"database"`
    User string `json:"user"`
    Password string `json:"password"`
}

type Config struct {
    AppID int `json:"app_id"`
    AppName string `json:"app_name"`
    AppTitle string `json:"app_title"`
    AppIcon string `json:"app_icon"`
    ResPath	string `json:"res_path"`
    DocPath	string `json:"doc_path"`
    Domain string `json:"domain"`
    Port int `json:"port"`
	Connections []DatabaseConnection `json:"connections"`
}

// func (c *Config) MarshalJSON() ([]byte, error) {
	// type Alias Config
	// return json.Marshal(&struct {
		// *Alias
	// }{
		// Alias: (*Alias)(c),
	// })
// }

func (c *Config) UnmarshalJSON(data []byte) error {
	type Alias Config
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(c),
	}
	
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	c.init()
	
	return nil
}

func (c *Config) init() {
	for i := 0; i < len(c.Connections); i++ {
		cn := c.Connections[i]
		dbase.NewConnection(cn.Name, cn.Server, cn.Port, cn.Database, cn.User, cn.Password)
	}
}

func (c *Config) ShowConnections() {
	for i := 0; i < len(c.Connections); i++ {
		cn := c.Connections[i]
		// fmt.Printf("Connection to %s (%s:%d) is ready, total commands %d\n", cn.Name, cn.Server, cn.Port, len(cn.commands))
		// log.Println(fmt.Sprintf("Connection to %s (%s:%d) is ready", cn.Name, cn.Server, cn.Port))
		log.Println(fmt.Sprintf("%s is connected to database %s (%s:%d)", cn.Name, cn.Database, cn.Server, cn.Port))
	}
}

// Simple rewrite to force user to login
func Rewriter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if sn := session.GetSession(r); sn != nil {
			// to-do: Testing, auto login, remove later
			if !sn.Authenticated {
				sn.Vars["user_name"] = "super"
				sn.Vars["password"] = "pass1234"
				sn.Authenticated = true
				session.OnLogin(r, sn)
			}
			
			// to-do: Restore later
			// if match, _ := regexp.MatchString("^/app/([a-z]+)", r.URL.Path); !sn.Authenticated && match {
				// r.URL.Path = "/login"
			// } else if sn.Authenticated && r.URL.Path == "/" {
				// r.URL.Path = sn.Home
			// }
		}
		next.ServeHTTP(w, r)
	})
}

// var config *Config
var Settings *Config
var Router *mux.Router
var FinalHandler http.Handler

func init() {
	// log.SetFlags(0)
	// log.SetOutput(ioutil.Discard)
	
	Router = mux.NewRouter()
	
	// Matches a dynamic subdomain.
	// r.Host("{subdomain:[a-z]+}.domain.com")
	
	h1 := handlers.CombinedLoggingHandler(os.Stdout, Router)
    // FinalHandler = handlers.CompressHandler(h1)
    // FinalHandler = session.VerifyCookieHandler(handlers.CompressHandler(h1))
    FinalHandler = Rewriter(session.VerifyCookieHandler(handlers.CompressHandler(h1)))

    raw, err := ioutil.ReadFile("config.json")
    if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }

    err = json.Unmarshal(raw, &Settings)
    if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }

	// Settings.AppID = Settings.AppID
	
	template.Router = Router
	// template.Domain = config.Domain
	template.AppIcon = Settings.AppIcon
	// crud.Router = Router
	// crud.Domain = config.Domain
	
	// initialize the resource path of the bundler
	bundle.ResourcePath = Settings.ResPath
	
	// initialize the OnSessionStart event and log the session to the database
	// OnSessionStart is in ibsi.cookie.go
	session.OnSessionStart = func(r *http.Request, sessionId string) (int64) {
		visits, err := dbase.Connections["DBSecure"].Execute("AddVisit", dbase.TParameters{
			"application_id": Settings.AppID,
			"session_id": sessionId,
			"request_url": r.RequestURI,
			"method": r.Method,
			"remote_host": r.RemoteAddr,
			"user_agent": r.Header.Get("User-Agent"),
			"referrer_url": r.Referer(),
		})
		
		if err != nil {
			log.Println(err.Error())
		} else {
			out, _ := visits.GetOutput()
			return out.Get("visit_id").(int64)
		}

		return 0
	}
	
	session.OnLogin = func(r *http.Request, session *session.Session) (bool) {
		login, err := dbase.Connections["DBSecure"].Execute("Login", dbase.TParameters{
			"user_name": session.Vars["user_name"],
			"password": session.Vars["password"],
			"visit_id": session.VisitorId,
		})

		if err != nil {
			log.Println(err.Error())
		} else {						
			if out, _ := login.GetOutput(); out.Get("user_id").(int64) > 0 {
				session.Home = "/app/home"
				session.Vars["user_id"] = out.Get("user_id")
				
				if dbase.Connections["DBApp"].HasCommand("System_ManageSession") {
					dbase.Connections["DBApp"].Execute("System_ManageSession", dbase.TParameters{
						"id": session.VisitorId,
						"session_id": session.SessionId,
						"user_id": session.Vars["user_id"],
						"user_name": session.Vars["user_name"],
						"action": 10,
					})
				}			
				
				return true
			} 
		}

		return false
	}
	
	session.OnLogout = func(r *http.Request, session *session.Session) (bool) {
		_, err := dbase.Connections["DBSecure"].Execute("Logout", dbase.TParameters{
			"visit_id": session.VisitorId,
		})

		if err != nil {
			log.Println(err.Error())
		} else {						
			if dbase.Connections["DBApp"].HasCommand("System_ManageSession") {
				dbase.Connections["DBApp"].Execute("System_ManageSession", dbase.TParameters{
					"id": session.VisitorId,
					"action": 11,
				})
			}			
			
			return true
		}

		return false
	}
	
	// Router.Host(config.Domain)
}
