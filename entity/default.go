package entity

type StaticConfiguration struct {
	StaticPort string
	StaticPath string
	FilePath   string
}

type File struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	Size    int64  `json:"size"`
	ModTime int64  `json:"modTime"`
}

type ResponseData struct {
	StatusCode int64  `json:"statusCode"`
	Msg        string `json:"msg"`
	Data       []File `json:"data"`
}