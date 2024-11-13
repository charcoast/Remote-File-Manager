package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"rfm.com/executors/create/model"
)

// CreateFile
// @Summary Create a file
// @Description Create a file int a given path
// @ID create-file
// @Accept json
// @Produce json
// @Param CreateRequest body model.CreateRequest true "The information about the file to be created"
// @Success 200 {string} string "ok"
// @Failure 400 {object} model.CreateException "Can not create the file"
// @Failure 404 {object} model.CreateException "Can not create the file"
// @Router /create/file [post]
func CreateFile(w http.ResponseWriter, r *http.Request) {
	var request model.CreateRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		_ = returnJsonObject(&w, throwNewListException("GetRequestBodyException", "Ocorreu um erro ao recuperar o body a requisição"))
	}
	filePath := filepath.Join(request.Path, request.Name)
	err = createDirectory(request.Path)
	if err != nil {
		_ = returnJsonObject(&w, "Ocorreu um erro ao crear crear o path")
	}
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Erro ao criar o arquivo:", err)
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)
	decodedBytes, err := base64.StdEncoding.DecodeString(request.Content)
	if err != nil {
		fmt.Println("Erro ao decodificar:", err)
		return
	}
	_, err = file.Write(decodedBytes)
	if err != nil {
		fmt.Println("Erro ao escrever no arquivo:", err)
		return
	}
}

// CreateDirectory
// @Summary Create a file
// @Description Create a file int a given path
// @ID create-directory
// @Accept json
// @Produce json
// @Param CreateRequest body model.CreateRequest true "The information about the file to be created"
// @Success 200 {string} string "ok"
// @Failure 400 {object} model.CreateException "Can not create the directory"
// @Failure 404 {object} model.CreateException "Can not create the directory"
// @Router /create/directory [post]
func CreateDirectory(w http.ResponseWriter, r *http.Request) {
	var request model.CreateRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		_ = returnJsonObject(&w, throwNewListException("GetRequestBodyException", "Ocorreu um erro ao recuperar o body a requisição"))
	}
	err = createDirectory(filepath.Join(request.Path, request.Name))
	if err != nil {
		_ = returnJsonObject(&w, "Ocorreu um erro ao crear crear o path")
	}
}

func createDirectory(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

func returnJsonObject(w *http.ResponseWriter, data any) error {
	return json.NewEncoder(*w).Encode(data)
}

func throwNewListException(exception, message string) model.CreateException {
	return model.CreateException{Exception: exception, Details: message}
}
