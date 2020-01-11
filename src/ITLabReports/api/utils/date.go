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
	return date[12:16]+"-"+months[date[8:11]]+"-"+date[5:7]+"T"+date[17:25]  //Example: 2019-Dec-07T22:34:31
}
func FormatQueryDate(date string) string {
	return date[6:10]+"-"+date[3:5]+"-"+date[0:2] //05.01.2020
}