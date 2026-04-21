// Meal Plan - localStorage backed
// Structure: { breakfast: [{slug, title}, ...], lunch: [...], dinner: [...] }
// Each meal type has up to 4 slots
const MEALPLAN_KEY = 'recipe-mealplan';
const SLOTS_PER_MEAL = 4;

function getMealPlan() {
    const data = localStorage.getItem(MEALPLAN_KEY);
    if (!data) return { breakfast: [], lunch: [], dinner: [] };
    const plan = JSON.parse(data);
    if (!plan.breakfast) plan.breakfast = [];
    if (!plan.lunch) plan.lunch = [];
    if (!plan.dinner) plan.dinner = [];
    return plan;
}

function saveMealPlan(plan) {
    localStorage.setItem(MEALPLAN_KEY, JSON.stringify(plan));
}

function addToMealPlan(slug, title) {
    const meal = document.getElementById('mealplan-meal').value;
    const plan = getMealPlan();

    if (plan[meal].length >= SLOTS_PER_MEAL) {
        alert('All 4 ' + meal + ' slots are full. Clear one first.');
        return;
    }

    plan[meal].push({ slug: slug, title: title });
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

function renderMealPlan() {
    const plan = getMealPlan();

    ['breakfast', 'lunch', 'dinner'].forEach(meal => {
        const slots = document.querySelectorAll('.meal-slot[data-meal="' + meal + '"]');
        slots.forEach((slot, i) => {
            const recipeSpan = slot.querySelector('.slot-recipe');
            const clearBtn = slot.querySelector('.slot-clear');
            const entry = plan[meal][i];

            if (entry) {
                recipeSpan.textContent = entry.title;
                recipeSpan.closest('.meal-slot').classList.add('filled');
                clearBtn.style.display = 'inline';
            } else {
                recipeSpan.textContent = '';
                recipeSpan.closest('.meal-slot').classList.remove('filled');
                clearBtn.style.display = 'none';
            }
        });
    });

    // Update slot count display
    ['breakfast', 'lunch', 'dinner'].forEach(meal => {
        const count = plan[meal].length;
        const counter = document.getElementById(meal + '-count');
        if (counter) counter.textContent = count + ' / ' + SLOTS_PER_MEAL;
    });
}

function generateShoppingList() {
    const plan = getMealPlan();
    const slugs = [];
    const seen = new Set();

    ['breakfast', 'lunch', 'dinner'].forEach(meal => {
        plan[meal].forEach(entry => {
            if (!seen.has(entry.slug)) {
                slugs.push(entry.slug);
                seen.add(entry.slug);
            }
        });
    });

    if (slugs.length === 0) {
        alert('Add recipes to your meal plan first.');
        return;
    }

    window.location.href = '/shopping-list?slugs=' + slugs.join(',');
}

function getRecipeData() {
    const el = document.getElementById('recipe-data');
    if (!el) return null;
    const data = JSON.parse(el.textContent);
    // Remove trailing nulls from template comma hack
    ['breakfast', 'lunch', 'dinner', 'all'].forEach(key => {
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

    const plan = { breakfast: [], lunch: [], dinner: [] };
    const allRecipes = data.all || [];

    ['breakfast', 'lunch', 'dinner'].forEach(meal => {
        let pool = data[meal] || [];

        // Fall back to all recipes if category pool is too small
        if (pool.length < SLOTS_PER_MEAL) {
            pool = allRecipes.map(r => ({ slug: r.slug, title: r.title }));
        }

        const picked = shuffle(pool).slice(0, SLOTS_PER_MEAL);
        plan[meal] = picked.map(r => ({ slug: r.slug, title: r.title }));
    });

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
