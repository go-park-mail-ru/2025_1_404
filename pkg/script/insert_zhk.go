//go:build script

package main

import (
	"context"
	"log"

	"github.com/go-park-mail-ru/2025_1_404/pkg/database/postgres"
	"github.com/joho/godotenv"
)

type zhk struct {
	Images      []string
	Name        string
	Class       string
	Description string
	Developer   string
	Phone       string
	Address     string
}

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("⚠️ .env файл не найден, переменные будут браться из окружения")
	}
}

func main() {
	ctx := context.Background()

	dbpool, err := database.NewPool(ctx)
	if err != nil {
		log.Fatalf("не удалось подключиться к базе данных: %v", err)
	}
	defer dbpool.Close()

	zhk := zhk{
		Name:        "ЖК Цветочные Поляны",
		Class:       "Комфорт",
		Description: "Есть такие места, куда всегда хочется возвращаться и где с первого мгновения накрывает мысль, от которой становится тепло внутри: «Я дома». Таким местом силы для вас станет квартал «Цветочные Поляны» в Новой Москве. Гуляйте с семьей по собственному лесопарку или устраивайте с детьми веселые игры в плейхабе. Заводите новых друзей среди соседей, ведь многие квартиры уже нашли своих хозяев. Не правда ли, приятнее выпить капучино в кафе на первом этаже дома или провести тренировку во дворе в компании близких по духу людей?",
		Developer:   "ООО 'Самолет'",
		Phone:       "+7(123)-456-78-90",
		Images: []string{
			"cvetochnye-polyany-moskva-jk-2390440080-6.jpg",
			"cvetochnye-polyany-moskva-jk-2390440116-6.jpg",
			"cvetochnye-polyany-moskva-jk-2390440137-6.jpg",
			"cvetochnye-polyany-moskva-jk-2390440270-6.jpg",
			"cvetochnye-polyany-moskva-jk-2390440482-6.jpg",
		},
		Address: "Москва, поселение Филимонковское, деревня Староселье",
	}

	var classID int
	err = dbpool.QueryRow(ctx, `
	SELECT id FROM kvartirum.HousingComplexClass
	WHERE name = $1;
	`, zhk.Class).Scan(&classID)

	if err != nil {
		log.Fatalf("ошибка при получении ID класса: %v", err)
	}

	var zhkID int
	err = dbpool.QueryRow(ctx, `
	INSERT INTO kvartirum.HousingComplex (class_id, name, description,
		developer, phone_number, address)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id;
	`, classID, zhk.Name, zhk.Description, zhk.Developer, zhk.Phone, zhk.Address).Scan(&zhkID)

	if err != nil {
		log.Fatalf("ошибка при вставке ЖК: %v", err)
	}

	for _, img := range zhk.Images {
		var imageID int
		err = dbpool.QueryRow(ctx, `
		INSERT INTO kvartirum.Image (uuid)
		VALUES ($1)
		RETURNING id;
		`, img).Scan(&imageID)

		if err != nil {
			log.Fatalf("ошибка при вставке image: %v", err)
		}

		dbpool.Exec(ctx, `
		INSERT INTO kvartirum.HousingComplexImages (housing_complex_id, image_id)
		VALUES ($1, $2);
		`, zhkID, imageID)
	}
}
