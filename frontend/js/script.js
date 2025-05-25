document.addEventListener('DOMContentLoaded', () => {
    // Переключение темы
    const themeToggleBtn = document.getElementById('theme-toggle-btn');
    themeToggleBtn.addEventListener('click', () => {
        document.body.classList.toggle('dark-theme');
        document.body.classList.toggle('light-theme');
        // Можно сохранять выбор темы в localStorage
        localStorage.setItem('theme', document.body.classList.contains('dark-theme') ? 'dark' : 'light');
    });

    // Загрузка сохраненной темы
    const savedTheme = localStorage.getItem('theme');
    if (savedTheme === 'dark') {
        document.body.classList.add('dark-theme');
        document.body.classList.remove('light-theme');
    } else {
        document.body.classList.add('light-theme');
        document.body.classList.remove('dark-theme');
    }


    // Модальное окно авторизации/регистрации
    const authModal = document.getElementById('auth-modal');
    const openModalBtnHero = document.getElementById('book-now-hero-btn');
    const openModalBtnHeader = document.getElementById('login-register-btn-header');
    const closeModalBtn = document.querySelector('.modal .close-btn');

    function openModal() {
        if(authModal) authModal.style.display = 'block';
    }
    function closeModal() {
        if(authModal) authModal.style.display = 'none';
    }

    if(openModalBtnHero) openModalBtnHero.addEventListener('click', openModal);
    if(openModalBtnHeader) openModalBtnHeader.addEventListener('click', openModal);
    if(closeModalBtn) closeModalBtn.addEventListener('click', closeModal);

    window.addEventListener('click', (event) => {
        if (event.target == authModal) {
            closeModal();
        }
    });

    // Переключение форм в модальном окне
    const loginForm = document.getElementById('login-form');
    const registerForm = document.getElementById('register-form');
    const showRegisterLink = document.getElementById('show-register-form-link');
    const showLoginLink = document.getElementById('show-login-form-link');

    if(showRegisterLink) {
        showRegisterLink.addEventListener('click', (e) => {
            e.preventDefault();
            if(loginForm) loginForm.style.display = 'none';
            if(registerForm) registerForm.style.display = 'block';
        });
    }

    if(showLoginLink) {
        showLoginLink.addEventListener('click', (e) => {
            e.preventDefault();
            if(registerForm) registerForm.style.display = 'none';
            if(loginForm) loginForm.style.display = 'block';
        });
    }

    // Обработка формы поиска бронирования (пример)
    const bookingSearchForm = document.getElementById('booking-search-form');
    if (bookingSearchForm) {
        bookingSearchForm.addEventListener('submit', async (event) => {
            event.preventDefault();
            const formData = new FormData(bookingSearchForm);
            const establishmentType = formData.get('establishment_type');
            const bookingTime = formData.get('booking_time');
            const peopleCount = formData.get('people_count');

            console.log('Поиск:', { establishmentType, bookingTime, peopleCount });
            // Здесь будет логика отправки запроса на бэкенд и отображения результатов
            // Например: window.location.href = `/search-results?type=${establishmentType}&time=${bookingTime}...`;
            alert('Функционал поиска в разработке!');
        });
    }

    // Обработка формы быстройго поиска в шапке
    const quickSearchFormHeader = document.getElementById('quick-search-form-header');
    if (quickSearchFormHeader) {
        quickSearchFormHeader.addEventListener('submit', (event) => {
            event.preventDefault();
            // Логика аналогична основной форме поиска
            alert('Быстрый поиск в разработке!');
        });
    }

    // Обработка форм входа и регистрации (пример)
    if (loginForm) {
        loginForm.addEventListener('submit', async (event) => {
            event.preventDefault();
            // TODO: Валидация данных
            const formData = new FormData(loginForm);
            const phone = formData.get('phone');
            const password = formData.get('password');

            try {
                const response = await fetch('/api/v1/auth/login', { // Пример URL API
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ phone, password })
                });
                const data = await response.json();
                if (response.ok) {
                    alert('Вход успешен! Токен: ' + data.token); // Сохранить токен
                    closeModal();
                    // Обновить UI для авторизованного пользователя
                } else {
                    alert('Ошибка входа: ' + (data.error || 'Неверные данные'));
                }
            } catch (error) {
                console.error('Ошибка при входе:', error);
                alert('Произошла ошибка сети.');
            }
        });
    }

    if (registerForm) {
        registerForm.addEventListener('submit', async (event) => {
            event.preventDefault();
            // TODO: Валидация данных
            const formData = new FormData(registerForm);
            const name = formData.get('name');
            const phone = formData.get('phone');
            const email = formData.get('email');
            const password = formData.get('password');

             try {
                const response = await fetch('/api/v1/auth/register', { // Пример URL API
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ name, phone, email, password })
                });
                const data = await response.json();
                if (response.ok) {
                    alert('Регистрация успешна! ' + (data.message || 'Теперь вы можете войти.'));
                    // Опционально: запросить подтверждение по SMS/Email
                    loginForm.style.display = 'block'; // Показать форму входа
                    registerForm.style.display = 'none';
                } else {
                    alert('Ошибка регистрации: ' + (data.error || 'Попробуйте снова'));
                }
            } catch (error) {
                console.error('Ошибка при регистрации:', error);
                alert('Произошла ошибка сети.');
            }
        });
    }

    // Установка текущего года в подвале
    const currentYearSpan = document.getElementById('current-year');
    if (currentYearSpan) {
        currentYearSpan.textContent = new Date().getFullYear();
    }

    // Анимации при наведении (пример)
    document.querySelectorAll('.btn, .partner-item, .review-card').forEach(element => {
        element.addEventListener('mouseenter', () => {
            // element.style.transform = 'translateY(-2px)';
            // element.style.boxShadow = '0 4px 8px rgba(0,0,0,0.15)';
        });
        element.addEventListener('mouseleave', () => {
            // element.style.transform = 'translateY(0)';
            // element.style.boxShadow = '0 2px 5px rgba(0,0,0,0.1)'; // Вернуть исходную тень
        });
    });

});