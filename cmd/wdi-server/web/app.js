/* ==========================================================================
   WDI w Go — frontend SPA (vanilla JS, hash routing)
   8 kafelków → 8 graficznych implementacji zadań.
   ========================================================================== */

"use strict";

/* ---------------------------- helpers ---------------------------- */

const $ = (sel, el = document) => el.querySelector(sel);
const $$ = (sel, el = document) => Array.from(el.querySelectorAll(sel));

const esc = (s) =>
  String(s).replace(/[&<>"']/g, (c) => ({
    "&": "&amp;", "<": "&lt;", ">": "&gt;", '"': "&quot;", "'": "&#39;",
  }[c]));

function fmt(x, d = 6) {
  if (typeof x !== "number" || !isFinite(x)) return "∞";
  if (x === 0) return "0";
  const a = Math.abs(x);
  if (a >= 1e6 || a < 1e-4) return x.toExponential(2);
  return parseFloat(x.toFixed(d)).toString();
}

const fmtInt = (n) => n.toLocaleString("pl-PL");

function fmtSec(s) {
  if (!isFinite(s)) return "—";
  if (s < 0.001) return "< 1 ms";
  if (s < 1) return Math.round(s * 1000) + " ms";
  if (s < 100) return s.toFixed(2) + " s";
  return s.toFixed(1) + " s";
}

function digitCount(n) {
  return n === 0 ? 1 : Math.floor(Math.log10(n)) + 1;
}

async function api(method, url, body) {
  const opt = { method, headers: {} };
  if (body !== undefined) {
    opt.headers["Content-Type"] = "application/json";
    opt.body = JSON.stringify(body);
  }
  const r = await fetch(url, opt);
  const j = await r.json().catch(() => ({}));
  if (!r.ok) throw new Error(j.error || "błąd " + r.status);
  return j;
}

/* ------- rejestr dyspozytorów bieżącego widoku -------
   Wszystko, co widok zostawia poza sobą (intery, słuchacze dokumentu,
   pollery, aktywne joby), musi być wyrejestrowane przy zmianie routu.
   Dysponenci są wołani w route() — bez timerów globalnych i bez wycieków. */
let disposers = [];

/* onLeave — rejestruje sprzątanie bieżącego widoku. */
function onLeave(fn) {
  disposers.push(fn);
}

function runDisposers() {
  const list = disposers;
  disposers = [];
  for (const fn of list) {
    try {
      fn();
    } catch (_) {
      /* sprzątanie nie może wywrócić zmiany routu */
    }
  }
}

/* poll — interwał z uchwytem: zwraca stop() i sam się wypina z rejestru.
   Callback dostaje stop(), żeby mógł zakończyć pętlę (job zakończony, błąd). */
function poll(fn, ms) {
  let stopped = false;
  let id = 0;
  const stop = () => {
    if (stopped) return;
    stopped = true;
    clearInterval(id);
  };
  onLeave(stop);
  id = setInterval(() => {
    if (stopped) {
      stop();
      return;
    }
    fn(stop);
  }, ms);
  return stop;
}

/* ---------------------------- meta ----------------------------
   Katalog zadań ma JEDNO źródło: rejestr w Go (GET /api/tasks). Frontend nie
   trzyma własnej kopii tytułów ani opisów — pobiera listę raz na czas życia
   strony (cache w domknięciu modułu) i z niej buduje zarówno kafelki startowe,
   jak i nagłówki podstron. W JS zostaje to, co jest kodem, a nie danymi:
   mapowanie trasa → funkcja renderująca (VIEWS) oraz akcent kolorystyczny. */

// Akcenty — wyłącznie prezentacja; te dane nie istnieją w rejestrze Go.
const ACCENTS = {
  z01: "#4f9cf7", z02: "#f7b94f", z03: "#4fd78f", z04: "#e36ea8",
  z05: "#b478f7", z06: "#4fd7d7", z07: "#f78f4f", z08: "#8f9cf7",
};
const ACCENT_FALLBACK = "#4f9cf7";

const accentOf = (id) => ACCENTS[id] || ACCENT_FALLBACK;

let catalog = null;           // cache odpowiedzi GET /api/tasks
const catalogById = new Map();

// loadCatalog — jeden request na czas życia strony; każdy widok dostaje
// te same tytuły/opisy, co kafelki. Błąd sieci idzie do wołającego (banner).
async function loadCatalog() {
  if (catalog) return catalog;
  const list = (await api("GET", "/api/tasks")).tasks;
  catalog = list;
  for (const t of list) catalogById.set(t.id, t);
  return catalog;
}

// taskTitle — tytuł z katalogu API. Brak wpisu => sam identyfikator trasy:
// świadomie NIE kopiujemy tu tekstu, bo duplikat zawsze by się rozjechał.
function taskTitle(id) {
  const t = catalogById.get(id);
  return t ? t.title : id;
}

// taskHead — nagłówek podstrony: plakietka z id + tytuł z katalogu API.
function taskHead(id) {
  return `<div class="task-head"><span class="badge">${esc(id)}</span><h2>${esc(taskTitle(id))}</h2></div>`;
}

function setAccent(view, id) {
  view.style.setProperty("--accent", accentOf(id));
}

/* Stan pomocniczy widoków mieszka w domknięciu modułu — dawniej był to
   niejawny duet globali na oknie przeglądarki ("__capKeys" i "__sol"), czyli
   dwa anonimowe kanały zapisywalne z konsoli. Zastępuje je JEDNO nazwane API
   deweloperskie — window.wdiDebug — tylko do odczytu; writerem pozostaje UI.
     capKeys  — sufit przyjęty przez ostatnie zlecenie (z02, z05),
     solution — ostatnie wygenerowane rozwiązanie sudoku (z08),
     tasks    — id zadań z katalogu API. */
const devState = { capKeys: 2e6, solution: null };

window.wdiDebug = Object.freeze({
  get capKeys() { return devState.capKeys; },
  get solution() { return devState.solution; },
  tasks: () => (catalog ? catalog.map((t) => t.id) : []),
});

function banner(cls, html) {
  return `<div class="banner ${cls}">${html}</div>`;
}

function panel(html, h4) {
  return `<section class="panel">${h4 ? `<h4>${esc(h4)}</h4>` : ""}${html}</section>`;
}

function jobBox() {
  return `
  <div class="jobbox" hidden>
    <div class="pbar" role="progressbar" aria-label="postęp obliczeń"
         aria-valuemin="0" aria-valuemax="100" aria-valuenow="0"><i></i></div>
    <div class="jobstatus" role="status" aria-live="polite"><span class="what">…</span><span class="pct"></span></div>
    <div class="jobmsg"></div>
    <div class="actions">
      <button class="btn danger small cancel-btn" type="button">Przerwij</button>
    </div>
  </div>`;
}

/* startJob — pasek postępu + odpytywanie statusu joba.
   · Przy zmianie routu (onLeave): zatrzymuje poller, anuluje aktywny job
     (DELETE na serwerze) i blokuje onSettled — widok już nie żyje.
   · Błąd fetch w pętli: komunikat w UI + zatrzymanie pętli (inaczej pasek
     postępu wisi w nieskończoność bez żadnego sygnału). */
