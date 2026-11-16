(function () {
  var themeButton = document.querySelector(".theme-toggle");
  function syncThemeButton() {
    if (!themeButton) return;
    var dark = document.documentElement.getAttribute("data-theme") === "dark";
    themeButton.setAttribute("aria-pressed", dark ? "true" : "false");
    themeButton.setAttribute("aria-label", dark ? themeButton.getAttribute("data-label-light") : themeButton.getAttribute("data-label-dark"));
  }
  if (themeButton) {
    syncThemeButton();
    themeButton.addEventListener("click", function () {
      var next = document.documentElement.getAttribute("data-theme") === "dark" ? "light" : "dark";
      document.documentElement.setAttribute("data-theme", next);
      try { localStorage.setItem("theme", next); } catch (e) {}
      syncThemeButton();
    });
  }

  var copyLabel = document.body.getAttribute("data-copy-label") || "Copy";
  var copiedLabel = document.body.getAttribute("data-copied-label") || "Copied";
  document.querySelectorAll(".content .highlight").forEach(function (block) {
    var button = document.createElement("button");
    button.type = "button";
    button.className = "copy-code";
    button.textContent = copyLabel;
    button.setAttribute("aria-label", copyLabel);
    block.insertBefore(button, block.firstChild);
    button.addEventListener("click", function () {
      var lines = block.querySelectorAll(".cl");
      var text = "";
      if (lines.length) {
        text = Array.prototype.map.call(lines, function (line) {
          return line.textContent.replace(/\n$/, "");
        }).join("\n");
      } else {
        var code = block.querySelector("code");
        text = code ? code.textContent.replace(/\n$/, "") : "";
      }
      var done = function () {
        button.textContent = copiedLabel;
        button.classList.add("is-copied");
        window.setTimeout(function () {
          button.textContent = copyLabel;
          button.classList.remove("is-copied");
        }, 1500);
      };
      var fallback = function () {
        var area = document.createElement("textarea");
        area.value = text;
        area.setAttribute("readonly", "");
        area.style.position = "fixed";
        area.style.left = "-9999px";
        document.body.appendChild(area);
        area.select();
        var ok = document.execCommand("copy");
        document.body.removeChild(area);
        return ok ? Promise.resolve() : Promise.reject();
      };
      var write = navigator.clipboard && navigator.clipboard.writeText
        ? navigator.clipboard.writeText(text).catch(fallback)
        : fallback();
      write.then(done);
    });
  });

  var lightbox = document.getElementById("lightbox");
  if (lightbox) {
    var view = lightbox.querySelector("img");
    document.querySelectorAll(".content img").forEach(function (el) {
      if (el.closest("a")) return;
      el.addEventListener("click", function () {
        view.src = el.currentSrc || el.src;
        view.alt = el.alt || "";
        lightbox.showModal();
      });
    });
    lightbox.addEventListener("click", function () {
      lightbox.close();
    });
  }

  var toTop = document.querySelector(".to-top");
  if (toTop) {
    var toggleTop = function () {
      toTop.classList.toggle("is-visible", window.scrollY > 400);
    };
    toggleTop();
    window.addEventListener("scroll", toggleTop, { passive: true });
    toTop.addEventListener("click", function () {
      var reduce = window.matchMedia("(prefers-reduced-motion: reduce)").matches;
      window.scrollTo({ top: 0, behavior: reduce ? "auto" : "smooth" });
    });
  }

  var searchDialog = document.getElementById("search");
  var openButton = document.querySelector(".search-open");
  if (!searchDialog || !openButton) return;

  var form = searchDialog.querySelector("form");
  var input = document.getElementById("search-input");
  var results = searchDialog.querySelector(".search-results");
  var empty = searchDialog.querySelector(".search-empty");
  var closeButton = searchDialog.querySelector(".search-close");
  var index = null;
  var loading = false;

  function values(value) {
    if (!value) return [];
    return Array.isArray(value) ? value : [String(value)];
  }

  function render() {
    var query = input.value.trim().toLowerCase();
    results.innerHTML = "";
    empty.hidden = true;
    if (!index || query.length < 1) return;
    var found = 0;
    for (var i = 0; i < index.length && found < 20; i++) {
      var item = index[i];
      var haystack = [item.title, item.summary]
        .concat(values(item.categories), values(item.tags))
        .join(" ")
        .toLowerCase();
      if (haystack.indexOf(query) === -1) continue;
      found++;
      var li = document.createElement("li");
      var link = document.createElement("a");
      link.href = item.url;
      link.textContent = item.title || item.url;
      var date = document.createElement("small");
      date.textContent = item.date || "";
      li.appendChild(link);
      li.appendChild(date);
      results.appendChild(li);
    }
    empty.hidden = found !== 0;
  }

  function loadIndex() {
    if (index || loading) return;
    loading = true;
    var url = document.body.getAttribute("data-search-url");
    fetch(url)
      .then(function (response) {
        if (!response.ok) throw new Error("search");
        return response.json();
      })
      .then(function (data) {
        index = data || [];
        loading = false;
        render();
      })
      .catch(function () {
        loading = false;
        index = [];
        render();
      });
  }

  openButton.addEventListener("click", function () {
    searchDialog.showModal();
    input.focus();
    loadIndex();
  });
  input.addEventListener("input", render);
  form.addEventListener("submit", function (event) {
    event.preventDefault();
  });
  closeButton.addEventListener("click", function () {
    searchDialog.close();
  });
  searchDialog.addEventListener("click", function (event) {
    if (event.target === searchDialog) searchDialog.close();
  });
})();
