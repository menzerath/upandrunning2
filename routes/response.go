package routes

type BasicResponse struct {
	Success bool   `json:"requestSuccess"`
	Message string `json:"message"`
}

type SiteData struct {
	Title string
}

type AdminSiteData struct {
	Title      string
	Interval   int
	PbKey      string
	AppVersion string
	GoVersion  string
	GoArch     string
}
