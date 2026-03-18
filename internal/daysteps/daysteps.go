package daysteps

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	// Делим строку на слайс строк
	split := strings.Split(data, ",")
	// Проверяем длину слайса и преобразуем его значения
	if len(split) == 2 {
		steps, err := strconv.Atoi(split[0]) // Преобразуем число шагов в int
		if err != nil {
			return 0, 0, fmt.Errorf("ошибка преобразования количества шагов: %v", err)
		}
		duration, err := time.ParseDuration(split[1]) // Преобразуем продолжительность прогулки в time.Duration
		if err != nil {
			return 0, 0, fmt.Errorf("ошибка преобразования продолжительности прогулки: %v", err)
		}
		// Проверяем введенные данные
		if steps <= 0 || duration <= 0 {
			return 0, 0, fmt.Errorf("количество шагов и продолжительность должно быть больше 0, введено %d, %v", steps, duration)
		}
		return steps, duration, nil
	} else {
		return 0, 0, fmt.Errorf("ошибка: введите данные с двумя значениями, введено %d", len(split))
	}
}

func DayActionInfo(data string, weight, height float64) string {
	// Получаем данные о количестве шагов и продолжительности прогулки
	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Println(err)
		return ""
	}
	if steps <= 0 {
		return ""
	}
	distance := float64(steps) * stepLength // Вычесляем пройденую дистанцию в метрах
	distanceKm := distance / mInKm          // Переводим метры в километры
	// Вычесляем количество квлорий
	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		return ""
	}
	resultString := fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, distanceKm, calories)
	return resultString
}
