// Cargar canciones en explorar.html
document.addEventListener("DOMContentLoaded", () => {
  const contenedor = document.getElementById("canciones-container");
  if (!contenedor) return; // solo en páginas que lo tengan

  fetch("/api/canciones")
    .then(res => res.json())
    .then(canciones => {
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
      contenedor.innerHTML = "<p>Error al cargar canciones.</p>";
    });

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
