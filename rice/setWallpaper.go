package main

import (
	"fmt"
	"github.com/pelletier/go-toml"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var configDir = filepath.Join(os.Getenv("HOME"), ".config/awesome/rice")
var dataFile string
var wallpaperDir string

type WallpaperData struct {
	AvailableWallpapers []string `toml:"availableWallpapers"`
	CurrentWallpaper    string   `toml:"currentWallpaper"`
}

func main() {
	if len(os.Args) != 1 {
		fmt.Println("Usage: setWallpaper")
		return
	}

	numMonitors, err := getNumberOfMonitors()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var mode string
	if numMonitors == 1 {
		mode = "single"
		dataFile = filepath.Join(configDir, "singleData.toml")
		wallpaperDir = filepath.Join(configDir, "wallpapers/single")
	} else {
		mode = "dual"
		dataFile = filepath.Join(configDir, "dualData.toml")
		wallpaperDir = filepath.Join(configDir, "wallpapers/dual")
	}

	// mode := os.Args[1]

	// switch mode {
	// case "dual":
	// 	dataFile = filepath.Join(configDir, "dualData.toml")
	// 	wallpaperDir = filepath.Join(configDir, "wallpapers/dual")
	// case "single":
	// 	dataFile = filepath.Join(configDir, "singleData.toml")
	// 	wallpaperDir = filepath.Join(configDir, "wallpapers/single")
	// default:
	// 	fmt.Println("Usage: setWallpaper [dual|single]")
	// 	return
	// }

	// Initialize data structure
	data := &WallpaperData{}

	// Read data from TOML file
	readData(data)

	// Get and print the currently set wallpaper
	fmt.Printf("Currently set wallpaper: %s\n", data.CurrentWallpaper)

	// Update available wallpapers list
	updateAvailableWallpapersList(data)

	// Cycle through wallpapers and set wallpaper
	cycleWallpapers(data, mode)

	setColorsRasi(data.CurrentWallpaper, filepath.Join(os.Getenv("HOME"), ".config", "awesome", "rice", "rofi"))
	// Save updated data to TOML file
	saveData(data)
}

// Function to read data from the TOML file
func readData(data *WallpaperData) {
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
func saveData(data *WallpaperData) {
	dataBytes, err := toml.Marshal(data)
	if err != nil {
		fmt.Println("Error encoding TOML:", err)
		return
	}

	if err := ioutil.WriteFile(dataFile, dataBytes, 0644); err != nil {
		fmt.Println("Error writing data file:", err)
	}
}

func getNumberOfMonitors() (int, error) {
	cmd := exec.Command("xrandr", "--listmonitors")
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	// Split the output into lines
	lines := strings.Split(string(output), "\n")

	// The number of monitors is the number of lines excluding the header line
	numMonitors := len(lines) - 2

	return numMonitors, nil
}

// Function to set a new wallpaper
func setWallpaper(data *WallpaperData, wallpaperName, mode string) {
	wallpaperPath := filepath.Join(wallpaperDir, wallpaperName, "wallpaper.png")
	var cmd *exec.Cmd

	// cmd := exec.Command("feh", "--bg-fill", wallpaperPath)
	if mode == "dual" {
		cmd = exec.Command("feh", "--bg-fill", "--no-xinerama", wallpaperPath) // Replace "dual-wallpaper-command" with the actual command you want to run for "dual" mode.
	} else {
		cmd = exec.Command("feh", "--bg-fill", wallpaperPath) // Default command for other modes
	}
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error setting wallpaper:", err)
		return
	}

	data.CurrentWallpaper = wallpaperName
	fmt.Printf("Set wallpaper to: %s\n", data.CurrentWallpaper)
}

// Function to copy file colors.rasi from wallpaper directory to the specified destination directory
func setColorsRasi(wallpaperName, destDir string) error {
	colorsRasiSourcePath := filepath.Join(wallpaperDir, wallpaperName, "colors.rasi")
	colorsRasiDestPath := filepath.Join(destDir, "colors.rasi")

	data, err := ioutil.ReadFile(colorsRasiSourcePath)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(colorsRasiDestPath, data, 0644); err != nil {
		return err
	}

	fmt.Printf("Copied colors.rasi to %s\n", colorsRasiDestPath)
	return nil
}

// Function to update the list of available wallpapers
func updateAvailableWallpapersList(data *WallpaperData) {
	availableWallpapers := make([]string, 0)

	files, err := ioutil.ReadDir(wallpaperDir)
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
func cycleWallpapers(data *WallpaperData, mode string) {
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
	setWallpaper(data, nextWallpaper, mode)
}
