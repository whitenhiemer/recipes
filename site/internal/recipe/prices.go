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
}

var ingredientTable = map[string]ingredientInfo{
	// Proteins
	"ground beef":       {5.99, "🥩", DeptMeat},
	"lean ground beef":  {6.49, "🥩", DeptMeat},
	"beef":              {5.99, "🥩", DeptMeat},
	"brisket":           {8.99, "🥩", DeptMeat},
	"chicken breast":    {7.99, "🍗", DeptMeat},
	"chicken thighs":    {5.49, "🍗", DeptMeat},
	"chicken thigh":     {5.49, "🍗", DeptMeat},
	"chicken":           {5.49, "🍗", DeptMeat},
	"pork loin":         {12.99, "🥩", DeptMeat},
	"pork chops":        {7.99, "🥩", DeptMeat},
	"pork chop":         {7.99, "🥩", DeptMeat},
	"pork":              {7.99, "🥩", DeptMeat},
	"bacon":             {6.99, "🥓", DeptMeat},
	"pancetta":          {4.99, "🥓", DeptMeat},
	"breakfast sausage": {4.99, "🌭", DeptMeat},
	"sausage":           {4.99, "🌭", DeptMeat},
	"tuna":              {1.99, "🐟", DeptMeat},
	"shrimp":            {9.99, "🦐", DeptMeat},

	// Dairy & Eggs
	"eggs":              {4.49, "🥚", DeptDairy},
	"egg":               {4.49, "🥚", DeptDairy},
	"egg yolk":          {4.49, "🥚", DeptDairy},
	"butter":            {4.99, "🧈", DeptDairy},
	"unsalted butter":   {4.99, "🧈", DeptDairy},
	"cold butter":       {4.99, "🧈", DeptDairy},
	"milk":              {3.99, "🥛", DeptDairy},
	"whole milk":        {3.99, "🥛", DeptDairy},
	"cream":             {4.49, "🥛", DeptDairy},
	"sour cream":        {2.49, "🥛", DeptDairy},
	"yogurt":            {3.99, "🥛", DeptDairy},
	"cheddar cheese":    {4.49, "🧀", DeptDairy},
	"cheddar":           {4.49, "🧀", DeptDairy},
	"american cheese":   {3.99, "🧀", DeptDairy},
	"swiss cheese":      {4.49, "🧀", DeptDairy},
	"cheese":            {4.49, "🧀", DeptDairy},
	"parmesan":          {5.99, "🧀", DeptDairy},
	"parmesan cheese":   {5.99, "🧀", DeptDairy},

	// Bakery
	"bread":    {3.49, "🍞", DeptBakery},
	"biscuits": {2.99, "🍞", DeptBakery},

	// Grains & Pasta
	"flour":            {3.99, "🌾", DeptGrains},
	"all-purpose flour": {3.99, "🌾", DeptGrains},
	"bread flour":      {4.49, "🌾", DeptGrains},
	"elbow macaroni":   {1.49, "🍝", DeptGrains},
	"macaroni":         {1.49, "🍝", DeptGrains},
	"spaghetti":        {1.49, "🍝", DeptGrains},
	"pasta":            {1.49, "🍝", DeptGrains},
	"tortillas":        {3.49, "🫓", DeptGrains},
	"flour tortillas":  {3.49, "🫓", DeptGrains},
	"tortilla":         {3.49, "🫓", DeptGrains},
	"taco shells":      {2.49, "🌮", DeptGrains},
	"rolled oats":      {3.99, "🌾", DeptGrains},
	"oats":             {3.99, "🌾", DeptGrains},
	"masa harina":      {3.99, "🌽", DeptGrains},
	"rice":             {2.99, "🍚", DeptGrains},

	// Canned goods
	"crushed tomatoes": {1.99, "🍅", DeptCanned},
	"diced tomatoes":   {1.49, "🍅", DeptCanned},
	"tomato sauce":     {1.29, "🍅", DeptCanned},
	"tomato paste":     {0.99, "🍅", DeptCanned},
	"marinara sauce":   {3.49, "🍅", DeptCanned},
	"chicken broth":    {2.49, "🫙", DeptCanned},
	"beef broth":       {2.49, "🫙", DeptCanned},
	"broth":            {2.49, "🫙", DeptCanned},
	"kidney beans":     {1.29, "🫙", DeptCanned},
	"pumpkin":          {2.99, "🎃", DeptCanned},

	// Produce
	"onion":        {0.99, "🧅", DeptProduce},
	"onions":       {0.99, "🧅", DeptProduce},
	"red onion":    {0.99, "🧅", DeptProduce},
	"garlic":       {0.69, "🧄", DeptProduce},
	"tomato":       {1.29, "🍅", DeptProduce},
	"tomatoes":     {1.29, "🍅", DeptProduce},
	"lettuce":      {1.99, "🥬", DeptProduce},
	"bell pepper":  {1.29, "🫑", DeptProduce},
	"broccoli":     {2.49, "🥦", DeptProduce},
	"snap peas":    {2.99, "🫛", DeptProduce},
	"celery":       {1.99, "🥬", DeptProduce},
	"potatoes":     {3.99, "🥔", DeptProduce},
	"potato":       {0.99, "🥔", DeptProduce},
	"green onions": {0.99, "🧅", DeptProduce},
	"cilantro":     {0.99, "🌿", DeptProduce},
	"parsley":      {0.99, "🌿", DeptProduce},
	"lemon":        {0.69, "🍋", DeptProduce},
	"lemon juice":  {2.49, "🍋", DeptProduce},
	"ginger":       {0.99, "🫚", DeptProduce},
	"green chili":  {0.49, "🌶️", DeptProduce},
	"jalapenos":    {0.49, "🌶️", DeptProduce},
	"lime":         {0.49, "🍋", DeptProduce},

	// Condiments & Sauces
	"olive oil":       {6.99, "🫒", DeptCondiments},
	"vegetable oil":   {3.99, "🫒", DeptCondiments},
	"sesame oil":      {4.49, "🫒", DeptCondiments},
	"mayo":            {4.49, "🫙", DeptCondiments},
	"mayonnaise":      {4.49, "🫙", DeptCondiments},
	"dijon mustard":   {3.49, "🫙", DeptCondiments},
	"mustard":         {2.49, "🫙", DeptCondiments},
	"yellow mustard":  {2.49, "🫙", DeptCondiments},
	"soy sauce":       {3.49, "🫙", DeptCondiments},
	"oyster sauce":    {3.99, "🫙", DeptCondiments},
	"maple syrup":     {6.99, "🍁", DeptCondiments},
	"honey":           {5.99, "🍯", DeptCondiments},
	"salsa":           {3.49, "🫙", DeptCondiments},
	"hot sauce":       {2.99, "🌶️", DeptCondiments},
	"ketchup":         {3.49, "🫙", DeptCondiments},
	"vanilla extract": {4.99, "🫙", DeptBaking},
	"vanilla":         {4.99, "🫙", DeptBaking},

	// Baking
	"sugar":         {3.49, "🧂", DeptBaking},
	"brown sugar":   {3.49, "🧂", DeptBaking},
	"baking powder": {2.99, "🧂", DeptBaking},
	"baking soda":   {1.49, "🧂", DeptBaking},
	"cornstarch":    {2.49, "🧂", DeptBaking},
	"chia seeds":    {5.99, "🌱", DeptBaking},

	// Spices
	"salt":              {1.99, "🧂", DeptSpices},
	"kosher salt":       {3.49, "🧂", DeptSpices},
	"pepper":            {3.99, "🧂", DeptSpices},
	"black pepper":      {3.99, "🧂", DeptSpices},
	"cinnamon":          {3.49, "🧂", DeptSpices},
	"nutmeg":            {4.49, "🧂", DeptSpices},
	"cloves":            {4.49, "🧂", DeptSpices},
	"ground cloves":     {4.49, "🧂", DeptSpices},
	"paprika":           {3.49, "🧂", DeptSpices},
	"chili powder":      {3.49, "🌶️", DeptSpices},
	"cumin":             {3.49, "🧂", DeptSpices},
	"ground cumin":      {3.49, "🧂", DeptSpices},
	"garlic powder":     {3.49, "🧂", DeptSpices},
	"onion powder":      {3.49, "🧂", DeptSpices},
	"cayenne":           {3.49, "🌶️", DeptSpices},
	"red pepper flakes": {3.49, "🌶️", DeptSpices},
	"italian seasoning": {3.49, "🌿", DeptSpices},
	"italian herbs":     {3.49, "🌿", DeptSpices},
	"oregano":           {3.49, "🌿", DeptSpices},
	"dried oregano":     {3.49, "🌿", DeptSpices},
	"thyme":             {2.99, "🌿", DeptSpices},
	"rosemary":          {2.99, "🌿", DeptSpices},
	"bay leaves":        {3.49, "🌿", DeptSpices},
	"cumin seeds":       {3.49, "🧂", DeptSpices},
	"asafoetida":        {5.99, "🧂", DeptSpices},
	"coriander powder":  {3.49, "🧂", DeptSpices},
	"turmeric":          {3.49, "🧂", DeptSpices},
	"red chili powder":  {3.49, "🌶️", DeptSpices},

	// Other
	"corn husks": {3.99, "🌽", DeptOther},
	"lard":       {3.99, "🫙", DeptOther},
	"granola":    {4.99, "🥣", DeptOther},
	"water":      {0.00, "💧", DeptOther},
	"cold water": {0.00, "💧", DeptOther},
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
	Name       string
	Amounts    []string
	Price      float64
	Icon       string
	Department string
	IsPantry   bool
}

type ShoppingDepartment struct {
	Name  string
	Items []PricedShoppingItem
}

func PriceShoppingList(list *ShoppingList) ([]PricedShoppingItem, float64) {
	var items []PricedShoppingItem
	var total float64

	for _, item := range list.Items {
		price := EstimatePrice(item.Name)
		pantry := IsPantryStaple(item.Name)
		items = append(items, PricedShoppingItem{
			Name:       item.Name,
			Amounts:    item.Amounts,
			Price:      price,
			Icon:       IngredientIcon(item.Name),
			Department: IngredientDepartment(item.Name).String(),
			IsPantry:   pantry,
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
