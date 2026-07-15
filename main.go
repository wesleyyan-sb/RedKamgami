package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/klauspost/compress/zstd"
)

var (

	colorReset  string
	colorRed    string
	colorGreen  string
	colorYellow string
	colorBlue   string
	colorPurple string
	colorCyan   string
	colorWhite  string
	colorBold   string
	colorsEnabled bool
)

func initColors() {

	if runtime.GOOS == "windows" {

		colorsEnabled = true
	} else {

		colorsEnabled = true
	}

	if colorsEnabled {
		colorReset = "\033[0m"
		colorRed = "\033[31m"
		colorGreen = "\033[32m"
		colorYellow = "\033[33m"
		colorBlue = "\033[34m"
		colorPurple = "\033[35m"
		colorCyan = "\033[36m"
		colorWhite = "\033[37m"
		colorBold = "\033[1m"
	} else {

		colorReset = ""
		colorRed = ""
		colorGreen = ""
		colorYellow = ""
		colorBlue = ""
		colorPurple = ""
		colorCyan = ""
		colorWhite = ""
		colorBold = ""
	}
}

const (
	version         = "2.0.0"
	configFile      = "mirrors.json"
	wordlistVersion = "version.txt"
	logFile         = "launcher.log"
	timeout         = 30 * time.Second
	updateCheckURL  = "https://raw.githubusercontent.com/your-repo/wordlist-update/main/update.json"
	author          = "https://github.com/wesleyyan-sb"
	license         = "CC-0.1"
)

type Mirror struct {
	Name     string `json:"name"`
	URL      string `json:"url"`
	Priority int    `json:"priority"`
}

type UpdateInfo struct {
	Version     string   `json:"version"`
	Size        int64    `json:"size"`
	ReleaseDate string   `json:"release_date"`
	Description string   `json:"description"`
	Mirrors     []Mirror `json:"mirrors"`
}

type Logger struct {
	file *os.File
	mu   sync.Mutex
}

type Launcher struct {
	config     []Mirror
	updateInfo *UpdateInfo
	logger     *Logger
	ctx        context.Context
	cancel     context.CancelFunc
}

