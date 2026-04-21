package recipe

import "strings"

type Department int

const (
	DeptProduce Department = iota
	DeptMeat
	DeptDairy
	DeptBakery
	DeptGrains
	DeptCanned
	DeptCondiments
	DeptSpices
	DeptBaking
	DeptOther
)

func (d Department) String() string {
	switch d {
	case DeptProduce:
		return "Produce"
	case DeptMeat:
		return "Meat & Seafood"
	case DeptDairy:
		return "Dairy & Eggs"
	case DeptBakery:
		return "Bakery"
	case DeptGrains:
		return "Grains & Pasta"
	case DeptCanned:
		return "Canned Goods"
	case DeptCondiments:
		return "Condiments & Sauces"
	case DeptSpices:
		return "Spices & Seasonings"
	case DeptBaking:
		return "Baking"
	case DeptOther:
		return "Other"
	}
	return "Other"
}

var DepartmentOrder = []Department{
	DeptProduce,
	DeptMeat,
	DeptDairy,
	DeptBakery,
	DeptGrains,
	DeptCanned,
	DeptCondiments,
	DeptSpices,
	DeptBaking,
	DeptOther,
}

type ingredientInfo struct {
	Price      float64
	Icon       string
	Department Department
	BuyUnit    string
}

var ingredientTable = map[string]ingredientInfo{
	// Proteins
	"ground beef":       {5.99, "🥩", DeptMeat, "1 lb"},
	"lean ground beef":  {6.49, "🥩", DeptMeat, "1 lb"},
	"beef":              {5.99, "🥩", DeptMeat, "1 lb"},
	"brisket":           {8.99, "🥩", DeptMeat, "per lb"},
	"chicken breast":    {7.99, "🍗", DeptMeat, "1.5 lb / 3-4 breasts"},
	"chicken thighs":    {5.49, "🍗", DeptMeat, "2 lb / 6-8 thighs"},
	"chicken thigh":     {5.49, "🍗", DeptMeat, "2 lb / 6-8 thighs"},
	"chicken":           {5.49, "🍗", DeptMeat, "per lb"},
	"pork loin":         {12.99, "🥩", DeptMeat, "3-4 lb roast"},
	"pork chops":        {7.99, "🥩", DeptMeat, "4 chops"},
	"pork chop":         {7.99, "🥩", DeptMeat, "4 chops"},
	"pork":              {7.99, "🥩", DeptMeat, "per lb"},
	"bacon":             {6.99, "🥓", DeptMeat, "1 lb / 12-16 slices"},
	"pancetta":          {4.99, "🥓", DeptMeat, "4 oz pkg"},
	"breakfast sausage": {4.99, "🌭", DeptMeat, "1 lb roll"},
	"sausage":           {4.99, "🌭", DeptMeat, "1 lb"},
	"tuna":              {1.99, "🐟", DeptMeat, "5 oz can"},
	"shrimp":            {9.99, "🦐", DeptMeat, "1 lb / 21-25 ct"},

	// Dairy & Eggs
	"eggs":              {4.49, "🥚", DeptDairy, "1 dozen"},
	"egg":               {4.49, "🥚", DeptDairy, "1 dozen"},
	"egg yolk":          {4.49, "🥚", DeptDairy, "1 dozen"},
	"butter":            {4.99, "🧈", DeptDairy, "1 lb / 4 sticks"},
	"unsalted butter":   {4.99, "🧈", DeptDairy, "1 lb / 4 sticks"},
	"cold butter":       {4.99, "🧈", DeptDairy, "1 lb / 4 sticks"},
	"milk":              {3.99, "🥛", DeptDairy, "1 gallon"},
	"whole milk":        {3.99, "🥛", DeptDairy, "1 gallon"},
	"cream":             {4.49, "🥛", DeptDairy, "1 pint"},
	"sour cream":        {2.49, "🥛", DeptDairy, "16 oz"},
	"yogurt":            {3.99, "🥛", DeptDairy, "32 oz"},
	"cheddar cheese":    {4.49, "🧀", DeptDairy, "8 oz block"},
	"cheddar":           {4.49, "🧀", DeptDairy, "8 oz block"},
	"american cheese":   {3.99, "🧀", DeptDairy, "12 slices"},
	"swiss cheese":      {4.49, "🧀", DeptDairy, "8 oz block"},
	"cheese":            {4.49, "🧀", DeptDairy, "8 oz"},
	"parmesan":          {5.99, "🧀", DeptDairy, "5 oz wedge"},
	"parmesan cheese":   {5.99, "🧀", DeptDairy, "5 oz wedge"},

	// Bakery
	"bread":    {3.49, "🍞", DeptBakery, "1 loaf"},
	"biscuits": {2.99, "🍞", DeptBakery, "8 ct can"},

	// Grains & Pasta
	"flour":             {3.99, "🌾", DeptGrains, "5 lb bag"},
	"all-purpose flour": {3.99, "🌾", DeptGrains, "5 lb bag"},
	"bread flour":       {4.49, "🌾", DeptGrains, "5 lb bag"},
	"elbow macaroni":    {1.49, "🍝", DeptGrains, "1 lb box"},
	"macaroni":          {1.49, "🍝", DeptGrains, "1 lb box"},
	"spaghetti":         {1.49, "🍝", DeptGrains, "1 lb box"},
	"pasta":             {1.49, "🍝", DeptGrains, "1 lb box"},
	"tortillas":         {3.49, "🫓", DeptGrains, "10 ct pkg"},
	"flour tortillas":   {3.49, "🫓", DeptGrains, "10 ct pkg"},
	"tortilla":          {3.49, "🫓", DeptGrains, "10 ct pkg"},
	"taco shells":       {2.49, "🌮", DeptGrains, "12 ct box"},
	"rolled oats":       {3.99, "🌾", DeptGrains, "42 oz canister"},
	"oats":              {3.99, "🌾", DeptGrains, "42 oz canister"},
	"masa harina":       {3.99, "🌽", DeptGrains, "4 lb bag"},
	"rice":              {2.99, "🍚", DeptGrains, "2 lb bag"},

	// Canned goods
	"crushed tomatoes": {1.99, "🍅", DeptCanned, "28 oz can"},
	"diced tomatoes":   {1.49, "🍅", DeptCanned, "14.5 oz can"},
	"tomato sauce":     {1.29, "🍅", DeptCanned, "15 oz can"},
	"tomato paste":     {0.99, "🍅", DeptCanned, "6 oz can"},
	"marinara sauce":   {3.49, "🍅", DeptCanned, "24 oz jar"},
	"chicken broth":    {2.49, "🫙", DeptCanned, "32 oz carton"},
	"beef broth":       {2.49, "🫙", DeptCanned, "32 oz carton"},
	"broth":            {2.49, "🫙", DeptCanned, "32 oz carton"},
	"kidney beans":     {1.29, "🫙", DeptCanned, "15 oz can"},
	"pumpkin":          {2.99, "🎃", DeptCanned, "15 oz can"},

	// Produce
	"onion":        {0.99, "🧅", DeptProduce, "1 each"},
	"onions":       {0.99, "🧅", DeptProduce, "3 lb bag"},
	"red onion":    {0.99, "🧅", DeptProduce, "1 each"},
	"garlic":       {0.69, "🧄", DeptProduce, "1 head"},
	"tomato":       {1.29, "🍅", DeptProduce, "1 each"},
	"tomatoes":     {1.29, "🍅", DeptProduce, "per lb"},
	"lettuce":      {1.99, "🥬", DeptProduce, "1 head"},
	"bell pepper":  {1.29, "🫑", DeptProduce, "1 each"},
	"broccoli":     {2.49, "🥦", DeptProduce, "1 crown"},
	"snap peas":    {2.99, "🫛", DeptProduce, "8 oz bag"},
	"celery":       {1.99, "🥬", DeptProduce, "1 bunch"},
	"potatoes":     {3.99, "🥔", DeptProduce, "5 lb bag"},
	"potato":       {0.99, "🥔", DeptProduce, "1 each"},
	"green onions": {0.99, "🧅", DeptProduce, "1 bunch"},
	"cilantro":     {0.99, "🌿", DeptProduce, "1 bunch"},
	"parsley":      {0.99, "🌿", DeptProduce, "1 bunch"},
	"lemon":        {0.69, "🍋", DeptProduce, "1 each"},
	"lemon juice":  {2.49, "🍋", DeptProduce, "8 oz bottle"},
	"ginger":       {0.99, "🫚", DeptProduce, "1 knob"},
	"green chili":  {0.49, "🌶️", DeptProduce, "1 each"},
	"jalapenos":    {0.49, "🌶️", DeptProduce, "1 each"},
	"lime":         {0.49, "🍋", DeptProduce, "1 each"},

	// Condiments & Sauces
	"olive oil":       {6.99, "🫒", DeptCondiments, "16 oz bottle"},
	"vegetable oil":   {3.99, "🫒", DeptCondiments, "48 oz bottle"},
	"sesame oil":      {4.49, "🫒", DeptCondiments, "5 oz bottle"},
	"mayo":            {4.49, "🫙", DeptCondiments, "30 oz jar"},
	"mayonnaise":      {4.49, "🫙", DeptCondiments, "30 oz jar"},
	"dijon mustard":   {3.49, "🫙", DeptCondiments, "12 oz jar"},
	"mustard":         {2.49, "🫙", DeptCondiments, "14 oz bottle"},
	"yellow mustard":  {2.49, "🫙", DeptCondiments, "14 oz bottle"},
	"soy sauce":       {3.49, "🫙", DeptCondiments, "15 oz bottle"},
	"oyster sauce":    {3.99, "🫙", DeptCondiments, "9 oz bottle"},
	"maple syrup":     {6.99, "🍁", DeptCondiments, "12 oz bottle"},
	"honey":           {5.99, "🍯", DeptCondiments, "12 oz bottle"},
	"salsa":           {3.49, "🫙", DeptCondiments, "16 oz jar"},
	"hot sauce":       {2.99, "🌶️", DeptCondiments, "5 oz bottle"},
	"ketchup":         {3.49, "🫙", DeptCondiments, "20 oz bottle"},
	"vanilla extract": {4.99, "🫙", DeptBaking, "2 oz bottle"},
	"vanilla":         {4.99, "🫙", DeptBaking, "2 oz bottle"},

	// Baking
	"sugar":         {3.49, "🧂", DeptBaking, "4 lb bag"},
	"brown sugar":   {3.49, "🧂", DeptBaking, "2 lb bag"},
	"baking powder": {2.99, "🧂", DeptBaking, "8.1 oz can"},
	"baking soda":   {1.49, "🧂", DeptBaking, "16 oz box"},
	"cornstarch":    {2.49, "🧂", DeptBaking, "16 oz box"},
	"chia seeds":    {5.99, "🌱", DeptBaking, "12 oz bag"},

	// Spices
	"salt":              {1.99, "🧂", DeptSpices, "26 oz canister"},
	"kosher salt":       {3.49, "🧂", DeptSpices, "3 lb box"},
	"pepper":            {3.99, "🧂", DeptSpices, "2 oz grinder"},
	"black pepper":      {3.99, "🧂", DeptSpices, "2 oz grinder"},
	"cinnamon":          {3.49, "🧂", DeptSpices, "2.37 oz jar"},
	"nutmeg":            {4.49, "🧂", DeptSpices, "1.1 oz jar"},
	"cloves":            {4.49, "🧂", DeptSpices, "0.9 oz jar"},
	"ground cloves":     {4.49, "🧂", DeptSpices, "0.9 oz jar"},
	"paprika":           {3.49, "🧂", DeptSpices, "2.1 oz jar"},
	"chili powder":      {3.49, "🌶️", DeptSpices, "2.5 oz jar"},
	"cumin":             {3.49, "🧂", DeptSpices, "1.7 oz jar"},
	"ground cumin":      {3.49, "🧂", DeptSpices, "1.7 oz jar"},
	"garlic powder":     {3.49, "🧂", DeptSpices, "3.1 oz jar"},
	"onion powder":      {3.49, "🧂", DeptSpices, "2.6 oz jar"},
	"cayenne":           {3.49, "🌶️", DeptSpices, "1.75 oz jar"},
	"red pepper flakes": {3.49, "🌶️", DeptSpices, "1.5 oz jar"},
	"italian seasoning": {3.49, "🌿", DeptSpices, "0.87 oz jar"},
	"italian herbs":     {3.49, "🌿", DeptSpices, "0.87 oz jar"},
	"oregano":           {3.49, "🌿", DeptSpices, "0.75 oz jar"},
	"dried oregano":     {3.49, "🌿", DeptSpices, "0.75 oz jar"},
	"thyme":             {2.99, "🌿", DeptSpices, "0.63 oz jar"},
	"rosemary":          {2.99, "🌿", DeptSpices, "0.75 oz jar"},
	"bay leaves":        {3.49, "🌿", DeptSpices, "0.12 oz jar"},
	"cumin seeds":       {3.49, "🧂", DeptSpices, "1.7 oz jar"},
	"asafoetida":        {5.99, "🧂", DeptSpices, "3.5 oz jar"},
	"coriander powder":  {3.49, "🧂", DeptSpices, "1.5 oz jar"},
	"turmeric":          {3.49, "🧂", DeptSpices, "1.8 oz jar"},
	"red chili powder":  {3.49, "🌶️", DeptSpices, "2.5 oz jar"},

	// Other
	"corn husks": {3.99, "🌽", DeptOther, "1 pkg"},
	"lard":       {3.99, "🫙", DeptOther, "1 lb"},
	"granola":    {4.99, "🥣", DeptOther, "12 oz bag"},
	"water":      {0.00, "💧", DeptOther, ""},
	"cold water": {0.00, "💧", DeptOther, ""},
}

