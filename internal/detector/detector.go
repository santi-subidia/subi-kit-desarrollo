package detector

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// DetectedStack contiene la información del stack detectado en un directorio.
type DetectedStack struct {
	ProjectName  string
	RootPath     string
	Technologies []string
	Tags         []string
}

type packageJSON struct {
	Name            string            `json:"name"`
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

// Subdirectorios comunes a inspeccionar en proyectos monorepo / multi-capa
var commonSubdirs = []string{
	"frontend", "backend", "client", "server", "web", "api", "app", "src",
	"apps", "packages", "services", "infra", "deploy",
}

// Detect inspecciona el directorio provisto y subdirectorios clave para extraer el stack tecnológico completo.
func Detect(rootDir string) (*DetectedStack, error) {
	absPath, err := filepath.Abs(rootDir)
	if err != nil {
		absPath = rootDir
	}

	projectName := filepath.Base(absPath)
	techs := make(map[string]bool)
	tags := make(map[string]bool)

	// Siempre agregar tags base
	tags["core"] = true
	tags["general"] = true

	// Carpetas a explorar: raíz + subdirectorios relevantes (hasta 2 niveles)
	dirsToInspect := findInspectableDirs(absPath)

	for _, dir := range dirsToInspect {
		inspectDir(dir, techs, tags, &projectName)
	}

	var techList []string
	for t := range techs {
		techList = append(techList, t)
	}
	sort.Strings(techList)

	var tagList []string
	for t := range tags {
		tagList = append(tagList, t)
	}
	sort.Strings(tagList)

	return &DetectedStack{
		ProjectName:  projectName,
		RootPath:     absPath,
		Technologies: techList,
		Tags:         tagList,
	}, nil
}

func findInspectableDirs(root string) []string {
	result := []string{root}
	seen := map[string]bool{root: true}

	entries, err := os.ReadDir(root)
	if err != nil {
		return result
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasPrefix(name, ".") || name == "node_modules" || name == "dist" || name == "build" || name == "bin" || name == "obj" {
			continue
		}

		subPath := filepath.Join(root, name)
		if isRelevantSubdir(name) {
			if !seen[subPath] {
				seen[subPath] = true
				result = append(result, subPath)
			}

			// Nivel 2 para monorepos tipo apps/* o packages/* o backend/*
			if subEntries, err := os.ReadDir(subPath); err == nil {
				for _, subEntry := range subEntries {
					if subEntry.IsDir() && !strings.HasPrefix(subEntry.Name(), ".") && subEntry.Name() != "bin" && subEntry.Name() != "obj" && subEntry.Name() != "node_modules" {
						subSubPath := filepath.Join(subPath, subEntry.Name())
						if !seen[subSubPath] {
							seen[subSubPath] = true
							result = append(result, subSubPath)
						}
					}
				}
			}
		}
	}

	return result
}

func isRelevantSubdir(name string) bool {
	lower := strings.ToLower(name)
	for _, target := range commonSubdirs {
		if lower == target || strings.HasPrefix(lower, target) {
			return true
		}
	}
	return true
}

func inspectDir(dir string, techs map[string]bool, tags map[string]bool, projectName *string) {
	// 1. Inspeccionar package.json
	pkgPath := filepath.Join(dir, "package.json")
	if data, err := os.ReadFile(pkgPath); err == nil {
		var pkg packageJSON
		if err := json.Unmarshal(data, &pkg); err == nil {
			if pkg.Name != "" && (*projectName == "" || *projectName == filepath.Base(filepath.Dir(pkgPath))) {
				// No sobrescribir el projectName raíz a menos que sea más específico
			}

			allDeps := make(map[string]string)
			for k, v := range pkg.Dependencies {
				allDeps[k] = v
			}
			for k, v := range pkg.DevDependencies {
				allDeps[k] = v
			}

			// Detecciones en package.json
			if hasDep(allDeps, "next") {
				techs["Next.js"] = true
				tags["nextjs"] = true
				tags["react"] = true
				tags["frontend"] = true
			}
			if hasDep(allDeps, "react") {
				techs["React"] = true
				tags["react"] = true
				tags["frontend"] = true
			}

			if hasDep(allDeps, "typescript") {
				techs["TypeScript"] = true
				tags["typescript"] = true
				tags["types"] = true
			}

			if hasDep(allDeps, "tailwindcss") || hasDep(allDeps, "@tailwindcss/postcss") {
				techs["Tailwind CSS"] = true
				tags["tailwind"] = true
				tags["styling"] = true
				tags["frontend"] = true
			}

			if hasDep(allDeps, "@supabase/supabase-js") || hasDep(allDeps, "@supabase/ssr") {
				techs["Supabase"] = true
				tags["supabase"] = true
				tags["database"] = true
				tags["postgres"] = true
				tags["backend"] = true
			}

			if hasDep(allDeps, "zod") {
				techs["Zod"] = true
				tags["zod"] = true
			}

			if hasDep(allDeps, "express") || hasDep(allDeps, "fastify") || hasDep(allDeps, "@nestjs/core") || hasDep(allDeps, "hono") {
				techs["Node.js Backend"] = true
				tags["node"] = true
				tags["api"] = true
				tags["backend"] = true
			}

			if hasDep(allDeps, "jest") || hasDep(allDeps, "vitest") {
				techs["Testing (Jest/Vitest)"] = true
				tags["testing"] = true
			}
		}
	}

	// 2. Inspeccionar tsconfig.json
	if fileExists(filepath.Join(dir, "tsconfig.json")) {
		techs["TypeScript"] = true
		tags["typescript"] = true
		tags["types"] = true
	}

	// 3. Inspeccionar .NET / C# (.sln, .slnx, .csproj, .fsproj)
	if entries, err := os.ReadDir(dir); err == nil {
		for _, e := range entries {
			if !e.IsDir() {
				ext := strings.ToLower(filepath.Ext(e.Name()))
				if ext == ".sln" || ext == ".slnx" || ext == ".csproj" || ext == ".fsproj" {
					techs[".NET / C#"] = true
					tags["dotnet"] = true
					tags["csharp"] = true
					tags["backend"] = true
					tags["clean-architecture"] = true
					break
				}
			}
		}
	}

	// 4. Inspeccionar Go
	if fileExists(filepath.Join(dir, "go.mod")) {
		techs["Go"] = true
		tags["go"] = true
		tags["backend"] = true
	}

	// 5. Inspeccionar Python
	if fileExists(filepath.Join(dir, "requirements.txt")) ||
		fileExists(filepath.Join(dir, "pyproject.toml")) ||
		fileExists(filepath.Join(dir, "Pipfile")) {
		techs["Python"] = true
		tags["python"] = true
		tags["backend"] = true
	}

	// 6. Inspeccionar Supabase directory
	if dirExists(filepath.Join(dir, "supabase")) {
		techs["Supabase"] = true
		tags["supabase"] = true
		tags["postgres"] = true
		tags["database"] = true
	}

	// 7. Inspeccionar Docker
	if fileExists(filepath.Join(dir, "Dockerfile")) ||
		fileExists(filepath.Join(dir, "docker-compose.yml")) ||
		fileExists(filepath.Join(dir, "docker-compose.dev.yml")) ||
		fileExists(filepath.Join(dir, "compose.yaml")) {
		techs["Docker"] = true
		tags["docker"] = true
	}
}

func hasDep(deps map[string]string, key string) bool {
	_, ok := deps[key]
	return ok
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}
