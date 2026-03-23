
const API = {
    BASE_URL: 'http://localhost:5080',
    ENDPOINTS: {
        PRODUCTS:   '/products',
        CATEGORIES: '/categories',
        CHECKOUT:   '/orders',   
    }
};

let state = {
    products: [],
    filtered: [],
    cart: {}       
};

async function request(method, endpoint, body = null) {
    const url = API.BASE_URL + endpoint;
    const headers = { 'Content-Type': 'application/json' };

    const opts = { method, headers };
    if (body) opts.body = JSON.stringify(body);

    const res = await fetch(url, opts);
    const data = await res.json().catch(() => ({}));

    if (!res.ok) {
        throw new Error(data.message || data.error || `HTTP ${res.status}`);
    }
    return data;
}

async function fetchProducts() {
    showLoadingState();
    try {
        const data = await request('GET', API.ENDPOINTS.PRODUCTS);
        state.products = data.products || [];
        populateCategoryFilter();
        applyFilterAndSort();
    } catch (err) {
        console.error('fetchProducts error:', err);
        showToast('ไม่สามารถโหลดสินค้าได้: ' + err.message, 'error');
        renderProducts([]);
    }
}

async function checkout() {
    const items = Object.entries(state.cart).map(([id, qty]) => {
        const product = state.products.find(p => p.id == id);
        return {
            productId: parseInt(id),
            quantity:  qty,
            price:     product?.price ?? 0
        };
    });
    const totalAmount = items.reduce((sum, item) => sum + item.price * item.quantity, 0);

    setCheckoutLoading(true);
    try {
        await request('POST', API.ENDPOINTS.CHECKOUT, { items, totalAmount });

        // อัปเดต stock ฝั่ง frontend ให้ตรงทันที
        items.forEach(({ productId, quantity }) => {
            const p = state.products.find(pr => pr.id == productId);
            if (p && p.stock !== undefined) p.stock = Math.max(0, p.stock - quantity);
        });

        showSuccessModal(items, totalAmount);
        clearCart();
        applyFilterAndSort(); // re-render เพื่อแสดง stock ใหม่

    } catch (err) {
        showToast('ชำระเงินไม่สำเร็จ: ' + err.message, 'error');
    } finally {
        setCheckoutLoading(false);
    }
}

function addToCart(productId) {
    const product = state.products.find(p => p.id == productId);
    if (!product) return;

    const currentQty = state.cart[productId] || 0;
    const maxQty = product.stock ?? 999;

    if (product.stock !== undefined && product.stock <= 0) {
        showToast('สินค้าหมดแล้ว', 'error');
        return;
    }
    if (currentQty >= maxQty) {
        showToast(`สินค้าเหลือเพียง ${maxQty} ชิ้น`, 'info');
        return;
    }

    state.cart[productId] = currentQty + 1;
    updateCartUI();
    renderCardControls(productId);
    showToast(`เพิ่ม "${product.name}" ลงตะกร้าแล้ว`, 'success');
}

function removeFromCart(productId) {
    const qty = state.cart[productId] || 0;
    if (qty <= 1) delete state.cart[productId];
    else state.cart[productId] = qty - 1;
    updateCartUI();
    renderCardControls(productId);
}

function setCartQty(productId, qty) {
    const product = state.products.find(p => p.id == productId);
    const maxQty = product?.stock ?? 999;
    const newQty = Math.max(0, Math.min(qty, maxQty));
    if (newQty === 0) delete state.cart[productId];
    else state.cart[productId] = newQty;
    updateCartUI();
    renderCardControls(productId);
}

function clearCart() {
    state.cart = {};
    updateCartUI();
    state.products.forEach(p => renderCardControls(p.id));
}

function getTotalQty()   { return Object.values(state.cart).reduce((s, q) => s + q, 0); }
function getTotalPrice() {
    return Object.entries(state.cart).reduce((sum, [id, qty]) => {
        const p = state.products.find(pr => pr.id == id);
        return sum + (p?.price ?? 0) * qty;
    }, 0);
}

function renderProducts(products) {
    const grid  = document.getElementById('product-list');
    const empty = document.getElementById('emptyState');

    if (products.length === 0) {
        grid.innerHTML = '';
        empty.classList.remove('hidden');
        return;
    }
    empty.classList.add('hidden');

    grid.innerHTML = products.map((p, i) => `
        <div class="product-card ${p.stock === 0 ? 'out-of-stock' : ''}"
             id="card-${p.id}"
             style="animation-delay: ${i * 0.04}s">

            <div class="card-image-placeholder">◫</div>

            <div class="card-body">
                ${p.category ? `<span class="card-category">${escHtml(p.category)}</span>` : ''}
                <div class="card-name">${escHtml(p.name)}</div>
            </div>

            <div class="card-footer">
                <span class="card-price">฿${formatPrice(p.price)}</span>
                ${p.stock !== undefined
                    ? `<span class="card-stock ${p.stock === 0 ? 'out' : p.stock < 5 ? 'low' : ''}">
                        ${p.stock === 0 ? 'หมด' : `เหลือ ${p.stock}`}
                       </span>`
                    : ''
                }
            </div>

            <div class="card-controls" id="ctrl-${p.id}">
                ${renderControlsHTML(p.id)}
            </div>
        </div>
    `).join('');
}

