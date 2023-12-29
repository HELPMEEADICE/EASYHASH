package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash/crc32"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func calculateHashes(filePath string) map[string]string {
	md5Hash := md5.New()
	sha1Hash := sha1.New()
	sha224Hash := sha256.New224()
	sha256Hash := sha256.New()
	sha384Hash := sha512.New384()
	sha512Hash := sha512.New()

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
	defer file.Close()

	buffer := make([]byte, 65536)
	for {
		n, err := file.Read(buffer)
		if n > 0 {
			md5Hash.Write(buffer[:n])
			sha1Hash.Write(buffer[:n])
			sha224Hash.Write(buffer[:n])
			sha256Hash.Write(buffer[:n])
			sha384Hash.Write(buffer[:n])
			sha512Hash.Write(buffer[:n])
		}
		if err != nil {
			break
		}
	}

	file.Seek(0, 0)
	crc32Hash := crc32.NewIEEE()
	_, err = crc32Hash.Write([]byte(file.Name()))
	if err != nil {
		fmt.Println("Error calculating CRC32:", err)
		os.Exit(1)
	}
	crc32Value := fmt.Sprintf("%08x", crc32Hash.Sum32())

	hashValues := map[string]string{
		"CRC":    crc32Value,
		"MD5":    hex.EncodeToString(md5Hash.Sum(nil)),
		"SHA1":   hex.EncodeToString(sha1Hash.Sum(nil)),
		"SHA224": hex.EncodeToString(sha224Hash.Sum(nil)),
		"SHA256": hex.EncodeToString(sha256Hash.Sum(nil)),
		"SHA384": hex.EncodeToString(sha384Hash.Sum(nil)),
		"SHA512": hex.EncodeToString(sha512Hash.Sum(nil)),
	}
	return hashValues
}

func getTerminalWidth() int {
	cmd := exec.Command("mode", "con")
	out, err := cmd.Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting console width: %s", err)
		return 80
	}

	re := regexp.MustCompile(`\d+`)
	rs := re.FindAllString(string(out), -1)
	if len(rs) >= 2 {
		width, _ := strconv.Atoi(rs[1])
		return width
	}
	return 80
}

func printSeparator(width int) {
	separator := strings.Repeat("-", width)
	fmt.Printf("\033[36m%s\033[0m\n", separator)
}

func printColored(text string, colorCode int) {
	fmt.Printf("\033[%dm%s\033[0m\n", colorCode, text)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run hash.go <file_path>")
	} else {
		filePath := os.Args[1]
		hashValues := calculateHashes(filePath)
		width := getTerminalWidth()

		printColored("File Hash Data:", 32)
		keys := make([]string, 0, len(hashValues))
		for k := range hashValues {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		printSeparator(width)
		for _, algorithm := range keys {
			fmt.Printf("%s: %s\n", algorithm, hashValues[algorithm])
		}

		printSeparator(width)
	}
}
