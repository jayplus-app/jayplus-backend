package models

type BusinessHours struct {
	Monday    DayHours `json:"monday"`
	Tuesday   DayHours `json:"tuesday"`
	Wednesday DayHours `json:"wednesday"`
	Thursday  DayHours `json:"thursday"`
	Friday    DayHours `json:"friday"`
	Saturday  DayHours `json:"saturday"`
	Sunday    DayHours `json:"sunday"`
}

type DayHours struct {
	Start  string `json:"start"`
	End    string `json:"end"`
	Closed bool   `json:"closed"`
}
