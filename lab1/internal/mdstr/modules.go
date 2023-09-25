package mdstr

type Card struct {
	Index       int
	Title       string
	Image       string
	Description string
	LaunchDate  string
}

func GetModule() []Card {
	return []Card{
		{Index: 1, Title: "Электродвигательный модуль", Image: "img1.jpeg", Description: "Модуль, который обеспечивает движение космического корабля.", LaunchDate: "Ноябрь, 2025 год"},
		{Index: 2, Title: "Модуль ESPRIT", Image: "img2.jpg", Description: "Модуль, предназначенный для транспортировки и хранения грузов, а также для дозаправки станции.", LaunchDate: "2029 год"},
		{Index: 3, Title: "Модуль снабжения", Image: "img3.jpg", Description: "Модуль, который обеспечивает снабжение космического корабля необходимыми ресурсами.", LaunchDate: "2024 год"},
		{Index: 4, Title: "Аванпост жилья и логистики", Image: "img4.jpg", Description: "Также называемый модулем минимального жилья и ранее известный как утилизационный модуль. Представляет собой уменьшенный жилой модуль станции.", LaunchDate: "Ноябрь, 2025 год"},
		{Index: 5, Title: "Международный жилой модуль", Image: "img5.jpg", Description: "Жилой модуль, создаваемый международными партнёрами США.", LaunchDate: "2028 год"},
		{Index: 6, Title: "Американский жилой модуль", Image: "img6.jpg", Description: "Жилой модуль, создаваемый США.", LaunchDate: "Неизвестно"},
		{Index: 7, Title: "Шлюзовой модуль", Image: "img7.jpg", Description: "Модуль, предназначенный для выполнения внекорабельных действий за пределами станции.", LaunchDate: "2030 год"},
		{Index: 8, Title: `Пилотируемый модуль "Орион"`, Image: "img8.jpg", Description: "Многоразовый космический корабль, предназначенный для доставки экипажа на станцию.", LaunchDate: "Постоянно"},
	}
}
