package api

import (
	"encoding/json"
	"net/http"
	"os"
	"rfm.com/executors/list/model"
	"strings"
)

// GetDirectories
// @Summary Get the directories of the path
// @Description Get the directories
// @ID get-directories
// @Accept json
// @Produce json
// @Success 200 {string} []string "ok"
// @Failure 400 {object} model.ListException "Can not get the directories"
// @Failure 404 {object} model.ListException "Can not get the directories"
// @Router /list [get]
func GetDirectories(w http.ResponseWriter, r *http.Request) {
	getDirectories(w, "./")
}

// GetDirectoriesByBody
// @Summary Get the directories of a given path
// @Description Get the directories by the path specified
// @ID get-directories
// @Accept json
// @Produce json
// @Param listRequest body model.ListRequest true "The informations about the directories"
// @Success 200 {string} []string "ok"
// @Failure 400 {object} model.ListException "Can not get the directories"
// @Failure 404 {object} model.ListException "Can not get the directories"
// @Router /list [post]
func GetDirectoriesByBody(w http.ResponseWriter, r *http.Request) {
	var request model.ListRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		_ = returnJsonObject(&w, throwNewListException("GetRequestBodyException", "Ocorreu um erro ao recuperar o body a requisição"))
	}

	if strings.Index(request.Path, "./") > -1 || strings.Index(request.Path, "../") > -1 {
		_ = returnJsonObject(&w, throwNewListException("PathException", "O caminho não pode ser relativo"))
	}

	if !strings.HasPrefix(request.Path, "/") {
		_ = returnJsonObject(&w, throwNewListException("PathException", "O caminho deve começar com /"))
	}

	getDirectories(w, request.Path)
}

func getDirectories(w http.ResponseWriter, path string) {
	path = "/app" + path
	dirs, err := os.ReadDir(path)
	if err != nil {
		_ = returnJsonObject(&w, throwNewListException("GetDirectoriesException", "Ocorreu um erro ao consultar os diretórios"))
	}
	var directories = make([]model.Entry, 0)
	for _, dir := range dirs {
		dirName := strings.TrimSpace(dir.Name())
		if dirName == "" || len(dirName) == 0 {
			continue
		}
		directories = append(directories, model.Entry{Name: dirName, IsDir: dir.IsDir()})
	}
	err = returnJsonObject(&w, directories)
	if err != nil {
		_ = returnJsonObject(&w, throwNewListException("ConversionException", "Erro ao executar a codificação do objeto"))
	}
}

func returnJsonObject(w *http.ResponseWriter, data any) error {
	return json.NewEncoder(*w).Encode(data)
}

func throwNewListException(exception, message string) model.ListException {
	return model.ListException{Exception: exception, Details: message}
}
