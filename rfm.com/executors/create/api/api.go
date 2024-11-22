package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"rfm.com/common"
	"rfm.com/executors/create/model"
)

// CreateFile
// @Summary Create a file
// @Description Create a file int a given path
// @ID create-file
// @Accept json
// @Produce json
// @Param CreateFileRequest body model.CreateFileRequest true "The information about the file to be created"
// @Success 200 {string} string "ok"
// @Failure 400 {object} model.CreateException "Can not create the file"
// @Failure 404 {object} model.CreateException "Can not create the file"
// @Router /create/file [post]
func CreateFile(w http.ResponseWriter, r *http.Request) {
	var request model.CreateFileRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		_ = returnErrorJsonObject(w, throwNewListException("GetRequestBodyException", "Ocorreu um erro ao recuperar o body a requisição"))
	}
	path := filepath.Join(common.BasePath, request.Path)
	err = createDirectory(path)
	if err != nil {
		_ = returnErrorJsonObject(w, "Ocorreu um erro ao criar o path")
		return
	}
	filePath := filepath.Join(path, request.Name)
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
// @Param CreateFileRequest body model.CreateFileRequest true "The information about the file to be created"
// @Success 200 {string} string "ok"
// @Failure 400 {object} model.CreateException "Can not create the directory"
// @Failure 404 {object} model.CreateException "Can not create the directory"
// @Router /create/directory [post]
func CreateDirectory(w http.ResponseWriter, r *http.Request) {
	var request model.CreateDirRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		_ = returnErrorJsonObject(w, throwNewListException("GetRequestBodyException", "Ocorreu um erro ao recuperar o body a requisição"))
	}

	request.Path = filepath.Join(common.BasePath, request.Path)

	err = createDirectory(request.Path)
	if err != nil {
		_ = returnErrorJsonObject(w, "Ocorreu um erro ao criar o path")
	}
}

func createDirectory(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

func returnErrorJsonObject(w http.ResponseWriter, data any) error {
	w.WriteHeader(http.StatusBadRequest)
	return json.NewEncoder(w).Encode(data)
}

func throwNewListException(exception, message string) model.CreateException {
	return model.CreateException{Exception: exception, Details: message}
}
