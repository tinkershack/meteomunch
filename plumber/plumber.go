// Package plumber models meteo data
//
// It defines structures and methods to filter and move meteo data from multiple providers.
package plumber

type Coordinates struct {
	Latitude  float64 // in degrees
	Longitude float64 // in degrees
}

// Models like Copernicus GLO-90 provide elevation with an accuracy of <4 meters. So, it may not be necessary to capture accuracy separately.
type Elevation uint // in meters

// Location represents a geographical location with various attributes for geocoding
type Location struct {
	ID          int         `json:"id"`           // Unique ID for this location
	Name        string      `json:"name"`         // Location name, possibly localized
	Coordinates Coordinates `json:"coordinates"`  // Geographical WGS84 coordinates of this location
	Elevation   float64     `json:"elevation"`    // Elevation above mean sea level of this location
	Timezone    string      `json:"timezone"`     // Time zone using time zone database definitions https://en.wikipedia.org/wiki/List_of_tz_database_time_zones
	FeatureCode string      `json:"feature_code"` // Type of this location. Following the GeoNames feature_code definitions https://www.geonames.org/export/codes.html
	CountryCode string      `json:"country_code"` // 2-Character FIPS country code, example IN for India
	Country     string      `json:"country"`      // Country name
	CountryID   int         `json:"country_id"`   // Unique ID for this country
	Population  int         `json:"population"`   // Number of inhabitants
	Postcodes   []string    `json:"postcodes"`    // List of postcodes for this location
}

// BaseData is the main structure that holds rudimentary meteo data
type BaseData struct {
	Latitude             float64     `json:"latitude"`
	Longitude            float64     `json:"longitude"`
	UTCOffsetSeconds     int         `json:"utc_offset_seconds"`
	Timezone             string      `json:"timezone"`
	TimezoneAbbreviation string      `json:"timezone_abbreviation"`
	Elevation            int         `json:"elevation"`
	Current              CurrentData `json:"current"`
	Hourly               HourlyData  `json:"hourly"`
	Daily                DailyData   `json:"daily"`
}

/*
// SampleBaseData provides sample values for BaseData to serve as a reference
var SampleBaseData = BaseData{
    Latitude:             11.25, // Doesn't matter even if it's provided as 11.2855, most weather models normalize for 0.25 degree grid resolution
    Longitude:            77, // 76.9661
    UTCOffsetSeconds:     0,
    Timezone:             "GMT",
    TimezoneAbbreviation: "GMT",
    Elevation:            422,
    Current:              SampleCurrentData,
    Hourly:               HourlyData{}, // Placeholder for Hourly field
    Daily:                SampleDailyData,
}
*/

// CurrentData holds current meteo data for the time period that was requested
type CurrentData struct {
	Time                int64   `json:"time"`                 // Unix timestamp
	Interval            int     `json:"interval"`             // Interval in seconds for the data
	Temperature2M       float64 `json:"temperature_2m"`       // Temperature at 2 meters above ground in °C
	RelativeHumidity2M  int     `json:"relative_humidity_2m"` // Relative humidity at 2 meters above ground in %
	ApparentTemperature float64 `json:"apparent_temperature"` // Apparent temperature in °C
	IsDay               int     `json:"is_day"`               // 1 if it is day, 0 if it is night
	Precipitation       float64 `json:"precipitation"`        // Precipitation in mm
	Rain                float64 `json:"rain"`                 // Rain in mm
	Showers             float64 `json:"showers"`              // Showers in mm
	Snowfall            float64 `json:"snowfall"`             // Snowfall in cm
	WeatherCode         int     `json:"weather_code"`         // Weather code according to WMO
	CloudCover          int     `json:"cloud_cover"`          // Cloud cover in %
	PressureMSL         float64 `json:"pressure_msl"`         // Mean sea level pressure in hPa
	SurfacePressure     float64 `json:"surface_pressure"`     // Surface pressure in hPa
	WindSpeed10M        float64 `json:"wind_speed_10m"`       // Wind speed at 10 meters above ground in km/h
	WindDirection10M    int     `json:"wind_direction_10m"`   // Wind direction at 10 meters above ground in degrees
	WindGusts10M        float64 `json:"wind_gusts_10m"`       // Wind gusts at 10 meters above ground in km/h
}

