package typechat

import (
	"bufio"
	"fmt"
	"os"
)

func ProcessRequests(interactivePrompt string, inputFileName string, processRequest func(request string) error) error {
	file, err := os.Open(inputFileName)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Printf("%s%s\n", interactivePrompt, line)
		err = processRequest(line)
		if err != nil {
			return err
		}
	}

	return scanner.Err()
}
