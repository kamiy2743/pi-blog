package session

type Flash struct {
	Error   string
	Success string
}

func (f Flash) IsEmpty() bool {
	return f.Error == "" && f.Success == ""
}

func FlashToMap(flash *Flash) map[string]string {
	if flash == nil || flash.IsEmpty() {
		return map[string]string{}
	}

	flashMap := map[string]string{}
	if flash.Success != "" {
		flashMap["success"] = flash.Success
	}
	if flash.Error != "" {
		flashMap["error"] = flash.Error
	}
	return flashMap
}
