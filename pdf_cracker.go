package main

// pdf_cracker.go
// Author: Sanoop Thomas
// GitHub: https://github.com/s4n7h0
// Description: A PDF password brute-force tool supporting charset ranges, prefixes, suffixes, and min/max password length.

import (
        "flag"
        "fmt"
        "log"
        "strings"
        "time"

        "github.com/pdfcpu/pdfcpu/pkg/api"
        pdfmodel "github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)


var (
        pdfPath string
        charsetInput string
        minLen int
        maxLen int
        prefix string
        suffix string
        verbose bool
)

func init() {
        flag.StringVar(&pdfPath, "pdf", "", "Path to PDF file")
        flag.StringVar(&charsetInput, "charset", "", "Character set (e.g., a-z, 0-9, abcdef)")
        flag.IntVar(&minLen, "min", 1, "Minimum password length")
        flag.IntVar(&maxLen, "max", 4, "Maximum password length")
        flag.StringVar(&prefix, "prefix", "", "Fixed prefix for all generated passwords")
        flag.StringVar(&suffix, "suffix", "", "Fixed suffix for all generated passwords")
        flag.BoolVar(&verbose, "v", false, "Verbose mode")
}

func main() {
        flag.Parse()

        if pdfPath == "" {
                log.Fatal("Error: -pdf is required")
        }

        if charsetInput == "" {
                log.Fatal("Error: -charset is required")
        }

        // Expand charset ranges
        charset := expandCharset(charsetInput)
        if len(charset) == 0 {
                log.Fatal("Error: Invalid charset")
        }

        fmt.Printf("[*] PDF Cracker\n")
        fmt.Printf("[*] File: %s\n", pdfPath)
        fmt.Printf("[*] Charset: %s\n", charset)
        fmt.Printf("[*] Prefix: '%s' | Suffix: '%s'\n", prefix, suffix)
        fmt.Printf("[*] Length range: %d - %d\n\n", minLen, maxLen)

        start := time.Now()

        // Begin brute force
        for length := minLen; length <= maxLen; length++ {
                if bruteForce(charset, length, tryPassword) {
                        fmt.Printf("\n[+] Password cracked in %s\n", time.Since(start))
                        return
                }
        }

        fmt.Printf("\n[-] Password not found. Time: %s\n", time.Since(start))
}

func expandCharset(input string) []rune {
        var output []rune
        parts := strings.Split(input, ",")

        for _, part := range parts {
                part = strings.TrimSpace(part)

                if strings.Contains(part, "-") && len(part) == 3 {
                        start := rune(part[0])
                        end := rune(part[2])
                        for c := start; c <= end; c++ {
                                output = append(output, c)
                        }
                } else {
                        for _, c := range part {
                                output = append(output, c)
                        }
                }
        }

        return output
}

func bruteForce(charset []rune, length int, callback func(string) bool) bool {
        current := make([]int, length)
        maxIndex := len(charset) - 1

        for {
                // Build password string
                var sb strings.Builder
                sb.WriteString(prefix)
                for _, idx := range current {
                        sb.WriteRune(charset[idx])
                }
                sb.WriteString(suffix)

                password := sb.String()

                if callback(password) {
                        return true
                }

                // Increment character indices (base-n counting)
                for pos := length - 1; pos >= 0; pos-- {
                        if current[pos] < maxIndex {
                                current[pos]++
                                break
                        } else {
                                if pos == 0 {
                                        return false
                                }
                                current[pos] = 0
                        }
                }
        }
}

func tryPassword(password string) bool {
        if verbose {
                fmt.Printf("Trying: %s\n", password)
        }

        // Create a model.Configuration via the model package
        conf := pdfmodel.NewDefaultConfiguration()
        conf.UserPW = password
        conf.OwnerPW = password

        // ValidateFile accepts a *model.Configuration (pass conf here)
        err := api.ValidateFile(pdfPath, conf)
        if err == nil {
                fmt.Printf("\n[+] SUCCESS! Password = %s\n", password)
                return true
        }

        // Validation failed (wrong password or other issue)
        return false
