package parser

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gabriel-aranha/qk/internal/types"
	"go.uber.org/zap"
)

const (
	worldKiller = "<world>"
)

type Parser struct {
	logger *zap.Logger
}

func NewParser(logger *zap.Logger) Parser {
	var parser Parser
	parser.logger = logger

	return parser
}

func (p *Parser) formatGameNumber(gameNumber int) string {
	return fmt.Sprintf("game_%d", gameNumber)
}

func (p *Parser) Parse(arrayLines []string) (types.Games, error) {
	var startLineIndex int
	gameNumber := 0
	games := types.Games{Games: make(map[string]types.Game)}
	for i, line := range arrayLines {
		if p.isInitGameLine(line) {
			if i != 0 {
				game, err := p.processNewGame(gameNumber, arrayLines[startLineIndex:i])
				if err != nil {
					p.logger.Error("error processing new game", zap.Error(err))
					return games, err
				}
				games.Games[p.formatGameNumber(gameNumber)] = game
			}
			startLineIndex = i
			gameNumber++
		}
	}
	if startLineIndex != len(arrayLines) {
		game, err := p.processNewGame(gameNumber, arrayLines[startLineIndex:])
		if err != nil {
			p.logger.Error("error processing new game", zap.Error(err))
			return games, err
		}
		games.Games[p.formatGameNumber(gameNumber)] = game
	}
	return games, nil
}

func (p *Parser) newGame() types.Game {
	return types.Game{
		TotalKills:   0,
		Players:      []string{},
		PlayerList:   []types.Player{},
		Kills:        make(map[string]int),
		KillsByMeans: make(map[string]int),
	}
}

func (p *Parser) processNewGame(gameNumber int, gameLines []string) (types.Game, error) {
	game := p.newGame()

	for _, line := range gameLines {
		if p.isKillLine(line) {
			var err error
			game, err = p.processKillLine(line, game)
			if err != nil {
				p.logger.Error("error processing kill line", zap.Error(err))
				return game, err
			}
		} else if p.isUserInfoLine(line) {
			var err error
			game, err = p.processUserInfoLine(line, game)
			if err != nil {
				p.logger.Error("error processing client user info line", zap.Error(err))
				return game, err
			}
		}
	}

	// Add all players with kills to the Kills field
	for _, player := range game.PlayerList {
		if player.Kills > 0 {
			game.Kills[player.CurrentUsername] = player.Kills
		}
	}

	// Add all current usernames of all players to the Players field
	for _, player := range game.PlayerList {
		game.Players = append(game.Players, player.CurrentUsername)
	}

	return game, nil
}

func (p *Parser) processUserInfoLine(line string, game types.Game) (types.Game, error) {
	userID, currentUsername, err := p.extractUserDetails(line)
	if err != nil {
		p.logger.Error("error extracting client user info line", zap.Error(err))
		return game, err
	}

	newPlayer := types.Player{
		CurrentUsername: currentUsername,
		UserID:          userID,
	}

	// Check if player is already in the game
	for i, existingPlayer := range game.PlayerList {
		if existingPlayer.UserID == userID {
			// If the player has changed their username, update the player struct
			if existingPlayer.CurrentUsername != currentUsername {
				game.PlayerList[i].PreviousUsernames = append(game.PlayerList[i].PreviousUsernames, existingPlayer.CurrentUsername)
				game.PlayerList[i].CurrentUsername = currentUsername
			}
			return game, nil
		}
	}

	// Check if player with the same username exists in the game
	for i, existingPlayer := range game.PlayerList {
		if existingPlayer.CurrentUsername == currentUsername {
			// If the player has reconnected with a new userID, update the player struct
			game.PlayerList[i].UserID = userID
			return game, nil
		}
	}

	// If not, add player to the game
	game.PlayerList = append(game.PlayerList, newPlayer)

	return game, nil
}

func (p *Parser) extractUserDetails(line string) (userId string, username string, err error) {
	parts := strings.Split(line, "\\")
	if len(parts) < 2 {
		return "", "", fmt.Errorf("could not parse line: %s", line)
	}

	username = parts[1]

	r := regexp.MustCompile(`.*ClientUserinfoChanged: (.*\d+)`)
	matches := r.FindStringSubmatch(parts[0])
	if len(matches) < 2 {
		return "", "", fmt.Errorf("could not parse userId: %s", parts[0])
	}
	userId = strings.TrimSpace(matches[1])
	return userId, username, nil
}

func (p *Parser) processKillLine(line string, game types.Game) (types.Game, error) {
	killer, killed, means, err := p.extractKillDetails(line)
	if err != nil {
		p.logger.Error("error extracting kill line", zap.Error(err))
		return game, err
	}

	if killer == worldKiller || killer == killed {
		for i, player := range game.PlayerList {
			if player.CurrentUsername == killed && player.Kills > 0 {
				game.PlayerList[i].Kills--
				break
			}
		}
	} else {
		for i, player := range game.PlayerList {
			if player.CurrentUsername == killer {
				game.PlayerList[i].Kills++
				break
			}
		}
	}

	game.TotalKills++
	game.KillsByMeans[means]++

	return game, nil
}

func (p *Parser) extractKillDetails(line string) (killer, killed, means string, err error) {
	parts := strings.Split(line, " killed ")
	if len(parts) != 2 {
		return "", "", "", fmt.Errorf("could not parse line: %s", line)
	}

	r := regexp.MustCompile(`.*Kill: \d+ \d+ \d+: (.*)`)
	matches := r.FindStringSubmatch(parts[0])
	if len(matches) < 2 {
		return "", "", "", fmt.Errorf("could not parse killer: %s", parts[0])
	}
	killer = strings.TrimSpace(matches[1])

	parts = strings.Split(parts[1], " by ")
	if len(parts) != 2 {
		return "", "", "", fmt.Errorf("could not parse line: %s", line)
	}

	killed = strings.TrimSpace(parts[0])
	means = strings.TrimSpace(parts[1])

	return killer, killed, means, nil
}

func (p *Parser) isInitGameLine(line string) bool {
	pattern := `\d+:\d+ InitGame`
	r := regexp.MustCompile(pattern)
	matches := r.FindAllString(line, -1)
	return len(matches) > 0
}

func (p *Parser) isKillLine(line string) bool {
	pattern := `\d+:\d+ Kill`
	r := regexp.MustCompile(pattern)
	matches := r.FindAllString(line, -1)
	return len(matches) > 0
}

func (p *Parser) isUserInfoLine(line string) bool {
	pattern := `\d+:\d+ ClientUserinfoChanged`
	r := regexp.MustCompile(pattern)
	matches := r.FindAllString(line, -1)
	return len(matches) > 0
}
