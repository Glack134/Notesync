document.addEventListener('DOMContentLoaded', function() {
    // Автоматически заполняем поля для тестирования
    document.getElementById('username').value = "polyыуlwыk0s05";
    document.getElementById('password').value = "q2weerцыыty";
});

// Остальной код остается прежним
document.getElementById('registration-form').addEventListener('submit', function(event) {
    event.preventDefault(); // Предотвращаем отправку формы

    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;

    // Здесь вы можете отправить данные на сервер
    fetch('/auth/sign-up', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ username, password }),
    })
    .then(response => {
        if (response.ok) {
            alert('Регистрация успешна!');
            document.body.style.transition = "background-color 0.5s";
            document.body.style.backgroundColor = "#d4edda"; // Зеленый фон для успешной регистрации
            setTimeout(() => {
                window.location.href = '/main'; // Перенаправление на главную страницу
            }, 1000);
        } else {
            alert('Ошибка регистрации. Пожалуйста, попробуйте еще раз.');
        }
    })
    .catch(error => {
        console.error('Ошибка:', error);
        alert('Произошла ошибка. Пожалуйста, попробуйте еще раз.');
    });
});
