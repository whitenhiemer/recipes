package recipe

import "strings"

var priceTable = map[string]float64{
	// Proteins
	"ground beef":       5.99,
	"lean ground beef":  6.49,
	"chicken breast":    7.99,
	"chicken thighs":    5.49,
	"chicken thigh":     5.49,
	"pork loin":         12.99,
	"pork chops":        7.99,
	"pork chop":         7.99,
	"bacon":             6.99,
	"pancetta":          4.99,
	"breakfast sausage": 4.99,
	"sausage":           4.99,
	"tuna":              1.99,
	"eggs":              4.49,
	"egg":               4.49,
	"egg yolk":          4.49,

	// Dairy
	"butter":          4.99,
	"unsalted butter": 4.99,
	"cold butter":     4.99,
	"milk":            3.99,
	"whole milk":      3.99,
	"cream":           4.49,
	"sour cream":      2.49,
	"yogurt":          3.99,
	"cheddar cheese":  4.49,
	"cheddar":         4.49,
	"american cheese": 3.99,
	"swiss cheese":    4.49,
	"cheese":          4.49,
	"parmesan":        5.99,
	"parmesan cheese": 5.99,

	// Grains & Pasta
	"flour":            3.99,
	"all-purpose flour": 3.99,
	"bread flour":      4.49,
	"bread":            3.49,
	"elbow macaroni":   1.49,
	"macaroni":         1.49,
	"spaghetti":        1.49,
	"pasta":            1.49,
	"tortillas":        3.49,
	"flour tortillas":  3.49,
	"tortilla":         3.49,
	"taco shells":      2.49,
	"biscuits":         2.99,
	"rolled oats":      3.99,
	"oats":             3.99,
	"masa harina":      3.99,
	"rice":             2.99,

	// Canned goods
	"crushed tomatoes":  1.99,
	"diced tomatoes":    1.49,
	"tomato sauce":      1.29,
	"tomato paste":      0.99,
	"marinara sauce":    3.49,
	"chicken broth":     2.49,
	"beef broth":        2.49,
	"broth":             2.49,
	"kidney beans":      1.29,
	"pumpkin":           2.99,

	// Produce
	"onion":       0.99,
	"onions":      0.99,
	"garlic":      0.69,
	"tomato":      1.29,
	"tomatoes":    1.29,
	"lettuce":     1.99,
	"bell pepper": 1.29,
	"broccoli":    2.49,
	"snap peas":   2.99,
	"celery":      1.99,
	"potatoes":    3.99,
	"potato":      0.99,
	"green onions": 0.99,
	"cilantro":    0.99,
	"parsley":     0.99,
	"lemon":       0.69,
	"lemon juice": 2.49,
	"ginger":      0.99,
	"green chili":  0.49,
	"jalapenos":   0.49,
	"lime":        0.49,

	// Oils & Condiments
	"olive oil":        6.99,
	"vegetable oil":    3.99,
	"sesame oil":       4.49,
	"mayo":             4.49,
	"mayonnaise":       4.49,
	"dijon mustard":    3.49,
	"mustard":          2.49,
	"yellow mustard":   2.49,
	"soy sauce":        3.49,
	"oyster sauce":     3.99,
	"maple syrup":      6.99,
	"honey":            5.99,
	"salsa":            3.49,
	"hot sauce":        2.99,
	"ketchup":          3.49,
	"vanilla extract":  4.99,
	"vanilla":          4.99,

	// Baking
	"sugar":            3.49,
	"brown sugar":      3.49,
	"baking powder":    2.99,
	"baking soda":      1.49,
	"cornstarch":       2.49,
	"chia seeds":       5.99,

	// Spices
	"salt":             1.99,
	"kosher salt":      3.49,
	"pepper":           3.99,
	"black pepper":     3.99,
	"cinnamon":         3.49,
	"nutmeg":           4.49,
	"cloves":           4.49,
	"ground cloves":    4.49,
	"paprika":          3.49,
	"chili powder":     3.49,
	"cumin":            3.49,
	"ground cumin":     3.49,
	"garlic powder":    3.49,
	"onion powder":     3.49,
	"cayenne":          3.49,
	"red pepper flakes": 3.49,
	"italian seasoning": 3.49,
	"italian herbs":    3.49,
	"oregano":          3.49,
	"dried oregano":    3.49,
	"thyme":            2.99,
	"rosemary":         2.99,
	"bay leaves":       3.49,
	"cumin seeds":      3.49,
	"asafoetida":       5.99,
	"coriander powder": 3.49,
	"turmeric":         3.49,
	"red chili powder":  3.49,

	// Other
	"corn husks":       3.99,
	"lard":             3.99,
	"granola":          4.99,
	"water":            0.00,
	"cold water":       0.00,
}

func EstimatePrice(ingredientName string) float64 {
	name := strings.ToLower(strings.TrimSpace(ingredientName))

	if price, ok := priceTable[name]; ok {
		return price
	}

	for key, price := range priceTable {
		if strings.Contains(name, key) || strings.Contains(key, name) {
			return price
		}
	}

	return 2.99
}

