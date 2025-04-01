package cli

import (
	"bufio"
	json2 "encoding/json"
	"github.com/andreposman/capital-gains/internal/application"
	"github.com/andreposman/capital-gains/internal/infra/json"
	"log"
	"os"
)

func Handle() {
	scanner := bufio.NewScanner(os.Stdin)
	processor := application.OperationProcessor{}

	for scanner.Scan() {
		line := scanner.Bytes()

		if len(line) == 0 {
			break
		}

		operations, err := json.ParseInput(line)
		if err != nil {
			log.Fatal(err)
		}

		result := processor.ProcessOperations(operations)

		output, err := json2.Marshal(result)
		if err != nil {
			log.Fatal(err)
		}

		log.Println(string(output))
		os.Stdout.Write(output)
		println()
	}
}
