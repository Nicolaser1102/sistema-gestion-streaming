function render(items) {
  const el = document.getElementById("list");

  if (!items || items.length === 0) {
    el.innerHTML = "<p>No tienes contenidos en tu lista.</p>";
    return;
  }

  el.innerHTML = items.map(c => `
    <div class="card">
      <h3>${c.title}</h3>
      <p><b>Tipo:</b> ${c.type} | <b>Género:</b> ${c.genre} | <b>Año:</b> ${c.year}</p>
      <a href="content_detail.html?id=${encodeURIComponent(c.id)}">Ver detalle</a>
      <button data-id="${c.id}" class="btnRemove">Quitar</button>
    </div>
  `).join("");

  document.querySelectorAll(".btnRemove").forEach(btn => {
    btn.addEventListener("click", async () => {
      const id = btn.getAttribute("data-id");
      try {
        await apiDelete(`/my-list/${encodeURIComponent(id)}`, true);
        // recargar
        load().catch(console.error);
      } catch (e) {
        console.error(e);
        alert("No se pudo quitar.");
      }
    });
  });
}

async function load() {
  const token = localStorage.getItem("token");
  if (!token) {
    document.getElementById("list").innerHTML = "<p>Debes iniciar sesión.</p>";
    return;
  }

  const data = await apiGet("/my-list", true);
  render(data);
}

load().catch(err => {
  console.error(err);
  document.getElementById("list").innerHTML = "<p>Error cargando Mi lista.</p>";
});
