package libgenapi

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	colly "github.com/gocolly/colly"
)

type Book struct {
	ID           string
	MD5          string
	Title        string
	Author       string
	Publisher    string
	Year         string
	Language     string
	Pages        string
	Size         string
	Extension    string
	DownloadLink string
}

type Query struct {
	QueryType string
	Query     string
	SearchURL string
	Results   []Book
}

func NewQuery(queryType, query string) *Query {
	return &Query{
		QueryType: queryType,
		Query:     query,
	}
}

func (q *Query) Search() error {
	q.SearchURL = searchURLHandler(q.Query, q.QueryType)
	results, err := scrapeURL(q.SearchURL)
	if err != nil {
		return err
	}
	q.Results = results
	return nil
}

func searchURLHandler(query, queryType string) string {
	return fmt.Sprintf("https://libgen.is/search.php?req=%s&column=%s&res=100", query, queryType)
}

func generateDownloadLink(md5, bookID, bookTitle, bookFiletype string) string {
	var newBookID string
	if len(bookID) == 4 {
		newBookID = string(bookID[:1]) + "000"
	} else if len(bookID) == 5 {
		newBookID = string(bookID[:2]) + "000"
	}

	md5 = strings.ToLower(md5)
	bookTitle = strings.Replace(bookTitle, " ", "_", -1)
	return fmt.Sprintf("https://download.library.lol/main/%s/%s/%s.%s", newBookID, md5, bookTitle, bookFiletype)
}

func scrapeURL(searchURL string) ([]Book, error) {
	c := colly.NewCollector(
		colly.AllowedDomains("libgen.is"),
	)

	var books []Book
	skipFirstBook := true

	c.OnHTML("tr", func(e *colly.HTMLElement) {
		book := Book{}

		// to skip the first row
		isValidBook := false

		e.ForEach("td", func(index int, el *colly.HTMLElement) {
			switch index {
			case 0:
				if _, err := strconv.Atoi(el.Text); err == nil {
					book.ID = el.Text
					isValidBook = true
				}
			case 1:
				book.Author = el.Text
			case 2:
				book.Title = el.ChildText("a")
				md5 := strings.Split(el.ChildAttr("a", "href"), "md5=")
				if len(md5) == 2 {
					book.MD5 = md5[1]
				} else {
					book.MD5 = ""
				}
			case 3:
				book.Publisher = el.Text
			case 4:
				book.Year = el.Text
			case 5:
				book.Pages = el.Text
			case 6:
				book.Language = el.Text
			case 7:
				book.Size = el.Text
			case 8:
				book.Extension = el.Text
			}
		})

		if isValidBook && skipFirstBook {
			skipFirstBook = false
			return
		}

		if isValidBook && book.MD5 == "" {
			return
		}

		if isValidBook && book.Title != "" {
			book.DownloadLink = generateDownloadLink(book.MD5, book.ID, book.Title, book.Extension)
			books = append(books, book)
		}
	})

	err := c.Visit(searchURL)
	if err != nil {
		log.Println("Failed to visit target page:", err)
		return nil, err
	}

	return books, nil
}
