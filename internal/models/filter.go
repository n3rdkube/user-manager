package models

// ListOptions is the data model used to filter data by list shared by processor and manager
type ListOptions struct {

	// include elements with these values
	Include User `json:"include,omitempty"`

	PageNumber  int `json:"page_number,omitempty"`
	RowsPerPage int `json:"rows_per_page,omitempty"`
}
