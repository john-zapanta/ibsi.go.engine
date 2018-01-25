package utils

import (
	"encoding/json"
)

type InitNavigator func(*Navigator) 
type InitNavigatorItem func(*NavigatorItem) 
type InitMenuItem func(*MenuItem) 
	
type Navigator struct {
	// Status int `json:"status"`
	// Message string `json:"message"`
	Pid string `json:"pid"`
	PageTitle string `json:"page_title"`
	WindowTitle string `json:"window_title"`
	Items []NavigatorItem `json:"menu_items"`
	CustomData interface{} `json:"custom_data"`
}

type NavigatorItem struct {
	ID string `json:"id"`
	Title string `json:"title"`
	Items []MenuItem `json:"subItems"`
}

type MenuItem struct {
	ID string `json:"id"`
	Title string `json:"title"`
	Description string `json:"description"`
	Icon string `json:"icon"`
	Url string `json:"url"`
	Content string `json:"content"`
	Scripts string `json:"scripts"`
	Css string `json:"css"`
	Run string `json:"run"`
	Enabled bool `json:"enabled"`
	UseFrame bool `json:"frame"`
	Params map[string]interface{} `json:"params"`
	
	Action string `json:"-"`
}

func (c *Navigator) MarshalJSON() ([]byte, error) {
	type Alias Navigator
	return json.Marshal(&struct {
		*Alias
		// Commands []*Command `json:"commands"`
	}{
		Alias: (*Alias)(c),
		// Commands: c.commands,
	})
}

// func NewNavigator(pid string, data InitNavigatorCustomData, init InitNavigator) Navigator {
func NewNavigator(pid string, init InitNavigator) *Navigator {
    nav := &Navigator{
		// Status: 0, 
		// Message: "",
		Pid: pid,
		Items: make([]NavigatorItem, 0),
		CustomData: struct{}{},
	}
	
	init(nav)
	
	return nav
}

func (n *Navigator) Add(item NavigatorItem) {
	n.Items = append(n.Items, item)
}

func (n *NavigatorItem) Add(menu MenuItem) {
	n.Items = append(n.Items, menu)
}

func NewNavigatorItem(nav *Navigator, id string, title string, init InitNavigatorItem) *NavigatorItem {
    item := NavigatorItem{
		ID: id, 
		Title: title,
		Items: make([]MenuItem, 0),
	}
	
	init(&item)
	nav.Add(item)
	
	return &item
}

func NewMenuItem(ni *NavigatorItem, init InitMenuItem) *MenuItem {
	menu := MenuItem{
		ID: "",
		Title: "",
		Description: "",
		Icon: "",
		Url: "",
		Content: "",
		Scripts: "",
		Css: "",
		Run: "",
		Enabled: true,
		UseFrame: false,
		Params: make(map[string]interface{}),
	}
	
	init(&menu)
	
	if menu.Description == "" {
		menu.Description = menu.Title
	}
	
	ni.Add(menu)
	
	return &menu
}
