package fitbit2gcal

type Schedule struct {
	Title    string
	Location string
	Year     string
	Month    string
	Day      string
	Start    string
	End      string
}

// Sleep : Sleep data
type Sleep struct {
	Sleep []SleepData `json:"sleep"`
}

type SleepData struct {
	DateOfSleep         string      `json:"dateOfSleep"`
	Duration            int         `json:"duration"`
	Efficiency          int         `json:"efficiency"`
	IsMainSleep         bool        `json:"isMainSleep"`
	Levels              SleepLevels `json:"levels"`
	LogID               int         `json:"logId"`
	MinutesAfterWakeup  int         `json:"minutesAfterWakeup"`
	MinutesAsleep       int         `json:"minutesAsleep"`
	MinutesAwake        int         `json:"minutesAwake"`
	MinutesToFallAsleep int         `json:"minutesToFallAsleep"`
	StartTime           string      `json:"startTime"`
	TimeInBed           int         `json:"timeInBed"`
	Type                string      `json:"type"`
}

type SleepLevels struct {
	Summary   SleepLevelsSummary     `json:"summary"`
	Data      []SleepLevelsDatapoint `json:"data"`
	ShortData []SleepLevelsDatapoint `json:"shortData"`
}

type SleepLevelsSummary struct {
	Deep  SleepLevel `json:"deep"`
	Light SleepLevel `json:"light"`
	Rem   SleepLevel `json:"rem"`
	Wake  SleepLevel `json:"wake"`
}

type SleepLevel struct {
	Count               int `json:"count"`
	Minutes             int `json:"minutes"`
	ThirtyDayAvgMinutes int `json:"thirtyDayAvgMinutes"`
}

type SleepLevelsDatapoint struct {
	Datetime string `json:"datetime"`
	Level    string `json:"level"`
	Seconds  int    `json:"seconds"`
}

type Activity struct {
	Activities []ActivityData `json:"activities"`
}

type ActivityData struct {
	ActivityID       int     `json:"activityId"`
	ActivityParentID int     `json:"activityParentId"`
	Calories         int     `json:"calories"`
	Description      string  `json:"description"`
	Distance         float32 `json:"distance"`
	Duration         int     `json:"duration"`
	HasStartTime     bool    `json:"hasStartTime"`
	IsFavorite       bool    `json:"isFavorite"`
	LogID            int     `json:"logId"`
	Name             string  `json:"name"`
	StartTime        string  `json:"startTime"`
	Steps            int     `json:"steps"`
}

type ActivityGoals struct {
	CaloriesOut int     `json:"caloriesOut"`
	Distance    float32 `json:"distance"`
	Floors      int     `json:"floors"`
	Steps       int     `json:"steps"`
}

type ActivitySummary struct {
	ActivityCalories     int                `json:"activityCalories"`
	CaloriesBMR          int                `json:"caloriesBMR"`
	CaloriesOut          int                `json:"caloriesOut"`
	Distances            []ActivityDistance `json:"distances"`
	Elevation            float32            `json:"elevation"`
	FairlyActiveMinutes  int                `json:"fairlyActiveMinutes"`
	Floors               int                `json:"floors"`
	LightlyActiveMinutes int                `json:"lightlyActiveMinutes"`
	MarginalCalories     int                `json:"marginalCalories"`
	SedentaryMinutes     int                `json:"sedentaryMinutes"`
	Steps                int                `json:"steps"`
	VeryActiveMinutes    int                `json:"veryActiveMinutes"`
}

type ActivityDistance struct {
	Activity string  `json:"activity"`
	Distance float32 `json:"distance"`
}
