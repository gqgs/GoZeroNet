package event

type SiteChanged struct {
	Params interface{} `json:"params"`
	Cmd    string      `json:"cmd"`
}

func (e *SiteChanged) String() string {
	return "siteChanged"
}

func BroadcastSiteChanged(site string, broadcaster EventBroadcaster, siteChanged *SiteChanged) {
	broadcaster.Broadcast(site, siteChanged)
}

type SiteUpdate struct {
	InnerPath string `json:"inner_path"`
	Body      []byte `json:"body"`
}

func (e *SiteUpdate) String() string {
	return "SiteUpdate"
}

func BroadcastSiteUpdate(site string, broadcaster EventBroadcaster, siteUpdate *SiteUpdate) {
	broadcaster.Broadcast(site, siteUpdate)
}
