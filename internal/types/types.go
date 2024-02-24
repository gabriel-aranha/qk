package types

type Game struct {
	TotalKills   int            `json:"total_kills"`
	Players      []string       `json:"players"`
	Kills        map[string]int `json:"kills"`
	KillsByMeans map[string]int `json:"kills_by_means"`
	PlayerList   []Player       `json:"-"`
}

type Games struct {
	Games map[string]Game `json:"games"`
}

type Player struct {
	CurrentUsername   string   `json:"current_username"`
	UserID            string   `json:"user_id"`
	PreviousUsernames []string `json:"previous_usernames"`
	Kills             int      `json:"kills"`
}