function startJob(el, jobId, jobUrl) {
  const box = $(".jobbox", el);
  const acts = $(".actions", box);
  const msg = $(".jobmsg", box);
  const bar = $(".pbar > i", box);
  const pbar = $(".pbar", box);
  const what = $(".what", box);
  const pct = $(".pct", box);
  box.hidden = false;
  acts.hidden = false; // box jest współdzielony między uruchomieniami
  msg.innerHTML = "";
  /* setProgress — JEDNO miejsce ustawiające pasek: szerokość wypełnienia,
     procent przy pasku i aria-valuenow wskaźnika role="progressbar".
     Bez tego czytnik ekranu słyszał zawsze „0”. */
  const setProgress = (frac) => {
    const v = Math.round(100 * Math.min(1, Math.max(0, frac)));
    bar.style.width = v + "%";
    pct.textContent = v + " %";
    pbar.setAttribute("aria-valuenow", String(v));
    return v;
  };
  $(".cancel-btn", box).onclick = async () => {
    try { await fetch(jobUrl, { method: "DELETE" }); } catch (_) {}
  };

  let disposed = false; // wyszliśmy z widoku: nie dotykamy DOM i nie wołamy onSettled
  let settled = false;  // job zakończony: przy odejściu nie ma czego anulować
  let onJobSettled = () => {};

  const stop = poll(async () => {
    if (disposed) { stop(); return; }
    let s;
    try {
      s = await api("GET", jobUrl);
    } catch (e) {
      stop();
      if (disposed) return;
      what.textContent = "błąd połączenia";
      acts.hidden = true;
      msg.innerHTML = banner("err",
        "Utracono połączenie z serwerem — odpytywanie przerwane: " + esc(e.message) +
        ". Odśwież stronę i spróbuj ponownie.");
      return;
    }
    if (disposed) return;
    if (s.status === "done" || s.status === "error" || s.status === "cancelled") {
      stop();
      settled = true;
      // job zakończony: chowamy „Przerwij”, pokazujemy stan końcowy
      acts.hidden = true;
      if (s.status === "done") {
        setProgress(1);
        what.textContent = "zakończono";
      } else {
        what.textContent = s.status === "cancelled" ? "przerwano" : "błąd";
      }
      onJobSettled(s);
    } else if (s.progress) {
      renderProgress(s.progress);
    }
  }, 600);

  onLeave(() => {
    disposed = true;
    stop();
    // odjeżdżamy z routu — aktywny job nie ma odbiorcy: anulujemy go na serwerze
    if (!settled) {
      try {
        fetch(jobUrl, { method: "DELETE", keepalive: true }).catch(() => {});
      } catch (_) {}
    }
  });

  function renderProgress(p) {
    let frac = 0, whatTxt = "";
    if ("processed" in p && p.limit) {
      frac = p.processed / (p.limit + 1);
      whatTxt = `przetworzono ${fmtInt(p.processed)} / ${fmtInt(p.limit + 1)} liczb`;
    } else if ("done" in p && p.total) {
      frac = p.done / p.total;
      whatTxt = "pomiary: " + p.done + " / " + p.total +
        (p.last ? ` · ostatni: ${p.last.structure} n=${fmtInt(p.last.n)} ${fmtSec(p.last.seconds)}` : "");
    } else if ("keysTried" in p) {
      frac = Math.min(1, p.keysTried / (devState.capKeys || 2e6));
      whatTxt = `wypróbowano ${fmtInt(p.keysTried)} kluczy (długość ${p.len})`;
    }
    setProgress(frac);
    what.textContent = whatTxt;
  }
  return {
    set onSettled(fn) { onJobSettled = fn; },
  };
}

/* ---------------------------- router ---------------------------- */

const VIEWS = {};

/* route — zmiana routu. Przed renderem dokłada katalog z API (nagłówek
   podstrony bierze tytuł z jednego źródła). routeSeq pilnuje, żeby wolniejsza
   odpowiedź starszego routu nie nadpisała nowszego widoku. */
let routeSeq = 0;

