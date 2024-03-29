package board

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sea-battle/internal/ip"
	"sea-battle/internal/stats"
	"strconv"
	"time"

	"sea-battle/internal/boats"
	"sea-battle/internal/utils"
)

/*
	Overview of an empty sea battle board:

		A   B   C   D   E   F   G   H   I    J
	   -----------------------------------------
	01 |   |   |   |   |   |   |   |   |   |   |
	   -----------------------------------------
	02 |   |   |   |   |   |   |   |   |   |   |
	   -----------------------------------------
	03 |   |   |   |   |   |   |   |   |   |   |
	   -----------------------------------------
	04 |   |   |   |   |   |   |   |   |   |   |
	   -----------------------------------------
	05 |   |   |   |   |   |   |   |   |   |   |
	   -----------------------------------------
	06 |   |   |   |   |   |   |   |   |   |   |
	   -----------------------------------------
	07 |   |   |   |   |   |   |   |   |   |   |
	   -----------------------------------------
	08 |   |   |   |   |   |   |   |   |   |   |
	   -----------------------------------------
	09 |   |   |   |   |   |   |   |   |   |   |
	   -----------------------------------------
	10 |   |   |   |   |   |   |   |   |   |   |
	   -----------------------------------------
*/

type Shot struct {
	Position utils.Position
	Hit      bool
}

var BoatsBoard [5]boats.Boat

var BoatsDestroyedMap map[int]bool

var AllShots []Shot

func GetBoatAt(position utils.Position) boats.Boat {
	for _, boat := range BoatsBoard {
		for _, pos := range boat.Position {
			if pos.X == position.X && pos.Y == position.Y {
				return boat
			}
		}
	}
	panic("POSITION DOES NOT CORRESPOND TO A BOAT")
}

/*
		Prints an empty board for demonstration purposes (eg: tutorial)

	 	IMPORTANT: if user's terminal is less wide than 44 cols, the board will not
		be printed correctly
*/
func PrintEmptyBoard() {
	fmt.Println("\n     A   B   C   D   E   F   G   H   I   J")

	for i := 1; i <= 10; i++ {
		fmt.Println("   -----------------------------------------")
		fmt.Printf("%02d |   |   |   |   |   |   |   |   |   |   |\n", i)
	}

	fmt.Printf("   -----------------------------------------\n\n")
}

/*
Prints a board with shots & boats

IMPORTANT: if user's terminal is less wide than 44 cols, the board will not
be printed correctly
*/
func PrintBoard(boats [5]boats.Boat, isEnemyBoard bool, challengeSentence string) string {
	var result bytes.Buffer
	if isEnemyBoard {
		result.WriteString("\n     A   B   C   D   E   F   G   H   I   J \n")
	} else {
		fmt.Println("\n     A   B   C   D   E   F   G   H   I   J")
	}

	// Get all alive & destroyed boats positions
	var aliveBoatsPositions []utils.Position
	var destroyedBoatsPositions []utils.Position
	for _, boat := range boats {
		if BoatsDestroyedMap[boat.Id] {
			destroyedBoatsPositions = append(destroyedBoatsPositions, boat.Position...)
		} else {
			aliveBoatsPositions = append(aliveBoatsPositions, boat.Position...)
		}
	}

	for i := 1; i <= 10; i++ {
		if isEnemyBoard {
			result.WriteString("   ----------------------------------------- \n")
		} else {
			fmt.Println("   -----------------------------------------")
		}
		for j := 0; j <= 10; j++ {
			if j == 0 {
				if isEnemyBoard {
					result.WriteString(fmt.Sprintf("%02d |", i))
				} else {
					fmt.Printf("%02d |", i)
				}
			} else {
				/*
					Symbols:
					■ -> boat
					O -> missed shot
					X -> hit shot
					# -> destroyed boat
				*/

				symbol := " "

				if !isEnemyBoard {
					// Check if there is a boat alive at this position
					for _, boatPosition := range aliveBoatsPositions {
						if boatPosition.X == uint8(j) && boatPosition.Y == uint8(i) {
							symbol = "■"
						}
					}
				}

				// Check if there is a shot at this position
				for _, shot := range AllShots {
					if shot.Hit && shot.Position.X == uint8(j) && shot.Position.Y == uint8(i) {
						symbol = "X"
					} else if shot.Position.X == uint8(j) && shot.Position.Y == uint8(i) {
						symbol = "O"
					}
				}

				// Check if there is a destroyed boat at this position
				for _, boatPosition := range destroyedBoatsPositions {
					if boatPosition.X == uint8(j) && boatPosition.Y == uint8(i) {
						symbol = "#"
					}
				}
				if isEnemyBoard {
					result.WriteString(fmt.Sprintf(" %s |", symbol))
				} else {
					fmt.Printf(" %s |", symbol)
				}
			}
		}
		if isEnemyBoard {
			result.WriteString("\n")
		} else {
			fmt.Println()
		}
	}

	if isEnemyBoard {
		result.WriteString("   -----------------------------------------\n")
		result.WriteString("Message de l'adversaire :\n" + challengeSentence)
		return result.String()
	} else {
		fmt.Printf("   -----------------------------------------\n\n")
		return ""
	}
}

