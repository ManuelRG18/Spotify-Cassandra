document.addEventListener("DOMContentLoaded", () => {
  fetch("/api/canciones")
    .then(res => res.json())
    .then(canciones => {
      const contenedor = document.getElementById("canciones-container");
      if (canciones.length === 0) {
        contenedor.innerHTML = "<p>No hay canciones disponibles.</p>";
        return;
      }

      const lista = document.createElement("ul");
      canciones.forEach(c => {
        const item = document.createElement("li");
        item.innerHTML = `<strong>${c.titulo}</strong> – ${c.artista} <em>(${c.genero})</em>`;
        lista.appendChild(item);
      });
      contenedor.appendChild(lista);
    })
    .catch(error => {
      console.error("Error cargando canciones:", error);
      document.getElementById("canciones-container").innerHTML = "<p>Error al cargar canciones.</p>";
    });
});