/*
// SampleCurrentData provides sample values for CurrentData to sreve as a reference
var SampleCurrentData = CurrentData{
    Time:                1633036800, // Unix timestamp
    Interval:            3600,       // Interval in seconds
    Temperature2M:       15.3,       // Temperature at 2 meters above ground in °C
    RelativeHumidity2M:  80,         // Relative humidity at 2 meters above ground in %
    ApparentTemperature: 14.8,       // Apparent temperature in °C
    IsDay:               1,          // 1 if it is day, 0 if it is night
    Precipitation:       0.0,        // Precipitation in mm
    Rain:                0.0,        // Rain in mm
    Showers:             0.0,        // Showers in mm
    Snowfall:            0.0,        // Snowfall in cm
    WeatherCode:         1,          // Weather code according to WMO
    CloudCover:          20,         // Cloud cover in %
    PressureMSL:         1015.0,     // Mean sea level pressure in hPa
    SurfacePressure:     1013.0,     // Surface pressure in hPa
    WindSpeed10M:        5.0,        // Wind speed at 10 meters above ground in km/h
    WindDirection10M:    180,        // Wind direction at 10
}
*/

// HourlyData holds forecasted meteo data for a set of hourly intervals, normally 24 intervals, that was requested.
type HourlyData struct {
	Time                             []int64   `json:"time"` // Time intervals for which rest of the array fields' values are populated
	Temperature2M                    []float64 `json:"temperature_2m"`
	RelativeHumidity2M               []int     `json:"relative_humidity_2m"`
	DewPoint2M                       []float64 `json:"dew_point_2m"`
	ApparentTemperature              []float64 `json:"apparent_temperature"`
	PrecipitationProbability         []int     `json:"precipitation_probability"`
	Precipitation                    []float64 `json:"precipitation"`
	WeatherCode                      []int     `json:"weather_code"`
	PressureMSL                      []float64 `json:"pressure_msl"`
	SurfacePressure                  []float64 `json:"surface_pressure"`
	CloudCover                       []int     `json:"cloud_cover"`
	CloudCoverLow                    []int     `json:"cloud_cover_low"`
	CloudCoverMid                    []int     `json:"cloud_cover_mid"`
	CloudCoverHigh                   []int     `json:"cloud_cover_high"`
	Visibility                       []int     `json:"visibility"`
	Evapotranspiration               []float64 `json:"evapotranspiration"`
	ET0FAOEvapotranspiration         []float64 `json:"et0_fao_evapotranspiration"`
	VapourPressureDeficit            []float64 `json:"vapour_pressure_deficit"`
	WindSpeed10M                     []float64 `json:"wind_speed_10m"`
	WindSpeed80M                     []float64 `json:"wind_speed_80m"`
	WindSpeed120M                    []float64 `json:"wind_speed_120m"`
	WindSpeed180M                    []float64 `json:"wind_speed_180m"`
	WindDirection10M                 []int     `json:"wind_direction_10m"`
	WindDirection80M                 []int     `json:"wind_direction_80m"`
	WindDirection120M                []int     `json:"wind_direction_120m"`
	WindDirection180M                []int     `json:"wind_direction_180m"`
	WindGusts10M                     []float64 `json:"wind_gusts_10m"`
	Temperature80M                   []float64 `json:"temperature_80m"`
	Temperature120M                  []float64 `json:"temperature_120m"`
	Temperature180M                  []float64 `json:"temperature_180m"`
	UVIndex                          []float64 `json:"uv_index"`
	UVIndexClearSky                  []float64 `json:"uv_index_clear_sky"`
	IsDay                            []int     `json:"is_day"`
	SunshineDuration                 []float64 `json:"sunshine_duration"`
	TotalColumnIntegratedWaterVapour []float64 `json:"total_column_integrated_water_vapour"`
	Cape                             []float64 `json:"cape"`
	LiftedIndex                      []float64 `json:"lifted_index"`
	ConvectiveInhibition             []float64 `json:"convective_inhibition"`
	FreezingLevelHeight              []float64 `json:"freezing_level_height"`
	BoundaryLayerHeight              []float64 `json:"boundary_layer_height"`
	Temperature1000hPa               []float64 `json:"temperature_1000hpa"`
	Temperature975hPa                []float64 `json:"temperature_975hpa"`
	Temperature950hPa                []float64 `json:"temperature_950hpa"`
	Temperature925hPa                []float64 `json:"temperature_925hpa"`
	Temperature900hPa                []float64 `json:"temperature_900hpa"`
	Temperature850hPa                []float64 `json:"temperature_850hpa"`
	Temperature800hPa                []float64 `json:"temperature_800hpa"`
	Temperature700hPa                []float64 `json:"temperature_700hpa"`
	Temperature600hPa                []float64 `json:"temperature_600hpa"`
	Temperature500hPa                []float64 `json:"temperature_500hpa"`
	Temperature400hPa                []float64 `json:"temperature_400hpa"`
	RelativeHumidity1000hPa          []int     `json:"relative_humidity_1000hpa"`
	RelativeHumidity975hPa           []int     `json:"relative_humidity_975hpa"`
	RelativeHumidity950hPa           []int     `json:"relative_humidity_950hpa"`
	RelativeHumidity925hPa           []int     `json:"relative_humidity_925hpa"`
	RelativeHumidity900hPa           []int     `json:"relative_humidity_900hpa"`
	RelativeHumidity850hPa           []int     `json:"relative_humidity_850hpa"`
	RelativeHumidity800hPa           []int     `json:"relative_humidity_800hpa"`
	RelativeHumidity700hPa           []int     `json:"relative_humidity_700hpa"`
	RelativeHumidity600hPa           []int     `json:"relative_humidity_600hpa"`
	RelativeHumidity500hPa           []int     `json:"relative_humidity_500hpa"`
	RelativeHumidity400hPa           []int     `json:"relative_humidity_400hpa"`
	CloudCover1000hPa                []int     `json:"cloud_cover_1000hpa"`
	CloudCover975hPa                 []int     `json:"cloud_cover_975hpa"`
	CloudCover950hPa                 []int     `json:"cloud_cover_950hpa"`
	CloudCover925hPa                 []int     `json:"cloud_cover_925hpa"`
	CloudCover900hPa                 []int     `json:"cloud_cover_900hpa"`
	CloudCover850hPa                 []int     `json:"cloud_cover_850hpa"`
	CloudCover800hPa                 []int     `json:"cloud_cover_800hpa"`
	CloudCover700hPa                 []int     `json:"cloud_cover_700hpa"`
	CloudCover600hPa                 []int     `json:"cloud_cover_600hpa"`
	CloudCover500hPa                 []int     `json:"cloud_cover_500hpa"`
	CloudCover400hPa                 []int     `json:"cloud_cover_400hpa"`
	WindSpeed1000hPa                 []float64 `json:"wind_speed_1000hpa"`
	WindSpeed975hPa                  []float64 `json:"wind_speed_975hpa"`
	WindSpeed950hPa                  []float64 `json:"wind_speed_950hpa"`
	WindSpeed925hPa                  []float64 `json:"wind_speed_925hpa"`
	WindSpeed900hPa                  []float64 `json:"wind_speed_900hpa"`
	WindSpeed850hPa                  []float64 `json:"wind_speed_850hpa"`
	WindSpeed800hPa                  []float64 `json:"wind_speed_800hpa"`
	WindSpeed700hPa                  []float64 `json:"wind_speed_700hpa"`
	WindSpeed600hPa                  []float64 `json:"wind_speed_600hpa"`
	WindSpeed500hPa                  []float64 `json:"wind_speed_500hpa"`
	WindSpeed400hPa                  []float64 `json:"wind_speed_400hpa"`
	WindDirection1000hPa             []int     `json:"wind_direction_1000hpa"`
	WindDirection975hPa              []int     `json:"wind_direction_975hpa"`
	WindDirection950hPa              []int     `json:"wind_direction_950hpa"`
	WindDirection925hPa              []int     `json:"wind_direction_925hpa"`
	WindDirection900hPa              []int     `json:"wind_direction_900hpa"`
	WindDirection850hPa              []int     `json:"wind_direction_850hpa"`
	WindDirection800hPa              []int     `json:"wind_direction_800hpa"`
	WindDirection700hPa              []int     `json:"wind_direction_700hpa"`
	WindDirection600hPa              []int     `json:"wind_direction_600hpa"`
	WindDirection500hPa              []int     `json:"wind_direction_500hpa"`
	WindDirection400hPa              []int     `json:"wind_direction_400hpa"`
	GeopotentialHeight1000hPa        []float64 `json:"geopotential_height_1000hpa"`
	GeopotentialHeight975hPa         []float64 `json:"geopotential_height_975hpa"`
	GeopotentialHeight950hPa         []float64 `json:"geopotential_height_950hpa"`
	GeopotentialHeight925hPa         []float64 `json:"geopotential_height_925hpa"`
	GeopotentialHeight900hPa         []float64 `json:"geopotential_height_900hpa"`
	GeopotentialHeight850hPa         []float64 `json:"geopotential_height_850hpa"`
	GeopotentialHeight800hPa         []float64 `json:"geopotential_height_800hpa"`
	GeopotentialHeight700hPa         []float64 `json:"geopotential_height_700hpa"`
	GeopotentialHeight600hPa         []float64 `json:"geopotential_height_600hpa"`
	GeopotentialHeight500hPa         []float64 `json:"geopotential_height_500hpa"`
	GeopotentialHeight400hPa         []float64 `json:"geopotential_height_400hpa"`
}

