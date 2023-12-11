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

	fmt.Println("Fetching version...")
	/* */
	var version string
	resp, err := http.Get(versionUrl)
	if err != nil {
		fmt.Printf("Error while fetching version\n╰> %s\n", err)
		version = noVersion
	} else {
		buf, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Error while reading body\n╰> %s\n", err)
			version = noVersion
		} else {
			version = string(buf)
			version = strings.ReplaceAll(version, "\r", "")
			version = strings.ReplaceAll(version, "\n", "")
		}
	}
	/* */

	fmt.Printf("Fetching Player.exe %s...\n", version)
	/* */
	resp, err = http.Get(url)
	if err != nil {
		fmt.Printf("Error while fetching Player.exe\n╰> %s\n", err)
		return
	}
	defer resp.Body.Close()
	/* */

	fmt.Println("Reading body...")
	/* */
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error while reading body.\n╰> %s\n", err)
	}
	/* */

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

	fmt.Printf("Install %s at %s?\n╰[Y/n]: ", version, installPath)
	/* */
	in, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Error while reading string.\n╰> %s", err)
		return
	}

	if !agrees(in, true) {
		return
	}
	/* */

	fmt.Println("Opening file...")
	/* */
	file, err := os.Create(installPath)
	if err != nil {
		fmt.Printf("Error while opening file.\n╰> %s\n", err)
		return
	}
	defer file.Close()
	/* */

	fmt.Println("Creating ~/.cache/ if it doesn't already exist...")
	/* */
	if err := os.MkdirAll(filepath.Join(currentUser.HomeDir, ".cache"), os.ModePerm); err != nil {
		fmt.Printf("Error while creating ~/.cache/\n╰> %s\n", err)
	}
	/* */

	fmt.Println("Creating backup Player.exe at ~/.cache/Player.exe...")
	/* */
	dest, err := os.Create(filepath.Join(currentUser.HomeDir, ".cache", "Player.exe"))
	if err != nil {
		fmt.Printf("Error while creating backup Player.exe\n╰> %s\n", err)
		return
	}
	defer dest.Close()
	/* */

	fmt.Println("Writing to backup Player.exe...")
	/* */
	if _, err := io.Copy(dest, file); err != nil {
		fmt.Printf("Error while writing to backup Player.exe\n╰> %s\n", err)
		return
	}
	/* */

	fmt.Printf("Installing Gamma Client %s...\n", version)
	/* */
	if _, err := file.Write(body); err != nil {
		fmt.Printf("Error while overwriting Player.exe")
	}
	/* */

	fmt.Printf("\nGamma Client %s is now installed!\n", version)
}

func agrees(response string, favor bool) (agrees bool) {
	const (
		yes bool = true
		no  bool = false
	)
	const (
		agree    bool = true
		disagree bool = false
	)

	response = strings.ToLower(response)
	response = strings.TrimSpace(response)
	response = strings.TrimRight(response, "\r\n")

	if favor == yes {
		if response == "" || response == "y" || response == "yes" {
			return agree
		}
		return disagree
	}
	if response == "" || response == "n" || response == "no" {
		return disagree
	}
	return agree
}
