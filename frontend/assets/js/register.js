document.getElementById("registerForm").addEventListener("submit", async (e) => {
  e.preventDefault(); // evita recargar la página

  const name = document.getElementById("name").value.trim();
  const email = document.getElementById("email").value.trim();
  const password = document.getElementById("password").value.trim();
  const msg = document.getElementById("msg");

  msg.textContent = "";

  // Validación mínima frontend
  if (!name) {
    msg.textContent = "Nombre requerido.";
    return;
  }
  if (!email || !email.includes("@")) {
    msg.textContent = "Email inválido.";
    return;
  }
  if (!password || password.length < 6) {
    msg.textContent = "Password mínimo 6 caracteres.";
    return;
  }

  try {
    await apiPost("/auth/register", { name, email, password }, false);

    msg.textContent = "✅ Registrado correctamente. Volviendo al inicio...";

    // opcional: limpiar campos
    document.getElementById("registerForm").reset();

    setTimeout(() => {
      window.location.href = "login.html";
    }, 1200);

  } catch (err) {
    console.error(err);

    // Mensaje amigable (sin parsear el error del backend)
    if (String(err).includes("409")) {
      msg.textContent = "⚠️ Ese email ya está registrado.";
    } else {
      msg.textContent = "❌ No se pudo registrar. Revisa los datos.";
    }
  }
});