func main() {

	initColors()

	ctx, cancel := context.WithCancel(context.Background())
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	go func() {
		<-sigChan
		cancel()
	}()

	launcher := &Launcher{
		ctx:    ctx,
		cancel: cancel,
	}

	if err := launcher.initLogger(); err != nil {
		fmt.Printf("Error initializing logger: %v\n", err)
		os.Exit(1)
	}
	defer launcher.logger.Close()

	if err := launcher.loadConfig(); err != nil {
		launcher.logger.Log("ERROR", "Failed to load configuration: "+err.Error())
		fmt.Printf("Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	launcher.showHeader()

	launcher.showDisclaimer()

	for {
		launcher.showMenu()
		choice := launcher.getMenuChoice()

		switch choice {
		case 1:
			launcher.downloadWordlist()
		case 2:
			launcher.splitWordlist()
		case 3:
			launcher.removeDuplicates()
		case 4:
			launcher.countLines()
		case 5:
			launcher.mergeWordlists()
		case 6:
			launcher.estimateProbability()
		case 7:
			launcher.compressWordlist()
		case 8:
			launcher.consolidateWordlists()
		case 9:
			fmt.Println("\n👋 Goodbye!")
			return
		default:
			fmt.Println("Invalid option. Please try again.")
		}

		if choice != 9 {
			fmt.Println("\nPress Enter to continue...")
			bufio.NewReader(os.Stdin).ReadBytes('\n')
		}
	}
}

func (l *Launcher) initLogger() error {
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	l.logger = &Logger{file: file}
	return nil
}

func (lg *Logger) Log(level, message string) {
	lg.mu.Lock()
	defer lg.mu.Unlock()
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logLine := fmt.Sprintf("[%s] [%s] %s\n", timestamp, level, message)
	lg.file.WriteString(logLine)
}

func (lg *Logger) Close() error {
	return lg.file.Close()
}

func (l *Launcher) loadConfig() error {

	if _, err := os.Stat(configFile); os.IsNotExist(err) {

		defaultConfig := []Mirror{
			{Name: "Google Drive", URL: "https://drive.google.com/file/d/1YWlEzLnMiRXtQmihPRB2ZJnZUV2CJ3VF/view?usp=sharing", Priority: 1},
			{Name: "MEGA", URL: "https://mega.nz/file/YOUR_FILE_KEY", Priority: 2},
			{Name: "MediaFire", URL: "https://www.mediafire.com/file/nh64tkn5gb2uown/redkamgami.txt.zst/file", Priority: 3},
			{Name: "OneDrive", URL: "https://onedrive.live.com/download?cid=YOUR_CID&resid=YOUR_RESID", Priority: 4},
			{Name: "TeraBox", URL: "https://terabox.com/sharing/link?surl=YOUR_LINK", Priority: 5},
		}

		data, err := json.MarshalIndent(defaultConfig, "", "  ")
		if err != nil {
			return err
		}
		if err := os.WriteFile(configFile, data, 0644); err != nil {
			return err
		}
		l.config = defaultConfig
		return nil
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, &l.config); err != nil {
		return err
	}

	return nil
}

func (l *Launcher) showHeader() {
	fmt.Printf("%s%s", colorBold, colorCyan)
	fmt.Println("╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║                                                            ║")
	fmt.Println("║" + colorReset + colorBold + colorGreen + "         ███╗   ██╗███████╗██╗  ██╗██╗   ██╗██╗███╗   ███╗         " + colorReset + colorBold + colorCyan + "║")
	fmt.Println("║" + colorReset + colorBold + colorGreen + "         ████╗  ██║██╔════╝╚██╗██╔╝██║   ██║██║████╗ ████║         " + colorReset + colorBold + colorCyan + "║")
	fmt.Println("║" + colorReset + colorBold + colorGreen + "         ██╔██╗ ██║█████╗   ╚███╔╝ ██║   ██║██║██╔████╔██║         " + colorReset + colorBold + colorCyan + "║")
	fmt.Println("║" + colorReset + colorBold + colorGreen + "         ██║╚██╗██║██╔══╝   ██╔██╗ ██║   ██║██║██║╚██╔╝██║         " + colorReset + colorBold + colorCyan + "║")
	fmt.Println("║" + colorReset + colorBold + colorGreen + "         ██║ ╚████║███████╗██╔╝ ██╗╚██████╔╝██║██║ ╚═╝ ██║         " + colorReset + colorBold + colorCyan + "║")
	fmt.Println("║" + colorReset + colorBold + colorGreen + "         ╚═╝  ╚═══╝╚══════╝╚═╝  ╚═╝ ╚═════╝ ╚═╝╚═╝     ╚═╝         " + colorReset + colorBold + colorCyan + "║")
	fmt.Println("║                                                            ║")
	fmt.Println("║" + colorReset + colorBold + colorYellow + "                    RED KAMGAMI LAUNCHER                   " + colorReset + colorBold + colorCyan + "║")
	fmt.Println("║                                                            ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝" + colorReset)
	fmt.Println()
	fmt.Printf("%sVersion:%s %s%s%s\n", colorBold, colorReset, colorGreen, version, colorReset)
	fmt.Printf("%sAuthor:%s %s%s%s\n", colorBold, colorReset, colorCyan, author, colorReset)
	fmt.Printf("%sLicense:%s %s%s%s\n", colorBold, colorReset, colorPurple, license, colorReset)
	fmt.Println()
}

func (l *Launcher) showDisclaimer() {
	fmt.Printf("%s%s⚠️  LEGAL DISCLAIMER%s\n", colorBold, colorYellow, colorReset)
	fmt.Printf("%s", colorYellow)
	fmt.Println("─────────────────────────────────────────────────────────────")
	fmt.Println("This tool is for authorized security testing, research, and")
	fmt.Println("educational purposes only. All passwords are generated using")
	fmt.Println("the BlackKamgen generator and existing wordlists from GitHub.")
	fmt.Println()
	fmt.Println("The author (https://github.com/wesleyyan-sb) is NOT responsible")
	fmt.Println("for any misuse or illegal use of this software.")
	fmt.Println()
	fmt.Println("By using this tool, you agree to use it responsibly and in")
	fmt.Println("compliance with applicable laws and regulations.")
	fmt.Println("─────────────────────────────────────────────────────────────")
	fmt.Printf("%s\n", colorReset)
	fmt.Println()
}

func (l *Launcher) showMenu() {
	fmt.Printf("%s%s╔════════════════════════════════════════════════════════════╗%s\n", colorBold, colorCyan, colorReset)
	fmt.Printf("%s%s║%s                        MAIN MENU                        %s%s║%s\n", colorBold, colorCyan, colorReset, colorBold, colorCyan, colorReset)
	fmt.Printf("%s%s╠════════════════════════════════════════════════════════════╣%s\n", colorBold, colorCyan, colorReset)
	fmt.Printf("%s%s║%s  %s[1]%s Download RedKamgami                                   %s%s║%s\n", colorBold, colorCyan, colorReset, colorGreen, colorReset, colorBold, colorCyan, colorReset)
	fmt.Printf("%s%s║%s  %s[2]%s Split RedKamgami                                      %s%s║%s\n", colorBold, colorCyan, colorReset, colorGreen, colorReset, colorBold, colorCyan, colorReset)
	fmt.Printf("%s%s║%s  %s[3]%s Remove Duplicates                                    %s%s║%s\n", colorBold, colorCyan, colorReset, colorGreen, colorReset, colorBold, colorCyan, colorReset)
	fmt.Printf("%s%s║%s  %s[4]%s Count Lines                                          %s%s║%s\n", colorBold, colorCyan, colorReset, colorGreen, colorReset, colorBold, colorCyan, colorReset)
	fmt.Printf("%s%s║%s  %s[5]%s Merge Wordlists                                      %s%s║%s\n", colorBold, colorCyan, colorReset, colorGreen, colorReset, colorBold, colorCyan, colorReset)
	fmt.Printf("%s%s║%s  %s[6]%s Estimate Probability                                  %s%s║%s\n", colorBold, colorCyan, colorReset, colorGreen, colorReset, colorBold, colorCyan, colorReset)
	fmt.Printf("%s%s║%s  %s[7]%s Compress File                                         %s%s║%s\n", colorBold, colorCyan, colorReset, colorGreen, colorReset, colorBold, colorCyan, colorReset)
	fmt.Printf("%s%s║%s  %s[8]%s Consolidate All Wordlists                             %s%s║%s\n", colorBold, colorCyan, colorReset, colorGreen, colorReset, colorBold, colorCyan, colorReset)
	fmt.Printf("%s%s║%s  %s[9]%s Exit                                                 %s%s║%s\n", colorBold, colorCyan, colorReset, colorRed, colorReset, colorBold, colorCyan, colorReset)
	fmt.Printf("%s%s╚════════════════════════════════════════════════════════════╝%s\n", colorBold, colorCyan, colorReset)
	fmt.Println()
}

func (l *Launcher) getMenuChoice() int {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Select an option (1-9): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		choice, err := strconv.Atoi(input)
		if err != nil || choice < 1 || choice > 9 {
			fmt.Println("Invalid option. Please enter a number between 1 and 9.")
			continue
		}
		return choice
	}
}

func (l *Launcher) downloadWordlist() {

	if err := l.checkUpdates(); err != nil {
		l.logger.Log("ERROR", "Failed to check for updates: "+err.Error())
		fmt.Printf("Error checking for updates: %v\n", err)
		if !l.promptContinue("Do you want to continue with the current version? (Y/N)") {
			return
		}
	}

	if l.updateInfo != nil {
		currentVersion := l.getCurrentVersion()
		if currentVersion == l.updateInfo.Version {
			fmt.Println("\n✅ You already have the latest version!")
			return
		}
	}

	l.showMirrors()

	mirror := l.selectMirror()
	if mirror == nil {
		fmt.Println("No mirror selected.")
		return
	}

	if err := l.openInBrowser(mirror.URL); err != nil {
		l.logger.Log("ERROR", "Failed to open browser: "+err.Error())
		fmt.Printf("Error opening browser: %v\n", err)
		fmt.Printf("\n📋 Please manually open this URL in your browser:\n%s\n", mirror.URL)
	}

	if l.promptContinue("\nAfter downloading the file, place it in the current directory as 'redkamgami.txt.zst'.\nHas the download been completed? (Y/N)") {

		fmt.Println("\n📦 Decompressing RedKamgami...")
		if err := l.decompressFile("redkamgami.txt.zst", "redkamgami.txt"); err != nil {
			l.logger.Log("ERROR", "Failed to decompress RedKamgami: "+err.Error())
			fmt.Printf("❌ Error decompressing RedKamgami: %v\n", err)
			return
		}
		fmt.Println("✅ RedKamgami decompressed successfully!")
		

		if err := os.Remove("redkamgami.txt.zst"); err != nil {
			fmt.Printf("⚠️  Warning: Could not remove compressed file: %v\n", err)
		}

		if l.updateInfo != nil {
			if err := os.WriteFile(wordlistVersion, []byte(l.updateInfo.Version), 0644); err != nil {
				l.logger.Log("ERROR", "Failed to save version: "+err.Error())
				fmt.Printf("Error saving version: %v\n", err)
			} else {
				fmt.Println("✅ Version saved successfully!")
				l.logger.Log("INFO", fmt.Sprintf("Version updated to: %s", l.updateInfo.Version))
			}
		}
	}
}

func (l *Launcher) checkUpdates() error {
	fmt.Println("Checking for updates...")

	client := &http.Client{Timeout: timeout}

	resp, err := client.Get(updateCheckURL)
	if err != nil {
		return fmt.Errorf("error checking for updates: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned status: %d", resp.StatusCode)
	}

	var updateInfo UpdateInfo
	if err := json.NewDecoder(resp.Body).Decode(&updateInfo); err != nil {
		return fmt.Errorf("error decoding information: %v", err)
	}

	l.updateInfo = &updateInfo

	currentVersion := l.getCurrentVersion()

	if currentVersion == "" {
		fmt.Println("No installed version found.")
		fmt.Printf("Available version: %s\n", updateInfo.Version)
		fmt.Printf("Size: %s\n", formatSize(updateInfo.Size))
		fmt.Printf("Release date: %s\n", updateInfo.ReleaseDate)
		if updateInfo.Description != "" {
			fmt.Printf("Description: %s\n", updateInfo.Description)
		}
	} else if currentVersion != updateInfo.Version {
		fmt.Printf("\n🔔 New version found!\n")
		fmt.Printf("Current: %s\n", currentVersion)
		fmt.Printf("New: %s\n", updateInfo.Version)
		fmt.Printf("Size: %s\n", formatSize(updateInfo.Size))
		fmt.Printf("Release date: %s\n", updateInfo.ReleaseDate)
		if updateInfo.Description != "" {
			fmt.Printf("Description: %s\n", updateInfo.Description)
		}
	} else {
		fmt.Println("✅ You already have the latest version.")
		return nil
	}

	return nil
}

func (l *Launcher) getCurrentVersion() string {
	if data, err := os.ReadFile(wordlistVersion); err == nil {
		return strings.TrimSpace(string(data))
	}
	return ""
}

func (l *Launcher) showMirrors() {
	fmt.Println("\n📋 Available mirrors:")
	fmt.Println("---------------------")

	mirrors := l.config
	if l.updateInfo != nil && len(l.updateInfo.Mirrors) > 0 {
		mirrors = l.updateInfo.Mirrors
	}

	for i, mirror := range mirrors {
		fmt.Printf("%d. %s\n", i+1, mirror.Name)
	}
	fmt.Println("---------------------")
}

func (l *Launcher) selectMirror() *Mirror {
	mirrors := l.config
	if l.updateInfo != nil && len(l.updateInfo.Mirrors) > 0 {
		mirrors = l.updateInfo.Mirrors
	}

	if len(mirrors) == 0 {
		fmt.Println("No mirrors available.")
		return nil
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("\nSelect mirror number (or 'q' to quit): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "q" || input == "Q" {
			return nil
		}

		var index int
		if _, err := fmt.Sscanf(input, "%d", &index); err != nil {
			fmt.Println("Invalid input. Please enter a number.")
			continue
		}

		if index < 1 || index > len(mirrors) {
			fmt.Printf("Invalid number. Please enter 1 to %d.\n", len(mirrors))
			continue
		}

		return &mirrors[index-1]
	}
}

func (l *Launcher) openInBrowser(url string) error {
	fmt.Printf("\n🌐 Opening URL in your browser: %s\n", url)

	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	if err != nil {
		return fmt.Errorf("failed to open browser: %v", err)
	}

	fmt.Println("✅ Browser opened. Please download the file and save it as 'wordlist.txt' in the current directory.")
	return nil
}

func (l *Launcher) promptContinue(message string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(message)
	response, _ := reader.ReadString('\n')
	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes"
}

func (l *Launcher) showLoading(message string, duration time.Duration) {
	spinner := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	start := time.Now()
	i := 0
	
	for time.Since(start) < duration {
		fmt.Printf("\r%s%s%s %s%s", colorCyan, spinner[i%len(spinner)], colorReset, colorBold, message)
		time.Sleep(100 * time.Millisecond)
		i++
	}
	fmt.Printf("\r%s✓%s %s%s\n", colorGreen, colorReset, colorBold, message)
}

func (l *Launcher) showProgressBar(message string, current, total int) {
	barWidth := 40
	progress := float64(current) / float64(total)
	filled := int(progress * float64(barWidth))
	
	bar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)
	percentage := int(progress * 100)
	
	fmt.Printf("\r%s%s%s [%s] %d%% (%d/%d)", colorBold, colorCyan, message, colorReset+colorGreen+bar+colorReset, percentage, current, total)
	if current == total {
		fmt.Println()
	}
}

