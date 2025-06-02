function mostrarRegistro() {
  document.getElementById("login-form").style.display = "none";
  document.getElementById("registro-form").style.display = "block";
}

function mostrarLogin() {
  document.getElementById("registro-form").style.display = "none";
  document.getElementById("login-form").style.display = "block";
}

function login() {
  const usuario = document.getElementById("login-usuario").value;
  const password = document.getElementById("login-password").value;

  alert("Intentando iniciar sesión con: " + usuario);
  // Aquí después se conectará con el backend
}

function registrar() {
  const nombre = document.getElementById("registro-nombre").value;
  const usuario = document.getElementById("registro-usuario").value;
  const password = document.getElementById("registro-password").value;

  alert("Registrando usuario: " + nombre);
  // Aquí después se conectará con el backend
}