async function route() {
  const seq = ++routeSeq;
  runDisposers();
  const h = (location.hash || "#/").replace(/^#\//, "");
  const view = $("#view");
  view.innerHTML = "";
  const ctx = $("#ctx-title");
  if (!h || VIEWS[h] === undefined) {
    ctx.textContent = "zadania w Go";
    renderHome(view);
    return;
  }
  ctx.textContent = "zadanie " + h;
  setAccent(view, h);
  try {
    await loadCatalog();
  } catch (e) {
    if (seq !== routeSeq) return;
    view.innerHTML = banner("err", "Nie udało się pobrać spisu zadań: " + esc(e.message));
    return;
  }
  if (seq !== routeSeq) return; // nowszy rout wygrał — nie nadpisujemy
  VIEWS[h](view);
}

window.addEventListener("hashchange", route);
window.addEventListener("DOMContentLoaded", route);

/* ---------------------------- home ---------------------------- */

async function renderHome(view) {
  view.innerHTML = `
    <div class="muted" style="margin-bottom:18px">
      Wybierz jedno z ośmiu zadań — każde ma własną, interaktywną implementację po stronie Go.
    </div>
    <div class="tiles" id="tiles"><div class="muted">Ładowanie…</div></div>`;
  let tasks;
  try {
    tasks = await loadCatalog();
  } catch (e) {
    $("#tiles", view).innerHTML = banner("err", "Nie udało się połączyć z API: " + esc(e.message));
    return;
  }
  $("#tiles", view).innerHTML = tasks
    .map(
      (t) => `
    <a class="tile" href="#/${t.id}" style="--tile-accent:${accentOf(t.id)}">
      <div class="num">ZADANIE ${esc(t.num)}</div>
      <h3>${esc(t.title)}</h3>
      <p>${esc(t.tag)}</p>
      <div class="go">otwórz →</div>
    </a>`
    )
    .join("");
}

/* ==================================================================
   z01 — układ równań liniowych
   ================================================================== */

VIEWS.z01 = (view) => {
  view.innerHTML = `
  <a class="back" href="#/">← wszystkie zadania</a>
  ${taskHead("z01")}
  ${panel(`
  <p class="desc mt0">Rozwiązujemy układ <b>N równań liniowych o N niewiadomych</b> metodą
  <b>eliminacji Gaussa z wyborem wiodącego (Gauss-Jordan)</b>. Dane wczytywane są z pliku tekstowego:
  w pierwszym wierszu <code>N</code>, w kolejnych współczynniki i wyrazy wolne (N+1 liczb).
  Program klasyfikuje układ na: <b>jednoznaczny</b>, <b>sprzeczny</b> (brak rozwiązań — wiersz <code>0 = c ≠ 0</code>)
  lub <b>nieoznaczony</b> (nieskończenie wiele rozwiązań, podanych parametrycznie przez niewiadome wolne).
  Krok po kroku pokazujemy operacje na wierszach (animacja redukcji do postaci RREF).</p>`, "Opis rozwiązania")}
  ${panel(`
  <label class="field" for="in-text">Dane układu (format pliku: N, potem N wierszy po N+1 liczb)</label>
  <textarea id="in-text" rows="5">3
1 2 3 7
-1 2 4 6
2 1 1 13</textarea>
  <div class="preset-row">
    <button data-p="3\\n1 2 3 7\\n-1 2 4 6\\n2 1 1 13" type="button">przykład z zadania</button>
    <button data-p="2\\n1 1 3\\n2 2 6" type="button">nieoznaczony</button>
    <button data-p="2\\n2 1 1\\n4 2 6" type="button">sprzeczny</button>
  </div>
  <div class="actions"><button class="btn" id="run" type="button">Rozwiąż</button></div>
  <div id="out"></div>`)}
  `
  ;
  const out = $("#out", view);
  let steps = [], stepIdx = 0, playTimer = null;
  // animacja kroców musi umrzeć razem z widokiem (inaczej interwał tyka
  // po zmianie routu i pisze do martwego DOM)
  onLeave(() => {
    if (playTimer) { clearInterval(playTimer); playTimer = null; }
  });

  const drawMatrix = (m, rrefPivots) => {
    const n = m.length;
    let h = `<div class="matrix-wrap"><table class="matrix"><tr>`;
    for (let j = 0; j < n; j++) h += `<th>x${j + 1}</th>`;
    h += `<th>b</th></tr>`;
    for (let i = 0; i < n; i++) {
      h += `<tr>`;
      for (let j = 0; j < n + 1; j++) {
        const v = m[i][j];
        const cls = [];
        if (j === n) cls.push("bar");
        if (rrefPivots && j < n && rrefPivots[i] === j) cls.push("pv");
        h += `<td class="${cls.join(" ")}">${fmt(v, 4)}</td>`;
      }
      h += `</tr>`;
    }
    h += `</table></div>`;
    return h;
  };

  const drawStep = () => {
    const st = steps[stepIdx];
    if (!st) return;
    let pivots = null;
    if (stepIdx === steps.length - 1 && stMatrixIsRREF) {
      pivots = st.matrix.map((row) => {
        for (let j = 0; j < row.length - 1; j++) if (Math.abs(row[j]) > 0.5) return j;
        return -1;
      });
    }
    $("#mx", out).innerHTML = drawMatrix(st.matrix, pivots);
    $("#stepdesc", out).textContent = st.description;
    $("#counter", out).textContent = `${stepIdx + 1} / ${steps.length}`;
    $("#prev", out).disabled = stepIdx === 0;
    $("#next", out).disabled = stepIdx === steps.length - 1;
  };

  let stMatrixIsRREF = false;

  $("#run", view).onclick = async () => {
    out.innerHTML = `<div class="muted mt14">Rozwiązywanie…</div>`;
    try {
      const res = await api("POST", "/api/z01/solve", { text: $("#in-text", view).value });
      steps = res.steps;
      stepIdx = 0;
      stMatrixIsRREF = true;
      let cls, msg;
      if (res.status === "unique") {
        cls = "ok";
        msg = "Układ jednoznaczny — dokładnie jedno rozwiązanie:";
      } else if (res.status === "none") {
        cls = "err";
        msg = "Układ sprzeczny — nie ma rozwiązania.";
      } else {
        cls = "warn";
        msg = "Układ nieoznaczony — nieskończenie wiele rozwiązań (parametrycznie):";
      }
      const lines = res.solutionLines
        .map((l) => `<span class="chip"><b>${esc(l.var)}</b> = ${esc(l.expr)}</span>`)
        .join(" ");
      out.innerHTML =
        banner(cls, msg + ` <div class="chips">${lines}</div>`) +
        `<div class="kv"><span>rank = <b>${res.rank}</b> / ${res.n}</span>` +
        (res.freeVars ? `<span>niewiadome wolne: <b>${res.freeVars.map(esc).join(", ")}</b></span>` : "") +
        `<span>operacji wierszowych: <b>${res.steps.length - 1}</b></span></div>` +
        (res.stepsTruncated
          ? `<div class="muted mt8">Animacja pokazuje pierwsze ${res.steps.length} kroków (limit pamięci dla dużych układów) — wynik końcowy i rozwiązania są pełne.</div>`
          : "") +
        `<div class="flex-between mt14">
           <div class="muted">Redukcja do postaci RREF (krok po kroku):</div>
           <div class="stepnav" style="margin:0">
             <button class="btn ghost small" id="prev" type="button">‹ poprzedni</button>
             <span class="counter" id="counter"></span>
             <button class="btn ghost small" id="next" type="button">następny ›</button>
             <button class="btn small" id="play" type="button">▶ auto</button>
           </div>
         </div>
         <div id="mx"></div>
         <div class="stepdesc" id="stepdesc"></div>`;
      drawStep();
      $("#prev", out).onclick = () => { if (stepIdx > 0) { stepIdx--; drawStep(); } };
      $("#next", out).onclick = () => { if (stepIdx < steps.length - 1) { stepIdx++; drawStep(); } };
      $("#play", out).onclick = () => {
        if (playTimer) { clearInterval(playTimer); playTimer = null; $("#play", out).textContent = "▶ auto"; return; }
        $("#play", out).textContent = "⏸ pauza";
        playTimer = setInterval(() => {
          if (stepIdx < steps.length - 1) { stepIdx++; drawStep(); }
          else { clearInterval(playTimer); playTimer = null; $("#play", out).textContent = "▶ auto"; }
        }, 700);
      };
    } catch (e) {
      out.innerHTML = banner("err", esc(e.message));
    }
  };

  $$(".preset-row button", view).forEach((b) => {
    b.onclick = () => {
      $("#in-text", view).value = b.dataset.p.replace(/\\n/g, "\n");
    };
  });
};

/* ==================================================================
   z02 — liczby Armstronga
   ================================================================== */

VIEWS.z02 = (view) => {
  view.innerHTML = `
  <a class="back" href="#/">← wszystkie zadania</a>
  ${taskHead("z02")}
  ${panel(`
  <p class="desc mt0">Liczba Armstronga to N-cyfrowa liczba naturalna równa sumie swoich cyfr
  podniesionych do potęgi N (np. <code>153 = 1³+5³+3³</code>). Program przeszukuje przedział
  <code>[0, limit]</code> metodą wyczerpującą, ale z <b>przyrostowym obliczaniem sumy potęg cyfr</b>:
  przechodząc od n do n+1 suma zmienia się o stały przyrost <code>(b+1)ᵈ − bᵈ</code> (b — ostatnia cyfra),
  a przeniesienie zdarza się co 10 liczb. Przeszukiwanie jest <b>równoległe</b> — przedział dzielony
  jest na bloki po 4 mln liczb przetwarzanych przez gorutyny (runtime.GOMAXPROCS). Limit do 10¹².</p>`, "Opis rozwiązania")}
  ${panel(`
  <div class="row">
    <div>
      <label class="field" for="limit">Limit przeszukiwania (do 10¹²)</label>
      <select id="limit">
        <option value="1000000">1 000 000 (10⁶) — ok. 1 s</option>
        <option value="10000000">10 000 000 (10⁷) — ok. 2 s</option>
        <option value="100000000" selected>100 000 000 (10⁸) — zakres z wersji Python, ok. 5 s</option>
        <option value="1000000000">1 000 000 000 (10⁹) — ok. 10 s</option>
        <option value="10000000000">10 000 000 000 (10¹⁰) — ok. 30–60 s</option>
        <option value="100000000000">100 000 000 000 (10¹¹) — ok. 5–10 min (job w tle)</option>
        <option value="1000000000000">1 000 000 000 000 (10¹²) — zmierzono 31 min (job w tle)</option>
      </select>
    </div>
  </div>
  <div class="actions">
    <button class="btn" id="run" type="button">Szukaj</button>
    <span class="muted">wyniki pojawiają się na bieżąco</span>
  </div>
  ${jobBox()}
  <div id="found"></div>`)}
  ${panel(`<div class="muted">Ładowanie zweryfikowanej listy…</div>`, "Znane liczby Armstronga (do 10¹²)")}
  <section class="panel" id="known-panel"></section>
  `
  ;
  const foundEl = $("#found", view);
  let found = [];
  let cap = 0;

  $("#run", view).onclick = async () => {
    const limit = parseInt($("#limit", view).value, 10);
    cap = 2e6;
    devState.capKeys = 2e6;
    found = [];
    foundEl.innerHTML = "";
    let handle;
    try {
      const { jobId } = await api("POST", "/api/z02/search", { limit });
      handle = startJob(view, jobId, "/api/z02/search/" + jobId);
    } catch (e) {
      foundEl.innerHTML = banner("err", esc(e.message));
      return;
    }
    handle.onSettled = (s) => {
      if (s.status === "done" && s.result) {
        found = s.result.found || [];
        renderFound();
      } else if (s.status === "error") {
        foundEl.insertAdjacentHTML("beforeend", banner("err", esc(s.error)));
      }
    };
  };

  const renderFound = () => {
    if (!found.length) {
      foundEl.innerHTML = `<div class="muted mt14">Brak wyników.</div>`;
      return;
    }
    foundEl.innerHTML =
      `<div class="kv"><span>znaleziono: <b>${found.length}</b> liczb Armstronga</span></div>
       <div class="chips">` +
      found
        .map((n) => `<span class="chip">${fmtInt(n)}<span class="d">${digitCount(n)} cyfr</span></span>`)
        .join("") +
      `</div>`;
  };

  // znane liczby
  api("GET", "/api/z02/known")
    .then((d) => {
      const byDigits = {};
      d.numbers.forEach((n) => {
        const k = digitCount(n);
        (byDigits[k] = byDigits[k] || []).push(n);
      });
      const keys = Object.keys(byDigits).map(Number).sort((a, b) => a - b);
      let html = `<table class="data"><tr><th>liczba cyfr</th><th>ile</th><th>liczby</th></tr>`;
      for (const k of keys) {
        html += `<tr><td class="num">${k}</td><td class="num">${byDigits[k].length}</td>
          <td style="white-space:normal">${byDigits[k].map((n) => fmtInt(n)).join(", ")}</td></tr>`;
      }
      html += `</table><p class="muted mt8">Zweryfikowane wyczerpującym przeszukaniem Go do 10¹² (12 gorutyn).</p>`;
      $("#known-panel", view).innerHTML =
        `<h4>Znane liczby Armstronga (do 10¹²)</h4>` + html;
    })
    .catch((e) => {
      $("#known-panel", view).innerHTML = banner("err", esc(e.message));
    });
};

/* ==================================================================
   z03 — sito Eratostenesa
   ================================================================== */

VIEWS.z03 = (view) => {
  const Ns = [1e3, 1e4, 1e5, 1e6, 1e7, 1e8];
  view.innerHTML = `
  <a class="back" href="#/">← wszystkie zadania</a>
  ${taskHead("z03")}
  ${panel(`
  <p class="desc mt0">Trzy implementacje sita, jak w zadaniu (w Pythonie: lista, słownik, tablica numpy —
  tutaj ich odpowiedniki w Go):
  <ul>
    <li><code>[]bool</code> — <b>wycinek</b> (analog listy): prosta tablica flag;</li>
    <li><code>map[int]bool</code> — <b>mapa</b> (analog słownika): klucz = liczba, wartość = flaga;</li>
    <li>bitmapa zbitowa — <b>tablica</b> (analog tablicy numpy): gęsty zbiór bitów, 8 flag w 1 bajcie.</li>
  </ul>
  Dla każdej struktury mierzony jest czas i liczba liczb pierwszych dla n = 10³…10⁸;
  wyniki porównujemy w tabeli i na wykresie (skala logarytmiczna).</p>`, "Opis rozwiązania")}
  ${panel(`
  <label class="field" for="ns-1000">Wartości n do pomiaru</label>
  <div class="preset-row" id="ns">
    ${Ns.map((n) => `<label class="chip" for="ns-${n}" style="cursor:pointer"><input type="checkbox" id="ns-${n}" value="${n}" ${n <= 1e7 ? "checked" : ""} style="margin-right:6px">10${sup(Math.log10(n))}</label>`).join("")}
  </div>
  <div class="actions"><button class="btn" id="run" type="button">Zmierz czasy</button></div>
  ${jobBox()}
  <div id="out"></div>`)}
  `
  ;
  const out = $("#out", view);

  function sup(n) {
    const m = { 3: "³", 4: "⁴", 5: "⁵", 6: "⁶", 7: "⁷", 8: "⁸" };
    return m[n] || String(n);
  }

  $("#run", view).onclick = async () => {
    const ns = $$("#ns input:checked", view).map((c) => parseInt(c.value, 10));
    if (!ns.length) { out.innerHTML = banner("warn", "Zaznacz co najmniej jedno n."); return; }
    ns.sort((a, b) => a - b);
    out.innerHTML = "";
    let handle;
    try {
      const { jobId } = await api("POST", "/api/z03/benchmark", { ns });
      handle = startJob(view, jobId, "/api/z03/benchmark/" + jobId);
    } catch (e) {
      out.innerHTML = banner("err", esc(e.message));
      return;
    }
    handle.onSettled = (s) => {
      if (s.status === "done" && s.result) renderResults(s.result);
      else if (s.status === "error") out.insertAdjacentHTML("beforeend", banner("err", esc(s.error)));
    };
  };

  const structures = ["slice (lista)", "map (słownik)", "bitmapa (tablica)"];
  const structColor = { "slice (lista)": "#4f9cf7", "map (słownik)": "#e36ea8", "bitmapa (tablica)": "#4fd78f" };

  function renderResults(rows) {
    const byN = {};
    rows.forEach((r) => { (byN[r.n] = byN[r.n] || {})[r.structure] = r; });
    const ns = Object.keys(byN).map(Number).sort((a, b) => a - b);
    let html = `<table class="data mt14"><tr><th>struktura</th>` +
      ns.map((n) => `<th class="num">n = 10${sup(Math.log10(n))}</th>`).join("") + `</tr>`;
    for (const sName of structures) {
      html += `<tr><td>${esc(sName)}</td>` + ns.map((n) => {
        const r = byN[n][sName];
        if (!r) return `<td class="num">—</td>`;
        return `<td class="num" title="${r.count} liczb pierwszych">${fmtSec(r.seconds)}</td>`;
      }).join("") + `</tr>`;
    }
    html += `</table>`;

    // liczby pierwsze dla największego n
    const maxN = ns[ns.length - 1];
    const counts = structures
      .map((sName) => (byN[maxN][sName] ? byN[maxN][sName].count : null))
      .filter((x) => x !== null);
    if (counts.length) {
      html += `<div class="kv"><span>dla n = 10${sup(Math.log10(maxN))}: <b>${fmtInt(counts[0])}</b> liczb pierwskich (wszystkie struktury zgadzają się)</span></div>`;
    }

    // wykres porównawczy — skala logarytmiczna, 100% = najwolniejsza struktura
    html += `<h4 class="mt14">Czasy działania</h4>
      <div class="muted">Dla każdego n paski odniesione są do <b>najwolniejszej</b> struktury (= 100%);
      najszybsza = 0%. Skala jest logarytmiczna, bo różnice sięgają dwóch rzędów wielkości
      (np. 10⁷: mapa 3,87 s vs wycinek 40 ms) — w skali liniowej szybsze struktury byłyby
      niewidoczne. Wartości bezwzględne podane po prawej.</div>
      <div class="bars">`;
    for (const n of ns) {
      html += `<div class="muted" style="margin-top:14px">n = 10${sup(Math.log10(n))}</div>`;
      const vals = structures
        .map((sName) => ({ name: sName, v: byN[n][sName] ? byN[n][sName].seconds : 0 }))
        .filter((x) => x.v > 0);
      if (!vals.length) continue;
      const vmin = Math.min(...vals.map((x) => x.v));
      const vmax = Math.max(...vals.map((x) => x.v));
      for (const v of vals) {
        // 0% = najszybsza, 100% = najwolniejsza (odniesienie), interpolacja logarytmiczna
        const frac = vmax > vmin ? Math.log10(v.v / vmin) / Math.log10(vmax / vmin) : 1;
        const col = structColor[v.name];
        html += `<div class="barrow logscale">
          <span class="pos"></span>
          <span class="word" style="color:${col}">${esc(v.name)}</span>
          <span class="track"><span class="fill" style="width:${(100 * Math.max(0.02, Math.min(1, frac))).toFixed(1)}%;background:linear-gradient(90deg, color-mix(in srgb, ${col} 35%, #10141c), ${col})"></span></span>
          <span class="cnt">${fmtSec(v.v)}</span>
        </div>`;
      }
    }
    html += `</div>`;
    out.innerHTML = html;
  }
};

/* ==================================================================
   z04 — częstość słów
   ================================================================== */

VIEWS.z04 = (view) => {
  view.innerHTML = `
  <a class="back" href="#/">← wszystkie zadania</a>
  ${taskHead("z04")}
  ${panel(`
  <p class="desc mt0">Program wczytuje tekst, <b>ignoruje słowa krótsze niż 5 liter</b> (typu „i”, „lub”,
  „się”, „aby”, „żeby”) i wypisuje <b>20 najczęściej występujących słów</b> (zliczanie bez różnicowania
  wielkości liter). Słowo = maksymalny ciąg liter (Unicode). Tekst <i>„Pana Tadeusza”</i> jest
  <b>wbudowany w binarię Go</b> (go:embed, ~476 kB) — wynik domyślny nie wymaga żadnego pliku.
  Można też wkleić własny tekst i policzyć jego słowa.</p>`, "Opis rozwiązania")}
  ${panel(`<div class="muted" id="builtin">Ładowanie top 20 dla „Pana Tadeusza”…</div>`, "Top 20 — „Pan Tadeusz” (Mickiewicz)")}
  ${panel(`
  <label class="field" for="text">Własny tekst</label>
  <textarea id="text" rows="5" placeholder="Wklej dowolny tekst…"></textarea>
  <div class="actions"><button class="btn" id="run" type="button">Zlicz słowa</button></div>
  <div id="out"></div>`)}
  `
  ;
  const renderBars = (el, words, title) => {
    if (!words || !words.length) {
      el.innerHTML = banner("warn", "Brak słów o długości ≥ 5 liter.");
      return;
    }
    const max = words[0].count;
    el.innerHTML = (title ? `<div class="muted mt14">${esc(title)}</div>` : "") +
      `<div class="bars">` +
      words
        .map(
          (w, i) => `<div class="barrow">
          <span class="pos">${i + 1}</span>
          <span class="word">${esc(w.word)}</span>
          <span class="track"><span class="fill" style="width:${(100 * w.count / max).toFixed(1)}%"></span></span>
          <span class="cnt">${fmtInt(w.count)}×</span>
        </div>`
        )
        .join("") +
      `</div>`;
  };

  api("GET", "/api/z04/top20")
    .then((d) => renderBars($("#builtin", view), d.words))
    .catch((e) => ($("#builtin", view).innerHTML = banner("err", esc(e.message))));

  $("#run", view).onclick = async () => {
    try {
      const d = await api("POST", "/api/z04/top20", { text: $("#text", view).value, n: 20 });
      renderBars($("#out", view), d.words, "Wynik dla własnego tekstu:");
    } catch (e) {
      $("#out", view).innerHTML = banner("err", esc(e.message));
    }
  };
};

/* ==================================================================
   z05 — szyfr Gronsfelda
   ================================================================== */

VIEWS.z05 = (view) => {
  view.innerHTML = `
  <a class="back" href="#/">← wszystkie zadania</a>
  ${taskHead("z05")}
  ${panel(`
  <p class="desc mt0">Szyfr Gronsfelda to modyfikacja szyfru Cezara z <b>kluczem liczbowym</b>: kolejne litery
  przesuwane są o kolejne cyfry klucza (klucz cykluje się po swojej długości). Alfabet ma 27 znaków:
  <span class="inline-code">␣ a b c d e f g h i j k l m n o p q r s t u v w x y z</span>
  (spacja + 26 małych liter). Np. klucz <code>312063</code> zamienia „computer science” w „fpop whsbsilhoee”.
  <b>Atak słownikowy</b> generuje kandydackie klucze (długości 1…9, kolejność leksykograficzna),
  odszyfrowuje nimi szyfrogram i sprawdza, czy <b>&gt; 80% słów</b> wyniku znajduje się w wbudowanym
  słowniku angielskim (25 991 słów, go:embed). Pierwszy pasujący klucz wygrywa.</p>`, "Opis rozwiązania")}
  ${panel(`
  <div class="tabs" id="tabs">
    <button data-t="enc" class="active" type="button">Szyfruj</button>
    <button data-t="dec" type="button">Deszyfruj</button>
    <button data-t="break" type="button">Złam szyfrogram</button>
  </div>
  <div id="tab-enc">
    <div class="row">
      <div><label class="field" for="enc-text">Tekst jawny</label><input type="text" id="enc-text" value="computer science"></div>
      <div><label class="field" for="enc-key">Klucz (cyfry 0–9, 1–9 cyfr)</label><input type="text" id="enc-key" value="312063"></div>
    </div>
    <div class="actions"><button class="btn" id="enc-run" type="button">Szyfruj</button></div>
    <div id="enc-out"></div>
  </div>
  <div id="tab-dec" hidden>
    <div class="row">
      <div><label class="field" for="dec-text">Szyfrogram</label><input type="text" id="dec-text" value="fpop whsbsilhoee"></div>
      <div><label class="field" for="dec-key">Klucz</label><input type="text" id="dec-key" value="312063"></div>
    </div>
    <div class="actions"><button class="btn" id="dec-run" type="button">Deszyfruj</button></div>
    <div id="dec-out"></div>
  </div>
  <div id="tab-break" hidden>
    <div class="row">
      <div><label class="field" for="br-cipher">Szyfrogram do złamania</label>
        <input type="text" id="br-cipher" value="axjmrlamuhnvayum jsbkvggpzgz zqmna tjpmh"></div>
      <div><label class="field" for="br-len">Długość klucza</label>
        <select id="br-len">
          <option value="6" selected>do 6 cyfr (domyślnie — do 111 111 kluczy, szybkie)</option>
          <option value="7">do 7 cyfr</option>
          <option value="8">do 8 cyfr (wolne)</option>
          <option value="9">do 9 cyfr (max, do 100 mln kluczy — job w tle)</option>
        </select></div>
    </div>
    <div class="preset-row">
      <button data-c="fpop whsbsilhoee" type="button">przykład z zadania (computer science)</button>
      <button data-c="axjmrlamuhnvayum jsbkvggpzgz zqmna tjpmh" type="button">przysłowie z zadania</button>
    </div>
    <div class="actions"><button class="btn" id="br-run" type="button">Rozpocznij atak</button></div>
    ${jobBox()}
    <div id="br-out"></div>
  </div>
  `)}
  `
  ;

  const tabs = { enc: $("#tab-enc", view), dec: $("#tab-dec", view), break: $("#tab-break", view) };
  $$("#tabs button", view).forEach((b) => {
    b.onclick = () => {
      $$("#tabs button", view).forEach((x) => x.classList.remove("active"));
      b.classList.add("active");
      Object.entries(tabs).forEach(([k, el]) => (el.hidden = k !== b.dataset.t));
    };
  });

  /* skippedNote — normalizacja wycina wszystko spoza alfabetu a–z + spacja
     i (od P2) mówi, co wycięła. Notatka jest obowiązkowa w UI: bez niej
     „ZAŻÓŁĆ GĘŚLĄ” wygląda jak poprawne szyfrowanie dwuznakowego wyniku. */
  const skippedNote = (skipped) =>
    skipped && skipped.length
      ? banner("warn",
          `pominięte znaki: <b class="mono">${esc(skipped.join(" "))}</b> — ` +
          `alfabet szyfru to a–z i spacja, więc te znaki nie weszły do wyniku.`)
      : "";

  $("#enc-run", view).onclick = async () => {
    try {
      const d = await api("POST", "/api/z05/crypt", { mode: "enc", text: $("#enc-text", view).value, key: $("#enc-key", view).value });
      $("#enc-out", view).innerHTML = skippedNote(d.skipped) + `<div class="banner ok">Szyfrogram: <span class="mono" style="font-size:15px">${esc(d.out)}</span></div>`;
    } catch (e) { $("#enc-out", view).innerHTML = banner("err", esc(e.message)); }
  };

  $("#dec-run", view).onclick = async () => {
    try {
      const d = await api("POST", "/api/z05/crypt", { mode: "dec", text: $("#dec-text", view).value, key: $("#dec-key", view).value });
      $("#dec-out", view).innerHTML = skippedNote(d.skipped) + `<div class="banner ok">Tekst jawny: <span class="mono" style="font-size:15px">${esc(d.out)}</span></div>`;
    } catch (e) { $("#dec-out", view).innerHTML = banner("err", esc(e.message)); }
  };

  $$("#tab-break .preset-row button", view).forEach((b) => {
    b.onclick = () => { $("#br-cipher", view).value = b.dataset.c; };
  });

  $("#br-run", view).onclick = async () => {
    const len = parseInt($("#br-len", view).value, 10);
    devState.capKeys = len >= 7 ? 1e8 : 2e6;
    $("#br-out", view).innerHTML = "";
    let handle;
    try {
      const { jobId } = await api("POST", "/api/z05/break", { cipher: $("#br-cipher", view).value, maxLen: len, capKeys: devState.capKeys });
      handle = startJob(view, jobId, "/api/z05/break/" + jobId);
    } catch (e) {
      $("#br-out", view).innerHTML = banner("err", esc(e.message));
      return;
    }
    handle.onSettled = (s) => {
      if (s.status === "done" && s.result) {
        const r = s.result;
        if (r.found) {
          $("#br-out", view).innerHTML =
            skippedNote(r.skipped) +
            banner("ok", `Złamano! <div class="chips">
              <span class="chip">klucz: <b>${esc(r.key)}</b></span>
              <span class="chip">tekst jawny: <b>${esc(r.plain.trim())}</b></span></div>`) +
            `<div class="kv"><span>wypróbowano kluczy: <b>${fmtInt(r.keysTried)}</b></span>
             <span>czas: <b>${fmtSec(r.timeMs / 1000)}</b></span></div>`;
        } else {
          $("#br-out", view).innerHTML = skippedNote(r.skipped) + banner("warn",
            `Nie znaleziono klucza po wypróbowaniu <b>${fmtInt(r.keysTried)}</b> kluczy. Spróbuj większej długości lub dłuższego szyfrogramu (więcej słów = mocniejsza weryfikacja).`);
        }
      } else if (s.status === "error") {
        $("#br-out", view).innerHTML = banner("err", esc(s.error));
      }
    };
  };
};

/* ==================================================================
   z06 — wykresy funkcji
   ================================================================== */

VIEWS.z06 = (view) => {
  view.innerHTML = `
  <a class="back" href="#/">← wszystkie zadania</a>
  ${taskHead("z06")}
  ${panel(`
  <p class="desc mt0">Wzór funkcji i przedział dla x podaje użytkownik. Wzór parsuje
  <b>własny parser wstępu rekurencyjnego</b> (analog sympy z wersji Python) — operatory
  <code>+ − · / ^</code>, nawiasy, implikowane mnożenie (<code>2x</code>), funkcje
  <code>sin cos tan asin acos atan sinh cosh tanh sqrt cbrt ln log log10 exp abs floor ceil round min max</code>
  oraz stałe <code>x, pi, e</code>. Funkcja jest próbkowana numerycznie i wykres
  <b>przerywany w miejscach asymptot pionowych</b> (wykrywanie: „wyrzuty” |y| powyżej
  mediany + 15·IQR oraz zmiana znaku z dużym skokiem między próbkami).
  Np. <code>1/x</code> rysuje się jako dwie gałęzie z przerywaną asymptotą w x = 0.</p>`, "Opis rozwiązania")}
  ${panel(`
  <div class="row">
    <div style="flex:2"><label class="field" for="expr">Wzór funkcji y = f(x)</label><input type="text" id="expr" value="x*x-6*x+3"></div>
    <div><label class="field" for="xmin">x min</label><input type="number" id="xmin" value="0" step="any"></div>
    <div><label class="field" for="xmax">x max</label><input type="number" id="xmax" value="6" step="any"></div>
    <div><label class="field" for="samples">próbki</label><input type="number" id="samples" value="1601" min="501" max="20001"></div>
  </div>
  <div class="preset-row">
    <button data-e="x*x-6*x+3" data-a="0" data-b="6" type="button">x²−6x+3 (zadanie)</button>
    <button data-e="1/x" data-a="-5" data-b="5" type="button">1/x (asymptota)</button>
    <button data-e="tan(x)" data-a="-6.3" data-b="6.3" type="button">tan(x)</button>
    <button data-e="sin(x)/x" data-a="-8" data-b="8" type="button">sin(x)/x</button>
    <button data-e="1/(x-2)" data-a="0" data-b="4" type="button">1/(x−2)</button>
    <button data-e="ln(x)" data-a="0.1" data-b="8" type="button">ln(x)</button>
    <button data-e="x^3 - 3*x" data-a="-2.5" data-b="2.5" type="button">x³−3x</button>
  </div>
  <div class="actions"><button class="btn" id="run" type="button">Rysuj</button></div>
  <div id="out"></div>`)}
  `
  ;
  const out = $("#out", view);

  const draw = (plot) => {
    const W = 780, H = 470, M = { l: 60, r: 18, t: 20, b: 44 };
    const pw = W - M.l - M.r, ph = H - M.t - M.b;
    let ymin = Infinity, ymax = -Infinity;
    for (const seg of plot.segments)
      for (const p of seg)
        if (Math.abs(p.y) <= plot.yclamp) { ymin = Math.min(ymin, p.y); ymax = Math.max(ymax, p.y); }
    if (!isFinite(ymin)) { ymin = -10; ymax = 10; }
    let span = ymax - ymin;
    if (span < 1e-9) { ymin -= 5; ymax += 5; span = 10; }
    ymin -= span * 0.08; ymax += span * 0.08;
    const X = (v) => M.l + ((v - plot.xmin) / (plot.xmax - plot.xmin)) * pw;
    const Y = (v) => M.t + (1 - (v - ymin) / (ymax - ymin)) * ph;
    const xt = niceTicks(plot.xmin, plot.xmax, 8);
    const yt = niceTicks(ymin, ymax, 6);
    let s = `<svg viewBox="0 0 ${W} ${H}" class="chart" role="img" aria-label="wykres">`;
    s += `<defs><clipPath id="clip${W}"><rect x="${M.l}" y="${M.t}" width="${pw}" height="${ph}"/></clipPath></defs>`;
    s += `<rect x="${M.l}" y="${M.t}" width="${pw}" height="${ph}" fill="#10141c" stroke="#1e2635"/>`;
    for (const t of xt) {
      const x = X(t);
      s += `<line class="grid-line" x1="${x}" y1="${M.t}" x2="${x}" y2="${M.t + ph}"/>`;
      s += `<text class="tick-label" x="${x}" y="${M.t + ph + 16}" text-anchor="middle">${fmt(t, 3)}</text>`;
    }
    for (const t of yt) {
      const y = Y(t);
      s += `<line class="grid-line" x1="${M.l}" y1="${y}" x2="${M.l + pw}" y2="${y}"/>`;
      s += `<text class="tick-label" x="${M.l - 8}" y="${y + 3.5}" text-anchor="end">${fmt(t, 3)}</text>`;
    }
    if (ymin < 0 && ymax > 0) s += `<line class="zero" x1="${M.l}" y1="${Y(0)}" x2="${M.l + pw}" y2="${Y(0)}"/>`;
    if (plot.xmin < 0 && plot.xmax > 0) s += `<line class="zero" x1="${X(0)}" y1="${M.t}" x2="${X(0)}" y2="${M.t + ph}"/>`;
    for (const a of plot.asymptotes) {
      const x = X(a);
      s += `<line class="asym" x1="${x}" y1="${M.t}" x2="${x}" y2="${M.t + ph}"/>`;
      s += `<text class="asym-label" x="${x + 5}" y="${M.t + 15}">x ≈ ${niceAsym(a)}</text>`;
    }
    s += `<g clip-path="url(#clip${W})">`;
    for (const seg of plot.segments) {
      if (seg.length < 2) continue;
      const pts = seg.map((p) => `${X(p.x).toFixed(2)},${Y(p.y).toFixed(2)}`).join(" ");
      s += `<polyline class="fn" points="${pts}"/>`;
    }
    s += `</g></svg>`;
    out.innerHTML =
      `<div class="kv"><span>wzór: <b>y = ${esc(plot.formula)}</b></span>
       <span>przedział: <b>[${fmt(plot.xmin, 3)}, ${fmt(plot.xmax, 3)}]</b></span>
       <span>próbek: <b>${fmtInt(plot.samples)}</b>${
         plot.samplesRequested && plot.samplesRequested !== plot.samples
           ? ` <span class="muted">(żądano ${fmtInt(plot.samplesRequested)} — przycinamy do 501…20001)</span>`
           : ""}</span>
       <span>odcinków: <b>${plot.segments.length}</b></span>
       <span>asymptoty: <b>${plot.asymptotes.length ? plot.asymptotes.map((a) => "x≈" + niceAsym(a)).join(", ") : "brak"}</b></span></div>
       <div class="chart-wrap">${s}</div>`;
  };

  async function run() {
    try {
      const plot = await api("POST", "/api/z06/plot", {
        expr: $("#expr", view).value,
        xmin: parseFloat($("#xmin", view).value),
        xmax: parseFloat($("#xmax", view).value),
        samples: parseInt($("#samples", view).value, 10) || 1601,
      });
      draw(plot);
    } catch (e) {
      out.innerHTML = banner("err", "Błąd: " + esc(e.message));
    }
  }

  $("#run", view).onclick = run;
  $$(".preset-row button", view).forEach((b) => {
    b.onclick = () => {
      $("#expr", view).value = b.dataset.e;
      $("#xmin", view).value = b.dataset.a;
      $("#xmax", view).value = b.dataset.b;
      run();
    };
  });
  run();
};

// niceAsym — „przyciąga” wykrytą asymptotę do ładnej wartości (0, liczba
// całkowita, półliczba), jeśli jest w jej pobliżu — m.in. 1/x → x≈0.
function niceAsym(a) {
  for (const c of [Math.round(a), Math.round(a * 2) / 2]) {
    if (Math.abs(a - c) < 0.02) return String(c);
  }
  return fmt(a, 3);
}

function niceTicks(min, max, n) {
  const span = max - min;
  if (span <= 0 || !isFinite(span)) return [min];
  const step0 = span / n;
  const mag = Math.pow(10, Math.floor(Math.log10(step0)));
  const norm = step0 / mag;
  let step;
  if (norm >= 5) step = 10 * mag;
  else if (norm >= 2) step = 5 * mag;
  else if (norm >= 1) step = 2 * mag;
  else step = mag;
  const start = Math.ceil(min / step) * step;
  const out = [];
  for (let v = start; v <= max + step * 1e-6; v += step) out.push(Math.abs(v) < step * 1e-9 ? 0 : v);
  return out;
}

/* ==================================================================
   z07 — zagadki literowo-cyfrowe
   ================================================================== */

VIEWS.z07 = (view) => {
  view.innerHTML = `
  <a class="back" href="#/">← wszystkie zadania</a>
  ${taskHead("z07")}
  ${panel(`
  <p class="desc mt0">Litery A–J w równaniu zastępują cyfry (np. <code>ABC · BD = EFGAH</code> →
  <code>623 · 29 = 18067</code>). Zasady: litera ↔ dokładnie jedna cyfra, liczba wielocyfrowa
  nie zaczyna się od zera, dzielenie jest dokładne. Program wczytuje równanie w postaci
  <code>&lt;liczba&gt; &lt;operator&gt; &lt;liczba&gt; = &lt;liczba&gt;</code> (operator: <code>+ − · /</code>)
  i <b>wypisuje wszystkie rozwiązania</b> (niektóre zagadki mają ich wiele).
  Metoda: wyczerpujące przeszukiwanie permutacji cyfr (10P<sub>k</sub>) z cięciami
  (np. liczba cyfr iloczynu) — w Go 10P₁₀ ≈ 3,6 mln permutacji to ułamek sekundy.</p>`, "Opis rozwiązania")}
  ${panel(`
  <label class="field" for="eq">Równanie w postaci literowej</label>
  <input type="text" id="eq" value="ABC * BD = EFGAH" style="font-family:var(--mono)">
  <div class="preset-row">
    <button data-e="ABC * BD = EFGAH" type="button">z zadania: ABC·BD = EFGAH</button>
    <button data-e="SEND + MORE = MONEY" type="button">SEND + MORE = MONEY</button>
    <button data-e="AB / C = DE" type="button">AB / C = DE</button>
    <button data-e="A * BC = DEFG" type="button">A · BC = DEFG</button>
  </div>
  <div class="actions"><button class="btn" id="run" type="button">Rozwiąż</button></div>
  <div id="out"></div>`)}
  `
  ;
  const out = $("#out", view);

  $$(".preset-row button", view).forEach((b) => {
    b.onclick = () => { $("#eq", view).value = b.dataset.e; };
  });

  $("#run", view).onclick = async () => {
    out.innerHTML = `<div class="muted mt14">Przeszukiwanie permutacji…</div>`;
    try {
      const r = await api("POST", "/api/z07/solve", { equation: $("#eq", view).value });
      if (!r.solutions || !r.solutions.length) {
        out.innerHTML = banner("warn", "Brak rozwiązań dla tej zagadki.") +
          `<div class="kv"><span>przestrzeń: <b>${fmtInt(r.spaceSize)}</b> permutacji</span>
           <span>czas: <b>${r.timeMs} ms</b></span></div>`;
        return;
      }
      let html = banner("ok", `Znaleziono <b>${fmtInt(r.count)}</b> ${r.count === 1 ? "rozwiązanie" : r.count < 5 ? "rozwiązania" : "rozwiązań"}` +
        (r.truncated
          ? ` (pokazano ${fmtInt(r.solutions.length)} z ${fmtInt(r.count)} — pełna liczba jest w liczniku)`
          : "") + ``) +
        `<div class="kv"><span>litery: <b>${r.letters}</b></span>
         <span>przestrzeń permutacji: <b>${fmtInt(r.spaceSize)}</b></span>
         <span>czas: <b>${r.timeMs} ms</b></span></div>
         <div class="chips" style="max-height:340px;overflow:auto">`;
      r.solutions.forEach((s) => {
        const map = Object.entries(s.mapping).sort((a, b) => a[0].localeCompare(b[0]));
        html += `<span class="chip" title="${esc(map.map(([k, v]) => k + "=" + v).join(", "))}">${esc(s.equation)}</span>`;
      });
      html += `</div>`;
      // szczegółowe przypisania dla pierwszego rozwiązania
      const s0 = r.solutions[0];
      html += `<table class="data mt14"><tr><th>równanie</th><th>` +
        Object.keys(s0.mapping).sort().map((k) => `<th class="num">${k}</th>`).join("") + `</tr><tr>
        <td class="mono">${esc(s0.equation)}</td>` +
        Object.keys(s0.mapping).sort().map((k) => `<td class="num"><b style="color:var(--accent)">${s0.mapping[k]}</b></td>`).join("") +
        `</tr></table>`;
      out.innerHTML = html;
    } catch (e) {
      out.innerHTML = banner("err", esc(e.message));
    }
  };
};

/* ==================================================================
   z08 — sudoku
   ================================================================== */

VIEWS.z08 = (view) => {
  view.innerHTML = `
  <a class="back" href="#/">← wszystkie zadania</a>
  ${taskHead("z08")}
  ${panel(`
  <p class="desc mt0">Generator: (1) losowo wypełnia całą planszę 9×9 <b>backtrackingiem z losową
  kolejnością kandydatów</b>; (2) „kopie dołki” — pola w losowej kolejności są wycierane, a po każdym
  usunięciu zliczacz rozwiązań (backtracking + heurystyka MRV, maski bitowe wiersza/kolumny/bloku,
  przerwanie po 2 rozwiązaniach) sprawdza <b>unikalność</b>. Jeśli unikalność znika — pole wraca.
  Powtarzamy, aż żadne pole nie da się usunąć: zagadka ma <b>dokładnie jedno rozwiązanie i wiele
  pustych pól</b> — typowo 55–58 (dla seed 1–8 i 42); „maksimum” jest umowne, bo teoretyczny sufit
  dla jednoznacznej zagadki jest wyższy (~64). Kliknij puste pole i wpisz cyfrę (klawiatura), a potem
  sprawdź lub poproś o rozwiązanie.</p>`, "Opis rozwiązania")}
  ${panel(`
  <div class="toolbar">
    <div class="fld">
      <label class="field" for="seed">seed (0 = losowy)</label>
      <input type="number" id="seed" value="0">
    </div>
    <div class="actions">
      <button class="btn" id="gen" type="button">Generuj</button>
      <button class="btn ghost" id="show" type="button">Pokaż rozwiązanie</button>
      <button class="btn ghost" id="clear" type="button">Wyczyść wpisy</button>
      <button class="btn ghost" id="check" type="button">Sprawdź / rozwiąż</button>
    </div>
  </div>
  <div class="kv" id="stats"></div>
  <div id="grid"></div>
  <div class="legend">kliknij puste pole → cyfra 1–9 (Backspace usuwa) · szary = wskazówka, akcent = Twój wpis</div>`, "Zagadka 9×9 (generowana w Go)")}
  ${panel(`
  <label class="field" for="own">Własna plansza: 9 wierszy × 9 cyfr (0 = puste pole)</label>
  <textarea id="own" rows="9" style="min-height:150px" placeholder="np. 000700000&#10;000005002&#10;700980003&#10;000000000&#10;010500090&#10;006430805&#10;090650000&#10;004009360&#10;008003010"></textarea>
  <div class="actions"><button class="btn" id="own-run" type="button">Zweryfikuj i rozwiąż</button></div>
  <div id="own-out"></div>`)}
  `
  ;

  let puzzle = new Array(81).fill(0);
  let userVal = new Array(81).fill(0);
  let sel = -1;

  const drawGrid = (showSol = false, solution = null) => {
    let html = `<div class="sudoku">`;
    for (let i = 0; i < 81; i++) {
      const r = Math.floor(i / 9), c = i % 9;
      const classes = ["cell"];
      if (c === 2 || c === 5) classes.push("boxr");
      if (r === 2 || r === 5) classes.push("boxb");
      let v = userVal[i];
      let cls = "";
      if (showSol && solution && solution[i]) {
        v = solution[i];
        cls = "sol";
      } else if (puzzle[i]) {
        v = puzzle[i];
        cls = "given";
      } else if (v) {
        cls = "user";
      }
      if (i === sel && !puzzle[i]) classes.push("sel");
      html += `<div class="${classes.concat(cls).join(" ")}" data-i="${i}">${v || ""}</div>`;
    }
    html += `</div>`;
    $("#grid", view).innerHTML = html;
    $$("#grid .cell", view).forEach((el) => {
      el.onclick = () => {
        const i = parseInt(el.dataset.i, 10);
        if (!puzzle[i]) { sel = i; drawGrid(); }
      };
    });
  };

  const onKey = (e) => {
    if (sel < 0) return;
    if (e.key >= "1" && e.key <= "9") {
      userVal[sel] = parseInt(e.key, 10);
      sel = -1;
      drawGrid();
    } else if (e.key === "Backspace" || e.key === "Delete") {
      userVal[sel] = 0;
      sel = -1;
      drawGrid();
    } else if (e.key === "Escape") {
      sel = -1;
      drawGrid();
    }
  };

  // Słuchacz klawiatury wieszanym na dokumencie musi mieć uchwyt:
  // addEventListener + AbortController, wypinany przy zmianie routu.
  // (document.body.onkeydown nadpisywał każdy globalny handler i został
  //  zawieszony na całym dokumencie po wyjściu z widoku.)
  let keyCtrl = null;
  const attachKeys = () => {
    if (keyCtrl) return;
    keyCtrl = new AbortController();
    document.addEventListener("keydown", onKey, { signal: keyCtrl.signal });
    onLeave(() => { keyCtrl.abort(); keyCtrl = null; });
  };

  const gen = async () => {
    const seed = parseInt($("#seed", view).value, 10) || 0;
    $("#stats", view).innerHTML = `<span>generowanie…</span>`;
    try {
      const r = await api("POST", "/api/z08/generate", { seed });
      puzzle = r.puzzle.slice();
      userVal = new Array(81).fill(0);
      sel = -1;
      devState.solution = r.solution;
      $("#stats", view).innerHTML =
        `<span>puste pola: <b>${r.blanks}</b> / 81</span>
         <span>wskazówki: <b>${81 - r.blanks}</b></span>
         <span>czas generacji: <b>${r.timeMs} ms</b></span>`;
      drawGrid();
      attachKeys();
    } catch (e) {
      $("#stats", view).innerHTML = banner("err", esc(e.message));
    }
  };

  $("#gen", view).onclick = gen;
  $("#clear", view).onclick = () => { userVal = new Array(81).fill(0); sel = -1; drawGrid(); };
  $("#show", view).onclick = () => drawGrid(true, devState.solution);
  $("#check", view).onclick = async () => {
    const grid = puzzle.map((v, i) => v || userVal[i]);
    try {
      const r = await api("POST", "/api/z08/solve", { grid });
      if (r.solved) {
        const mine = grid.filter((v, i) => !puzzle[i] && v).length;
        const correct = grid.filter((v, i) => !puzzle[i] && v && r.solution[i] === v).length;
        $("#stats", view).innerHTML =
          banner("ok", `Plansza ma <b>${r.capped ? "≥" + r.count : r.count}</b> ${r.count === 1 ? "rozwiązanie" : "rozwiązań"} (czas: ${r.timeMs} ms).`) +
          (mine ? `<div class="kv"><span>Twoje wpisy: <b>${mine}</b>, zgodne z rozwiązaniem: <b>${correct}</b></span></div>` : "");
        drawGrid(true, r.solution);
      } else {
        $("#stats", view).innerHTML = banner("err", "Plansza <b>nie ma żadnego rozwiązania</b> — sprawdź wpisane cyfry.");
      }
    } catch (e) {
      $("#stats", view).innerHTML = banner("err", esc(e.message));
    }
  };

  $("#own-run", view).onclick = async () => {
    const lines = $("#own", view).value.trim().split(/\s+/);
    const grid = [];
    for (const ln of lines) for (const ch of ln) if (/[0-9]/.test(ch)) grid.push(parseInt(ch, 10));
    try {
      const r = await api("POST", "/api/z08/solve", { grid });
      let html;
      if (r.solved) {
        html = banner("ok", `Rozwiązań: <b>${r.capped ? "≥" + r.count : r.count}</b> (czas: ${r.timeMs} ms).`) +
          `<div class="kv"><span>rozwiązanie (wierszami):</span></div>
           <div class="mono mt8" style="line-height:1.9;white-space:pre">${r.solution.map((v, i) => (i % 9 === 0 ? "\n" : "") + v).join("")}</div>`;
      } else {
        html = banner("warn", "Ta plansza <b>nie ma rozwiązania</b>.");
      }
      $("#own-out", view).innerHTML = html;
    } catch (e) {
      $("#own-out", view).innerHTML = banner("err", esc(e.message));
    }
  };

  gen();
};
