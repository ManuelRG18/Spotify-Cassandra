package basedata

import (
	"fmt"
	"log"

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
		email TEXT
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
		anio INT
	)`).Exec()
	if err != nil {
		log.Fatal("No se pudo crear la tabla musica:", err)
	}

	fmt.Println("Keyspace y tablas creados correctamente en Cassandra")
}

// InsertUsuario inserta un nuevo usuario en la base de datos
func InsertUsuario(nombre, email string) error {
	id := gocql.TimeUUID() // Genera un UUID único

	query := Session.Query(`INSERT INTO usuarios (id, nombre, email) VALUES (?, ?, ?)`,
		id, nombre, email)

	if err := query.Exec(); err != nil {
		return fmt.Errorf("error al insertar usuario: %v", err)
	}

	fmt.Printf("Usuario insertado con ID: %v\n", id)
	return nil
}

// InsertCancion inserta una nueva canción en la base de datos
func InsertCancion(titulo, artista, album string, anio int) error {
	id := gocql.TimeUUID() // Genera un UUID único

	query := Session.Query(`INSERT INTO musica (id, titulo, artista, album, anio) VALUES (?, ?, ?, ?, ?)`,
		id, titulo, artista, album, anio)

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
	}{
		{"Bohemian Rhapsody", "Queen", "A Night at the Opera", 1975},
		{"Imagine", "John Lennon", "Imagine", 1971},
		{"Hotel California", "Eagles", "Hotel California", 1976},
		// Agrega más canciones aquí
	}

	for _, cancion := range canciones {
		err := InsertCancion(cancion.titulo, cancion.artista, cancion.album, cancion.anio)
		if err != nil {
			return fmt.Errorf("error al insertar canción %s: %v", cancion.titulo, err)
		}
	}
	return nil
}
