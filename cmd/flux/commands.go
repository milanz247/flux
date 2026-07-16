package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// cmdServe runs the Go HTTP server and the Vite dev server side by side,
// streaming both outputs — the Flux equivalent of `php artisan serve` +
// `npm run dev`.
func cmdServe() error {
	vite := command("npm", "run", "dev")
	if err := vite.Start(); err != nil {
		return fmt.Errorf("starting vite dev server (is npm installed?): %w", err)
	}
	defer func() {
		if vite.Process != nil {
			_ = vite.Process.Kill()
		}
	}()

	server := command("go", "run", ".")
	if err := server.Run(); err != nil {
		return fmt.Errorf("go server exited: %w", err)
	}
	return nil
}

// cmdBuild produces a production build: the Vite bundle into public/build
// and a compiled Go binary in bin/.
func cmdBuild() error {
	fmt.Println("→ building frontend (vite)")
	if err := command("npm", "run", "build").Run(); err != nil {
		return fmt.Errorf("frontend build failed: %w", err)
	}

	fmt.Println("→ building backend (go)")
	binary := filepath.Join("bin", "app")
	if runtime.GOOS == "windows" {
		binary += ".exe"
	}
	if err := command("go", "build", "-o", binary, ".").Run(); err != nil {
		return fmt.Errorf("backend build failed: %w", err)
	}

	fmt.Printf("✓ build complete: %s (set APP_ENV=production to serve the bundle)\n", binary)
	return nil
}

// cmdNew scaffolds a fresh Flux project by cloning the starter template.
func cmdNew(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: flux new <project-name>")
	}
	name := args[0]

	if _, err := os.Stat(name); err == nil {
		return fmt.Errorf("directory %q already exists", name)
	}

	template := os.Getenv("FLUX_TEMPLATE_REPO")
	if template == "" {
		return fmt.Errorf(
			"no template repository configured.\n"+
				"Set FLUX_TEMPLATE_REPO to your Flux starter repo, e.g.:\n"+
				"  set FLUX_TEMPLATE_REPO=https://github.com/you/flux-starter\n"+
				"  flux new %s\n"+
				"or clone this project directly and run `flux serve`", name)
	}

	fmt.Printf("→ creating %s from %s\n", name, template)
	clone := command("git", "clone", "--depth", "1", template, name)
	if err := clone.Run(); err != nil {
		return fmt.Errorf("git clone failed: %w", err)
	}
	_ = os.RemoveAll(filepath.Join(name, ".git"))

	fmt.Printf("✓ project created.\n\nNext steps:\n  cd %s\n  copy .env.example .env   (set APP_KEY and DB_*)\n  npm install\n  flux migrate\n  flux serve\n", name)
	return nil
}

// command builds an exec.Cmd that streams output and works on Windows
// (npm/npx are .cmd shims there).
func command(name string, args ...string) *exec.Cmd {
	if runtime.GOOS == "windows" && (name == "npm" || name == "npx") {
		args = append([]string{"/c", name}, args...)
		name = "cmd"
	}
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd
}
