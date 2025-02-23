package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/user"
	"path"
	"strings"

	"github.com/spf13/viper"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath("$HOME/.nyxpad/")

	viper.SetDefault("font-family", "sans-serif")
	viper.SetDefault("font-size", "16px")
	viper.SetDefault("background-color", "#ffffff")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("Config file does not exist. Creating...")

			usr, err := user.Current()
			if err != nil {
				log.Fatal(err)
			}

			f, err := os.Create(path.Join(usr.HomeDir, ".nyxpad", "config.toml"))
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()

			err = viper.WriteConfig()
			if err != nil {
				log.Fatal(err)
			}

			err = viper.ReadInConfig()
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("Created!")
		} else {
			log.Fatal(err)
		}
	}
}

//? FILE FUNCTIONS ?//

// Output of ReadFile
type ReadFileOutput struct {
	Content string `json:"content"`
	File    string `json:"file"`
}

// Opens a file dialog, then reads the selected file and returns
// the file & file content.
func (a *App) ReadFile() (output ReadFileOutput, _ error) {
	openFileOptions := runtime.OpenDialogOptions{
		Title:           "Open File",
		DefaultFilename: "Untitled.txt",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Text Files (*.txt;*.md)",
				Pattern:     "*.txt;*.md",
			},
		},
	}

	f, err := runtime.OpenFileDialog(a.ctx, openFileOptions)
	if err != nil {
		fmt.Println("ln93 ", err)
		return ReadFileOutput{}, err
	}

	c, err := os.ReadFile(f)
	if err != nil {
		fmt.Println("ln99 ", err)
		return ReadFileOutput{}, err
	}

	if f == "" {
		fmt.Println("user cancelled")
		return ReadFileOutput{
			File:    "",
			Content: "",
		}, nil
	}

	return ReadFileOutput{
		File:    f,
		Content: string(c),
	}, nil
}

// Save a new file
func (a *App) SaveFile(content string) (file string, _ error) {
	saveFileOptions := runtime.SaveDialogOptions{
		Title:           "Save File",
		DefaultFilename: "Untitled.txt",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Text Files (*.txt;*.md)",
				Pattern:     "*.txt;*.md",
			},
		},
	}

	f, err := runtime.SaveFileDialog(a.ctx, saveFileOptions)
	if err != nil {
		return "", err
	}

	c, err := os.Create(f)
	if err != nil {
		return "", err
	}
	defer c.Close()

	_, err = c.Write([]byte(content))
	if err != nil {
		return "", err
	}

	return f, nil
}

// Write to a file
func (a *App) WriteFile(file string, content string) error {
	content = strings.ReplaceAll(content, "\n", "\r\n")
	return os.WriteFile(file, []byte(content), 0644)
}

//? CONFIG ?//

func (a *App) ConfigGet(key string) string {
	return viper.GetString(key)
}
