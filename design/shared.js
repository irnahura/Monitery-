(function () {
  function setMessage(node, text, type) {
    if (!node) return;
    node.textContent = text;
    node.className = "message" + (type ? " " + type : "");
  }

  function validateEmail(value) {
    return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value);
  }

  function initLogin() {
    var form = document.getElementById("login-form");
    if (!form) return;
    form.addEventListener("submit", function (event) {
      event.preventDefault();
      var email = document.getElementById("login-email").value.trim();
      var password = document.getElementById("login-password").value.trim();
      var message = document.getElementById("login-message");
      if (!validateEmail(email)) {
        setMessage(message, "Enter a valid email address.", "error");
        return;
      }
      if (!password) {
        setMessage(message, "Enter your password.", "error");
        return;
      }
      setMessage(message, "Signed in. Redirecting to dashboard...", "success");
      window.setTimeout(function () {
        window.location.href = "dashboard.html";
      }, 500);
    });
  }

  function initRegister() {
    var form = document.getElementById("register-form");
    if (!form) return;
    form.addEventListener("submit", function (event) {
      event.preventDefault();
      var name = document.getElementById("register-name").value.trim();
      var email = document.getElementById("register-email").value.trim();
      var password = document.getElementById("register-password").value;
      var confirm = document.getElementById("register-confirm").value;
      var message = document.getElementById("register-message");
      if (!name) {
        setMessage(message, "Enter your full name.", "error");
        return;
      }
      if (!validateEmail(email)) {
        setMessage(message, "Enter a valid email address.", "error");
        return;
      }
      if (password.length < 8) {
        setMessage(message, "Password must be at least 8 characters.", "error");
        return;
      }
      if (password !== confirm) {
        setMessage(message, "Passwords do not match.", "error");
        return;
      }
      setMessage(message, "Account created. Redirecting to login...", "success");
      window.setTimeout(function () {
        window.location.href = "login.html";
      }, 600);
    });
  }

  function initMonitorForm() {
    var form = document.getElementById("monitor-form");
    if (!form) return;
    form.addEventListener("submit", function (event) {
      event.preventDefault();
      var name = document.getElementById("monitor-name").value.trim();
      var url = document.getElementById("monitor-url").value.trim();
      var message = document.getElementById("monitor-message");
      if (!name) {
        setMessage(message, "Enter a monitor name.", "error");
        return;
      }
      try {
        var parsed = new URL(url);
        if (!/^https?:$/.test(parsed.protocol)) {
          throw new Error("Bad protocol");
        }
      } catch (error) {
        setMessage(message, "Enter a valid http or https URL.", "error");
        return;
      }
      setMessage(message, "Monitor created. Returning to dashboard...", "success");
      window.setTimeout(function () {
        window.location.href = "dashboard.html";
      }, 550);
    });
  }

  function initSettingsTabs() {
    var tabs = Array.prototype.slice.call(document.querySelectorAll(".settings-tab"));
    if (!tabs.length) return;
    var panels = Array.prototype.slice.call(document.querySelectorAll(".settings-panel"));

    function activate(targetId) {
      tabs.forEach(function (tab) {
        tab.classList.toggle("active", tab.getAttribute("data-target") === targetId);
      });
      panels.forEach(function (panel) {
        panel.classList.toggle("active", panel.id === targetId);
      });
      if (window.location.hash !== "#" + targetId) {
        history.replaceState(null, "", "#" + targetId);
      }
    }

    tabs.forEach(function (tab) {
      tab.addEventListener("click", function () {
        activate(tab.getAttribute("data-target"));
      });
    });

    var hash = window.location.hash.replace("#", "");
    if (hash) {
      activate(hash);
    }
  }

  function initInlineForms() {
    Array.prototype.slice.call(document.querySelectorAll(".inline-form")).forEach(function (form) {
      form.addEventListener("submit", function (event) {
        event.preventDefault();
        var message = form.querySelector(".message");
        var inputs = Array.prototype.slice.call(form.querySelectorAll("input[type='text'], input[type='email'], input[type='password']"));
        var empty = inputs.some(function (input) {
          return !input.value.trim();
        });
        if (inputs.length && empty) {
          setMessage(message, "Complete all fields before saving.", "error");
          return;
        }
        if (form.querySelector("#confirm-password")) {
          var newPassword = document.getElementById("new-password").value;
          var confirmPassword = document.getElementById("confirm-password").value;
          if (newPassword !== confirmPassword) {
            setMessage(message, "New password fields must match.", "error");
            return;
          }
        }
        setMessage(message, form.getAttribute("data-save-message") || "Saved.", "success");
      });
    });
  }

  function initApiTable() {
    var apiMessage = document.getElementById("api-message");
    var generateButton = document.getElementById("generate-key");
    var table = document.querySelector(".api-table");
    if (!table) return;

    if (generateButton) {
      generateButton.addEventListener("click", function () {
        var row = document.createElement("div");
        row.className = "table-row static-row api-row";
        row.setAttribute("role", "row");
        row.innerHTML =
          '<div><strong>New generated key</strong><span class="mono">pk_live_new27c</span></div>' +
          '<span class="mono">2026-07-07</span>' +
          '<span class="mono">2026-10-07</span>' +
          '<div class="table-actions">' +
          '<button class="button button-tertiary copy-key" type="button" data-key="pk_live_new27c">Copy key</button>' +
          '<button class="button button-tertiary revoke-key" type="button">Revoke key</button>' +
          "</div>";
        table.appendChild(row);
        setMessage(apiMessage, "New API key generated.", "success");
      });
    }

    table.addEventListener("click", function (event) {
      var target = event.target;
      if (!(target instanceof HTMLElement)) return;
      if (target.classList.contains("copy-key")) {
        var key = target.getAttribute("data-key") || "";
        if (navigator.clipboard && navigator.clipboard.writeText) {
          navigator.clipboard.writeText(key);
        }
        setMessage(apiMessage, "Key copied.", "success");
      }
      if (target.classList.contains("revoke-key")) {
        var row = target.closest(".api-row");
        if (row) {
          row.remove();
          setMessage(apiMessage, "Key revoked.", "success");
        }
      }
    });
  }

  initLogin();
  initRegister();
  initMonitorForm();
  initSettingsTabs();
  initInlineForms();
  initApiTable();
})();
