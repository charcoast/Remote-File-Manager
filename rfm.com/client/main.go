package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

const DiscoveryPort = "9090"

func main() {

	arguments := os.Args     // Pega os argumentos da linha de comando
	if len(arguments) == 1 { // Se não tiver argumentos retorna erro
		fmt.Println("Enter with arguments host.")
		return
	}
	fmt.Print("Digite um comando e pressione ENTER")

	for {
		c, err := net.Dial("tcp", arguments[1]+":"+DiscoveryPort)
		if err != nil {
			fmt.Println(err)
			return
		}

		reader := bufio.NewReader(os.Stdin) // Prepara o buffer de leitura
		fmt.Println("\nMSG: ")
		text, _ := reader.ReadString('\n') // Le um texto do teclado
		fmt.Fprintf(c, text+"\n")          // Envia o texto pela conexão

		message, _ := bufio.NewReader(c).ReadString('\n') // Aguarda resposta do servidor
		fmt.Print("RCV: " + message + "\n\n")
	}

}
