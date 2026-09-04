package ui

import (
	"fmt"
	"strings"
)

const (
	colorReset  = "\033[0m"
	colorBold   = "\033[1m"
	colorGreen  = "\033[32m"
	colorCyan   = "\033[36m"
	colorYellow = "\033[33m"
	colorRed    = "\033[31m"
	colorDim    = "\033[2m"
	colorMagenta = "\033[35m"
)

// PrintBanner imprime el encabezado visual de subikit.
func PrintBanner() {
	fmt.Println(colorCyan + colorBold + `
   ______       __     _ __   _ __ 
  / __/ /  __  / /    (_) /__(_) /_
 _\ \/ _ \/ /_/ /    / /  '_/ / __/
/___/____/\____/____/_/_/\_\_/\__/ 
              /_____/              ` + colorReset + colorMagenta + `(v0.5.1)` + colorReset)
	fmt.Println(colorDim + " SubiKit: Dev-Kit Personalizado para Desarrollo con IA" + colorReset)
	fmt.Println()
}

// Success imprime un mensaje de éxito.
func Success(msg string) {
	fmt.Printf("%s[OK]%s %s\n", colorGreen+colorBold, colorReset, msg)
}

// Info imprime un mensaje informativo.
func Info(msg string) {
	fmt.Printf("%s[INFO]%s %s\n", colorCyan, colorReset, msg)
}

// Warn imprime una advertencia.
func Warn(msg string) {
	fmt.Printf("%s[WARN]%s %s\n", colorYellow+colorBold, colorReset, msg)
}

// Error imprime un error.
func Error(msg string) {
	fmt.Printf("%s[ERROR]%s %s\n", colorRed+colorBold, colorReset, msg)
}

// Section imprime un separador de sección.
func Section(title string) {
	fmt.Println()
	fmt.Printf("%s=== %s ===%s\n", colorBold, strings.ToUpper(title), colorReset)
}

// Bullet imprime un elemento de lista con viñeta.
func Bullet(label string, value string) {
	fmt.Printf("  %s•%s %s%s%s: %s\n", colorCyan, colorReset, colorBold, label, colorReset, value)
}
