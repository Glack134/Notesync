document.getElementById('change-greeting').addEventListener('click', function() {
    const greetings = [
        "Добро пожаловать на наш сайт!",
        "Рады видеть вас!",
        "Надеемся, вам понравится!",
        "Приветствуем вас!",
        "Спасибо, что заглянули!"
    ];
    
    // Выбор случайного приветствия
    const randomGreeting = greetings[Math.floor(Math.random() * greetings.length)];
    document.getElementById('greeting').innerText = randomGreeting;
});
