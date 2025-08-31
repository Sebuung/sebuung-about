package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/adrg/frontmatter"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

type Author struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

type PostData struct {
	Title       string
	Author      Author
	Date        string
	ReadingTime string
	Content     template.HTML
}

type IndexPostData struct {
	Title         string
	Author        Author
	Date          string
	FormattedDate string
	Excerpt       string
	FileName      string
}

type FrontMatter struct {
	Title       string `yaml:"title"`
	Author      Author `yaml:"author"`
	Date        string `yaml:"date"`
	ReadingTime string `yaml:"readingTime"`
	Excerpt     string `yaml:"excerpt"`
}

func main() {
	postsDir := "posts"
	publicDir := "public"
	distDir := "dist"
	blogDir := filepath.Join(distDir, "blog")
	postsOutputDir := filepath.Join(blogDir, "posts")

	if err := os.MkdirAll(postsOutputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directories: %v", err)
	}

	if err := copyDir(publicDir, filepath.Join(distDir, "public")); err != nil {
		log.Fatalf("Failed to copy public directory: %v", err)
	}

	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
		goldmark.WithRendererOptions(html.WithHardWraps(), html.WithUnsafe()),
	)

	postTemplate, err := template.ParseFiles("templates/post.html")
	if err != nil {
		log.Fatalf("Failed to parse post template: %v", err)
	}

	blogIndexTemplate, err := template.ParseFiles("templates/blog.html")
	if err != nil {
		log.Fatalf("Failed to parse blog index template: %v", err)
	}

	generateStaticPage("templates/index.html", filepath.Join(distDir, "index.html"))

	var posts []IndexPostData

	err = filepath.Walk(postsDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".md") {
			post, err := processFile(path, postsOutputDir, md, postTemplate)
			if err != nil {
				log.Printf("Error processing %s: %v", path, err)
			} else {
				posts = append(posts, post)
			}
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Failed to walk posts directory: %v", err)
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Date > posts[j].Date
	})

	generateBlogIndexPage(posts, blogDir, blogIndexTemplate)

	fmt.Println("Blog generation complete!")
}

func processFile(path, outputDir string, md goldmark.Markdown, tpl *template.Template) (IndexPostData, error) {
	source, err := os.ReadFile(path)
	if err != nil {
		return IndexPostData{}, fmt.Errorf("failed to read file %s: %w", path, err)
	}

	var fm FrontMatter
	content, err := frontmatter.Parse(bytes.NewReader(source), &fm)
	if err != nil {
		return IndexPostData{}, fmt.Errorf("failed to parse front matter for %s: %w", path, err)
	}

	var buf bytes.Buffer
	if err := md.Convert(content, &buf); err != nil {
		return IndexPostData{}, fmt.Errorf("failed to convert markdown for %s: %w", path, err)
	}

	parsedDate, err := time.Parse("2006-01-02", fm.Date)
	formattedDate := fm.Date
	if err != nil {
		log.Printf("Warning: could not parse date for %s. Using original string.", path)
	} else {
		formattedDate = parsedDate.Format("2006년 1월 2일")
	}

	postData := PostData{
		Title:       fm.Title,
		Author:      fm.Author,
		Date:        formattedDate,
		ReadingTime: fm.ReadingTime,
		Content:     template.HTML(buf.String()),
	}

	outputFileName := strings.TrimSuffix(filepath.Base(path), ".md") + ".html"
	outputFilePath := filepath.Join(outputDir, outputFileName)
	file, err := os.Create(outputFilePath)
	if err != nil {
		return IndexPostData{}, fmt.Errorf("failed to create output file %s: %w", outputFilePath, err)
	}
	defer file.Close()

	if err := tpl.Execute(file, postData); err != nil {
		return IndexPostData{}, fmt.Errorf("failed to execute template for %s: %w", outputFilePath, err)
	}

	fmt.Printf("Generated post: %s\n", outputFilePath)

	indexData := IndexPostData{
		Title:         fm.Title,
		Author:        fm.Author,
		Date:          fm.Date,
		FormattedDate: formattedDate,
		Excerpt:       fm.Excerpt,
		FileName:      outputFileName,
	}

	return indexData, nil
}

func generateBlogIndexPage(posts []IndexPostData, outputDir string, tpl *template.Template) {
	indexPath := filepath.Join(outputDir, "index.html")
	file, err := os.Create(indexPath)
	if err != nil {
		log.Fatalf("Failed to create blog index page: %v", err)
	}
	defer file.Close()

	if err := tpl.Execute(file, posts); err != nil {
		log.Fatalf("Failed to execute blog index template: %v", err)
	}

	fmt.Printf("Generated blog index: %s\n", indexPath)
}

func generateStaticPage(templatePath, outputPath string) {
	tpl, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Fatalf("Failed to parse static template %s: %v", templatePath, err)
	}

	file, err := os.Create(outputPath)
	if err != nil {
		log.Fatalf("Failed to create static page %s: %v", outputPath, err)
	}
	defer file.Close()

	if err := tpl.Execute(file, nil); err != nil {
		log.Fatalf("Failed to execute static template %s: %v", outputPath, err)
	}
	fmt.Printf("Generated static page: %s\n", outputPath)
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		dstPath := filepath.Join(dst, strings.TrimPrefix(path, src))

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		dstFile, err := os.Create(dstPath)
		if err != nil {
			return err
		}
		defer dstFile.Close()

		_, err = io.Copy(dstFile, srcFile)
		return err
	})
}