// DailyData holds forecasted meteo data that provides a higher level trend for the day
type DailyData struct {
	Time                        []int64   `json:"time"`                          // Unix timestamps
	WeatherCode                 []int     `json:"weather_code"`                  // Weather codes according to WMO
	Temperature2MMax            []float64 `json:"temperature_2m_max"`            // Maximum temperature at 2 meters above ground in °C
	Temperature2MMin            []float64 `json:"temperature_2m_min"`            // Minimum temperature at 2 meters above ground in °C
	ApparentTemperatureMax      []float64 `json:"apparent_temperature_max"`      // Maximum apparent temperature in °C
	ApparentTemperatureMin      []float64 `json:"apparent_temperature_min"`      // Minimum apparent temperature in °C
	Sunrise                     []int64   `json:"sunrise"`                       // Unix timestamps for sunrise
	Sunset                      []int64   `json:"sunset"`                        // Unix timestamps for sunset
	DaylightDuration            []float64 `json:"daylight_duration"`             // Duration of daylight in seconds
	SunshineDuration            []float64 `json:"sunshine_duration"`             // Duration of sunshine in seconds
	UVIndexMax                  []float64 `json:"uv_index_max"`                  // Maximum UV index
	UVIndexClearSkyMax          []float64 `json:"uv_index_clear_sky_max"`        // Maximum UV index under clear sky
	PrecipitationSum            []float64 `json:"precipitation_sum"`             // Total precipitation in mm
	PrecipitationHours          []float64 `json:"precipitation_hours"`           // Hours of precipitation
	PrecipitationProbabilityMax []int     `json:"precipitation_probability_max"` // Maximum probability of precipitation in %
	WindSpeed10MMax             []float64 `json:"wind_speed_10m_max"`            // Maximum wind speed at 10 meters above ground in km/h
	WindGusts10MMax             []float64 `json:"wind_gusts_10m_max"`            // Maximum wind gusts at 10 meters above ground in km/h
	WindDirection10MDominant    []int     `json:"wind_direction_10m_dominant"`   // Dominant wind direction at 10 meters above ground in degrees
	ShortwaveRadiationSum       []float64 `json:"shortwave_radiation_sum"`       // Sum of shortwave radiation in MJ/m²
	ET0FAOEvapotranspiration    []float64 `json:"et0_fao_evapotranspiration"`    // Evapotranspiration in mm
}