func (l *Launcher) printSuccess(message string) {
	fmt.Printf("%s✓%s %s%s\n", colorGreen, colorReset, colorBold, message)
}

func (l *Launcher) printError(message string) {
	fmt.Printf("%s✗%s %s%s\n", colorRed, colorReset, colorBold, message)
}

func (l *Launcher) printInfo(message string) {
	fmt.Printf("%sℹ%s %s%s\n", colorCyan, colorReset, colorBold, message)
}

func (l *Launcher) printWarning(message string) {
	fmt.Printf("%s⚠%s %s%s\n", colorYellow, colorReset, colorBold, message)
}

func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func (l *Launcher) splitWordlist() {
	wordlistPath := "redkamgami.txt"
	

	if _, err := os.Stat(wordlistPath); os.IsNotExist(err) {
		l.printError("redkamgami.txt not found in current directory.")
		return
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%sEnter the number of parts to split into:%s ", colorBold, colorReset)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	
	numParts, err := strconv.Atoi(input)
	if err != nil || numParts < 1 {
		l.printError("Invalid number. Please enter a positive integer.")
		return
	}

	l.printInfo("Counting lines...")
	totalLines, err := l.countLinesInFile(wordlistPath)
	if err != nil {
		l.printError(fmt.Sprintf("Error counting lines: %v", err))
		return
	}

	linesPerPart := totalLines / numParts
	remainder := totalLines % numParts

	fmt.Printf("\n%sTotal lines:%s %d\n", colorBold, colorReset, totalLines)
	fmt.Printf("%sLines per part:%s %d\n", colorBold, colorReset, linesPerPart)
	fmt.Printf("%sSplitting into %d parts...%s\n", colorBold, numParts, colorReset)

	file, err := os.Open(wordlistPath)
	if err != nil {
		l.printError(fmt.Sprintf("Error opening file: %v", err))
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	partNum := 1
	lineCount := 0
	currentPartLines := linesPerPart
	if partNum <= remainder {
		currentPartLines++
	}

	var partFile *os.File
	var partWriter *bufio.Writer

	for scanner.Scan() {

		if partFile == nil {
			partFileName := fmt.Sprintf("redkamgami_part_%d.txt", partNum)
			partFile, err = os.Create(partFileName)
			if err != nil {
				l.printError(fmt.Sprintf("Error creating part file: %v", err))
				return
			}
			partWriter = bufio.NewWriter(partFile)
			l.printInfo(fmt.Sprintf("Created: %s", partFileName))
		}

		partWriter.WriteString(scanner.Text() + "\n")
		lineCount++

		l.showProgressBar("Splitting", lineCount, totalLines)

		if lineCount >= currentPartLines {
			partWriter.Flush()
			partFile.Close()
			partFile = nil
			partWriter = nil
			lineCount = 0
			partNum++
			if partNum <= numParts {
				currentPartLines = linesPerPart
				if partNum <= remainder {
					currentPartLines++
				}
			}
		}
	}

	if partFile != nil {
		partWriter.Flush()
		partFile.Close()
	}

	if err := scanner.Err(); err != nil {
		l.printError(fmt.Sprintf("Error reading file: %v", err))
		return
	}

	l.printSuccess(fmt.Sprintf("RedKamgami split successfully into %d parts!", numParts))
	l.logger.Log("INFO", fmt.Sprintf("RedKamgami split into %d parts", numParts))
}

func (l *Launcher) removeDuplicates() {
	wordlistPath := "redkamgami.txt"
	

	if _, err := os.Stat(wordlistPath); os.IsNotExist(err) {
		l.printError("redkamgami.txt not found in current directory.")
		return
	}

	l.printInfo("Removing duplicates... This may take a while for large files.")

	file, err := os.Open(wordlistPath)
	if err != nil {
		l.printError(fmt.Sprintf("Error opening file: %v", err))
		return
	}
	defer file.Close()

	uniqueLines := make(map[string]bool)
	var lines []string
	scanner := bufio.NewScanner(file)
	totalLines := 0
	
	for scanner.Scan() {
		line := scanner.Text()
		totalLines++
		if !uniqueLines[line] {
			uniqueLines[line] = true
			lines = append(lines, line)
		}
		

		if totalLines%10000 == 0 {
			l.showProgressBar("Processing", totalLines, totalLines)
		}
	}

	if err := scanner.Err(); err != nil {
		l.printError(fmt.Sprintf("Error reading file: %v", err))
		return
	}

	outputFile, err := os.Create(wordlistPath)
	if err != nil {
		l.printError(fmt.Sprintf("Error creating output file: %v", err))
		return
	}
	defer outputFile.Close()

	writer := bufio.NewWriter(outputFile)
	for _, line := range lines {
		writer.WriteString(line + "\n")
	}
	writer.Flush()

	duplicatesRemoved := totalLines - len(lines)
	l.printSuccess(fmt.Sprintf("Duplicates removed! Removed: %d, Remaining: %d", duplicatesRemoved, len(lines)))
	l.logger.Log("INFO", fmt.Sprintf("Duplicates removed. Total unique lines: %d", len(lines)))
}

func (l *Launcher) countLines() {
	wordlistPath := "redkamgami.txt"
	

	if _, err := os.Stat(wordlistPath); os.IsNotExist(err) {
		l.printError("redkamgami.txt not found in current directory.")
		return
	}

	l.printInfo("Counting lines...")
	totalLines, err := l.countLinesInFile(wordlistPath)
	if err != nil {
		l.printError(fmt.Sprintf("Error counting lines: %v", err))
		return
	}

	fileInfo, _ := os.Stat(wordlistPath)
	fileSize := fileInfo.Size()

	fmt.Printf("\n%s╔════════════════════════════════════════════════════════════╗%s\n", colorBold+colorCyan, colorReset)
	fmt.Printf("%s%s║%s                    RED KAMGAMI STATS                    %s%s║%s\n", colorBold, colorCyan, colorReset, colorBold, colorCyan, colorReset)
	fmt.Printf("%s%s╠════════════════════════════════════════════════════════════╣%s\n", colorBold, colorCyan, colorReset)
	fmt.Printf("%s%s║%s  %sTotal lines:%s %s%-50d%s %s║%s\n", colorBold, colorCyan, colorReset, colorBold, colorReset, colorGreen, totalLines, colorReset, colorBold+colorCyan, colorReset)
	fmt.Printf("%s%s║%s  %sFile size:%s  %s%-50s%s %s║%s\n", colorBold, colorCyan, colorReset, colorBold, colorReset, colorGreen, formatSize(fileSize), colorReset, colorBold+colorCyan, colorReset)
	fmt.Printf("%s%s╚════════════════════════════════════════════════════════════╝%s\n", colorBold, colorCyan, colorReset)
	
	l.logger.Log("INFO", fmt.Sprintf("RedKamgami has %d lines", totalLines))
}

func (l *Launcher) countLinesInFile(filePath string) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	count := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		count++
	}
	return count, scanner.Err()
}

