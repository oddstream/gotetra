// Copyright ©️ 2020 oddstream.games

package tetra

import (
	"encoding/json"
	"log"
	"os"
	"path"
)

// UserData contains the level the user is on
type UserData struct {
	Copyright string
	Level     int
}

func fullPath() string {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatal(err)
	}
	// println("UserConfigDir", userConfigDir) // /home/gilbert/.config
	return path.Join(userConfigDir, "oddstream.games", "tetra", "userdata.json")
}

func makeConfigDir() {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatal(err)
	}

	dir := path.Join(userConfigDir, "oddstream.games", "tetra")
	err = os.MkdirAll(dir, 0755) // https://stackoverflow.com/questions/14249467/os-mkdir-and-os-mkdirall-permission-value
	if err != nil {
		log.Fatal(err)
	}
	// if path is already a directory, MkdirAll does nothing and returns nil
}

// NewUserData create a new UserData object and tries to load it's content from file
// it always returns an object, even if file does not exist
func NewUserData() *UserData {
	ud := &UserData{Copyright: "Copyright ©️ 2020 oddstream.games", Level: 1}

	file, err := os.Open(fullPath())
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

	makeConfigDir()

	file, err := os.Create(fullPath())
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = file.Write(bytes)
	if err != nil {
		log.Fatal(err)
	}
}
