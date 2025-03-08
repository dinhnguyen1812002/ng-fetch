package system

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"strings"
	"unicode/utf8"
)

// SystemInfo holds all system information
type SystemInfo struct {
	Platform    string
	Kernel      string
	Hostname    string
	CPU         string
	Memory      float64
	Disk        float64
	Uptime      float64
	NetworkSent float64
	NetworkRecv float64
}

// PrintSystemInfo displays system information in an enhanced format
func PrintSystemInfo(noColor bool) error {
	// If noColor is true, disable color output
	color.NoColor = noColor

	// Collect system information
	info, err := collectSystemInfo()
	if err != nil {
		return fmt.Errorf("failed to collect system information: %v", err)
	}

	// Create color schemes
	schemes := createColorSchemes()

	// Print dashboard
	//printDashboardHeader(schemes.header)
	printSystemDetails(info, schemes)
	//printLanguageSection(schemes)

	return nil
}

func collectSystemInfo() (*SystemInfo, error) {
	hostInfo, err := host.Info()
	if err != nil {
		return nil, fmt.Errorf("failed to get host info: %v", err)
	}

	cpuInfo, err := cpu.Info()
	if err != nil {
		return nil, fmt.Errorf("failed to get CPU info: %v", err)
	}

	cpuCount, err := cpu.Counts(true)
	if err != nil {
		return nil, fmt.Errorf("failed to get CPU count: %v", err)
	}

	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return nil, fmt.Errorf("failed to get memory info: %v", err)
	}

	diskInfo, err := disk.Usage("/")
	if err != nil {
		return nil, fmt.Errorf("failed to get disk info: %v", err)
	}

	netInfo, err := net.IOCounters(false)
	if err != nil {
		return nil, fmt.Errorf("failed to get network info: %v", err)
	}

	return &SystemInfo{
		Platform:    fmt.Sprintf("%s %s", hostInfo.Platform, hostInfo.PlatformVersion),
		Kernel:      hostInfo.KernelVersion,
		Hostname:    hostInfo.Hostname,
		CPU:         fmt.Sprintf("%s (%d cores)", cpuInfo[0].ModelName, cpuCount),
		Memory:      float64(memInfo.Total) / (1 << 30),
		Disk:        float64(diskInfo.Total) / (1 << 30),
		Uptime:      float64(hostInfo.Uptime) / 3600,
		NetworkSent: float64(netInfo[0].BytesSent) / (1 << 20),
		NetworkRecv: float64(netInfo[0].BytesRecv) / (1 << 20),
	}, nil
}

type colorSchemes struct {
	header  *color.Color
	section *color.Color
	value   *color.Color
	border  *color.Color
}

func createColorSchemes() colorSchemes {
	return colorSchemes{
		header:  color.New(color.FgHiGreen, color.Bold),
		section: color.New(color.FgHiBlue, color.Bold),
		value:   color.New(color.FgWhite),
		border:  color.New(color.FgHiBlack, color.Bold),
	}
}

//
//func printDashboardHeader(headerColor *color.Color) {
//	timestamp := time.Now().Format("2006-01-02 15:04:05")
//	borderLine := strings.Repeat("═", 60)
//
//	fmt.Printf("╔%s╗\n", borderLine)
//	fmt.Printf("║ %s ║\n", centerText("SYSTEM INFORMATION DASHBOARD", 58))
//	fmt.Printf("║ %s ║\n", centerText(timestamp, 58))
//	fmt.Printf("╠%s╣\n", borderLine)
//}

func getDisplayWidth(s string) int {
	return utf8.RuneCountInString(s)
}

func getPadding(content string, totalWidth int) string {
	displayWidth := getDisplayWidth(content)
	paddingWidth := totalWidth - displayWidth
	if paddingWidth < 0 {
		paddingWidth = 0
	}
	return strings.Repeat(" ", paddingWidth)
}

func printSystemDetails(info *SystemInfo, schemes colorSchemes) {
	const totalWidth = 58 // Total width of the display area

	metrics := []struct {
		icon  string
		name  string
		value interface{}
		unit  string
	}{
		{"\uF17C", "Platform", info.Platform, ""},
		{"\uE70F", "Kernel", info.Kernel, ""},
		{"\uE795", "Hostname", info.Hostname, ""},
		{"\uF4BC", "CPU", info.CPU, ""},
		{"\uF85A", "Memory", info.Memory, "GB"},
		{"\uF0A0", "Disk", info.Disk, "GB"},
		{"\uF43A", "Uptime", info.Uptime, "hours"},
		{"\uF6FF", "Network", fmt.Sprintf("↑%.2f MB | ↓%.2f MB", info.NetworkSent, info.NetworkRecv), ""},
	}

	for _, metric := range metrics {
		var valueStr string
		if v, ok := metric.value.(float64); ok {
			valueStr = fmt.Sprintf("%.2f %s", v, metric.unit)
		} else {
			valueStr = fmt.Sprintf("%v", metric.value)
		}

		line := fmt.Sprintf("%s %s: %s",
			metric.icon,
			schemes.header.Sprint(metric.name),
			schemes.value.Sprint(valueStr))

		padding := getPadding(line, totalWidth)
		fmt.Printf(" %s%s \n", line, padding)
	}
}

//
//func printLanguageSection(schemes colorSchemes) {
//	borderLine := strings.Repeat("═", 60)
//	fmt.Printf("╠%s╣\n", borderLine)
//	fmt.Printf("║ %s ║\n", centerText("INSTALLED PROGRAMMING LANGUAGES", 58))
//	fmt.Printf("╠%s╣\n", borderLine)
//
//	languages := GetProgrammingLanguages()
//	for _, lang := range languages {
//		line := fmt.Sprintf("%s %s: %s",
//			lang.Icon,
//			schemes.section.Sprint(lang.Name),
//			schemes.value.Sprint(lang.Version))
//
//		padding := getPadding(line, 56) // 58 - 2 for the border spaces
//		fmt.Printf("║ %s%s ║\n", line, padding)
//	}
//
//	fmt.Printf("╚%s╝\n", borderLine)
//}

func centerText(text string, width int) string {
	displayWidth := getDisplayWidth(text)
	if displayWidth >= width {
		return text
	}

	padding := width - displayWidth
	leftPad := padding / 2
	rightPad := padding - leftPad

	return strings.Repeat(" ", leftPad) + text + strings.Repeat(" ", rightPad)
}
