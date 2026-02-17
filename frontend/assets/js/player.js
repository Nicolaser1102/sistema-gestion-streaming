function getContentIdFromUrl() {
  const params = new URLSearchParams(window.location.search);
  return params.get("id");
}

function requireLogin() {
  const token = localStorage.getItem("token");
  if (!token) {
    window.location.href = "login.html";
    return false;
  }
  return true;
}

function showMsg(text) {
  document.getElementById("msg").textContent = text || "";
}

async function loadPlayer() {
  if (!requireLogin()) return;

  const id = getContentIdFromUrl();
  if (!id) {
    document.getElementById("info").innerHTML = "<p>Falta el parámetro id.</p>";
    return;
  }

  const video = document.getElementById("video");

  // DEBUG: ver eventos
  video.addEventListener("stalled", () => console.warn("VIDEO stalled"));
  video.addEventListener("waiting", () => console.warn("VIDEO waiting"));
  video.addEventListener("pause", () => console.warn("VIDEO paused at", video.currentTime));
  video.addEventListener("error", () => console.error("VIDEO error", video.error));

  // 1) Cargar contenido
  const content = await apiGet(`/contents/${encodeURIComponent(id)}`);

  document.getElementById("info").innerHTML = `
    <h3>${content.title}</h3>
    <p><b>Género:</b> ${content.genre} | <b>Año:</b> ${content.year} | <b>Tipo:</b> ${content.type}</p>
    <p>${content.synopsis || ""}</p>
  `;

  // 2) Asignar fuente de video (SOLO UNA VEZ)
  const BACKEND_BASE = "http://localhost:8080"; // fijo para evitar errores
  if (!content.videoURL) {
    showMsg("⚠️ Este contenido no tiene videoURL configurado.");
    return;
  }

  video.src = `${BACKEND_BASE}${content.videoURL}`;
  video.preload = "auto";

  // 3) Obtener progreso backend + local (antes de loadedmetadata)
  let backendSeconds = 0;
  let backendPercent = 0;

  try {
    const prog = await apiGet(`/playback/${encodeURIComponent(id)}`, true);
    backendSeconds = prog?.seconds || 0;
    backendPercent = prog?.percent || 0;
  } catch (e) {
    console.error("Get progress failed:", e);
    const msg = String(e);
    if (msg.includes("403")) {
      showMsg("⚠️ Suscripción no activa. No puedes reproducir.");
      video.controls = false;
      return;
    }
  }

  const localSeconds = parseInt(localStorage.getItem(`player_pos_${id}`) || "0", 10) || 0;
  const startSeconds = Math.max(backendSeconds, localSeconds);

  console.log("ID:", id);
  console.log("backendSeconds:", backendSeconds, "localSeconds:", localSeconds, "startSeconds:", startSeconds);

  // 4) Esperar metadata y recién setear currentTime
  video.addEventListener("loadedmetadata", () => {
    if (startSeconds > 0 && startSeconds < video.duration) {
      video.currentTime = startSeconds;
      showMsg(`Continuando desde ${startSeconds}s (${Math.floor((startSeconds / video.duration) * 100)}%).`);
      console.log("✅ currentTime set to:", startSeconds);
    } else {
      showMsg("Iniciando desde el inicio.");
      console.log("ℹ️ startSeconds inválido o 0:", startSeconds, "duration:", video.duration);
    }
  }, { once: true });

  // 5) Guardado de progreso (backend + local)
  let lastSavedSecond = 0;
  let saving = false;

  async function saveProgress(force = false) {
    if (!video.duration || isNaN(video.duration)) return;

    const seconds = Math.floor(video.currentTime);
    if (!force && (seconds - lastSavedSecond < 15)) return;
    if (saving) return;

    // guardar local inmediato
    localStorage.setItem(`player_pos_${id}`, String(seconds));

    const percent = Math.min(100, Math.max(0, (video.currentTime / video.duration) * 100));

    saving = true;
    try {
      await apiPut(`/playback/${encodeURIComponent(id)}/progress`, { seconds, percent }, true);
      lastSavedSecond = seconds;
    } catch (e) {
      console.error("Save progress failed:", e);
    } finally {
      saving = false;
    }
  }

  // Guardar cada 15s mientras reproduce
  setInterval(() => {
    if (!video.paused && !video.ended) saveProgress(false);
  }, 15000);

  // Guardar al pausar y evitar reset visual a 0
  video.addEventListener("pause", () => {
    const t = video.currentTime;
    localStorage.setItem(`player_pos_${id}`, String(Math.floor(t)));
    saveProgress(true);

    setTimeout(() => {
      // si UI se fue a 0, lo devolvemos
      if (t > 1 && video.currentTime < 1) {
        video.currentTime = t;
      }
    }, 150);
  });

  // Guardar al terminar
  video.addEventListener("ended", () => {
    localStorage.setItem(`player_pos_${id}`, "0");
    saveProgress(true);
    showMsg("✅ Reproducción finalizada. Progreso guardado.");
  });
}

loadPlayer().catch(err => {
  console.error(err);
  showMsg("Error cargando el reproductor.");
});

// Helper PUT si no existe en api.js
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

  return res.json().catch(() => ({}));
}
