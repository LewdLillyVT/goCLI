package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/robertkrimen/otto" // For JavaScript plugins
)

// Plugin struct to store plugin name and download URL
type Plugin struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// Welcome message and menu
func displayWelcomeMessage() {
	fmt.Println("=======================================================")
	fmt.Println("      Welcome to goCLI a tool made by LewdLillyVT      ")
	fmt.Println("=======================================================")
}

func displayMenu() {
	fmt.Println("\nPlease select an option:")
	fmt.Println("1. Help/Info")
	fmt.Println("2. List Available Plugins")
	fmt.Println("3. Load and Execute a Plugin")
	fmt.Println("4. Install Plugin from Library")
	fmt.Println("0. Exit")
	fmt.Print("Enter your choice: ")
}

// Function to handle Help/Info command
func displayHelp() {
	fmt.Println("\n---- Help/Info ----")
	fmt.Println("This is a modular CLI tool written in Go.")
	fmt.Println("You can load and execute plugins written in JavaScript, Python, or PowerShell.")
	fmt.Println("Choose an option from the menu to explore the tool’s features.")
}

// Function to list all plugins in the plugins folder
func listPlugins() []string {
	files, err := ioutil.ReadDir("./plugins")
	if err != nil {
		log.Fatal(err)
	}

	var plugins []string
	fmt.Println("\n---- Available Plugins ----")
	for i, f := range files {
		if !f.IsDir() && (strings.HasSuffix(f.Name(), ".js") || strings.HasSuffix(f.Name(), ".py") || strings.HasSuffix(f.Name(), ".ps1")) {
			plugins = append(plugins, f.Name())
			fmt.Printf("%d. %s\n", i+1, f.Name()) // Print with numbering
		}
	}
	if len(plugins) == 0 {
		fmt.Println("No plugins found.")
	}
	fmt.Println()
	return plugins
}

// Function to download and install a plugin
func installPlugin() {
	// URL of the server-side plugin list
	const pluginListURL = "https://dev.lewdlilly.tv/PluginLib/pliuginlib.json" // Update with your server URL

	// Fetch the plugin list
	resp, err := http.Get(pluginListURL)
	if err != nil {
		log.Println("Failed to retrieve plugin list:", err)
		return
	}
	defer resp.Body.Close()

	// Read and decode the JSON list
	var plugins []Plugin
	err = json.NewDecoder(resp.Body).Decode(&plugins)
	if err != nil {
		log.Println("Failed to parse plugin list:", err)
		return
	}

	// Display available plugins
	fmt.Println("\n---- Available Plugins for Download ----")
	for i, plugin := range plugins {
		fmt.Printf("%d. %s\n", i+1, plugin.Name)
	}

	// Ask user to choose a plugin to download
	fmt.Print("Enter the plugin number to download: ")
	reader := bufio.NewReader(os.Stdin)
	pluginChoice, _ := reader.ReadString('\n')
	pluginChoice = strings.TrimSpace(pluginChoice)

	pluginIndex, err := strconv.Atoi(pluginChoice)
	if err != nil || pluginIndex < 1 || pluginIndex > len(plugins) {
		fmt.Println("Invalid plugin number.")
		return
	}

	// Get selected plugin info
	selectedPlugin := plugins[pluginIndex-1]

	// Download the plugin file
	fmt.Println("Downloading plugin:", selectedPlugin.Name)
	resp, err = http.Get(selectedPlugin.URL)
	if err != nil {
		log.Println("Failed to download plugin:", err)
		return
	}
	defer resp.Body.Close()

	// Save the downloaded file to the plugins folder
	filePath := "./plugins/" + selectedPlugin.Name
	out, err := os.Create(filePath)
	if err != nil {
		log.Println("Failed to save plugin:", err)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Println("Failed to write plugin file:", err)
		return
	}

	fmt.Println("Plugin installed successfully:", selectedPlugin.Name)
}

// Function to load and execute a plugin based on the file extension
func loadAndRunPlugin(pluginName string) {
	pluginPath := "./plugins/" + pluginName

	if strings.HasSuffix(pluginName, ".js") {
		loadJSPlugin(pluginPath)
	} else if strings.HasSuffix(pluginName, ".py") {
		loadPythonPlugin(pluginPath)
	} else if strings.HasSuffix(pluginName, ".ps1") {
		loadPowerShellPlugin(pluginPath)
	} else {
		fmt.Println("Unsupported plugin type.")
	}
}

// Function to load JavaScript plugins
func loadJSPlugin(path string) {
	vm := otto.New()

	script, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println("Failed to read JS plugin:", err)
		return
	}

	_, err = vm.Run(string(script))
	if err != nil {
		log.Println("Failed to run JS plugin:", err)
		return
	}

	value, err := vm.Call("run", nil)
	if err != nil {
		log.Println("Failed to call 'run' function in JS plugin:", err)
		return
	}

	fmt.Println("JavaScript Plugin output:", value.String())
}

// Function to load Python plugins
func loadPythonPlugin(path string) {
	cmd := exec.Command("python", path) // Change "python3" to "python" if needed
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Println("Failed to execute Python plugin:", err)
		return
	}

	fmt.Println("Python Plugin output:", string(output))
}

// Function to load and execute PowerShell plugins
func loadPowerShellPlugin(path string) {
	cmd := exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-File", path)

	// Redirect stdin, stdout, and stderr to allow the PowerShell script to handle its own prompts
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Execute the PowerShell command
	err := cmd.Run()
	if err != nil {
		log.Println("Failed to execute PowerShell plugin:", err)
	}
}

func main() {
	cmd := exec.Command("cmd", "/c", "title goCLI Tool by LewdLillyVT") // Replace with your desired title
	err := cmd.Run()
	if err != nil {
		log.Println("Failed to set CMD window title:", err)
		return
	}
	displayWelcomeMessage()

	reader := bufio.NewReader(os.Stdin)
	for {
		displayMenu()

		// Read the user's choice
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			displayHelp()
		case "2":
			listPlugins()
		case "3":
			plugins := listPlugins()
			if len(plugins) > 0 {
				fmt.Print("Enter the plugin number to load and execute: ")
				pluginChoice, _ := reader.ReadString('\n')
				pluginChoice = strings.TrimSpace(pluginChoice)

				// Convert the user input to an integer
				pluginIndex, err := strconv.Atoi(pluginChoice)
				if err != nil || pluginIndex < 1 || pluginIndex > len(plugins) {
					fmt.Println("Invalid plugin number.")
				} else {
					// Load and execute the selected plugin
					loadAndRunPlugin(plugins[pluginIndex-1])
				}
			}
		case "4":
			installPlugin()
		case "0":
			fmt.Println("Exiting goCLI. Goodbye!")
			return
		default:
			fmt.Println("Invalid choice, please select a valid option.")
		}
	}
}
