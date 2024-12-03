package api

import (
	"backend/services"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func PostUsersHandler(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		IdSala      string `json:"idsala"`
		TokenSesion string `json:"tokensession"`
		Nickname    string `json:"nickname"`
	}

	// Decodificar la solicitud
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		log.Printf("Error al decodificar la solicitud: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Obtener la instancia del singleton
	secMod := services.GetSecModServidorChat()
	log.Printf("Singleton de servidor de chat obtenido: %v", secMod)

	// Validar IdSala
	idSala, err := uuid.Parse(requestData.IdSala)
	if err != nil {
		log.Printf("Error al parsear IdSala: %v", err)
		http.Error(w, fmt.Sprintf("IdSala inválido: %v", err), http.StatusBadRequest)
		return
	}
	log.Printf("IdSala validado correctamente: %v", idSala)

	// Obtener usuarios activos
	usuarios, err := secMod.GestionUsuarios.ObtenerUsuariosActivos(idSala)
	if err != nil {
		log.Printf("Error al obtener usuarios activos: %v", err)
		response := map[string]interface{}{
			"status":  "nok",
			"message": err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}
	log.Printf("Usuarios activos obtenidos para la sala %v: %v", idSala, usuarios)

	// Preparar la respuesta
	response := map[string]interface{}{
		"status":   "ok",
		"usuarios": usuarios,
	}
	log.Printf("Respuesta preparada: %v", response)

	// Enviar la respuesta
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
	log.Println("Respuesta enviada correctamente.")
}
