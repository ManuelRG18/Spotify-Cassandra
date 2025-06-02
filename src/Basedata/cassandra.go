// filepath: d:\programacion\proyecto bd2\src\basedata\cassandra.go
package basedata

import (
	"fmt"
	"log"
	"time"

	"github.com/gocql/gocql"
)

var Session *gocql.Session

// InitCassandra inicializa la conexión y crea el keyspace y las tablas necesarias
func InitCassandra() {
	cluster := gocql.NewCluster("127.0.0.1") // Cambia la IP si es necesario
	cluster.Keyspace = "system"
	cluster.Consistency = gocql.Quorum

	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatal("No se pudo conectar a Cassandra:", err)
	}
	defer session.Close()

	// Crear keyspace
	err = session.Query(`CREATE KEYSPACE IF NOT EXISTS spotify WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1}`).Exec()
	if err != nil {
		log.Fatal("No se pudo crear el keyspace:", err)
	}

	// Conectarse al keyspace
	cluster.Keyspace = "spotify"
	Session, err = cluster.CreateSession()
	if err != nil {
		log.Fatal("No se pudo conectar al keyspace spotify:", err)
	}

	// Crear tabla usuarios
	err = Session.Query(`CREATE TABLE IF NOT EXISTS usuarios (
	id UUID PRIMARY KEY,
	nombre TEXT,
	email TEXT,
	password TEXT
	)`).Exec()
	if err != nil {
		log.Fatal("No se pudo crear la tabla usuarios:", err)
	}

	// Crear tabla musica
	err = Session.Query(`CREATE TABLE IF NOT EXISTS musica (
		id UUID PRIMARY KEY,
		titulo TEXT,
		artista TEXT,
		album TEXT,
		anio INT,
		genero TEXT
	)`).Exec()
	if err != nil {
		log.Fatal("No se pudo crear la tabla musica:", err)
	}
	// Crear tabla escuchas
	err = Session.Query(`CREATE TABLE IF NOT EXISTS escuchas (
	usuario_id UUID,
	cancion_id UUID,
	fecha_escucha DATE,
	PRIMARY KEY (usuario_id, fecha_escucha, cancion_id)
) WITH CLUSTERING ORDER BY (fecha_escucha DESC);`).Exec()
	if err != nil {
		log.Fatal("No se pudo crear la tabla escuchas:", err)
	}

	// Crear tabla OLAP escuchas por género y mes
	err = Session.Query(`CREATE TABLE IF NOT EXISTS escuchas_por_genero_mes (
	genero TEXT,
	anio INT,
	mes INT,
	total_escuchas COUNTER,
	PRIMARY KEY ((genero), anio, mes)
)`).Exec()
	if err != nil {
		log.Fatal("No se pudo crear la tabla escuchas_por_genero_mes:", err)
	}
	fmt.Println("Keyspace y tablas creados correctamente en Cassandra")
}

// InsertUsuario inserta un nuevo usuario en la base de datos
func InsertUsuario(nombre, email, password string) (gocql.UUID, error) {
	id := gocql.TimeUUID()

	query := Session.Query(`INSERT INTO usuarios (id, nombre, email, password) VALUES (?, ?, ?, ?)`,
		id, nombre, email, password)

	if err := query.Exec(); err != nil {
		return gocql.UUID{}, fmt.Errorf("error al insertar usuario: %v", err)
	}

	fmt.Printf("Usuario insertado con ID: %v\n", id)
	return id, nil
}

// InsertCancion inserta una nueva canción en la base de datos
func InsertCancion(titulo, artista, album, genero string, anio int) error {
	id := gocql.TimeUUID()
	query := Session.Query(`INSERT INTO musica (id, titulo, artista, album, anio, genero) VALUES (?, ?, ?, ?, ?, ?)`,
		id, titulo, artista, album, anio, genero)
	if err := query.Exec(); err != nil {
		return fmt.Errorf("error al insertar canción: %v", err)
	}
	fmt.Printf("Canción insertada con ID: %v\n", id)
	return nil
}
func SeedMusicData() error {
	canciones := []struct {
		titulo  string
		artista string
		album   string
		anio    int
		genero  string
	}{
		{"Bohemian Rhapsody", "Queen", "A Night at the Opera", 1975, "Rock"},
		{"Imagine", "John Lennon", "Imagine", 1971, "Pop"},
		{"Hotel California", "Eagles", "Hotel California", 1976, "Rock"},
		{"Like a Prayer", "Madonna", "Like a Prayer", 1989, "Pop"},
		{"Smells Like Teen Spirit", "Nirvana", "Nevermind", 1991, "Grunge"},
		{"Hey Jude", "The Beatles", "Hey Jude", 1968, "Rock"},
		{"Billie Jean", "Michael Jackson", "Thriller", 1982, "Pop"},
	}

	for _, c := range canciones {
		// Elimina la canción si ya existe (por título y artista)
		// Esto es para evitar duplicados al reiniciar
		iter := Session.Query("SELECT id FROM musica WHERE titulo = ? AND artista = ? ALLOW FILTERING", c.titulo, c.artista).Iter()
		var idExistente gocql.UUID
		existe := false
		for iter.Scan(&idExistente) {
			existe = true
		}
		iter.Close()
		if existe {
			// Si existe, elimina la canción
			Session.Query("DELETE FROM musica WHERE id = ?", idExistente).Exec()
		}
		// Inserta la canción
		err := InsertCancion(c.titulo, c.artista, c.album, c.genero, c.anio)
		if err != nil {
			return fmt.Errorf("error al insertar canción %s: %v", c.titulo, err)
		}
	}
	return nil
}

