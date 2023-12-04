package models

type UIConfig struct {
	PrimaryColorLight       string `json:"primaryColorLight"`
	PrimaryColorDark        string `json:"primaryColorDark"`
	SecondaryColorLight     string `json:"secondaryColorLight"`
	SecondaryColorDark      string `json:"secondaryColorDark"`
	SecondaryColorDarker    string `json:"secondaryColorDarker"`
	SecondaryColorDarkest   string `json:"secondaryColorDarkest"`
	ComplementaryColorLight string `json:"complementaryColorLight"`
	ComplementaryColorDark  string `json:"complementaryColorDark"`
	DisableColor            string `json:"disableColor"`
}
