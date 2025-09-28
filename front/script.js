const API_BASE_URL = 'http://localhost:8081';

function showLoading(show) {
    document.getElementById('loading').style.display = show ? 'block' : 'none';
}

function showError(message) {
    const errorDiv = document.getElementById('error');
    errorDiv.textContent = message;
    errorDiv.style.display = 'block';
    document.getElementById('success').style.display = 'none';
}

function showSuccess(message) {
    const successDiv = document.getElementById('success');
    successDiv.textContent = message;
    successDiv.style.display = 'block';
    document.getElementById('error').style.display = 'none';
}

function hideMessages() {
    document.getElementById('error').style.display = 'none';
    document.getElementById('success').style.display = 'none';
}

function formatDate(dateString) {
    if (!dateString) return 'Не указано';
    try {
        const date = new Date(dateString);
        return date.toLocaleString('ru-RU');
    } catch {
        return dateString;
    }
}

function formatCurrency(amount, currency) {
    if (!amount) return 'Не указано';
    const formatter = new Intl.NumberFormat('ru-RU', {
        style: 'currency',
        currency: currency || 'USD'
    });
    return formatter.format(amount);
}

function searchOrder() {
    const orderId = document.getElementById('orderIdInput').value.trim();
    
    if (!orderId) {
        showError('Пожалуйста, введите ID заказа');
        return;
    }

    hideMessages();
    showLoading(true);
    document.getElementById('resultSection').style.display = 'none';

    fetch(`${API_BASE_URL}/order/${orderId}`)
        .then(response => {
            if (!response.ok) {
                if (response.status === 404) {
                    throw new Error('Заказ не найден');
                } else {
                    throw new Error('Ошибка сервера: ' + response.status);
                }
            }
            return response.json();
        })
        .then(data => {
            displayOrderData(data);
            showSuccess('Информация о заказе успешно загружена!');
        })
        .catch(error => {
            showError(error.message);
        })
        .finally(() => {
            showLoading(false);
        });
}

function displayOrderData(order) {
    const orderDetails = document.getElementById('orderDetails');
    orderDetails.innerHTML = `
        <div class="info-item">
            <span class="info-label">ID заказа</span>
            <span class="info-value">${order.order_uid || 'Не указано'}</span>
        </div>
        <div class="info-item">
            <span class="info-label">Трек номер</span>
            <span class="info-value">${order.track_number || 'Не указано'}</span>
        </div>
        <div class="info-item">
            <span class="info-label">Точка входа</span>
            <span class="info-value">${order.entry || 'Не указано'}</span>
        </div>
        <div class="info-item">
            <span class="info-label">Локаль</span>
            <span class="info-value">${order.locale || 'Не указано'}</span>
        </div>
        <div class="info-item">
            <span class="info-label">ID клиента</span>
            <span class="info-value">${order.customer_id || 'Не указано'}</span>
        </div>
        <div class="info-item">
            <span class="info-label">Служба доставки</span>
            <span class="info-value">${order.delivery_service || 'Не указано'}</span>
        </div>
        <div class="info-item">
            <span class="info-label">Дата создания</span>
            <span class="info-value">${formatDate(order.date_created)}</span>
        </div>
    `;
    const deliveryInfo = document.getElementById('deliveryInfo');
    if (order.delivery) {
        deliveryInfo.innerHTML = `
            <div class="info-item">
                <span class="info-label">Имя</span>
                <span class="info-value">${order.delivery.name || 'Не указано'}</span>
            </div>
            <div class="info-item">
                <span class="info-label">Телефон</span>
                <span class="info-value">${order.delivery.phone || 'Не указано'}</span>
            </div>
            <div class="info-item">
                <span class="info-label">Email</span>
                <span class="info-value">${order.delivery.email || 'Не указано'}</span>
            </div>
            <div class="info-item">
                <span class="info-label">Адрес</span>
                <span class="info-value">${order.delivery.address || 'Не указано'}</span>
            </div>
            <div class="info-item">
                <span class="info-label">Город</span>
                <span class="info-value">${order.delivery.city || 'Не указано'}</span>
            </div>
            <div class="info-item">
                <span class="info-label">Регион</span>
                <span class="info-value">${order.delivery.region || 'Не указано'}</span>
            </div>
            <div class="info-item">
                <span class="info-label">Почтовый индекс</span>
                <span class="info-value">${order.delivery.zip || 'Не указано'}</span>
            </div>
        `;
    } else {
        deliveryInfo.innerHTML = '<div class="info-item"><span class="info-value">Информация о доставке отсутствует</span></div>';
    }
    const paymentInfo = document.getElementById('paymentInfo');
    if (order.payment) {
        paymentInfo.innerHTML = `
            <div class="info-item">
                <span class="info-label">Транзакция</span>
                <span class="info-value">${order.payment.transaction || 'Не указано'}</span>
            </div>
            <div class="info-item">
                <span class="info-label">Сумма</span>
                <span class="info-value">${formatCurrency(order.payment.amount, order.payment.currency)}</span>
            </div>
            <div class="info-item">
                <span class="info-label">Валюта</span>
                <span class="info-value">${order.payment.currency || 'Не указано'}</span>
            </div>
            <div class="info-item">
                <span class="info-label">Провайдер</span>
                <span class="info-value">${order.payment.provider || 'Не указано'}</span>
            </div>
            <div class="info-item">
                <span class="info-label">Банк</span>
                <span class="info-value">${order.payment.bank || 'Не указано'}</span>
            </div>
            <div class="info-item">
                <span class="info-label">Стоимость доставки</span>
                <span class="info-value">${formatCurrency(order.payment.delivery_cost, order.payment.currency)}</span>
            </div>
            <div class="info-item">
                <span class="info-label">Стоимость товаров</span>
                <span class="info-value">${formatCurrency(order.payment.goods_total, order.payment.currency)}</span>
            </div>
            <div class="info-item">
                <span class="info-label">Дата оплаты</span>
                <span class="info-value">${formatDate(new Date(order.payment.payment_dt * 1000))}</span>
            </div>
        `;
    } else {
        paymentInfo.innerHTML = '<div class="info-item"><span class="info-value">Информация об оплате отсутствует</span></div>';
    }
    const itemsBody = document.getElementById('itemsBody');
    if (order.items && order.items.length > 0) {
        itemsBody.innerHTML = order.items.map(item => `
            <tr>
                <td>${item.name || 'Не указано'}</td>
                <td>${item.brand || 'Не указано'}</td>
                <td>${formatCurrency(item.price, order.payment?.currency)}</td>
                <td>${item.quantity || Math.round(item.total_price / item.price) || 1}</td>
                <td>${formatCurrency(item.total_price, order.payment?.currency)}</td>
                <td>${item.status || 'Не указано'}</td>
            </tr>
        `).join('');
    } else {
        itemsBody.innerHTML = '<tr><td colspan="6" style="text-align: center;">Товары не найдены</td></tr>';
    }
    document.getElementById('resultSection').style.display = 'block';
}

document.getElementById('orderIdInput').addEventListener('keypress', function(e) {
    if (e.key === 'Enter') {
        searchOrder();
    }
});

document.getElementById('orderIdInput').focus();