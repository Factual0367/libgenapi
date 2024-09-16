# LibgenAPI

Search Library Genesis programmatically using an this Go package. It supports using Library Genesis' default, author, and title searches.

## Install the package

    go get github.com/onurhanak/libgenapi

## Example Usage

### Default Search
  
    query := libgenapi.NewQuery("default", "libraries", 25) // or 50, 100
    err := query.Search()
    fmt.Println(query.Results)

### Title Search

    query := libgenapi.NewQuery("title", "archaeology", 25) // or 50, 100
    err := query.Search()
    fmt.Println(query.Results)

### Author Search
  
    query := libgenapi.NewQuery("author", "Foucault", 25) // or 50, 100
    err := query.Search()
    fmt.Println(query.Results)

### Add Alternative Download Link to Book

    query := libgenapi.NewQuery("author", "Foucault", 25) // or 50, 100
    err := query.Search()
    book := query.Results[10]
    err := book.AddSecondDownloadLink()
