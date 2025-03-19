// JavaScript для клиентских функций

// Автоматическое закрытие алертов через 5 секунд
document.addEventListener("DOMContentLoaded", function () {
  // Находим все алерты на странице
  const alerts = document.querySelectorAll(".alert");

  // Устанавливаем таймер для каждого алерта
  alerts.forEach(function (alert) {
    setTimeout(function () {
      // Создаем экземпляр Bootstrap Alert и вызываем метод close
      const bsAlert = new bootstrap.Alert(alert);
      bsAlert.close();
    }, 5000);
  });

  // Обработка форм с валидацией
  const forms = document.querySelectorAll(".needs-validation");

  Array.from(forms).forEach(function (form) {
    form.addEventListener(
      "submit",
      function (event) {
        if (!form.checkValidity()) {
          event.preventDefault();
          event.stopPropagation();
        }
        form.classList.add("was-validated");
      },
      false
    );
  });

  // Инициализация всплывающих подсказок
  const tooltipTriggerList = [].slice.call(
    document.querySelectorAll('[data-bs-toggle="tooltip"]')
  );
  tooltipTriggerList.map(function (tooltipTriggerEl) {
    return new bootstrap.Tooltip(tooltipTriggerEl);
  });
});

// Функция для подтверждения действий (например, удаления)
function confirmAction(message) {
  return confirm(message || "Вы уверены, что хотите выполнить это действие?");
}

// Функция для обновления статуса заказа
function updateOrderStatus(orderId, status) {
  // Проверка подтверждения
  if (!confirmAction("Вы уверены, что хотите изменить статус заказа?")) {
    return false;
  }

  // Отправляем AJAX запрос
  fetch(`/client/orders/${orderId}/status`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ status: status }),
  })
    .then((response) => {
      if (!response.ok) {
        throw new Error("Ошибка при обновлении статуса");
      }
      return response.json();
    })
    .then((data) => {
      // Перезагружаем страницу после успешного обновления
      window.location.reload();
    })
    .catch((error) => {
      alert("Произошла ошибка: " + error.message);
    });

  return false;
}

// Функция для принятия заказа провайдером
function acceptOrder(orderId) {
  // Получаем данные из формы
  const serviceId = document.getElementById("service_id").value;
  const providerMessage = document.getElementById("provider_message").value;

  // Валидация
  if (!serviceId) {
    alert("Выберите услугу");
    return false;
  }

  if (!providerMessage) {
    alert("Введите сообщение");
    return false;
  }

  // Отправляем AJAX запрос
  fetch(`/provider/orders/${orderId}/accept`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      service_id: parseInt(serviceId),
      provider_message: providerMessage,
    }),
  })
    .then((response) => {
      if (!response.ok) {
        throw new Error("Ошибка при принятии заказа");
      }
      return response.json();
    })
    .then((data) => {
      // Перенаправляем на страницу заказов после успешного принятия
      window.location.href = "/provider/orders";
    })
    .catch((error) => {
      alert("Произошла ошибка: " + error.message);
    });

  return false;
}
