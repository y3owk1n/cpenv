package core

import (
	"cpenv/utils"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
)

func GetProjectsList() ([]utils.Directory, error) {
	config, err := LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	utils.Logger.Debugf("Vault directory: %s", config.VaultDir)

	dir, err := CreateEnvFilesDirectoryIfNotFound(config.VaultDir)
	if err != nil {
		return nil, fmt.Errorf("error creating directory: %w", err)
	}
	utils.Logger.Debugf("Directory: %s", dir)

	directories, err := utils.GetDirectories(dir)
	if err != nil {
		return nil, fmt.Errorf("error getting directories: %w", err)
	}
	utils.Logger.Debugf("Directories: %s", directories)

	return directories, nil
}

func SelectProject(projects []utils.Directory) (string, error) {
	var selectedProject string

	if len(projects) == 0 {
		return "", fmt.Errorf("no projects found in the vault")
	}

	catppuccin := huh.ThemeCatppuccin()

	form := huh.NewForm(
		huh.NewGroup(
			// Ask the user for a base burger and toppings.
			huh.NewSelect[string]().
				Title("Choose a project").
				Options(generateProjectOptions(projects)...).
				Value(&selectedProject), // store the chosen project in the "selectedProject" variable
		)).WithTheme(catppuccin)

	err := form.Run()
	if err != nil {
		if err.Error() == "user aborted" {
			fmt.Println("Until next time!")
			os.Exit(0)
		}

		return "", fmt.Errorf("error starting the selection form: %w", err)
	}

	if selectedProject == "" {
		return "", fmt.Errorf("no project selected")
	}

	return selectedProject, nil
}

func generateProjectOptions(projects []utils.Directory) []huh.Option[string] {
	options := []huh.Option[string]{}
	for _, project := range projects {
		options = append(options, huh.NewOption(project.Name, project.Value))
	}
	return options
}

func CopyEnvFilesToProject(project string, currentPath string) error {
	// Fetch environment files directory
	config, err := LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	utils.Logger.Debugf("Vault directory: %s", config.VaultDir)

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	envFilesDirectory := filepath.Join(homeDir, config.VaultDir)

	// Construct the project path
	projectPath := filepath.Join(envFilesDirectory, project, currentPath)

	// Read files from the project path
	filesInProject, err := utils.ReadDirRecursive(projectPath)
	if err != nil {
		return fmt.Errorf("error reading project path: %w", err)
	}

	for _, file := range filesInProject {
		sourcePath := file
		relativePath := strings.TrimPrefix(file, projectPath+"/")
		destinationPath := filepath.Join(utils.GetCurrentWorkingDirectory(), currentPath)

		destinationPathWithFile := filepath.Join(destinationPath, relativePath)

		if strings.Contains(file, ".env") {
			// Check if file already exists in destination
			fileExists, err := utils.CheckFileExists(destinationPath, relativePath)
			if err != nil {
				return fmt.Errorf("error checking if file exists: %w", err)
			}

			// If file doesn't exist, copy it
			if !fileExists {
				action := func() {
					destDir := filepath.Dir(destinationPathWithFile)

					err = os.MkdirAll(destDir, os.ModePerm)
					if err != nil {
						utils.Logger.Error("Failed to create directory", "message", err)
					}

					err = utils.CopyFile(sourcePath, destinationPathWithFile)
					if err != nil {
						utils.Logger.Error("Failed to copy file", "message", err)
					}
				}
				_ = spinner.New().Title(fmt.Sprintf("Copying %s to %s", file, destinationPathWithFile)).Action(action).Run()

				fmt.Println(utils.SuccessMessage.Render(" Copied", file, "", destinationPathWithFile))
			} else {
				var confirm bool

				catppuccin := huh.ThemeCatppuccin()

				form := huh.NewForm(huh.NewGroup(huh.NewConfirm().
					Title(fmt.Sprintf("%s exists, do you want to overwrite?", destinationPathWithFile)).
					Affirmative("Yes!").
					Negative("No.").
					Value(&confirm))).WithTheme(catppuccin)

				err := form.Run()
				if err != nil {
					if err.Error() == "user aborted" {
						fmt.Println("Until next time!")
						os.Exit(0)
					}
					return fmt.Errorf("error confirming cwd form: %w", err)
				}

				if !confirm {
					fmt.Println(utils.WarningMessage.Render("Skipped", destinationPathWithFile))
				} else {
					destDir := filepath.Dir(destinationPathWithFile)

					if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
						utils.Logger.Error("Failed to create directory", "message", err)
					}

					err := utils.CopyFile(sourcePath, destinationPathWithFile)
					if err != nil {
						utils.Logger.Error("Failed to copy file", "message", err)
					}

					fmt.Println(utils.SuccessMessage.Render(" Copied", file, "", destinationPathWithFile))
				}
			}
		}
	}
	return nil
}

