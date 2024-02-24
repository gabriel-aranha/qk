# Quake 3 Log Parser

This project reads Quake 3 game logs, parses all the relevant kill, user and games data, and output a report in the following JSON format:
```json
{
    "games": {
        "game_0": {
            "total_kills": 7,
            "players": [
                "Isgalamido",
                "Assasinu Credi",
                "Dono da Bola",
                "Oootsimo",
                "Assasinu Credi"
            ],
            "kills": {
                "Isgalamido": 2,
                "Dono da Bola": 1,
            },
            "kills_by_means": {
                "MOD_TRIGGER_HURT": 3,
                "MOD_ROCKET": 4
            }
        },
        ...
    }
}
```

## Dependencies  
```bash
Go 1.22
Git
```

## Installation

To install this project, you need to have Go 1.22 installed on your machine. If you don't have Go installed, you can download it from the [official Go website](https://golang.org/dl/).

Once you have Go installed, you can clone this repository to your local machine:
```bash
git clone https://github.com/gabriel-aranha/qk.git
cd qk
```

## Testing
To run the tests for this project, navigate to the project directory and run the following command:
```bash
go test ./...
```
This will run all tests in the project.

## Running the Project
By default the project will use the built-in file located in the following directory:
```bash
qk/input/games.log
```
If you want, just replace the file with your own log file.

To run the project, use the following command:
```bash
go run main.go
```
This will process the input games log file and output a report file to the following directory:
```bash
qk/output/report.json
```