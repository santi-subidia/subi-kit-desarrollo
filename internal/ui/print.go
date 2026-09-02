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
)

// PrintBanner imprime el encabezado visual de devkit.
func PrintBanner() {
	fmt.Println(colorCyan + colorBold + `
   ___            _  ___ _   
  / _ \_____   __/ |/ (_) |_ 
 / / |/ _ \ \ / /   /| | __|
/ /_/ /  __/\ V /   \| | |_ 
\____/ \___| \_/|_|\_\_|\__| (IA Dev-Kit)` + colorReset)
	fmt.Println(colorDim + " Kit de Desarrollo con IA & Reglas de Arquitectura" + colorReset)
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
