package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	defer func() {
		fmt.Print("\nPress enter to exit.\n╰: ")
		reader.ReadLine()
	}()

	if runtime.GOOS != "windows" {
		fmt.Printf("Invalid OS (%s)\n", runtime.GOOS)
		return
	}

	fmt.Println("Getting current user...")
	/* */
	currentUser, err := user.Current()
	if err != nil {
		fmt.Printf("Error while getting current user.\n╰> %s\n", err)
		return
	}
	/* */

	fmt.Println("Fetching EXE...")
	/* */
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error while fetching EXE.\n╰> %s\n", err)
		return
	}
	/* */

	var body []byte
	resp.Body.Read(body)

	fmt.Println("Parsing install path...")
	/* */
	var paths []string
	paths = strings.Split(path, "/")
	for i, v := range paths {
		if v == "~" {
			paths[i] = currentUser.HomeDir
		}
	}
	/* */

	var installPath string = filepath.Join(paths...)

	fmt.Printf("Install at %s?\n╰[Y/n]: ", installPath)

	in, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Error while reading string.\n╰> %s", err)
		return
	}

	if !agrees(in, true) {
		return
	}

	fmt.Println("Opening file...")
	/* */
	file, err := os.OpenFile(installPath, os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Printf("Error while opening file.\n╰> %s\n", err)
		return
	}
	/* */

	fmt.Println("Creating backup Player.exe at ~/.cache/Player.exe...")
	/* */
	dest, err := os.Create(filepath.Join(currentUser.HomeDir, ".cache", "Player.exe"))
	if err != nil {
		fmt.Printf("Error while creating backup Player.exe\n╰> %s\n", err)
		return
	}
	/* */

	fmt.Println("Writing to backup Player.exe...")
	/* */
	if _, err := io.Copy(dest, file); err != nil {
		fmt.Printf("Error while writing to backup Player.exe\n╰> %s\n", err)
		return
	}
	/* */

	fmt.Println("Overwriting Player.exe...")
	/* */
	if _, err := file.Write(body); err != nil {
		fmt.Printf("Error while overwriting Player.exe")
	}
	/* */

	fmt.Println("\nGamma Client is now installed!")
}

func agrees(response string, favor bool) (agrees bool) {
	const (
		yes bool = true
		no  bool = false
	)

	response = strings.ToLower(response)
	response = strings.TrimRight(response, "\n")

	if favor == yes {
		if response == "" || response == "y" || response == "yes" {
			return true
		}
		return false
	}
	if response == "" || response == "n" || response == "no" {
		return false
	}
	return true
}
