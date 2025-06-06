package core

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/y3owk1n/cpenv/utils"

	"github.com/briandowns/spinner"
	"github.com/manifoldco/promptui"
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

var selectProjectRun = func(prompt promptui.Select) (int, string, error) {
	return prompt.Run()
}

func SelectProject(projects []utils.Directory) (string, error) {
	if len(projects) == 0 {
		return "", fmt.Errorf("no projects found in the vault")
	}

	prompt := promptui.Select{
		Label: "Choose a project to copy from",
		Items: generateProjectOptions(projects),
	}

	_, selectedProject, err := selectProjectRun(prompt)
	if err != nil {
		if err == promptui.ErrInterrupt {
			logrus.Debug("User aborted project selection")
			fmt.Printf("%s %s\n", utils.WarningIcon(), utils.WhiteText("Selection cancelled."))
			exitFunc(0)
		}
		return "", fmt.Errorf("error starting the selection form: %w", err)
	}

	if selectedProject == "" {
		return "", fmt.Errorf("no project selected")
	}

	logrus.Debugf("Project selected: %s", selectedProject)
	return selectedProject, nil
}

func generateProjectOptions(projects []utils.Directory) []string {
	var projectOptions []string

	for _, project := range projects {
		projectOptions = append(projectOptions, project.Name)
	}
	return projectOptions
}

func CopyEnvFilesToProject(project string, currentPath string, vaultDir string) error {
	logrus.Debugf("Vault directory details: vault_dir: %s, project: %s, current_path: %s", vaultDir, project, currentPath)

	projectPath := filepath.Join(vaultDir, project, currentPath)
	filesInProject, err := utils.ReadDirRecursiveFunc(projectPath)
	if err != nil {
		return fmt.Errorf("error reading project path: %w", err)
	}

	for _, file := range filesInProject {
		if err := processCopyEnvFileToProject(file, projectPath, currentPath, vaultDir); err != nil {
			logrus.Errorf("Error processing env file: file: %s, error: %v", file, err)
		}
	}
	return nil
}

var copyFileWithSpinnerFunc = copyFileWithSpinner

func processCopyEnvFileToProject(file, projectPath, currentPath string, vaultDir string) error {
	relativePath, err := filepath.Rel(projectPath, file)
	if err != nil {
		return fmt.Errorf("failed to compute relative path: %w", err)
	}
	destinationPath := filepath.Join(utils.GetCurrentWorkingDirectory(), currentPath)
	destinationPathWithFile := filepath.Join(destinationPath, relativePath)

	fileExists, err := utils.CheckFileExists(destinationPath, relativePath)
	if err != nil {
		return fmt.Errorf("error checking file existence: %w", err)
	}

	if !fileExists {
		logrus.Debugf("File does not exist, proceeding to copy: %s", file)
		return copyFileWithSpinnerFunc(file, destinationPathWithFile, vaultDir)
	}

	logrus.Debugf("File exists, prompting for overwrite: %s", destinationPathWithFile)
	return handleExistingFile(file, destinationPathWithFile, vaultDir)
}

func prettifiedPath(path, vaultDir string) string {
	// Resolve all paths to absolute and follow symlinks
	resolvePath := func(p string) string {
		abs, err := filepath.Abs(p)
		if err != nil {
			return p
		}
		resolved, err := filepath.EvalSymlinks(abs)
		if err != nil {
			return abs
		}
		return resolved
	}

	resolvedPath := resolvePath(path)
	resolvedVault := resolvePath(vaultDir)

	cwd, err := utils.GetWdFunc()
	if err != nil {
		logrus.Errorf("Failed to get current working directory: %v", err)
		return ""
	}
	resolvedCWD := resolvePath(cwd)

	// Check against CWD first
	if rel, err := filepath.Rel(resolvedCWD, resolvedPath); err == nil && !strings.HasPrefix(rel, "..") {
		if rel == "." {
			return "{project}"
		}
		return filepath.Join("{project}", rel)
	}

	// Check against vault directory
	if rel, err := filepath.Rel(resolvedVault, resolvedPath); err == nil && !strings.HasPrefix(rel, "..") {
		if rel == "." {
			return "{vault}"
		}
		return filepath.Join("{vault}", rel)
	}

	// Return original path if no match (not prettified)
	return path
}

func copyFileWithSpinner(sourcePath, destinationPath string, vaultDir string) error {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Prefix = fmt.Sprintf("Copying %s to %s", sourcePath, destinationPath)
	s.Start()

	destDir := filepath.Dir(destinationPath)
	logrus.Debugf("Creating directory: %s", destDir)
	if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
		logrus.Errorf("Failed to create directory %s: %v", destDir, err)
	}

	logrus.Debugf("Copying file from %s to %s", sourcePath, destinationPath)
	if err := utils.CopyFileFunc(sourcePath, destinationPath); err != nil {
		logrus.Errorf("Failed to copy file from %s to %s: %v", sourcePath, destinationPath, err)
	}
	logrus.Debugf("File copied successfully: %s", destinationPath)

	s.Stop()

	fmt.Printf("%s %s %s %s %s\n", utils.SuccessIcon(), utils.WhiteText("Copied"), utils.CyanText(prettifiedPath(sourcePath, vaultDir)), utils.WhiteText("to"), utils.CyanText(prettifiedPath(destinationPath, vaultDir)))
	return nil
}

func handleExistingFile(sourcePath, destinationPath string, vaultDir string) error {
	fmt.Printf("\n%s %s\n", utils.InfoIcon(), fmt.Sprintf("Processing for: %s", utils.CyanText(destinationPath)))
	fmt.Printf("%s ", "File exists! Do you want to overwrite? (y/N): ")

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		logrus.Fatalf("Failed to read input: %v", err)
	}
	input = strings.TrimSpace(input)
	if strings.ToLower(input) != "y" {
		logrus.Debugf("User chose not to overwrite file: %s", destinationPath)
		fmt.Printf("%s %s\n", utils.WarningIcon(), utils.WhiteText("Skipped."))
		return nil
	}

	return copyFileWithSpinnerFunc(sourcePath, destinationPath, vaultDir)
}

var exitFunc = os.Exit

func ConfirmCwd() error {
	dir := utils.GetCurrentWorkingDirectory()
	logrus.Debugf("Current working directory: %s", dir)

	fmt.Printf("%s %s\n\n", utils.InfoIcon(), fmt.Sprintf("Current working directory: %s", utils.CyanText(dir)))
	fmt.Printf("%s ", "Is this your root directory to perform the backup? (y/N): ")

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		logrus.Fatalf("Failed to read input: %v", err)
	}
	input = strings.TrimSpace(input)
	if strings.ToLower(input) != "y" {
		logrus.Debugf("User chose not to backup to: %s", dir)
		fmt.Printf("%s %s\n", utils.WarningIcon(), utils.WhiteText("Aborted... cd to your desired directory and restart the backup again."))
		exitFunc(0)
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

	filesInProject, err := utils.ReadDirRecursiveFunc(dir)
	if err != nil {
		return fmt.Errorf("error reading project path: %w", err)
	}

	for _, file := range filesInProject {
		if err := processCopyEnvFileToVault(file, dir, destinationPath, vaultDir); err != nil {
			logrus.Errorf("Error processing env file: file: %s, error: %v", file, err)
		}
	}
	return nil
}

func processCopyEnvFileToVault(file, cwd, destinationPath string, vaultDir string) error {
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
	return copyFileWithSpinnerFunc(sourcePath, destinationPathWithFile, vaultDir)
}
