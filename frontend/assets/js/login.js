document.getElementById("btnLogin").addEventListener("click", async () => {
  const email = document.getElementById("email").value.trim();
  const password = document.getElementById("password").value.trim();
  const msg = document.getElementById("msg");

  msg.textContent = "";

  try {
    const data = await apiPost("/auth/login", { email, password }, false);

    // Guardar token y user
    localStorage.setItem("token", data.token);
    localStorage.setItem("user", JSON.stringify(data.user));

    msg.textContent = "Login OK. Redirigiendo...";
    window.location.href = "catalog.html";
  } catch (e) {
    console.error(e);
    msg.textContent = "Login fallido. Revisa credenciales.";
  }
});