func (l *Launcher) mergeWordlists() {
	wordlistPath := "redkamgami.txt"
	

	if _, err := os.Stat(wordlistPath); os.IsNotExist(err) {
		l.printError("redkamgami.txt not found in current directory.")
		return
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%sEnter the path to the wordlist file to merge:%s ", colorBold, colorReset)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if _, err := os.Stat(input); os.IsNotExist(err) {
		l.printError(fmt.Sprintf("File not found: %s", input))
		return
	}

	l.printInfo("Merging wordlists and removing duplicates...")

	currentLines := make(map[string]bool)
	file, err := os.Open(wordlistPath)
	if err != nil {
		l.printError(fmt.Sprintf("Error opening current RedKamgami: %v", err))
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		currentLines[scanner.Text()] = true
	}

	mergeFile, err := os.Open(input)
	if err != nil {
		l.printError(fmt.Sprintf("Error opening merge file: %v", err))
		return
	}
	defer mergeFile.Close()

	newLines := 0
	mergeScanner := bufio.NewScanner(mergeFile)
	for mergeScanner.Scan() {
		line := mergeScanner.Text()
		if !currentLines[line] {
			currentLines[line] = true
			newLines++
		}
	}

	outputFile, err := os.Create(wordlistPath)
	if err != nil {
		l.printError(fmt.Sprintf("Error creating output file: %v", err))
		return
	}
	defer outputFile.Close()

	writer := bufio.NewWriter(outputFile)
	for line := range currentLines {
		writer.WriteString(line + "\n")
	}
	writer.Flush()

	l.printSuccess(fmt.Sprintf("Wordlists merged! Added %d new unique lines.", newLines))
	fmt.Printf("%sTotal unique lines:%s %d\n", colorBold, colorGreen, len(currentLines))
	l.logger.Log("INFO", fmt.Sprintf("Merged wordlist. Added %d new unique lines", newLines))
}

