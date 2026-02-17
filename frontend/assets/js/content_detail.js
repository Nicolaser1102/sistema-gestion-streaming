function getContentIdFromUrl() {
  const params = new URLSearchParams(window.location.search);
  return params.get("id");
}

function renderDetail(c) {
  const detail = document.getElementById("detail");
  detail.innerHTML = `
    <div class="card">
      <h3>${c.title}</h3>
      <p><b>Tipo:</b> ${c.type}</p>
      <p><b>Género:</b> ${c.genre}</p>
      <p><b>Año:</b> ${c.year}</p>
      <p><b>Duración (min):</b> ${c.durationMin}</p>
      <p><b>Sinopsis:</b> ${c.synopsis}</p>
      <p><b>Activo:</b> ${c.active ? "Sí" : "No"}</p>
    </div>
  `;
}

(async () => {
  const id = getContentIdFromUrl();
  if (!id) {
    document.getElementById("detail").innerHTML = "<p>Falta el parámetro id.</p>";
    return;
  }

  try {
    const content = await apiGet(`/contents/${encodeURIComponent(id)}`);
    renderDetail(content);

      await setupMyListButton(id);
  } catch (e) {
    console.error(e);
    document.getElementById("detail").innerHTML = "<p>No se pudo cargar el detalle.</p>";
  }
})();

async function isInMyList(contentId) {
  try {
    const token = localStorage.getItem("token");
    if (!token) return false;

    const list = await apiGet("/my-list", true);
    return Array.isArray(list) && list.some(c => String(c.id) === String(contentId));
  } catch {
    return false;
  }
}

async function setupMyListButton(contentId) {
  const btn = document.getElementById("btnAdd");
  const msg = document.getElementById("msg");

  if (!btn) return;

  const token = localStorage.getItem("token");
  if (!token) {
    btn.disabled = true;
    msg.textContent = "Inicia sesión para usar Mi lista.";
    return;
  }

  const already = await isInMyList(contentId);
  if (already) {
    btn.textContent = "Ya está en Mi lista ✅";
    btn.disabled = true;
    return;
  }

  btn.addEventListener("click", async () => {
    msg.textContent = "";
    try {
      await apiPost(`/my-list/${encodeURIComponent(contentId)}`, null, true);
      btn.textContent = "Agregado ✅";
      btn.disabled = true;
      msg.textContent = "Contenido agregado a Mi lista.";
    } catch (e) {
      console.error(e);
      msg.textContent = "No se pudo agregar. Revisa sesión/token.";
    }
  });
}
