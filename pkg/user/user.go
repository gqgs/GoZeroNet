package user

import (
	"encoding/json"
	"os"
	"path"

	"github.com/gqgs/go-zeronet/pkg/config"
)

type User interface {
	AuthAddress(addr string) string
	CertUserID(addr string) string
	SiteSettings(addr string) map[string]interface{}
}

type userManager struct {
	users map[string]*user
}

type UserManager interface {
	User() User
}

func (m *userManager) User() User {
	// For now return the first user
	for _, user := range m.users {
		return user
	}
	return nil
}

func NewUserManager() (*userManager, error) {
	userFilePath := path.Join(config.DataDir, "users.json")
	file, err := os.Open(userFilePath)
	if err != nil {
		// TODO: create if it doesn't exist
		return nil, err
	}
	defer file.Close()

	users := make(map[string]*user)
	if err := json.NewDecoder(file).Decode(&users); err != nil {
		return nil, err
	}

	return &userManager{users}, nil
}

type Cert struct {
	AuthAddress    string `json:"auth_address"`
	AuthPrivatekey string `json:"auth_privatekey"`
	AuthType       string `json:"auth_type"`
	AuthuserName   string `json:"auth_user_name"`
	CertSign       string `json:"cert_sign"`
}

type Site struct {
	AuthAddress    string                 `json:"auth_address"`
	AuthPrivatekey string                 `json:"auth_privatekey"`
	Settings       map[string]interface{} `json:"settings"`
}

type user struct {
	Certs      map[string]Cert        `json:"certs"`
	MasterSeed string                 `json:"master_seed"`
	Settings   map[string]interface{} `json:"settings"`
	Sites      map[string]Site        `json:"sites"`
}

func (u *user) AuthAddress(addr string) string {
	return u.Sites[addr].AuthAddress
}

func (u *user) CertUserID(addr string) string {
	// TODO: implement me
	return ""
}

func (u *user) SiteSettings(addr string) map[string]interface{} {
	return u.Sites[addr].Settings
}
