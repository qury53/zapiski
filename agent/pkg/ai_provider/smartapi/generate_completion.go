package ai_provider_smartapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/qury53/GC_learn/agent/internal/domain"
	core_errors "github.com/qury53/GC_learn/agent/internal/errors"
)

const (
	baseGenerateCompletionURL = "https://api.smartapi.pro/v1/chat/completions"
)

var testMsg = `ЛЕКЦИЯ № 5 • ТЕМА 3: ЭЛЕМЕНТЫ АНАЛИТИЧЕСКОЙ ГЕОМЕТРИИ, КРИВЫХ ВТОРОГО ПОРЯДКА И ТЕОРИИ
ПОВЕРХНОСТЕЙ
Прямая на плоскости
краткая выжимка для изучения: уравнение линии, виды уравнений прямой, взаимное расположение прямых, расстояние от
точки до прямой
Нумерация формул (1)–(7) и порядок разделов соответствуют исходной лекции; полужирным выделены определения и клю-
чевые условия.
1. Уравнение линии на плоскости
Линией L на плоскости называется множество точек, обладающих некоторым только им присущим геометриче-
ским свойством. Классический пример — окружность радиуса R с центром P0(x0; y0): это множество всех точек
плоскости, удалённых от центра ровно на расстояние R. Основная идея аналитической геометрии состоит в том,
чтобы описывать такие множества алгебраически — уравнением, связывающим координаты произвольной точки
множества.
Определение. Уравнением линии L на плоскости Oxy называется уравнение F (x; y) = 0, которому удовлетво-
ряют координаты любой точки данной линии и только они. Переменные x и y называются текущими координата-
ми. Аналогично определяются уравнение линии в полярной системе координат F (ρ; φ) = 0 и параметрические
уравнения x = x(t), y = y(t), где параметр t определяет положение точки на линии.
Четыре правила позволяют решать простейшие задачи о расположении линии на плоскости.`

func (p *SmartAPIProvider) GenerateCompletion(
	ctx context.Context,
	request domain.ChatRequest,
) (domain.ChatResponse, error) {
	dtoReq := domainReqToDto(request)
	// dtoReq.Messages = append(dtoReq.Messages, MessageDto{
	// 	Role:    "user",
	// 	Content: testMsg,
	// })

	p.log.Debug(
		"starting completion generation",
		slog.Any("request", dtoReq),
	)

	payload, err := json.Marshal(dtoReq)

	if err != nil {
		return domain.ChatResponse{}, fmt.Errorf("marshaling failed: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		baseGenerateCompletionURL,
		bytes.NewReader(payload),
	)

	if err != nil {
		return domain.ChatResponse{}, fmt.Errorf(
			"requred initialization failed: %w",
			err,
		)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.secretKey)

	response, err := p.client.Do(req)
	if err != nil {
		return domain.ChatResponse{}, fmt.Errorf(
			"failed to do requrest: %w",
			err,
		)
	}
	p.log.Debug(response.Status)
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	p.log.Info(string(body))
	if err != nil {
		return domain.ChatResponse{}, fmt.Errorf("read body error : %w", err)
	}

	if err := caseStatus(response.StatusCode); err != nil {
		return domain.ChatResponse{}, err
	}

	var chatResponseDto ChatResponseDto
	if err := json.Unmarshal(body, &chatResponseDto); err != nil {
		return domain.ChatResponse{}, fmt.Errorf("failed to unmarshal: %w", err)
	}

	return dtoResponseToDomain(chatResponseDto), nil
}

func caseStatus(statusCode int) error {
	switch statusCode {
	case http.StatusOK:
		return nil
	case http.StatusPaymentRequired:
		return core_errors.InsufficientBalance
	case http.StatusTooManyRequests:
		return core_errors.ErrTooMantRequest
	case http.StatusInternalServerError:
		return core_errors.ErrInternalServerError
	default:
		return core_errors.ErrUndefinedError
	}
}
