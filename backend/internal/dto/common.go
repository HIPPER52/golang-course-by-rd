package dto

type GetDialogMessagesDTO struct {
	RoomID string `params:"room_id" validate:"required"`
}
