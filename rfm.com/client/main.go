package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const DiscoveryPort = "9090"

func main() {

	arguments := os.Args     // Pega os argumentos da linha de comando
	if len(arguments) == 1 { // Se não tiver argumentos retorna erro
		fmt.Println("Enter with arguments host.")
		return

	}
	discoveryHost := arguments[1]
	fmt.Println("Digite um comando e pressione ENTER")

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("\nCOMANDO: ")

		text, _ := reader.ReadString('\n')
		body, _ := json.Marshal(Command{Command: text})
		resp, _ := http.Post("http://"+discoveryHost+":"+DiscoveryPort+"/command", "application/json", bytes.NewBuffer(body))

		buf := new(strings.Builder)
		_, _ = io.Copy(buf, resp.Body)

		fmt.Println("RESULTADO: " + buf.String())
	}

}
