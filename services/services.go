package service
import (
    "fmt"
    "log"
    "time"
    "net/http"
    "net/url"
    "github.com/ChimeraCoder/anaconda"
)

// Revisa DMs en base a la fecha actual
func RevisarDM(w http.ResponseWriter, r *http.Request){
    // Verifica si ha obtenido nuevos DMs
    anaconda.SetConsumerKey("qFwSZ5VTSw2E23u166iWPAZj1")
    anaconda.SetConsumerSecret("AgJPsgKrP1OW3fx9TpDwZ55beUMGdBLcybJLyUvgMJM0tQj5PD")
    api := anaconda.NewTwitterApi("151138456-QnyorQfe59zld78ATJdfwNpcpsywVPZqcqo5uXiZ", "TpG0upWDAAmXF81K0HscQOxXCqHfpWEYd66c7fVdLX1k8")
    // Revisamos que ejecute bien la conexion con las credenciales
    if api.Credentials == nil {
        fmt.Println("Twitter Api client has empty (nil) credentials")
    }
    uri := url.Values{}
    dmResults, err := api.GetDirectMessages(uri)
    if err != nil {
        log.Printf("Error: %s", err.Error())
    }
    currentDate := time.Now()
    fmt.Printf(currentDate.Format("dd mm yyyy"))
    for _, message := range dmResults {
        fmt.Println(message.CreatedAt)
    }
    
    fmt.Fprintf(w, "hi")
}