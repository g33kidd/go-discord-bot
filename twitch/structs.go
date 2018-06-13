package twitch

// TODO: Add more fields to this...
type TwitchChannel struct {
	Status      string `json:"status"`
	Game        string `json:"game"`
	Followers   int    `json:"followers"`
	Views       int    `json:"views"`
	DisplayName string `json:"display_name"`
	ID          string `json:"_id"`
}

type TwitchChannelEditData struct {
	Status string `json:"status"`
	Game   string `json:"game"`
}

type TwitchChannelEdit struct {
	Channel *TwitchChannelEditData `json:"channel"`
}

type TwitchStream struct {
	Data *TwitchStreamData `json:"stream,omitempty"`
}

type TwitchStreamPreview struct {
	Small  string `json:"small"`
	Medium string `json:"medium"`
}

type TwitchStreamData struct {
	Game      string               `json:"game"`
	Viewers   int                  `json:"viewers"`
	Preview   *TwitchStreamPreview `json:"preview"`
	Channel   *TwitchChannel       `json:"channel"`
	CreatedAt string               `json:"created_at"`
	Delay     int                  `json:"delay"`
}

type TwitchUser struct {
	ID   string `json:"_id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type TwitchUsers struct {
	Total int           `json:"_total"`
	Users []*TwitchUser `json:"users"`
}
