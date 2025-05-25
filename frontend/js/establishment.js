// frontend/js/establishment.js
document.addEventListener('DOMContentLoaded', () => {
    const establishmentIdParam = new URLSearchParams(window.location.search).get('id');
    const establishmentNameTitle = document.querySelector('title');

    const establishmentImageEl = document.getElementById('establishmentImage');
    const establishmentNameEl = document.getElementById('establishmentName');
    const establishmentTypeBadgeEl = document.getElementById('establishmentTypeBadge');
    const establishmentDescriptionEl = document.getElementById('establishmentDescription');
    const establishmentAddressEl = document.getElementById('establishmentAddress');
    const establishmentHoursEl = document.getElementById('establishmentHours');
    const establishmentPhoneEl = document.getElementById('establishmentPhone');

    const tabLinks = document.querySelectorAll('.tab-link');
    const tabPanes = document.querySelectorAll('.tab-pane');
    const menuTabButton = document.getElementById('menuTabButton');

    const hallSelect = document.getElementById('hallSelect');
    const currentHallMapEl = document.getElementById('currentHallMap');
    const hallMapPlaceholder = document.querySelector('.hall-map-placeholder');
    
    const bookingDateTimeEl = document.getElementById('bookingDateTime');
    const bookingPeopleCountEl = document.getElementById('bookingPeopleCount');
    const selectedPlacesCountEl = document.getElementById('selectedPlacesCount');
    const selectedPlacesListEl = document.getElementById('selectedPlacesList');
    const confirmBookingBtn = document.getElementById('confirmBookingBtn');

    let currentEstablishmentData = null; 
    let currentHallData = null; 
    let selectedPlaces = [];
    
    const API_BASE_URL_PUBLIC = ''; 

    // ИКОНКИ (будут установлены в loadEstablishmentData)
    let ICON_PLACE_FREE = 'https://i.ibb.co/RkjbtWSC/2-1.png'; // Старая дефолтная "свободно"
    let ICON_PLACE_BOOKED = 'https://i.ibb.co/QvZW3Hrh/1-1.png'; // Старая дефолтная "занято"

    // Новые иконки по типам
    const RESTAURANT_ICON_FREE = 'https://i.ibb.co/ch1cN4hZ/image.png';
    const RESTAURANT_ICON_BOOKED = 'https://i.ibb.co/HpzjmX9D/free-icon-dinning-table-1606240.png';
    const COWORKING_ICON_FREE = 'https://i.ibb.co/Y7Ft2mXW/free-icon-workspace-6135610.png';
    const COWORKING_ICON_BOOKED = 'https://i.ibb.co/V0LcfrbM/free-icon-coworking-6148084.png';

    // --- Функции ---

    async function loadEstablishmentData(id) {
        if (!id) {
            showErrorMessage('ID заведения не указан в URL.');
            return;
        }
        try {
            const response = await fetch(`${API_BASE_URL_PUBLIC}/api/v1/establishments/${id}`);
            
            if (!response.ok) {
                let errorMsg = `Ошибка загрузки данных заведения: ${response.statusText}`;
                try { const errorData = await response.json(); errorMsg = `Ошибка: ${errorData.error || response.statusText}`; } catch (e) {}
                showErrorMessage(errorMsg);
                return;
            }

            currentEstablishmentData = await response.json();
            console.log('Fetched establishment data (public):', currentEstablishmentData);

            if (!currentEstablishmentData) {
                showErrorMessage('Данные о заведении не получены или некорректны.');
                return;
            }
            
            // УСТАНОВКА ПРАВИЛЬНЫХ ИКОНОК В ЗАВИСИМОСТИ ОТ ТИПА ЗАВЕДЕНИЯ
            if (currentEstablishmentData.type === 'restaurant') {
                ICON_PLACE_FREE = RESTAURANT_ICON_FREE;
                ICON_PLACE_BOOKED = RESTAURANT_ICON_BOOKED;
                console.log("Using RESTAURANT icons");
            } else if (currentEstablishmentData.type === 'coworking') {
                ICON_PLACE_FREE = COWORKING_ICON_FREE;
                ICON_PLACE_BOOKED = COWORKING_ICON_BOOKED;
                console.log("Using COWORKING icons");
            } else { 
                // Если тип не определен или другой, используем дефолтные (старые)
                ICON_PLACE_FREE = 'https://i.ibb.co/RkjbtWSC/2-1.png'; 
                ICON_PLACE_BOOKED = 'https://i.ibb.co/QvZW3Hrh/1-1.png';
                console.log("Using DEFAULT icons");
            }

            updatePageWithData();

        } catch (error) {
            console.error('Сетевая ошибка или ошибка парсинга при загрузке данных заведения:', error);
            showErrorMessage('Не удалось загрузить информацию о заведении. Проверьте ваше интернет-соединение.');
        }
    }

    function showErrorMessage(message) { 
        const container = document.querySelector('.establishment-page .container');
        if (container) {
            let mainContentParent = container;
            const introSection = container.querySelector('section.establishment-intro');
            const contentSection = container.querySelector('section.establishment-content');
            if(introSection && contentSection) { // Если есть обе основные секции, чистим их родителя
                mainContentParent = introSection.parentNode;
            }
            mainContentParent.innerHTML = ''; // Очищаем
            const errorCard = document.createElement('div');
            errorCard.className = 'card';
            errorCard.style.padding = '20px';
            errorCard.style.textAlign = 'center';
            errorCard.innerHTML = `<p class="error-message">${message} Пожалуйста, вернитесь на <a href="index.html">главную страницу</a>.</p>`;
            mainContentParent.appendChild(errorCard);
        }
        if(establishmentNameTitle) establishmentNameTitle.textContent = "Ошибка - Minity";
    }
    
    function updatePageWithData() {
        if (!currentEstablishmentData) return;
        if(establishmentNameTitle) establishmentNameTitle.textContent = `${currentEstablishmentData.name} - Minity`;
        if(establishmentImageEl) {
            establishmentImageEl.src = currentEstablishmentData.photos && currentEstablishmentData.photos.length > 0 
                                       ? currentEstablishmentData.photos[0] 
                                       : 'https://via.placeholder.com/450x300.png?text=No+Image';
            establishmentImageEl.alt = currentEstablishmentData.name;
        }
        if(establishmentNameEl) establishmentNameEl.textContent = currentEstablishmentData.name;
        if(establishmentTypeBadgeEl) establishmentTypeBadgeEl.textContent = currentEstablishmentData.type === 'restaurant' ? 'Ресторан' : 'Коворкинг';
        if(establishmentDescriptionEl) establishmentDescriptionEl.textContent = currentEstablishmentData.description;
        if(establishmentAddressEl) establishmentAddressEl.textContent = currentEstablishmentData.address;
        if(establishmentHoursEl) establishmentHoursEl.textContent = currentEstablishmentData.working_hours; 
        if(establishmentPhoneEl) establishmentPhoneEl.textContent = currentEstablishmentData.phone || 'Не указан'; 
        if (currentEstablishmentData.type === 'restaurant' && currentEstablishmentData.menu && menuTabButton) {
            menuTabButton.style.display = 'inline-block'; renderMenu(currentEstablishmentData.menu);
        } else if (menuTabButton) { menuTabButton.style.display = 'none';}
        if (hallSelect && currentEstablishmentData.halls) {
            hallSelect.innerHTML = '<option value="">-- Выберите зал --</option>';
            currentEstablishmentData.halls.forEach(hall => { const option = document.createElement('option'); option.value = hall.id; option.textContent = hall.name; hallSelect.appendChild(option); });
            if (currentEstablishmentData.halls.length > 0) { hallSelect.value = String(currentEstablishmentData.halls[0].id); renderHallMap(currentEstablishmentData.halls[0].id);
            } else { renderHallMap(null); }
        } else if (hallSelect) { hallSelect.innerHTML = '<option value="">Залы не найдены</option>'; renderHallMap(null); }
    }

    function renderMenu(menuData) { 
        const menuCategoriesEl = document.getElementById('menuCategories'); if (!menuCategoriesEl) return; menuCategoriesEl.innerHTML = ''; 
        menuData.forEach(category => {
            const categoryDiv = document.createElement('div'); categoryDiv.className = 'menu-category';
            const categoryTitle = document.createElement('h3'); categoryTitle.textContent = category.category; categoryDiv.appendChild(categoryTitle);
            category.items.forEach(item => { const itemDiv = document.createElement('div'); itemDiv.className = 'menu-item';
                itemDiv.innerHTML = `<div class="menu-item-info"><h4>${item.name}</h4>${item.description ? `<p class="menu-item-description">${item.description}</p>` : ''}</div><p class="menu-item-price">${item.price} ₽</p>`; 
                categoryDiv.appendChild(itemDiv);
            });
            menuCategoriesEl.appendChild(categoryDiv);
        });
    }

    function sanitizeClassName(name) { if (!name) return ''; return name.toLowerCase().replace(/[^a-z0-9_]+/g, '-').replace(/^-+|-+$/g, ''); }

    function renderHallMap(hallIdToRender) {
        const hallId = parseInt(hallIdToRender); 
        if (!currentEstablishmentData || !currentEstablishmentData.halls || !currentHallMapEl) { if(hallMapPlaceholder) hallMapPlaceholder.textContent = 'Ошибка загрузки.'; if(currentHallMapEl) currentHallMapEl.innerHTML = ''; if(currentHallMapEl && hallMapPlaceholder) currentHallMapEl.appendChild(hallMapPlaceholder); return; }
        currentHallMapEl.innerHTML = ''; selectedPlaces = []; updateSelectedPlacesInfo();
        currentHallData = currentEstablishmentData.halls.find(h => h.id === hallId);
        if (!currentHallData || !currentHallData.places || currentHallData.places.length === 0) { if(hallMapPlaceholder) { currentHallMapEl.appendChild(hallMapPlaceholder); hallMapPlaceholder.textContent = hallId ? 'Мест нет.' : 'Выберите зал.';} return; }
        if(hallMapPlaceholder && hallMapPlaceholder.parentNode === currentHallMapEl) { currentHallMapEl.removeChild(hallMapPlaceholder); }

        currentHallData.places.forEach(place => {
            const placeEl = document.createElement('div'); placeEl.className = 'place';
            placeEl.dataset.placeId = String(place.id); 
            if (place.type) { placeEl.classList.add(`place-type-${sanitizeClassName(place.type)}`); }
            placeEl.title = `${place.type || 'Место'} (ID: ${place.id}) - ${place.status === 'free' ? 'свободно' : 'занято'}`;
            let xPos = '10%', yPos = '10%'; 
            if (place.visual_info) { try { const vi = JSON.parse(place.visual_info); if (vi.x !== undefined) xPos = `${vi.x}%`; if (vi.y !== undefined) yPos = `${vi.y}%`; } catch (e) {} } 
            else { xPos = `${Math.random() * 80 + 5}%`; yPos = `${Math.random() * 80 + 5}%`; }
            placeEl.style.top = yPos; placeEl.style.left = xPos;
            const iconImg = document.createElement('img'); iconImg.className = 'place-icon';
            const placeLabel = document.createElement('span'); placeLabel.className = 'place-label'; placeLabel.textContent = place.name; 

            // Иконки ICON_PLACE_FREE и ICON_PLACE_BOOKED теперь устанавливаются в loadEstablishmentData
            if (place.status === 'booked' || place.status === 'occupied') { 
                iconImg.src = ICON_PLACE_BOOKED; 
                iconImg.alt = 'Занято'; placeEl.classList.add('place-booked');
            } else { 
                iconImg.src = ICON_PLACE_FREE; 
                iconImg.alt = 'Свободно'; placeEl.classList.add('place-free');
                if (place.status !== 'unavailable') { placeEl.addEventListener('click', handlePlaceClick); } 
                else { placeEl.classList.add('place-unavailable'); placeEl.title += ' (недоступно)'; }
            }
            placeEl.appendChild(iconImg); placeEl.appendChild(placeLabel); currentHallMapEl.appendChild(placeEl);
        });
    }

    function handlePlaceClick(event) { 
        const placeEl = event.currentTarget; const placeId = placeEl.dataset.placeId;
        if (placeEl.classList.contains('place-booked') || placeEl.classList.contains('place-unavailable')) return; 
        const maxSelectInput = bookingPeopleCountEl ? parseInt(bookingPeopleCountEl.value) : 1;
        const maxSelect = isNaN(maxSelectInput) || maxSelectInput < 1 ? 1 : maxSelectInput;
        const isSelected = placeEl.classList.contains('selected');
        if (!isSelected && selectedPlaces.length >= maxSelect) { alert(`Можно выбрать не более ${maxSelect} ${maxSelect === 1 ? 'места' : 'мест'}.`); return; }
        placeEl.classList.toggle('selected');
        if (placeEl.classList.contains('selected')) { selectedPlaces.push(placeId); } else { selectedPlaces = selectedPlaces.filter(id => id !== placeId); }
        updateSelectedPlacesInfo();
    }

    function updateSelectedPlacesInfo() {
        if(selectedPlacesCountEl) selectedPlacesCountEl.textContent = selectedPlaces.length;
        if(selectedPlacesListEl) { if (selectedPlaces.length > 0) { selectedPlacesListEl.innerHTML = `Выбраны: <strong>${selectedPlaces.join(', ')}</strong>`; } else { selectedPlacesListEl.innerHTML = 'Места не выбраны'; }}
        if(confirmBookingBtn) { const dateTime = bookingDateTimeEl ? bookingDateTimeEl.value : ''; const peopleCountInput = bookingPeopleCountEl ? bookingPeopleCountEl.value : '0'; const peopleCount = parseInt(peopleCountInput); confirmBookingBtn.disabled = selectedPlaces.length === 0 || !dateTime || isNaN(peopleCount) || peopleCount < 1; }
    }

    tabLinks.forEach(link => { link.addEventListener('click', (e) => { e.preventDefault(); const tabId = link.dataset.tab; tabLinks.forEach(item => item.classList.remove('active')); tabPanes.forEach(pane => pane.classList.remove('active')); link.classList.add('active'); const activePane = document.getElementById(tabId); if (activePane) { activePane.classList.add('active'); }}); });
    if (hallSelect) { hallSelect.addEventListener('change', (event) => { const selectedHallId = event.target.value; if (selectedHallId) { renderHallMap(parseInt(selectedHallId)); } else { renderHallMap(null); }}); }
    if (bookingDateTimeEl) bookingDateTimeEl.addEventListener('change', updateSelectedPlacesInfo);
    if (bookingPeopleCountEl) bookingPeopleCountEl.addEventListener('input', updateSelectedPlacesInfo);
    if (confirmBookingBtn) { 
        confirmBookingBtn.addEventListener('click', () => {
            if (selectedPlaces.length === 0) { alert('Выберите место.'); return; } const dateTime = bookingDateTimeEl.value; const peopleCount = bookingPeopleCountEl.value;
            if (!dateTime) { alert('Выберите дату и время.'); bookingDateTimeEl.focus(); return; } const peopleCountNum = parseInt(peopleCount);
            if (isNaN(peopleCountNum) || peopleCountNum < 1) { alert('Укажите кол-во человек.'); bookingPeopleCountEl.focus(); return; }
            if (selectedPlaces.length > peopleCountNum) { alert('Мест больше, чем кол-во человек.'); return; }
            const bookingData = { establishmentId: currentEstablishmentData.id, hallId: currentHallData ? currentHallData.id : null, placeIds: selectedPlaces.map(id => parseInt(id.replace(/[^0-9]/g, ''))), dateTime: dateTime, peopleCount: peopleCountNum };
            console.log('Бронирование:', bookingData);
            setTimeout(() => {
                alert(`(ДЕМО) Забронировано! \nЗаведение: ${currentEstablishmentData.name} \nЗал: ${currentHallData ? currentHallData.name : 'N/A'} \nМеста: ${selectedPlaces.join(', ')} \nВремя: ${new Date(dateTime).toLocaleString()} \nКол-во человек: ${peopleCount}`);
                if (currentEstablishmentData && currentHallData && currentHallData.places) { selectedPlaces.forEach(bookedPlaceId => { const placeInMock = currentHallData.places.find(p => String(p.id) === bookedPlaceId); if (placeInMock) { placeInMock.status = 'booked'; }}); }
                renderHallMap(currentHallData ? currentHallData.id : null);
            }, 500);
        });
    }

    if (establishmentIdParam) { loadEstablishmentData(establishmentIdParam); } else { showErrorMessage('ID заведения не указан.'); }
    const currentYearSpan = document.getElementById('current-year'); if (currentYearSpan) { currentYearSpan.textContent = new Date().getFullYear(); }
    const themeToggleButton = document.getElementById('theme-toggle-btn');
    if (themeToggleButton) { const applyTheme = (theme) => { if (theme === 'dark') { document.body.classList.remove('light-theme'); document.body.classList.add('dark-theme');} else { document.body.classList.remove('dark-theme'); document.body.classList.add('light-theme');}}; const savedTheme = localStorage.getItem('theme'); applyTheme(savedTheme || 'light'); themeToggleButton.addEventListener('click', () => { const newTheme = document.body.classList.contains('dark-theme') ? 'light' : 'dark'; applyTheme(newTheme); localStorage.setItem('theme', newTheme); }); }
    updateSelectedPlacesInfo();
});