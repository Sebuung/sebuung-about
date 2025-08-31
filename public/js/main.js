// ===== Scroll reveal =====
const revealObserver = new IntersectionObserver(
  (entries) => {
    for (const e of entries) {
      if (e.isIntersecting) {
        e.target.classList.add("in");
        revealObserver.unobserve(e.target);
      }
    }
  },
  { threshold: 0.1 }
);
document.querySelectorAll("[data-animate]").forEach((el) => revealObserver.observe(el));

// ===== Counter (number & currency) =====
function animateCount(el, to, format = "plain", duration = 1200) {
  const start = 0;
  const startTime = performance.now();
  function frame(now) {
    const p = Math.min(1, (now - startTime) / duration);
    const val = Math.floor(start + (to - start) * (1 - Math.pow(1 - p, 3)));
    el.textContent = format === "money" ? val.toLocaleString("ko-KR") : val.toLocaleString("ko-KR");
    if (p < 1) requestAnimationFrame(frame);
  }
  requestAnimationFrame(frame);
}
const countObserver = new IntersectionObserver(
  (entries) => {
    for (const e of entries) {
      if (!e.isIntersecting) continue;
      const el = e.target;
      const to = Number(el.getAttribute("data-to")) || 0;
      const format = el.getAttribute("data-format") || "plain";
      animateCount(el, to, 1400);
      countObserver.unobserve(el);
    }
  },
  { threshold: 0.4 }
);
document.querySelectorAll("[data-count]").forEach((el) => countObserver.observe(el));

// ===== Members: total & "전체" 탭 구성 =====
const panelEls = [...document.querySelectorAll("[data-panel]")];
const allNamesWrap = document.getElementById("all-names");

function collectNames() {
  const names = [];
  panelEls.forEach((panel) => {
    if (panel.dataset.panel === "genAll") return;
    panel.querySelectorAll(".name").forEach((n) => names.push(n.textContent.trim()));
  });
  return names;
}
function renderAllNames() {
  const names = collectNames();
  allNamesWrap.innerHTML = "";
  names.forEach((n) => {
    const span = document.createElement("span");
    span.className = "name";
    span.textContent = n;
    allNamesWrap.appendChild(span);
  });

  // 멤버 총원 16명으로 표기
  const TARGET_MEMBER_TOTAL = 16;
  const totalEl = document.getElementById("member-total");
  totalEl.setAttribute("data-to", String(TARGET_MEMBER_TOTAL));
  if (totalEl.classList.contains("in")) {
    totalEl.textContent = TARGET_MEMBER_TOTAL.toLocaleString("ko-KR");
  }
}
renderAllNames();

// ===== Tabs =====
const tabEls = document.querySelectorAll(".tab");
function activateTab(key) {
  tabEls.forEach((t) => t.setAttribute("aria-selected", t.dataset.tab === key ? "true" : "false"));
  panelEls.forEach((p) => (p.hidden = p.dataset.panel !== key));
}
tabEls.forEach((btn) => btn.addEventListener("click", () => activateTab(btn.dataset.tab)));

// ===== Accessibility: keyboard nav for tabs =====
document.querySelector(".tabs").addEventListener("keydown", (e) => {
  const order = [...tabEls];
  const idx = order.findIndex((b) => b.getAttribute("aria-selected") === "true");
  if (["ArrowRight", "ArrowLeft", "Home", "End"].includes(e.key)) e.preventDefault();
  let next = idx;
  if (e.key === "ArrowRight") next = (idx + 1) % order.length;
  if (e.key === "ArrowLeft") next = (idx - 1 + order.length) % order.length;
  if (e.key === "Home") next = 0;
  if (e.key === "End") next = order.length - 1;
  if (next !== idx) {
    order[next].focus();
    order[next].click();
  }
});

// Optional: activate "전체" by default on narrow screens
if (matchMedia("(max-width: 720px)").matches) activateTab("genAll");

// Smooth scroll for internal links
document.addEventListener("click", (e) => {
  const a = e.target.closest('a[href^="#"]');
  if (!a) return;
  const id = a.getAttribute("href").slice(1);
  const target = document.getElementById(id);
  if (target) {
    e.preventDefault();
    target.scrollIntoView({ behavior: "smooth", block: "start" });
    history.pushState(null, "", `#${id}`);
  }
});

// ===== Mobile nav (hamburger) =====
const navHeader = document.querySelector("header.nav");
const navToggle = document.getElementById("nav-toggle");
const navMenu = document.getElementById("nav-menu");
function setNavOpen(open) {
  if (!navHeader || !navToggle) return;
  navHeader.dataset.open = open ? "true" : "false";
  navToggle.setAttribute("aria-expanded", open ? "true" : "false");
  navToggle.setAttribute("aria-label", open ? "메뉴 닫기" : "메뉴 열기");
}
navToggle?.addEventListener("click", () => {
  const isOpen = navHeader?.dataset.open === "true";
  setNavOpen(!isOpen);
});
// Close when clicking a link in the menu
navMenu?.addEventListener("click", (e) => {
  if (e.target.closest("a")) setNavOpen(false);
});
// Close with Escape key
document.addEventListener("keydown", (e) => {
  if (e.key === "Escape") setNavOpen(false);
});
// Ensure closed when switching to desktop
matchMedia("(min-width: 881px)").addEventListener("change", (e) => {
  if (e.matches) setNavOpen(false);
});