function renderControlsHTML(productId) {
    const qty     = state.cart[productId] || 0;
    const product = state.products.find(p => p.id == productId);

    if (product?.stock === 0) {
        return `<button class="btn-add" disabled style="opacity:0.4;cursor:not-allowed">หมดสต็อก</button>`;
    }
    if (qty === 0) {
        return `<button class="btn-add" onclick="addToCart(${productId})"><span>＋</span> ใส่ตะกร้า</button>`;
    }
    return `<div class="qty-stepper">
                <button class="qty-btn" onclick="removeFromCart(${productId})">−</button>
                <span class="qty-display">${qty}</span>
                <button class="qty-btn" onclick="addToCart(${productId})">＋</button>
            </div>`;
}

function renderCardControls(productId) {
    const ctrl = document.getElementById(`ctrl-${productId}`);
    if (ctrl) ctrl.innerHTML = renderControlsHTML(productId);
}

function showLoadingState() {
    document.getElementById('product-list').innerHTML =
        `<div class="loading-state"><div class="spinner"></div><p>กำลังโหลดสินค้า...</p></div>`;
    document.getElementById('emptyState').classList.add('hidden');
}

function updateCartUI() {
    const totalQty = getTotalQty();

    // Badge
    const badge = document.getElementById('cartBadge');
    if (totalQty > 0) {
        badge.textContent = totalQty > 99 ? '99+' : totalQty;
        badge.classList.remove('hidden');
    } else {
        badge.classList.add('hidden');
    }

    // Summary
    document.getElementById('totalQty').textContent   = `${totalQty} ชิ้น`;
    document.getElementById('totalPrice').textContent = `฿${formatPrice(getTotalPrice())}`;
    document.getElementById('checkoutBtn').disabled   = totalQty === 0;

    renderCartItems();
}

function renderCartItems() {
    const cartItems = document.getElementById('cartItems');
    const cartEmpty = document.getElementById('cartEmpty');
    const entries   = Object.entries(state.cart);

    if (entries.length === 0) {
        cartEmpty.classList.remove('hidden');
        cartItems.classList.add('hidden');
        cartItems.innerHTML = '';
        return;
    }
    cartEmpty.classList.add('hidden');
    cartItems.classList.remove('hidden');

    cartItems.innerHTML = entries.map(([id, qty]) => {
        const p = state.products.find(pr => pr.id == id);
        if (!p) return '';
        return `
            <div class="cart-item">
                <div class="cart-item-img-placeholder">◫</div>
                <div class="cart-item-info">
                    <div class="cart-item-name">${escHtml(p.name)}</div>
                    <div class="cart-item-price">฿${formatPrice(p.price)} / ชิ้น</div>
                    <div class="cart-item-subtotal">รวม ฿${formatPrice(p.price * qty)}</div>
                </div>
                <div class="cart-item-controls">
                    <div class="mini-stepper">
                        <button class="mini-btn" onclick="removeFromCart(${id})">−</button>
                        <span class="mini-qty">${qty}</span>
                        <button class="mini-btn" onclick="addToCart(${id})">＋</button>
                    </div>
                    <button class="btn-remove" onclick="setCartQty(${id}, 0)">✕ นำออก</button>
                </div>
            </div>`;
    }).join('');
}

function showSuccessModal(items, totalAmount) {
    const rows = items.map(({ productId, quantity }) => {
        const p = state.products.find(pr => pr.id == productId);
        return `<div class="order-row">
                    <span>${escHtml(p?.name ?? 'สินค้า')} × ${quantity}</span>
                    <span>฿${formatPrice((p?.price ?? 0) * quantity)}</span>
                </div>`;
    }).join('');

    document.getElementById('orderSummary').innerHTML = `
        ${rows}
        <div class="order-row total">
            <span>รวมทั้งหมด</span>
            <span>฿${formatPrice(totalAmount)}</span>
        </div>`;
    closeModal('confirmModal');
    openModal('successModal');
}

function applyFilterAndSort() {
    const query    = document.getElementById('searchInput')?.value.toLowerCase().trim() || '';
    const sortVal  = document.getElementById('sortSelect')?.value || '';
    const category = document.getElementById('categoryFilter')?.value || '';

    let result = [...state.products];

    if (query) {
        result = result.filter(p =>
            p.name?.toLowerCase().includes(query) ||
            p.category?.toLowerCase().includes(query)
        );
    }
    if (category) {
        result = result.filter(p => p.category === category);
    }

    if (sortVal === 'price_asc')  result.sort((a, b) => a.price - b.price);
    if (sortVal === 'price_desc') result.sort((a, b) => b.price - a.price);
    if (sortVal === 'name_asc')   result.sort((a, b) => a.name?.localeCompare(b.name, 'th'));

    // สินค้าหมดไปท้าย
    result.sort((a, b) => (a.stock === 0 ? 1 : 0) - (b.stock === 0 ? 1 : 0));

    state.filtered = result;
    renderProducts(result);

    const countEl = document.getElementById('productCount');
    if (countEl) countEl.textContent = `${result.length} รายการ`;
}