func RequestBoard(clientIP ip.IP) {
	port := strconv.Itoa(int(clientIP.Port))
	url := "http://" + clientIP.Ip + ":" + port + "/board"

	client := http.Client{
		Timeout: 2 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println("Une erreur est survenue.")
		return
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Une erreur est survenue.")
		return
	}
	result := string(body)
	fmt.Println(result)
}

// This function get a string in parameter (ex: "J6") and return a Position struct
func GetPositionFromString(inputPos string) utils.Position {
	var pos utils.Position
	YtoInt, _ := strconv.Atoi(inputPos[1:])
	pos.Y = uint8(YtoInt)

	mapOfCord := map[string]byte{"A": 1, "B": 2, "C": 3, "D": 4, "E": 5, "F": 6, "G": 7, "H": 8, "I": 9, "J": 10, "a": 1, "b": 2, "c": 3, "d": 4, "e": 5, "f": 6, "g": 7, "h": 8, "i": 9, "j": 10}
	pos.X = mapOfCord[inputPos[:1]]

	return pos
}

// Returns the number of alive boats
func GetAliveBoats() uint8 {
	var aliveBoats uint8

	for _, boat := range BoatsBoard {
		if !BoatsDestroyedMap[boat.Id] {
			aliveBoats++
		}
	}

	return aliveBoats
}

func InitBoatsBoard(bBoard [5]boats.Boat) {
	BoatsBoard = bBoard
	BoatsDestroyedMap = make(map[int]bool)
	for _, boat := range BoatsBoard {
		BoatsDestroyedMap[boat.Id] = false
	}
}

func GetBoatsBoard() [5]boats.Boat {
	return BoatsBoard
}

func AddShot(position utils.Position) bool {
	isShot := checkShot(position)

	actualShot := Shot{Position: position, Hit: isShot}

	if !alreadyShooted(position) {
		AllShots = append(AllShots, actualShot)
	}

	if isShot {
		checkDestroyed(GetBoatAt(position))
	}

	return actualShot.Hit
}

func alreadyShooted(position utils.Position) bool {
	for _, bol := range AllShots {
		if bol.Position.X == position.X && bol.Position.Y == position.Y {
			return true
		}
	}
	return false
}

func checkDestroyed(boat boats.Boat) {
	count := boat.Size
	for _, pos := range boat.Position {
		for _, shot := range AllShots {
			if pos.X == shot.Position.X && pos.Y == shot.Position.Y {
				count--
			}
		}
	}
	if count <= 0 {
		BoatsDestroyedMap[boat.Id] = true
	}
}

// Function to check if a shot is a hit or not and return a boolean
func checkShot(position utils.Position) bool {

	// Concatenate all boats' positions
	var allBoatsPositions []utils.Position
	for _, boat := range BoatsBoard {
		allBoatsPositions = append(allBoatsPositions, boat.Position...)
	}

	// Check if there is a boat at this position
	for _, boatPosition := range allBoatsPositions {
		if boatPosition.X == position.X && boatPosition.Y == position.Y {
			return true
		}
	}
	return false
}

func RequestHit(clientIP ip.IP, pos utils.Position) bool {

	port := strconv.Itoa(int(clientIP.Port))
	url := "http://" + clientIP.Ip + ":" + port + "/hit"

	jsonValue, _ := json.Marshal(pos)

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	request, err := client.Post(url, "application/json", bytes.NewBuffer(jsonValue))
	//set HTTP request header Content-Type (optional)
	//req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	if err != nil {
		//fmt.Println(err)
		fmt.Println("On dirait que votre adversaire est parti, tant pis !")
		return false
	}
	defer request.Body.Close()
	body, err := io.ReadAll(request.Body)

	if err != nil {
		fmt.Printf("Reading body failed: %s", err)
		return false
	}
	result := string(body)

	if result == "true\n" {
		fmt.Print("\nTouché ! 😎️ \n")
		stats.AddShotHit()

		// Request opponents's alive boats
		request, err := client.Get("http://" + clientIP.Ip + ":" + port + "/boats")
		if err != nil {
			panic(err)
		}
		defer request.Body.Close()

		body, err := io.ReadAll(request.Body)
		if err != nil {
			panic(err)
		}

		aliveBoatsInt, _ := strconv.Atoi(string(body))

		// Check if all boats are destroyed
		if aliveBoatsInt == 0 {
			// Notify player that he won
			fmt.Print("\nBravo, vous avez gagné ! 🎉\n")
			fmt.Print("Appuyez sur Entrée pour continuer...")
			fmt.Scanln()
			stats.AddGameWon()
		}
	} else {
		fmt.Print("\nRaté ! ☹️ \n")
		stats.AddShotMissed()
	}

	return true
}
