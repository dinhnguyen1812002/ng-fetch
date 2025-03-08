package ascii

import (
	"fmt"
	"os"
	"path/filepath"
)

func PrintASCIIArt(filename string) {
	path := filepath.Join("ascii", "assets", filename+".txt")
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Error loading ASCII art:", err)
		return
	}
	fmt.Println(string(data))

}

//func PrintASCIIArt(filename string) {
//	// Get the directory of the currently running executable
//	execPath, err := os.Executable()
//	if err != nil {
//		fmt.Println("Error getting executable path:", err)
//		return
//	}
//
//	// Construct the absolute path to the ASCII art file
//	execDir := filepath.Dir(execPath)
//	path := filepath.Join(execDir, "ascii", "assets", filename+".txt")
//
//	// Read and display the ASCII art
//	data, err := os.ReadFile(path)
//	if err != nil {
//		fmt.Println("Error loading ASCII art:", err)
//		return
//	}
//	fmt.Println(string(data))
//}
