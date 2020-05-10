package Process
import (
	//"encoding/json"
	//"io/ioutil"
	"os"
)

func WriteJson(m []Movie_info){
	//file, _ := json.MarshalIndent(m, "", " ")

	file, _ := os.OpenFile("test.txt",os.O_APPEND|os.O_WRONLY,0644)

	defer file.Close()
	for _,value := range m{
		file.WriteString(value.Title)
	}

	file.Close()

}