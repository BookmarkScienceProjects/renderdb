package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
)

func main() {
	directories, err := glideNoVendor()
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("Could not run 'glide nv': %s\n", err))
		os.Exit(1)
	}

	var wg sync.WaitGroup
	wg.Add(len(directories))
	for _, d := range directories {
		dir := d
		go func() {
			defer wg.Done()
			fmt.Printf("Generating mocks for '%s'...\n", dir)
			err = generateMocks(dir)
			if err != nil {
				os.Stderr.WriteString(fmt.Sprintf("Could not generate mocks for directory '%s': %s\n", dir, err))
				os.Exit(2)
			}
		}()
	}
	wg.Wait()
}

func glideNoVendor() ([]string, error) {
	cmd := exec.Command("glide", "nv")
	bytes, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(bytes), "\n")
	filtered := make([]string, 0, len(lines))
	for _, l := range lines {
		line := strings.TrimSuffix(l, "...")
		if line != "" && line != "." && line != "./mocks/" {
			filtered = append(filtered, line)
		}
	}
	return filtered, err
}

func generateMocks(dir string) error {
	args := []string{"-dir", dir, "-recursive", "-all", "-inpkg"}
	cmd := exec.Command("mockery", args...)
	return cmd.Run()
}
