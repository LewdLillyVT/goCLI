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
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/robertkrimen/otto" // For JavaScript plugins
)

// Plugin struct to store plugin name and download URL
type Plugin struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// Paths for plugins, logs, and dependencies
var (
	baseDir         = filepath.Join(os.Getenv("LOCALAPPDATA"), "goCLI")
	pluginsDir      = filepath.Join(baseDir, "plugins")
	logsDir         = filepath.Join(baseDir, "logs")
	logFilePrefix   = "errorlog"
	dependenciesDir = filepath.Join(baseDir, "dependencies") // New directory for dependencies
)

// Initialize directories and log file
func initializeAppDirectories() {
	// Create goCLI, plugins, logs, and dependencies directories if they don't exist
	dirs := []string{baseDir, pluginsDir, logsDir, dependenciesDir}
	for _, dir := range dirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err := os.MkdirAll(dir, 0755); err != nil {
				log.Fatalf("Failed to create directory %s: %v", dir, err)
			}
		}
	}

	// Initialize log file
	currentDate := time.Now().Format("2006-01-02")
	logFilePath := filepath.Join(logsDir, fmt.Sprintf("%s_%s.log", logFilePrefix, currentDate))
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to create log file: %v", err)
	}
	log.SetOutput(logFile)
	log.Println("Application started")
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
	fmt.Println("5. Load Plugin Info") // New option for loading plugin info
	fmt.Println("0. Exit")
	fmt.Print("Enter your choice: ")
}

// Function to handle Help/Info command
func displayHelp() {
	fmt.Println("\n---- Help/Info ----")
	fmt.Println("This is a modular CLI tool written in Go.")
	fmt.Println("You can load and execute plugins written in JavaScript, Python, or PowerShell.")
	fmt.Println("Choose an option from the menu to explore the tool’s features.")
	fmt.Println("\n---- Plugin Information ----")
	fmt.Printf("Plugins Folder Location: %s\n", pluginsDir)
	fmt.Println("Place your plugins in the above folder. Supported plugins have the extensions .js, .py, and .ps1.")

	fmt.Println("\n---- Error Reporting ----")
	fmt.Printf("Error logs are stored in the logs folder located at: %s\n", logsDir)
	fmt.Println("If you encounter issues, please follow these steps:")
	fmt.Println("1. Locate the latest error log file in the logs folder (e.g., errorlog_YYYY-MM-DD.log).")
	fmt.Println("2. Copy the contents of the log file and paste it to a service like Pastebin.")
	fmt.Println("3. Visit the GitHub repository for this project and create a new issue.")
	fmt.Println("4. Include the Pastebin link & a description of what caused the error in the issue report for troubleshooting assistance.")
	fmt.Println("\nGitHub Repository for issues: https://github.com/LewdLillyVT/goCLI/issues")
}

