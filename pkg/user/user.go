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
	GlobalSettings() GlobalSettings
}

type manager struct {
	users map[string]*user
}

type Manager interface {
	User() User
}

func (m *manager) User() User {
	// For now return the first user
	for _, user := range m.users {
		return user
	}
	return nil
}

func NewUserManager() (*manager, error) {
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

	return &manager{users}, nil
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

type GlobalSettings struct {
	Theme          string `json:"theme"`
	UseSystemTheme bool   `json:"use_system_theme"`
}

type user struct {
	Certs      map[string]Cert `json:"certs"`
	MasterSeed string          `json:"master_seed"`
	Settings   GlobalSettings  `json:"settings"`
	Sites      map[string]Site `json:"sites"`
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

func (u *user) GlobalSettings() GlobalSettings {
	return u.Settings
}
