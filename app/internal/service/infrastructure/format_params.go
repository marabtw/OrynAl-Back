package infrastructure

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alibekabdrakhman1/orynal/internal/model"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"strconv"
	"strings"
	"time"
)

type FormatParams struct{}

func NewFormatParams() *FormatParams { return &FormatParams{} }

func (obj *FormatParams) RestaurantsSearchFormatting(params *model.Params, ctx echo.Context) (*model.Params, error) {
	err := obj.SearchFormat(params, ctx)
	if err != nil {
		return nil, err
	}

	err = obj.LimitFormat(params, ctx)
	if err != nil {
		return nil, err
	}

	err = obj.PageIndexFormat(params, ctx)
	if err != nil {
		return nil, err
	}

	return params, nil
}

func (obj *FormatParams) ReviewsSearchFormatting(params *model.Params, ctx echo.Context) (*model.Params, error) {
	err := obj.LimitFormat(params, ctx)
	if err != nil {
		return nil, err
	}

	err = obj.PageIndexFormat(params, ctx)
	if err != nil {
		return nil, err
	}

	return params, nil
}

func (obj *FormatParams) TablesSearchFormatting(params *model.Params, ctx echo.Context) (*model.Params, error) {
	err := obj.OrderFormat(params, ctx, model.TablesOrderKeyList)
	if err != nil {
		return nil, err
	}

	err = obj.OrderVectorFormat(params, ctx)
	if err != nil {
		return nil, err
	}

	err = obj.LimitFormat(params, ctx)
	if err != nil {
		return nil, err
	}

	err = obj.SearchFormat(params, ctx)
	if err != nil {
		return nil, err
	}

	err = obj.PageIndexFormat(params, ctx)
	if err != nil {
		return nil, err
	}

	err = obj.DateFormat(params, ctx)
	if err != nil {
		return nil, err
	}

	return params, nil
}

func (obj *FormatParams) UserSearchFormatting(params *model.Params, ctx echo.Context) (*model.Params, error) {
	err := obj.SearchFormat(params, ctx)
	if err != nil {
		return nil, err
	}

	err = obj.OrderFormat(params, ctx, model.UsersOrderKeyList)
	if err != nil {
		return nil, err
	}

	err = obj.OrderVectorFormat(params, ctx)
	if err != nil {
		return nil, err
	}

	err = obj.LimitFormat(params, ctx)
	if err != nil {
		return nil, err
	}

	err = obj.PageIndexFormat(params, ctx)
	if err != nil {
		return nil, err
	}

	return params, nil
}

func (obj *FormatParams) OrderSearchFormatting(params *model.Params, ctx echo.Context) (*model.Params, error) {
	err := obj.OrderFormat(params, ctx, model.UsersOrderKeyList)
	if err != nil {
		return nil, err
	}

	err = obj.OrderVectorFormat(params, ctx)
	if err != nil {
		return nil, err
	}

	err = obj.LimitFormat(params, ctx)
	if err != nil {
		return nil, err
	}

	err = obj.PageIndexFormat(params, ctx)
	if err != nil {
		return nil, err
	}

	return params, nil
}

func (obj *FormatParams) MenuSearchFormatting(params *model.Params, ctx echo.Context) (*model.Params, error) {
	err := obj.SearchFormat(params, ctx)
	if err != nil {
		return nil, err
	}

	err = obj.LimitFormat(params, ctx)
	if err != nil {
		return nil, err
	}

	err = obj.PageIndexFormat(params, ctx)
	if err != nil {
		return nil, err
	}

	return params, nil
}

func (obj *FormatParams) OrderVectorFormat(paramsModel *model.Params, ctx echo.Context) error {
	orderVector := ctx.QueryParam("order_vector")

	if orderVector == "" {
		paramsModel.SortVector = nil
	} else {
		if lo.Contains([]string{"asc", "desc"}, orderVector) {
			paramsModel.SortVector = orderVector
		} else {
			return errors.New("your order_vector param is not accepted")
		}
	}

	return nil
}

func (obj *FormatParams) FilterFormat(paramsModel *model.Params, ctx echo.Context, keyFilterList map[string]string) error {
	filter := ctx.QueryParam("filter")

	if filter != "" {
		var unmarshalledFilter interface{}
		err := json.Unmarshal([]byte(filter), &unmarshalledFilter)
		if err != nil {
			return err
		}

		paramsModel.Filter = unmarshalledFilter.(map[string]interface{})

		for name, _ := range paramsModel.Filter {
			if !lo.Contains(lo.Values(keyFilterList), name) {
				return errors.New(fmt.Sprintf("filter value not found %s", name))
			}
		}
	} else {
		return errors.New("filter param not found")
	}

	return nil
}

func (obj *FormatParams) OrderFormat(paramsModel *model.Params, ctx echo.Context, keyFilterList []string) error {
	order := ctx.QueryParam("order")
	if order != "" {
		var unmarshalledOrder interface{}
		err := json.Unmarshal([]byte(order), &unmarshalledOrder)
		if err != nil {
			return err
		}

		var newList []string

		for _, item := range unmarshalledOrder.([]interface{}) {
			if !lo.Contains(keyFilterList, item.(string)) {
				return errors.New(fmt.Sprintf("%v param is not accepted", item))
			}

			newList = append(newList, item.(string))
		}
		paramsModel.Order = strings.Join(newList, ",")
	}

	return nil
}

func (obj *FormatParams) LimitFormat(paramsModel *model.Params, ctx echo.Context) error {
	limit := ctx.QueryParam("limit")

	if limit == "" {
		paramsModel.Limit = 20
		return nil
	}

	convertedInt, err := strconv.Atoi(limit)
	if err != nil {
		return err
	}

	switch {
	case convertedInt > 20 && convertedInt < 0:
		paramsModel.Limit = 20
	default:
		paramsModel.Limit = convertedInt
	}

	return nil
}

func (obj *FormatParams) PageIndexFormat(paramsModel *model.Params, ctx echo.Context) error {
	page := ctx.QueryParam("page")

	if page != "" {
		intQueryLimit, err := strconv.Atoi(page)
		if err != nil {
			return err
		}

		if intQueryLimit <= 1 {
			paramsModel.Offset = 0
			paramsModel.PageIndex = 1
		} else {
			paramsModel.Offset = (intQueryLimit - 1) * paramsModel.Limit
			paramsModel.PageIndex = intQueryLimit
		}
	} else {
		paramsModel.Offset = 0
		paramsModel.PageIndex = 1
	}

	return nil
}

func (obj *FormatParams) SearchFormat(paramsModel *model.Params, ctx echo.Context) error {
	q := ctx.QueryParam("q")

	if q != "" {
		paramsModel.Query = q
	}

	return nil
}

func (obj *FormatParams) DateFormat(paramsModel *model.Params, ctx echo.Context) error {
	date := ctx.QueryParam("date")

	if date != "" {
		layout := "2006-01-02T15:04:05"
		t, err := time.Parse(layout, date)

		if err != nil {
			return fmt.Errorf("error parsing date string: %w", err)
		}

		paramsModel.Date = &t
	}

	return nil
}
