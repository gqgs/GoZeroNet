package site

type SiteChangedEvent struct {
	Info
	Cmd    string   `json:"cmd"`
	Events []string `json:"events"`
}
