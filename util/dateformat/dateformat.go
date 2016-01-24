package dateformat

import (
    "fmt"
    "strings"
    "errors"
)



func getMonth(month string)(string){
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
    return Months[month]
}
// Regresa el formato de la fecha para comparar vs el formato de Twitter
func DateFormat(str string)(string, error)  {
    if str == "" {
        return "", errors.New("Empty string date")
    }
    strslice := strings.Split(str, " ")
    return fmt.Sprintf("%s-%s-%s", strslice[2], getMonth(strslice[1]), strslice[5]), nil
}