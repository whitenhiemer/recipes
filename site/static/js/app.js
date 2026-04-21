// Meal Plan - localStorage backed
// Structure: { breakfast: [...], lunch: [...], dinner: [...], treats: [...] }
// Main meals have 4 slots, treats has 2
const MEALPLAN_KEY = 'recipe-mealplan';
const SLOTS_PER_MEAL = 4;
const SLOTS_PER_TREATS = 2;

function getMealPlan() {
    const data = localStorage.getItem(MEALPLAN_KEY);
    if (!data) return { breakfast: [], lunch: [], dinner: [], treats: [] };
    const plan = JSON.parse(data);
    if (!plan.breakfast) plan.breakfast = [];
    if (!plan.lunch) plan.lunch = [];
    if (!plan.dinner) plan.dinner = [];
    if (!plan.treats) plan.treats = [];
    return plan;
}

function saveMealPlan(plan) {
    localStorage.setItem(MEALPLAN_KEY, JSON.stringify(plan));
}

function addToMealPlan(slug, title) {
    const meal = document.getElementById('mealplan-meal').value;
    const plan = getMealPlan();
    const maxSlots = meal === 'treats' ? SLOTS_PER_TREATS : SLOTS_PER_MEAL;

    if (plan[meal].length >= maxSlots) {
        alert('All ' + maxSlots + ' ' + meal + ' slots are full. Clear one first.');
        return;
    }

    plan[meal].push({ slug: slug, title: title });
    saveMealPlan(plan);
    renderMealPlan();
}

function selectSlot(meal, slotIndex, selectEl) {
    const plan = getMealPlan();
    const slug = selectEl.value;

    if (!slug) {
        if (plan[meal][slotIndex]) {
            plan[meal].splice(slotIndex, 1);
        }
    } else {
        const data = getRecipeData();
        const all = data ? data.all : [];
        const recipe = all.find(function(r) { return r.slug === slug; });
        const title = recipe ? recipe.title : slug;

        while (plan[meal].length < slotIndex) {
            plan[meal].push(null);
        }
        plan[meal][slotIndex] = { slug: slug, title: title };
    }

    plan[meal] = plan[meal].filter(function(e) { return e !== null; });
    saveMealPlan(plan);
    renderMealPlan();
}

function clearSlot(meal, index) {
    const plan = getMealPlan();
    plan[meal].splice(index, 1);
    saveMealPlan(plan);
    renderMealPlan();
}

function clearMealPlan() {
    localStorage.removeItem(MEALPLAN_KEY);
    renderMealPlan();
}

function buildSlotOptions(select, meal, selectedSlug) {
    while (select.firstChild) select.removeChild(select.firstChild);

    var placeholder = document.createElement('option');
    placeholder.value = '';
    placeholder.textContent = 'Choose a recipe...';
    select.appendChild(placeholder);

    var data = getRecipeData();
    if (!data) return;

    var catRecipes = data[meal] || [];
    var allRecipes = data.all || [];
    var catSlugs = new Set(catRecipes.map(function(r) { return r.slug; }));
    var otherRecipes = allRecipes.filter(function(r) { return !catSlugs.has(r.slug); });

    if (catRecipes.length) {
        var catGroup = document.createElement('optgroup');
        catGroup.label = meal.charAt(0).toUpperCase() + meal.slice(1);
        catRecipes.forEach(function(r) {
            var opt = document.createElement('option');
            opt.value = r.slug;
            opt.textContent = r.title;
            if (r.slug === selectedSlug) opt.selected = true;
            catGroup.appendChild(opt);
        });
        select.appendChild(catGroup);
    }

    if (otherRecipes.length) {
        var otherGroup = document.createElement('optgroup');
        otherGroup.label = 'Other Recipes';
        otherRecipes.forEach(function(r) {
            var opt = document.createElement('option');
            opt.value = r.slug;
            opt.textContent = r.title + ' (' + r.category + ')';
            if (r.slug === selectedSlug) opt.selected = true;
            otherGroup.appendChild(opt);
        });
        select.appendChild(otherGroup);
    }
}

function renderMealPlan() {
    var plan = getMealPlan();

    ['breakfast', 'lunch', 'dinner', 'treats'].forEach(function(meal) {
        var slots = document.querySelectorAll('.meal-slot[data-meal="' + meal + '"]');
        slots.forEach(function(slot, i) {
            var select = slot.querySelector('.slot-select');
            if (!select) return;
            var entry = plan[meal][i];
            var selectedSlug = entry ? entry.slug : '';

            buildSlotOptions(select, meal, selectedSlug);
            slot.classList.toggle('filled', !!entry);
        });
    });
}

