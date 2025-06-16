package structs

type FileSortConfig struct {
	Type string
	Dir  string
}

type Config struct {
	Dir               string   `json:"sortDir"`
	AcceptedLanguages []string `json:"acceptedLanguages"`
}
