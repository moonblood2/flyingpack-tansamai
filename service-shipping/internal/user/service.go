package user

import (
	"errors"
	"github.com/go-ozzo/ozzo-validation/is"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/jna-distribution/service-shipping/internal"
	"github.com/jna-distribution/service-shipping/internal/contact"
	"github.com/jna-distribution/service-shipping/internal/entity"
	"github.com/jna-distribution/service-shipping/pkg/password"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

// Service interface
type Service interface {
	Login(input LoginInput) (LoginOutput, error)
	Register(input RegisterInput) (RegisterOutput, error)
	GetAll() ([]entity.User, error)
	GetById(id string) (entity.User, error)
	Edit(input EditUserInput) (entity.User, error)
	Remove(id string) error
}

// service implement Service interface.
type service struct {
	userRepository    Repository
	contactRepository contact.Repository
}

// NewService create new service.
func NewService(userRepository Repository, contactRepository contact.Repository) *service {
	return &service{
		userRepository:    userRepository,
		contactRepository: contactRepository,
	}
}

// LoginInput use when want to call Login() and pass input to it.
type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginOutput struct {
	User    entity.User
	Contact entity.Contact
}

// Validate method of LoginInput, for implement Validatable interface of ozzo-validation.
func (i LoginInput) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.Email, validation.Required),
		validation.Field(&i.Password, validation.Required),
	)
}

// Login if success return userId, nil, if not return 0, error.
func (s *service) Login(input LoginInput) (LoginOutput, error) {
	//Constant for error message.
	emailPassErr := "e-mail or password that you submitted not match with any account."
	//Validate
	if err := input.Validate(); err != nil {
		return LoginOutput{}, internal.ErrInvalidInput{InternalError: err}
	}
	exist, err := s.userRepository.DoesEmailExist(input.Email)
	if err != nil {
		return LoginOutput{}, err
	}
	//If e-mail doesn't exists, go ahead.
	if !exist {
		return LoginOutput{}, internal.ErrInvalidInput{Details: emailPassErr}
	}
	output := LoginOutput{}
	//Find user
	u, err := s.userRepository.FindByEmail(input.Email)
	if err != nil {
		return LoginOutput{}, err
	}
	//Compare user password with input password.
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(input.Password))
	if err != nil {
		return LoginOutput{}, internal.ErrInvalidInput{InternalError: err, Details: emailPassErr}
	}
	output.User = u
	//If user is contact.
	if u.Role == entity.RoleShop {
		sh, err := s.contactRepository.FindByUserId(u.Id)
		if err != nil {
			return LoginOutput{}, err
		}
		output.Contact = sh
	}

	return output, nil
}

// RegisterInput use when want to call Register() and pass input to it.
type RegisterInput struct {
	Name            string         `json:"name"`
	Email           string         `json:"email"`
	Role            entity.Role    `json:"role"`
	Password        string         `json:"password"`
	ConfirmPassword string         `json:"confirm_password"`
	Contact         entity.Contact `json:"contact"`
}

type RegisterOutput struct {
	User    entity.User
	Contact entity.Contact
}

func comparePassword(password string) validation.RuleFunc {
	return func(value interface{}) error {
		confirmPass, _ := value.(string)
		if confirmPass != password {
			return errors.New("not match with password. ")
		}
		return nil
	}
}

// Validate method of RegisterInput, for implement Validatable interface of ozzo-validation.
func (input RegisterInput) Validate() error {
	fieldRules := []*validation.FieldRules{
		validation.Field(&input.Name, validation.Required, validation.Length(6, 64)),
		validation.Field(&input.Email, validation.Required, is.Email),
		validation.Field(&input.Role,
			validation.Required,
			validation.In(entity.RoleAdmin, entity.RoleShop, entity.RoleAgentNetworkMember, entity.RoleAccounting),
		),
		validation.Field(&input.Password, validation.Required, validation.Length(10, 0)),
		validation.Field(&input.ConfirmPassword, validation.Required, validation.Length(10, 0), validation.By(comparePassword(input.Password))),
	}
	//Admin not validate contact.
	if input.Role == entity.RoleShop {
		fieldRules = append(fieldRules, validation.Field(&input.Contact))
		return validation.ValidateStruct(&input, fieldRules...)
	}
	return validation.ValidateStruct(&input, fieldRules...)
}

