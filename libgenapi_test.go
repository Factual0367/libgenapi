package libgenapi

import (
	"fmt"
	"testing"
)

func TestSearchURLHandler(t *testing.T) {
	query := "Marx"
	queryType := "author"
	expectedURL := "https://libgen.is/search.php?req=Marx&column=author&res=100"

	result := searchURLHandler(query, queryType)
	if result != expectedURL {
		t.Errorf("searchURLHandler() = %v; want %v", result, expectedURL)
	}
}

func TestGenerateDownloadLink(t *testing.T) {
	md5 := "abcd1234"
	bookID := "1234"
	bookTitle := "Das Kapital"
	bookFiletype := "pdf"
	expectedLink := "https://download.library.lol/main/1000/abcd1234/Das_Kapital.pdf"

	result := generateDownloadLink(md5, bookID, bookTitle, bookFiletype)
	if result != expectedLink {
		t.Errorf("generateDownloadLink() = %v; want %v", result, expectedLink)
	}
}

func TestQuerySearch(t *testing.T) {
	query := NewQuery("default", "Marx")
	err := query.Search()
	if err != nil {
		t.Fatalf("Query.Search() error: %v", err)
	}

	if len(query.Results) == 0 {
		t.Errorf("Query.Search() returned 0 results; want > 0")
	} else {
		book := query.Results[3]
		fmt.Printf("ID: %s\n", book.ID)
		fmt.Printf("MD5: %s\n", book.MD5)
		fmt.Printf("Title: %s\n", book.Title)
		fmt.Printf("Author: %s\n", book.Author)
		fmt.Printf("Publisher: %s\n", book.Publisher)
		fmt.Printf("Year: %s\n", book.Year)
		fmt.Printf("Language: %s\n", book.Language)
		fmt.Printf("Pages: %s\n", book.Pages)
		fmt.Printf("Size: %s\n", book.Size)
		fmt.Printf("Extension: %s\n", book.Extension)
		fmt.Printf("DownloadLink: %s\n", book.DownloadLink)
		fmt.Printf("CoverLink: %s\n", book.CoverLink)

	}
}