func (l *Launcher) estimateProbability() {
	wordlistPath := "redkamgami.txt"
	

	if _, err := os.Stat(wordlistPath); os.IsNotExist(err) {
		l.printError("redkamgami.txt not found in current directory.")
		return
	}

	fmt.Printf("\n%s╔════════════════════════════════════════════════════════════╗%s\n", colorBold+colorCyan, colorReset)
	fmt.Printf("%s%s║%s              PROBABILITY ESTIMATION TOOL                  %s%s║%s\n", colorBold, colorCyan, colorReset, colorBold, colorCyan, colorReset)
	fmt.Printf("%s%s╚════════════════════════════════════════════════════════════╝%s\n", colorBold, colorCyan, colorReset)
	fmt.Println()
	l.printInfo("This estimation is based on common password patterns and statistics.")
	l.printWarning("Note: This is a rough estimation and not guaranteed accuracy.")
	fmt.Println()

	totalLines, err := l.countLinesInFile(wordlistPath)
	if err != nil {
		l.printError(fmt.Sprintf("Error counting lines: %v", err))
		return
	}

	reader := bufio.NewReader(os.Stdin)
	
	fmt.Printf("%sEnter target's name (optional, press Enter to skip):%s ", colorBold, colorReset)
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)
	
	fmt.Printf("%sEnter target's birthday (YYYY-MM-DD, optional, press Enter to skip):%s ", colorBold, colorReset)
	birthday, _ := reader.ReadString('\n')
	birthday = strings.TrimSpace(birthday)
	
	fmt.Printf("%sEnter known keywords (comma-separated, optional, press Enter to skip):%s ", colorBold, colorReset)
	keywordsInput, _ := reader.ReadString('\n')
	keywordsInput = strings.TrimSpace(keywordsInput)
	keywords := strings.Split(keywordsInput, ",")
	for i, kw := range keywords {
		keywords[i] = strings.TrimSpace(kw)
	}

	baseProbability := math.Min(0.95, float64(totalLines)/10000000.0)
	

	adjustment := 0.0
	
	if name != "" {
		adjustment += 0.15
		l.printSuccess("Target name provided: +15% probability")
	}
	
	if birthday != "" {
		adjustment += 0.10
		l.printSuccess("Target birthday provided: +10% probability")
	}
	
	if len(keywords) > 0 && keywords[0] != "" {
		adjustment += float64(len(keywords)) * 0.05
		l.printSuccess(fmt.Sprintf("%d keywords provided: +%.1f%% probability", len(keywords), float64(len(keywords))*5.0))
	}

	finalProbability := math.Min(0.99, baseProbability + adjustment)
	
	fmt.Printf("\n%s╔════════════════════════════════════════════════════════════╗%s\n", colorBold+colorCyan, colorReset)
	fmt.Printf("%s%s║%s                   ESTIMATION RESULTS                     %s%s║%s\n", colorBold, colorCyan, colorReset, colorBold, colorCyan, colorReset)
	fmt.Printf("%s%s╠════════════════════════════════════════════════════════════╣%s\n", colorBold, colorCyan, colorReset)
	fmt.Printf("%s%s║%s  %sWordlist size:%s %s%-45d%s %s║%s\n", colorBold, colorCyan, colorReset, colorBold, colorReset, colorGreen, totalLines, colorReset, colorBold+colorCyan, colorReset)
	fmt.Printf("%s%s║%s  %sBase probability:%s %s%-40.2f%%%s %s║%s\n", colorBold, colorCyan, colorReset, colorBold, colorReset, colorGreen, baseProbability*100, colorReset, colorBold+colorCyan, colorReset)
	fmt.Printf("%s%s║%s  %sInformation bonus:%s %s%-37.2f%%%s %s║%s\n", colorBold, colorCyan, colorReset, colorBold, colorReset, colorGreen, adjustment*100, colorReset, colorBold+colorCyan, colorReset)
	fmt.Printf("%s%s║%s  %sFinal probability:%s %s%-39.2f%%%s %s║%s\n", colorBold, colorCyan, colorReset, colorBold, colorReset, colorYellow, finalProbability*100, colorReset, colorBold+colorCyan, colorReset)
	
	var status string
	var statusColor string
	if finalProbability > 0.8 {
		status = "🔥 HIGH PROBABILITY"
		statusColor = colorRed
	} else if finalProbability > 0.5 {
		status = "⚠️  MODERATE PROBABILITY"
		statusColor = colorYellow
	} else {
		status = "❄️  LOW PROBABILITY"
		statusColor = colorCyan
	}
	fmt.Printf("%s%s║%s  %sStatus:%s %s%s%-44s%s %s║%s\n", colorBold, colorCyan, colorReset, colorBold, colorReset, statusColor, colorBold, status, colorReset, colorBold+colorCyan, colorReset)
	fmt.Printf("%s%s╚════════════════════════════════════════════════════════════╝%s\n", colorBold, colorCyan, colorReset)
	
	fmt.Println()
	l.printWarning("Disclaimer: This is a statistical estimation and does not guarantee results.")
	l.printInfo("Actual success depends on many factors including password complexity,")
	l.printInfo("security measures, and the quality of your wordlist.")
	
	l.logger.Log("INFO", fmt.Sprintf("Probability estimation: %.2f%%", finalProbability*100))
}