/*
// SampleDailyData provides sample values for DailyData
var SampleDailyData = DailyData{
    Time:                        []int64{1633036800}, // Unix timestamp
    WeatherCode:                 []int{1},            // Weather code according to WMO
    Temperature2MMax:            []float64{20.0},     // Maximum temperature at 2 meters above ground in °C
    Temperature2MMin:            []float64{10.0},     // Minimum temperature at 2 meters above ground in °C
    ApparentTemperatureMax:      []float64{19.5},     // Maximum apparent temperature in °C
    ApparentTemperatureMin:      []float64{9.5},      // Minimum apparent temperature in °C
    Sunrise:                     []int64{1633065600}, // Unix timestamp for sunrise
    Sunset:                      []int64{1633108800}, // Unix timestamp for sunset
    DaylightDuration:            []float64{43200},    // Daylight duration in seconds
    SunshineDuration:            []float64{36000},    // Sunshine duration in seconds
    UVIndexMax:                  []float64{5.0},      // Maximum UV index
    UVIndexClearSkyMax:          []float64{6.0},      // Maximum UV index under clear sky
    PrecipitationSum:            []float64{0.0},      // Total precipitation in mm
    PrecipitationHours:          []float64{0.0},      // Hours of precipitation
    PrecipitationProbabilityMax: []int{0},            // Maximum probability of precipitation in %
    WindSpeed10MMax:             []float64{10.0},     // Maximum wind speed at 10 meters above ground in km/h
    WindGusts10MMax:             []float64{15.0},     // Maximum wind gusts at 10 meters above ground in km/h
    WindDirection10MDominant:    []int{180},          // Dominant wind direction at 10 meters above ground in degrees
    ShortwaveRadiationSum:       []float64{20.0},     // Sum of shortwave radiation in MJ/m²
    ET0FAOEvapotranspiration:    []float64{3.0},      // Evapotranspiration in mm
}
*/

