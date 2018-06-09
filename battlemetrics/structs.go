package battlemetrics

// Server is a server..
type Server struct {
	Data struct {
		Type       string            `json:"type,omitempty"`
		ID         string            `json:"id,omitempty"`
		Attributes *ServerAttributes `json:"attributes,omitempty"`
	} `json:"data,omitempty"`
}

// ServerAttributes is the data, yo!
type ServerAttributes struct {
	ID         string         `json:"id,omitempty"`
	Name       string         `json:"name,omitempty"`
	IP         string         `json:"ip,omitempty"`
	Port       int            `json:"port,omitempty"`
	Players    int            `json:"players,omitempty"`
	MaxPlayers int            `json:"maxPlayers,omitempty"`
	Rank       int            `json:"rank,omitempty"`
	Status     string         `json:"status,omitempty"`
	CreatedAt  string         `json:"createdAt,omitempty"`
	UpdatedAt  string         `json:"updatedAt,omitempty"`
	PortQuery  int            `json:"portQuery,omitempty"`
	Country    string         `json:"country,omitempty"`
	Details    *ServerDetails `json:"details,omitempty"`
}

// ServerDetails are details, yo!
type ServerDetails struct {
	Map         string `json:"map,omitempty"`
	Environment string `json:"environment,omitempty"`
}
