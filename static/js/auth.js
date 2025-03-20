/**
 * Обработчик авторизации пользователей и отображения данных пользователя
 */
document.addEventListener("DOMContentLoaded", function () {
  // Проверяем, есть ли данные пользователя в localStorage
  const userDataJson = localStorage.getItem("user_data");
  const authToken = localStorage.getItem("auth_token");

  if (userDataJson && authToken) {
    try {
      const userData = JSON.parse(userDataJson);

      // Получаем элемент из DOM, который показывает, авторизован ли пользователь
      const authStatusElement = document.getElementById("auth-status");
      const isAuthenticated = authStatusElement
        ? authStatusElement.dataset.auth === "true"
        : false;

      // Если на странице нет данных пользователя, но они есть в localStorage,
      // обновляем интерфейс
      if (!isAuthenticated && userData) {
        // Обновляем навигационное меню
        const authSection = document.querySelector(".d-flex");
        if (authSection) {
          let menuItems = "";

          if (userData.role === "provider") {
            menuItems += `
              <li><a class="dropdown-item" href="/provider/services">Мои услуги</a></li>
              <li><a class="dropdown-item" href="/provider/orders">Заказы</a></li>
            `;
          } else if (userData.role === "client") {
            menuItems += `
              <li><a class="dropdown-item" href="/client/orders">Мои заказы</a></li>
            `;
          }

          authSection.innerHTML = `
            <div class="dropdown">
              <button class="btn btn-outline-light dropdown-toggle" type="button" data-bs-toggle="dropdown">
                ${userData.username || userData.email}
              </button>
              <ul class="dropdown-menu dropdown-menu-end">
                <li><a class="dropdown-item" href="/profile">Профиль</a></li>
                ${menuItems}
                <li><hr class="dropdown-divider"></li>
                <li><a class="dropdown-item text-danger" href="#" id="logoutBtn">Выйти</a></li>
              </ul>
            </div>
          `;

          // Добавляем функциональность для кнопки выхода
          const logoutBtn = document.getElementById("logoutBtn");
          if (logoutBtn) {
            logoutBtn.addEventListener("click", function (e) {
              e.preventDefault();

              // Отправляем запрос на сервер для выхода
              axios
                .post("/user/logout")
                .then(function () {
                  // Очищаем локальные данные
                  localStorage.removeItem("auth_token");
                  localStorage.removeItem("user_data");

                  // Перезагружаем страницу после успешного выхода
                  window.location.reload();
                })
                .catch(function (error) {
                  console.error("Ошибка при выходе из системы:", error);
                  // Даже в случае ошибки все равно выполняем локальный выход
                  localStorage.removeItem("auth_token");
                  localStorage.removeItem("user_data");
                  window.location.reload();
                });
            });
          }
        }
      }
    } catch (e) {
      console.error("Ошибка при обработке данных пользователя:", e);
    }
  }

  // Находим все кнопки выхода (для уже авторизованных пользователей)
  const allLogoutBtns = document.querySelectorAll("#logoutBtn");
  allLogoutBtns.forEach((btn) => {
    if (btn) {
      btn.addEventListener("click", function (e) {
        e.preventDefault();

        // Отправляем запрос на сервер для выхода
        axios
          .post("/user/logout")
          .then(function () {
            // Очищаем локальные данные
            localStorage.removeItem("auth_token");
            localStorage.removeItem("user_data");

            // Перезагружаем страницу после успешного выхода
            window.location.reload();
          })
          .catch(function (error) {
            console.error("Ошибка при выходе из системы:", error);
            // Даже в случае ошибки все равно выполняем локальный выход
            localStorage.removeItem("auth_token");
            localStorage.removeItem("user_data");
            window.location.reload();
          });
      });
    }
  });
});