// Function to list all plugins in the plugins folder
func listPlugins() []string {
	files, err := ioutil.ReadDir(pluginsDir)
	if err != nil {
		log.Println("Failed to read plugins directory:", err)
		return nil
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
	const pluginListURL = "https://dev.lewdlilly.tv/PluginLib/pliuginlib.json"

	// Create an HTTP client with a timeout
	client := &http.Client{Timeout: 10 * time.Second}

	// Fetch the plugin list
	resp, err := client.Get(pluginListURL)
	if err != nil {
		log.Println("Failed to retrieve plugin list:", err)
		fmt.Println("Error: Could not retrieve plugin list. Check your internet connection or server URL.")
		return
	}
	defer resp.Body.Close()

	// Check if the server returned a successful HTTP status code
	if resp.StatusCode != http.StatusOK {
		log.Printf("Server returned non-200 status: %d %s\n", resp.StatusCode, resp.Status)
		fmt.Println("Error: Unable to retrieve the plugin list. Server returned an error.")
		return
	}

	// Verify the content type to ensure it's JSON
	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		log.Println("Unexpected content type:", contentType)
		fmt.Println("Error: Server response was not JSON. Please check the plugin list URL.")
		return
	}

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Failed to read response body:", err)
		fmt.Println("Error: Failed to read server response.")
		return
	}

	// Attempt to parse the response as JSON
	var plugins []Plugin
	err = json.Unmarshal(body, &plugins)
	if err != nil {
		log.Println("Failed to parse plugin list as JSON:", err)
		fmt.Println("Error: Failed to parse plugin list. Server might be returning an unexpected format.")
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

	// Validate user input
	pluginIndex, err := strconv.Atoi(pluginChoice)
	if err != nil || pluginIndex < 1 || pluginIndex > len(plugins) {
		fmt.Println("Invalid plugin number. Please enter a valid option.")
		return
	}

	// Get selected plugin info
	selectedPlugin := plugins[pluginIndex-1]

	// Ensure the plugins directory exists
	if _, err := os.Stat(pluginsDir); os.IsNotExist(err) {
		err := os.MkdirAll(pluginsDir, os.ModePerm)
		if err != nil {
			log.Println("Failed to create plugins directory:", err)
			fmt.Println("Error: Could not create plugins directory.")
			return
		}
	}

	// Download the plugin file with timeout-enabled client
	fmt.Println("Downloading plugin:", selectedPlugin.Name)
	resp, err = client.Get(selectedPlugin.URL)
	if err != nil {
		log.Println("Failed to download plugin:", err)
		fmt.Println("Error: Could not download the plugin.")
		return
	}
	defer resp.Body.Close()

	// Save the downloaded file to the plugins folder in AppData
	filePath := filepath.Join(pluginsDir, selectedPlugin.Name)
	out, err := os.Create(filePath)
	if err != nil {
		log.Println("Failed to save plugin:", err)
		fmt.Println("Error: Could not save the downloaded plugin.")
		return
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Println("Failed to write plugin file:", err)
		fmt.Println("Error: Could not write the downloaded plugin to disk.")
		return
	}

	fmt.Println("Plugin installed successfully:", selectedPlugin.Name)
}

// Function to load and execute a plugin based on the file extension
func loadAndRunPlugin(pluginName string) {
	pluginPath := filepath.Join(pluginsDir, pluginName)

	switch {
	case strings.HasSuffix(pluginName, ".js"):
		loadJSPlugin(pluginPath)
	case strings.HasSuffix(pluginName, ".py"):
		loadPythonPlugin(pluginPath)
	case strings.HasSuffix(pluginName, ".ps1"):
		loadPowerShellPlugin(pluginPath)
	default:
		log.Println("Unsupported plugin type:", pluginName)
		fmt.Println("Error: Unsupported plugin type.")
	}
}

// Function to load and run a JavaScript plugin
func loadJSPlugin(pluginPath string) {
	// Read the plugin file
	data, err := ioutil.ReadFile(pluginPath)
	if err != nil {
		log.Println("Failed to read JS plugin file:", err)
		fmt.Println("Error: Unable to read the plugin file.")
		return
	}

	// Execute the JavaScript code
	vm := otto.New()
	_, err = vm.Run(string(data))
	if err != nil {
		log.Println("Error executing JS plugin:", err)
		fmt.Println("Error: Failed to execute the plugin.")
		return
	}

	fmt.Println("JavaScript plugin executed successfully.")
}

// Function to load and run a Python plugin
func loadPythonPlugin(pluginPath string) {
	cmd := exec.Command("python", pluginPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Println("Error executing Python plugin:", err)
		fmt.Println("Error: Failed to execute the Python plugin.")
		return
	}

	fmt.Println("Python plugin output:\n", string(output))
}

// Function to load and run a PowerShell plugin
func loadPowerShellPlugin(pluginPath string) {
	cmd := exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-File", pluginPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Println("Error executing PowerShell plugin:", err)
		fmt.Println("Error: Failed to execute the PowerShell plugin.")
		return
	}

	fmt.Println("PowerShell plugin output:\n", string(output))
}

