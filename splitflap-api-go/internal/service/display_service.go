package service

import "splitflap-api-go/internal/model"

type DisplayService struct {
	demoDisplay model.Display
}

func NewDisplayService() *DisplayService {
	return &DisplayService{
		demoDisplay: model.Display{
			ID: "demo",
			Content: model.DisplayContent{
				Rows: [][]string{
					{"H", "E", "L", "L", "O", " ", "W", "O", "R", "L"},
					{"D", " ", "S", "P", "L", "I", "T", "F", "L", "A"},
					{"P", " ", "D", "I", "S", "P", "L", "A", "Y", " "},
					{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"},
					{"-", ":", ".", ",", " ", "A", "Z", "a", "z", "!"},
				},
			},
			Config: model.DisplayConfig{
				RowCount:    5,
				ColumnCount: 10,
			},
		},
	}
}

func (s *DisplayService) GetDisplay(id string) *model.Display {
	if id != s.demoDisplay.ID {
		return nil
	}

	return &s.demoDisplay
}