var iconTable = map[string]string{
	// Proteins
	"ground beef": "🥩", "lean ground beef": "🥩", "beef": "🥩", "brisket": "🥩",
	"chicken breast": "🍗", "chicken thighs": "🍗", "chicken thigh": "🍗", "chicken": "🍗",
	"pork loin": "🥩", "pork chops": "🥩", "pork chop": "🥩", "pork": "🥩",
	"bacon": "🥓", "pancetta": "🥓",
	"breakfast sausage": "🌭", "sausage": "🌭",
	"tuna": "🐟", "shrimp": "🦐",

	// Dairy & Eggs
	"egg": "🥚", "eggs": "🥚", "egg yolk": "🥚",
	"butter": "🧈", "unsalted butter": "🧈", "cold butter": "🧈",
	"milk": "🥛", "whole milk": "🥛", "cream": "🥛",
	"cheese": "🧀", "cheddar": "🧀", "cheddar cheese": "🧀", "american cheese": "🧀",
	"swiss cheese": "🧀", "parmesan": "🧀", "parmesan cheese": "🧀",
	"sour cream": "🥛", "yogurt": "🥛",

	// Grains & Pasta
	"flour": "🌾", "all-purpose flour": "🌾", "bread flour": "🌾",
	"bread": "🍞", "tortillas": "🫓", "flour tortillas": "🫓", "tortilla": "🫓",
	"taco shells": "🌮",
	"pasta": "🍝", "spaghetti": "🍝", "elbow macaroni": "🍝", "macaroni": "🍝",
	"rice": "🍚", "rolled oats": "🌾", "oats": "🌾",
	"biscuits": "🍞", "masa harina": "🌽",

	// Produce
	"onion": "🧅", "onions": "🧅", "red onion": "🧅",
	"garlic": "🧄",
	"tomato": "🍅", "tomatoes": "🍅",
	"lettuce": "🥬",
	"bell pepper": "🫑", "green chili": "🌶️", "jalapenos": "🌶️",
	"broccoli": "🥦",
	"potato": "🥔", "potatoes": "🥔",
	"celery": "🥬",
	"snap peas": "🫛",
	"ginger": "🫚",
	"lemon": "🍋", "lime": "🍋",
	"green onions": "🧅",
	"cilantro": "🌿", "parsley": "🌿",
	"pumpkin": "🎃",

	// Oils & Condiments
	"olive oil": "🫒", "vegetable oil": "🫒", "sesame oil": "🫒",
	"mayo": "🫙", "mayonnaise": "🫙",
	"mustard": "🫙", "dijon mustard": "🫙", "yellow mustard": "🫙",
	"soy sauce": "🫙", "oyster sauce": "🫙",
	"maple syrup": "🍁", "honey": "🍯",
	"salsa": "🫙", "hot sauce": "🌶️", "ketchup": "🫙",
	"vanilla extract": "🫙", "vanilla": "🫙",

	// Baking
	"sugar": "🧂", "brown sugar": "🧂",
	"baking powder": "🧂", "baking soda": "🧂", "cornstarch": "🧂",

	// Spices
	"salt": "🧂", "kosher salt": "🧂",
	"pepper": "🧂", "black pepper": "🧂",
	"cinnamon": "🧂", "nutmeg": "🧂", "cloves": "🧂", "ground cloves": "🧂",
	"paprika": "🧂", "chili powder": "🌶️", "cumin": "🧂", "ground cumin": "🧂",
	"garlic powder": "🧂", "onion powder": "🧂", "cayenne": "🌶️",
	"red pepper flakes": "🌶️", "italian seasoning": "🌿", "italian herbs": "🌿",
	"oregano": "🌿", "dried oregano": "🌿", "thyme": "🌿", "rosemary": "🌿",
	"bay leaves": "🌿", "turmeric": "🧂",

	// Other
	"corn husks": "🌽", "lard": "🫙", "water": "💧", "cold water": "💧",
	"granola": "🥣", "chia seeds": "🌱",
}

func IngredientIcon(ingredientName string) string {
	name := strings.ToLower(strings.TrimSpace(ingredientName))

	if icon, ok := iconTable[name]; ok {
		return icon
	}

	for key, icon := range iconTable {
		if strings.Contains(name, key) || strings.Contains(key, name) {
			return icon
		}
	}

	return "🛒"
}

type PricedShoppingItem struct {
	Name    string
	Amounts []string
	Price   float64
	Icon    string
}

func PriceShoppingList(list *ShoppingList) ([]PricedShoppingItem, float64) {
	var items []PricedShoppingItem
	var total float64

	for _, item := range list.Items {
		price := EstimatePrice(item.Name)
		items = append(items, PricedShoppingItem{
			Name:    item.Name,
			Amounts: item.Amounts,
			Price:   price,
			Icon:    IngredientIcon(item.Name),
		})
		total += price
	}

	return items, total
}
