package handlers

import (
	"encoding/json"
	"net/http"
	"proyectobd2/src/basedata"
)

// GetCanciones maneja la ruta /api/canciones y devuelve todas las canciones
func GetCanciones(w http.ResponseWriter, r *http.Request) {
	canciones, err := basedata.GetAllCanciones()
	if err != nil {
		http.Error(w, "Error al obtener canciones: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(canciones)
}
