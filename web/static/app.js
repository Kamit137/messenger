const form = document.querySelector("#auth-form");
const title = document.querySelector("#form-title");
const submit = form.querySelector("button");
const toggle = document.querySelector("#toggle-mode");
const message = document.querySelector("#message");
let mode = "login";

async function showCurrentUser() {
  const response = await fetch("/api/me");
  if (!response.ok) return;
  const user = await response.json();
  form.hidden = true;
  toggle.hidden = true;
  title.textContent = `Вы вошли как ${user.email}`;
  message.textContent = "Регистрация и вход готовы. Следующим шагом подключим чат по WebSocket.";
}

toggle.addEventListener("click", () => {
  mode = mode === "login" ? "register" : "login";
  title.textContent = mode === "login" ? "Войти" : "Создать аккаунт";
  submit.textContent = mode === "login" ? "Войти" : "Зарегистрироваться";
  toggle.textContent = mode === "login" ? "Нет аккаунта? Зарегистрироваться" : "Уже есть аккаунт? Войти";
  message.textContent = "";
});

form.addEventListener("submit", async (event) => {
  event.preventDefault();
  message.textContent = "";
  const values = new FormData(form);
  const response = await fetch(`/api/${mode}`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email: values.get("email"), password: values.get("password") }),
  });
  const body = await response.json();
  if (!response.ok) {
    message.textContent = body.error || "Что-то пошло не так";
    return;
  }
  if (mode === "register") {
    message.textContent = "Аккаунт создан. Теперь войди с этим email и паролем.";
    toggle.click();
    return;
  }
  await showCurrentUser();
});

showCurrentUser();
