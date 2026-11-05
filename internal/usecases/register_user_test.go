package usecases

import (
	"errors"
	"testing"

	"github.com/slashispy/quicktask/internal/models"
	"go.mongodb.org/mongo-driver/mongo"
)

// fakeUserRepo is an in-memory UserRepository for tests.
type fakeUserRepo struct {
	saved      []models.User
	existing   *models.User
	findErr    error
	saveErr    error
	findByMail map[string]*models.User
}

func (f *fakeUserRepo) Save(user models.User) error {
	if f.saveErr != nil {
		return f.saveErr
	}
	f.saved = append(f.saved, user)
	return nil
}

func (f *fakeUserRepo) FindByEmail(email string) (*models.User, error) {
	if f.findErr != nil {
		return nil, f.findErr
	}
	if f.findByMail != nil {
		return f.findByMail[email], nil
	}
	return f.existing, nil
}

// fakeEncrypter returns a deterministic, recognizable hash.
type fakeEncrypter struct {
	err error
}

func (f *fakeEncrypter) Encrypt(password string) (string, error) {
	if f.err != nil {
		return "", f.err
	}
	return "hashed:" + password, nil
}

func newDTO() models.UserDTO {
	return models.UserDTO{Name: "Ada", Email: "ada@example.com", Password: "supersecret"}
}

func TestRegisterUser_Success(t *testing.T) {
	repo := &fakeUserRepo{}
	uc := &UserUseCase{UserRepo: repo, Encrypter: &fakeEncrypter{}}

	if err := uc.RegisterUser(newDTO()); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(repo.saved) != 1 {
		t.Fatalf("expected 1 saved user, got %d", len(repo.saved))
	}
	got := repo.saved[0]
	if got.Email != "ada@example.com" || got.Name != "Ada" {
		t.Errorf("unexpected saved user: %+v", got)
	}
	if got.Password != "hashed:supersecret" {
		t.Errorf("password was not encrypted before saving: %q", got.Password)
	}
}

func TestRegisterUser_AlreadyExists(t *testing.T) {
	repo := &fakeUserRepo{existing: &models.User{Email: "ada@example.com"}}
	uc := &UserUseCase{UserRepo: repo, Encrypter: &fakeEncrypter{}}

	err := uc.RegisterUser(newDTO())
	if err == nil {
		t.Fatal("expected error for existing user, got nil")
	}
	if len(repo.saved) != 0 {
		t.Errorf("should not save when user already exists, saved=%d", len(repo.saved))
	}
}

func TestRegisterUser_NoDocumentsIsNotAnError(t *testing.T) {
	// mongo.ErrNoDocuments must be treated as "user not found", not a failure.
	repo := &fakeUserRepo{findErr: mongo.ErrNoDocuments}
	uc := &UserUseCase{UserRepo: repo, Encrypter: &fakeEncrypter{}}

	if err := uc.RegisterUser(newDTO()); err != nil {
		t.Fatalf("ErrNoDocuments should be ignored, got %v", err)
	}
	if len(repo.saved) != 1 {
		t.Errorf("expected user to be saved, saved=%d", len(repo.saved))
	}
}

func TestRegisterUser_FindError(t *testing.T) {
	repo := &fakeUserRepo{findErr: errors.New("db down")}
	uc := &UserUseCase{UserRepo: repo, Encrypter: &fakeEncrypter{}}

	if err := uc.RegisterUser(newDTO()); err == nil {
		t.Fatal("expected error when FindByEmail fails, got nil")
	}
}

func TestRegisterUser_EncryptError(t *testing.T) {
	repo := &fakeUserRepo{}
	uc := &UserUseCase{UserRepo: repo, Encrypter: &fakeEncrypter{err: errors.New("boom")}}

	if err := uc.RegisterUser(newDTO()); err == nil {
		t.Fatal("expected error when Encrypt fails, got nil")
	}
	if len(repo.saved) != 0 {
		t.Errorf("should not save when encryption fails, saved=%d", len(repo.saved))
	}
}

func TestRegisterUser_SaveError(t *testing.T) {
	repo := &fakeUserRepo{saveErr: errors.New("write failed")}
	uc := &UserUseCase{UserRepo: repo, Encrypter: &fakeEncrypter{}}

	if err := uc.RegisterUser(newDTO()); err == nil {
		t.Fatal("expected error when Save fails, got nil")
	}
}
