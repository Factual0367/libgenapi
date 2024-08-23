# LibgenAPI-Go

Search Library Genesis programmatically using an this Go package. It supports using Library Genesis' default, author, and title searches.

## Install the package

    go get https://github.com/onurhanak/LibgenAPI-Go

## Example Usage

### Default Search
  
    query := libgenapi.NewQuery("default", "libraries")
    err := query.Search()
    fmt.Println(query.Results)

### Title Search

    query := libgenapi.NewQuery("title", "archaeology")
    err := query.Search()
    fmt.Println(query.Results)

### Author Search
  
    query := libgenapi.NewQuery("author", "Foucault")
    err := query.Search()
    fmt.Println(query.Results)

