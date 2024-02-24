package parser

import (
	"reflect"
	"testing"

	"github.com/gabriel-aranha/qk/internal/types"
)

func TestParse(t *testing.T) {
	p := NewParser(nil)

	tests := []struct {
		description   string
		arrayLines    []string
		expectedGames int
	}{
		{
			description: "one game",
			arrayLines: []string{
				"20:00 InitGame: \\sv_floodProtect\\1\\sv_maxPing\\0",
				"20:37 ClientUserinfoChanged: 3 n\\Dono da Bola\\t\\0",
				"20:40 ClientUserinfoChanged: 2 n\\Isgalamido\\t\\0",
				"20:34 Kill: 1022 2 22: Isgalamido killed Dono da Bola by MOD_TRIGGER_HURT",
			},
			expectedGames: 1,
		},
		{
			description: "two games",
			arrayLines: []string{
				"20:00 InitGame: \\sv_floodProtect\\1\\sv_maxPing\\0",
				"20:37 ClientUserinfoChanged: 3 n\\Dono da Bola\\t\\0",
				"20:40 ClientUserinfoChanged: 2 n\\Isgalamido\\t\\0",
				"20:34 Kill: 1022 2 22: Isgalamido killed Dono da Bola by MOD_TRIGGER_HURT",
				"20:00 InitGame: \\sv_floodProtect\\1\\sv_maxPing\\0",
				"20:37 ClientUserinfoChanged: 3 n\\Dono da Bola\\t\\0",
				"20:40 ClientUserinfoChanged: 2 n\\Isgalamido\\t\\0",
				"20:34 Kill: 1022 2 22: Isgalamido killed Dono da Bola by MOD_TRIGGER_HURT",
			},
			expectedGames: 2,
		},
	}

	for _, test := range tests {
		games, _ := p.Parse(test.arrayLines)
		countGames := len(games.Games)
		if countGames != test.expectedGames {
			t.Errorf("%s: Expected number of games %v, got %v", test.description, test.expectedGames, countGames)
		}
	}
}

func TestProcessNewGame(t *testing.T) {
	p := NewParser(nil)

	tests := []struct {
		description         string
		gameLines           []string
		expectedKills       map[string]int
		expectedPlayers     []string
		expectedKillByMeans map[string]int
		expectedTotalKills  int
	}{
		{
			description: "game with 1 kill",
			gameLines: []string{
				"20:37 ClientUserinfoChanged: 3 n\\Dono da Bola\\t\\0",
				"20:40 ClientUserinfoChanged: 2 n\\Isgalamido\\t\\0",
				"20:44 Kill: 1022 2 22: Isgalamido killed Dono da Bola by MOD_TRIGGER_HURT",
			},
			expectedKills: map[string]int{
				"Isgalamido": 1,
			},
			expectedPlayers: []string{
				"Dono da Bola",
				"Isgalamido",
			},
			expectedKillByMeans: map[string]int{
				"MOD_TRIGGER_HURT": 1,
			},
			expectedTotalKills: 1,
		},
		{
			description: "game with 1 kill by <world>",
			gameLines: []string{
				"20:37 ClientUserinfoChanged: 3 n\\Dono da Bola\\t\\0",
				"20:40 ClientUserinfoChanged: 2 n\\Isgalamido\\t\\0",
				"20:44 Kill: 1022 2 22: <world> killed Dono da Bola by MOD_TRIGGER_HURT",
			},
			expectedKills: map[string]int{},
			expectedPlayers: []string{
				"Dono da Bola",
				"Isgalamido",
			},
			expectedKillByMeans: map[string]int{
				"MOD_TRIGGER_HURT": 1,
			},
			expectedTotalKills: 1,
		},
		{
			description: "game with 1 kill by suicide",
			gameLines: []string{
				"20:37 ClientUserinfoChanged: 3 n\\Dono da Bola\\t\\0",
				"20:40 ClientUserinfoChanged: 2 n\\Isgalamido\\t\\0",
				"20:44 Kill: 1022 2 22: Dono da Bola killed Dono da Bola by MOD_ROCKET_SPLASH",
			},
			expectedKills: map[string]int{},
			expectedPlayers: []string{
				"Dono da Bola",
				"Isgalamido",
			},
			expectedKillByMeans: map[string]int{
				"MOD_ROCKET_SPLASH": 1,
			},
			expectedTotalKills: 1,
		},
		{
			description: "game with 1 kill and username change",
			gameLines: []string{
				"20:37 ClientUserinfoChanged: 3 n\\Dono da Bola\\t\\0",
				"20:40 ClientUserinfoChanged: 2 n\\Isgalamido\\t\\0",
				"20:44 Kill: 1022 2 22: Isgalamido killed Dono da Bola by MOD_ROCKET_SPLASH",
				"20:42 ClientUserinfoChanged: 2 n\\NewIsgalamido\\t\\0",
			},
			expectedKills: map[string]int{
				"NewIsgalamido": 1,
			},
			expectedPlayers: []string{
				"Dono da Bola",
				"NewIsgalamido",
			},
			expectedKillByMeans: map[string]int{
				"MOD_ROCKET_SPLASH": 1,
			},
			expectedTotalKills: 1,
		},
	}

	for _, test := range tests {
		game, _ := p.processNewGame(1, test.gameLines)
		if !reflect.DeepEqual(game.Kills, test.expectedKills) {
			t.Errorf("%s: Expected kills %v, got %v", test.description, test.expectedKills, game.Kills)
		}
		if !reflect.DeepEqual(game.Players, test.expectedPlayers) {
			t.Errorf("%s: Expected players %v, got %v", test.description, test.expectedPlayers, game.Players)
		}
		if !reflect.DeepEqual(game.KillsByMeans, test.expectedKillByMeans) {
			t.Errorf("%s: Expected kill by means %v, got %v", test.description, test.expectedKillByMeans, game.KillsByMeans)
		}
		if game.TotalKills != test.expectedTotalKills {
			t.Errorf("%s: Expected total kills %v, got %v", test.description, test.expectedTotalKills, game.TotalKills)
		}
	}
}

