package backend

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

type Config struct {
	AudioBackend      string `json:"audio_backend"`
	WindowsPluginsDir string `json:"windows_plugins_dir"`
	Theme             string `json:"theme"`
	FontFamily        string `json:"font_family"`
}

type Environment struct {
	JDAWDir            string
	WinePrefix         string
	WineBin            string
	YabridgeDir        string
	YabridgeLink       string
	ReaperDir          string
	ConfigPath         string
	Config             Config
}

func NewEnvironment() (*Environment, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	jdawDir := filepath.Join(homeDir, ".local", "share", "jdaw")
	configPath := filepath.Join(jdawDir, "config.json")

	env := &Environment{
		JDAWDir:      jdawDir,
		WinePrefix:   filepath.Join(jdawDir, "data"),
		WineBin:      filepath.Join(jdawDir, "wine-7.22-staging-amd64", "bin"),
		YabridgeDir:  filepath.Join(jdawDir, "yabridge"),
		YabridgeLink: filepath.Join(homeDir, ".local", "share", "yabridge"),
		ReaperDir:    filepath.Join(jdawDir, "REAPER"),
		ConfigPath:   configPath,
		Config: Config{
			AudioBackend:      "Automatic",
			WindowsPluginsDir: filepath.Join(jdawDir, "plugins", "Windows"),
			Theme:             "Dark",
			FontFamily:        "Roboto",
		},
	}

	env.LoadConfig()
	return env, nil
}

func (e *Environment) LoadConfig() {
	data, err := os.ReadFile(e.ConfigPath)
	if err == nil {
		json.Unmarshal(data, &e.Config)
	}
}

func (e *Environment) SaveConfig() error {
	os.MkdirAll(e.JDAWDir, 0755)
	data, err := json.MarshalIndent(e.Config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(e.ConfigPath, data, 0644)
}

func (e *Environment) SetupYabridge() error {
	os.MkdirAll(filepath.Dir(e.YabridgeLink), 0755)
	
	// Create/Update symlink
	if err := os.Remove(e.YabridgeLink); err != nil && !os.IsNotExist(err) {
		return err
	}
	if err := os.Symlink(e.YabridgeDir, e.YabridgeLink); err != nil {
		return err
	}

	// Add dynamic path
	yabridgectl := filepath.Join(e.YabridgeLink, "yabridgectl")
	if _, err := os.Stat(yabridgectl); err == nil {
		cmd := exec.Command(yabridgectl, "add", e.Config.WindowsPluginsDir)
		cmd.Run()

		// Sync
		cmdSync := exec.Command(yabridgectl, "sync", "--prune")
		cmdSync.Run()
	}

	return nil
}

func (e *Environment) LaunchReaper() error {
	if err := e.SetupYabridge(); err != nil {
		fmt.Printf("Warning: yabridge setup failed: %v\n", err)
	}

	env := os.Environ()
	path := fmt.Sprintf("%s:%s:%s", e.WineBin, e.ReaperDir, os.Getenv("PATH"))
	
	// Add updated path and wineprefix
	var newEnv []string
	for _, v := range env {
		newEnv = append(newEnv, v)
	}
	newEnv = append(newEnv, "PATH="+path)
	newEnv = append(newEnv, "WINEPREFIX="+e.WinePrefix)

	isWayland := os.Getenv("XDG_SESSION_TYPE") == "wayland"

	var command []string

	usePWJack := false
	if e.Config.AudioBackend == "Pipewire" {
		usePWJack = true
	} else if e.Config.AudioBackend == "Automatic" {
		if _, err := exec.LookPath("pw-jack"); err == nil {
			usePWJack = true
		}
	}

	if isWayland {
		newEnv = append(newEnv, "DISPLAY=:0")
		newEnv = append(newEnv, "GDK_BACKEND=x11")
		newEnv = append(newEnv, "QT_QPA_PLATFORM=xcb")
		
		// Remove WAYLAND_DISPLAY
		var filteredEnv []string
		for _, v := range newEnv {
			if len(v) < 15 || v[:15] != "WAYLAND_DISPLAY" {
				filteredEnv = append(filteredEnv, v)
			}
		}
		newEnv = filteredEnv
	}

	if usePWJack {
		command = []string{"pw-jack", "reaper"}
	} else {
		command = []string{"reaper"}
	}

	cmd := exec.Command(command[0], command[1:]...)
	cmd.Env = newEnv
	
	// Start in background
	return cmd.Start()
}

func (e *Environment) DownloadWine(progressCallback func(float64)) error {
	os.MkdirAll(e.JDAWDir, 0755)
	
	url := "https://github.com/Kron4ek/Wine-Builds/releases/download/7.22/wine-7.22-staging-amd64.tar.xz"
	tarPath := filepath.Join(e.JDAWDir, "wine.tar.xz")

	if err := downloadFile(url, tarPath); err != nil {
		return err
	}

	// Extract
	cmd := exec.Command("tar", "-xf", tarPath, "-C", e.JDAWDir)
	if err := cmd.Run(); err != nil {
		return err
	}

	os.Remove(tarPath)
	return nil
}

func (e *Environment) DownloadReaper(progressCallback func(float64)) error {
	os.MkdirAll(e.JDAWDir, 0755)

	url := "https://www.reaper.fm/files/7.x/reaper782_linux_x86_64.tar.xz"
	tarPath := filepath.Join(e.JDAWDir, "reaper.tar.xz")

	if err := downloadFile(url, tarPath); err != nil {
		return err
	}

	// Extract
	cmd := exec.Command("tar", "-xf", tarPath, "-C", e.JDAWDir)
	if err := cmd.Run(); err != nil {
		return err
	}

	os.Remove(tarPath)
	return nil
}

func (e *Environment) DownloadYabridge(progressCallback func(float64)) error {
	os.MkdirAll(e.JDAWDir, 0755)

	url := "https://github.com/robbert-vdh/yabridge/releases/download/5.1.1/yabridge-5.1.1.tar.gz"
	tarPath := filepath.Join(e.JDAWDir, "yabridge.tar.gz")

	if err := downloadFile(url, tarPath); err != nil {
		return err
	}

	// Extract
	cmd := exec.Command("tar", "-xf", tarPath, "-C", e.JDAWDir)
	if err := cmd.Run(); err != nil {
		return err
	}

	os.Remove(tarPath)
	return nil
}

// downloadFile downloads a URL to a local file path.
func downloadFile(url, destPath string) error {
	out, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	_, err = io.Copy(out, resp.Body)
	return err
}

