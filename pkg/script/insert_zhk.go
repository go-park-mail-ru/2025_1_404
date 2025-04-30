//go:build script

package main

import (
	"context"
	"log"

	"github.com/go-park-mail-ru/2025_1_404/config"
	database "github.com/go-park-mail-ru/2025_1_404/pkg/database/postgres"
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
	Station     string
}

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("не удалось загрузить конфиг: %v", err)
	}

	ctx := context.Background()

	dbpool, err := database.NewPool(&cfg.Postgres, ctx)
	if err != nil {
		log.Fatalf("не удалось подключиться к базе данных: %v", err)
	}
	defer dbpool.Close()

	var zhks []zhk
	zhks = append(zhks, zhk{
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
		Station: "Филатов луг",
	})

	zhks = append(zhks, zhk{
		Name:        "ЖК Era(Эра)",
		Class:       "Бизнес",
		Description: "Жилой комплекс ERA — современный премиум-квартал, расположенный в Даниловском районе Москвы, на Дербеневской набережной. Проект выполнен в стиле ар-деко и включает шесть высотных башен с элегантной архитектурой. ERA сочетает удобство городской жизни и комфорт загородного отдыха. На территории комплекса создано уникальное благоустройство: английский сад, цветущие поляны, зеленый лабиринт, зоны для отдыха и современные детские площадки. Первые этажи корпусов и стилобат займут коммерческие помещения, включая магазины и кафе",

		Developer: "ООО 'TEKTA GROUP'",
		Phone:     "+7(495)-157-76-47",
		Images: []string{
			"era-moskva-jk-2346863144-10.jpg",
			"era-moskva-jk-2346863179-10.jpg",
			"era-moskva-jk-2346863332-10.jpg",
			"era-moskva-jk-2346863556-10.jpg",
			"era-moskva-jk-2346863752-10.jpg",
		},
		Address: "Москва, жилой комплекс Эра, 3",
		Station: "Технопарк",
	})

	zhks = append(zhks, zhk{
		Name:        "ЖК Vesper Tverskaya",
		Class:       "Бизнес",
		Description: "Исключительный проект на главной улице Москвы, меняющий эстетику жизни в современном мегаполисе.Компания Vesper представляет проект на Тверской — апартаменты класса de luxe с сервисом и инфраструктурой. Проект Vesper Tverskaya завершает архитектурный ансамбль Тверской, занимая последний свободный участок главной улицы города — в самом центре культурной, деловой и ночной жизни Москвы.Два здания, в которых расположены апартаменты, создают единый ансамбль.Они соединены общими этажами, где расположены ресторан, wellness с бассейном, торговая галерея и паркинг. При этом у каждого из зданий отдельный вход.",

		Developer: "ООО 'VESPER'",
		Phone:     "+7(495)-138-19-37",
		Images: []string{
			"vesper-tverskaya-moskva-jk-1826317936-10.jpg",
			"vesper-tverskaya-moskva-jk-2304796811-10.jpg",
			"vesper-tverskaya-moskva-jk-1826317738-10.jpg",
			"vesper-tverskaya-moskva-jk-1826317681-10.jpg",
			"vesper-tverskaya-moskva-jk-2304795911-10.jpg",
		},
		Address: "Москва, 1-я Тверская-Ямская улица, 2А",
		Station: "Тверская",
	})

	for _, zhk := range zhks {
		var classID int
		err = dbpool.QueryRow(ctx, `
		SELECT id FROM kvartirum.HousingComplexClass
		WHERE name = $1;
		`, zhk.Class).Scan(&classID)


		if err != nil {
			log.Fatalf("ошибка при получении ID класса: %v", err)
		}

		var metroId int
		err = dbpool.QueryRow(ctx, `
			SELECT id from kvartirum.MetroStation
			WHERE name = $1;
		`, zhk.Station).Scan(&metroId)

		if err != nil {
			log.Fatalf("Ошибка при получение id станции: %v", err)
		}

		var zhkID int
		err = dbpool.QueryRow(ctx, `
		INSERT INTO kvartirum.HousingComplex (class_id, name, description,
			developer, phone_number, address, metro_station_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id;
		`, classID, zhk.Name, zhk.Description, zhk.Developer, zhk.Phone, zhk.Address, metroId).Scan(&zhkID)

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
}
