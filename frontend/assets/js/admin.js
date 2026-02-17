// ✅ 1) Verificar acceso admin
function checkAdminAccess() {
  const userStr = localStorage.getItem("user");
  if (!userStr) {
    window.location.href = "login.html";
    return false;
  }

  const user = JSON.parse(userStr);
  if (user.role !== "ADMIN") {
    document.body.innerHTML = "<h2>Acceso restringido. Solo administradores.</h2>";
    return false;
  }

  return true;
}

if (!checkAdminAccess()) {
  throw new Error("No autorizado");
}

function showMsg(text) {
  document.getElementById("msg").textContent = text || "";
}

// ✅ 2) Cargar contenidos en lista admin
async function loadAdminContents() {
  const el = document.getElementById("contentList");
  el.innerHTML = "<p>Cargando...</p>";

  try {
    const contents = await apiGet("/admin/contents", true);

    if (!contents || contents.length === 0) {
      el.innerHTML = "<p>No hay contenidos.</p>";
      return;
    }

    el.innerHTML = contents.map(c => `
      <div class="card">
        <h3>${c.title}</h3>
        <p>
          <b>ID:</b> ${c.id} |
          <b>Tipo:</b> ${c.type} |
          <b>Género:</b> ${c.genre} |
          <b>Año:</b> ${c.year} |
          <b>Activo:</b> ${c.active ? "Sí" : "No"}
        </p>

        <button class="btnEdit" data-id="${c.id}">Editar</button>
        <button class="btnDelete" data-id="${c.id}">Desactivar</button>
        <button class="btnActivate" data-id="${c.id}">Activar</button>
      </div>
    `).join("");

    // Eventos botones
    document.querySelectorAll(".btnEdit").forEach(btn => {
      btn.addEventListener("click", () => {
        const id = btn.getAttribute("data-id");
        const item = contents.find(x => String(x.id) === String(id));
        if (item) openEditForm(item);
      });
    });

    document.querySelectorAll(".btnDelete").forEach(btn => {
      btn.addEventListener("click", async () => {
        const id = btn.getAttribute("data-id");
        if (!confirm("¿Seguro que quieres desactivar este contenido?")) return;

        try {
          await apiDelete(`/admin/contents/${encodeURIComponent(id)}`, true);
          showMsg("✅ Contenido desactivado.");
          await loadAdminContents();
        } catch (e) {
          console.error(e);
          showMsg("❌ No se pudo desactivar.");
        }
      });
    });

    document.querySelectorAll(".btnActivate").forEach(btn => {
      btn.addEventListener("click", async () => {
        const id = btn.getAttribute("data-id");
        const item = contents.find(x => String(x.id) === String(id));
        if (!item) return;

        try {
          const payload = {
            title: item.title,
            synopsis: item.synopsis,
            genre: item.genre,
            year: item.year,
            durationMin: item.durationMin || 0,
            type: item.type,
            coverURL: item.coverURL || "",
            active: true
          };

          await apiPut(`/admin/contents/${encodeURIComponent(id)}`, payload, true);
          showMsg("✅ Contenido activado.");
          await loadAdminContents();
        } catch (e) {
          console.error(e);
          showMsg("❌ No se pudo activar.");
        }
      });
    });

  } catch (e) {
    console.error(e);
    el.innerHTML = "<p>Error cargando contenidos (¿token admin?).</p>";
  }
}

// ✅ 3) Crear contenido
document.getElementById("createForm").addEventListener("submit", async (e) => {
  e.preventDefault();

  const content = {
    title: document.getElementById("title").value.trim(),
    genre: document.getElementById("genre").value.trim(),
    year: parseInt(document.getElementById("year").value),
    type: document.getElementById("type").value,
    synopsis: document.getElementById("synopsis").value.trim(),
    durationMin: parseInt(document.getElementById("duration").value || "0"),
    coverURL: ""
  };

  try {
    await apiPost("/admin/contents", content, true);
    showMsg("✅ Contenido creado correctamente.");
    document.getElementById("createForm").reset();
    await loadAdminContents();
  } catch (err) {
    console.error(err);
    showMsg("❌ Error al crear contenido.");
  }
});

// ✅ 4) Abrir formulario de edición con datos
function openEditForm(c) {
  const form = document.getElementById("editForm");
  form.style.display = "block";

  document.getElementById("editId").value = c.id;
  document.getElementById("editTitle").value = c.title || "";
  document.getElementById("editGenre").value = c.genre || "";
  document.getElementById("editYear").value = c.year || "";
  document.getElementById("editType").value = c.type || "MOVIE";
  document.getElementById("editSynopsis").value = c.synopsis || "";
  document.getElementById("editDuration").value = c.durationMin || 0;
  document.getElementById("editActive").value = String(!!c.active);

  // bajar/mostrar el form
  form.scrollIntoView({ behavior: "smooth" });
}

// ✅ 5) Guardar edición (PUT)
document.getElementById("editForm").addEventListener("submit", async (e) => {
  e.preventDefault();

  const id = document.getElementById("editId").value;

  const payload = {
    title: document.getElementById("editTitle").value.trim(),
    genre: document.getElementById("editGenre").value.trim(),
    year: parseInt(document.getElementById("editYear").value),
    type: document.getElementById("editType").value,
    synopsis: document.getElementById("editSynopsis").value.trim(),
    durationMin: parseInt(document.getElementById("editDuration").value || "0"),
    coverURL: "",
    active: document.getElementById("editActive").value === "true"
  };

  try {
    await apiPut(`/admin/contents/${encodeURIComponent(id)}`, payload, true);
    showMsg("✅ Contenido actualizado.");
    document.getElementById("editForm").style.display = "none";
    await loadAdminContents();
  } catch (err) {
    console.error(err);
    showMsg("❌ No se pudo actualizar.");
  }
});

// ✅ 6) Cancelar edición
document.getElementById("btnCancelEdit").addEventListener("click", () => {
  document.getElementById("editForm").style.display = "none";
  showMsg("");
});

// ✅ 7) Helper PUT usando fetch (si no existe en api.js)
async function apiPut(path, body, auth = false) {
  const headers = { "Content-Type": "application/json" };
  if (auth) {
    const token = localStorage.getItem("token") || "";
    if (token) headers["Authorization"] = `Bearer ${token}`;
  }

  const res = await fetch(`${API_BASE}${path}`, {
    method: "PUT",
    headers,
    body: JSON.stringify(body)
  });

  if (!res.ok) {
    const text = await res.text();
    throw new Error(`HTTP ${res.status}: ${text}`);
  }

  // tu backend responde {"message": "..."} en update
  return res.json().catch(() => ({}));
}

// ✅ Ejecutar al cargar página
loadAdminContents();