// Function to load and display plugin information
func loadPluginInfo() {
	files, err := ioutil.ReadDir(pluginsDir)
	if err != nil {
		log.Println("Failed to read plugins directory:", err)
		return
	}

	// Map to store plugin info
	plugins := make([]string, 0)

	for _, f := range files {
		if !f.IsDir() && (strings.HasSuffix(f.Name(), ".js") || strings.HasSuffix(f.Name(), ".py") || strings.HasSuffix(f.Name(), ".ps1")) {
			plugins = append(plugins, f.Name())
		}
	}

	// Display the list of available plugins
	if len(plugins) == 0 {
		fmt.Println("No plugins found.")
		return
	}

	fmt.Println("\n---- Available Plugins ----")
	for i, name := range plugins {
		fmt.Printf("%d. %s\n", i+1, name) // Print with numbering
	}

	// Ask user to choose a plugin to view its information
	fmt.Print("Enter the plugin number to view its information: ")
	reader := bufio.NewReader(os.Stdin)
	pluginChoice, _ := reader.ReadString('\n')
	pluginChoice = strings.TrimSpace(pluginChoice)

	// Validate user input
	pluginIndex, err := strconv.Atoi(pluginChoice)
	if err != nil || pluginIndex < 1 || pluginIndex > len(plugins) {
		fmt.Println("Invalid plugin number. Please enter a valid option.")
		return
	}

	// Get selected plugin info
	selectedPlugin := plugins[pluginIndex-1]
	info, err := readPluginInfo(filepath.Join(pluginsDir, selectedPlugin))
	if err != nil {
		fmt.Printf("Failed to read info for %s: %v\n", selectedPlugin, err)
		return
	}

	// Display the plugin information
	fmt.Printf("\n---- Information for %s ----\n", selectedPlugin)
	fmt.Println(info)
}

// readPluginInfo reads the first three lines of comments from a plugin file
func readPluginInfo(pluginPath string) (string, error) {
	file, err := os.Open(pluginPath)
	if err != nil {
		return "", fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	var commentLines []string
	scanner := bufio.NewScanner(file)

	// Read the file line by line
	for scanner.Scan() {
		line := scanner.Text()

		// Check if the line is a comment
		if isCommentLine(line, pluginPath) {
			// Append the comment line to our list
			commentLines = append(commentLines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading file: %w", err)
	}

	// Display up to the last three comments, in the order they appear in the file
	if len(commentLines) > 3 {
		commentLines = commentLines[len(commentLines)-3:]
	}

	if len(commentLines) > 0 {
		return strings.Join(commentLines, "\n"), nil
	}

	return "No plugin information found.", nil
}

// isCommentLine checks if a line is a comment based on the file type
func isCommentLine(line string, pluginPath string) bool {
	line = strings.TrimSpace(line)             // Trim whitespace
	if strings.HasSuffix(pluginPath, ".ps1") { // PowerShell
		return strings.HasPrefix(line, "#")
	} else if strings.HasSuffix(pluginPath, ".py") { // Python
		return strings.HasPrefix(line, "#")
	} else if strings.HasSuffix(pluginPath, ".js") { // JavaScript
		return strings.HasPrefix(line, "//")
	}
	return false
}

// Main function
func main() {
	initializeAppDirectories()
	cmd := exec.Command("cmd", "/c", "title goCLI Tool by LewdLillyVT") // Replace with your desired title
	err := cmd.Run()
	if err != nil {
		log.Println("Failed to set CMD window title:", err)
		return
	}
	displayWelcomeMessage()

	for {
		displayMenu()
		reader := bufio.NewReader(os.Stdin)
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			displayHelp()
		case "2":
			listPlugins()
		case "3":
			fmt.Print("Enter the plugin name to load: ")
			pluginName, _ := reader.ReadString('\n')
			loadAndRunPlugin(strings.TrimSpace(pluginName))
		case "4":
			installPlugin()
		case "5":
			loadPluginInfo()
		case "0":
			fmt.Println("Exiting the application. Goodbye!")
			return
		default:
			fmt.Println("Invalid choice. Please select again.")
		}
	}
}
