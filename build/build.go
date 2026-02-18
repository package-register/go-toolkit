package build

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/duke-git/lancet/v2/fileutil"
)

// Platform represents a build platform with OS and Architecture.
type Platform struct {
	OS   string
	Arch string
}

type Option struct {
	Path      string
	ZipMode   bool
	Platforms []Platform
}

type OptionFunc func(*Option)

func WithPath(path string) OptionFunc {
	return func(o *Option) {
		o.Path = path
	}
}

func WithZip(zip bool) OptionFunc {
	return func(o *Option) {
		o.ZipMode = zip
	}
}

func WithPlaftforms(p Platform) OptionFunc {
	return func(o *Option) {
		o.Platforms = append(o.Platforms, p)
	}
}

// Platforms contains the list of platforms to build for.
var Platforms = []Platform{
	{OS: "windows", Arch: "amd64"},
	{OS: "linux", Arch: "amd64"},
	{OS: "darwin", Arch: "amd64"},
}

var defaultOpt = &Option{
	Path:      "bin",
	ZipMode:   false,
	Platforms: Platforms,
}

func Builder(opts ...OptionFunc) error {
	for _, opt := range opts {
		opt(defaultOpt)
	}

	return Build(defaultOpt)
}

func Build(option *Option) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(option.Platforms))

	if len(option.Platforms) == 0 {
		option.Platforms = defaultOpt.Platforms
	}

	if err := fileutil.CreateDir(option.Path); err != nil {
		return fmt.Errorf("failed to create %s directory: %w", option.Path, err)
	}

	for _, platform := range Platforms {
		wg.Add(1)
		go func(p Platform) {
			defer wg.Done()
			if err := buildForPlatform(option.Path, p.OS, p.Arch); err != nil {
				errChan <- fmt.Errorf("failed to build for %s/%s: %w", p.OS, p.Arch, err)
			}
		}(platform)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	if option.ZipMode {
		if err := fileutil.Zip(option.Path, fmt.Sprintf("%s.zip", option.Path)); err != nil {
			return fmt.Errorf("failed to zip %s: %w", option.Path, err)
		}
	}

	fmt.Println("All builds completed successfully.")
	return nil
}

// buildForPlatform builds the application for a specific platform.
func buildForPlatform(path, goos, goarch string) error {
	fmt.Printf("Starting build for %s/%s\n", goos, goarch)
	outputName := fmt.Sprintf("app_%s_%s", goos, goarch)
	if goos == "windows" {
		outputName += ".exe"
	}

	outputPath := filepath.Join(path, outputName)

	cmd := exec.Command("go", "build", "-o", outputPath)
	cmd.Env = append(os.Environ(), fmt.Sprintf("GOOS=%s", goos), fmt.Sprintf("GOARCH=%s", goarch))
	fmt.Printf("Running build command: %v\n", cmd.Args)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("build command failed: %w", err)
	}
	fmt.Printf("Completed build for %s/%s\n", goos, goarch)
	return nil
}
