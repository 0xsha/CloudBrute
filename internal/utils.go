package internal

import (
	"bufio"
	"fmt"
	"github.com/rs/zerolog/log"
	"math/rand"
	"os"
	"time"
)

func ReadTextFile(path string)  ([]string , error)  {

	 var buffer []string

     file , err := os.Open(path)
     if err!=nil{
		 log.Fatal().Err(err).Msg("Exiting ...")
	 }
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		buffer = append(buffer, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil,err

	}
	return buffer , nil

}

func SelectRandomItem(agents []string)  string {

	rand.Seed(time.Now().UnixNano())

	randomIndex := rand.Intn(len(agents))
	chosen := agents[randomIndex]

	return chosen

}

func WriteResultsToFile(results []string, output string)  {


	file, err := os.OpenFile(output+".txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer file.Close()

	if err != nil {
		log.Fatal().Err(err).Msg("failed creating file")
	}

	lineWriter := bufio.NewWriter(file)

	for _, result := range results {
		_, _ = lineWriter.WriteString(result + "\n")
	}

	lineWriter.Flush()


}

func Unique(input []string) []string {
	unique := make(map[string]bool, len(input))
	list := make([]string, len(unique))
	for _, el := range input {
		if len(el) != 0 {
			if !unique[el] {
				list = append(list, el)
				unique[el] = true
			}
		}
	}
	return list
}


func GenerateOutputName(output string) string  {

	t := time.Now()
	result := fmt.Sprintf("%s-%d-%02d-%02dT%02d-%02d-%02d",
		output, t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())

	return result
}