func lookupIngredient(name string) (ingredientInfo, bool) {
	lower := strings.ToLower(strings.TrimSpace(name))

	if info, ok := ingredientTable[lower]; ok {
		return info, true
	}

	for key, info := range ingredientTable {
		if strings.Contains(lower, key) || strings.Contains(key, lower) {
			return info, true
		}
	}

	return ingredientInfo{}, false
}

func EstimatePrice(ingredientName string) float64 {
	if info, ok := lookupIngredient(ingredientName); ok {
		return info.Price
	}
	return 2.99
}

func IngredientIcon(ingredientName string) string {
	if info, ok := lookupIngredient(ingredientName); ok {
		return info.Icon
	}
	return "🛒"
}

func IngredientDepartment(ingredientName string) Department {
	if info, ok := lookupIngredient(ingredientName); ok {
		return info.Department
	}
	return DeptOther
}

func IngredientBuyUnit(ingredientName string) string {
	if info, ok := lookupIngredient(ingredientName); ok {
		return info.BuyUnit
	}
	return ""
}

var defaultPantryStaples = map[string]bool{
	"salt":              true,
	"kosher salt":       true,
	"pepper":            true,
	"black pepper":      true,
	"olive oil":         true,
	"vegetable oil":     true,
	"flour":             true,
	"all-purpose flour": true,
	"sugar":             true,
	"brown sugar":       true,
	"baking powder":     true,
	"baking soda":       true,
	"cornstarch":        true,
	"garlic powder":     true,
	"onion powder":      true,
	"paprika":           true,
	"chili powder":      true,
	"cumin":             true,
	"ground cumin":      true,
	"cayenne":           true,
	"red pepper flakes": true,
	"italian seasoning": true,
	"oregano":           true,
	"dried oregano":     true,
	"cinnamon":          true,
	"bay leaves":        true,
	"soy sauce":         true,
	"vanilla extract":   true,
	"vanilla":           true,
	"honey":             true,
	"vinegar":           true,
	"water":             true,
	"cold water":        true,
}