// CommonUnits is a map of common units used in meteo data
// Providers may have different preferences for units, this map is used to standardize the units
var CommonUnits = map[string]string{
	"Temperature":               "°C",
	"WindSpeed":                 "km/h",
	"Humidity":                  "%",
	"Pressure":                  "hPa",
	"Visibility":                "km",
	"Precipitation":             "mm",
	"CloudCover":                "%",
	"SunshineHours":             "s",
	"IsDay":                     "",
	"Time":                      "unixtime",
	"Interval":                  "seconds",
	"WindGust":                  "km/h",
	"WindDirection":             "°",
	"DewPoint":                  "°C",
	"UVIndex":                   "",
	"WeatherCode":               "wmo code",
	"Snowfall":                  "cm",
	"SnowDepth":                 "cm",
	"Rain":                      "mm",
	"FreezingLevel":             "m",
	"SoilTemperature":           "°C",
	"SoilMoisture":              "m³/m³",
	"Sunrise":                   "unixtime",
	"Sunset":                    "unixtime",
	"DaylightDuration":          "s",
	"PrecipitationProbability":  "%",
	"Evapotranspiration":        "mm",
	"ET0FAOEvapotranspiration":  "mm",
	"VapourPressureDeficit":     "kPa",
	"GeopotentialHeight":        "m",
	"ShortwaveRadiationSum":     "MJ/m²",
	"ConvectiveInhibition":      "J/kg",
	"LiftedIndex":               "",
	"BoundaryLayerHeight":       "m",
	"Temperature1000hPa":        "°C",
	"Temperature975hPa":         "°C",
	"Temperature950hPa":         "°C",
	"Temperature925hPa":         "°C",
	"Temperature900hPa":         "°C",
	"Temperature850hPa":         "°C",
	"Temperature800hPa":         "°C",
	"Temperature700hPa":         "°C",
	"Temperature600hPa":         "°C",
	"Temperature500hPa":         "°C",
	"Temperature400hPa":         "°C",
	"RelativeHumidity1000hPa":   "%",
	"RelativeHumidity975hPa":    "%",
	"RelativeHumidity950hPa":    "%",
	"RelativeHumidity925hPa":    "%",
	"RelativeHumidity900hPa":    "%",
	"RelativeHumidity850hPa":    "%",
	"RelativeHumidity800hPa":    "%",
	"RelativeHumidity700hPa":    "%",
	"RelativeHumidity600hPa":    "%",
	"RelativeHumidity500hPa":    "%",
	"RelativeHumidity400hPa":    "%",
	"CloudCover1000hPa":         "%",
	"CloudCover975hPa":          "%",
	"CloudCover950hPa":          "%",
	"CloudCover925hPa":          "%",
	"CloudCover900hPa":          "%",
	"CloudCover850hPa":          "%",
	"CloudCover800hPa":          "%",
	"CloudCover700hPa":          "%",
	"CloudCover600hPa":          "%",
	"CloudCover500hPa":          "%",
	"CloudCover400hPa":          "%",
	"WindSpeed1000hPa":          "km/h",
	"WindSpeed975hPa":           "km/h",
	"WindSpeed950hPa":           "km/h",
	"WindSpeed925hPa":           "km/h",
	"WindSpeed900hPa":           "km/h",
	"WindSpeed850hPa":           "km/h",
	"WindSpeed800hPa":           "km/h",
	"WindSpeed700hPa":           "km/h",
	"WindSpeed600hPa":           "km/h",
	"WindSpeed500hPa":           "km/h",
	"WindSpeed400hPa":           "km/h",
	"WindDirection1000hPa":      "°",
	"WindDirection975hPa":       "°",
	"WindDirection950hPa":       "°",
	"WindDirection925hPa":       "°",
	"WindDirection900hPa":       "°",
	"WindDirection850hPa":       "°",
	"WindDirection800hPa":       "°",
	"WindDirection700hPa":       "°",
	"WindDirection600hPa":       "°",
	"WindDirection500hPa":       "°",
	"WindDirection400hPa":       "°",
	"GeopotentialHeight1000hPa": "m",
	"GeopotentialHeight975hPa":  "m",
	"GeopotentialHeight950hPa":  "m",
	"GeopotentialHeight925hPa":  "m",
	"GeopotentialHeight900hPa":  "m",
	"GeopotentialHeight850hPa":  "m",
	"GeopotentialHeight800hPa":  "m",
	"GeopotentialHeight700hPa":  "m",
	"GeopotentialHeight600hPa":  "m",
	"GeopotentialHeight500hPa":  "m",
	"GeopotentialHeight400hPa":  "m",
}