// GetAllCanciones retorna todas las canciones de la base de datos
func GetAllCanciones() ([]map[string]interface{}, error) {
	var canciones []map[string]interface{}

	iter := Session.Query("SELECT id, titulo, artista, genero FROM musica").Iter()
	var id gocql.UUID
	var titulo, artista, genero string

	for iter.Scan(&id, &titulo, &artista, &genero) {
		cancion := map[string]interface{}{
			"id":      id.String(),
			"titulo":  titulo,
			"artista": artista,
			"genero":  genero,
		}
		canciones = append(canciones, cancion)
	}

	if err := iter.Close(); err != nil {
		return nil, fmt.Errorf("error al obtener canciones: %v", err)
	}

	return canciones, nil
}

// RegistrarEscucha inserta un registro en la tabla escuchas y actualiza la tabla OLAP
func RegistrarEscucha(usuarioID, cancionID gocql.UUID, fecha string) error {
	// Insertar en escuchas
	query := Session.Query(`INSERT INTO escuchas (usuario_id, cancion_id, fecha_escucha) VALUES (?, ?, ?)`,
		usuarioID, cancionID, fecha)
	if err := query.Exec(); err != nil {
		return fmt.Errorf("error al registrar escucha: %v", err)
	}

	// Obtener género y año/mes de la canción
	var genero string
	var anio int
	var mes int
	// Obtener datos de la canción
	err := Session.Query("SELECT genero, anio FROM musica WHERE id = ?", cancionID).Scan(&genero, &anio)
	if err != nil {
		return fmt.Errorf("error al obtener datos de la canción: %v", err)
	}
	// Parsear mes desde la fecha (YYYY-MM-DD)
	t, err := time.Parse("2006-01-02", fecha)
	if err != nil {
		return fmt.Errorf("error al parsear fecha: %v", err)
	}
	mes = int(t.Month())

	// Actualizar contador OLAP
	err = Session.Query(`UPDATE escuchas_por_genero_mes SET total_escuchas = total_escuchas + 1 WHERE genero = ? AND anio = ? AND mes = ?`,
		genero, anio, mes).Exec()
	if err != nil {
		return fmt.Errorf("error al actualizar OLAP: %v", err)
	}

	fmt.Printf("Escucha registrada y OLAP actualizado: usuario %v, canción %v, fecha %s\n", usuarioID, cancionID, fecha)
	return nil
}

// GetTopCancionesPorGenero retorna las canciones más escuchadas por género
func GetTopCancionesPorGenero(genero string, limite int) ([]map[string]interface{}, error) {
	// Cassandra no soporta JOINs ni agregaciones complejas, así que lo haremos en dos pasos
	// 1. Obtener todas las canciones del género
	canciones, err := GetAllCancionesPorGenero(genero)
	if err != nil {
		return nil, err
	}
	// 2. Contar escuchas por cancion_id
	type cancionContada struct {
		id      string
		titulo  string
		artista string
		total   int
	}
	var resultados []cancionContada
	for _, c := range canciones {
		var total int
		err := Session.Query(`SELECT COUNT(*) FROM escuchas WHERE cancion_id = ? ALLOW FILTERING`, c["id"]).Scan(&total)
		if err == nil && total > 0 {
			resultados = append(resultados, cancionContada{
				id:      c["id"].(string),
				titulo:  c["titulo"].(string),
				artista: c["artista"].(string),
				total:   total,
			})
		}
	}
	// Ordenar por total descendentemente
	for i := 0; i < len(resultados)-1; i++ {
		for j := i + 1; j < len(resultados); j++ {
			if resultados[j].total > resultados[i].total {
				resultados[i], resultados[j] = resultados[j], resultados[i]
			}
		}
	}
	// Limitar resultados
	top := []map[string]interface{}{}
	for i, c := range resultados {
		if i >= limite {
			break
		}
		top = append(top, map[string]interface{}{
			"id":      c.id,
			"titulo":  c.titulo,
			"artista": c.artista,
			"total":   c.total,
		})
	}
	return top, nil
}

// GetAllCancionesPorGenero retorna todas las canciones de un género
func GetAllCancionesPorGenero(genero string) ([]map[string]interface{}, error) {
	var canciones []map[string]interface{}
	iter := Session.Query("SELECT id, titulo, artista, genero FROM musica WHERE genero = ? ALLOW FILTERING", genero).Iter()
	var id, titulo, artista, generoStr string
	for iter.Scan(&id, &titulo, &artista, &generoStr) {
		canciones = append(canciones, map[string]interface{}{
			"id":      id,
			"titulo":  titulo,
			"artista": artista,
			"genero":  generoStr,
		})
	}
	if err := iter.Close(); err != nil {
		return nil, fmt.Errorf("error al obtener canciones por género: %v", err)
	}
	return canciones, nil
}
