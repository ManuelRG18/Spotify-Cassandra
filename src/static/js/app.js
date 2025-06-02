// Registrar escucha de una canción
function registrarEscucha(usuario_id, cancion_id) {
  const fecha = new Date().toISOString().slice(0, 10); // YYYY-MM-DD
  fetch("/api/escuchar", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ usuario_id, cancion_id, fecha_escucha: fecha }),
  })
    .then(res => res.json())
    .then(data => {
      if (data.error) {
        console.error("❌ Error al registrar escucha:", data.error);
      } else {
        console.log("✅ Escucha registrada");
      }
    })
    .catch(() => console.error("❌ Error al registrar escucha (fetch)"));
}
// Cargar canciones en explorar.html
document.addEventListener("DOMContentLoaded", () => {
  const contenedor = document.getElementById("canciones-container");
  if (contenedor) {
    fetch("/api/canciones")
      .then(res => res.json())
      .then(canciones => {
        if (canciones.length === 0) {
          contenedor.innerHTML = "<p>No hay canciones disponibles.</p>";
          return;
        }

        const usuario_id = localStorage.getItem("usuario_id");
        const lista = document.createElement("ul");
        canciones.forEach(c => {
          const item = document.createElement("li");
          item.innerHTML = `<strong>${c.titulo}</strong> – ${c.artista} <em>(${c.genero})</em>`;
          // Registrar escucha al hacer click en la canción
          item.style.cursor = "pointer";
          item.addEventListener("click", function() {
            if (usuario_id && c.id) {
              registrarEscucha(usuario_id, c.id);
            }
          });
          lista.appendChild(item);
        });
        contenedor.appendChild(lista);
      })
      .catch(error => {
        console.error("Error cargando canciones:", error);
        contenedor.innerHTML = "<p>Error al cargar canciones.</p>";
      });
  }

  // Login handler
  const formLogin = document.getElementById("form-login");
  if (formLogin) {
    formLogin.addEventListener("submit", function(e) {
      e.preventDefault();
      const email = document.getElementById("login-usuario").value;
      const password = document.getElementById("login-password").value;

      fetch("/api/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email, password }),
      })
        .then(res => res.json())
        .then(data => {
          if (data.error) {
            alert("❌ " + data.error);
          } else {
            alert("✅ Bienvenido, " + data.nombre);
            localStorage.setItem("usuario_id", data.usuario_id);
            window.location.href = "/dashboard";
          }
        })
        .catch(() => alert("❌ Error al iniciar sesión."));
    });
  }
});

function registrar() {
  const nombre = document.getElementById("registro-nombre").value;
  const email = document.getElementById("registro-usuario").value;
  const password = document.getElementById("registro-password").value;
  const confirm = document.getElementById("confirm-password").value;

  if (password !== confirm) {
    alert("❌ Las contraseñas no coinciden.");
    return;
  }

  fetch("/api/usuarios", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ nombre, email, password }),
  })
    .then(res => res.json())
    .then(data => {
      if (data.error) {
        alert("❌ " + data.error);
      } else {
        alert("✅ Registro exitoso.");
        localStorage.setItem("usuario_id", data.usuario_id);
        window.location.href = "/";
      }
    })
    .catch(() => alert("❌ Error al registrar usuario."));
}
