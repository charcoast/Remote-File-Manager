package api

import (
	"bufio"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"rfm.com/executors/read/model"
)

// ReadFile
// @Summary Create a file
// @Description Create a file int a given path
// @ID create-file
// @Accept json
// @Produce json
// @Param ReadRequest body model.ReadRequest true "The information about the file to be created"
// @Success 200 {string} string "ok"
// @Failure 400 {object} model.ReadException "Can not create the file"
// @Failure 404 {object} model.ReadException "Can not create the file"
// @Router /create/file [post]
func ReadFile(w http.ResponseWriter, r *http.Request) {
	var request model.ReadRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		_ = returnJsonObject(&w, throwNewListException("GetRequestBodyException", "Ocorreu um erro ao recuperar o body a requisição"))
	}
	filePath := filepath.Join(request.Path, request.Name)
	file, err := os.ReadFile(filePath)

	if err != nil {
		_ = returnJsonObject(&w, "Ocorreu um erro ao ler o path")
	}

	bufio.NewWriter(w).Write(file)
	w.Header().Set("Content-Disposition", "inline; filename="+request.Name)
}

// CreateDirectory
// @Summary Create a file
// @Description Create a file int a given path
// @ID create-directory
// @Accept json
// @Produce json
// @Param ReadRequest body model.ReadRequest true "The information about the file to be created"
// @Success 200 {string} string "ok"
// @Failure 400 {object} model.ReadException "Can not create the directory"
// @Failure 404 {object} model.ReadException "Can not create the directory"
// @Router /create/directory [post]
func CreateDirectory(w http.ResponseWriter, r *http.Request) {
	var request model.ReadRequest
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

func throwNewListException(exception, message string) model.ReadException {
	return model.ReadException{Exception: exception, Details: message}
}
