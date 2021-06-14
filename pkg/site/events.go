package site

type SiteChangedEvent struct {
	Params *Info  `json:"params"`
	Cmd    string `json:"cmd"`
}
