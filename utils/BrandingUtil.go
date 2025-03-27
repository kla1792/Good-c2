package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gliderlabs/ssh"
)

// Global gradients loaded from gradient.json
var globalGradients map[string]map[string]interface{}

// Initialize gradients during package load
func Init() {
	gradientFile := "./assets/gradient.json"
	file, err := os.ReadFile(gradientFile)
	if err != nil {
		fmt.Printf("Error loading gradient file: %v\n", err)
		globalGradients = make(map[string]map[string]interface{})
		return
	}

	err = json.Unmarshal(file, &globalGradients)
	if err != nil {
		fmt.Printf("Error parsing gradient file: %v\n", err)
		globalGradients = make(map[string]map[string]interface{})
	}
}

// Convert hex color (#RRGGBB) to RGB
func hexToRGB(hex string) (int, int, int, error) {
	if len(hex) != 7 || hex[0] != '#' {
		return 0, 0, 0, fmt.Errorf("invalid hex color: %s", hex)
	}

	r, err := strconv.ParseInt(hex[1:3], 16, 0)
	if err != nil {
		return 0, 0, 0, err
	}

	g, err := strconv.ParseInt(hex[3:5], 16, 0)
	if err != nil {
		return 0, 0, 0, err
	}

	b, err := strconv.ParseInt(hex[5:7], 16, 0)
	if err != nil {
		return 0, 0, 0, err
	}

	return int(r), int(g), int(b), nil
}

func applyGradient(text, gradientName string) string {
	gradient, exists := globalGradients[gradientName]
	if !exists {
		return text
	}

	fromColor, fromExists := gradient["from_color"].(string)
	toColor, toExists := gradient["to_color"].(string)
	background := gradient["background"].(bool)
	if !fromExists || !toExists {
		return text
	}

	// Parse colors
	r1, g1, b1, err := hexToRGB(fromColor)
	if err != nil {
		return text
	}

	r2, g2, b2, err := hexToRGB(toColor)
	if err != nil {
		return text
	}

	// Apply gradient
	var result strings.Builder
	length := len([]rune(text)) // Obsługa wielobajtowych znaków

	for i, char := range []rune(text) {
		t := float64(i) / float64(length-1) // Interpolacja koloru
		r := int(float64(r1) + t*float64(r2-r1))
		g := int(float64(g1) + t*float64(g2-g1))
		b := int(float64(b1) + t*float64(b2-b1))

		if background {
			// Sekwencja ANSI dla koloru tła
			result.WriteString(fmt.Sprintf("\x1b[48;2;%d;%d;%dm%c\u200B", r, g, b, char))
		} else {
			// Sekwencja ANSI dla koloru tekstu
			result.WriteString(fmt.Sprintf("\x1b[38;2;%d;%d;%dm%c\u200B", r, g, b, char))
		}
	}

	// Reset ANSI
	result.WriteString("\x1b[0m")

	return result.String()
}

func Branding(session ssh.Session, filename string, content map[string]interface{}) string {
	name := "./assets/branding/" + filename + ".tfx"
	file, err := os.ReadFile(name)
	if err != nil {
		return ""
	}
	fileContent := string(file)

	// Process functions like <<$sleep(1000)>>
	re := regexp.MustCompile(`<<\$(\w+)\((\d+)\)>>`)
	fileContent = re.ReplaceAllStringFunc(fileContent, func(match string) string {
		submatches := re.FindStringSubmatch(match)
		if len(submatches) < 3 {
			return match
		}

		funcName := submatches[1]
		argStr := submatches[2]

		arg, err := strconv.Atoi(argStr)
		if err != nil {
			return match
		}

		if funcName == "sleep" {
			return fmt.Sprintf("<<SLEEP(%d)>>", arg)
		}

		return match
	})

	// Replace placeholders with content values
	for key, value := range content {
		if v, ok := value.(string); ok {
			fileContent = strings.ReplaceAll(fileContent, "<<$"+key+">>", v)
		}
	}

	// Apply gradients to text wrapped in <gradient>...</gradient>
	gradientRegex := regexp.MustCompile(`<gradient name="([^"]+)">(.+?)</gradient>`)
	fileContent = gradientRegex.ReplaceAllStringFunc(fileContent, func(match string) string {
		submatches := gradientRegex.FindStringSubmatch(match)
		if len(submatches) < 3 {
			return match
		}

		gradientName := submatches[1]
		text := submatches[2]
		return applyGradient(text, gradientName)
	})

	var result strings.Builder
	lastIndex := 0
	sleepRegex := regexp.MustCompile(`<<SLEEP\((\d+)\)>>`)

	// Process SLEEP directives
	for _, match := range sleepRegex.FindAllStringSubmatchIndex(fileContent, -1) {
		result.WriteString(fileContent[lastIndex:match[0]])

		SendMessage(session, result.String(), false)

		durationStr := fileContent[match[2]:match[3]]
		duration, err := strconv.Atoi(durationStr)
		if err != nil {
			SendMessage(session, fileContent[match[0]:match[1]], false)
			continue
		}

		time.Sleep(time.Duration(duration) * time.Millisecond)

		result.Reset()
		lastIndex = match[1]
	}

	result.WriteString(fileContent[lastIndex:])

	return result.String()
}
