const API = {
    BASE_URL: 'http://localhost:5000',  // TODO: เปลี่ยนเป็น URL backend ของคุณ
    ENDPOINTS: {
        LOGIN:    '/login', 
        REGISTER: '/register',
        PRODUCTS: '/products',              // GET list, POST create
        PRODUCT:  (id) => `/product/${id}`  // GET one, PUT update, DELETE
    }
};

let state = {
    user: null,
    products: [],
    filtered: [],
    isLoading: false
};

const savedUser = localStorage.getItem('sm_user');
if (savedUser) {
    try { state.user = JSON.parse(savedUser); } catch(e) {}
}

async function request(method, endpoint, body = null) {
    const url = API.BASE_URL + endpoint;
    const headers = { 'Content-Type': 'application/json' };

    if (state.user?.token) {
        headers['Authorization'] = `Bearer ${state.user.token}`;
    }

    const opts = { method, headers };
    if (body) opts.body = JSON.stringify(body);

    const res = await fetch(url, opts);
    const data = await res.json().catch(() => ({}));

    if (!res.ok) {
        throw new Error(data.message || `HTTP ${res.status}`);
    }
    return data;
}

/**
 * TODO: BACKEND — GET /products
 * คาดหวัง response: { products: [...], total: 0 } หรือ array ตรงๆ: [...]
 */
async function fetchProducts() {
    showLoadingState();
    try {
        const data = await request('GET', API.ENDPOINTS.PRODUCTS);
        state.products = Array.isArray(data) ? data : (data.products || []);
        applyFilterAndSort();
        updateProductCount();
    } catch (err) {
        console.error('fetchProducts error:', err);
        showToast('ไม่สามารถโหลดสินค้าได้: ' + err.message, 'error');
        renderProducts([]);
    }
}

/**
 * TODO: BACKEND — POST /products
 * body: { name, price, category, stock, description, image }
 * คาดหวัง response: { product: {...} } หรือ object ตรงๆ
 */
async function createProduct(formData) {
    const data = await request('POST', API.ENDPOINTS.PRODUCTS, formData);
    const newProduct = data.product || data;
    state.products.unshift(newProduct);
    applyFilterAndSort();
    updateProductCount();
    return newProduct;
}

/**
 * TODO: BACKEND — PUT /product/:id
 * body: { name, price, category, stock, description, image }
 * คาดหวัง response: { product: {...} } หรือ object ตรงๆ
 */
async function updateProduct(id, formData) {
    const data = await request('PUT', API.ENDPOINTS.PRODUCT(id), formData);
    const updated = data.product || data;
    const idx = state.products.findIndex(p => p.id == id);
    if (idx !== -1) state.products[idx] = updated;
    applyFilterAndSort();
    return updated;
}

/**
 * TODO: BACKEND — DELETE /product/:id
 * คาดหวัง response: { message: "Deleted" } หรือ status 204
 */
async function deleteProduct(id) {
    await request('DELETE', API.ENDPOINTS.PRODUCT(id));
    state.products = state.products.filter(p => p.id != id);
    applyFilterAndSort();
    updateProductCount();
}

/**
 * TODO: BACKEND — POST /login
 * body: { username, password }
 * คาดหวัง response: { user: { id, username }, token: "jwt..." }
 */
async function loginUser(username, password) {
    const data = await request('POST', API.ENDPOINTS.LOGIN, { username, password });
    state.user = {
        id:       data.user?.id || data.id,
        username: data.user?.username || data.username,
        token:    data.token || data.accessToken
    };
    localStorage.setItem('sm_user', JSON.stringify(state.user));
    return state.user;
}

/**
 * TODO: BACKEND — POST /register
 * body: { username, email, password }
 * คาดหวัง response: { user: {...}, token: "..." } หรือแค่ { message: "success" }
 */
async function registerUser(username, email, password) {
    const data = await request('POST', API.ENDPOINTS.REGISTER, { username, email, password });
    if (data.token || data.accessToken) {
        state.user = {
            id:       data.user?.id || data.id,
            username: data.user?.username || data.username || username,
            token:    data.token || data.accessToken
        };
        localStorage.setItem('sm_user', JSON.stringify(state.user));
    }
    return data;
}