func ConfirmCwd() error {
	dir := utils.GetCurrentWorkingDirectory()

	var confirm bool

	catppuccin := huh.ThemeCatppuccin()

	form := huh.NewForm(huh.NewGroup(huh.NewConfirm().
		Title(fmt.Sprintf("Is this your root directory to perform the backup? (%s)", dir)).
		Affirmative("Yes!").
		Negative("No.").
		Value(&confirm))).WithTheme(catppuccin)

	err := form.Run()
	if err != nil {
		if err.Error() == "user aborted" {
			fmt.Println("Until next time!")
			os.Exit(0)
		}
		return fmt.Errorf("error confirming cwd form: %w", err)
	}

	if !confirm {
		fmt.Println(utils.WarningMessage.Render("cd to your desired directory and restart the backup."))
		os.Exit(0)
	}

	return nil
}

func CopyEnvFilesToVault() error {
	// Fetch environment files directory
	config, err := LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	utils.Logger.Debugf("Vault directory: %s", config.VaultDir)

	dir := utils.GetCurrentWorkingDirectory()

	currentProjectFolderName := filepath.Base(dir)
	if currentProjectFolderName == "" {
		return fmt.Errorf("failed to parse the folder name, try again...")
	}

	currentProjectFolderNameWithTimestamp := fmt.Sprintf("%s-%s", currentProjectFolderName, utils.GetBackupTimestamp())

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	envFilesDirectory := filepath.Join(homeDir, config.VaultDir)

	destinationPath := filepath.Join(envFilesDirectory, currentProjectFolderNameWithTimestamp)

	err = os.MkdirAll(destinationPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create destination path: %w", err)
	}

	filesInProject, err := utils.ReadDirRecursive(dir)
	if err != nil {
		return fmt.Errorf("error reading project path: %w", err)
	}

	for _, file := range filesInProject {
		fileName := filepath.Base(file)
		relativePath := strings.TrimPrefix(file, dir+"/")

		// Skip unwanted files or directories
		if strings.Contains(fileName, "node_modules") || strings.Contains(fileName, ".template") || strings.Contains(fileName, ".example") {
			continue
		}

		// Check if the file is a .env file
		if strings.HasSuffix(fileName, ".env") {
			sourcePath := filepath.Join(dir, relativePath)
			destinationPathWithFile := filepath.Join(destinationPath, relativePath)

			// Copy file to destination
			action := func() {
				destDir := filepath.Dir(destinationPathWithFile)

				if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
					utils.Logger.Error("Failed to create directory", "message", err)
				}

				err := utils.CopyFile(sourcePath, destinationPathWithFile)
				if err != nil {
					utils.Logger.Error("Failed to copy file", "message", err)
				}
			}
			_ = spinner.New().Title(fmt.Sprintf("Copying %s to %s", file, destinationPathWithFile)).Action(action).Run()

			fmt.Println(utils.SuccessMessage.Render(" Copied", file, "", destinationPathWithFile))
		}

	}
	return nil
}
