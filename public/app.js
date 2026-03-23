const API = {
    BASE_URL: 'http://localhost:5080',
    ENDPOINTS: {
        LOGIN:    '/login', 
        REGISTER: '/register',
        PRODUCTS: '/products',
        CATEGORIES: '/categories', 
        PRODUCT:  (id) => `/product/${id}` 
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
        throw new Error(data.message || data.error || `HTTP ${res.status}`);
    }
    return data;
}

async function fetchProducts() {
    showLoadingState();
    try {
        const data = await request('GET', API.ENDPOINTS.PRODUCTS);
        // 📌 Go ส่งมาในรูปแบบ {"products": [...]}
        state.products = data.products || [];
        applyFilterAndSort();
        updateProductCount();
    } catch (err) {
        console.error('fetchProducts error:', err);
        showToast('ไม่สามารถโหลดสินค้าได้: ' + err.message, 'error');
        renderProducts([]);
    }
}

async function createProduct(formData) {
    const data = await request('POST', API.ENDPOINTS.PRODUCTS, formData);
    showToast(data.message || 'เพิ่มสินค้าเรียบร้อย', 'success');
    await fetchProducts(); // โหลดข้อมูลใหม่หลังจากเพิ่มเสร็จ
}

async function updateProduct(id, formData) {
    const data = await request('PUT', API.ENDPOINTS.PRODUCT(id), formData);
    showToast(data.message || 'แก้ไขสินค้าเรียบร้อย', 'success');
    await fetchProducts(); // โหลดข้อมูลใหม่หลังจากแก้เสร็จ
}

async function deleteProduct(id) {
    const data = await request('DELETE', API.ENDPOINTS.PRODUCT(id));
    showToast(data.message || 'ลบสินค้าเรียบร้อย', 'success');
    await fetchProducts(); // โหลดข้อมูลใหม่หลังจากลบเสร็จ
}

// ... ฟังก์ชัน Login / Register เหมือนเดิม ...
async function loginUser(username, password) {
    const data = await request('POST', API.ENDPOINTS.LOGIN, { username, password });
    state.user = {
        id:       data.user?.id || data.id,
        username: data.user?.username || data.username || username,
        token:    data.token || data.accessToken
    };
    localStorage.setItem('sm_user', JSON.stringify(state.user));
    return state.user;
}

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

