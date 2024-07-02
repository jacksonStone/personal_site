package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gomarkdown/markdown"
)

// Struct to hold dynamic data
type PageData struct {
	Content template.HTML
	Title   string
	Style   template.CSS
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run build_markdown.go <input_directory> <output_directory>")
		os.Exit(1)
	}

	inputDir := os.Args[1]
	outputDir := os.Args[2]
	tmpl, err := template.ParseFiles("./web_template.html")
	if err != nil {
		fmt.Printf("Failed to compile Template: %v\n", err)
		os.Exit(1)
	}
	style, err := os.ReadFile("./style.css")
	if err != nil {
		fmt.Printf("Failed to load style sheet: %v\n", err)
		os.Exit(1)
	}

	err = filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".md") {
			relPath, err := filepath.Rel(inputDir, path)
			if err != nil {
				return err
			}

			outputPath := filepath.Join(outputDir, strings.TrimSuffix(relPath, ".md")+".html")

			htmlContent, err := convertMarkdownToHTML(path)
			if err != nil {
				fmt.Printf("Error converting %s: %v\n", path, err)
				return err
			}
			pageData := PageData{
				Title: strings.TrimSuffix(info.Name(), ".md"),
				// Covert htmlContent to template.HTML to prevent escaping
				Content: template.HTML(htmlContent),
				Style:   template.CSS(style)}
			// Buffer to hold the rendered template
			var tplBuffer bytes.Buffer
			// Execute the template with the data and write it to the buffer
			err = tmpl.Execute(&tplBuffer, pageData)
			if err != nil {
				log.Fatalf("Error executing template: %v", err)
			}

			// Get the rendered template as a string
			err = os.MkdirAll(filepath.Dir(outputPath), 0755)
			if err != nil {
				return err
			}

			err = os.WriteFile(outputPath, tplBuffer.Bytes(), 0644)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error walking directory: %v\n", err)
		os.Exit(1)
	}
}

func convertMarkdownToHTML(inputPath string) ([]byte, error) {
	mdContent, err := os.ReadFile(inputPath)
	if err != nil {
		return nil, err
	}

	htmlContent := markdown.ToHTML(mdContent, nil, nil)
	return htmlContent, nil
}