func IsPantryStaple(ingredientName string) bool {
	name := strings.ToLower(strings.TrimSpace(ingredientName))
	return defaultPantryStaples[name]
}

type PricedShoppingItem struct {
	Name        string
	Amounts     []string
	Price       float64
	Icon        string
	Department  string
	BuyUnit     string
	IsPantry    bool
	IsPerishable bool
	Trip        int
}

type ShoppingDepartment struct {
	Name  string
	Items []PricedShoppingItem
}

func IsPerishableDepartment(dept Department) bool {
	return dept == DeptProduce || dept == DeptMeat || dept == DeptDairy
}

func PriceShoppingList(list *ShoppingList) ([]PricedShoppingItem, float64) {
	var items []PricedShoppingItem
	var total float64

	for _, item := range list.Items {
		price := EstimatePrice(item.Name)
		pantry := IsPantryStaple(item.Name)
		dept := IngredientDepartment(item.Name)
		items = append(items, PricedShoppingItem{
			Name:         item.Name,
			Amounts:      item.Amounts,
			Price:        price,
			Icon:         IngredientIcon(item.Name),
			Department:   dept.String(),
			BuyUnit:      IngredientBuyUnit(item.Name),
			IsPantry:     pantry,
			IsPerishable: IsPerishableDepartment(dept),
		})
		if !pantry {
			total += price
		}
	}

	return items, total
}

func GroupByDepartment(items []PricedShoppingItem) []ShoppingDepartment {
	grouped := make(map[string][]PricedShoppingItem)
	for _, item := range items {
		grouped[item.Department] = append(grouped[item.Department], item)
	}

	var departments []ShoppingDepartment
	for _, dept := range DepartmentOrder {
		name := dept.String()
		if items, ok := grouped[name]; ok {
			departments = append(departments, ShoppingDepartment{
				Name:  name,
				Items: items,
			})
		}
	}

	return departments
}
