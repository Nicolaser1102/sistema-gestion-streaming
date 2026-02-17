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

  // 1) cargar info del contenido (público)
  const content = await apiGet(`/contents/${encodeURIComponent(id)}`);
  document.getElementById("info").innerHTML = `
    <h3>${content.title}</h3>
    <p><b>Género:</b> ${content.genre} | <b>Año:</b> ${content.year} | <b>Tipo:</b> ${content.type}</p>
    <p>${content.synopsis || ""}</p>
  `;

  // 2) asignar video demo (si no tienes coverURL/videoURL real)
  // Si luego agregas videoURL al modelo, aquí solo cambias a content.videoURL
  const video = document.getElementById("video");

  const BACKEND_BASE = "http://localhost:8080";
video.src = `${BACKEND_BASE}${content.videoURL}`;
  //deo.src = "https://www.w3schools.com/html/mov_bbb.mp4";

  // 3) obtener progreso
  try {
    const prog = await apiGet(`/playback/${encodeURIComponent(id)}`, true);

    // si existe seconds, continuar
    if (prog && prog.seconds && prog.seconds > 0) {
      showMsg(`Continuando desde ${prog.seconds}s (${prog.percent || 0}%).`);
      // esperar metadata para poder setear currentTime
      video.addEventListener("loadedmetadata", () => {
        video.currentTime = prog.seconds;
      }, { once: true });
    } else {
      showMsg("Iniciando desde el inicio.");
    }
  } catch (e) {
    // si es 403 es suscripción no activa
    console.error(e);
    const msg = String(e);
    if (msg.includes("403")) {
      showMsg("⚠️ Suscripción no activa. No puedes reproducir.");
      video.controls = false;
    } else {
      showMsg("No se pudo cargar el progreso.");
    }
  }

  // 4) Guardar progreso cada 5s mientras reproduce
  let lastSavedSecond = 0;

  async function saveProgress() {
    if (!video.duration || isNaN(video.duration)) return;

    const seconds = Math.floor(video.currentTime);
    if (seconds === lastSavedSecond) return;

    const percent = Math.min(100, Math.max(0, (video.currentTime / video.duration) * 100));

    try {
      await apiPut(`/playback/${encodeURIComponent(id)}/progress`, {
        seconds,
        percent
      }, true);
      lastSavedSecond = seconds;
      // no spamear mensaje cada vez
    } catch (e) {
      console.error(e);
    }
  }

  // cada 5 segundos
  setInterval(() => {
    if (!video.paused && !video.ended) saveProgress();
  }, 5000);

  // guardar al pausar
  video.addEventListener("pause", saveProgress);

  // guardar al terminar (y marcar completado por percent>=90)
  video.addEventListener("ended", async () => {
    try {
      await apiPut(`/playback/${encodeURIComponent(id)}/progress`, {
        seconds: Math.floor(video.duration),
        percent: 100
      }, true);
      showMsg("✅ Reproducción finalizada. Progreso guardado (completado).");
    } catch (e) {
      console.error(e);
      showMsg("Finalizó, pero no se pudo guardar.");
    }
  });
}

loadPlayer().catch(err => {
  console.error(err);
  showMsg("Error cargando el reproductor.");
});

// Helper PUT (si no lo tienes en api.js)
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
