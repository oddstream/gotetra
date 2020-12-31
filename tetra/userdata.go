// Copyright ©️ 2020 oddstream.games

package tetra

import (
	"encoding/json"
	"log"
	"os"
)

const fname string = "tetra.json"

// UserData contains the level the user is on
type UserData struct {
	Level int
}

// NewUserData create a new UserData object and tries to load it's content from file
// it always returns an object, even if file does not exist
func NewUserData() *UserData {
	ud := &UserData{Level: 1}

	file, err := os.Open(fname)
	if err == nil && file != nil {
		defer file.Close()

		bytes := make([]byte, 256)
		var count int
		count, err = file.Read(bytes)
		if err != nil {
			log.Fatal(err)
		}
		if count > 0 {
			// golang gotcha reslice buffer to number of bytes actually read
			err = json.Unmarshal(bytes[:count], ud)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	return ud
}

// Save writes the UserData object to file
func (ud *UserData) Save() {
	bytes, err := json.Marshal(ud)
	if err != nil {
		log.Fatal(err)
	}
	file, err := os.Create(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	bytesWritten, err := file.Write(bytes)
	if err != nil {
		log.Fatal(err)
	}
	println("Wrote bytes", bytesWritten)
}
