package main

import (
	"encoding/json"
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func createWallpaperDirectory(destDir, wallpaperName, setupType string) (string, error) {
	setupDir := "single" // Default to single setup
	if setupType == "dual" {
		setupDir = "dual"
	}

	destDir = filepath.Join(destDir, setupDir)
	destSubDir := filepath.Join(destDir, wallpaperName)

	if err := os.MkdirAll(destSubDir, os.ModePerm); err != nil {
		return "", err
	}

	return destSubDir, nil
}

func createColorsRasiFile(destSubDir string, dominantColors map[string]string) error {
	colorsRasiContent := fmt.Sprintf(`*{
	background: %s90;
	background-alt: %sff;
	text: %sff;
	text-dark: %sff;
	selected: %sff;
	active: %s90;
	urgent: #f7768eff;
}`, dominantColors["color1"], dominantColors["color1"], dominantColors["color2"], dominantColors["color1"], dominantColors["color3"], dominantColors["color3"])

	colorsRasiFilePath := filepath.Join(destSubDir, "colors.rasi")
	return ioutil.WriteFile(colorsRasiFilePath, []byte(colorsRasiContent), 0644)
}

func copyWallpaper(sourcePath, destFilePath string) error {
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(destFilePath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	return nil
}

func getDominantColors(wallpaperDir, destDir string) error {
	wallpaperPath := filepath.Join(wallpaperDir, "wallpaper.png")
	outputPath := filepath.Join(wallpaperDir, "dominantcolors.json")

	cmd := exec.Command("dominantcolors", wallpaperPath)
	output, err := cmd.Output()
	if err != nil {
		return err
	}

	colors := strings.Split(string(output), "\n")
	dominantColors := make(map[string]string)

	for i, color := range colors[:3] {
		key := fmt.Sprintf("color%d", i+1)
		dominantColors[key] = color
	}

	colorsJSON, err := json.MarshalIndent(dominantColors, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(outputPath, colorsJSON, 0644)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	if len(os.Args) < 2 || len(os.Args) > 3 {
		fmt.Println("Usage: go run wallpapers.go <source_picture_path> [single|dual]")
		os.Exit(1)
	}

	sourcePath := os.Args[1]
	setupType := "single" // Default to single setup

	if len(os.Args) == 3 {
		setupType = os.Args[2]
		if setupType != "single" && setupType != "dual" {
			fmt.Println("Invalid setup type. Please use 'single' or 'dual'.")
			os.Exit(1)
		}
	}

	destDir := filepath.Join(os.Getenv("HOME"), ".config", "awesome", "rice", "wallpapers")

	_, fileNameWithoutExt := filepath.Split(sourcePath)
	ext := filepath.Ext(fileNameWithoutExt)
	wallpaperName := fileNameWithoutExt[:len(fileNameWithoutExt)-len(ext)]

	if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
		fmt.Printf("Error creating destination directory: %v\n", err)
		os.Exit(1)
	}

	destSubDir, err := createWallpaperDirectory(destDir, wallpaperName, setupType)
	if err != nil {
		fmt.Printf("Error creating subdirectory: %v\n", err)
		os.Exit(1)
	}

	destFilePath := filepath.Join(destSubDir, "wallpaper.png")

	if err := copyWallpaper(sourcePath, destFilePath); err != nil {
		fmt.Printf("Error copying file contents: %v\n", err)
		os.Exit(1)
	}

	if err := getDominantColors(destSubDir, destDir); err != nil {
		fmt.Printf("Error getting dominant colors: %v\n", err)
		os.Exit(1)
	}

	dominantColorsFilePath := filepath.Join(destSubDir, "dominantcolors.json")
	dominantColorsContent, err := ioutil.ReadFile(dominantColorsFilePath)
	if err != nil {
		fmt.Printf("Error reading dominant colors JSON: %v\n", err)
		os.Exit(1)
	}

	var dominantColors map[string]string
	if err := json.Unmarshal(dominantColorsContent, &dominantColors); err != nil {
		fmt.Printf("Error parsing dominant colors JSON: %v\n", err)
		os.Exit(1)
	}

	if err := createColorsRasiFile(destSubDir, dominantColors); err != nil {
		fmt.Printf("Error creating colors.rasi file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Copied %s to %s\n", sourcePath, destFilePath)
}
