package core

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/y3owk1n/cpenv/utils"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
)

func GetProjectsList(vaultDir string) ([]utils.Directory, error) {
	logrus.Debugf("Entering GetProjectsList, function_type: %v", reflect.TypeOf(GetProjectsList))
	logrus.Debugf("Vault directory details: %s", vaultDir)

	directories, err := utils.GetDirectories(vaultDir)
	if err != nil {
		return nil, fmt.Errorf("error retrieving directories: %w", err)
	}

	if len(directories) == 0 {
		return nil, fmt.Errorf("no projects found in the vault")
	}

	logrus.Debugf("Projects retrieved, project_count: %d, projects: %v", len(directories), directoriesToStringSlice(directories))
	return directories, nil
}

func directoriesToStringSlice(dirs []utils.Directory) []string {
	names := make([]string, len(dirs))
	for i, dir := range dirs {
		names[i] = dir.Name
	}
	return names
}

func SelectProject(projects []utils.Directory) (string, error) {
	if len(projects) == 0 {
		return "", fmt.Errorf("no projects found in the vault")
	}

	var selectedProject string
	baseTheme := huh.ThemeBase()

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose a project to copy from").
				Options(generateProjectOptions(projects)...).
				Value(&selectedProject),
		),
	).WithTheme(baseTheme)

	err := form.Run()
	if err != nil {
		if err == huh.ErrUserAborted {
			logrus.Debug("User aborted project selection")
			fmt.Println("Until next time!")
			os.Exit(0)
		}
		return "", fmt.Errorf("error starting the selection form: %w", err)
	}

	if selectedProject == "" {
		return "", fmt.Errorf("no project selected")
	}

	logrus.Debugf("Project selected: %s", selectedProject)
	return selectedProject, nil
}

func generateProjectOptions(projects []utils.Directory) []huh.Option[string] {
	options := make([]huh.Option[string], len(projects))
	for i, project := range projects {
		options[i] = huh.NewOption(project.Name, project.Value)
	}
	return options
}

func CopyEnvFilesToProject(project string, currentPath string, vaultDir string) error {
	logrus.Debugf("Vault directory details: vault_dir: %s, project: %s, current_path: %s", vaultDir, project, currentPath)

	projectPath := filepath.Join(vaultDir, project, currentPath)
	filesInProject, err := utils.ReadDirRecursive(projectPath)
	if err != nil {
		return fmt.Errorf("error reading project path: %w", err)
	}

	for _, file := range filesInProject {
		if err := processCopyEnvFileToProject(file, projectPath, currentPath); err != nil {
			logrus.Errorf("Error processing env file: file: %s, error: %v", file, err)
		}
	}
	return nil
}

func processCopyEnvFileToProject(file, projectPath, currentPath string) error {
	relativePath := strings.TrimPrefix(file, projectPath+"/")
	destinationPath := filepath.Join(utils.GetCurrentWorkingDirectory(), currentPath)
	destinationPathWithFile := filepath.Join(destinationPath, relativePath)

	fileExists, err := utils.CheckFileExists(destinationPath, relativePath)
	if err != nil {
		return fmt.Errorf("error checking file existence: %w", err)
	}

	if !fileExists {
		logrus.Debugf("File does not exist, proceeding to copy: %s", file)
		return copyFileWithSpinner(file, destinationPathWithFile)
	}

	logrus.Debugf("File exists, prompting for overwrite: %s", destinationPathWithFile)
	return handleExistingFile(file, destinationPathWithFile)
}

func copyFileWithSpinner(sourcePath, destinationPath string) error {
	action := func() {
		destDir := filepath.Dir(destinationPath)
		logrus.Debugf("Creating directory: %s", destDir)
		if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
			logrus.Errorf("Failed to create directory %s: %v", destDir, err)
			return
		}

		logrus.Debugf("Copying file from %s to %s", sourcePath, destinationPath)
		if err := utils.CopyFile(sourcePath, destinationPath); err != nil {
			logrus.Errorf("Failed to copy file from %s to %s: %v", sourcePath, destinationPath, err)
			return
		}
		logrus.Debugf("File copied successfully: %s", destinationPath)
	}

	spinnerTitle := fmt.Sprintf("Copying %s to %s", sourcePath, destinationPath)
	logrus.Debugf("Starting spinner with title: %s", spinnerTitle)
	_ = spinner.New().
		Title(spinnerTitle).
		Action(action).
		Run()

	fmt.Println(utils.SuccessMessage.Render("Copied", sourcePath, "", destinationPath))
	return nil
}

