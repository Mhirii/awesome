package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml"
)

var configDir = filepath.Join(os.Getenv("HOME"), ".config/awesome/rice")
var dataFile = filepath.Join(configDir, "singleData.toml")

type singleData struct {
	AvailableWallpapers []string `toml:"availableWallpapers"`
	CurrentWallpaper    string   `toml:"currentWallpaper"`
}

func main() {
	data := &singleData{}

	readData(data)

	fmt.Printf("Currently set wallpaper: %s\n", data.CurrentWallpaper)

	updateAvailableWallpapersList(data)

	cycleWallpapers(data)

	saveData(data)
}

// Function to read data from the TOML file
func readData(data *singleData) {
	file, err := ioutil.ReadFile(dataFile)
	if err != nil {
		fmt.Println("Error reading data file:", err)
		return
	}

	if err := toml.Unmarshal(file, data); err != nil {
		fmt.Println("Error decoding TOML:", err)
		return
	}
}

// Function to save data to the TOML file
func saveData(data *singleData) {
	dataBytes, err := toml.Marshal(data)
	if err != nil {
		fmt.Println("Error encoding TOML:", err)
		return
	}

	if err := ioutil.WriteFile(dataFile, dataBytes, 0644); err != nil {
		fmt.Println("Error writing data file:", err)
	}
}

// Function to set a new wallpaper
func setWallpaper(data *singleData, wallpaperName string) {
	wallpaperDir := filepath.Join(configDir, "wallpapers/single", wallpaperName)
	wallpaperPath := filepath.Join(wallpaperDir, "wallpaper.png")

	cmd := exec.Command("feh", "--bg-fill", wallpaperPath)
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error setting wallpaper:", err)
		return
	}

	data.CurrentWallpaper = wallpaperName
	fmt.Printf("Set wallpaper to: %s\n", data.CurrentWallpaper)
}

// Function to get the currently set wallpaper
func getCurrentWallpaper(data *singleData) {
	output, err := exec.Command("feh", "--bg-list").Output()
	if err != nil {
		fmt.Println("Error getting current wallpaper:", err)
		return
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) > 0 {
		data.CurrentWallpaper = strings.TrimSpace(lines[0])
	}
}

// Function to update the list of available wallpapers
func updateAvailableWallpapersList(data *singleData) {
	availableWallpapers := make([]string, 0)

	files, err := ioutil.ReadDir(filepath.Join(configDir, "wallpapers/single"))
	if err != nil {
		fmt.Println("Error reading wallpaper directory:", err)
		return
	}

	for _, file := range files {
		if file.IsDir() {
			availableWallpapers = append(availableWallpapers, file.Name())
		}
	}

	data.AvailableWallpapers = availableWallpapers
}

// Function to cycle through wallpapers
func cycleWallpapers(data *singleData) {
	if len(data.AvailableWallpapers) == 0 {
		fmt.Println("No wallpapers available.")
		return
	}

	currentIdx := -1
	for i, wp := range data.AvailableWallpapers {
		if wp == data.CurrentWallpaper {
			currentIdx = i
			break
		}
	}

	if currentIdx == -1 {
		// Current wallpaper not found, start with the first one
		currentIdx = 0
	} else {
		// Cycle to the next wallpaper
		currentIdx = (currentIdx + 1) % len(data.AvailableWallpapers)
	}

	nextWallpaper := data.AvailableWallpapers[currentIdx]
	setWallpaper(data, nextWallpaper)
}