func TestProcessUserInfoLine(t *testing.T) {
	p := NewParser(nil)

	tests := []struct {
		description      string
		line             string
		previousUsername string
		expectedPlayerID string
		expectedUsername string
	}{
		{
			description:      "new user",
			line:             "20:34 ClientUserinfoChanged: 2 n\\Isgalamido\\t\\0",
			previousUsername: "",
			expectedPlayerID: "2",
			expectedUsername: "Isgalamido",
		},
		{
			description:      "username change",
			line:             "20:37 ClientUserinfoChanged: 3 n\\Dono da Bola\\t\\0",
			previousUsername: "Isgalamido",
			expectedPlayerID: "3",
			expectedUsername: "Dono da Bola",
		},
	}

	for _, test := range tests {
		game := p.newGame()
		if test.previousUsername != "" {
			player := types.Player{CurrentUsername: test.previousUsername, UserID: test.expectedPlayerID}
			game.PlayerList = append(game.PlayerList, player)
		}
		game, err := p.processUserInfoLine(test.line, game)
		if err != nil {
			t.Errorf("%s: Unexpected error: %v", test.description, err)
		}
		if test.previousUsername != "" {
			if game.PlayerList[0].PreviousUsernames[0] != test.previousUsername {
				t.Errorf("%s: Expected previous username %v, got %v", test.description, test.previousUsername, game.PlayerList[0].PreviousUsernames[0])
			}
		}
	}
}

func TestProcessKillLine(t *testing.T) {
	p := NewParser(nil)

	tests := []struct {
		description        string
		line               string
		expectedTotalKills int
	}{
		{
			description:        "kill by <world>",
			line:               "20:34 Kill: 1022 2 22: <world> killed Isgalamido by MOD_TRIGGER_HURT",
			expectedTotalKills: 1,
		},
		{
			description:        "kill by other player",
			line:               "20:34 Kill: 2 3 10: Isgalamido killed Dono da Bola by MOD_RAILGUN",
			expectedTotalKills: 1,
		},
		{
			description:        "kill by self damage",
			line:               "20:34 Kill: 2 3 10: Isgalamido killed Isgalamido by MOD_RAILGUN",
			expectedTotalKills: 1,
		},
	}

	for _, test := range tests {
		game := p.newGame()
		game, err := p.processKillLine(test.line, game)
		if err != nil {
			t.Errorf("%s: Unexpected error: %v", test.description, err)
		}
		if game.TotalKills != test.expectedTotalKills {
			t.Errorf("%s: Expected total kills %v, got %v", test.description, test.expectedTotalKills, game.TotalKills)
		}
	}
}

func TestExtractUserDetails(t *testing.T) {
	p := NewParser(nil)
	tests := []struct {
		description  string
		line         string
		expectedID   string
		expectedName string
	}{
		{
			description:  "user and single digit id check",
			line:         "20:34 ClientUserinfoChanged: 2 n\\Isgalamido\\t\\0",
			expectedID:   "2",
			expectedName: "Isgalamido",
		},
		{
			description:  "user and double digits id check",
			line:         "20:34 ClientUserinfoChanged: 16 n\\Isgalamido\\t\\0",
			expectedID:   "16",
			expectedName: "Isgalamido",
		},
	}

	for _, test := range tests {
		id, name, err := p.extractUserDetails(test.line)
		if err != nil {
			t.Errorf("%s: Unexpected error: %v", test.description, err)
		}
		if id != test.expectedID {
			t.Errorf("%s: Expected ID %v, got %v", test.description, test.expectedID, id)
		}
		if name != test.expectedName {
			t.Errorf("%s: Expected name %v, got %v", test.description, test.expectedName, name)
		}
	}
}

