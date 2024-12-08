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
