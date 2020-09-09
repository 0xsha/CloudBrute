package internal

import (
	"bufio"
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
