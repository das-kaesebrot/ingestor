//go:build ignore

package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
)

type webDeps struct {
	Dependencies map[string]string `json:"dependencies"`
}

type downloadObject struct {
	rootPackage string
	version     string
	file        string
}

func (o *downloadObject) buildURL() string {
	return fmt.Sprintf(unpkgUrl, o.rootPackage, o.version, o.file)
}

var (
	unpkgUrl = "https://unpkg.com/%s@%s/%s"
	outDir   = "web/static/libs"
)

func getPackageFiles(pkg string, version string) (map[string]downloadObject, error) {
	switch pkg {
	case "bootstrap":
		return map[string]downloadObject{
			path.Join(outDir, "js/bootstrap.bundle.min.js"): {rootPackage: pkg, version: version, file: "dist/js/bootstrap.bundle.min.js"},
			path.Join(outDir, "css/bootstrap.min.css"):      {rootPackage: pkg, version: version, file: "dist/css/bootstrap.min.css"},
		}, nil
	case "htmx.org":
		return map[string]downloadObject{
			path.Join(outDir, "js/htmx.min.js"): {rootPackage: pkg, version: version, file: "dist/htmx.min.js"},
		}, nil
	default:
		return nil, fmt.Errorf("unknown pkg: '%s'", pkg)
	}
}

func main() {
	// we still read package.json so renovate thinks it's managing npm packages and we can have automatic updates
	packageJson, err := os.ReadFile("package.json")
	if err != nil {
		log.Fatalf("%v", err)
	}

	var d webDeps
	err = json.Unmarshal(packageJson, &d)
	if err != nil {
		log.Fatalf("%v", err)
	}

	packageFiles := make(map[string]downloadObject)

	for name, version := range d.Dependencies {
		f, err := getPackageFiles(name, version)

		for k, v := range f {
			packageFiles[k] = v
		}

		if err != nil {
			log.Fatalf("%v", err)
		}
	}

	err = downloadToWebDir(packageFiles)
	if err != nil {
		log.Fatalf("%v", err)
	}
}

func downloadToWebDir(packageFiles map[string]downloadObject) error {
	for outputFile, dObject := range packageFiles {
		resp, err := http.Get(dObject.buildURL())
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		downloadedHash, downloadedContent, err := hashReader(resp.Body)
		if err != nil {
			return err
		}

		if fileExists(outputFile) {
			existingHash, err := hashFile(outputFile)
			if err != nil {
				return err
			}

			if downloadedHash == existingHash {
				log.Printf("UNCHANGED '%s' (old=%s) (skipping download)", outputFile, existingHash)
				continue
			}

			log.Printf("UPDATED '%s' (old=%s, new=%s)", outputFile, existingHash, downloadedHash)
		} else {
			log.Printf("NEW '%s' (new=%s)", outputFile, downloadedHash)
		}

		err = os.MkdirAll(path.Dir(outputFile), 0750)
		if err != nil {
			return err
		}

		out, err := os.Create(outputFile)
		if err != nil {
			return err
		}
		defer out.Close()

		_, err = io.Copy(out, downloadedContent)
		if err != nil {
			return err
		}
	}

	return nil
}

func hashReader(r io.Reader) (string, io.Reader, error) {
	h := sha256.New()
	buf := new(bytes.Buffer)

	mw := io.MultiWriter(h, buf)
	_, err := io.Copy(mw, r)
	if err != nil {
		return "", nil, err
	}

	hash := fmt.Sprintf("%x", h.Sum(nil))
	return hash, bytes.NewReader(buf.Bytes()), nil
}

func hashFile(filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
