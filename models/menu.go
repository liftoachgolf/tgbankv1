package models

type MenuButton struct {
	MenuButtonCommands *MenuButtonCommands `json:"MenuButtonCommands"`
}
type MenuButtonCommands struct {
	Type string `json:"type"`
}

type MenuButtonDefault struct {
	Type string `json:"type"`
}
type MenuButtonWebApp struct {
	Type   string     `json:"type"`
	Text   string     `json:"text"`
	WebApp WebAppInfo `json:"web_app"`
}

type WebAppInfo struct {
	Url string `json:"url"`
}
