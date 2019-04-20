package aumodels

// Categories are information about categories and subcategories
type Categories struct {
    ID                      uint                `json:"id"`
    Name                    string              `json:"name"`
    Slug                    string              `json:"slug"`
    ParentID                *uint               `json:"parent_id"`
    ParentSlug              *string             `json:"parent_slug"`
    ChildrenCount           uint                `json:"children_count"`
    Subcategories           []Categories        `json:"subcategories"`
    IsRootCategory          bool                `json:"is_root"`
}
