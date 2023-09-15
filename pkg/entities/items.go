package entities

type Item struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Picture     string `json:"picture"`
	Description string `json:"description"`
	CategoryID  int    `json:"category_id"`
	Price       int    `json:"price"`
	MinNumber   int    `json:"min_number"`
	MaxNumber   int    `json:"max_number"`
	Command     string `json:"command"`
}
