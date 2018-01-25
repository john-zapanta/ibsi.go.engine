package template

import (
	"fmt"
	"math"
	"strconv"
	"path/filepath"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"html/template"
	"github.com/gorilla/mux"
	// "github.com/leekchan/accounting"
	"ibsi/utils"
)

var Router *mux.Router
var AppIcon string

type Page struct {
	Pid string
	Ver string
	Uid string
	Callback string
	Title string
	Icon string
	Nav *utils.Navigator
	Data string
	Custom interface{}
}

type Controller struct {
	Pid string
	Root string
	Template string
	Params string
	Format string
	ValueParam string
	OnInitHandlers InitHandlers
	OnInitPageData InitPageData
	OnInitTemplateCustomData InitTemplateCustomData
	OnInitCallbackData InitCallbackData
	OnInitCustomData InitCustomData
}

type InitTemplateCustomData func(*http.Request, *Page) (interface{})
type InitCallbackData func(nav *utils.Navigator)
type InitCustomData func(*http.Request, *utils.Navigator) (interface{})
type InitPageData func(*http.Request, *Page)
type InitHandlers func(*Controller)

func (s *Controller) Add(url string)  {
	Router.HandleFunc(url, s.Handler)
}

func (s Controller) Handler(w http.ResponseWriter, r *http.Request)  {
	s.LoadTemplate(w, r)
}

// thanks to: https://stackoverflow.com/questions/9320427/best-practice-for-embedding-arbitrary-json-in-the-dom
func NewController(ts Controller) {
	ts.OnInitHandlers(&ts)
}

func (ts Controller) LoadTemplate(w http.ResponseWriter, r *http.Request) {

	// fmt.Println("here 1...")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	
	// load the html file to memory
	// ex, err := os.Executable()
    // if err != nil {
        // panic(err)
    // }
    // exPath := filepath.Dir(ex)
	
	content, err := ioutil.ReadFile(filepath.Join("", "templates", fmt.Sprintf("%s.html", ts.Template)))
	// content, err := ioutil.ReadFile(fmt.Sprintf("%s.html", ts.Template))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	// create the template object
	// ac := accounting.Accounting{Symbol: "", Precision: 2}
	
	// fmt.Println(ts.Pid)
	t := template.New(ts.Pid).Funcs(template.FuncMap{
		// "Money": func(v float64) string {
			// return ac.FormatMoney(v)
		// }
		"Money": func(v float64) string {
			Part := func (v float64) (float64, string) {
				x := math.Mod(v, 1000)
				return (v - x) / 1000, strconv.FormatFloat(x, 'G', -1, 64)
			}
			
			f, money, s := v, "", ""

			for f != 0 {
				f, s = Part(math.Floor(f))
				money = s +","+ money
					
			}
			
			return money[0:len(money)-1]
		},
	})
	
	// parse the html content
	t, err = t.Parse(string(content[:]))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	
	page := Page {
		Pid: ts.Pid,
		Ver: "go",
		Uid: "CLkRHQZj+dTKcoB+rXWeag==",
		// Callback: callback,
		Icon: AppIcon,
		// Custom: ts.OnInitTemplateCustomData(r),
	}
	
	page.Nav = utils.NewNavigator(ts.Pid, func(nav *utils.Navigator) {
	})
	
	data, _ := json.MarshalIndent(page.Nav, "", "\t") // formatted
	page.Data = string(data[:])
	
	ts.OnInitPageData(r, &page)
	
	// merge the html content and custom data
	t.Execute(w, page)
}