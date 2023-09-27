package mdstr

type Module struct {
	Index       int
	Title       string
	Image       string
	Description string
	Mass        string
	Diameter    string
	Length      string
	LaunchDate  string
}

func GetModule() []Module {
	return []Module{
		{Index: 1, Title: "Электродвигательный модуль (PPE)", Image: "img1.jpeg", Description: "Модуль, который обеспечивает движение космического корабля.", Mass: "5 тонн", Length: "11 метров", Diameter: "4,5 метра", LaunchDate: "Ноябрь, 2025 год"},
		{Index: 2, Title: "Модуль ESPRIT", Image: "img2.jpg", Description: "Модуль, предназначенный для транспортировки и хранения грузов, а также для дозаправки станции.", Mass: "10 тонн", Length: "12 метров", Diameter: "3,91 метра", LaunchDate: "2029 год"},
		{Index: 3, Title: "Модуль снабжения", Image: "img3.jpg", Description: "Модуль, который обеспечивает снабжение космического корабля необходимыми ресурсами.", Mass: "10,4 тонны", Length: "9 метров", Diameter: "4,5 метра", LaunchDate: "2024 год"},
		{Index: 4, Title: "Аванпост жилья и логистики (HALO)", Image: "img4.jpg", Description: "Также называемый модулем минимального жилья и ранее известный как утилизационный модуль. Представляет собой уменьшенный жилой модуль станции.", Mass: "6,6 тонн", Length: "10 метров", Diameter: "4 метра", LaunchDate: "Ноябрь, 2025 год"},
		{Index: 5, Title: "Международный жилой модуль (I-HAB)", Image: "img5.jpg", Description: "Жилой модуль, создаваемый международными партнёрами США.", Mass: "10 тонн", Length: "13 метров", Diameter: "5,4 метра", LaunchDate: "2028 год"},
		{Index: 6, Title: "Американский жилой модуль (US-HB)", Image: "img6.jpg", Description: "Жилой модуль, создаваемый США.", Mass: "11 тонн", Length: "13 метров", Diameter: "5 метров", LaunchDate: "2028 год"},
		{Index: 7, Title: "Шлюзовой модуль", Image: "img7.jpg", Description: "Модуль, предназначенный для выполнения внекорабельных действий за пределами станции.", Mass: "5 тонн", Length: "5 метров", Diameter: "4,5 метра", LaunchDate: "2030 год"},
		{Index: 8, Title: `Пилотируемый корабль "Орион"`, Image: "img8.jpg", Description: "Многоразовый космический корабль, предназначенный для доставки экипажа на станцию.", Mass: "25 тонн", Length: "3,3 метра", Diameter: "5,3 метра", LaunchDate: "Постоянно"},
	}
}