func (l *Launcher) decompressFile(inputPath, outputPath string) error {

	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer inputFile.Close()

	decoder, err := zstd.NewReader(inputFile)
	if err != nil {
		return fmt.Errorf("failed to create zstd decoder: %w", err)
	}
	defer decoder.Close()

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	bufferedWriter := bufio.NewWriterSize(outputFile, 65536)
	defer bufferedWriter.Flush()

	if _, err := bufferedWriter.ReadFrom(decoder); err != nil {
		return fmt.Errorf("failed to decompress data: %w", err)
	}

	return nil
}

func (l *Launcher) compressFile(inputPath, outputPath string) error {

	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer inputFile.Close()

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	encoder, err := zstd.NewWriter(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create zstd encoder: %w", err)
	}
	defer encoder.Close()

	bufferedReader := bufio.NewReaderSize(inputFile, 65536)

	if _, err := bufferedReader.WriteTo(encoder); err != nil {
		return fmt.Errorf("failed to compress data: %w", err)
	}

	return nil
}

func (l *Launcher) compressWordlist() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%sEnter the file path to compress (or press Enter for 'redkamgami.txt'):%s ", colorBold, colorReset)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		input = "redkamgami.txt"
	}

	if _, err := os.Stat(input); os.IsNotExist(err) {
		l.printError(fmt.Sprintf("File not found: %s", input))
		return
	}

	outputPath := input + ".zst"

	if _, err := os.Stat(outputPath); err == nil {
		if !l.promptContinue(fmt.Sprintf("File %s already exists. Overwrite? (Y/N)", outputPath)) {
			return
		}
	}

	l.printInfo("Compressing file...")
	startTime := time.Now()

	if err := l.compressFile(input, outputPath); err != nil {
		l.logger.Log("ERROR", "Failed to compress file: "+err.Error())
		l.printError(fmt.Sprintf("Error compressing file: %v", err))
		return
	}

	elapsed := time.Since(startTime)

	inputInfo, _ := os.Stat(input)
	outputInfo, _ := os.Stat(outputPath)
	
	inputSize := inputInfo.Size()
	outputSize := outputInfo.Size()
	ratio := float64(outputSize) / float64(inputSize) * 100

	l.printSuccess("File compressed successfully!")
	fmt.Printf("\n%s╔════════════════════════════════════════════════════════════╗%s\n", colorBold+colorCyan, colorReset)
	fmt.Printf("%s%s║%s                  COMPRESSION STATS                       %s%s║%s\n", colorBold, colorCyan, colorReset, colorBold, colorCyan, colorReset)
	fmt.Printf("%s%s╠════════════════════════════════════════════════════════════╣%s\n", colorBold, colorCyan, colorReset)
	fmt.Printf("%s%s║%s  %sOriginal size:%s %s%-45s%s %s║%s\n", colorBold, colorCyan, colorReset, colorBold, colorReset, colorGreen, formatSize(inputSize), colorReset, colorBold+colorCyan, colorReset)
	fmt.Printf("%s%s║%s  %sCompressed size:%s %s%-43s%s %s║%s\n", colorBold, colorCyan, colorReset, colorBold, colorReset, colorGreen, formatSize(outputSize), colorReset, colorBold+colorCyan, colorReset)
	fmt.Printf("%s%s║%s  %sCompression ratio:%s %s%-40.1f%%%s %s║%s\n", colorBold, colorCyan, colorReset, colorBold, colorReset, colorYellow, ratio, colorReset, colorBold+colorCyan, colorReset)
	fmt.Printf("%s%s║%s  %sTime elapsed:%s %s%-43.2f s%s %s║%s\n", colorBold, colorCyan, colorReset, colorBold, colorReset, colorGreen, elapsed.Seconds(), colorReset, colorBold+colorCyan, colorReset)
	fmt.Printf("%s%s╚════════════════════════════════════════════════════════════╝%s\n", colorBold, colorCyan, colorReset)
	
	l.logger.Log("INFO", fmt.Sprintf("Compressed %s to %s (%.1f%% ratio)", input, outputPath, ratio))
}

