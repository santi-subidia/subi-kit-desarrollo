package updater

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	GitHubOwner = "santi-subidia"
	GitHubRepo  = "subi-kit-desarrollo"
	UserAgent   = "SubiKit-Updater/1.0"
)

// ReleaseInfo representa la información de un release en GitHub.
type ReleaseInfo struct {
	TagName     string         `json:"tag_name"`
	Name        string         `json:"name"`
	Body        string         `json:"body"`
	HTMLURL     string         `json:"html_url"`
	PublishedAt time.Time      `json:"published_at"`
	Assets      []ReleaseAsset `json:"assets"`
}

// ReleaseAsset representa un archivo binario adjunto al release.
type ReleaseAsset struct {
	Name               string `json:"name"`
	Size               int64  `json:"size"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

// UpdateResult contiene los detalles del resultado de una comprobación de versión.
type UpdateResult struct {
	CurrentVersion string
	LatestVersion  string
	UpdateAvail    bool
	Release        *ReleaseInfo
}

// CheckLatest consulta la versión más reciente en GitHub usando la API REST con fallback a redirección web (sin rate-limits).
func CheckLatest(currentVersion string) (*UpdateResult, error) {
	// Intentar primero con la API REST de GitHub
	res, err := checkViaAPI(currentVersion)
	if err == nil {
		return res, nil
	}

	// Si la API falla por 403 (rate limit de IP) o red, intentar vía Web Redirect (inmune a rate limits)
	webRes, webErr := checkViaWebRedirect(currentVersion)
	if webErr == nil {
		return webRes, nil
	}

	return nil, fmt.Errorf("no se pudo verificar actualizaciones: %v (fallback web: %v)", err, webErr)
}

func checkViaAPI(currentVersion string) (*UpdateResult, error) {
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", GitHubOwner, GitHubRepo)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	// Si existe token en el entorno (ej. GITHUB_TOKEN o GH_TOKEN), utilizarlo para elevar cuotas
	if token := getGitHubToken(); token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return &UpdateResult{
			CurrentVersion: currentVersion,
			LatestVersion:  currentVersion,
			UpdateAvail:    false,
			Release:        nil,
		}, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("código HTTP %d (%s)", resp.StatusCode, resp.Status)
	}

	var release ReleaseInfo
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, err
	}

	latestTag := release.TagName
	isNewer := CompareSemver(latestTag, currentVersion) > 0

	return &UpdateResult{
		CurrentVersion: currentVersion,
		LatestVersion:  latestTag,
		UpdateAvail:    isNewer,
		Release:        &release,
	}, nil
}

// checkViaWebRedirect consulta la URL web /releases/latest que redirige a /releases/tag/vX.Y.Z
// Este método no tiene límites de rate limit de API y funciona de manera pública e inmediata.
func checkViaWebRedirect(currentVersion string) (*UpdateResult, error) {
	targetURL := fmt.Sprintf("https://github.com/%s/%s/releases/latest", GitHubOwner, GitHubRepo)

	client := &http.Client{
		Timeout: 10 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Detener la redirección para capturar la cabecera Location
			return http.ErrUseLastResponse
		},
	}

	req, err := http.NewRequest("HEAD", targetURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", UserAgent)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	location := resp.Header.Get("Location")
	if location == "" || strings.HasSuffix(strings.TrimRight(location, "/"), "/releases") || !strings.Contains(location, "/tag/") {
		// Si es 200 o redirige a /releases, no hay releases publicados aún
		return &UpdateResult{
			CurrentVersion: currentVersion,
			LatestVersion:  currentVersion,
			UpdateAvail:    false,
			Release:        nil,
		}, nil
	}

	// Extraer el tag de la URL (ej. https://github.com/.../releases/tag/v0.5.0)
	parts := strings.Split(location, "/tag/")
	if len(parts) < 2 {
		return &UpdateResult{
			CurrentVersion: currentVersion,
			LatestVersion:  currentVersion,
			UpdateAvail:    false,
			Release:        nil,
		}, nil
	}
	latestTag := strings.TrimSpace(parts[1])

	// Construir información de release sintética con los assets esperados
	ext := "tar.gz"
	if runtime.GOOS == "windows" {
		ext = "zip"
	}
	assetName := fmt.Sprintf("subikit-%s-%s.%s", runtime.GOOS, runtime.GOARCH, ext)
	downloadURL := fmt.Sprintf("https://github.com/%s/%s/releases/download/%s/%s", GitHubOwner, GitHubRepo, latestTag, assetName)

	release := &ReleaseInfo{
		TagName: latestTag,
		Name:    latestTag,
		HTMLURL: location,
		Assets: []ReleaseAsset{
			{
				Name:               assetName,
				BrowserDownloadURL: downloadURL,
			},
		},
	}

	isNewer := CompareSemver(latestTag, currentVersion) > 0

	return &UpdateResult{
		CurrentVersion: currentVersion,
		LatestVersion:  latestTag,
		UpdateAvail:    isNewer,
		Release:        release,
	}, nil
}

func getGitHubToken() string {
	if t := os.Getenv("GITHUB_TOKEN"); t != "" {
		return t
	}
	if t := os.Getenv("GH_TOKEN"); t != "" {
		return t
	}
	return ""
}

// CompareSemver compara dos cadenas semver (ej. "v0.5.0" vs "0.4.0").
// Retorna 1 si v1 > v2, -1 si v1 < v2, y 0 si v1 == v2.
func CompareSemver(v1, v2 string) int {
	clean1 := strings.TrimPrefix(strings.TrimSpace(v1), "v")
	clean2 := strings.TrimPrefix(strings.TrimSpace(v2), "v")

	base1 := strings.Split(clean1, "-")[0]
	base2 := strings.Split(clean2, "-")[0]

	parts1 := strings.Split(base1, ".")
	parts2 := strings.Split(base2, ".")

	maxLen := len(parts1)
	if len(parts2) > maxLen {
		maxLen = len(parts2)
	}

	for i := 0; i < maxLen; i++ {
		num1 := 0
		num2 := 0

		if i < len(parts1) {
			num1, _ = strconv.Atoi(parts1[i])
		}
		if i < len(parts2) {
			num2, _ = strconv.Atoi(parts2[i])
		}

		if num1 > num2 {
			return 1
		}
		if num1 < num2 {
			return -1
		}
	}

	return 0
}

// FindAssetForPlatform busca el asset correspondiente para el OS y arquitectura actuales.
func FindAssetForPlatform(assets []ReleaseAsset, goos, goarch string) (*ReleaseAsset, error) {
	expectedExt := ".tar.gz"
	if goos == "windows" {
		expectedExt = ".zip"
	}

	targetPattern := fmt.Sprintf("%s-%s", goos, goarch)

	for _, asset := range assets {
		name := strings.ToLower(asset.Name)
		if strings.Contains(name, targetPattern) && strings.HasSuffix(name, expectedExt) {
			return &asset, nil
		}
	}

	for _, asset := range assets {
		name := strings.ToLower(asset.Name)
		if strings.Contains(name, goos) && strings.Contains(name, goarch) {
			return &asset, nil
		}
	}

	// Si solo hay un asset o coincidencia sintética, usarlo
	if len(assets) == 1 {
		return &assets[0], nil
	}

	return nil, fmt.Errorf("no se encontró binario compatible para %s/%s en este release", goos, goarch)
}

// ApplyUpdate descarga el paquete del release, extrae el binario y reemplaza el ejecutable actual de forma segura.
func ApplyUpdate(release *ReleaseInfo, progressCb func(percent float64, status string)) error {
	if release == nil {
		return fmt.Errorf("información de release no disponible")
	}

	asset, err := FindAssetForPlatform(release.Assets, runtime.GOOS, runtime.GOARCH)
	if err != nil {
		return err
	}

	if progressCb != nil {
		progressCb(0.1, fmt.Sprintf("Descargando %s...", asset.Name))
	}

	// 1. Descargar paquete
	client := &http.Client{Timeout: 60 * time.Second}
	req, err := http.NewRequest("GET", asset.BrowserDownloadURL, nil)
	if err != nil {
		return fmt.Errorf("error creando petición de descarga: %w", err)
	}
	req.Header.Set("User-Agent", UserAgent)
	if token := getGitHubToken(); token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error al descargar paquete: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error al descargar binario (%s)", resp.Status)
	}

	archiveBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error al leer contenido del paquete: %w", err)
	}

	if progressCb != nil {
		progressCb(0.5, "Extrayendo ejecutable...")
	}

	// 2. Extraer el binario del archivo comprimido
	targetBinaryName := "subikit"
	if runtime.GOOS == "windows" {
		targetBinaryName = "subikit.exe"
	}

	binaryBytes, err := extractBinary(archiveBytes, asset.Name, targetBinaryName)
	if err != nil {
		return fmt.Errorf("error al extraer binario: %w", err)
	}

	if progressCb != nil {
		progressCb(0.8, "Reemplazando ejecutable actual...")
	}

	// 3. Obtener ruta del ejecutable en curso
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("no se pudo determinar la ruta del ejecutable actual: %w", err)
	}
	execPath, err = filepath.EvalSymlinks(execPath)
	if err != nil {
		return fmt.Errorf("error al resolver symlink del ejecutable: %w", err)
	}

	// 4. Reemplazo seguro según el Sistema Operativo
	if runtime.GOOS == "windows" {
		err = replaceWindowsExecutable(execPath, binaryBytes)
	} else {
		err = replaceUnixExecutable(execPath, binaryBytes)
	}

	if err != nil {
		return fmt.Errorf("error al reemplazar binario: %w", err)
	}

	if progressCb != nil {
		progressCb(1.0, fmt.Sprintf("¡Actualizado exitosamente a %s!", release.TagName))
	}

	return nil
}

// extractBinary extrae el binario desde un archivo .zip o .tar.gz en memoria.
func extractBinary(archiveData []byte, assetName, targetName string) ([]byte, error) {
	if strings.HasSuffix(strings.ToLower(assetName), ".zip") {
		reader, err := zip.NewReader(bytes.NewReader(archiveData), int64(len(archiveData)))
		if err != nil {
			return nil, fmt.Errorf("error abriendo zip: %w", err)
		}

		for _, f := range reader.File {
			baseName := filepath.Base(f.Name)
			if strings.EqualFold(baseName, targetName) {
				rc, err := f.Open()
				if err != nil {
					return nil, err
				}
				defer rc.Close()
				return io.ReadAll(rc)
			}
		}
		return nil, fmt.Errorf("el archivo %s no contiene '%s'", assetName, targetName)
	}

	if strings.HasSuffix(strings.ToLower(assetName), ".tar.gz") || strings.HasSuffix(strings.ToLower(assetName), ".tgz") {
		gzReader, err := gzip.NewReader(bytes.NewReader(archiveData))
		if err != nil {
			return nil, fmt.Errorf("error descomprimiendo gzip: %w", err)
		}
		defer gzReader.Close()

		tarReader := tar.NewReader(gzReader)
		for {
			header, err := tarReader.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				return nil, fmt.Errorf("error leyendo tar: %w", err)
			}

			baseName := filepath.Base(header.Name)
			if strings.EqualFold(baseName, targetName) {
				return io.ReadAll(tarReader)
			}
		}
		return nil, fmt.Errorf("el archivo %s no contiene '%s'", assetName, targetName)
	}

	return archiveData, nil
}

// replaceWindowsExecutable realiza un swap seguro renombrando el ejecutable en curso a .old
func replaceWindowsExecutable(execPath string, newBytes []byte) error {
	oldPath := execPath + ".old"

	_ = os.Remove(oldPath)

	if err := os.Rename(execPath, oldPath); err != nil {
		return fmt.Errorf("no se pudo renombrar ejecutable actual a .old: %w", err)
	}

	if err := os.WriteFile(execPath, newBytes, 0755); err != nil {
		_ = os.Rename(oldPath, execPath)
		return fmt.Errorf("error al escribir nuevo binario: %w", err)
	}

	return nil
}

// replaceUnixExecutable reemplaza atómicamente el ejecutable en Linux/macOS
func replaceUnixExecutable(execPath string, newBytes []byte) error {
	dir := filepath.Dir(execPath)
	tmpFile, err := os.CreateTemp(dir, ".subikit-update-*")
	if err != nil {
		return fmt.Errorf("error al crear archivo temporal: %w", err)
	}
	tmpPath := tmpFile.Name()

	defer func() {
		_ = os.Remove(tmpPath)
	}()

	if _, err := tmpFile.Write(newBytes); err != nil {
		tmpFile.Close()
		return fmt.Errorf("error al escribir temporal: %w", err)
	}
	if err := tmpFile.Chmod(0755); err != nil {
		tmpFile.Close()
		return fmt.Errorf("error configurando permisos 0755: %w", err)
	}
	tmpFile.Close()

	if err := os.Rename(tmpPath, execPath); err != nil {
		return fmt.Errorf("error al reemplazar binario con atomic rename: %w", err)
	}

	return nil
}

// CleanupOld elimina archivos residuales .old generados durante actualizaciones previas en Windows.
func CleanupOld() {
	if runtime.GOOS != "windows" {
		return
	}
	execPath, err := os.Executable()
	if err != nil {
		return
	}
	oldPath := execPath + ".old"
	if _, err := os.Stat(oldPath); err == nil {
		_ = os.Remove(oldPath)
	}
}
