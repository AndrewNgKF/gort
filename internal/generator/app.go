package generator

import (
	"fmt"
	"os"
	"path/filepath"
)

// NewApp creates a new Gort application with the complete directory structure
func NewApp(name string, apiMode bool) error {
	// Create base directories
	dirs := []string{
		"app/controllers",
		"app/models",
		"app/middleware",
		"config/environments",
		"db/migrations",
		"test/controllers",
		"test/models",
		"test/integration",
		"tmp/pids",
	}

	// Add view and asset directories only if not in API mode
	if !apiMode {
		dirs = append(dirs, []string{
			"app/views/layouts",
			"app/views/home",
			"public/assets/css",
			"public/assets/js",
			"public/assets/images",
			"public/uploads",
		}...)
	}

	for _, dir := range dirs {
		path := filepath.Join(name, dir)
		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", path, err)
		}
	}

	// Create files
	if err := createMainGo(name); err != nil {
		return err
	}
	if err := createGoMod(name); err != nil {
		return err
	}
	if err := createReadme(name); err != nil {
		return err
	}
	if err := createGitignore(name); err != nil {
		return err
	}
	if err := createConfigFiles(name); err != nil {
		return err
	}
	if err := createControllerFiles(name); err != nil {
		return err
	}

	// Create view and asset files only if not in API mode
	if !apiMode {
		if err := createViewFiles(name); err != nil {
			return err
		}
		if err := createCSSFile(name); err != nil {
			return err
		}
	}

	if err := createSeedsFile(name); err != nil {
		return err
	}

	return nil
}

func writeFile(name, path, content string) error {
	fullPath := filepath.Join(name, path)
	return os.WriteFile(fullPath, []byte(content), 0644)
}
