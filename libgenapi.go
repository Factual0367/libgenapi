package libgenapi

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
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
	CoverLink    string
}

type Query struct {
	QueryType string
	Query     string
	SearchURL string
	QuerySize int
	Results   []Book
}

func NewQuery(queryType, query string, querySize int) *Query {
	return &Query{
		QueryType: queryType,
		Query:     query,
		QuerySize: querySize,
	}
}

func (q *Query) Search() error {
	q.SearchURL = searchURLHandler(q.Query, q.QueryType, q.QuerySize)
	results, err := scrapeURL(q.SearchURL)
	if err != nil {
		return err
	}
	q.Results = results
	return nil
}

func searchURLHandler(query, queryType string, querySize int) string {
	query = strings.ReplaceAll(query, " ", "%20")
	return fmt.Sprintf("https://libgen.is/search.php?req=%s&column=%s&res=%d", query, queryType, querySize)
}

func generateDownloadLink(md5, bookID, bookTitle, bookFiletype string) string {
	var newBookID string
	if len(bookID) == 4 {
		newBookID = string(bookID[:1]) + "000"
	} else if len(bookID) == 5 {
		newBookID = string(bookID[:2]) + "000"
	} else if len(bookID) == 6 {
		newBookID = string(bookID[:3]) + "000"
	} else if len(bookID) == 7 {
		newBookID = string(bookID[:4]) + "000"
	}

	md5 = strings.ToLower(md5)
	bookTitle = strings.Replace(bookTitle, " ", "_", -1)
	return fmt.Sprintf("https://download.library.lol/main/%s/%s/%s.%s", newBookID, md5, bookTitle, bookFiletype)
}

func getOpenLibraryId(idsJoined string) []map[string]string {
	url := fmt.Sprintf("https://libgen.is/json.php?ids=%s&fields=id,openlibraryid", idsJoined)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var jsonMap []map[string]string

	err = json.Unmarshal(body, &jsonMap)
	if err != nil {
		log.Fatal(err)
	}

	return jsonMap
}

func addBookCoverLinks(books []Book) []Book {

	ids := make([]string, len(books))
	for i, book := range books {
		ids[i] = book.ID
	}
	idsJoined := strings.Join(ids, ",")
	openLibraryIds := getOpenLibraryId(idsJoined)

	for i, book := range books {
		for _, id := range openLibraryIds {
			if book.ID == id["id"] {
				if id["openlibraryid"] != "" {
					books[i].CoverLink = fmt.Sprintf("https://covers.openlibrary.org/b/olid/%s-M.jpg", id["openlibraryid"])
				} else {
					books[i].CoverLink = ""
				}
			}
		}
	}

	return books
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
	books = addBookCoverLinks(books)
	return books, nil
}
