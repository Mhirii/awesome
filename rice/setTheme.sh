#!/bin/bash

# Check if an argument is provided
if [ $# -ne 1 ]; then
  echo "Usage: $0 <theme_name>"
  exit 1
fi


# Define the theme directory

theme_name="$1"
theme_dir="$HOME/.config/awesome/rice/themes/$theme_name"

# Check if the theme directory exists
if [ -d "$theme_dir" ]; then
  # Copy colors.Xresources to $HOME/.Xresources
  cp "$theme_dir/colors.Xresources" "$HOME/.Xresources"
  echo "Copied colors.Xresources to $HOME/.Xresources"

  # Copy colors-rofi-dark.rasi to $HOME/.config/awesome/rice/rofi/theme.rasi
  cp "$theme_dir/colors-rofi-dark.rasi" "$HOME/.config/awesome/rice/rofi/theme.rasi"
  echo "Copied colors-rofi-dark.rasi to $HOME/.config/awesome/rice/rofi/theme.rasi"

  # Copy colors.yml to $HOME/.config/alacritty/colors.yml
  cp "$theme_dir/colors.yml" "$HOME/.config/alacritty/colors.yml"
  echo "Copied colors.yml to $HOME/.config/alacritty/colors.yml"
else
  echo "Theme directory '$theme_dir' does not exist."
fi