func handleExistingFile(sourcePath, destinationPath string) error {
	var confirm bool
	baseTheme := huh.ThemeBase()

	form := huh.NewForm(huh.NewGroup(
		huh.NewConfirm().
			Title("File exists! Do you want to overwrite?").
			Description(fmt.Sprintf("File: %s", destinationPath)).
			Affirmative("Yes!").
			Negative("No.").
			Value(&confirm),
	)).WithTheme(baseTheme)

	err := form.Run()
	if err != nil {
		if err == huh.ErrUserAborted {
			logrus.Debug("User aborted project selection")
			fmt.Println("Until next time!")
			os.Exit(0)
		}
		return fmt.Errorf("error confirming overwrite: %w", err)
	}

	if !confirm {
		logrus.Debugf("User chose not to overwrite file: %s", destinationPath)
		fmt.Println(utils.WarningMessage.Render("Skipped", destinationPath))
		return nil
	}

	return copyFileWithSpinner(sourcePath, destinationPath)
}

func ConfirmCwd() error {
	dir := utils.GetCurrentWorkingDirectory()
	logrus.Debugf("Current working directory: %s", dir)

	var confirm bool
	baseTheme := huh.ThemeBase()

	form := huh.NewForm(huh.NewGroup(
		huh.NewConfirm().
			Title("Is this your root directory to perform the backup?").
			Description(fmt.Sprintf("Current Root Directory: %s", dir)).
			Affirmative("Yes!").
			Negative("No.").
			Value(&confirm),
	)).WithTheme(baseTheme)

	err := form.Run()
	if err != nil {
		if err == huh.ErrUserAborted {
			logrus.Debug("User aborted project selection")
			fmt.Println("Until next time!")
			os.Exit(0)
		}
		return fmt.Errorf("error confirming cwd form: %w", err)
	}

	if !confirm {
		logrus.Debug("User did not confirm the current working directory")
		fmt.Println(utils.WarningMessage.Render("cd to your desired directory and restart the backup."))
		os.Exit(0)
		return nil
	}

	logrus.Debug("Current working directory confirmed by user")
	return nil
}

func CopyEnvFilesToVault(vaultDir string) error {
	logrus.Debugf("Vault directory details: %s", vaultDir)

	dir := utils.GetCurrentWorkingDirectory()
	logrus.Debugf("Current working directory for backup: %s", dir)

	currentProjectFolderName := filepath.Base(dir)
	if currentProjectFolderName == "" {
		return fmt.Errorf("failed to parse the folder name, try again")
	}

	currentProjectFolderNameWithTimestamp := fmt.Sprintf("%s-%s", currentProjectFolderName, utils.GetBackupTimestamp())
	logrus.Debugf("Current project folder with timestamp: %s", currentProjectFolderNameWithTimestamp)

	destinationPath := filepath.Join(vaultDir, currentProjectFolderNameWithTimestamp)
	logrus.Debugf("Destination path for backup: %s", destinationPath)

	if err := os.MkdirAll(destinationPath, os.ModePerm); err != nil {
		logrus.Errorf("Failed to create destination path: %v", err)
		return fmt.Errorf("failed to create destination path: %w", err)
	}
	logrus.Debugf("Destination path created: %s", destinationPath)

	filesInProject, err := utils.ReadDirRecursive(dir)
	if err != nil {
		return fmt.Errorf("error reading project path: %w", err)
	}

	for _, file := range filesInProject {
		if err := processCopyEnvFileToVault(file, dir, destinationPath); err != nil {
			logrus.Errorf("Error processing env file: file: %s, error: %v", file, err)
		}
	}
	return nil
}

func processCopyEnvFileToVault(file, cwd, destinationPath string) error {
	fileName := filepath.Base(file)
	fullPath, _ := filepath.Abs(file)

	if strings.Contains(fullPath, "node_modules") || strings.Contains(fileName, ".template") || strings.Contains(fileName, ".example") {
		logrus.Debugf("Skipping file: %s", file)
		return nil
	}

	if !strings.HasSuffix(fileName, ".env") {
		logrus.Debugf("Skipping non-env file: %s", file)
		return nil
	}

	relativePath := strings.TrimPrefix(file, cwd+"/")
	sourcePath := filepath.Join(cwd, relativePath)
	destinationPathWithFile := filepath.Join(destinationPath, relativePath)

	logrus.Debugf("Copying env file to vault: source: %s, destination: %s", sourcePath, destinationPathWithFile)
	return copyFileWithSpinner(sourcePath, destinationPathWithFile)
}
