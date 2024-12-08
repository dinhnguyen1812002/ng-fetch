package system

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// Language represents a programming language with its details
type Language struct {
	Icon    string
	Name    string
	Version string
}

// GetProgrammingLanguages returns a detailed list of installed languages
func GetProgrammingLanguages() []Language {
	// Cross-platform language detection
	languages := []Language{
		languageDetector("Python", "🐍", getPythonVersion()),
		languageDetector("Go", "🟢", getGoVersion()),
		languageDetector("Node.js", "🟨", getNodeVersion()),
		languageDetector("Java", "☕", getJavaVersion()),
		languageDetector("Ruby", "💎", getRubyVersion()),
		languageDetector("Rust", "🦀", getRustVersion()),
		languageDetector("PHP", "🐘", getPHPVersion()),
	}

	// Filter out languages with empty versions
	var installedLanguages []Language
	for _, lang := range languages {
		if lang.Version != "" {
			installedLanguages = append(installedLanguages, lang)
		}
	}

	return installedLanguages
}

// Helper function to create a language entry
func languageDetector(name, icon, version string) Language {
	return Language{
		Name:    name,
		Icon:    icon,
		Version: version,
	}
}

// Platform-specific version detection functions
func getPythonVersion() string {
	return runCommand("python", "--version")
}

func getGoVersion() string {
	return runCommand("go", "version")
}

func getNodeVersion() string {
	return runCommand("node", "--version")
}

func getJavaVersion() string {
	return runCommand("java", "-version")
}

func getRubyVersion() string {
	return runCommand("ruby", "--version")
}

func getRustVersion() string {
	return runCommand("rustc", "--version")
}

func getPHPVersion() string {
	return runCommand("php", "--version")
}

// Cross-platform command runner
func runCommand(command string, args ...string) string {
	// Construct the full command with args
	fullCmd := append([]string{command}, args...)

	// Different approaches for Windows and Unix-like systems
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		// On Windows, use cmd.exe to run commands
		winArgs := append([]string{"/c"}, fullCmd...)
		cmd = exec.Command("cmd", winArgs...)
	} else {
		// On Unix-like systems, use sh
		cmd = exec.Command(command, args...)
	}

	// Run the command and capture output
	output, err := cmd.CombinedOutput()
	if err != nil {
		return ""
	}

	// Clean and process the version string
	version := strings.TrimSpace(string(output))

	// Remove command name or additional text
	version = processVersionString(command, version)

	return version
}

// Process version string to extract clean version
func processVersionString(command, version string) string {
	// Different processing for different commands
	switch command {
	case "python", "ruby", "go", "rustc":

		parts := strings.Fields(version)
		if len(parts) > 1 {
			return parts[1]
		}
	case "node":
		// Node.js version starts with 'v'
		return strings.TrimPrefix(version, "v")
	case "java":
		// Java version is more complex
		lines := strings.Split(version, "\n")
		for _, line := range lines {
			if strings.Contains(line, "version") {
				// Extract version in quotes
				parts := strings.Split(line, "\"")
				if len(parts) > 1 {
					return parts[1]
				}
			}
		}
	case "php":
		parts := strings.Fields(version)
		if len(parts) > 0 {
			return parts[0]
		}
	}

	return version
}

// PrintSystemInfo displays system information in a single row with a language table
// PrintSystemInfo displays system information in a well-structured format
// PrintSystemInfo displays system information in an enhanced format
func PrintSystemInfo(noColor bool) {
	// If noColor is true, disable color output
	if noColor {
		color.NoColor = true
	}

	// Collect system information
	hostInfo, _ := host.Info()
	cpuInfo, _ := cpu.Info()
	cpuCount, _ := cpu.Counts(true)
	memInfo, _ := mem.VirtualMemory()
	diskInfo, _ := disk.Usage("/")
	netInfo, _ := net.IOCounters(false)

	// Define color functions
	titleColor := color.New(color.FgHiMagenta, color.Bold).SprintFunc()
	headerColor := color.New(color.FgHiCyan, color.Bold).SprintFunc()
	//valueColor := color.New(color.FgWhite).SprintFunc()
	//infoColor := color.New(color.FgHiGreen).SprintFunc()

	// Print decorative border
	borderLine := strings.Repeat("═", 50)
	fmt.Println(titleColor("╔" + borderLine + "╗"))
	fmt.Printf("%s %s %s\n",
		titleColor("║"),
		titleColor("🌟 SYSTEM INFORMATION DASHBOARD 🌟"),
		titleColor("║"),
	)
	fmt.Println(titleColor("╠" + borderLine + "╣"))

	// System Information Display
	systemInfoItems := []struct {
		icon   string
		header string
		value  string
		color  func(a ...interface{}) string
	}{
		{
			icon:   "🖥️ ",
			header: "Platform",
			value:  fmt.Sprintf("%s %s", hostInfo.Platform, hostInfo.PlatformVersion),
			color:  color.New(color.FgHiRed).SprintFunc(),
		},
		{
			icon:   "🧊 ",
			header: "Kernel",
			value:  hostInfo.KernelVersion,
			color:  color.New(color.FgHiBlue).SprintFunc(),
		},
		{
			icon:   "🏠 ",
			header: "Hostname",
			value:  hostInfo.Hostname,
			color:  color.New(color.FgHiYellow).SprintFunc(),
		},
		{
			icon:   "🧠 ",
			header: "CPU",
			value:  fmt.Sprintf("\t%s (%d cores)", cpuInfo[0].ModelName, cpuCount),
			color:  color.New(color.FgHiGreen).SprintFunc(),
		},
		{
			icon:   "💾 ",
			header: "Memory",
			value:  fmt.Sprintf("\t%.2f GB", float64(memInfo.Total)/(1<<30)),
			color:  color.New(color.FgHiMagenta).SprintFunc(),
		},
		{
			icon:   "📂 ",
			header: "Disk",
			value:  fmt.Sprintf("\t%.2f GB", float64(diskInfo.Total)/(1<<30)),
			color:  color.New(color.FgHiCyan).SprintFunc(),
		},
		{
			icon:   "⏳ ",
			header: "Uptime",
			value:  fmt.Sprintf("\t%.2f hrs", float64(hostInfo.Uptime)/3600),
			color:  color.New(color.FgHiWhite).SprintFunc(),
		},
		{
			icon:   "🌐 ",
			header: "Network",
			value: fmt.Sprintf("\t%.2f MB sent | %.2f MB received",
				float64(netInfo[0].BytesSent)/(1<<20),
				float64(netInfo[0].BytesRecv)/(1<<20)),
			color: color.New(color.FgHiYellow).SprintFunc(),
		},
	}

	// Print system information with enhanced formatting
	for _, item := range systemInfoItems {
		fmt.Printf("%s %s: %s\n",
			headerColor(item.icon+item.header),
			headerColor(":"),
			item.color(item.value),
		)
	}

	// Close information section
	fmt.Println(titleColor("╠" + borderLine + "╣"))

	// Get installed languages
	languages := GetProgrammingLanguages()

	// Create table for languages with improved styling
	fmt.Println(titleColor("║ 🚀 INSTALLED PROGRAMMING LANGUAGES ") + titleColor("║"))
	fmt.Println(titleColor("╠" + borderLine + "╣"))

	// Create table for languages
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Icon", "Language", "Version"})
	table.SetBorder(false)
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
	)

	// Populate table with language information
	for _, lang := range languages {
		table.Append([]string{lang.Icon, lang.Name, lang.Version})
	}

	// Render the table
	table.Render()

	// Close the dashboard
	fmt.Println(titleColor("╚" + borderLine + "╝"))
}