func (l *Launcher) consolidateWordlists() {

	files, err := os.ReadDir(".")
	if err != nil {
		l.printError(fmt.Sprintf("Error reading directory: %v", err))
		return
	}

	var txtFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".txt") {

			if file.Name() != "redkamgami.txt" {
				txtFiles = append(txtFiles, file.Name())
			}
		}
	}

	if len(txtFiles) == 0 {
		l.printWarning("No .txt files found in current directory.")
		return
	}

	fmt.Printf("\n%sFound %d wordlist files to consolidate:%s\n", colorBold, len(txtFiles), colorReset)
	for _, file := range txtFiles {
		fmt.Printf("  - %s\n", file)
	}

	if !l.promptContinue("\nProceed with consolidation? (Y/N)") {
		l.printInfo("Consolidation cancelled.")
		return
	}

	l.printInfo("Reading all wordlists and removing duplicates...")

	uniqueLines := make(map[string]bool)
	totalLines := 0
	duplicatesRemoved := 0

	for i, filename := range txtFiles {
		l.printInfo(fmt.Sprintf("Processing %s (%d/%d)...", filename, i+1, len(txtFiles)))
		
		file, err := os.Open(filename)
		if err != nil {
			l.printError(fmt.Sprintf("Error opening %s: %v", filename, err))
			continue
		}

		scanner := bufio.NewScanner(file)
		fileLines := 0
		fileDuplicates := 0

		for scanner.Scan() {
			line := scanner.Text()
			fileLines++
			totalLines++
			
			if !uniqueLines[line] {
				uniqueLines[line] = true
			} else {
				duplicatesRemoved++
				fileDuplicates++
			}
		}

		file.Close()

		if err := scanner.Err(); err != nil {
			l.printError(fmt.Sprintf("Error reading %s: %v", filename, err))
		}

		fmt.Printf("  Lines: %d, Duplicates: %d\n", fileLines, fileDuplicates)
	}

	l.printInfo("Writing consolidated wordlist to redkamgami.txt...")
	
	outputFile, err := os.Create("redkamgami.txt")
	if err != nil {
		l.printError(fmt.Sprintf("Error creating redkamgami.txt: %v", err))
		return
	}
	defer outputFile.Close()

	writer := bufio.NewWriter(outputFile)
	lineCount := 0
	
	for line := range uniqueLines {
		writer.WriteString(line + "\n")
		lineCount++
		

		if lineCount%10000 == 0 {
			l.showProgressBar("Writing", lineCount, len(uniqueLines))
		}
	}
	
	l.showProgressBar("Writing", lineCount, len(uniqueLines))
	writer.Flush()

	fmt.Printf("\n%s╔════════════════════════════════════════════════════════════╗%s\n", colorBold+colorCyan, colorReset)
	fmt.Printf("%s%s║%s              CONSOLIDATION RESULTS                      %s%s║%s\n", colorBold, colorCyan, colorReset, colorBold, colorCyan, colorReset)
	fmt.Printf("%s%s╠════════════════════════════════════════════════════════════╣%s\n", colorBold, colorCyan, colorReset)
	fmt.Printf("%s%s║%s  %sFiles processed:%s %s%-45d%s %s║%s\n", colorBold, colorCyan, colorReset, colorBold, colorReset, colorGreen, len(txtFiles), colorReset, colorBold+colorCyan, colorReset)
	fmt.Printf("%s%s║%s  %sTotal lines read:%s %s%-43d%s %s║%s\n", colorBold, colorCyan, colorReset, colorBold, colorReset, colorGreen, totalLines, colorReset, colorBold+colorCyan, colorReset)
	fmt.Printf("%s%s║%s  %sDuplicates removed:%s %s%-41d%s %s║%s\n", colorBold, colorCyan, colorReset, colorBold, colorReset, colorYellow, duplicatesRemoved, colorReset, colorBold+colorCyan, colorReset)
	fmt.Printf("%s%s║%s  %sUnique lines written:%s %s%-40d%s %s║%s\n", colorBold, colorCyan, colorReset, colorBold, colorReset, colorGreen, lineCount, colorReset, colorBold+colorCyan, colorReset)
	fmt.Printf("%s%s╚════════════════════════════════════════════════════════════╝%s\n", colorBold, colorCyan, colorReset)

	l.printSuccess("Consolidation completed successfully!")
	l.logger.Log("INFO", fmt.Sprintf("Consolidated %d files into redkamgami.txt. Total unique lines: %d", len(txtFiles), lineCount))
}