/**
 * TODO: BACKEND — POST /logout (optional)
 */
async function logoutUser() {
    try {
        if (state.user?.token) {
            await request('POST', API.ENDPOINTS.LOGOUT);
        }
    } catch(e) {}
    state.user = null;
    localStorage.removeItem('sm_user');
}

function renderProducts(products) {
    const grid = document.getElementById('product-list');
    const empty = document.getElementById('emptyState');

    if (products.length === 0) {
        grid.innerHTML = '';
        empty.classList.remove('hidden');
        return;
    }

    empty.classList.add('hidden');
    grid.innerHTML = products.map((p, i) => `
        <div class="product-card" style="animation-delay: ${i * 0.04}s">
            ${p.Image
                ? `<div class="card-image"><img src="${escHtml(p.Image)}" alt="${escHtml(p.Name)}" loading="lazy"></div>`
                : `<div class="card-image-placeholder">◫</div>`
            }
            <div class="card-body">
                ${p.Category ? `<span class="card-category">${escHtml(p.Category)}</span>` : ''}
                <div class="card-name">${escHtml(p.Name)}</div>
                ${p.Description ? `<div class="card-desc">${escHtml(p.Description)}</div>` : ''}
            </div>
            <div class="card-footer">
                <span class="card-price">฿${formatPrice(p.Price)}</span>
                ${p.Stock !== undefined
                    ? `<span class="card-stock ${p.Stock === 0 ? 'out' : p.Stock < 5 ? 'low' : ''}">
                        ${p.Stock === 0 ? 'หมด' : `สต็อก: ${p.Stock}`}
                       </span>`
                    : ''
                }
            </div>
            ${state.user ? `
            <div class="card-actions">
                <button class="card-btn card-btn-edit" onclick="openEditProduct(${p.ID})">✎ แก้ไข</button>
                <button class="card-btn card-btn-delete" onclick="openDeleteConfirm(${p.ID})">✕ ลบ</button>
            </div>` : ''}
        </div>
    `).join('');
}

function showLoadingState() {
    const grid = document.getElementById('product-list');
    grid.innerHTML = `
        <div class="loading-state">
            <div class="spinner"></div>
            <p>กำลังโหลดสินค้า...</p>
        </div>`;
    document.getElementById('emptyState').classList.add('hidden');
}

function updateProductCount() {
    const count = document.getElementById('productCount');
    if(count) count.textContent = `${state.filtered.length} รายการ`;
}

function updateNavUI() {
    const loginBtn     = document.getElementById('loginBtn');
    const registerBtn  = document.getElementById('registerBtn');
    const logoutBtn    = document.getElementById('logoutBtn');
    const greeting     = document.getElementById('userGreeting');
    const addBtn       = document.getElementById('addProductBtn');

    if (state.user) {
        if(loginBtn) loginBtn.classList.add('hidden');
        if(registerBtn) registerBtn.classList.add('hidden');
        if(logoutBtn) logoutBtn.classList.remove('hidden');
        if(greeting) {
            greeting.classList.remove('hidden');
            greeting.textContent = `สวัสดี, ${state.user.username}`;
        }
        if(addBtn) addBtn.style.display = '';
    } else {
        if(loginBtn) loginBtn.classList.remove('hidden');
        if(registerBtn) registerBtn.classList.remove('hidden');
        if(logoutBtn) logoutBtn.classList.add('hidden');
        if(greeting) greeting.classList.add('hidden');
        if(addBtn) addBtn.style.display = 'none';
    }
}

function applyFilterAndSort() {
    const searchInput = document.getElementById('searchInput');
    const sortSelect = document.getElementById('sortSelect');
    
    const query   = searchInput ? searchInput.value.toLowerCase().trim() : '';
    const sortVal = sortSelect ? sortSelect.value : '';

    let result = [...state.products];

    if (query) {
        result = result.filter(p =>
            p.name?.toLowerCase().includes(query) ||
            p.description?.toLowerCase().includes(query) ||
            p.category?.toLowerCase().includes(query)
        );
    }

    if (sortVal === 'price_asc')  result.sort((a, b) => a.price - b.price);
    if (sortVal === 'price_desc') result.sort((a, b) => b.price - a.price);
    if (sortVal === 'name_asc')   result.sort((a, b) => a.name?.localeCompare(b.name, 'th'));

    state.filtered = result;
    renderProducts(result);
    updateProductCount();
}

