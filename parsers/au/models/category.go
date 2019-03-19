package aumodels

// Category is the root element of a category output
type Category struct {
    Adverts                 []Advert            `json:"ads"`
    Pagination              *CategoryPagination `json:"pagination"`
}

// CategoryPagination is the page control of the category
type CategoryPagination struct {
    CurrentPage             uint                `json:"current"`
    PageSize                uint                `json:"page_size"`
    EntrySize               uint                `json:"entry_size"`
}
