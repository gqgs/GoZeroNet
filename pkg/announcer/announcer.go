package announcer

type Stats struct {
	Status        string  `json:"status"`
	NumRequest    int     `json:"num_request"`
	NumSuccess    int     `json:"num_success"`
	NumError      int     `json:"num_error"`
	TimeRequest   float64 `json:"time_request"`
	TimeLastError float64 `json:"time_last_error"`
	TimeStatus    float64 `json:"time_status"`
	LastError     string  `json:"last_error"`
}

func GetStats() map[string]Stats {
	return map[string]Stats{
		"http://h4.trakx.nibba.trade:80/announce": {
			Status:      "announced",
			NumRequest:  10,
			NumSuccess:  10,
			NumError:    0,
			TimeRequest: 1623398701.6069171,
			TimeStatus:  1623398702.3022082,
		},
		"udp://104.238.198.186:8000": {
			Status:        "error",
			LastError:     "could not connect",
			NumError:      7,
			NumRequest:    10,
			NumSuccess:    2,
			TimeLastError: 1623398710.6333466,
			TimeRequest:   1623398701.6066,
			TimeStatus:    1623398710.6333451,
		},
	}
}
