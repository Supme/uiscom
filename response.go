package uiscom

import "strconv"

// Error response
//
//	{
//	 "jsonrpc": "2.0",
//	 "id": null,
//	 "error": {
//	   "code": "number",
//	   "message": "string",
//	   "data": {
//	     "mnemonic": "string",
//	     "field": "string",
//	     "value": "string",
//	     "params": {
//	       "object": "string"
//	     },
//	     "extended_helper": "string",
//	     "metadata": {
//
//	     }
//	   }
//	 }
//	}
type Error struct {
	Jsonrpc      string `json:"jsonrpc"`
	ID           any    `json:"id"`
	ErrorContent struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    struct {
			Mnemonic string `json:"mnemonic"`
			Field    string `json:"field"`
			Value    string `json:"value"`
			Params   struct {
				Object string `json:"object"`
			} `json:"params"`
			ExtendedHelper string   `json:"extended_helper"`
			Metadata       Metadata `json:"metadata"`
		} `json:"data"`
	} `json:"error"`
}

func (e *Error) Error() string {
	return strconv.Itoa(e.ErrorContent.Code) + " " + e.ErrorContent.Message
}

// MetaData
//
//	{
//	 "metadata":{
//	   "api_version":{
//	     "current_version_deprecated":"boolean",
//	     "current_version":"string",
//	     "latest_version":"string"
//	   },
//	   "limits":{
//	     "day_limit":"number",
//	     "day_remaining":"number",
//	     "day_reset":"number",
//	     "minute_limit":"number",
//	     "minute_remaining":"number",
//	     "minute_reset":"number"
//	   },
//	   "total_items":"number"
//	 }
//	}
type Metadata struct {
	APIVersion struct {
		CurrentVersionDeprecated string `json:"current_version_deprecated"`
		CurrentVersion           string `json:"current_version"`
		LatestVersion            string `json:"latest_version"`
	} `json:"api_version"`
	Limits struct {
		DayLimit        string `json:"day_limit"`
		DayRemaining    string `json:"day_remaining"`
		DayReset        string `json:"day_reset"`
		MinuteLimit     string `json:"minute_limit"`
		MinuteRemaining string `json:"minute_remaining"`
		MinuteReset     string `json:"minute_reset"`
	} `json:"limits"`
	TotalItems string `json:"total_items"`
}