// CurrentUnits is a map of units used in CurrentData structure's meteo data
var CurrentUnits = map[string]string{
	"Temperature":     CommonUnits["Temperature"],
	"WindSpeed":       CommonUnits["WindSpeed"],
	"Humidity":        CommonUnits["Humidity"],
	"Pressure":        CommonUnits["Pressure"],
	"Visibility":      CommonUnits["Visibility"],
	"IsDay":           CommonUnits["IsDay"],
	"Time":            CommonUnits["Time"],
	"WindGust":        CommonUnits["WindGust"],
	"WindDirection":   CommonUnits["WindDirection"],
	"DewPoint":        CommonUnits["DewPoint"],
	"UVIndex":         CommonUnits["UVIndex"],
	"WeatherCode":     CommonUnits["WeatherCode"],
	"Snowfall":        CommonUnits["Snowfall"],
	"SnowDepth":       CommonUnits["SnowDepth"],
	"Rain":            CommonUnits["Rain"],
	"FreezingLevel":   CommonUnits["FreezingLevel"],
	"SoilTemperature": CommonUnits["SoilTemperature"],
	"SoilMoisture":    CommonUnits["SoilMoisture"],
}

// HourlyUnits is a map of units used in HourlyData structure's meteo data
var HourlyUnits = map[string]string{
	"Temperature":               CommonUnits["Temperature"],
	"WindSpeed":                 CommonUnits["WindSpeed"],
	"Humidity":                  CommonUnits["Humidity"],
	"Pressure":                  CommonUnits["Pressure"],
	"Visibility":                CommonUnits["Visibility"],
	"Precipitation":             CommonUnits["Precipitation"],
	"CloudCover":                CommonUnits["CloudCover"],
	"IsDay":                     CommonUnits["IsDay"],
	"Time":                      CommonUnits["Time"],
	"Interval":                  CommonUnits["Interval"],
	"WindGust":                  CommonUnits["WindGust"],
	"WindDirection":             CommonUnits["WindDirection"],
	"DewPoint":                  CommonUnits["DewPoint"],
	"UVIndex":                   CommonUnits["UVIndex"],
	"WeatherCode":               CommonUnits["WeatherCode"],
	"Snowfall":                  CommonUnits["Snowfall"],
	"SnowDepth":                 CommonUnits["SnowDepth"],
	"Rain":                      CommonUnits["Rain"],
	"FreezingLevel":             CommonUnits["FreezingLevel"],
	"SoilTemperature":           CommonUnits["SoilTemperature"],
	"SoilMoisture":              CommonUnits["SoilMoisture"],
	"PrecipitationProbability":  CommonUnits["PrecipitationProbability"],
	"Evapotranspiration":        CommonUnits["Evapotranspiration"],
	"ET0FAOEvapotranspiration":  CommonUnits["ET0FAOEvapotranspiration"],
	"VapourPressureDeficit":     CommonUnits["VapourPressureDeficit"],
	"GeopotentialHeight":        CommonUnits["GeopotentialHeight"],
	"ShortwaveRadiationSum":     CommonUnits["ShortwaveRadiationSum"],
	"ConvectiveInhibition":      CommonUnits["ConvectiveInhibition"],
	"LiftedIndex":               CommonUnits["LiftedIndex"],
	"BoundaryLayerHeight":       CommonUnits["BoundaryLayerHeight"],
	"Temperature1000hPa":        CommonUnits["Temperature1000hPa"],
	"Temperature975hPa":         CommonUnits["Temperature975hPa"],
	"Temperature950hPa":         CommonUnits["Temperature950hPa"],
	"Temperature925hPa":         CommonUnits["Temperature925hPa"],
	"Temperature900hPa":         CommonUnits["Temperature900hPa"],
	"Temperature850hPa":         CommonUnits["Temperature850hPa"],
	"Temperature800hPa":         CommonUnits["Temperature800hPa"],
	"Temperature700hPa":         CommonUnits["Temperature700hPa"],
	"Temperature600hPa":         CommonUnits["Temperature600hPa"],
	"Temperature500hPa":         CommonUnits["Temperature500hPa"],
	"Temperature400hPa":         CommonUnits["Temperature400hPa"],
	"RelativeHumidity1000hPa":   CommonUnits["RelativeHumidity1000hPa"],
	"RelativeHumidity975hPa":    CommonUnits["RelativeHumidity975hPa"],
	"RelativeHumidity950hPa":    CommonUnits["RelativeHumidity950hPa"],
	"RelativeHumidity925hPa":    CommonUnits["RelativeHumidity925hPa"],
	"RelativeHumidity900hPa":    CommonUnits["RelativeHumidity900hPa"],
	"RelativeHumidity850hPa":    CommonUnits["RelativeHumidity850hPa"],
	"RelativeHumidity800hPa":    CommonUnits["RelativeHumidity800hPa"],
	"RelativeHumidity700hPa":    CommonUnits["RelativeHumidity700hPa"],
	"RelativeHumidity600hPa":    CommonUnits["RelativeHumidity600hPa"],
	"RelativeHumidity500hPa":    CommonUnits["RelativeHumidity500hPa"],
	"RelativeHumidity400hPa":    CommonUnits["RelativeHumidity400hPa"],
	"CloudCover1000hPa":         CommonUnits["CloudCover1000hPa"],
	"CloudCover975hPa":          CommonUnits["CloudCover975hPa"],
	"CloudCover950hPa":          CommonUnits["CloudCover950hPa"],
	"CloudCover925hPa":          CommonUnits["CloudCover925hPa"],
	"CloudCover900hPa":          CommonUnits["CloudCover900hPa"],
	"CloudCover850hPa":          CommonUnits["CloudCover850hPa"],
	"CloudCover800hPa":          CommonUnits["CloudCover800hPa"],
	"CloudCover700hPa":          CommonUnits["CloudCover700hPa"],
	"CloudCover600hPa":          CommonUnits["CloudCover600hPa"],
	"CloudCover500hPa":          CommonUnits["CloudCover500hPa"],
	"CloudCover400hPa":          CommonUnits["CloudCover400hPa"],
	"WindSpeed1000hPa":          CommonUnits["WindSpeed1000hPa"],
	"WindSpeed975hPa":           CommonUnits["WindSpeed975hPa"],
	"WindSpeed950hPa":           CommonUnits["WindSpeed950hPa"],
	"WindSpeed925hPa":           CommonUnits["WindSpeed925hPa"],
	"WindSpeed900hPa":           CommonUnits["WindSpeed900hPa"],
	"WindSpeed850hPa":           CommonUnits["WindSpeed850hPa"],
	"WindSpeed800hPa":           CommonUnits["WindSpeed800hPa"],
	"WindSpeed700hPa":           CommonUnits["WindSpeed700hPa"],
	"WindSpeed600hPa":           CommonUnits["WindSpeed600hPa"],
	"WindSpeed500hPa":           CommonUnits["WindSpeed500hPa"],
	"WindSpeed400hPa":           CommonUnits["WindSpeed400hPa"],
	"WindDirection1000hPa":      CommonUnits["WindDirection1000hPa"],
	"WindDirection975hPa":       CommonUnits["WindDirection975hPa"],
	"WindDirection950hPa":       CommonUnits["WindDirection950hPa"],
	"WindDirection925hPa":       CommonUnits["WindDirection925hPa"],
	"WindDirection900hPa":       CommonUnits["WindDirection900hPa"],
	"WindDirection850hPa":       CommonUnits["WindDirection850hPa"],
	"WindDirection800hPa":       CommonUnits["WindDirection800hPa"],
	"WindDirection700hPa":       CommonUnits["WindDirection700hPa"],
	"WindDirection600hPa":       CommonUnits["WindDirection600hPa"],
	"WindDirection500hPa":       CommonUnits["WindDirection500hPa"],
	"WindDirection400hPa":       CommonUnits["WindDirection400hPa"],
	"GeopotentialHeight1000hPa": CommonUnits["GeopotentialHeight1000hPa"],
	"GeopotentialHeight975hPa":  CommonUnits["GeopotentialHeight975hPa"],
	"GeopotentialHeight950hPa":  CommonUnits["GeopotentialHeight950hPa"],
	"GeopotentialHeight925hPa":  CommonUnits["GeopotentialHeight925hPa"],
	"GeopotentialHeight900hPa":  CommonUnits["GeopotentialHeight900hPa"],
	"GeopotentialHeight850hPa":  CommonUnits["GeopotentialHeight850hPa"],
	"GeopotentialHeight800hPa":  CommonUnits["GeopotentialHeight800hPa"],
	"GeopotentialHeight700hPa":  CommonUnits["GeopotentialHeight700hPa"],
	"GeopotentialHeight600hPa":  CommonUnits["GeopotentialHeight600hPa"],
	"GeopotentialHeight500hPa":  CommonUnits["GeopotentialHeight500hPa"],
	"GeopotentialHeight400hPa":  CommonUnits["GeopotentialHeight400hPa"],
}

