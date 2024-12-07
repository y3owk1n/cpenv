package core

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"

	"github.com/y3owk1n/cpenv/utils"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
)

func ProjectManagementFunctions() map[string]interface{} {
	return map[string]interface{}{
		"GetProjectsList":       getFunctionInfo(GetProjectsList),
		"SelectProject":         getFunctionInfo(SelectProject),
		"CopyEnvFilesToProject": getFunctionInfo(CopyEnvFilesToProject),
		"ConfirmCwd":            getFunctionInfo(ConfirmCwd),
		"CopyEnvFilesToVault":   getFunctionInfo(CopyEnvFilesToVault),
	}
}

func getFunctionInfo(fn interface{}) map[string]interface{} {
	fnType := reflect.TypeOf(fn)

	info := map[string]interface{}{
		"name":         runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name(),
		"input_count":  fnType.NumIn(),
		"output_count": fnType.NumOut(),
		"input_types":  make([]string, fnType.NumIn()),
		"output_types": make([]string, fnType.NumOut()),
	}

	for i := 0; i < fnType.NumIn(); i++ {
		info["input_types"].([]string)[i] = fnType.In(i).String()
	}

	for i := 0; i < fnType.NumOut(); i++ {
		info["output_types"].([]string)[i] = fnType.Out(i).String()
	}

	return info
}

func GetProjectsList() ([]utils.Directory, error) {
	utils.Logger.Debug("Entering GetProjectsList",
		"function_type", reflect.TypeOf(GetProjectsList),
	)

	config, err := LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	utils.Logger.Debug("Vault directory details",
		"vault_dir", config.VaultDir,
	)

	dir, err := CreateEnvFilesDirectoryIfNotFound(config.VaultDir)
	if err != nil {
		return nil, fmt.Errorf("error creating directory: %w", err)
	}

	directories, err := utils.GetDirectories(dir)
	if err != nil {
		return nil, fmt.Errorf("error retrieving directories: %w", err)
	}

	utils.Logger.Debug("Projects retrieved",
		"project_count", len(directories),
		"projects", directoriesToStringSlice(directories),
	)

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

	catppuccin := huh.ThemeCatppuccin()

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose a project").
				Options(generateProjectOptions(projects)...).
				Value(&selectedProject),
		)).WithTheme(catppuccin)

	err := form.Run()
	if err != nil {
		if err.Error() == "user aborted" {
			utils.Logger.Debug("User aborted project selection")
			fmt.Println("Until next time!")
			os.Exit(0)
		}

		return "", fmt.Errorf("error starting the selection form: %w", err)
	}

	if selectedProject == "" {
		return "", fmt.Errorf("no project selected")
	}

	utils.Logger.Debug("Project selected",
		"project", selectedProject,
	)

	return selectedProject, nil
}

func generateProjectOptions(projects []utils.Directory) []huh.Option[string] {
	options := make([]huh.Option[string], len(projects))

	for i, project := range projects {
		options[i] = huh.NewOption(project.Name, project.Value)
	}
	return options
}

func CopyEnvFilesToProject(project string, currentPath string) error {
	config, err := LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	utils.Logger.Debug("Vault directory details",
		"vault_dir", config.VaultDir,
		"project", project,
		"current_path", currentPath,
	)

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	envFilesDirectory := filepath.Join(homeDir, config.VaultDir)
	projectPath := filepath.Join(envFilesDirectory, project, currentPath)

	filesInProject, err := utils.ReadDirRecursive(projectPath)
	if err != nil {
		return fmt.Errorf("error reading project path: %w", err)
	}

	for _, file := range filesInProject {
		if err := processCopyEnvFileToProject(file, projectPath, currentPath); err != nil {
			utils.Logger.Error("Error processing env file",
				"file", file,
				"error", err,
			)
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
		return copyFileWithSpinner(file, destinationPathWithFile)
	}

	// If file exists, prompt for overwrite
	return handleExistingFile(file, destinationPathWithFile)
}

func copyFileWithSpinner(sourcePath, destinationPath string) error {
	action := func() {
		destDir := filepath.Dir(destinationPath)

		if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
			utils.Logger.Error("Failed to create directory", "message", err)
			return
		}

		if err := utils.CopyFile(sourcePath, destinationPath); err != nil {
			utils.Logger.Error("Failed to copy file", "message", err)
			return
		}
	}

	_ = spinner.New().
		Title(fmt.Sprintf("Copying %s to %s", sourcePath, destinationPath)).
		Action(action).
		Run()

	fmt.Println(utils.SuccessMessage.Render("Copied", sourcePath, "", destinationPath))
	return nil
}

func handleExistingFile(sourcePath, destinationPath string) error {
	var confirm bool
	catppuccin := huh.ThemeCatppuccin()

	form := huh.NewForm(huh.NewGroup(
		huh.NewConfirm().
			Title(fmt.Sprintf("%s exists, do you want to overwrite?", destinationPath)).
			Affirmative("Yes!").
			Negative("No.").
			Value(&confirm),
	)).WithTheme(catppuccin)

	err := form.Run()
	if err != nil {
		if err.Error() == "user aborted" {
			utils.Logger.Debug("User aborted project selection")
			fmt.Println("Until next time!")
			os.Exit(0)
		}
		return fmt.Errorf("error confirming overwrite: %w", err)
	}

	if !confirm {
		fmt.Println(utils.WarningMessage.Render("Skipped", destinationPath))
		return nil
	}

	return copyFileWithSpinner(sourcePath, destinationPath)
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
			utils.Logger.Debug("User aborted project selection")
			fmt.Println("Until next time!")
			os.Exit(0)
		}
		return fmt.Errorf("error confirming cwd form: %w", err)
	}

	if !confirm {
		fmt.Println(utils.WarningMessage.Render("cd to your desired directory and restart the backup."))
		return nil
	}

	return nil
}

func CopyEnvFilesToVault() error {
	config, err := LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	utils.Logger.Debug("Vault directory details",
		"vault_dir", config.VaultDir,
	)

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
		if err := processCopyEnvFileToVault(file, dir, destinationPath); err != nil {
			utils.Logger.Error("Error processing env file",
				"file", file,
				"error", err,
			)
		}
	}
	return nil
}

func processCopyEnvFileToVault(file, cwd, destinationPath string) error {
	fileName := filepath.Base(file)

	if strings.Contains(fileName, "node_modules") || strings.Contains(fileName, ".template") || strings.Contains(fileName, ".example") {
		return nil
	}

	if !strings.HasSuffix(fileName, ".env") {
		return nil
	}

	relativePath := strings.TrimPrefix(file, cwd+"/")

	sourcePath := filepath.Join(cwd, relativePath)
	destinationPathWithFile := filepath.Join(destinationPath, relativePath)

	return copyFileWithSpinner(sourcePath, destinationPathWithFile)
}
