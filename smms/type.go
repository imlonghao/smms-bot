package smms

type UploadJSON struct {
	Success bool   `json:"success"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    struct {
		FileID    int    `json:"file_id"`
		Width     int    `json:"width"`
		Height    int    `json:"height"`
		Filename  string `json:"filename"`
		Storename string `json:"storename"`
		Size      int    `json:"size"`
		Path      string `json:"path"`
		Hash      string `json:"hash"`
		URL       string `json:"url"`
		Delete    string `json:"delete"`
		Page      string `json:"page"`
	} `json:"data"`
	RequestID string `json:"RequestId"`
}

type HistoryJSON struct {
	Success bool   `json:"success"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    []struct {
		FileID    int    `json:"file_id"`
		Width     int    `json:"width"`
		Height    int    `json:"height"`
		Filename  string `json:"filename"`
		Storename string `json:"storename"`
		Size      int    `json:"size"`
		Path      string `json:"path"`
		Hash      string `json:"hash"`
		URL       string `json:"url"`
		Delete    string `json:"delete"`
		Page      string `json:"page"`
	} `json:"data"`
	RequestID string `json:"RequestId"`
}

type DeleteJSON struct {
	Success   bool   `json:"success"`
	Code      string `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"RequestId"`
}

type ClearJSON struct {
	Success   bool          `json:"success"`
	Code      string        `json:"code"`
	Message   string        `json:"message"`
	Data      []interface{} `json:"data"`
	RequestID string        `json:"RequestId"`
}