func TestExtractKillDetails(t *testing.T) {
	p := NewParser(nil)
	tests := []struct {
		description    string
		line           string
		expectedKiller string
		expectedKilled string
		expectedMeans  string
	}{
		{
			description:    "kill by <world>",
			line:           "20:34 Kill: 1022 2 22: <world> killed Isgalamido by MOD_TRIGGER_HURT",
			expectedKiller: "<world>",
			expectedKilled: "Isgalamido",
			expectedMeans:  "MOD_TRIGGER_HURT",
		},
		{
			description:    "kill by self damage",
			line:           "20:34 Kill: 1022 2 22: Isgalamido killed Isgalamido by MOD_ROCKET_SPLASH",
			expectedKiller: "Isgalamido",
			expectedKilled: "Isgalamido",
			expectedMeans:  "MOD_ROCKET_SPLASH",
		},
		{
			description:    "kill by other player",
			line:           "20:34 Kill: 1022 2 22: Mocinha killed Isgalamido by MOD_ROCKET",
			expectedKiller: "Mocinha",
			expectedKilled: "Isgalamido",
			expectedMeans:  "MOD_ROCKET",
		},
	}

	for _, test := range tests {
		killer, killed, means, err := p.extractKillDetails(test.line)
		if err != nil {
			t.Errorf("%s: Unexpected error: %v", test.description, err)
		}
		if killer != test.expectedKiller {
			t.Errorf("%s: Expected killer %v, got %v", test.description, test.expectedKiller, killer)
		}
		if killed != test.expectedKilled {
			t.Errorf("%s: Expected killed %v, got %v", test.description, test.expectedKilled, killed)
		}
		if means != test.expectedMeans {
			t.Errorf("%s: Expected means %v, got %v", test.description, test.expectedMeans, means)
		}
	}
}

func TestIsInitGameLine(t *testing.T) {
	p := NewParser(nil)
	tests := []struct {
		description string
		line        string
		expected    bool
	}{
		{
			description: "init game line",
			line:        "20:34 InitGame: \\sv_floodProtect\\1\\sv_maxPing\\0",
			expected:    true,
		},
		{
			description: "kill line",
			line:        "20:34 Kill: 1022 2 22: <world> killed Isgalamido by MOD_TRIGGER_HURT",
			expected:    false,
		},
		{
			description: "client user info changed line",
			line:        "20:34 ClientUserinfoChanged: 2 n\\Isgalamido\\t\\0",
			expected:    false,
		},
	}

	for _, test := range tests {
		result := p.isInitGameLine(test.line)
		if result != test.expected {
			t.Errorf("%s: Expected %v, got %v", test.description, test.expected, result)
		}
	}
}

func TestIsKillLine(t *testing.T) {
	p := NewParser(nil)
	tests := []struct {
		description string
		line        string
		expected    bool
	}{
		{
			description: "init game line",
			line:        "20:34 InitGame: \\sv_floodProtect\\1\\sv_maxPing\\0",
			expected:    false,
		},
		{
			description: "kill line",
			line:        "20:34 Kill: 1022 2 22: <world> killed Isgalamido by MOD_TRIGGER_HURT",
			expected:    true,
		},
		{
			description: "client user info changed line",
			line:        "20:34 ClientUserinfoChanged: 2 n\\Isgalamido\\t\\0",
			expected:    false,
		},
	}

	for _, test := range tests {
		result := p.isKillLine(test.line)
		if result != test.expected {
			t.Errorf("%s: Expected %v, got %v", test.description, test.expected, result)
		}
	}
}

func TestIsUserInfoLine(t *testing.T) {
	p := NewParser(nil)
	tests := []struct {
		description string
		line        string
		expected    bool
	}{
		{
			description: "init game line",
			line:        "20:34 InitGame: \\sv_floodProtect\\1\\sv_maxPing\\0",
			expected:    false,
		},
		{
			description: "kill line",
			line:        "20:34 Kill: 1022 2 22: <world> killed Isgalamido by MOD_TRIGGER_HURT",
			expected:    false,
		},
		{
			description: "client user info changed line",
			line:        "20:34 ClientUserinfoChanged: 2 n\\Isgalamido\\t\\0",
			expected:    true,
		},
	}

	for _, test := range tests {
		result := p.isUserInfoLine(test.line)
		if result != test.expected {
			t.Errorf("%s: Expected %v, got %v", test.description, test.expected, result)
		}
	}
}

func TestFormatGameNumber(t *testing.T) {
	p := &Parser{}

	tests := []struct {
		description string
		gameNumber  int
		expected    string
	}{
		{
			description: "game number 1",
			gameNumber:  1,
			expected:    "game_1",
		},
		{
			description: "game number 10",
			gameNumber:  10,
			expected:    "game_10",
		},
	}

	for _, test := range tests {
		result := p.formatGameNumber(test.gameNumber)
		if result != test.expected {
			t.Errorf("%s: Expected %s, got %s", test.description, test.expected, result)
		}
	}
}
