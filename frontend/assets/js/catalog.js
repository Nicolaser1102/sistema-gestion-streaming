function buildQuery() {
  const q = document.getElementById("q").value.trim();
  const type = document.getElementById("type").value.trim();
  const genre = document.getElementById("genre").value.trim();
  const year = document.getElementById("year").value.trim();

  const params = new URLSearchParams();
  if (q) params.set("q", q);
  if (type) params.set("type", type);
  if (genre) params.set("genre", genre);
  if (year) params.set("year", year);

  const qs = params.toString();
  return qs ? `?${qs}` : "";
}

function renderList(items) {
  const list = document.getElementById("list");

  if (!items || items.length === 0) {
    list.innerHTML = "<p>No hay resultados.</p>";
    return;
  }

  list.innerHTML = items.map(c => `
    <div class="card">
      <h3>${c.title}</h3>
      <p><b>Tipo:</b> ${c.type} | <b>Género:</b> ${c.genre} | <b>Año:</b> ${c.year}</p>
      <p>${c.synopsis}</p>
      <a href="content_detail.html?id=${encodeURIComponent(c.id)}">Ver detalle</a>
    </div>
  `).join("");
}

async function loadContents() {
  const qs = buildQuery();
  const data = await apiGet(`/contents${qs}`);
  renderList(data);
}

/* ===========================
   CONTINÚA VIENDO
=========================== */

async function loadContinueWatching() {
  const container = document.getElementById("continueList");
  if (!container) return;

  const token = localStorage.getItem("token");
  if (!token) {
    container.innerHTML = "<p>Inicia sesión para ver tu progreso.</p>";
    return;
  }

  try {
    const items = await apiGet("/continue-watching", true);

    if (!items || items.length === 0) {
      container.innerHTML = "<p>No tienes contenido en progreso.</p>";
      return;
    }

    container.innerHTML = items.map(x => `
      <div class="card">
        <h4>${x.content.title}</h4>
        <p>${Math.floor(x.percent)}% • ${x.seconds}s</p>
        <a href="player.html?id=${encodeURIComponent(x.content.id)}">Continuar</a>
      </div>
    `).join("");

  } catch (err) {
    console.error(err);
    container.innerHTML = "<p>Error cargando 'Continúa viendo'.</p>";
  }
}

/* ===========================
   EVENTOS
=========================== */

document.getElementById("btnSearch").addEventListener("click", () => {
  loadContents().catch(err => {
    console.error(err);
    document.getElementById("list").innerHTML = "<p>Error cargando catálogo.</p>";
  });
});

// carga inicial
loadContents().catch(err => {
  console.error(err);
  document.getElementById("list").innerHTML = "<p>Error cargando catálogo.</p>";
});

loadContinueWatching();

/* ===========================
   NAVBAR
=========================== */

function renderNavbar() {
  const nav = document.getElementById("navbar");
  const userStr = localStorage.getItem("user");

  if (!userStr) {
    nav.innerHTML = `
      <a href="login.html">Login</a>
      <a href="register.html">Register</a>
    `;
    return;
  }

  const user = JSON.parse(userStr);

  nav.innerHTML = `
    <span>Bienvenido, ${user.name}</span>
    <a href="catalog.html">Catálogo</a>
    <a href="mylist.html">Mi Lista</a>
    ${user.role === "ADMIN" ? '<a href="admin.html">Panel Admin</a>' : ''}
    <button onclick="logout()">Cerrar sesión</button>
  `;
}

function logout() {
  localStorage.removeItem("token");
  localStorage.removeItem("user");
  window.location.href = "../index.html";
}

renderNavbar();
