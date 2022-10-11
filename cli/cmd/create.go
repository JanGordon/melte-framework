package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var cwd, err = os.Getwd()
var createCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"c"},
	Short:   "creates a new melte project",
	Run: func(cmd *cobra.Command, args []string) {
		// creat here?
		createHerePrompt := promptui.Prompt{
			Label:     "Create project in this directory? ",
			IsConfirm: false,
		}

		createHere, err := createHerePrompt.Run()

		if err != nil {
			color.Red(fmt.Sprintf("exited %v\n", err))
			return
		}
		projPath := cwd
		overwriteFiles := false
		if createHere != "y" {
			// get project path
			validate := func(input string) error {
				return nil
			}

			templates := &promptui.PromptTemplates{
				Prompt:  "{{ . | yellow}} ",
				Valid:   "{{ . | green }} ",
				Invalid: "{{ . | red }} ",
				Success: "{{ . | bold }} ",
			}

			pathPrompt := promptui.Prompt{

				Label:     "Create at (path) : ",
				Templates: templates,
				Validate:  validate,
			}

			projPath, err = pathPrompt.Run()
			if err != nil {
				color.Red(fmt.Sprintf("exited %v\n", err))
				return
			}
			if projPath == "" {
				projPath = cwd
			}
			_, err := os.Stat(projPath)
			if err != nil {
				color.Blue("Creating directory : ", projPath)
				err := os.MkdirAll(projPath, 0755)
				if err != nil {
					panic(err)
				}
			}

		}

		if !IsEmpty(projPath) {
			// creat here?
			overwritePrompt := promptui.Prompt{
				Label:     "There are files in this directory, would you like to overwrite them? (y/N)",
				IsConfirm: false,
			}

			overwrite, err := overwritePrompt.Run()

			if err != nil {
				color.Red(fmt.Sprintf("exited %v\n", err))
				return
			}
			if overwrite == "y" {
				// creat here?
				overwritepPrompt := promptui.Prompt{
					Label:     "Are you sure you want to overwrite files? ",
					IsConfirm: false,
				}

				overwritep, err := overwritepPrompt.Run()

				if err != nil {
					color.Red(fmt.Sprintf("exited %v\n", err))
					return
				}

				if overwritep == "y" {
					overwriteFiles = true
				} else {
					color.Red(fmt.Sprintf("exited"))
					return
				}

			}
		}
		// get project type
		prompt := promptui.Select{
			Label: "Select project type: ",
			Items: []string{"Skeleton Project", "Demo App", "Notes App"},
		}

		_, projType, err := prompt.Run()

		if err != nil {
			color.Red(fmt.Sprintf("exited %v\n", err))
			return
		}

		folders := []string{"routes", "public", "components"}
		for _, folder := range folders {
			err = os.Mkdir(filepath.Join(projPath, folder), 0755)
			if err != nil {
				panic(err)
			}
		}
		color.Green(fmt.Sprintf("Creating %q at %s \n overwrite files : %t", projType, projPath, overwriteFiles))

	},
}

func init() {
	// the commadn should build scaffold and cd into the project
	rootCmd.AddCommand(createCmd)
}

func IsEmpty(name string) bool {
	f, err := os.Open(name)
	if err != nil {
		return false
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	return err == io.EOF       // Either not empty or error, suits both cases
}
