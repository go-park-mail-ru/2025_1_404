package domain

type Answer struct {
	QuestionID int64
	Rating     int
}

type Question struct {
	ID     int64
	CsatID int64
	Text   string
}

type AnswerDTO struct {
	QuestionID int64 `json:"question_id"`
	Rating     int   `json:"rating"`
}

type QuestionDTO struct {
	ID     int64  `json:"id"`
	CsatID int64  `json:"csat_id"`
	Text   string `json:"text"`
}

type AnswersStat struct {
	TotalAnswers  int           `json:"total_answers"`
	AvgRating     float64       `json:"avg_rating"`
	OneStarStat   oneStarStat   `json:"one_star_stat"`
	TwoStarStat   twoStarStat   `json:"two_star_stat"`
	ThreeStarStat threeStarStat `json:"three_star_stat"`
	FourStarStat  fourStarStat  `json:"four_star_stat"`
	FiveStarStat  fiveStarStat  `json:"five_star_stat"`
}

type EventList struct {
	Events []string `json:"events"`
}

type oneStarStat struct {
	Amount     int
	Percentage float64
}

type twoStarStat struct {
	Amount     int
	Percentage float64
}

type threeStarStat struct {
	Amount     int
	Percentage float64
}

type fourStarStat struct {
	Amount     int
	Percentage float64
}

type fiveStarStat struct {
	Amount     int
	Percentage float64
}
