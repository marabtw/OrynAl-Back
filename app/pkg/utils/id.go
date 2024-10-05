package utils

import (
	"context"
	"errors"
	"fmt"
	"github.com/alibekabdrakhman1/orynal/internal/model"
	"strconv"
)

func ConvertIdToUint(in string) (uint, error) {
	id, err := strconv.ParseUint(in, 10, 32)
	if err != nil {
		return 0, errors.New(fmt.Sprintf("converting id to uint error: %v", err))
	}

	return uint(id), err
}

func GetIDFromContext(ctx context.Context) (uint, error) {
	userID, ok := ctx.Value(model.ContextUserIDKey).(*model.ContextUserID)
	if !ok {
		return 0, errors.New("not valid context userID")
	}

	return userID.ID, nil
}

func GetRoleFromContext(ctx context.Context) (string, error) {
	userRole, ok := ctx.Value(model.ContextUserRoleKey).(*model.ContextUserRole)
	if !ok {
		return "", errors.New("not valid context userRole")
	}

	return userRole.Role, nil
}
