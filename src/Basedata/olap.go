package basedata

import (
	"fmt"
)

type EscuchasCiudadMes struct {
	Ciudad string `json:"ciudad"`
	Anio   int    `json:"anio"`
	Mes    int    `json:"mes"`
	Total  int    `json:"total"`
}

// Consulta OLAP: escuchas por ciudad y mes
func GetEscuchasPorCiudadMes() ([]EscuchasCiudadMes, error) {
	var resultados []EscuchasCiudadMes
	query := "SELECT ciudad, anio, mes, total_escuchas FROM escuchas_por_ciudad_mes"
	iter := Session.Query(query).Iter()
	var ciudad string
	var anio, mes int
	var total int
	for iter.Scan(&ciudad, &anio, &mes, &total) {
		resultados = append(resultados, EscuchasCiudadMes{
			Ciudad: ciudad,
			Anio:   anio,
			Mes:    mes,
			Total:  total,
		})
	}
	if err := iter.Close(); err != nil {
		return nil, fmt.Errorf("error al obtener datos OLAP por ciudad: %v", err)
	}
	return resultados, nil
}

type EscuchasGeneroMes struct {
	Genero string `json:"genero"`
	Anio   int    `json:"anio"`
	Mes    int    `json:"mes"`
	Total  int    `json:"total"`
}

// Consulta OLAP: escuchas por género y mes
func GetEscuchasPorGeneroMes() ([]EscuchasGeneroMes, error) {
	var resultados []EscuchasGeneroMes
	query := "SELECT genero, anio, mes, total_escuchas FROM escuchas_por_genero_mes"
	iter := Session.Query(query).Iter()
	var genero string
	var anio, mes int
	var total int
	for iter.Scan(&genero, &anio, &mes, &total) {
		resultados = append(resultados, EscuchasGeneroMes{
			Genero: genero,
			Anio:   anio,
			Mes:    mes,
			Total:  total,
		})
	}
	if err := iter.Close(); err != nil {
		return nil, fmt.Errorf("error al obtener datos OLAP: %v", err)
	}
	return resultados, nil
}
