package logs

import (
	"fmt"
	"log"
	"os"
	"time"
)

func WriteLogFile(a ...any) {
	file, err := os.OpenFile("./logs.txt", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if !Check(err) {
		Terminate()
	}
	defer file.Close()
	log.SetOutput(file)

	for i := range a {
		log.Printf("%v", a[i])
		fmt.Printf("%v\n", a[i])
	}
}

func WriteErrorLogFile(mainErr error) {
	file, err := os.OpenFile("./errors.txt", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if !Check(err) {
		Terminate()
	}

	log.SetOutput(file)
	log.Printf("ERROR: %v\n", mainErr)
	file.Close()

	file, err = os.OpenFile("./logs.txt", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if !Check(err) {
		Terminate()
	}
	defer file.Close()

	log.SetOutput(file)
	log.Printf("ERROR: %v", mainErr)
	fmt.Printf("ERROR: %v\n", mainErr)
}

func CheckDataRetention(daysToKeep int, fileName string) {
	info, err := os.Stat(fileName)
	if !Check(err) {
		Terminate()
	}

	modTime := info.ModTime()
	deletionTime := time.Now().AddDate(0, 0, -daysToKeep)

	if modTime.Before(deletionTime) || daysToKeep == 0 {
		file, err := os.Create(fileName)
		Check(err)
		defer file.Close()
		WriteLogFile(fmt.Sprintf("The file %v has been cleared because the data stored is older than %v days.\n",
			fileName, daysToKeep))
	}
}

func Check(err error) bool {
	if err != nil {
		WriteErrorLogFile(err)
		return false
	}
	return true
}

func Terminate() {
	fmt.Print("Press 'ENTER' to exit the program.")
	fmt.Scanln()
	os.Exit(1)
}
