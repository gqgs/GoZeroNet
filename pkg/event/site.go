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
