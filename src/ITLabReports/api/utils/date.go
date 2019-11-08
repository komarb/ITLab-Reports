package utils

var months = map[string]string {
	"Jan":	"01",
	"Feb":	"02",
	"Mar":	"03",
	"Apr":	"04",
	"May":	"05",
	"Jun":	"06",
	"Jul":	"07",
	"Aug":	"08",
	"Sep":	"09",
	"Oct":	"10",
	"Nov":	"11",
	"Dec":	"12",
}

func FormatDate(date string) string {
	return date[12:16]+"-"+months[date[8:11]]+"-"+date[5:7]+"T"+date[17:25]
}
