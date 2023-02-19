package fitbit2gcal

import "golang.org/x/oauth2"

type FitbitConfig struct {
	OauthConfig *oauth2.Config
}

type GCalConfig struct {
	SleepCalendarID    string
	ActivityCalendarID string
	OauthConfig        *oauth2.Config
}

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
	Activities      []ActivityData  `json:"activities"`
	ActivityGoals   ActivityGoals   `json:"goals"`
	ActivitySummary ActivitySummary `json:"summary"`
}

type ActivityData struct {
	ActivityID        int             `json:"activityId"`
	ActiveDuration    int             `json:"activeDuration"`
	ActivityLevel     []ActivityLevel `json:"activityLevel"`
	ActivityName      string          `json:"activityName"`
	ActivityTypeID    int             `json:"activityTypeId"`
	Calories          int             `json:"calories"`
	Duration          int             `json:"duration"`
	ElevationGain     int             `json:"elevationGain"`
	HasGPS            bool            `json:"hasGps"`
	LastModified      string          `json:"lastModified"`
	LogID             int             `json:"logId"`
	LogType           string          `json:"logType"`
	OriginalDuration  int             `json:"originalDuration"`
	OriginalStartTime string          `json:"originakStartTime"`
	StartTime         string          `json:"startTime"`
	Steps             int             `json:"steps"`
}

type ActivityGoals struct {
	ActiveMinutes int     `json:"activeMinutes"`
	Calories      int     `json:"calories"`
	Distance      float32 `json:"distance"`
	DistanceUnit  string  `json:"distanceUnit"`
	Floors        int     `json:"floors"`
	Steps         int     `json:"steps"`
}

type ActivityLevel struct {
	Distance float32 `json:"distance"`
	Minutes  int     `json:"minutes"`
	Name     string  `json:"name"`
}

type ActivitySummary struct {
	ActivityLevels []ActivityLevel         `json:"activityLevels"`
	Calories       ActivityCalories        `json:"calories"`
	Distance       float32                 `json:"distance"`
	Elevation      float32                 `json:"elevation"`
	Floors         int                     `json:"floors"`
	HeartRateZones []ActivityHeartRateZone `json:"heartRateZones"`
	Steps          int                     `json:"steps"`
}

type ActivityCalories struct {
	BMR   int `json:"bmr"`
	Total int `json:"total"`
}

type ActivityHeartRateZone struct {
	CaloriesOut float32 `json:"caloriesOut"`
	Max         int     `json:"max"`
	Min         int     `json:"min"`
	Minutes     int     `json:"minutes"`
	Name        string  `json:"name"`
}
