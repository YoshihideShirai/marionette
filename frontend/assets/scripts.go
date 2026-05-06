package assets

const ThemeBootstrapJS = `(function() {
  var root = document.documentElement;
  var key = "marionette-theme";
  var storedTheme = null;
  try {
    storedTheme = localStorage.getItem(key);
  } catch (e) {}
  var prefersDark = window.matchMedia && window.matchMedia("(prefers-color-scheme: dark)").matches;
  var theme = storedTheme || (prefersDark ? "dark" : "corporate");
  root.setAttribute("data-theme", theme);
  window.mrnToggleTheme = function() {
    var next = root.getAttribute("data-theme") === "dark" ? "corporate" : "dark";
    root.setAttribute("data-theme", next);
    try {
      localStorage.setItem(key, next);
    } catch (e) {}
  };
})();`

const ChartBootstrapJS = `(function() {
  var charts = new WeakMap();

  function initCharts(root) {
    if (!window.Chart) return;
    var scope = root || document;
    var canvases = scope.querySelectorAll ? scope.querySelectorAll("[data-mrn-chart]") : [];
    canvases.forEach(function(canvas) {
      var container = canvas.closest("[data-mrn-chart-root]");
      if (!container) return;
      var configEl = container.querySelector("[data-mrn-chart-config]");
      if (!configEl) return;

      var config;
      try {
        config = JSON.parse(configEl.textContent || "{}");
      } catch (e) {
        return;
      }

      var existing = charts.get(canvas);
      if (existing) existing.destroy();
      var chart = new window.Chart(canvas, config);
      chart.options = chart.options || {};
      chart.options.onClick = function(_, elements) {
        if (!elements || !elements.length) return;
        var first = elements[0];
        var label = (chart.data && chart.data.labels && chart.data.labels[first.index]) || "";
        var stateName = container.getAttribute("data-mrn-query-state");
        var column = container.getAttribute("data-mrn-filter-column");
        if (!stateName || !column) return;
        var payload = {state: stateName, filters: [{column: column, op: "eq", value: String(label)}]};
        document.dispatchEvent(new CustomEvent("mrn:data-query-change", {detail: payload}));
        if (window.htmx) {
          window.htmx.trigger(document.body, "mrn:data-query-change", payload);
        }
      };
      chart.update();
      charts.set(canvas, chart);
    });
  }

  document.addEventListener("DOMContentLoaded", function() {
    initCharts(document);
  });
  document.addEventListener("htmx:afterSwap", function(event) {
    initCharts(event.detail && event.detail.elt ? event.detail.elt : document);
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
