package writer

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"

	"github.com/gabriel-aranha/qk/internal/types"
)

func TestWrite(t *testing.T) {
	w := NewWriter(nil)

	tests := []struct {
		description string
		games       types.Games
	}{
		{
			description: "full game report",
			games: types.Games{
				Games: map[string]types.Game{
					"game_1": {
						TotalKills: 10,
						Players:    []string{"player_1", "player_2"},
						Kills: map[string]int{
							"player_1": 5,
							"player_2": 5,
						},
						KillsByMeans: map[string]int{
							"MOD_TEST": 10,
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		err := w.Write(test.games)
		if err != nil {
			t.Errorf("%s: Error writing to output file: %v", test.description, err)
		}

		file, err := os.Open("output/report.json")
		if err != nil {
			t.Errorf("%s: Error opening output file: %v", test.description, err)
		}
		defer os.Remove("output/report.json")

		var games types.Games
		decoder := json.NewDecoder(file)
		err = decoder.Decode(&games)
		if err != nil {
			t.Errorf("%s: Error decoding output file: %v", test.description, err)
		}

		if !reflect.DeepEqual(games, test.games) {
			t.Errorf("%s: Expected %v, got %v", test.description, test.games, games)
		}

		file.Close()
	}
}
