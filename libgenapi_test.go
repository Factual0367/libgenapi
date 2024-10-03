package libgenapi

import (
	"fmt"
	"testing"
)

func TestSearchURLHandler(t *testing.T) {
	query := "Karl Marx"
	queryType := "author"
	querySize := 25
	expectedURL := "https://libgen.is/search.php?req=Karl%20Marx&column=author&res=25"
	result := searchURLHandler(query, queryType, querySize)
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

func TestQueryEmptySearch(t *testing.T) {
	query := NewQuery("default", "", 25)
	err := query.Search()
	if err != nil {
		t.Fatalf("Query.Search() error: %v", err)
	}
}

func TestQuerySearch(t *testing.T) {
	query := NewQuery("author", "Marx", 25)
	err := query.Search()
	if err != nil {
		t.Fatalf("Query.Search() error: %v", err)
	}

	if len(query.Results) == 0 {
		t.Errorf("Query.Search() returned 0 results; want > 0")
	} else {
		book := query.Results[10]

		err = book.AddSecondDownloadLink()
		if err != nil {
			t.Errorf("AddSecondDownloadLink() error: %v", err)
		}

		if book.AlternativeDownloadLink == "" {
			t.Errorf("AlternativeDownloadLink is empty; expected a valid link")
		} else {
			fmt.Printf("AlternativeDownloadLink: %s\n", book.AlternativeDownloadLink)
		}

	}
	t.Logf("Full book details: %+v\n", query.Results[10])
}
