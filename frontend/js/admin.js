// frontend/js/admin.js
document.addEventListener('DOMContentLoaded', () => {
    // --- Глобальные переменные и константы ---
    const API_BASE_URL = ''; // Тот же домен
    let currentAdminEstablishment = null; 
    let loggedInAdminName = '';
    let currentSelectedOwnerId = null; 
    let currentEditingHall = null;     
    let currentEditingPlaceId = null;  

    // --- Элементы DOM для "фейкового" логина ---
    const adminLoginSection = document.getElementById('adminLoginSection');
    const adminLoginError = document.getElementById('adminLoginError');
    const loginAsVistaAdminBtn = document.getElementById('loginAsVistaAdminBtn');
    const loginAsMostHubAdminBtn = document.getElementById('loginAsMostHubAdminBtn');
    const logoutAdminBtn = document.getElementById('logoutAdminBtn');

    // --- Элементы DOM для основной информации о заведении ---
    const establishmentManagementSection = document.getElementById('establishmentManagementSection');
    const adminWelcomeMessage = document.getElementById('adminWelcomeMessage');
    const currentEstablishmentIdStoreEl = document.getElementById('currentEstablishmentIdStore');
    const currentOwnerIdStoreEl = document.getElementById('currentOwnerIdStore');
    const adminEstNameInput = document.getElementById('adminEstNameInput');
    const adminEstTypeDisplay = document.getElementById('adminEstTypeDisplay');
    const adminEstAddressInput = document.getElementById('adminEstAddressInput');
    const adminEstHoursInput = document.getElementById('adminEstHoursInput');
    const adminEstDescriptionInput = document.getElementById('adminEstDescriptionInput');
    const saveEstablishmentChangesBtn = document.getElementById('saveEstablishmentChangesBtn');
    const adminEstUpdateStatus = document.getElementById('adminEstUpdateStatus');

    // --- Элементы DOM для управления залами ---
    const hallsListContainerEl = document.getElementById('hallsListContainer');
    const adminAddHallForm = document.getElementById('adminAddHallForm');
    const newHallNameInput = document.getElementById('newHallName');
    const newHallDescriptionInput = document.getElementById('newHallDescription');
    const newHallCapacityInput = document.getElementById('newHallCapacity');
    const newHallAirConditionerCheckbox = document.getElementById('newHallAirConditioner');
    const adminAddHallStatus = document.getElementById('adminAddHallStatus');

    // --- Элементы DOM для модального окна редактирования зала ---
    const editHallModal = document.getElementById('editHallModal');
    const closeEditHallModalBtn = document.getElementById('closeEditHallModalBtn');
    const editHallForm = document.getElementById('editHallForm');
    const editHallIdInput = document.getElementById('editHallId');
    const editHallNameInput = document.getElementById('editHallName');
    const editHallDescriptionInput = document.getElementById('editHallDescription');
    const editHallCapacityInput = document.getElementById('editHallCapacity');
    const editHallAirConditionerCheckbox = document.getElementById('editHallAirConditioner');
    const editHallStatus = document.getElementById('editHallStatus');

    // --- Элементы DOM для управления меню (пока только отображение) ---
    const menuManagementAdmin = document.querySelector('.menu-management-admin');
    const adminMenuList = document.getElementById('adminMenuList');

    // --- Элементы DOM для УПРАВЛЕНИЯ МЕСТАМИ ---
    const placeManagementSection = document.getElementById('placeManagementSection');
    const placeManagementTitle = document.getElementById('placeManagementTitle');
    const placeManagementHallName = document.getElementById('placeManagementHallName');
    const backToHallsBtn = document.getElementById('backToHallsBtn');
    const adminHallMapPreview = document.getElementById('adminHallMapPreview');
    const adminPlaceFormTitle = document.getElementById('adminPlaceFormTitle');
    const adminPlaceForm = document.getElementById('adminPlaceForm');
    const editingPlaceIdInput = document.getElementById('editingPlaceId'); 
    const placeNameInput = document.getElementById('placeName');
    const placeTypeInput = document.getElementById('placeType');
    const placeCoordXInput = document.getElementById('placeCoordX');
    const placeCoordYInput = document.getElementById('placeCoordY');
    const savePlaceBtn = document.getElementById('savePlaceBtn');
    const cancelEditPlaceBtn = document.getElementById('cancelEditPlaceBtn');
    const adminPlaceFormStatus = document.getElementById('adminPlaceFormStatus');
    const placesListContainerAdminEl = document.getElementById('placesListContainerAdmin');
    const currentEditingHallNameForList = document.getElementById('currentEditingHallNameForList');

    // НОВЫЕ ИКОНКИ ДЛЯ АДМИНКИ
    const ICON_PLACE_FREE_ADMIN = 'https://i.ibb.co/ch1cN4hZ/image.png'; 
    const ICON_PLACE_BOOKED_ADMIN = 'https://i.ibb.co/mVKyLxRQ/free-icon-armchair-2168247.png';


    // --- Вспомогательные функции ---
    function displayMessage(element, message, isSuccess) {
        if (element) {
            element.textContent = message;
            element.className = isSuccess ? 'status-message success' : 'status-message error';
            element.style.display = 'block';
            setTimeout(() => { element.style.display = 'none'; }, 4000);
        }
    }

    function openModal(modalElement) {
        if (modalElement) modalElement.style.display = 'flex';
    }
    function closeModal(modalElement) {
        if (modalElement) modalElement.style.display = 'none';
    }
    function sanitizeClassName(name) {
        if (!name) return '';
        return name.toLowerCase().replace(/[^a-z0-9_]+/g, '-').replace(/^-+|-+$/g, '');
    }

    // --- Логика "Фейкового" Логина и Отображения ---
    function checkLoginStatusMVP() { 
        currentSelectedOwnerId = sessionStorage.getItem('minityAdminMVPOwnerID');
        loggedInAdminName = sessionStorage.getItem('minityAdminMVPName');
        if (currentSelectedOwnerId && loggedInAdminName) { fetchAdminEstablishmentDataMVP(currentSelectedOwnerId); showMainManagementSection(loggedInAdminName); } 
        else { showLoginSection(); }
    }
    function showLoginSection() { 
        if (adminLoginSection) adminLoginSection.style.display = 'block';
        if (establishmentManagementSection) establishmentManagementSection.style.display = 'none';
        if (placeManagementSection) placeManagementSection.style.display = 'none'; 
        if (logoutAdminBtn) logoutAdminBtn.style.display = 'none';
        currentAdminEstablishment = null; 
        if (currentEstablishmentIdStoreEl) currentEstablishmentIdStoreEl.value = '';
        if (currentOwnerIdStoreEl) currentOwnerIdStoreEl.value = '';
    }
    function showMainManagementSection(adminName) { 
        if (adminLoginSection) adminLoginSection.style.display = 'none';
        if (establishmentManagementSection) establishmentManagementSection.style.display = 'block';
        if (placeManagementSection) placeManagementSection.style.display = 'none'; 
        if (logoutAdminBtn) logoutAdminBtn.style.display = 'inline-block';
        if (adminWelcomeMessage) { adminWelcomeMessage.textContent = `Добро пожаловать, ${adminName}! Панель управления:`; }
    }
    async function fetchAdminEstablishmentDataMVP(ownerId, callback) { 
        if (!ownerId) { displayMessage(adminLoginError, "Не выбран администратор.", false); showLoginSection(); return; }
        displayMessage(adminEstUpdateStatus, 'Загрузка данных...', false); 
        try {
            const response = await fetch(`${API_BASE_URL}/api/v1/admin/my-establishment?owner_id=${ownerId}`, { method: 'GET' });
            adminEstUpdateStatus.style.display = 'none'; 
            if (response.ok) {
                const data = await response.json(); currentAdminEstablishment = data; 
                if (currentEstablishmentIdStoreEl && data && data.id) currentEstablishmentIdStoreEl.value = data.id;
                if (currentOwnerIdStoreEl && data && data.owner_user_id) currentOwnerIdStoreEl.value = data.owner_user_id;
                console.log(`Est data for owner_id ${ownerId}:`, currentAdminEstablishment);
                displayEstablishmentData(currentAdminEstablishment);
                if (callback) callback(); 
            } else {
                const errorData = await response.json(); console.error(`Ошибка загрузки для owner_id ${ownerId}:`, errorData.error || response.statusText);
                displayMessage(adminLoginError, `Не удалось загрузить: ${errorData.error || response.statusText}`, false); showLoginSection(); 
            }
        } catch (error) {
            adminEstUpdateStatus.style.display = 'none'; console.error(`Workspace error for owner_id ${ownerId}:`, error);
            displayMessage(adminLoginError, 'Сетевая ошибка.', false); showLoginSection();
        }
    }
    function displayEstablishmentData(establishment) { 
        if (!establishment) {
            if (adminEstNameInput) adminEstNameInput.value = ''; if (adminEstTypeDisplay) adminEstTypeDisplay.textContent = '';
            if (adminEstAddressInput) adminEstAddressInput.value = ''; if (adminEstHoursInput) adminEstHoursInput.value = '';
            if (adminEstDescriptionInput) adminEstDescriptionInput.value = ''; if (menuManagementAdmin) menuManagementAdmin.style.display = 'none';
            if (hallsListContainerEl) hallsListContainerEl.innerHTML = '<li>Данные не загружены.</li>';
            if (adminMenuList) adminMenuList.innerHTML = ''; return;
        }
        if (adminEstNameInput) adminEstNameInput.value = establishment.name || '';
        if (adminEstTypeDisplay) adminEstTypeDisplay.textContent = establishment.type === 'restaurant' ? 'Ресторан' : (establishment.type === 'coworking' ? 'Коворкинг' : 'Не указан');
        if (adminEstAddressInput) adminEstAddressInput.value = establishment.address || '';
        if (adminEstHoursInput) adminEstHoursInput.value = establishment.working_hours || '';
        if (adminEstDescriptionInput) adminEstDescriptionInput.value = establishment.description || '';
        if (menuManagementAdmin) { menuManagementAdmin.style.display = establishment.type === 'restaurant' ? 'block' : 'none';}
        renderHallsList(establishment.halls, establishment.id, establishment.owner_user_id);
        if (adminMenuList && establishment.menu && establishment.type === 'restaurant') { renderAdminMenu(establishment.menu);
        } else if (adminMenuList && establishment.type === 'restaurant') { adminMenuList.innerHTML = '<h4>Позиции меню:</h4><p>Меню отсутствует.</p>';}
    }
    function handleAdminChoice(ownerId, adminName) { 
        if (adminLoginError) adminLoginError.style.display = 'none'; loggedInAdminName = adminName; currentSelectedOwnerId = ownerId;
        sessionStorage.setItem('minityAdminMVPOwnerID', ownerId); sessionStorage.setItem('minityAdminMVPName', adminName);
        fetchAdminEstablishmentDataMVP(ownerId); showMainManagementSection(adminName);
    }
    if (loginAsVistaAdminBtn) { loginAsVistaAdminBtn.addEventListener('click', () => { handleAdminChoice(loginAsVistaAdminBtn.dataset.ownerId, loginAsVistaAdminBtn.dataset.adminName); }); }
    if (loginAsMostHubAdminBtn) { loginAsMostHubAdminBtn.addEventListener('click', () => { handleAdminChoice(loginAsMostHubAdminBtn.dataset.ownerId, loginAsMostHubAdminBtn.dataset.adminName); }); }
    if (logoutAdminBtn) { logoutAdminBtn.addEventListener('click', () => { sessionStorage.removeItem('minityAdminMVPOwnerID'); sessionStorage.removeItem('minityAdminMVPName'); loggedInAdminName = ''; currentSelectedOwnerId = null; displayEstablishmentData(null); showLoginSection(); }); }
    if (saveEstablishmentChangesBtn) { 
        saveEstablishmentChangesBtn.addEventListener('click', async () => {
            const estIdToUpdate = currentEstablishmentIdStoreEl ? currentEstablishmentIdStoreEl.value : null; const ownerIdForUpdate = currentOwnerIdStoreEl ? currentOwnerIdStoreEl.value : null;
            if (!estIdToUpdate || !ownerIdForUpdate) { displayMessage(adminEstUpdateStatus, 'Ошибка: ID не определены.', false); return; }
            const updateData = {};
            if(currentAdminEstablishment && adminEstNameInput.value.trim() !== currentAdminEstablishment.name) updateData.name = adminEstNameInput.value.trim();
            if(currentAdminEstablishment && adminEstAddressInput.value.trim() !== currentAdminEstablishment.address) updateData.address = adminEstAddressInput.value.trim();
            if(currentAdminEstablishment && adminEstHoursInput.value.trim() !== currentAdminEstablishment.working_hours) updateData.working_hours = adminEstHoursInput.value.trim();
            if(currentAdminEstablishment && adminEstDescriptionInput.value.trim() !== currentAdminEstablishment.description) updateData.description = adminEstDescriptionInput.value.trim();
            if (Object.keys(updateData).length === 0) { displayMessage(adminEstUpdateStatus, 'Нет изменений.', true); return; }
            if (updateData.name !== undefined && !updateData.name) { alert("Название не может быть пустым."); adminEstNameInput.focus(); return; }
            displayMessage(adminEstUpdateStatus, 'Сохранение...', false); 
            try {
                const response = await fetch(`${API_BASE_URL}/api/v1/admin/my-establishment?establishment_id=${estIdToUpdate}&owner_id=${ownerIdForUpdate}`, { method: 'PUT', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(updateData) });
                if (response.ok) { const updatedEst = await response.json(); currentAdminEstablishment = updatedEst; displayEstablishmentData(currentAdminEstablishment); displayMessage(adminEstUpdateStatus, 'Сохранено!', true);
                } else { const errorData = await response.json(); displayMessage(adminEstUpdateStatus, `Ошибка: ${errorData.error || response.statusText}`, false); }
            } catch (error) { console.error('Update est error:', error); displayMessage(adminEstUpdateStatus, 'Сетевая ошибка.', false); }
        });
    }

    // --- Управление Залами (отображение, добавление, редактирование, удаление) ---
    function renderHallsList(halls, establishmentId, ownerId) { 
        if (!hallsListContainerEl) return; hallsListContainerEl.innerHTML = ''; 
        if (!halls || halls.length === 0) { const li = document.createElement('li'); li.textContent = 'Залы не добавлены.'; hallsListContainerEl.appendChild(li); return; }
        halls.forEach(hall => {
            const li = document.createElement('li'); li.className = 'hall-list-item';
            li.innerHTML = `<div class="hall-info"><strong>${hall.name}</strong> (ID: ${hall.id})<br><small>Описание: ${hall.description || 'нет'}, Вместимость: ${hall.capacity || '?'}, Кондиционер: ${hall.has_air_conditioner ? 'Да' : 'Нет'}</small></div><div class="hall-actions"><button class="btn btn-sm btn-edit-hall" data-hall-id="${hall.id}">Ред.</button><button class="btn btn-sm btn-delete-hall" data-hall-id="${hall.id}">Удал.</button><button class="btn btn-sm btn-manage-places" data-hall-id="${hall.id}" data-hall-name="${hall.name}">Места</button></div>`;
            const deleteBtn = li.querySelector('.btn-delete-hall'); if (deleteBtn) { deleteBtn.addEventListener('click', () => handleDeleteHallClick(establishmentId, hall.id, ownerId, hall.name)); }
            const editBtn = li.querySelector('.btn-edit-hall'); if (editBtn) { editBtn.addEventListener('click', () => openEditHallModal(hall)); } 
            const managePlacesBtn = li.querySelector('.btn-manage-places'); if (managePlacesBtn) { managePlacesBtn.addEventListener('click', () => openPlaceManagement(hall)); } 
            hallsListContainerEl.appendChild(li);
        });
    }
    async function handleDeleteHallClick(establishmentId, hallId, ownerId, hallName) { 
        if (!establishmentId || !hallId || !ownerId) { alert('Ошибка: не хватает данных.'); return; }
        if (confirm(`Удалить зал "${hallName}" (ID: ${hallId})?`)) {
            displayMessage(adminAddHallStatus, `Удаление зала...`, false); 
            try {
                const response = await fetch(`${API_BASE_URL}/api/v1/admin/establishments/${establishmentId}/halls/${hallId}?owner_id=${ownerId}`, {method: 'DELETE'});
                if (response.ok) { displayMessage(adminAddHallStatus, `Зал "${hallName}" удален.`, true); if (currentSelectedOwnerId) { fetchAdminEstablishmentDataMVP(currentSelectedOwnerId); }}
                else { const errorData = await response.json(); displayMessage(adminAddHallStatus, `Ошибка удаления: ${errorData.error || response.statusText}`, false); }
            } catch (error) { console.error('Delete hall err:', error); displayMessage(adminAddHallStatus, 'Сетевая ошибка.', false); }
        }
    }
    if (adminAddHallForm) { 
        adminAddHallForm.addEventListener('submit', async (e) => {
            e.preventDefault(); displayMessage(adminAddHallStatus, '', true); 
            const establishmentId = currentEstablishmentIdStoreEl ? currentEstablishmentIdStoreEl.value : null; const ownerId = currentOwnerIdStoreEl ? currentOwnerIdStoreEl.value : null;
            if (!establishmentId || !ownerId) { displayMessage(adminAddHallStatus, 'Ошибка: ID не определены.', false); return; }
            const hallData = { name: newHallNameInput.value.trim(), description: newHallDescriptionInput.value.trim() || null, capacity: parseInt(newHallCapacityInput.value) || 0, has_air_conditioner: newHallAirConditionerCheckbox.checked };
            if (!hallData.name) { alert("Название зала не может быть пустым."); newHallNameInput.focus(); return; }
            displayMessage(adminAddHallStatus, 'Добавление зала...', false);
            try {
                const response = await fetch(`${API_BASE_URL}/api/v1/admin/establishments/${establishmentId}/halls?owner_id=${ownerId}`, { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(hallData) });
                if (response.ok) { displayMessage(adminAddHallStatus, 'Зал добавлен!', true); adminAddHallForm.reset(); if (currentSelectedOwnerId) { fetchAdminEstablishmentDataMVP(currentSelectedOwnerId); }}
                else { const errorData = await response.json(); displayMessage(adminAddHallStatus, `Ошибка: ${errorData.error || response.statusText}`, false); }
            } catch (error) { console.error('Add hall err:', error); displayMessage(adminAddHallStatus, 'Сетевая ошибка.', false); }
        });
    }
    function openEditHallModal(hall) { 
        if (!hall || !editHallModal || !editHallForm) return;
        editHallIdInput.value = hall.id; editHallNameInput.value = hall.name || ''; editHallDescriptionInput.value = hall.description || '';
        editHallCapacityInput.value = hall.capacity || 0; editHallAirConditionerCheckbox.checked = hall.has_air_conditioner || false;
        displayMessage(editHallStatus, '', true); openModal(editHallModal);
    }
    if (closeEditHallModalBtn) { closeEditHallModalBtn.addEventListener('click', () => closeModal(editHallModal)); }
    if (editHallModal) { editHallModal.addEventListener('click', (event) => { if (event.target === editHallModal) { closeModal(editHallModal); } }); }
    if (editHallForm) { 
        editHallForm.addEventListener('submit', async (e) => {
            e.preventDefault(); displayMessage(editHallStatus, '', true);
            const hallId = editHallIdInput.value; const establishmentId = currentEstablishmentIdStoreEl ? currentEstablishmentIdStoreEl.value : null; const ownerId = currentOwnerIdStoreEl ? currentOwnerIdStoreEl.value : null;
            if (!hallId || !establishmentId || !ownerId) { displayMessage(editHallStatus, 'Ошибка: ID не определены.', false); return; }
            const hallUpdateData = { name: editHallNameInput.value.trim(), description: editHallDescriptionInput.value.trim() || null, capacity: parseInt(editHallCapacityInput.value) || 0, has_air_conditioner: editHallAirConditionerCheckbox.checked };
            if (!hallUpdateData.name) { alert("Название зала не может быть пустым."); editHallNameInput.focus(); return; }
            displayMessage(editHallStatus, 'Сохранение...', false);
            try {
                const response = await fetch(`${API_BASE_URL}/api/v1/admin/establishments/${establishmentId}/halls/${hallId}?owner_id=${ownerId}`, { method: 'PUT', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(hallUpdateData) });
                if (response.ok) { displayMessage(editHallStatus, 'Зал обновлен!', true); closeModal(editHallModal); if (currentSelectedOwnerId) { fetchAdminEstablishmentDataMVP(currentSelectedOwnerId); }}
                else { const errorData = await response.json(); displayMessage(editHallStatus, `Ошибка: ${errorData.error || response.statusText}`, false); }
            } catch (error) { console.error('Update hall err:', error); displayMessage(editHallStatus, 'Сетевая ошибка.', false); }
        });
    }

    // --- ЛОГИКА ДЛЯ УПРАВЛЕНИЯ МЕСТАМИ ---
    function openPlaceManagement(hall) {
        if (!hall || !placeManagementSection || !establishmentManagementSection) return;
        console.log("Управление местами для зала:", hall);
        currentEditingHall = hall; 

        if (placeManagementHallName) placeManagementHallName.textContent = hall.name;
        if (currentEditingHallNameForList) currentEditingHallNameForList.textContent = hall.name;

        establishmentManagementSection.style.display = 'none'; 
        placeManagementSection.style.display = 'block';      

        renderAdminHallMapPreview(hall.places || []); 
        renderAdminPlacesList(hall.places || []);   
        resetAdminPlaceForm(); 
    }

    if (backToHallsBtn) {
        backToHallsBtn.addEventListener('click', () => {
            if (placeManagementSection) placeManagementSection.style.display = 'none';
            if (establishmentManagementSection) establishmentManagementSection.style.display = 'block';
            currentEditingHall = null; 
            currentEditingPlaceId = null; 
        });
    }

    function renderAdminHallMapPreview(places) {
        if (!adminHallMapPreview) return;
        adminHallMapPreview.innerHTML = ''; 

        if (!places || places.length === 0) {
            const placeholder = document.createElement('p');
            placeholder.className = 'hall-map-placeholder';
            placeholder.textContent = 'В этом зале еще нет мест. Добавьте их с помощью формы.';
            adminHallMapPreview.appendChild(placeholder);
            return;
        }

        places.forEach(place => {
            const placeEl = document.createElement('div');
            placeEl.className = 'place';
            placeEl.dataset.placeId = place.id;
            if (place.type) placeEl.classList.add(`place-type-${sanitizeClassName(place.type)}`);
            
            let x = '10', y = '10'; 
            if (place.visual_info) {
                try {
                    const vi = JSON.parse(place.visual_info);
                    if (vi.x !== undefined) x = vi.x;
                    if (vi.y !== undefined) y = vi.y;
                } catch (e) { console.warn("Bad visual_info for place", place.id, place.visual_info); }
            }
            placeEl.style.top = `${y}%`;
            placeEl.style.left = `${x}%`;
            placeEl.title = `Место ${place.name} (ID: ${place.id})\nТип: ${place.type || 'не указан'}\nСтатус: ${place.status}`;

            const iconImg = document.createElement('img');
            iconImg.className = 'place-icon';
            // ИСПОЛЬЗУЕМ НОВЫЕ ИКОНКИ
            iconImg.src = (place.status === 'booked' || place.status === 'occupied') ? ICON_PLACE_BOOKED_ADMIN : ICON_PLACE_FREE_ADMIN;
            iconImg.alt = place.status;
            
            const placeLabel = document.createElement('span');
            placeLabel.className = 'place-label';
            placeLabel.textContent = place.name; 

            placeEl.appendChild(iconImg);
            placeEl.appendChild(placeLabel);
            
            placeEl.addEventListener('click', () => populatePlaceFormForEdit(place));
            adminHallMapPreview.appendChild(placeEl);
        });
    }

    function renderAdminPlacesList(places) {
        if (!placesListContainerAdminEl) return;
        placesListContainerAdminEl.innerHTML = '';

        if (!places || places.length === 0) {
            const li = document.createElement('li');
            li.textContent = 'Места в этом зале еще не добавлены.';
            placesListContainerAdminEl.appendChild(li);
            return;
        }

        places.forEach(place => {
            const li = document.createElement('li');
            li.className = 'place-list-item-admin';
            let visualX = '?', visualY = '?';
            try {
                if(place.visual_info) {
                    const vi = JSON.parse(place.visual_info);
                    visualX = vi.x !== undefined ? vi.x : '?';
                    visualY = vi.y !== undefined ? vi.y : '?';
                }
            } catch(e){}

            li.innerHTML = `
                <span>ID: ${place.id}, Имя: <strong>${place.name}</strong>, Тип: ${place.type || '-'}, X: ${visualX}, Y: ${visualY}, Статус: ${place.status}</span>
                <div class="actions">
                    <button class="btn btn-sm btn-edit-place" data-place-id="${place.id}">Ред.</button>
                    <button class="btn btn-sm btn-delete-place" data-place-id="${place.id}">Удал.</button>
                </div>
            `;
            li.querySelector('.btn-edit-place').addEventListener('click', () => populatePlaceFormForEdit(place));
            li.querySelector('.btn-delete-place').addEventListener('click', () => handleDeletePlaceClick(place));
            placesListContainerAdminEl.appendChild(li);
        });
    }

    function resetAdminPlaceForm() {
        if (!adminPlaceForm) return;
        adminPlaceForm.reset();
        if(editingPlaceIdInput) editingPlaceIdInput.value = ''; 
        if(adminPlaceFormTitle) adminPlaceFormTitle.textContent = 'Добавить новое место';
        if(savePlaceBtn) savePlaceBtn.textContent = 'Добавить место';
        if(cancelEditPlaceBtn) cancelEditPlaceBtn.style.display = 'none';
        if(adminHallMapPreview) { 
            adminHallMapPreview.querySelectorAll('.place.selected-for-edit').forEach(p => p.classList.remove('selected-for-edit'));
        }
        currentEditingPlaceId = null;
    }

    function populatePlaceFormForEdit(place) {
        if (!adminPlaceForm || !place) return;
        resetAdminPlaceForm(); 

        currentEditingPlaceId = place.id;
        if(editingPlaceIdInput) editingPlaceIdInput.value = place.id;
        if(placeNameInput) placeNameInput.value = place.name || '';
        if(placeTypeInput) placeTypeInput.value = place.type || '';
        
        let x = '', y = '';
        if (place.visual_info) {
            try {
                const vi = JSON.parse(place.visual_info);
                if (vi.x !== undefined) x = vi.x;
                if (vi.y !== undefined) y = vi.y;
            } catch (e) { console.warn("Bad visual_info for place", place.id, place.visual_info); }
        }
        if(placeCoordXInput) placeCoordXInput.value = x;
        if(placeCoordYInput) placeCoordYInput.value = y;

        if(adminPlaceFormTitle) adminPlaceFormTitle.textContent = `Редактировать место: ${place.name} (ID: ${place.id})`;
        if(savePlaceBtn) savePlaceBtn.textContent = 'Сохранить изменения';
        if(cancelEditPlaceBtn) cancelEditPlaceBtn.style.display = 'inline-block';

        if(adminHallMapPreview) {
            adminHallMapPreview.querySelectorAll('.place.selected-for-edit').forEach(p => p.classList.remove('selected-for-edit'));
            const placeElOnMap = adminHallMapPreview.querySelector(`.place[data-place-id="${place.id}"]`);
            if (placeElOnMap) placeElOnMap.classList.add('selected-for-edit');
        }
        placeNameInput.focus();
    }
    
    if (cancelEditPlaceBtn) {
        cancelEditPlaceBtn.addEventListener('click', resetAdminPlaceForm);
    }

    if (adminPlaceForm) {
        adminPlaceForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            if (!currentEditingHall || !currentAdminEstablishment) {
                displayMessage(adminPlaceFormStatus, 'Ошибка: Зал или заведение не определены.', false); return;
            }
            const establishmentId = currentAdminEstablishment.id;
            const hallId = currentEditingHall.id;
            const ownerId = currentAdminEstablishment.owner_user_id;

            const placeDataForRequest = {
                name: placeNameInput.value.trim(),
                type: placeTypeInput.value.trim() || undefined, // Отправляем undefined, если пусто, чтобы не перезаписать на пустую строку на бэке, если поле опционально
                visual_info: JSON.stringify({ 
                    x: parseInt(placeCoordXInput.value) || 0, 
                    y: parseInt(placeCoordYInput.value) || 0 
                })
            };
            // Для PUT запроса, если поле не изменилось, лучше его не отправлять (особенно если оно опционально на бэке)
            // Но для простоты MVP мы отправляем все поля, которые есть в форме. Бэкенд разберется.
            // Если поле status будет, то:
            // if (currentEditingPlaceId && someStatusSelect.value) placeDataForRequest.status = someStatusSelect.value;


            if (!placeDataForRequest.name) { alert("Имя/номер места не может быть пустым."); placeNameInput.focus(); return; }
            if (isNaN(parseInt(placeCoordXInput.value)) || isNaN(parseInt(placeCoordYInput.value))) {
                alert("Координаты X и Y должны быть числами."); return;
            }

            let url, method;
            const placeIdBeingEdited = editingPlaceIdInput.value; // Берем из скрытого поля

            if (placeIdBeingEdited) { // Редактирование
                url = `${API_BASE_URL}/api/v1/admin/establishments/${establishmentId}/halls/${hallId}/places/${placeIdBeingEdited}?owner_id=${ownerId}`;
                method = 'PUT';
                displayMessage(adminPlaceFormStatus, 'Обновление места...', false);
            } else { // Добавление
                url = `${API_BASE_URL}/api/v1/admin/establishments/${establishmentId}/halls/${hallId}/places?owner_id=${ownerId}`;
                method = 'POST';
                displayMessage(adminPlaceFormStatus, 'Добавление места...', false);
            }

            try {
                const response = await fetch(url, {
                    method: method, headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(placeDataForRequest)
                });
                if (response.ok) {
                    displayMessage(adminPlaceFormStatus, `Место успешно ${placeIdBeingEdited ? 'обновлено' : 'добавлено'}!`, true);
                    resetAdminPlaceForm();
                    if (currentSelectedOwnerId) {
                        fetchAdminEstablishmentDataMVP(currentSelectedOwnerId, () => {
                            if (placeManagementSection.style.display === 'block' && currentAdminEstablishment && currentAdminEstablishment.halls) {
                                const updatedHall = currentAdminEstablishment.halls.find(h => h.id === hallId);
                                if (updatedHall) {
                                    currentEditingHall = updatedHall; 
                                    renderAdminHallMapPreview(updatedHall.places || []);
                                    renderAdminPlacesList(updatedHall.places || []);
                                }
                            }
                        });
                    }
                } else {
                    const errorData = await response.json();
                    displayMessage(adminPlaceFormStatus, `Ошибка: ${errorData.error || response.statusText}`, false);
                }
            } catch (error) {
                console.error('Save place error:', error);
                displayMessage(adminPlaceFormStatus, 'Сетевая ошибка при сохранении места.', false);
            }
        });
    }
    
    async function handleDeletePlaceClick(place) {
        if (!currentEditingHall || !currentAdminEstablishment || !place) {
            alert('Ошибка: не хватает данных для удаления места.'); return;
        }
        const establishmentId = currentAdminEstablishment.id;
        const hallId = currentEditingHall.id;
        const ownerId = currentAdminEstablishment.owner_user_id;
        const placeId = place.id;

        if (confirm(`Вы уверены, что хотите удалить место "${place.name}" (ID: ${placeId})?`)) {
            displayMessage(adminPlaceFormStatus, `Удаление места "${place.name}"...`, false);
            try {
                const response = await fetch(`${API_BASE_URL}/api/v1/admin/establishments/${establishmentId}/halls/${hallId}/places/${placeId}?owner_id=${ownerId}`, {
                    method: 'DELETE'
                });
                if (response.ok) {
                    displayMessage(adminPlaceFormStatus, `Место "${place.name}" успешно удалено.`, true);
                    if (currentSelectedOwnerId) {
                         fetchAdminEstablishmentDataMVP(currentSelectedOwnerId, () => {
                            if (placeManagementSection.style.display === 'block' && currentAdminEstablishment && currentAdminEstablishment.halls) {
                                const updatedHall = currentAdminEstablishment.halls.find(h => h.id === hallId);
                                if (updatedHall) {
                                    currentEditingHall = updatedHall;
                                    renderAdminHallMapPreview(updatedHall.places || []);
                                    renderAdminPlacesList(updatedHall.places || []);
                                } else { 
                                    placeManagementSection.style.display = 'none';
                                    if(establishmentManagementSection) establishmentManagementSection.style.display = 'block';
                                }
                            }
                        });
                    }
                } else {
                    const errorData = await response.json();
                    displayMessage(adminPlaceFormStatus, `Ошибка удаления места: ${errorData.error || response.statusText}`, false);
                }
            } catch (error) {
                console.error('Delete place error:', error);
                displayMessage(adminPlaceFormStatus, 'Сетевая ошибка при удалении места.', false);
            }
        }
    }

    // --- Инициализация и общие обработчики ---
    checkLoginStatusMVP();
    const currentYearSpan = document.getElementById('current-year'); if (currentYearSpan) { currentYearSpan.textContent = new Date().getFullYear(); }
    const themeToggleButtonGlobal = document.getElementById('theme-toggle-btn');
    if (themeToggleButtonGlobal) {
        const applyTheme = (theme) => { if (theme === 'dark') { document.body.classList.remove('light-theme'); document.body.classList.add('dark-theme');} else { document.body.classList.remove('dark-theme'); document.body.classList.add('light-theme');}};
        const savedTheme = localStorage.getItem('theme'); applyTheme(savedTheme || 'light'); 
        themeToggleButtonGlobal.addEventListener('click', () => { const newTheme = document.body.classList.contains('dark-theme') ? 'light' : 'dark'; applyTheme(newTheme); localStorage.setItem('theme', newTheme); });
    }
});