package assets

const ThemeBootstrapJS = `(function() {
  var root = document.documentElement;
  var key = "marionette-theme";
  var systemQuery = window.matchMedia ? window.matchMedia("(prefers-color-scheme: dark)") : null;
  var labels = { system: "System", corporate: "Light", dark: "Dark" };
  var icons = { system: "◐", corporate: "☀", dark: "☾" };
  var cycle = { system: "corporate", corporate: "dark", dark: "system" };

  function normalizeTheme(value) {
    return value === "dark" || value === "corporate" || value === "system" ? value : "system";
  }

  function resolveTheme(mode) {
    if (mode === "system") {
      return systemQuery && systemQuery.matches ? "dark" : "corporate";
    }
    return mode;
  }

  function updateThemeControls(mode) {
    var controls = document.querySelectorAll ? document.querySelectorAll("[data-mrn-theme-toggle]") : [];
    controls.forEach(function(control) {
      control.setAttribute("aria-pressed", String(mode !== "system"));
      control.setAttribute("aria-label", "Toggle theme: " + (labels[mode] || labels.system));
      control.setAttribute("data-mrn-theme-mode", mode);
      var icon = control.querySelector("[data-mrn-theme-icon]");
      if (icon) icon.textContent = icons[mode] || icons.system;
    });
  }

  function applyTheme(mode, persist) {
    mode = normalizeTheme(mode);
    root.setAttribute("data-theme", resolveTheme(mode));
    root.setAttribute("data-mrn-theme-mode", mode);
    updateThemeControls(mode);
    if (persist) {
      try {
        localStorage.setItem(key, mode);
      } catch (e) {}
    }
  }

  var storedTheme = null;
  try {
    storedTheme = localStorage.getItem(key);
  } catch (e) {}

  applyTheme(normalizeTheme(storedTheme), false);

  window.mrnSetTheme = function(mode) {
    applyTheme(mode, true);
  };

  window.mrnToggleTheme = function() {
    var current = normalizeTheme(root.getAttribute("data-mrn-theme-mode"));
    applyTheme(cycle[current] || "system", true);
  };

  document.addEventListener("DOMContentLoaded", function() {
    updateThemeControls(normalizeTheme(root.getAttribute("data-mrn-theme-mode")));
  });

  if (systemQuery) {
    var handleSystemChange = function() {
      if (normalizeTheme(root.getAttribute("data-mrn-theme-mode")) === "system") {
        applyTheme("system", false);
      }
    };
    if (systemQuery.addEventListener) {
      systemQuery.addEventListener("change", handleSystemChange);
    } else if (systemQuery.addListener) {
      systemQuery.addListener(handleSystemChange);
    }
  }
})();`

const ChartBootstrapJS = `(function() {
  var charts = new WeakMap();
  var observed = new WeakSet();

  function observeChartRoot(container) {
    if (!window.ResizeObserver || observed.has(container)) return;
    observed.add(container);
    var observer = new ResizeObserver(function() {
      scheduleInitCharts(container);
    });
    observer.observe(container);
  }

  function initCharts(root) {
    if (!window.Chart) return;
    var scope = root || document;
    var canvases = scope.querySelectorAll ? scope.querySelectorAll("[data-mrn-chart]") : [];
    canvases.forEach(function(canvas) {
      var container = canvas.closest("[data-mrn-chart-root]");
      if (!container) return;
      observeChartRoot(container);
      var configEl = container.querySelector("[data-mrn-chart-config]");
      if (!configEl) return;
      var frame = canvas.parentElement;
      var rect = frame && frame.getBoundingClientRect ? frame.getBoundingClientRect() : null;
      if (rect && (rect.width < 1 || rect.height < 1)) return;
      canvas.style.display = "block";
      canvas.style.width = "100%";
      canvas.style.height = "100%";

      var config;
      try {
        config = JSON.parse(configEl.textContent || "{}");
      } catch (e) {
        return;
      }

      config.options = config.options || {};
      config.options.onClick = function(_, elements) {
        if (!elements || !elements.length) return;
        var first = elements[0];
        var label = (config.data && config.data.labels && config.data.labels[first.index]) || "";
        var stateName = container.getAttribute("data-mrn-query-state");
        var column = container.getAttribute("data-mrn-filter-column");
        if (!stateName || !column) return;
        var payload = {state: stateName, filters: [{column: column, op: "eq", value: String(label)}]};
        document.dispatchEvent(new CustomEvent("mrn:data-query-change", {detail: payload}));
        if (window.htmx) {
          window.htmx.trigger(document.body, "mrn:data-query-change", payload);
        }
      };

      var existing = charts.get(canvas);
      if (!existing && window.Chart.getChart) {
        existing = window.Chart.getChart(canvas);
      }
      if (existing) existing.destroy();

      var chart;
      try {
        chart = new window.Chart(canvas, config);
      } catch (e) {
        if (window.console && console.error) console.error("Marionette chart init failed", e);
        return;
      }
      charts.set(canvas, chart);
    });
  }

  function scheduleInitCharts(root) {
    initCharts(root);
    if (window.requestAnimationFrame) {
      window.requestAnimationFrame(function() {
        initCharts(root);
      });
    }
    window.setTimeout(function() {
      initCharts(root);
    }, 80);
    window.setTimeout(function() {
      initCharts(root);
    }, 250);
    window.setTimeout(function() {
      initCharts(root);
    }, 750);
  }

  document.addEventListener("DOMContentLoaded", function() {
    scheduleInitCharts(document);
  });
  window.addEventListener("load", function() {
    scheduleInitCharts(document);
  });
  document.addEventListener("htmx:afterSwap", function(event) {
    scheduleInitCharts(event.detail && event.detail.elt ? event.detail.elt : document);
  });
  document.addEventListener("mrn:data-query-change", function(event) {
    var detail = event.detail || {};
    var tables = document.querySelectorAll("[data-mrn-query-state]");
    tables.forEach(function(node) {
      if (node.getAttribute("data-mrn-query-state") !== detail.state) return;
      node.setAttribute("data-mrn-selected-filter", JSON.stringify(detail.filters || []));
    });
  });
  window.mrnInitCharts = initCharts;
})();`
