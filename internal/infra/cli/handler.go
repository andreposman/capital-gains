package cli

import (
	"bufio"
	"bytes"
	json2 "encoding/json"
	"github.com/andreposman/capital-gains/internal/application"
	"github.com/andreposman/capital-gains/internal/infra/json"
	"log"
	"os"
)

var newLine = []byte{'\n'}

func Handle() {
	scanner := bufio.NewScanner(os.Stdin)
	processor := application.OperationProcessor{}

	for scanner.Scan() {
		line := scanner.Bytes()
		currLine := line

		if len(currLine) == 0 {
			break
		}

		operations, err := json.ParseInput(line)
		if err != nil {
			log.Fatalf("Error parsing input JSON: %v", err)
		}

		// if json is empty "[]"
		if len(operations) == 0 {
			_, writeErr := os.Stdout.Write([]byte("[]\n")) // Write empty array with newline
			if writeErr != nil {
				log.Fatalf("Error writing empty array to stdout: %v", writeErr)
			}
			continue //move to next line
		}

		//valid, non-empty json
		result := processor.ProcessOperations(operations)

		output, err := json2.Marshal(result)
		if err != nil {
			log.Fatalf("Error marshalling output JSON: %v", err)
		}
		//log.Println(string(output))

		outputWithNewline := bytes.Join([][]byte{output, newLine}, []byte{})
		_, writeErr := os.Stdout.Write(outputWithNewline)
		if writeErr != nil {
			log.Fatalf("Error writing result to stdout: %v", writeErr)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading standard input: %v", err)
	}
}
