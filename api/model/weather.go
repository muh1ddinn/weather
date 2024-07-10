package model

type WeatherResponse struct {
	ResolvedAddress string       `json:"resolvedAddress"`
	Address         string       `json:"address"`
	Timezone        string       `json:"timezone"`
	Days            []Day        `json:"days"`
	Hours           []HourlyData `json:"hours"`
}

type Day struct {
	DatetimeEpoch int64        `json:"datetimeEpoch"`
	Datetime      string       `json:"datetime"`
	Temp          float64      `json:"temp"`
	TempMax       float64      `json:"tempmax"`
	TempMin       float64      `json:"tempmin"`
	Humidity      float64      `json:"humidity"`
	Conditions    string       `json:"conditions"`
	Hours         []HourlyData `json:"hours"`
}

type HourlyData struct {
	Datetime      string  `json:"datetime"`
	DatetimeEpoch int64   `json:"datetimeEpoch"`
	Temp          float64 `json:"temp"`
	Humidity      float64 `json:"humidity"`
	Icon          string  `json:"icon"`
}

type Responseweathers_hour struct {
	HourlyData []HourlyData `json:"hourly_data"`
}

type GetAllWeatherRequest struct {
	Location string `json:"location"`
}

type GetAllResponseWeatherHour struct {
	WeatherData []WeatherResponse `json:"weatherData"`
}

type GetAllContactRequest struct {
	Search string `json:"search"`
	Page   uint64 `json:"page"`
	Limit  uint64 `json:"limit"`
}

type GetAllWeathername struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type GetAllWeatherRequestt struct {
	City      string  `json:"city"`
	Condition string  `json:"condition"`
	MinTemp   float64 `json:"min_temp"`
	MaxTemp   float64 `json:"max_temp"`
	Page      uint64  `json:"page"`
	Limit     uint64  `json:"limit"`
	StartTime string  `form:"start_time"`
	EndTime   string  `form:"end_time"`
}

type GetAllWeatherResponse struct {
	TotalCount int       `json:"total_count"`
	Weathers   []Weather `json:"weathers"`
}

type Weather struct {
	Location struct {
		Id      string `json:"id"`
		Name    string `json:"name"`
		Country string `json:"country"`
	} `json:"location"`
	Hour struct {
		Time        string  `json:"time"`
		Temperature float64 `json:"temperature"`
		Humidity    float64 `json:"humidity"`
		Condition   string  `json:"condition"`
	} `json:"hour"`
	Day struct {
		Time           string  `json:"time"`
		Temperature    float64 `json:"temperature"`
		TemperatureMax float64 `json:"temperature_max"`
		TemperatureMin float64 `json:"temperature_min"`
		Humidity       float64 `json:"humidity"`
		Condition      string  `json:"condition"`
	} `json:"day"`
}
