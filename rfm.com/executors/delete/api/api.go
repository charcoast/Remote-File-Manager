package api

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"

	"rfm.com/common"
	"rfm.com/executors/delete/model"
)

// DeleteFile
// @Summary Delete a file
// @Description Delete a file int a given path
// @ID create-file
// @Accept json
// @Produce json
// @Param DeleteFileRequest body model.DeleteFileRequest true "The information about the file to be created"
// @Success 200 {string} string "ok"
// @Failure 400 {object} model.DeleteException "Can not create the file"
// @Failure 404 {object} model.DeleteException "Can not create the file"
// @Router /create/file [post]
func DeleteFile(w http.ResponseWriter, r *http.Request) {
	var request model.DeleteFileRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		_ = returnErrorJsonObject(w, throwNewDeleteException("GetRequestBodyException", "Ocorreu um erro ao recuperar o body a requisição"))
	}
	path := filepath.Join(common.BasePath, request.Path)
	filePath := filepath.Join(path, request.Name)
	err = os.Remove(filePath)
	if err != nil {
		_ = returnErrorJsonObject(w, throwNewDeleteException("DeleteFileException", "Ocorreu um erro ao remover o arquivo"))
	}
}

// DeleteDirectory
// @Summary Delete a file
// @Description Delete a file int a given path
// @ID create-directory
// @Accept json
// @Produce json
// @Param DeleteFileRequest body model.DeleteFileRequest true "The information about the file to be created"
// @Success 200 {string} string "ok"
// @Failure 400 {object} model.DeleteException "Can not create the directory"
// @Failure 404 {object} model.DeleteException "Can not create the directory"
// @Router /create/directory [post]
func DeleteDirectory(w http.ResponseWriter, r *http.Request) {
	var request model.DeleteDirRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		_ = returnErrorJsonObject(w, throwNewDeleteException("GetRequestBodyException", "Ocorreu um erro ao recuperar o body a requisição"))
	}

	request.Path = filepath.Join(common.BasePath, request.Path)

	err = os.RemoveAll(request.Path)
	if err != nil {
		_ = returnErrorJsonObject(w, "Ocorreu um erro ao criar o path")
	}
}

func returnErrorJsonObject(w http.ResponseWriter, data any) error {
	w.WriteHeader(http.StatusBadRequest)
	return json.NewEncoder(w).Encode(data)
}

func throwNewDeleteException(exception, message string) model.DeleteException {
	return model.DeleteException{Exception: exception, Details: message}
}