// Register method for user register.
func (s *service) Register(input RegisterInput) (RegisterOutput, error) {
	//Validate input.
	if err := input.Validate(); err != nil {
		return RegisterOutput{}, internal.ErrInvalidInput{Details: err}
	}
	//Check used email.
	exist, err := s.userRepository.DoesEmailExist(input.Email)
	if err != nil {
		return RegisterOutput{}, err
	}
	// If input e-mail already exists.
	if exist {
		return RegisterOutput{}, internal.ErrInvalidInput{Details: struct {
			Email string `json:"email"`
		}{"already used."}}
	}
	if input.Role == entity.RoleShop {
		//Check used phone number of contact.
		exist, err = s.contactRepository.DoesPhoneNumberExist(input.Contact.PhoneNumber)
		if err != nil {
			return RegisterOutput{}, err
		}
		if exist {
			return RegisterOutput{}, internal.ErrInvalidInput{Details: struct {
				PhoneNumber string `json:"phone_number"`
			}{"already used."}}
		}
	}
	//Generate userId
	userId := uuid.NewV4().String()
	//Hash password with Bcrypt.
	hash, err := password.BcryptHash(input.Password)
	if err != nil {
		return RegisterOutput{}, err
	}
	out := RegisterOutput{}
	u := entity.User{
		Id:       userId,
		Email:    input.Email,
		Name:     input.Name,
		Role:     input.Role,
		Password: hash,
	}
	//Check if role is member, let's create contact with user.
	if input.Role == entity.RoleShop {
		//Create new contact.
		u, sh, err := s.contactRepository.CreateWithUser(u, input.Contact)
		if err != nil {
			return RegisterOutput{}, err
		}
		out.User = u
		out.Contact = sh
	} else {
		u, err = s.userRepository.Create(u)
		if err != nil {
			return RegisterOutput{}, err
		}
		out.User = u
	}
	return out, nil
}

// GetAll method for get all users.
func (s *service) GetAll() ([]entity.User, error) {
	return s.userRepository.FindAll()
}

// GetById method for get user by id.
func (s *service) GetById(id string) (entity.User, error) {
	return s.userRepository.FindById(id)
}

// EditUserInput use when want to call Edit() and pass input to it.
type EditUserInput struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

// Validate method of EditUserInput, for implement Validatable interface of ozzo-validation.
func (input EditUserInput) Validate() error {
	//TODO 2, Regular expression for validate Password.
	return validation.ValidateStruct(&input,
		validation.Field(&input.Id, validation.Required),
		validation.Field(&input.Name, validation.When(input.Name != "", validation.Length(3, 64))),
		validation.Field(&input.Password, validation.When(input.Password != ""), validation.Length(10, 0)),
	)
}

// Edit user each filed. Not allowed to edit e-mail.
func (s *service) Edit(i EditUserInput) (entity.User, error) {
	//Validate
	if err := i.Validate(); err != nil {
		return entity.User{}, internal.ErrInvalidInput{InternalError: err}
	}
	//Find user by id.
	u, err := s.userRepository.FindById(i.Id)
	if err != nil {
		return entity.User{}, err
	}
	//Check empty Name, Role and password, if not empty use data i(input) to update.
	if len(i.Name) != 0 {
		u.Name = i.Name
	}
	if len(i.Password) != 0 {
		hash, err := password.BcryptHash(i.Password)
		if err != nil {
			return entity.User{}, err
		}
		u.Password = hash
	}
	//Call repository for update.
	return s.userRepository.Update(u)
}

// Remove user by id.
func (s *service) Remove(id string) error {
	//Check if id that received, user does exist.
	if _, err := s.userRepository.FindById(id); err != nil {
		return err
	}
	//Call repository for delete.
	if err := s.userRepository.Delete(id); err != nil {
		return err
	}
	return nil
}