async function populateCategoryFilter() {
    try {
        const data       = await request('GET', API.ENDPOINTS.CATEGORIES);
        const categories = data.categories || data.data || [];
        const select     = document.getElementById('categoryFilter');
        if (!select) return;

        select.innerHTML = `<option value="">หมวดหมู่ทั้งหมด</option>` +
            categories.map(cat =>
                `<option value="${escHtml(cat.name)}">${escHtml(cat.name)}</option>`
            ).join('');
    } catch (err) {
        console.error('โหลดหมวดหมู่ไม่สำเร็จ:', err);
    }
}

function openCheckoutConfirm() {
    if (getTotalQty() === 0) return;

    const rows = Object.entries(state.cart).map(([id, qty]) => {
        const p = state.products.find(pr => pr.id == id);
        return `<div class="order-row">
                    <span>${escHtml(p?.name ?? 'สินค้า')} × ${qty}</span>
                    <span>฿${formatPrice((p?.price ?? 0) * qty)}</span>
                </div>`;
    }).join('');

    document.getElementById('confirmSummary').innerHTML = `
        ${rows}
        <div class="order-row total" style="border-top:1px solid var(--border);margin-top:8px;padding-top:8px;font-weight:700;color:var(--success);">
            <span>รวม</span>
            <span>฿${formatPrice(getTotalPrice())}</span>
        </div>`;
    openModal('confirmModal');
}

function setCheckoutLoading(loading) {
    const btn = document.getElementById('checkoutBtn');
    btn.querySelector('.checkout-text')?.classList.toggle('hidden', loading);
    btn.querySelector('.checkout-loader')?.classList.toggle('hidden', !loading);
    btn.disabled = loading;
}

function openModal(id)  { const m = document.getElementById(id); if (m) { m.classList.add('active');    document.body.style.overflow = 'hidden'; } }
function closeModal(id) { const m = document.getElementById(id); if (m) { m.classList.remove('active'); document.body.style.overflow = '';       } }

function toggleCart() {
    const sidebar = document.getElementById('cartSidebar');
    const layout  = document.querySelector('.layout');
    const isOpen  = sidebar.classList.contains('open');
    sidebar.classList.toggle('open', !isOpen);
    layout.classList.toggle('cart-open', !isOpen);
}

document.querySelectorAll('.modal-overlay').forEach(overlay => {
    overlay.addEventListener('click', e => { if (e.target === overlay) closeModal(overlay.id); });
});

document.getElementById('logoutBtn')?.addEventListener('click', () => {

    localStorage.removeItem('sm_user');
    
    if (typeof showToast === 'function') {
        showToast('ออกจากระบบแล้ว', 'info');
    }
    setTimeout(() => {
        window.location.href = '../index.html'; 
    }, 500);
});

let toastTimer;
function showToast(message, type = 'info') {
    const toast = document.getElementById('toast');
    if (!toast) return;
    toast.textContent = message;
    toast.className = `toast ${type}`;
    clearTimeout(toastTimer);
    toastTimer = setTimeout(() => toast.className = 'toast hidden', 3500);
}
function escHtml(str)   { return String(str ?? '').replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;').replace(/"/g,'&quot;'); }
function formatPrice(p) { return Number(p).toLocaleString('th-TH', { minimumFractionDigits: 0, maximumFractionDigits: 2 }); }

function init() {
    // Cart
    document.getElementById('cartToggle')?.addEventListener('click', toggleCart);
    document.getElementById('cartClose')?.addEventListener('click', toggleCart);

    // Checkout
    document.getElementById('checkoutBtn')?.addEventListener('click', openCheckoutConfirm);
    document.getElementById('confirmPayBtn')?.addEventListener('click', checkout);
    document.getElementById('clearCartBtn')?.addEventListener('click', () => {
        if (getTotalQty() === 0) return;
        clearCart();
        showToast('ล้างตะกร้าแล้ว', 'info');
    });

    // Search / Filter / Sort — debounce เดิมจาก admin
    let searchTimer;
    document.getElementById('searchInput')?.addEventListener('input', () => {
        clearTimeout(searchTimer);
        searchTimer = setTimeout(applyFilterAndSort, 300);
    });
    document.getElementById('sortSelect')?.addEventListener('change', applyFilterAndSort);
    document.getElementById('categoryFilter')?.addEventListener('change', applyFilterAndSort);

    fetchProducts();
    updateCartUI();
}

document.addEventListener('DOMContentLoaded', init);