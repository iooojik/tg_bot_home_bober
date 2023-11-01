package service

type BotService struct {
}

func NewBotService() *BotService {
	return nil
}

func (bs *BotService) HandleMessage() {

}

func (bs *BotService) CheckUser(userId int) (string, error) {
	return "", nil
}

func (bs *BotService) ChangeDate(date, userId int) (string, error) {
	return "", nil
}
