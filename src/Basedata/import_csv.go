package basedata

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/gocql/gocql"
)

// Importa usuarios y retorna un mapa de id numérico a UUID
func ImportUsuariosCSV(path string) (map[string]gocql.UUID, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	r := csv.NewReader(file)
	r.Read() // saltar cabecera
	m := make(map[string]gocql.UUID)
	for {
		rec, err := r.Read()
		if err != nil {
			break
		}
		id := rec[0]
		nombre := rec[1]
		ciudad := rec[2]
		email := rec[3]
		password := rec[4]
		uuid, _ := InsertUsuario(nombre, ciudad, email, password)
		m[id] = uuid
	}
	return m, nil
}

// Importa canciones y retorna un mapa de id numérico a UUID
func ImportCancionesCSV(path string) (map[string]gocql.UUID, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	r := csv.NewReader(file)
	r.Read() // saltar cabecera
	m := make(map[string]gocql.UUID)
	for {
		rec, err := r.Read()
		if err != nil {
			break
		}
		id := rec[0]
		titulo := rec[1]
		artista := rec[2]
		album := rec[3]
		genero := rec[4]
		anio, _ := strconv.Atoi(rec[5])
		uuid, _ := InsertCancion(titulo, artista, album, genero, anio)
		m[id] = uuid
	}
	return m, nil
}

// Importa escuchas usando los mapas de uuid
func ImportEscuchasCSV(path string, userMap, songMap map[string]gocql.UUID) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	r := csv.NewReader(file)
	r.Read() // saltar cabecera
	for {
		rec, err := r.Read()
		if err != nil {
			break
		}
		usuarioID := rec[0]
		cancionID := rec[1]
		fecha := rec[2]
		uuidUsuario, okU := userMap[usuarioID]
		uuidCancion, okC := songMap[cancionID]
		if okU && okC {
			RegistrarEscucha(uuidUsuario, uuidCancion, fecha)
		}
	}
	return nil
}

// Importación completa
func ImportAllCSVs() error {
	userMap, err := ImportUsuariosCSV("src/csv/usuarios.csv")
	if err != nil {
		return fmt.Errorf("error usuarios: %v", err)
	}
	songMap, err := ImportCancionesCSV("src/csv/canciones.csv")
	if err != nil {
		return fmt.Errorf("error canciones: %v", err)
	}
	err = ImportEscuchasCSV("src/csv/escuchas.csv", userMap, songMap)
	if err != nil {
		return fmt.Errorf("error escuchas: %v", err)
	}
	return nil
}