async function logoutUser() {
    try {
        if (state.user?.token) await request('POST', API.ENDPOINTS.LOGOUT);
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
            <div class="card-image-placeholder">◫</div>
            <div class="card-body">
                ${p.category ? `<span class="card-category">${escHtml(p.category)}</span>` : ''}
                <div class="card-name">${escHtml(p.name)}</div>
            </div>
            <div class="card-footer">
                <span class="card-price">฿${formatPrice(p.price)}</span>
                ${p.stock !== undefined
                    ? `<span class="card-stock ${p.stock === 0 ? 'out' : p.stock < 5 ? 'low' : ''}">
                        ${p.stock === 0 ? 'หมด' : `สต็อก: ${p.stock}`}
                       </span>`
                    : ''
                }
            </div>
            ${state.user ? `
            <div class="card-actions">
                <button class="card-btn card-btn-edit" onclick="openEditProduct(${p.id})">✎ แก้ไข</button>
                <button class="card-btn card-btn-delete" onclick="openDeleteConfirm(${p.id})">✕ ลบ</button>
            </div>` : ''}
        </div>
    `).join('');
}

function showLoadingState() {
    const grid = document.getElementById('product-list');
    grid.innerHTML = `<div class="loading-state"><div class="spinner"></div><p>กำลังโหลดสินค้า...</p></div>`;
    document.getElementById('emptyState').classList.add('hidden');
}

function updateProductCount() {
    const count = document.getElementById('productCount');
    if(count) count.textContent = `${state.filtered.length} รายการ`;
}

function updateNavUI() {
    const loginBtn = document.getElementById('loginBtn');
    const registerBtn = document.getElementById('registerBtn');
    const logoutBtn = document.getElementById('logoutBtn');
    const greeting = document.getElementById('userGreeting');
    const addBtn = document.getElementById('addProductBtn');

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
    
    const query = searchInput ? searchInput.value.toLowerCase().trim() : '';
    const sortVal = sortSelect ? sortSelect.value : '';

    let result = [...state.products];

    if (query) {
        result = result.filter(p => p.name?.toLowerCase().includes(query) || p.category?.toLowerCase().includes(query));
    }

    if (sortVal === 'price_asc')  result.sort((a, b) => a.price - b.price);
    if (sortVal === 'price_desc') result.sort((a, b) => b.price - a.price);
    if (sortVal === 'name_asc')   result.sort((a, b) => a.name?.localeCompare(b.name, 'th'));

    state.filtered = result;
    renderProducts(result);
    updateProductCount();
}

async function fetchCategories() {
    try {
        const data = await request('GET', API.ENDPOINTS.CATEGORIES);
        const categories = data.categories || data.data || [];

        const select = document.getElementById('productCategory');
        
        select.innerHTML = '<option value="" disabled selected>-- กรุณาเลือกหมวดหมู่ --</option>';
        categories.forEach(cat => {
            select.innerHTML += `<option value="${cat.id}">${cat.name}</option>`;
        });
    } catch (err) {
        console.error('โหลดหมวดหมู่ไม่สำเร็จ:', err);
        document.getElementById('productCategory').innerHTML = '<option value="" disabled>โหลดหมวดหมู่ไม่สำเร็จ</option>';
    }
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
    document.getElementById('productId').value        = product.id;
    document.getElementById('productName').value      = product.name || '';
    document.getElementById('productPrice').value     = product.price || '';
    document.getElementById('productCategory').value  = product.category_id || ''; 
    document.getElementById('productStock').value     = product.stock ?? '';
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
        category_id: parseInt(document.getElementById('productCategory').value) || 0,
        stock:       parseInt(document.getElementById('productStock').value) || 0,
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
        } else {
            await createProduct(formData);
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
        closeModal('deleteModal');
    } catch (err) {
        showToast('ลบไม่สำเร็จ: ' + err.message, 'error');
    } finally {
        btn.disabled = false;
        btn.textContent = 'ลบสินค้า';
    }
});

// ... Event Listeners ของ Login/Register คงเดิมเหมือนโค้ดเก่าของคุณอ้นเลยครับ ...
document.getElementById('loginForm')?.addEventListener('submit', async (e) => {
    e.preventDefault();
    const username = document.getElementById('loginUsername').value.trim();
    const password = document.getElementById('loginPassword').value;
    const btn = document.getElementById('loginSubmitBtn');

    if (!username || !password) return showError('loginError', 'กรุณากรอกข้อมูลให้ครบ');
    hideError('loginError'); setLoading(btn, true);

    try {
        await loginUser(username, password);
        closeModal('loginModal');
        updateNavUI();
        renderProducts(state.filtered);
        showToast(`ยินดีต้อนรับ, ${state.user.username}!`, 'success');
    } catch (err) {
        showError('loginError', err.message || 'ชื่อผู้ใช้หรือรหัสผ่านไม่ถูกต้อง');
    } finally { setLoading(btn, false); }
});

document.getElementById('logoutBtn')?.addEventListener('click', async () => {
    await logoutUser();
    updateNavUI();
    renderProducts(state.filtered);
    showToast('ออกจากระบบแล้ว', 'info');
});

// ฟังก์ชันเปิด/ปิด Modal และ Utils
function openModal(id) { const m = document.getElementById(id); if(m){ m.classList.add('active'); document.body.style.overflow = 'hidden';} }
function closeModal(id) { const m = document.getElementById(id); if(m){ m.classList.remove('active'); document.body.style.overflow = '';} }
function switchModal(from, to) { closeModal(from); setTimeout(() => openModal(to), 150); }
document.querySelectorAll('.modal-overlay').forEach(overlay => { overlay.addEventListener('click', (e) => { if (e.target === overlay) closeModal(overlay.id); }); });
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
function setLoading(btn, loading) { if(!btn) return; btn.querySelector('.btn-text')?.classList.toggle('hidden', loading); btn.querySelector('.btn-loader')?.classList.toggle('hidden', !loading); btn.disabled = loading; }
function showError(id, msg) { const el = document.getElementById(id); if(el) { el.textContent = msg; el.classList.remove('hidden'); } }
function hideError(id) { document.getElementById(id)?.classList.add('hidden'); }
function escHtml(str) { return String(str ?? '').replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;'); }
function formatPrice(price) { return Number(price).toLocaleString('th-TH', { minimumFractionDigits: 0, maximumFractionDigits: 2 }); }

function init() {
    updateNavUI();
    fetchProducts(); 
    fetchCategories();

    let searchTimer;
    const searchInput = document.getElementById('searchInput');
    if(searchInput) searchInput.addEventListener('input', () => { clearTimeout(searchTimer); searchTimer = setTimeout(applyFilterAndSort, 300); });
    const sortSelect = document.getElementById('sortSelect');
    if(sortSelect) sortSelect.addEventListener('change', applyFilterAndSort);
}

document.addEventListener('DOMContentLoaded', init);