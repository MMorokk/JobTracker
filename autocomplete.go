package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/ollama/ollama/api"
	"golang.org/x/net/html"
)

// skipTags lists HTML elements whose content is irrelevant to job postings.
// Skipping them reduces noise and token count sent to the LLM.
var skipTags = map[string]bool{
	"script": true, "style": true, "svg": true, "img": true,
	"nav": true, "footer": true, "header": true, "noscript": true,
	"iframe": true, "meta": true, "link": true,
}

// JobPosting holds structured data extracted from a job listing page.
type JobPosting struct {
	Title        string
	Company      string
	Location     string
	Type         string
	WorkingMode  string
	Salary       string
	Description  string
	Summary      string
	Requirements []string
	URL          string
}

// AutoFill scrapes the page at url, sends its cleaned text to a local Ollama
// model, and returns a populated JobPosting. model should be an Ollama model
// name (e.g. "llama3:8b", "mistral:7b").
func AutoFill(url string, model string) (autoFilled JobPosting, err error) {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		return JobPosting{}, err
	}
	// 3-minute ceiling covers slow LLM inference on large pages.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	prompt := fmt.Sprintf(`You are a data extraction assistant. Extract structured job posting data from the HTML below.

Return ONLY a valid JSON object — no markdown, no explanation, no code fences. Use null for any field not found.

Schema:
{
  "title": "short job title (strip unneeded words)",
  "company": "company name",
  "location": "city, country (use office location even if role is remote, if city is not known, write only country)",
  "type": "full-time | part-time | internship | freelance | contract",
	"workingMode": "remote | hybrid | in-office",
  "salary": "salary range or null",
  "description": "full job description text",
  "summary": "2-3 sentence summary of the role written by you",
  "requirements": ["requirement 1", "requirement 2", "..."]
}

HTML:
%v`, cleanHTML(scrapeJS(url)))
	//println(prompt)
	req := &api.GenerateRequest{
		Model:  model,
		Prompt: prompt,
		Stream: new(false),
	}

	// client.Generate streams responses; with Stream=false we get one final chunk.
	var fullResponse string
	err = client.Generate(ctx, req, func(resp api.GenerateResponse) error {
		fullResponse = resp.Response
		return nil
	})
	if err != nil {
		return JobPosting{}, err
	}

	var jp JobPosting
	err = json.Unmarshal([]byte(fullResponse), &jp)
	if err != nil {
		return JobPosting{}, err
	}

	jp.URL = url
	return jp, nil
}

// scrapeJS navigates to url in a headless Chromium instance and returns the
// rendered inner HTML of <body>. Using a real browser ensures JavaScript-heavy
// job boards (LinkedIn, Greenhouse, etc.) are fully rendered before extraction.
func scrapeJS(url string) string {
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(),
		append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.Flag("no-sandbox", true),
			chromedp.Flag("headless", true),
			chromedp.Flag("disable-dev-shm-usage", true),
			chromedp.UserAgent("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"),
		)...,
	)
	defer cancel()
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()
	// Hard cap so a hanging page doesn't block the whole program.
	ctx, cancel = context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	var content string
	err := chromedp.Run(ctx,
		chromedp.ActionFunc(func(ctx context.Context) error {
			// Separate shorter timeout for navigation only — don't wait
			// 2 minutes just because the page is slow to respond.
			navCtx, cancel := context.WithTimeout(ctx, 8*time.Second)
			defer cancel()
			return chromedp.Navigate(url).Do(navCtx)
		}),
		chromedp.Sleep(4*time.Second), // wait for JS to finish rendering
		chromedp.InnerHTML("body", &content),
	)
	if err != nil {
		log.Fatal(err)
	}
	return content
}

// cleanHTML parses raw HTML, strips non-content tags, extracts visible text,
// and truncates to ~4000 characters to stay within local model context limits.
func cleanHTML(raw string) string {
	doc, err := html.Parse(strings.NewReader(raw))
	if err != nil {
		return raw
	}

	var buf bytes.Buffer
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		// Skip entire subtrees for tags that carry no useful text.
		if n.Type == html.ElementNode && skipTags[n.Data] {
			return
		}
		if n.Type == html.TextNode {
			text := strings.TrimSpace(n.Data)
			if text != "" {
				buf.WriteString(text)
				buf.WriteByte('\n')
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}

	walk(doc)

	// Collapse multiple blank lines
	lines := strings.Split(buf.String(), "\n")
	var result []string
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			result = append(result, line)
		}
	}

	output := strings.Join(result, "\n")

	// Truncate to ~4000 chars for local models
	if len(output) > 4000 {
		output = output[:4000]
	}

	return output
}
