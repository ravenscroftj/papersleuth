package sleuth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type DOIResolveClient struct {
}

func tidyDOI(doi string) string {

	doi = strings.TrimPrefix(doi, "doi:")

	return doi
}

type PaperResponse struct {
	Abstract    string
	Title       string
	DOI         string
	FullSources []string
	Authors     []string
}

func (client *DOIResolveClient) ResolveDOI(doi string) (*PaperResponse, error) {

	resp, err := http.Get(fmt.Sprintf("https://dx.doi.org/%s", tidyDOI(doi)))

	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		return nil, err
	}

	paper := &PaperResponse{}

	doc.Find("meta").Each(func(i int, el *goquery.Selection) {

		name := el.AttrOr("name", "")

		switch strings.ToLower(name) {
		case "citation_title":
			paper.Title = el.AttrOr("content", "")
			break

		case "citation_doi":
		case "dc.identifier":
			paper.DOI = el.AttrOr("content", "")
			break

		case "citation_abstract":
		case "description":
			paper.Abstract = el.AttrOr("content", "")
			break

		case "citation_pdf_url":
			paper.FullSources = append(paper.FullSources, el.AttrOr("content", ""))
			break

		case "citation_author":
			paper.Authors = append(paper.Authors, el.AttrOr("content", ""))
			break
		}

	})

	return paper, nil
}
