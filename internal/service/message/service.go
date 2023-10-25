package message

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/go-redis/redis"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"

	"tgtrello/config"
	"tgtrello/internal/model"
	"tgtrello/internal/pkg/crypto"
	"tgtrello/internal/pkg/utils"
	rdb "tgtrello/internal/redis"
	"tgtrello/internal/repository"
)

type Service struct {
	logger *zap.Logger
	texts  map[string]string
	bot    *tgbotapi.BotAPI
	rdb    *redis.Client
	repo   *repository.PGRepository
}

func NewMessageService(log *zap.Logger, rdb *redis.Client, repo *repository.PGRepository, bot *tgbotapi.BotAPI, texts map[string]string) *Service {
	return &Service{
		logger: log,
		bot:    bot,
		rdb:    rdb,
		repo:   repo,
		texts:  texts,
	}
}

func (m *Service) SignUp(s *model.Situation) error {
	userLogin, err := m.repo.CheckUserRegister(s.User.ID)
	if err != nil {
		return fmt.Errorf("check user: %w", err)
	}
	if userLogin != "" {
		err = m.SendMsgToUser(s.User.ID, utils.GetFormatText(m.texts, "already_registered"))
		if err != nil {
			return err
		}

		return nil
	}

	err = m.SendMsgToUser(s.User.ID, utils.GetFormatText(m.texts, "send_login"))
	if err != nil {
		return err
	}

	rdb.SetPath(m.logger, m.rdb, s.User.ID, "/login")

	return nil

}

func (m *Service) Login(s *model.Situation) error {
	exists, err := m.repo.CheckLogin(s.Message.Text)
	if err != nil {
		return err
	}

	if exists {
		err := m.SendMsgToUser(s.User.ID, utils.GetFormatText(m.texts, "login_exists"))
		if err != nil {
			return err
		}

		return nil
	}

	rdb.SetLogin(m.logger, m.rdb, s.User.ID, s.Message.Text)
	rdb.SetPath(m.logger, m.rdb, s.User.ID, "/password")

	err = m.SendMsgToUser(s.User.ID, utils.GetFormatText(m.texts, "send_password"))
	if err != nil {
		return err
	}

	return nil
}

func (m *Service) Password(s *model.Situation) error {
	login := rdb.GetLogin(m.logger, m.rdb, s.User.ID)
	if strings.Contains(login, rdb.EmptyLogin) {
		err := m.SendMsgToUser(s.User.ID, utils.GetFormatText(m.texts, "some_wrong"))
		if err != nil {
			return err
		}

		return nil
	}

	password, err := crypto.HashPassword(s.Message.Text)
	if err != nil {
		return err
	}

	user := &model.User{
		ID:         s.User.ID,
		Login:      login,
		Password:   password,
		TgName:     s.Message.Chat.FirstName,
		TgUsername: s.Message.Chat.UserName,
	}

	err = m.repo.AddNewUser(user)
	if err != nil {
		return err
	}

	return m.SendMsgToUser(s.User.ID, utils.GetFormatText(m.texts, "registration_successful"))
}

func (m *Service) Unrecognized(s *model.Situation) error {
	return m.SendMsgToUser(s.User.ID, utils.GetFormatText(m.texts, "unrecognized"))
}

func (m *Service) Start(s *model.Situation) error {
	rdb.SetPath(m.logger, m.rdb, s.User.ID, "main")
	userLogin, err := m.repo.CheckUserRegister(s.User.ID)
	if err != nil {
		return fmt.Errorf("check user: %w", err)
	}

	if userLogin == "" {
		err := m.SendMsgToUser(s.User.ID, utils.GetFormatText(m.texts, "not_registered"))
		if err != nil {
			return err
		}

		return nil
	}

	markUp := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(utils.GetFormatText(m.texts, "check_tasks")),
			tgbotapi.NewKeyboardButton(utils.GetFormatText(m.texts, "create_task"))),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(utils.GetFormatText(m.texts, "team"))))

	return m.SendMsgToUserWithMarkUp(s.User.ID, utils.GetFormatText(m.texts, "choose"), markUp)
}

