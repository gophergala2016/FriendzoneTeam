package util

import "fmt"
import "strings"
import "errors"



func getMonth(month string){
    Months := map[string]string{
        " ": "---",
        "Jan": "January",
        "Feb": "February", 
        "Mar": "March", 
        "Apr": "April", 
        "May": "May", 
        "Jun": "June", 
        "Jul": "July", 
        "Aug": "August", 
        "Sep": "September", 
        "Oct": "October", 
        "Nov": "November", 
        "Dec": "December",
    }
    fmt.Print(Months[month])
}
func DateFormat(str string)(string, error)  {
    if str == "" {
        return "", errors.New("Empty string date")
    }
    strslice := strings.Split(str, " ")
    getMonth(strslice[1])
    return "", nil
    
}