function openAddProduct() {
    if (!state.user) { showToast('กรุณาเข้าสู่ระบบก่อน', 'info'); return; }
    document.getElementById('productModalTitle').textContent = 'เพิ่มสินค้าใหม่';
    document.getElementById('productForm').reset();
    document.getElementById('productId').value = '';
    hideError('productError');
    openModal('addProductModal');
}

function openEditProduct(id) {
    const product = state.products.find(p => p.id == id);
    if (!product) return;

    document.getElementById('productModalTitle').textContent = 'แก้ไขสินค้า';
    document.getElementById('productId').value        = product.ID;
    document.getElementById('productName').value      = product.Name || '';
    document.getElementById('productPrice').value     = product.Price || '';
    document.getElementById('productCategory').value  = product.Category || '';
    document.getElementById('productStock').value     = product.Stock ?? '';
    document.getElementById('productDescription').value = product.Description || '';
    document.getElementById('productImage').value     = product.image || '';
    hideError('productError');
    openModal('addProductModal');
}

function openDeleteConfirm(id) {
    document.getElementById('deleteProductId').value = id;
    openModal('deleteModal');
}

document.getElementById('productForm')?.addEventListener('submit', async (e) => {
    e.preventDefault();
    const id = document.getElementById('productId').value;
    const formData = {
        name:        document.getElementById('productName').value.trim(),
        price:       parseFloat(document.getElementById('productPrice').value),
        category:    document.getElementById('productCategory').value.trim(),
        stock:       parseInt(document.getElementById('productStock').value) || 0,
        description: document.getElementById('productDescription').value.trim(),
        image:       document.getElementById('productImage').value.trim()
    };

    if (!formData.name || isNaN(formData.price)) {
        showError('productError', 'กรุณากรอกชื่อสินค้าและราคา');
        return;
    }

    const btn = document.getElementById('productSubmitBtn');
    setLoading(btn, true);
    hideError('productError');

    try {
        if (id) {
            await updateProduct(id, formData);
            showToast('แก้ไขสินค้าเรียบร้อย', 'success');
        } else {
            await createProduct(formData);
            showToast('เพิ่มสินค้าเรียบร้อย', 'success');
        }
        closeModal('addProductModal');
    } catch (err) {
        showError('productError', err.message);
    } finally {
        setLoading(btn, false);
    }
});

document.getElementById('confirmDeleteBtn')?.addEventListener('click', async () => {
    const id = document.getElementById('deleteProductId').value;
    const btn = document.getElementById('confirmDeleteBtn');
    btn.disabled = true;
    btn.textContent = '⟳ กำลังลบ...';
    try {
        await deleteProduct(id);
        showToast('ลบสินค้าเรียบร้อย', 'success');
        closeModal('deleteModal');
    } catch (err) {
        showToast('ลบไม่สำเร็จ: ' + err.message, 'error');
    } finally {
        btn.disabled = false;
        btn.textContent = 'ลบสินค้า';
    }
});

document.getElementById('loginForm')?.addEventListener('submit', async (e) => {
    e.preventDefault();
    const username = document.getElementById('loginUsername').value.trim();
    const password = document.getElementById('loginPassword').value;
    const btn = document.getElementById('loginSubmitBtn');

    if (!username || !password) {
        showError('loginError', 'กรุณากรอกข้อมูลให้ครบ');
        return;
    }

    hideError('loginError');
    setLoading(btn, true);

    try {
        await loginUser(username, password);
        closeModal('loginModal');
        updateNavUI();
        renderProducts(state.filtered);
        showToast(`ยินดีต้อนรับ, ${state.user.username}!`, 'success');
    } catch (err) {
        showError('loginError', err.message || 'ชื่อผู้ใช้หรือรหัสผ่านไม่ถูกต้อง');
    } finally {
        setLoading(btn, false);
    }
});