func (m *Service) CreateTeam(s *model.Situation) error {
	teamId, err := m.repo.CheckTeam(s.User.ID)
	if err != nil {
		return err
	}

	if teamId != 0 {
		err = m.SendMsgToUser(s.User.ID, utils.GetFormatText(m.texts, "team_exists"))
		if err != nil {
			return err
		}

		return nil
	}

	rdb.SetPath(m.logger, m.rdb, s.User.ID, "/team_created")
	err = m.SendMsgToUser(s.User.ID, utils.GetFormatText(m.texts, "team_name"))
	if err != nil {
		return err
	}

	return nil
}

func (m *Service) TeamCreated(s *model.Situation) error {
	rdb.SetPath(m.logger, m.rdb, s.User.ID, "created")
	err := m.repo.CreateTeam(s.User.ID, s.Message.Text)
	if err != nil {
		return err
	}

	return m.SendMsgToUser(s.User.ID, utils.GetFormatText(m.texts, "team_created_successfully"))
}

func (m *Service) AddUserTeam(s *model.Situation) error {
	teamName, err := m.repo.AddUserToTeam(s.TeamID, s.User.ID)
	if err != nil {
		return err
	}

	return m.SendMsgToUser(s.User.ID, utils.GetFormatText(m.texts, "you_added_to_team", teamName))
}

func (m *Service) AddUser(s *model.Situation) error {
	teamId, err := m.repo.CheckTeam(s.User.ID)
	if err != nil {
		return err
	}

	link := config.C.BotLink + "?start=new_team_user_" + strconv.Itoa(teamId)

	return m.SendMsgToUser(s.User.ID, utils.GetFormatText(m.texts, "send_link", link))
}

func (m *Service) YourTeam(s *model.Situation) error {
	rdb.SetPath(m.logger, m.rdb, s.User.ID, "your_team")
	teamId, err := m.repo.CheckTeam(s.User.ID)
	if err != nil {
		return err
	}

	if teamId == 0 {
		err = m.SendMsgToUser(s.User.ID, utils.GetFormatText(m.texts, "team_need_create"))
		if err != nil {
			return err
		}

		return nil
	}

	team, err := m.repo.YourTeam(teamId)
	if err != nil {
		return err
	}

	var text string
	for i, user := range team.Users {
		text = strconv.Itoa(i) + ". " + user.Login + "\n"
	}

	markUp := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(utils.GetFormatText(m.texts, "add_user")),
			tgbotapi.NewKeyboardButton(utils.GetFormatText(m.texts, "delete_user"))),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(utils.GetFormatText(m.texts, "exit_team"))))
	return m.SendMsgToUserWithMarkUp(s.User.ID, utils.GetFormatText(m.texts, "team_info", team.Name, text), markUp)
}

func (m *Service) Team(s *model.Situation) error {
	markUp := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(utils.GetFormatText(m.texts, "create_team")),
			tgbotapi.NewKeyboardButton(utils.GetFormatText(m.texts, "your_team"))))

	return m.SendMsgToUserWithMarkUp(s.User.ID, utils.GetFormatText(m.texts, "choose"), markUp)
}

func (m *Service) SendMsgToUser(userID int64, text string) error {
	msg := &tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID: userID,
		},
		Text: text,
	}

	_, err := m.bot.Send(msg)
	if err != nil {
		return fmt.Errorf("send msg to user: %w", err)
	}

	return nil
}

func (m *Service) SendMsgToUserWithMarkUp(userID int64, text string, markUp tgbotapi.ReplyKeyboardMarkup) error {
	msg := &tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID:      userID,
			ReplyMarkup: markUp,
		},
		Text:      text,
		ParseMode: "HTML",
	}

	_, err := m.bot.Send(msg)
	if err != nil {
		return fmt.Errorf("send msg to user: %w", err)
	}

	return nil
}