function generateShoppingList() {
    const plan = getMealPlan();
    const trip1 = new Set();
    const trip2 = new Set();
    const trip1Slugs = [];
    const trip2Slugs = [];

    ['breakfast', 'lunch', 'dinner'].forEach(meal => {
        plan[meal].forEach((entry, i) => {
            if (i < 2) {
                if (!trip1.has(entry.slug)) {
                    trip1Slugs.push(entry.slug);
                    trip1.add(entry.slug);
                }
            } else {
                if (!trip2.has(entry.slug)) {
                    trip2Slugs.push(entry.slug);
                    trip2.add(entry.slug);
                }
            }
        });
    });

    // Treats go to trip 1 (shelf-stable)
    (plan.treats || []).forEach(entry => {
        if (!trip1.has(entry.slug)) {
            trip1Slugs.push(entry.slug);
            trip1.add(entry.slug);
        }
    });

    if (trip1Slugs.length === 0 && trip2Slugs.length === 0) {
        alert('Add recipes to your meal plan first.');
        return;
    }

    if (trip2Slugs.length === 0) {
        window.location.href = '/shopping-list?slugs=' + trip1Slugs.join(',');
        return;
    }

    const params = [];
    if (trip1Slugs.length) params.push('trip1=' + trip1Slugs.join(','));
    if (trip2Slugs.length) params.push('trip2=' + trip2Slugs.join(','));
    window.location.href = '/shopping-list?' + params.join('&');
}

function getRecipeData() {
    const el = document.getElementById('recipe-data');
    if (!el) return null;
    const data = JSON.parse(el.textContent);
    // Remove trailing nulls from template comma hack
    ['breakfast', 'lunch', 'dinner', 'treats', 'all'].forEach(key => {
        if (data[key]) data[key] = data[key].filter(r => r !== null);
    });
    return data;
}

function shuffle(arr) {
    const a = [...arr];
    for (let i = a.length - 1; i > 0; i--) {
        const j = Math.floor(Math.random() * (i + 1));
        [a[i], a[j]] = [a[j], a[i]];
    }
    return a;
}

function randomizeMealPlan() {
    const data = getRecipeData();
    if (!data) return;

    const plan = { breakfast: [], lunch: [], dinner: [], treats: [] };
    const allRecipes = data.all || [];

    ['breakfast', 'lunch', 'dinner'].forEach(meal => {
        let pool = data[meal] || [];

        if (pool.length < SLOTS_PER_MEAL) {
            pool = allRecipes.map(r => ({ slug: r.slug, title: r.title }));
        }

        const picked = shuffle(pool).slice(0, SLOTS_PER_MEAL);
        plan[meal] = picked.map(r => ({ slug: r.slug, title: r.title }));
    });

    let treatsPool = data.treats || [];
    if (treatsPool.length < SLOTS_PER_TREATS) {
        treatsPool = allRecipes.map(r => ({ slug: r.slug, title: r.title }));
    }
    const treatsP = shuffle(treatsPool).slice(0, SLOTS_PER_TREATS);
    plan.treats = treatsP.map(r => ({ slug: r.slug, title: r.title }));

    saveMealPlan(plan);
    renderMealPlan();
}

function filterMealPlanRecipes(query) {
    const items = document.querySelectorAll('.mealplan-recipe-item');
    const q = query.toLowerCase();
    items.forEach(item => {
        const title = item.dataset.title.toLowerCase();
        const category = item.dataset.category.toLowerCase();
        item.style.display = (title.includes(q) || category.includes(q)) ? '' : 'none';
    });
}

// Shopping List
function toggleShoppingItem(checkbox) {
    // Visual only, state managed by the checkbox itself
}

function copyShoppingList() {
    const items = document.querySelectorAll('.shopping-item');
    let text = 'Shopping List\n============\n\n';

    items.forEach(item => {
        const name = item.querySelector('strong').textContent;
        const amounts = item.querySelectorAll('li');
        text += '- ' + name + '\n';
        amounts.forEach(a => {
            text += '  ' + a.textContent + '\n';
        });
    });

    navigator.clipboard.writeText(text).then(() => {
        alert('Copied to clipboard!');
    });
}

// Wake Lock - keeps screen on while cooking
let wakeLock = null;

async function toggleWakeLock(enabled) {
    if (!('wakeLock' in navigator)) return;
    if (enabled) {
        try {
            wakeLock = await navigator.wakeLock.request('screen');
            wakeLock.addEventListener('release', () => {
                const toggle = document.getElementById('wake-lock-toggle');
                if (toggle) toggle.checked = false;
            });
        } catch (e) {}
    } else if (wakeLock) {
        wakeLock.release();
        wakeLock = null;
    }
}

// Re-acquire wake lock when returning to the tab
document.addEventListener('visibilitychange', async () => {
    const toggle = document.getElementById('wake-lock-toggle');
    if (toggle && toggle.checked && document.visibilityState === 'visible') {
        toggleWakeLock(true);
    }
});

// Initialize on page load
document.addEventListener('DOMContentLoaded', () => {
    if (document.getElementById('breakfast-slots')) {
        renderMealPlan();
    }
    const wakeToggle = document.getElementById('wake-lock-toggle');
    if (wakeToggle && !('wakeLock' in navigator)) {
        wakeToggle.parentElement.style.display = 'none';
    }
});
