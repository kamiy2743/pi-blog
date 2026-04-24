package session

type Flash struct {
	Success string
	Error   string
}

func (f Flash) IsEmpty() bool {
	return f.Success == "" && f.Error == ""
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
