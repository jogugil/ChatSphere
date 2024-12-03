package api

import (
	"backend/services"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func PostUsersHandler(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		IdSala      string `json:"idsala"`
		TokenSesion string `json:"tokensession"`
		Nickname    string `json:"nickname"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Obtener la instancia del singleton
	secMod := services.NewSecModServidorChat()
	idSala, err := uuid.Parse(requestData.IdSala)
	usuarios, err := secMod.GestionUsuarios.ObtenerUsuariosActivos(idSala)
	if err != nil {
		response := map[string]interface{}{
			"status":  "nok",
			"message": err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := map[string]interface{}{
		"status":   "ok",
		"usuarios": usuarios,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
