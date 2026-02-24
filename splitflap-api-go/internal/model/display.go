package model

type Display struct {
	ID      string         `json:"id"`
	Content DisplayContent `json:"content"`
	Config  DisplayConfig  `json:"config"`
}

type DisplayContent struct {
	Rows [][]string `json:"rows"`
}

type DisplayConfig struct {
	RowCount    int `json:"rowCount"`
	ColumnCount int `json:"columnCount"`
}
