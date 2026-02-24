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
					{"TIME", "DESTINATION", "PLATFORM", "STATUS"},
					{"10:30", "OTTAWA", "3", "ON TIME"},
					{"10:45", "VANCOUVER", "5", "DELAYED"},
					{"11:00", "MONTREAL", "7", "BOARDING"},
					{"11:15", "HALIFAX", "2", "ON TIME"},
				},
			},
			Config: model.DisplayConfig{
				RowCount:    5,
				ColumnCount: 4,
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
