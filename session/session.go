package session

import (
	"fmt"
	"net/url"
	"net/http"
	"ibsi/utils"
)
	
func GetVisitorId(r *http.Request) int64 {
	var vid int64
	sn := GetSession(r)
	if (sn == nil) {
		vid = 1
	} else {
		vid = sn.VisitorId
	}
	
	return vid
}
	
func GetUserId(r *http.Request) int {
	var id int
	sn := GetSession(r)
	if (sn == nil) {
		id = 1
	} else {
		id = sn.Vars["user_id"].(int)
	}
	
	return id
}

func Login(w http.ResponseWriter, r *http.Request) {
	
	// w.Header().Set("Content-Type", "application/json; charset=utf-8")
	r.ParseForm()
	
	if sn := GetSession(r); sn != nil {
		sn.Vars["user_name"] = r.Form.Get("username")
		sn.Vars["password"] = r.Form.Get("password")
		sn.Vars["login_attempts"] = sn.Vars["login_attempts"].(int) + 1
		if sn.Authenticated = OnLogin(r, sn); sn.Authenticated {
			// http.Redirect(w, r, sn.Home, http.StatusMovedPermanently)
			u, _ := url.Parse(r.Referer())
			
			// if u.Path == "/app/login" {
			if u.Path == "/login" {
				http.Redirect(w, r, "/", http.StatusMovedPermanently)
			} else {
				http.Redirect(w, r, u.Path, http.StatusMovedPermanently)
			}
			// http.Redirect(w, r, "/", http.StatusMovedPermanently)
		} else {
			u, _ := url.Parse(r.Referer())
			
			// fmt.Println("Referer", u.Host, u.Path, u.Fragment, u.RawQuery)
			// http.Redirect(w, r, fmt.Sprintf("/app/login?error=%v", sn.Vars["login_attempts"]), http.StatusMovedPermanently)
			// http.Redirect(w, r, fmt.Sprintf("/app/login?error=%v", sn.Vars["login_attempts"]), http.StatusMovedPermanently)
			http.Redirect(w, r, fmt.Sprintf("%s?error=%v", u.Path, sn.Vars["login_attempts"]), http.StatusMovedPermanently)
		}
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	
	// r.ParseForm()
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	
	if sn := GetSession(r); sn != nil && sn.Authenticated {
		if OnLogout(r, sn) {
			delete(Sessions, sn.SessionId)
			kv := utils.NewKeyValue()
			kv.Set("Status", 0)
			w.Write([]byte(kv.Json()))
		}
	}
	// http.Redirect(w, r, "/", http.StatusMovedPermanently)
	
	// if sn := GetSession(r); sn != nil {
		
		// sn.Vars["user_name"] = r.Form.Get("username")
		// sn.Vars["password"] = r.Form.Get("password")
		// sn.Vars["login_attempts"] = sn.Vars["login_attempts"].(int) + 1
		
		// if sn.Authenticated = OnLogin(r, sn); sn.Authenticated {
			// u, _ := url.Parse(r.Referer())
			
			// if u.Path == "/login" {
				// http.Redirect(w, r, "/", http.StatusMovedPermanently)
			// } else {
				// http.Redirect(w, r, u.Path, http.StatusMovedPermanently)
			// }
		// } else {
			// u, _ := url.Parse(r.Referer())
			
			// http.Redirect(w, r, fmt.Sprintf("%s?error=%v", u.Path, sn.Vars["login_attempts"]), http.StatusMovedPermanently)
		// }
	// }
}