document.getElementById('registerForm')?.addEventListener('submit', async (e) => {
    e.preventDefault();
    const username  = document.getElementById('regUsername').value.trim();
    const email     = document.getElementById('regEmail') ? document.getElementById('regEmail').value.trim() : '';
    const password  = document.getElementById('regPassword').value;
    const confirm   = document.getElementById('regConfirmPassword') ? document.getElementById('regConfirmPassword').value : password;
    const btn       = document.getElementById('registerSubmitBtn');

    if (!username || !password) {
        showError('registerError', 'กรุณากรอกชื่อผู้ใช้และรหัสผ่าน');
        return;
    }
    if (password !== confirm) {
        showError('registerError', 'รหัสผ่านไม่ตรงกัน');
        return;
    }

    hideError('registerError');
    setLoading(btn, true);

    try {
        await registerUser(username, email, password);
        closeModal('registerModal');
        updateNavUI();
        showToast('สมัครสมาชิกเรียบร้อย!', 'success');
        if (!state.user) openModal('loginModal');
    } catch (err) {
        showError('registerError', err.message || 'ไม่สามารถสมัครสมาชิกได้');
    } finally {
        setLoading(btn, false);
    }
});

document.getElementById('logoutBtn')?.addEventListener('click', async () => {
    await logoutUser();
    updateNavUI();
    renderProducts(state.filtered);
    showToast('ออกจากระบบแล้ว', 'info');
});

function openModal(id) {
    const modal = document.getElementById(id);
    if(modal) {
        modal.classList.add('active');
        document.body.style.overflow = 'hidden';
    }
}

function closeModal(id) {
    const modal = document.getElementById(id);
    if(modal) {
        modal.classList.remove('active');
        document.body.style.overflow = '';
    }
}

function switchModal(from, to) {
    closeModal(from);
    setTimeout(() => openModal(to), 150);
}

document.querySelectorAll('.modal-overlay').forEach(overlay => {
    overlay.addEventListener('click', (e) => {
        if (e.target === overlay) closeModal(overlay.id);
    });
});

document.getElementById('loginBtn')?.addEventListener('click', () => openModal('loginModal'));
document.getElementById('registerBtn')?.addEventListener('click', () => openModal('registerModal'));
document.getElementById('addProductBtn')?.addEventListener('click', openAddProduct);

let toastTimer;
function showToast(message, type = 'info') {
    const toast = document.getElementById('toast');
    if(!toast) return;
    toast.textContent = message;
    toast.className = `toast ${type}`;
    clearTimeout(toastTimer);
    toastTimer = setTimeout(() => toast.className = 'toast hidden', 3500);
}

function setLoading(btn, loading) {
    if(!btn) return;
    const text = btn.querySelector('.btn-text');
    const loader = btn.querySelector('.btn-loader');
    if(text) text.classList.toggle('hidden', loading);
    if(loader) loader.classList.toggle('hidden', !loading);
    btn.disabled = loading;
}

function showError(id, msg) {
    const el = document.getElementById(id);
    if(el) {
        el.textContent = msg;
        el.classList.remove('hidden');
    }
}

function hideError(id) {
    const el = document.getElementById(id);
    if(el) el.classList.add('hidden');
}

function escHtml(str) {
    return String(str ?? '')
        .replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;')
        .replace(/"/g, '&quot;');
}

function formatPrice(price) {
    return Number(price).toLocaleString('th-TH', {
        minimumFractionDigits: 0,
        maximumFractionDigits: 2
    });
}

function init() {
    updateNavUI();
    fetchProducts(); 

    let searchTimer;
    const searchInput = document.getElementById('searchInput');
    if(searchInput) {
        searchInput.addEventListener('input', () => {
            clearTimeout(searchTimer);
            searchTimer = setTimeout(applyFilterAndSort, 300);
        });
    }
    
    const sortSelect = document.getElementById('sortSelect');
    if(sortSelect) {
        sortSelect.addEventListener('change', applyFilterAndSort);
    }
}

document.addEventListener('DOMContentLoaded', init);