// DailyUnits is a map of units used in DailyData structure's meteo data
var DailyUnits = map[string]string{
	"Time":                        CommonUnits["Time"],
	"WeatherCode":                 CommonUnits["WeatherCode"],
	"TemperatureMax":              CommonUnits["Temperature"],
	"TemperatureMin":              CommonUnits["Temperature"],
	"ApparentTemperatureMax":      CommonUnits["Temperature"],
	"ApparentTemperatureMin":      CommonUnits["Temperature"],
	"Sunrise":                     CommonUnits["Sunrise"],
	"Sunset":                      CommonUnits["Sunset"],
	"DaylightDuration":            CommonUnits["DaylightDuration"],
	"SunshineDuration":            CommonUnits["SunshineHours"],
	"UVIndexMax":                  CommonUnits["UVIndex"],
	"UVIndexClearSkyMax":          CommonUnits["UVIndex"],
	"PrecipitationSum":            CommonUnits["Precipitation"],
	"PrecipitationHours":          "h",
	"PrecipitationProbabilityMax": CommonUnits["PrecipitationProbability"],
	"WindSpeedMax":                CommonUnits["WindSpeed"],
	"WindGustsMax":                CommonUnits["WindGust"],
	"WindDirectionDominant":       CommonUnits["WindDirection"],
	"ShortwaveRadiationSum":       CommonUnits["ShortwaveRadiationSum"],
	"ET0FAOEvapotranspiration":    CommonUnits["ET0FAOEvapotranspiration"],
}
