// Meal Plan - localStorage backed
const MEALPLAN_KEY = 'recipe-mealplan';

function getMealPlan() {
    const data = localStorage.getItem(MEALPLAN_KEY);
    return data ? JSON.parse(data) : {};
}

function saveMealPlan(plan) {
    localStorage.setItem(MEALPLAN_KEY, JSON.stringify(plan));
}

function addToMealPlan(slug, title) {
    const day = document.getElementById('mealplan-day').value;
    const meal = document.getElementById('mealplan-meal').value;
    const plan = getMealPlan();

    if (!plan[day]) plan[day] = {};
    plan[day][meal] = { slug: slug, title: title };
    saveMealPlan(plan);
    renderMealPlan();
}

function clearSlot(day, meal) {
    const plan = getMealPlan();
    if (plan[day]) {
        delete plan[day][meal];
        if (Object.keys(plan[day]).length === 0) delete plan[day];
    }
    saveMealPlan(plan);
    renderMealPlan();
}

function clearMealPlan() {
    localStorage.removeItem(MEALPLAN_KEY);
    renderMealPlan();
}

function renderMealPlan() {
    const plan = getMealPlan();
    document.querySelectorAll('.meal-slot').forEach(slot => {
        const day = slot.dataset.day;
        const meal = slot.dataset.meal;
        const recipeSpan = slot.querySelector('.slot-recipe');
        const clearBtn = slot.querySelector('.slot-clear');

        if (plan[day] && plan[day][meal]) {
            recipeSpan.textContent = plan[day][meal].title;
            clearBtn.style.display = 'inline';
        } else {
            recipeSpan.textContent = '';
            clearBtn.style.display = 'none';
        }
    });
}

function generateShoppingList() {
    const plan = getMealPlan();
    const slugs = [];
    const seen = new Set();

    Object.values(plan).forEach(meals => {
        Object.values(meals).forEach(entry => {
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

function filterMealPlanRecipes(query) {
    const items = document.querySelectorAll('.mealplan-recipe-item');
    const q = query.toLowerCase();
    items.forEach(item => {
        const title = item.dataset.title.toLowerCase();
        item.style.display = title.includes(q) ? '' : 'none';
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

// Hide wake lock toggle if browser doesn't support it
document.addEventListener('DOMContentLoaded', () => {
    if (document.getElementById('mealplan-body')) {
        renderMealPlan();
    }
    const wakeToggle = document.getElementById('wake-lock-toggle');
    if (wakeToggle && !('wakeLock' in navigator)) {
        wakeToggle.parentElement.style.display = 'none';
    }
});
