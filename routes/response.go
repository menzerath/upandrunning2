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
	GoVersion  string
	AppVersion string
}
